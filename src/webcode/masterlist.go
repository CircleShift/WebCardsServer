package webcode

import (
	"card"
	"sync"
	"time"
	"log"
)

// Master lists for Players, Chats, and Games

type Player struct {
	Options UOptions
	gameID string
	chatList []string
	as *AsyncWS
	clean int
}

const (
	MAX_PLAYERS = 1000
	MAX_CHATS = 1000
	MAX_GAMES = 500
	DELETE_CYCLES = 2
)

var (
	player_lock sync.Mutex
	player_ml map[string]*Player
	chat_lock sync.Mutex
	chat_ml map[string]*Chat
	game_lock sync.Mutex
	game_ml map[string]*Game
	pub_game_msg GameListMessage
)

func newPlayer() string {
	player_lock.Lock()
	defer player_lock.Unlock()

	if len(player_ml) >= MAX_PLAYERS {
		return ""
	}

	out := Player{DefaultUserSettings, "", []string{}, nil, 0}

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

func newChat(name, id string) {
	chat_lock.Lock()
	defer chat_lock.Unlock()

	_, ok := chat_ml[id]
	if !ok {
		out := Chat{ Broadcast: make(chan ChatMessage), Name: name, clients: []*AsyncWS{}, shutdown: false}
		chat_ml[id] = &out
		go out.loop(&WebWG)
	}
}

func NewChatID(name string) string {
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

	newChat(name, cid)
	return cid
}

func newGame(o GOptions, p string) string {
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

	game_ml[gid] = InitGame(o, p)

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

func delPlayer(pid string) bool {
	player_lock.Lock()
	defer player_lock.Unlock()

	if _, ok := player_ml[pid]; ok {
		delete(player_ml, pid)
		return true
	}
	return false
}

func getPlayer(as *AsyncWS, pid string) *Player {
	if p, ok := player_ml[pid]; ok && p.as == as {
		return p
	}
	return nil
}

func getGame(gid string) *Game {
	if g, ok := game_ml[gid]; ok{
		return g
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
	for _, cid := range p.chatList {
		if cid == id {
			return true
		}
	}

	c, ok := chat_ml[id]
	if p == nil || !ok || p.as.isClosed() {
		return false
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

func delChat(cid string) bool {
	chat_lock.Lock()
	defer chat_lock.Unlock()

	if c, ok := chat_ml[cid]; ok {
		c.Shutdown()
		delete(chat_ml, cid)
		return true
	}
	return false
}

func delGame(gid string) bool {
	game_lock.Lock()
	defer game_lock.Unlock()

	if g, ok := game_ml[gid]; ok {
		delete(game_ml, gid)
		for _, cid := range g.Chats {
			delChat(cid)
		}
		return true
	}
	return false
}

func initWebcode() {
	player_ml = make(map[string]*Player)
	chat_ml = make(map[string]*Chat)
	game_ml = make(map[string]*Game)
	pub_game_msg = GameListMessage{GAME_NAME, len(card.Packs), []AddGameMessage{}}

	newChat("Global", "global")
}

func cleanML() {
	player_lock.Lock()

	log.Println("Clean cycle")
	for k, p := range player_ml {
		if p.as == nil || p.as.isClosed() {
			if p.clean < DELETE_CYCLES - 1 {
				p.clean += 1
			} else {
				delete(player_ml, k)
			}
		}
	}

	player_lock.Unlock()
	chat_lock.Lock()

	for k, c := range chat_ml {
		 c.Clean()
		 if (len(c.clients) == 0 && k != "global") || c.shutdown {
			delete(chat_ml, k)
		 }
	}

	chat_lock.Unlock()
}