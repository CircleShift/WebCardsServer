package webcode

import (
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"sync"
)

// HubLoop waits for either a shutdown signal or a new connection from a websocket through conchan.
// When a new *Conn is sent, LobbyLoop creates a new goroutine on lobby and adds one to wg.
// LobbyLoop calls wg.Done() on exit.
func HubLoop(conchan chan *websocket.Conn, shutdown chan bool, wg *sync.WaitGroup) {
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

// validates websocket connection.
func handshake(conn *websocket.Conn) bool {
	var m RecieveMessage
	err := conn.ReadJSON(&m)

	if err != nil || m.Type != "version" {
		return false
	}

	i, err := strconv.Atoi(m.Data)

	if err != nil || i < 1 || i > 1 {
		conn.WriteJSON(SendMessage{Type: "verr", Data: ""})
		return false
	}

	return true
}

// lobby is the hub for a websocket client. A client goes to lobby on connection.
// Once a client closes or errors, lobby exits.
// lobby calls wg.Done() on exit
func lobby(conn *websocket.Conn, shutdown chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	if !handshake(conn) {
		log.Println("Handshake error, closing connection.")
		conn.Close()
		return
	}

	var m msg.Message
	for {
		err := conn.ReadJSON(&m)
		if err != nil {
			log.Println("Error reading JSON, closing conn.")
			conn.Close()
			return
		}

		switch m.Type {
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
