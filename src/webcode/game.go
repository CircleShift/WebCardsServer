package webcode


type Game struct {
	// Must provide all 
	Chats []string
	// 
	Players []string
	end bool
}

// For modification: init code
func InitGame(options GOptions, player_one string) *Game {
	return nil
}

// Return true if player is added to the game
// return false otherwise
func (g *Game) AddPlayer(pid string, password string) bool {
	if g.end {
		return false
	}
	return false
}

// For modification: code for when a player makes a move
// Return false if move is illegal
func (g *Game) TryMove(player string) bool {
	return true
}