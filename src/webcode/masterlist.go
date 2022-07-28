package webcode

import (
	"game"
	"sync"
)

// Master lists for Players, Chats, and Games

type Player struct {
	Options game.UOptions
	gameID string
	chatList []string
}

const (
	MAX_PLAYERS = 1000
	MAX_CHATS = 1000
	MAX_GAMES = 500
)

var (
	player_ml map[string]*Player
	chat_ml map[string]*Chat
	game_ml map[string]*game.Game
)

func newPlayer() string {
	if len(player_ml) >= MAX_PLAYERS {
		return ""
	}

	out := Player{game.DefaultUserSettings, "", []string{"global"}}

	sid, err := generateID()
	_, ok := player_ml[sid]
	for i := 0; i < 100; i++ {
		if err != nil {
			return ""
		} else if !ok {
			break
		}
		sid, err = generateID()
	}

	_, ok = player_ml[sid]
	if !ok {
		return ""
	}

	player_ml[sid] = &out

	return sid
}

func newChat(name, id string, wg *sync.WaitGroup) {
	out := Chat{ Broadcast: make(chan ChatMessage), Name: name, clients: []*AsyncWS{}, shutdown: false}

	chats[id] = &out
	go out.loop(wg)
}

func newChatID(name string, wg *sync.WaitGroup) string {
	if len(chat_ml) >= MAX_CHATS {
		return ""
	}

	cid, err := generateID()
	_, ok := chat_ml[cid]
	for i := 0; i < 100; i++ {
		if err != nil {
			return ""
		} else if !ok {
			break
		}
		cid, err = generateID()
	}

	_, ok = chat_ml[cid]
	if !ok {
		return ""
	}

	newChat(name, cid, wg)
	return cid
}

func newGame() string {
	if len(game_ml) >= MAX_GAMES {
		return ""
	}

	gid, err := generateID()
	_, ok := game_ml[gid]
	for i := 0; i < 100; i++ {
		if err != nil {
			return ""
		} else if !ok {
			break
		}
		gid, err = generateID()
	}

	_, ok = game_ml[gid]
	if !ok {
		return ""
	}

	return gid
}

func getPlayer(pid string) *Player {
	if p, ok := player_ml[pid]; ok {
		return p
	}
	return nil
}

func (p *Player) hasChat(cid string) bool {
	for _, c := range p.chatList {
		if _, ok := chat_ml[c]; ok && c == cid {
			return true
		}
	}
	return false
}



func initWebcode() {

}