package card

import (
	"math/rand"
)

// Deck represents a deck of cards composed from one or more packs.
// A Desk should be created through the NewDeck function.
// Decks abstract a set of packs into a single pack with some extra features, but you can still ask for a pack based on it's name.
// If you want to cherry pick cards or add/remove cards in game, see RefDeck
type Deck struct {
	Packs [](*Pack)

	r *rand.Rand

	duplicates uint
	suits      []string
	cardCt     int
}

// NewDeck creates a new Deck struct and initializes it's random.
// packs sets the packs that are included in the deck.
// dups sets how many of the same card can be in play at once (-1 for infinite).
// The Rand is initialized with a source dependent on the current time.
func NewDeck(packs [](*Pack), dups uint) *Deck {
	ct := int(0)

	for i := range packs {
		ct += packs[i].CardCount()
	}

	var suits []string
	var inSuits bool
	for i := range packs {

		for s := range packs[i].Suits {
			inSuits = false

			for j := range suits {
				if s == suits[j] {
					inSuits = true
					break
				}
			}

			if inSuits {
				continue
			}

			suits = append(suits, s)
		}
	}

	return &Deck{Packs: packs, r: syncSafeRandom(), cardCt: ct, suits: suits}
}

// GetPack returns a pack based on it's name.
// Returns nil if the pack does not exist in the deck.
func (d *Deck) GetPack(name string) *Pack {
	for i := range d.Packs {
		if d.Packs[i].Name == name {
			return d.Packs[i]
		}
	}

	return nil
}

// CardCount returns the number of cards in the deck
func (d *Deck) CardCount() int {
	i := d.cardCt
	return i
}

// SuitList returns a slice containing all the unique suits in the deck.
func (d *Deck) SuitList() []string {
	s := d.suits
	return s
}

// SuitCount returns the number of unique suits in the deck.
func (d *Deck) SuitCount() int {
	return len(d.suits)
}

// CardsInSuit returns the number of cards in the suit (or 0 if the suit doesn't exist)
func (d *Deck) CardsInSuit(name string) int {
	ct := int(0)

	for i := range d.Packs {
		c := d.Packs[i].CardsInSuit(name)
		if c > 0 {
			ct += c
		}
	}

	return ct
}

// GetCardChance returns a chance value such that the GetRandomCard funcation has a good chance of outputing a card.
func (d *Deck) GetCardChance() float64 {
	return 1.0 / float64(d.cardCt)
}

// GetRandomCard returns a reference to a random card in the deck.
// The reference is garunteed to not be NilRef unless chance is zero or lower.
func (d *Deck) GetRandomCard(chance float64) Ref {
	if chance <= 0 {
		return NilRef
	}

	var c Ref
	for {
		for _, p := range d.Packs {
			c = p.GetRandomCard(d.r, chance)

			if c != NilRef {
				return c
			}
		}
	}
}

// GetCardInSuitChance returns a chance value such that the GetRandomCard funcation has a good chance of outputing a card.
// Returns -1 if no suit with name exists.
// This function is expensive and should be used sparingly.
func (d *Deck) GetCardInSuitChance(name string) float64 {
	i := d.CardsInSuit(name)
	if i < 1 {
		return -1
	}
	return 1 / float64(i)
}

// GetRandomCardInSuit gets a random card in the specified suit.
// Returns NilRef if the suit doesn't exist or chance is zero or lower.
func (d *Deck) GetRandomCardInSuit(name string, chance float64) Ref {
	if chance <= 0 || d.CardsInSuit(name) == 0 {
		return NilRef
	}

	var c Ref
	for {
		for _, p := range d.Packs {
			c = p.GetRandomCardInSuit(d.r, name, chance)

			if c != NilRef {
				return c
			}
		}
	}
}

// GetCardInPackChance returns a chance value such that the GetRandomCardInPack funcation has a good chance of outputing a card.
// Returns -1 if no pack with name exists.
// This function is expensive and should be used sparingly.
func (d *Deck) GetCardInPackChance(name string) float64 {
	p := d.GetPack(name)
	if p == nil {
		return -1
	}
	return 1 / float64(p.GetCardChance())
}

// GetRandomCardInPack gets a random card in the specified pack.
// Returns NilRef if the pack doesn't exist or chance is zero or lower.
func (d *Deck) GetRandomCardInPack(name string, chance float64) Ref {
	p := d.GetPack(name)

	if chance <= 0 || p == nil {
		return NilRef
	}

	var c Ref
	for {
		c = p.GetRandomCard(d.r, chance)

		if c != NilRef {
			return c
		}
	}
}
