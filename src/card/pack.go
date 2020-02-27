package card

import (
	"math/rand"
)

// Pack represents a pack of cards with multiple suits
type Pack struct {
	Name  string            `json:"name"`
	Suits map[string][]Card `json:"suits"`
}

// CardCount returns the number of cards in the pack
func (p *Pack) CardCount() int {
	cardct := int(0)

	for n := range p.Suits {
		cardct += len(p.Suits[n])
	}

	return cardct
}

// CardsInSuit returns the number of cards in the suit (or 0 if the suit doesn't exist)
func (p *Pack) CardsInSuit(name string) int {
	if p.Suits[name] == nil {
		return 0
	}

	return len(p.Suits[name])
}

// GetCardChance returns a chance value such that the GetRandomCard funcation has a good chance of outputing a card.
func (p *Pack) GetCardChance() float64 {
	return 1.0 / float64(p.CardCount())
}

// GetRandomCard returns a random card reference by generating random numbers with r.
// The chance variable is the percent chance of any given card being chosen and should be between 0 (inclusive) and 1 (exclusive).
// The function will return NilRef if chance is below zero or no card is selected by r.
func (p *Pack) GetRandomCard(r *rand.Rand, chance float64) Ref {
	if chance <= 0 {
		return NilRef
	}

	for u, cs := range p.Suits {
		for v := range cs {
			if r.Float64() < chance {
				return Ref{P: p, Suit: u, CardID: v}
			}
		}
	}

	return NilRef
}

// GetSuitChance returns a chance value such that the GetRandomSuit funcation has a good chance of outputing a suit.
func (p *Pack) GetSuitChance() float64 {
	return 1.0 / float64(len(p.Suits))
}

// GetRandomSuit returns a reference to the first card in a suit by generating random numbers with r.
// The chance variable is the percent chance of any given suit being chosen and should be between 0 (inclusive) and 1 (exclusive).
// The function will return NilRef if chance is below zero or no suit is selected by r.
func (p *Pack) GetRandomSuit(r *rand.Rand, chance float64) Ref {
	if chance <= 0 {
		return NilRef
	}

	for u := range p.Suits {
		if r.Float64() < chance {
			return Ref{P: p, Suit: u, CardID: 0}
		}
	}

	return NilRef
}

// GetSuits returns a slice of all the suit names in the pack
func (p *Pack) GetSuits() []string {
	out := []string{}

	for v := range p.Suits {
		out = append(out, v)
	}

	return out
}

// GetCardInSuitChance returns a chance value such that the GetRandomCardInSuit funcation has a good chance of outputing a card.
// Returns -1 if the suit doesn't exist
func (p *Pack) GetCardInSuitChance(name string) float64 {
	return 1.0 / float64(p.CardsInSuit(name))
}

// GetRandomCardInSuit returns a random card from a suit with the specified name.
// The chance variable is the percent chance of any given suit being chosen and should be between 0 (inclusive) and 1 (exclusive).
// The function will return NilRef if chance is below zero, the specified suit doesn't exist, or no card is selected by r.
func (p *Pack) GetRandomCardInSuit(r *rand.Rand, name string, chance float64) Ref {
	if chance <= 0 {
		return NilRef
	}

	if p.Suits[name] == nil {
		return NilRef
	}

	for c := range p.Suits[name] {
		if r.Float64() < chance {
			return Ref{P: p, Suit: name, CardID: c}
		}
	}

	return NilRef
}

// Validate checks the pack to see weather there are any suits in it.
// It also makes sure that each suit has at least one card.
func (p *Pack) Validate() bool {
	for len(p.Suits) == 0 {
		return false
	}

	for _, s := range p.Suits {
		if len(s) == 0 {
			return false
		}
	}

	return true
}
