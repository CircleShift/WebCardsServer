package card

// Ref represents a reference to a card.
type Ref struct {
	P      *Pack
	Suit   string
	CardID int
}

// GetCard returns a copy of the card that r points to.
func (r *Ref) GetCard() Card {
	return r.P.Suits[r.Suit][r.CardID]
}

// GetSuit returns a copy of the suit that r points to.
// Use sparingly as the copy required can be an expensive operation.
func (r *Ref) GetSuit() []Card {
	return r.P.Suits[r.Suit]
}

// NilRef represents a non-reference.
var NilRef = Ref{P: nil, Suit: "", CardID: 0}
