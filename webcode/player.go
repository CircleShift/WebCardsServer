package webcode

import (
	"log"
)

type Player struct {
	options UOptions
	gameID string
	chatList []string
	as *AsyncWS
	clean int
}

func (p *Player) getChat(cid string) *Chat {
	if p == nil {
		return nil
	}

	for _, c := range p.chatList {
		if c == cid {
			return getChat(cid)
		}
	}

	return nil
}

func (p *Player) addChat(id string) bool {
	c := getChat(id)
	if p == nil || c == nil {
		return false
	}

	replace := -1
	for i, cid := range p.chatList {
		if cid == id {
			return true
		} else if cid == "" {
			replace = i
		}
	}

	if p.as.trySend(SendMessage{"chat", SendMessage{"addChannel", AddChatMessage{c.Name, id, true}}}) {
		c.AddPlayer(p.as)
		
		if replace != -1 {
			p.chatList[replace] = id
		} else {
			p.chatList = append(p.chatList, id)
		}
		
		return true
	} else {
		log.Println("Failed to add chat to player")
	}

	return false
}

func (p *Player) delChat(id string) {
	c := getChat(id)
	if c != nil || p == nil {
		return
	}

	for i, cid := range p.chatList {
		if cid == id {
			p.chatList[i] = ""
			if p.as.isClosed() || p.as.trySend(SendMessage{"chat", SendMessage{"deleteChannel", id}}) {
				c.DelPlayer(p.as)
			} else {
				log.Println("Failed to remove channel from player");
			}

			return
		}
	}
}

func (p *Player) noJoinGame(reason string) {
	p.as.trySend(SendMessage{"game", SendMessage{"nojoin", reason}})
}

// Assuming that one player is only ever in one game.
func (p *Player) joinGame() {
	p.as.trySend(SendMessage{"game", SendMessage{"join", ""}})
}

func (p *Player) leaveGame() {
	p.gameID = ""
	p.as.trySend(SendMessage{"game", SendMessage{"leave", ""}})
}

func (p *Player) newDeck(d NewDeckMessage) {
	p.as.trySend(SendMessage{"game", SendMessage{"newDeck", d}})
}

func (p *Player) deleteDeck(id string) {
	p.as.trySend(SendMessage{"game", SendMessage{"deleteDeck", id}})
}

func (p *Player) newCard(c NewCardMessage) {
	p.as.trySend(SendMessage{"game", SendMessage{"newCard", c}})
}

func (p *Player) deleteCard(id string) {
	p.as.trySend(SendMessage{"game", SendMessage{"deleteCard", id}})
}

func (p *Player) moveCard(m MoveCardMessage) {
	p.as.trySend(SendMessage{"game", SendMessage{"moveCard", m}})
}

func (p *Player) replaceCard(s SwapCardMessage) {
	p.as.trySend(SendMessage{"game", SendMessage{"swapCard", s}})
}

