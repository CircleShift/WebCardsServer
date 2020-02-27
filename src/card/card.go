// Package card provides card manipulation types in the form of Decks, Sets, Suits, and Cards
package card

// Card represents a single card that a player can hold.
type Card struct {
	Name string `json:"name"`
	Text string `json:"text"`
}
