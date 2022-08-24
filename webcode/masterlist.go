package webcode

import (
	card "cshift.net/webcards/card"
	"sync"
	"log"
)

// Master lists for Players, Chats, and Games

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
	pub_game_lock sync.Mutex
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

func newChatID(name string) string {
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
		_, ok = chat_ml[cid]
	}

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
		_, ok = game_ml[gid]
	}
	
	if ok {
		return ""
	}

	game_ml[gid] = InitGame(o, gid, p)

	return gid
}

func getPlayer(pid string) *Player {
	if p, ok := player_ml[pid]; ok{
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

func getChat(cid string) *Chat {
	if c, ok := chat_ml[cid]; ok{
		return c
	}
	return nil
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

func addPublic(msg AddGameMessage) {
	pub_game_lock.Lock()
	pub_game_msg.Games = append(pub_game_msg.Games, msg)
	pub_game_lock.Unlock()
	player_lock.Lock()
	for _, p := range player_ml {
		p.as.trySend(SendMessage{"lobby", SendMessage{"addGame", msg}})
	}
	player_lock.Unlock()
}

func delPublic(id string) {
	pub_game_lock.Lock()
	for i, gm := range pub_game_msg.Games {
		if gm.GameID == id {
			pub_game_msg.Games[i] = pub_game_msg.Games[len(pub_game_msg.Games) - 1]
			pub_game_msg.Games = pub_game_msg.Games[0:len(pub_game_msg.Games) - 1]
			break
		}
	}
	pub_game_lock.Unlock()
	player_lock.Lock()
	for _, p := range player_ml {
		p.as.trySend(SendMessage{"lobby", SendMessage{"deleteGame", id}})
	}
	player_lock.Unlock()
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
		 if c.shutdown {
			delete(chat_ml, k)
		 }
	}

	chat_lock.Unlock()
	game_lock.Lock()

	for k, g := range game_ml {
		if g.end {
			log.Println("Deleted a game")
			delete(game_ml, k)
		}
	}

	game_lock.Unlock()
}