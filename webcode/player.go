package webcode

import (
	"log"
	"time"
)

type Player struct {
	Options UOptions
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
	for _, cid := range p.chatList {
		if cid == id {
			return true
		}
	}

	c := getChat(id)
	if p == nil || c == nil || p.as.isClosed() {
		return false
	}

	select {
	case p.as.O <- SendMessage{"chat", SendMessage{"addChannel", AddChatMessage{c.Name, id, true}}}:
		p.chatList = append(p.chatList, id)
		c.AddPlayer(p.as)
		return true
	case <-time.After(1*time.Second):
		log.Println("Failed to add chat to player")
	}
	return false
}

func (p *Player) delChat(id string) {
	if p == nil || p.as.isClosed() {
		return
	}

	for i, cid := range p.chatList {
		if cid == id {
			if p.as.isClosed() {
				return
			}

			select {
			case p.as.O <- SendMessage{"chat", SendMessage{"deleteChannel", id}}:
				if c := getChat(id); c != nil {
					c.DelPlayer(p.as)
					cpy := []string{}
					cpy = append(cpy, p.chatList[0:i]...)
					cpy = append(cpy, p.chatList[i+1:]...)
				}
			case <-time.After(1*time.Second):
				log.Println("Failed to remove channel from player");
				return
			}

			return
		}
	}
}

func (p *Player) noJoinGame(reason string) {
	if p == nil || p.as.isClosed() {
		return
	}

	select {
	case p.as.O <- SendMessage{"game", SendMessage{"nojoin", reason}}:
	case <-time.After(1*time.Second):
		log.Println("Failed to inform player of nojoin");
	}
}

// Assuming that one player is only ever in one game.
func (p *Player) joinGame() {
	if p == nil || p.as.isClosed() {
		return
	}

	select {
	case p.as.O <- SendMessage{"game", SendMessage{"join", ""}}:
	case <-time.After(1*time.Second):
		log.Println("Failed to inform player of nojoin");
	}
}

func (p *Player) leaveGame() {
	if p == nil || p.as.isClosed() {
		return
	}

	p.gameID = ""

	select {
	case p.as.O <- SendMessage{"game", SendMessage{"leave", ""}}:
	case <-time.After(1*time.Second):
		log.Println("Failed to inform player of leave");
	}
}