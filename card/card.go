// Package card provides card manipulation types in the form of Decks, Sets, Suits, and Cards
package card

// Card represents a single card that a player can hold.
// Data represents the text or images on the card, and is one to one with the server version.
type Card struct {
	Name string `json:"name"`
	Data interface{} `json:"data"`
}
