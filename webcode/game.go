package webcode

import (
	"sync"
	"log"

	card "cshift.net/webcards/card"
)

type Game struct {
	// lets the clean loop know that it can delete the game object
	end bool
	id string

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
	out.Players = []string{}
	out.end = false
	out.id = id

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
	log.Println("Game end")
	delPublic(g.id)
	g.end = true
}

// Return true if player is added to the game
// return false otherwise
func (g *Game) playerCanJoin(pid string, password string) bool {
	if g.end || (g.pass != password && g.pass != "") {
		return false
	}
	return true
}

// Setup player UI
// Might be called if a player re-joins a game
func (g *Game) setupUI(pid string) {
	if p := getPlayer(pid); p != nil {
		p.newDeck(NewDeckMessage{"0", DeckOptions{"inf", "one", 0, [4]float64{0.05, 0.05, 0.0, 1.0}}})
		p.newDeck(NewDeckMessage{"1", DeckOptions{"stack", "one", 0, [4]float64{0.95, 0.05, 0.0, 1.0}}})
		p.newDeck(NewDeckMessage{pid, DeckOptions{"strip-hr", "one", 0, [4]float64{0.5, 0.95, 0.0, 1.0}}})
		p.newCard(NewCardMessage{"0", "0", card.Packs[0].GetCard("draw", "draw").Data})
		p.newCard(NewCardMessage{"1", "1", card.Packs[0].GetCard("blue", "0").Data})
	}
}

// Once the player has been informed of game join, this is called.
// Should call setupUI and also 
// Deck/Card setup should happen here.
func (g *Game) playerJoin(pid string) {
	g.setupUI(pid)
	g.game_state.Lock()
	defer g.game_state.Unlock()
	if p := getPlayer(pid); p != nil && g.end {
		p.leaveGame()
		return
	}
	g.Players = append(g.Players, pid)
}

// Called when player asks to leave
func (g *Game) playerLeave(pid string) {
	g.game_state.Lock()
	defer g.game_state.Unlock()

	for i, s := range g.Players {
		if s == pid {
			g.Players[i] = g.Players[len(g.Players) - 1]
			g.Players = g.Players[:len(g.Players) - 1]
			break
		}
	}

	if len(g.Players) < 1 {
		g.EndGame()
	}
}

// For modification: code for when a player makes a move
// Return false if move is illegal
func (g *Game) tryMove(player string, msg MoveCardMessage) {
	g.game_state.Lock()
	defer g.game_state.Unlock()

	for _, pid := range g.Players {
		if p := getPlayer(pid); p != nil {
			p.moveCard(msg)
		}
	}
}