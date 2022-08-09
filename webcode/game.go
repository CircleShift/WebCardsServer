package webcode

import (
	"sync"
)

type Game struct {
	// Must provide all chats used in game
	Chats []string
	
	Players []string
	pass string
	end bool

	// The game object is going to be used by multiple threads,
	// so it's important to use a mutex to prevent race conditions
	game_state sync.Mutex
}

// For modification: init code
func InitGame(options GOptions, player_one string) *Game {
	var out Game

	out.Chats = []string{}
	out.Players = []string{player_one}
	out.end = false
	out.pass = options.Password
	


	return &out
}

// What happens when the game is done
func (g *Game) EndGame() {
	g.end = true
}

// Return true if player is added to the game
// return false otherwise
func (g *Game) playerAdd(pid string, password string) bool {
	if g.end || g.pass != password {
		return false
	}
	return true
}

// Once the player has been informed of game join, this is called.
// Deck/Card setup should happen here.
func (g *Game) playerJoin(pid string) {
	if p := getPlayer(pid); p != nil && !p.as.isClosed() {
		p.joinGame()
	}
}

// Called when player asks to leave
func (g *Game) playerLeave(pid string) {
	if p := getPlayer(pid); p != nil && !p.as.isClosed() {
		p.leaveGame()
	}
}

// For modification: code for when a player makes a move
// Return false if move is illegal
func (g *Game) TryMove(player string) bool {
	return true
}