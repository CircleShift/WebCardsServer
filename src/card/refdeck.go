package card

import (
	"math/rand"
	"time"
)

// RefDeck defines a deck of individual cards with no bearing on what pack each card belongs to.
// RefDeck supports removing and adding cards at runtime without damaging other decks or packs.
// RefDeck does not support functions to work on card suits out-of-the-box.
// NewRefDeck should be used to create a RefDeck.
type RefDeck struct {
	Cards  []Ref
	r      *rand.Rand
	cardCt int
}

// NewRefDeck creates a reference to a RefDeck based on card references and number of duplicates allowed in play.
func NewRefDeck(c []Ref, dups uint) *RefDeck {
	n := time.Now()
	d := time.Date(n.Year()-(n.Year()%32), time.January, 0, 0, 0, 0, 0, time.Local)
	rnd := rand.New(rand.NewSource(n.Sub(d).Nanoseconds()))

	return &RefDeck{Cards: c, r: rnd, cardCt: len(c)}
}

// CardCount returns the number of cards in the RefDeck.
func (rd *RefDeck) CardCount() int {
	i := rd.cardCt
	return i
}

// GetCardChance returns a chance value such that the GetRandomCard funcation has a good chance of outputing a card.
func (rd *RefDeck) GetCardChance() float64 {
	return 1.0 / float64(rd.cardCt)
}

// GetRandomCard returns a random card reference by generating random numbers with r.
// The chance variable is the percent chance of any given card being chosen and should be between 0 (inclusive) and 1 (exclusive).
// The function will return NilRef if chance is below zero or no card is selected by r.
func (rd *RefDeck) GetRandomCard(chance float64) Ref {
	if chance <= 0 {
		return NilRef
	}

	for c := range rd.Cards {
		if rd.r.Float64() < chance {
			return rd.Cards[c]
		}
	}

	return NilRef
}

// RemoveCard removes a card from the deck.  If it is successful, it returns true.
func (rd *RefDeck) RemoveCard(ref Ref) bool {
	for i := range rd.Cards {
		if rd.Cards[i] == ref {
			rd.RemoveCardAt(i)
			return true
		}
	}
	return false
}

// RemoveCardAt removes the card at index from the deck.  If it is successful, it returns true.
func (rd *RefDeck) RemoveCardAt(index int) bool {
	if rd.cardCt <= index || index < 0 {
		return false
	}

	rd.cardCt--
	rd.Cards[index] = rd.Cards[rd.cardCt]
	rd.Cards = rd.Cards[:rd.cardCt]

	return true
}

// AddCard adds a card to the deck.
func (rd *RefDeck) AddCard(ref Ref) {
	rd.cardCt++
	rd.Cards = append(rd.Cards, ref)
}
