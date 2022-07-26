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
	s <-chan bool
	I chan RecieveMessage
	O chan SendMessage
}

func (a *AsyncWS) close() {
	a.c.Close()

	select {
	case _, ok := <-a.I:
		if !ok {
			break
		} else {
			close(a.I)
		}
	default:
		close(a.I)
	}

	select {
	case _, ok := <-a.O:
		if !ok {
			break
		} else {
			close(a.O)
		}
	default:
		close(a.O)
	}
}

func (a *AsyncWS) readLoop(wg *sync.WaitGroup) {
	defer wg.Done()
	defer log.Println("Done reading.")

	var rm RecieveMessage
	for ;; {
		err := a.c.ReadJSON(&rm)
		if err != nil {
			a.close()
			break
		}

		if rm.Type == "pong" {
			a.c.SetReadDeadline(time.Now().Add(pongWait))
		} else {
			a.I <- rm
		}
	}
}

func (a *AsyncWS) writeLoop(wg *sync.WaitGroup) {
	defer wg.Done()
	defer log.Println("Done writing.")

	var sm SendMessage
	var ok bool

	for ;; {
		
		select {
		case <-a.s:
			a.close()
			return

		case sm, ok = (<-a.O):
			if ok {
				a.c.SetReadDeadline(time.Now().Add(pongWait))
			}

		case <-time.After(pongWait / 2):
			sm = pingMsg
		}

		err := a.c.WriteJSON(sm)
		if err != nil {
			a.close()
			break
		}
	}
}

func initAsyncWS(conn *websocket.Conn, shutdown <-chan bool, wg *sync.WaitGroup) AsyncWS {
	out := AsyncWS{ conn, shutdown, make(chan RecieveMessage), make(chan SendMessage) }

	go out.readLoop(wg)
	go out.writeLoop(wg)

	return out
}