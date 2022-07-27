package webcode

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
	"log"
)

var (
	pongWait = 30 * time.Second
	pingMsg = SendMessage{"ping", ""}
)

// AsyncWS Makes websocket read/writes into a channel operation, allowing for selection and timeouts
type AsyncWS struct {
	c *websocket.Conn
	I chan RecieveMessage
	O chan SendMessage
	closeSync sync.Mutex
	closed bool
}

// Totally safe and definately won't crash or cause catastrophic errors.
// Source: trust me bro
func (a *AsyncWS) close() {
	a.closeSync.Lock()
	defer a.closeSync.Unlock()
	if a.closed {
		return
	}
	a.closed = true
	a.c.Close()
	close(a.I)
	close(a.O)
}

func (a *AsyncWS) readLoop(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	var rm RecieveMessage
	for {
		err := a.c.ReadJSON(&rm)
		if err != nil || MasterShutdown {
			a.close()
			log.Printf("Done reading. %v\n", a.closed)
			break
		}

		if rm.Type == "pong" {
			a.c.SetReadDeadline(time.Now().Add(pongWait))
		} else if !a.closed {
			a.I <- rm
		} else {
			log.Printf("Done reading. %v\n", a.closed)
			break
		}
	}
}

func (a *AsyncWS) writeLoop(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	var sm SendMessage
	var ok bool

	for ;; {
		
		select {
		case sm, ok = (<-a.O):
			if ok {
				a.c.SetReadDeadline(time.Now().Add(pongWait))
			}

		case <-time.After(pongWait / 2):
			sm = pingMsg
		}

		err := a.c.WriteJSON(sm)
		if err != nil || a.closed || MasterShutdown {
			a.close()
			log.Printf("Done writing. %v\n", a.closed)
			break
		}
	}
}

func (a *AsyncWS) IsClosed() bool {
	return a.closed
}

func initAsyncWS(conn *websocket.Conn, wg *sync.WaitGroup) *AsyncWS {
	out := AsyncWS{ c: conn, I: make(chan RecieveMessage), O: make(chan SendMessage), closed: false }

	go out.readLoop(wg)
	go out.writeLoop(wg)

	return &out
}