package game

type Game struct {
	ChatID string
}

// For modification: init code
func Init(options GOptions) *Game {
	return nil
}

// For modification: code for when a player makes a move
// Return false if move is illegal
func (g *Game) TryMove(player string) bool {
	return true
}