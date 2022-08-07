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
			for _, as := range c.clients{
				if as.isClosed() {
					continue
				}

				select {
				case as.O <- send:
				case <-time.After(time.Second*1):
					log.Println("Dropped message to client.")
				}
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

func (c *Chat) Clean() {
	c.chatLock.Lock()
	e := []*AsyncWS{}
	for _, as := range c.clients {
		if !as.isClosed() {
			e = append(e, as)
		}
	}
	c.clients = e
	c.chatLock.Unlock()
}

func (c *Chat) Shutdown() {
	c.shutdown = true
	close(c.Broadcast)
}
