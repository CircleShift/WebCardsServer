package webcode

import (
	"time"
	"log"
	"sync"
)

// Todo: impl
type Chat struct {
	Broadcast chan ChatMessage
	Name string

	clients []*AsyncWS
	chatLock sync.Mutex
	shutdown bool
}

var (
	lockChats sync.Mutex
	chatCt int = 0
	chats map[string]*Chat
)

func (c *Chat) loop(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	for {
		select {
		case msg, ok := <-c.Broadcast:
			if !ok {
				return
			}
			send := SendMessage{"chat", SendMessage{"recieveMessage", msg}}
			e := []*AsyncWS{}
			for i := 0; i < len(c.clients); i++ {
				log.Println(*(c.clients[i]))
				if c.clients[i].IsClosed() {
					e = append(e, c.clients[i])
					continue
				}
				select {
				case c.clients[i].O <- send:
				case <-time.After(time.Second*1):
					log.Println("Dropped message to client.")
				}
			}
			for _, erase := range e {
				c.DelPlayer(erase)
			}
		case <-time.After(time.Second*1):
		}

		if MasterShutdown || c.shutdown {
			break
		}
	}
}

func (c *Chat) AddPlayer(a *AsyncWS) {
	c.chatLock.Lock()
	c.clients = append(c.clients, a)
	c.chatLock.Unlock()
}

func (c *Chat) DelPlayer(a *AsyncWS) {
	c.chatLock.Lock()
	for i := 0; i < len(c.clients); i++ {
		if c.clients[i] == a {
			cpy := []*AsyncWS{}
			cpy = append(cpy, c.clients[0:i]...)
			c.clients = append(cpy, c.clients[i+1:]...)
			break
		}
	}
	c.chatLock.Unlock()
}

func (c *Chat) Shutdown() {
	c.shutdown = true
}

func NewChat(name string, wg *sync.WaitGroup) string {
	out := Chat{ Broadcast: make(chan ChatMessage), Name: name, clients: []*AsyncWS{}, shutdown: false}

	id, _ := generateID()

	for _, ok := chats[id]; ok; _, ok = chats[id] {
		id, _ = generateID()
	}

	chats[id] = &out
	go out.loop(wg)

	return id
}

func InitGlobal(wg *sync.WaitGroup) {
	if _, ok := chats["global"]; ok {
		return
	}
	chats = make(map[string]*Chat)

	out := Chat{ Broadcast: make(chan ChatMessage), Name: "Global", clients: []*AsyncWS{}, shutdown: false}

	chats["global"] = &out
	go out.loop(wg)
}

func DelChat(id string) {
	if _, ok := chats[id]; ok {
		chats[id].Shutdown()
		delete(chats, id)
	}
}

func GetChat(id string) *Chat {
	if c, ok := chats[id]; ok {
		return c
	}
	return nil
}