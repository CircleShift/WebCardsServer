package webcode

import (
	"github.com/gorilla/websocket"
	"sync"
	"webcode/msg"
)

// LobbyLoop waits for either a shutdown signal or a new connection from a websocket through conchan.
// When a new *Conn is sent, LobbyLoop creates a new goroutine on lobby and adds one to wg.
// LobbyLoop calls wg.Done() on exit.
func LobbyLoop(conchan chan *websocket.Conn, shutdown chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case conn := <-conchan:
			wg.Add(1)
			go lobby(conn, shutdown, wg)
		case <-shutdown:
			return
		}
	}
}

// lobby is the hub for a websocket client. A client goes to lobby on connection.
// Once a client closes or errors, lobby exits.
// lobby calls wg.Done() on exit
func lobby(conn *websocket.Conn, shutdown chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	var m msg.Message
	for {
		err := conn.ReadJSON(&m)

		if err != nil {
			conn.Close()
			return
		}

		switch m.Name {
		case "games":
		}

		select {
		case <-shutdown:
			conn.Close()
			return

		default:
			continue
		}
	}
}
