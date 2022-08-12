package webcode

import (
	"sync"
	"log"
)

type Game struct {
	// lets the clean loop know that it can delete the game object
	end bool

	// chats used in game
	Chats []string
	// players in game
	Players []string
	// current player
	turn string

	// password for game ("" means no password)
	pass string

	// The game object is going to be used by multiple threads,
	// so it's important to use a mutex to prevent race conditions
	game_state sync.Mutex
}

// For modification: init code
func InitGame(options GOptions, id string, player_one string) *Game {
	var out Game

	out.Chats = []string{}
	out.Players = []string{player_one}
	out.end = false

	if options.UsePassword {
		out.pass = options.Password
	} else {
		out.pass = ""
	}
	
	if !options.Hidden {
		log.Println("Adding public game!")
		addPublic(AddGameMessage{options.Name, 0, id, options.UsePassword})
	}

	return &out
}

// What happens when the game is done
func (g *Game) EndGame() {
	g.end = true
}

// Return true if player is added to the game
// return false otherwise
func (g *Game) playerAdd(pid string, password string) bool {
	if g.end || (g.pass != password && g.pass != "") {
		return false
	}
	return true
}

// Setup player UI
// Might be called if a player re-joins a game
func (g *Game) setupUI(pid string) {
	if p := getPlayer(pid); p != nil {
		p.newDeck(NewDeckMessage{"0", DeckOptions{"stack", "one", 0, [4]float64{0.5, 0.5, 0.0, 1.0}}})
		p.newDeck(NewDeckMessage{pid, DeckOptions{"strip-hr", "one", 0, [4]float64{0.5, 1.0, 0.0, 1.0}}})
	}
}

// Once the player has been informed of game join, this is called.
// Should call setupUI and also 
// Deck/Card setup should happen here.
func (g *Game) playerJoin(pid string) {
	g.setupUI(pid)
	g.game_state.Lock()
	defer g.game_state.Unlock()
	g.Players = append(g.Players, pid)
}

// Called when player asks to leave
func (g *Game) playerLeave(pid string) {
	g.game_state.Lock()
	defer g.game_state.Unlock()
	if len(g.Players) <= 1 {
		g.EndGame()
	} else {

	}
}

// For modification: code for when a player makes a move
// Return false if move is illegal
func (g *Game) TryMove(player string) bool {
	return true
}