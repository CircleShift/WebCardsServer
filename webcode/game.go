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

	// Deck
	Cards map[int]card.Ref
	Decks map[string][]int
	CardPool *card.Deck
	NCID int

	// password for game ("" means no password)
	pass string

	// The game object is going to be used by multiple threads,
	// so it's important to use a mutex to prevent race conditions
	game_state sync.Mutex
}

// For modification: init code
func InitGame(options GOptions, id string, player_one string) *Game {
	var out Game
	
	out.Cards = make(map[int]card.Ref)
	out.Decks = make(map[string][]int)

	out.Chats = []string{}
	out.Players = []string{}
	out.end = false
	out.id = id
	out.NCID = 3
	out.turn = player_one

	if options.UsePassword {
		out.pass = options.Password
	} else {
		out.pass = ""
	}
	
	if !options.Hidden {
		log.Println("Adding public game!")
		addPublic(AddGameMessage{options.Name, 0, id, options.UsePassword})
	}

	out.CardPool = card.NewDeck([]*card.Pack{&(card.Packs[0])}, -1)

	out.newRandomCard(true)
	out.Decks["1"] = []int{1}

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
		p.newDeck(NewDeckMessage{"0", DeckOptions{"stack", "one", 0, [4]float64{0.05, 0.05, 0.0, 1.0}}})
		p.newDeck(NewDeckMessage{"turn", DeckOptions{"stack", "one", 0, [4]float64{0.5, 0.05, 0.0, 1.0}}})
		p.newDeck(NewDeckMessage{"1", DeckOptions{"stack", "one", 0, [4]float64{0.95, 0.05, 0.0, 1.0}}})
		p.newDeck(NewDeckMessage{pid, DeckOptions{"strip-hr", "one", 0, [4]float64{0.5, 0.95, 0.0, 1.0}}})
		p.newCard(NewCardMessage{2, "turn", card.Packs[0].GetCard("red", "0").Data})
		p.newCard(NewCardMessage{0, "0", card.Packs[0].GetCard("draw", "draw").Data})
		p.newCard(NewCardMessage{1, "1", g.TopCard().GetCard().Data})
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
	if len(g.Players) == 1 {
		if p := getPlayer(pid); p != nil {
			p.replaceCard(SwapCardMessage{2, 2, card.Packs[0].GetCard("green", "0").Data})
		}
		g.turn = pid
	}
	g.Decks[pid] = []int{}
	if p := getPlayer(pid); p != nil {
		for i := 0; i < 7; i++ {
			c := g.newRandomCard(false)
			p.newCard(NewCardMessage{c, pid, g.Cards[c].GetCard().Data})
			g.Decks[pid] = append(g.Decks[pid], c)
		}
	}
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

// Should only be called from functions that have state lock
func (g *Game) nextTurn() {
	//if len(g.Players) == 1 {
	//	return
	//}

	for i, s := range g.Players {
		if s == g.turn {
			if i == len(g.Players) - 1 {
				g.turn = g.Players[0]
			} else {
				g.turn = g.Players[i + 1]
			}
			return
		}
	}
	g.turn = g.Players[0]
}

func (g *Game) newRandomCard(s bool) int {
	if s {
		g.Cards[1] = g.CardPool.GetRandomCard(g.CardPool.GetCardChance())
		return 1
	}
	g.Cards[g.NCID] = g.CardPool.GetRandomCard(g.CardPool.GetCardChance())
	g.NCID = g.NCID + 1
	return g.NCID - 1
}

func (g *Game) hasCard(player string, cid int) int {
	for i, d := range g.Decks[player] {
		if d == cid {
			return i
		}
	}
	return -1
}

func (g *Game) CurrentDeck(cid int) (string, int) {
	for s, _ := range g.Decks {
		i := g.hasCard(s, cid)
		if i > -1 {
			return s, i
		}
	}
	return "", -1
}

// Should only be called from functions that have state lock
func (g *Game) ReturnCard(cid int, pid string) {
	log.Println("Returning card...")
	d, i := g.CurrentDeck(cid)
	if p := getPlayer(pid); p != nil {
		p.moveCard(MoveCardMessage{cid, d, i})
	}
}

func (g *Game) TopCard() card.Ref {
	return g.Cards[g.Decks["1"][len(g.Decks["1"]) - 1]]
}

func (g *Game) CanPlay(cid int) bool {
	top := g.TopCard()
	chk := g.Cards[cid]
	log.Println(chk.Suit, top.Suit)
	log.Println(chk.GetCard().Name, top.GetCard().Name)
	return chk.Suit == top.Suit || chk.GetCard().Name == top.GetCard().Name || chk.Suit == "wild"
}

// For modification: code for when a player makes a move
// Return false if move is illegal
func (g *Game) tryMove(player string, msg MoveCardMessage) {
	g.game_state.Lock()
	defer g.game_state.Unlock()

	if player != g.turn {
		g.ReturnCard(msg.CardID, player)
		return
	}

	if msg.DeckID == "1" {

		if g.hasCard(player, msg.CardID) < 0 || !g.CanPlay(msg.CardID) {
			g.ReturnCard(msg.CardID, player)
			return
		}

		d := g.Decks["1"][0]
		r := g.Cards[msg.CardID].GetCard().Data
		
		msg.Index = 1000
		
		for _, pid := range g.Players {
			if p := getPlayer(pid); p != nil {
				if pid == player {
					p.moveCard(msg)
				} else {
					p.newCard(NewCardMessage{msg.CardID, msg.DeckID, r})
				}
				p.deleteCard(d)
			}
		}

		g.Decks["1"] = append(g.Decks["1"][1:], msg.CardID)
		delete(g.Cards, d)
	} else if msg.CardID == 0 {
		c := g.CardPool.GetRandomCard(g.CardPool.GetCardChance())

		log.Println(g.NCID)
		g.Decks[player] = append(g.Decks[player], g.NCID)
		g.Cards[g.NCID] = c

		if p := getPlayer(player); p != nil {
			p.moveCard(MoveCardMessage{0, player, 0})
			p.replaceCard(SwapCardMessage{0, g.NCID, c.GetCard().Data})
			p.newCard(NewCardMessage{0, "0", card.Packs[0].GetCard("draw", "draw").Data})
		}

		g.NCID = g.NCID + 1
	} else {
		g.ReturnCard(msg.CardID, player)
	}

	if p := getPlayer(g.turn); p != nil {
		p.replaceCard(SwapCardMessage{2, 2, card.Packs[0].GetCard("red", "0").Data})
	}
	g.nextTurn()
	if p := getPlayer(g.turn); p != nil {
		p.replaceCard(SwapCardMessage{2, 2, card.Packs[0].GetCard("green", "0").Data})
	}
}