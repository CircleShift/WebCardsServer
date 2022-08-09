package card

// DeckState catalogues the cards played and will let you know of a card is in play.
// InPlay: Cards in play. These may be in hand or in discard.
// Duplicates: Number of times a ref can be duplicated before it is considered in play (-1 for infinite).
type DeckState struct {
	InHand    []Ref
	InDiscard []Ref

	Duplicates int
}

// CountHand returns the number of times ref appears in InHand.
func (ds *DeckState) CountHand(ref Ref) int {
	ct := int(0)

	for i := range ds.InHand {
		if ds.InHand[i] == ref {
			ct++
		}
	}

	return ct
}

// CountDiscarded returns the number of times ref appears in InDiscard.
func (ds *DeckState) CountDiscarded(ref Ref) int {
	ct := int(0)

	for i := range ds.InDiscard {
		if ds.InDiscard[i] == ref {
			ct++
		}
	}

	return ct
}

// CountPlayed returns the number of times ref appears in InDiscard or InHand.
func (ds *DeckState) CountPlayed(ref Ref) int {
	return ds.CountDiscarded(ref) + ds.CountHand(ref)
}

// IsInPlay returns if the card is considered in play (no more of the same ref can be added).
func (ds *DeckState) IsInPlay(ref Ref) bool {
	if ds.Duplicates == -1 {
		return false
	}

	return ds.CountPlayed(ref) > ds.Duplicates
}

// PutInHand checks if a card is in play and puts it in InHand if it isn't.
// Returns false if the card is already in play.
func (ds *DeckState) PutInHand(ref Ref) bool {
	if ds.IsInPlay(ref) {
		return false
	}

	ds.InHand = append(ds.InHand, ref)

	return true
}

// PutInDiscard checks if a card is in play and puts it in InDiscard if it isn't.
// Returns false if the card is already in play.
func (ds *DeckState) PutInDiscard(ref Ref) bool {
	if ds.IsInPlay(ref) {
		return false
	}

	ds.InDiscard = append(ds.InDiscard, ref)

	return true
}

// Function helped by https://stackoverflow.com/a/37335777

// MoveToDiscard checks if a card is in InHand and moves it to InDiscard if it is.
// Returns true if the card was in hand.
func (ds *DeckState) MoveToDiscard(ref Ref) bool {
	index := int(-1)

	for i := range ds.InHand {
		if ds.InHand[i] == ref {
			index = i
			break
		}
	}

	if index == -1 {
		return false
	}

	ds.InHand[index] = ds.InHand[len(ds.InHand)-1]
	ds.InHand = ds.InHand[:len(ds.InHand)-1]

	ds.InDiscard = append(ds.InDiscard, ref)

	return true
}

// Function helped by https://stackoverflow.com/a/37335777

// MoveToHand checks if a card is in InDiscard and moves it to InHand if it is.
// Returns true if the card was in hand.
func (ds *DeckState) MoveToHand(ref Ref) bool {
	index := int(-1)

	for i := range ds.InDiscard {
		if ds.InDiscard[i] == ref {
			index = i
			break
		}
	}

	if index == -1 {
		return false
	}

	ds.InDiscard[index] = ds.InDiscard[len(ds.InDiscard)-1]
	ds.InDiscard = ds.InDiscard[:len(ds.InDiscard)-1]

	ds.InHand = append(ds.InHand, ref)

	return true
}
