package webcode

import (
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"sync"
	"time"
)

const MAX_PLAYERS = 1000

var (
	plcSync sync.Mutex
	plc int = 0
)

func addPlayer() bool {
	defer plcSync.Unlock()
	plcSync.Lock()

	if plc == MAX_PLAYERS {
		return false
	}
	plc += 1
	return true
}

func delPlayer() {
	plcSync.Lock()
	if plc <= 0 {
		plc = 1
	}
	plc -= 1
	plcSync.Unlock()
}

// HubLoop waits for either a shutdown signal or a new connection from a websocket through conchan.
// When a new *Conn is sent, LobbyLoop creates a new goroutine on lobby and adds one to wg.
// LobbyLoop calls wg.Done() on exit.
func HubLoop(conchan <-chan *websocket.Conn, shutdown <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case conn := <-conchan:
			wg.Add(3)
			
			if addPlayer() {
				go lobby(initAsyncWS(conn, shutdown, wg), wg)
			}

		case <-shutdown:
			return
		}
	}
}

// validates websocket connection.
func handshake(msg RecieveMessage) bool {

	if msg.Type != "version" {
		log.Println("Not a version.")
		return false
	}

	i, err := strconv.Atoi(msg.Data)

	if err != nil || i != 1 {
		log.Println(i)
		return false
	}

	return true
}

// lobby is the hub for a websocket client. A client goes to lobby on connection.
// Once a client closes or errors, lobby exits.
// lobby calls wg.Done() on exit
// Only write to a conn in this group
func lobby(async AsyncWS, wg *sync.WaitGroup) {
	defer wg.Done()
	defer delPlayer()
	defer log.Println("Exit lobby")

	// First msg should be a handshake
	select {
	case msg, _ := (<-async.I):
		if !handshake(msg) {
			log.Println("Handshake error, closing connection.")
			return
		}
	case <-time.After(time.Second * 1):
		async.close()
	}
	
	log.Println("Connection ok")
	async.O <- SendMessage{"ready", ""}

	// Todo: fix
	//game_id := ""
	for ;; {
		select {
		case _, ok := <-async.I:
			if !ok {
				return
			}


		}
	}
}
