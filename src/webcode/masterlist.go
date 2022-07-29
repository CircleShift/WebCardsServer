package webcode

import (
	"game"
	"sync"
	"time"
)

// Master lists for Players, Chats, and Games

type Player struct {
	Options game.UOptions
	gameID string
	chatList []string
	as *AsyncWS
}

const (
	MAX_PLAYERS = 1000
	MAX_CHATS = 1000
	MAX_GAMES = 500
)

var (
	player_lock sync.Mutex
	player_ml map[string]*Player
	chat_lock sync.Mutex
	chat_ml map[string]*Chat
	game_lock sync.Mutex
	game_ml map[string]*game.Game
)

func newPlayer() string {
	player_lock.Lock()
	defer player_lock.Unlock()

	if len(player_ml) >= MAX_PLAYERS {
		return ""
	}

	out := Player{game.DefaultUserSettings, "", []string{}, nil}

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
	if ok {
		return ""
	}

	player_ml[sid] = &out

	return sid
}

func newChat(name, id string, wg *sync.WaitGroup) {
	chat_lock.Lock()
	defer chat_lock.Unlock()

	_, ok := chat_ml[id]
	if !ok {
		out := Chat{ Broadcast: make(chan ChatMessage), Name: name, clients: []*AsyncWS{}, shutdown: false}
		chat_ml[id] = &out
		go out.loop(wg)
	}
}

func newChatID(name string, wg *sync.WaitGroup) string {
	chat_lock.Lock()
	defer chat_lock.Unlock()

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
	if ok {
		return ""
	}

	newChat(name, cid, wg)
	return cid
}

func newGame() string {
	game_lock.Lock()
	defer game_lock.Unlock()

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
	if ok {
		return ""
	}

	return gid
}

func becomePlayer(as *AsyncWS, pid string) bool {
	player_lock.Lock()
	defer player_lock.Unlock()

	if p, ok := player_ml[pid]; ok {
		if p.as == nil || p.as.isClosed() {
			p.as = as
			return true
		}
		return false
	}
	return false
}

func getPlayer(as *AsyncWS, pid string) *Player {
	if p, ok := player_ml[pid]; ok && p.as == as {
		return p
	}
	return nil
}

func (p *Player) getChat(cid string) *Chat {
	if p == nil {
		return nil
	}

	for _, c := range p.chatList {
		if c == cid {
			return chat_ml[c]
		}
	}

	return nil
}

func (p *Player) addChat(id string) bool {
	c, ok := chat_ml[id]
	if p == nil || !ok || p.as.isClosed() {
		return false
	}

	for _, cid := range p.chatList {
		if cid == id {
			return true
		}
	}

	select {
	case p.as.O <- SendMessage{"chat", SendMessage{"addChannel", AddChatMessage{c.Name, id, true}}}:
		p.chatList = append(p.chatList, id)
		c.AddPlayer(p.as)
		return true
	case <-time.After(1*time.Second):
	}
	return false
}

func initWebcode(wg *sync.WaitGroup) {
	player_ml = make(map[string]*Player)
	chat_ml = make(map[string]*Chat)
	game_ml = make(map[string]*game.Game)

	newChat("Global", "global", wg)
}