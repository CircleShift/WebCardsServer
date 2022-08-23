package webcode

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
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
	chanSync sync.Mutex
	closed bool
}

// Totally safe and definately won't crash or cause catastrophic errors.
// Source: trust me bro
func (a *AsyncWS) close() {
	a.chanSync.Lock()
	defer a.chanSync.Unlock()
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
			break
		}

		if rm.Type == "pong" {
			a.c.SetReadDeadline(time.Now().Add(pongWait))
		} else if !a.closed {
			a.chanSync.Lock()
			select {
			case a.I <- rm:
			case <-time.After(time.Second):
			}
			a.chanSync.Unlock()
		} else {
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
			break
		}
	}
}

func (a *AsyncWS) trySend(sm SendMessage) bool {
	if a.isClosed() {
		return false
	}
	a.chanSync.Lock()
	defer a.chanSync.Unlock()

	select {
	case a.O <- sm:
		return true
	case <-time.After(time.Second):
	}

	return false
}

func (a *AsyncWS) isClosed() bool {
	return a == nil || a.closed
}

func newAsyncWS(conn *websocket.Conn) *AsyncWS {
	out := AsyncWS{ c: conn, I: make(chan RecieveMessage), O: make(chan SendMessage), closed: false }

	go out.readLoop(&WebWG)
	go out.writeLoop(&WebWG)

	return &out
}