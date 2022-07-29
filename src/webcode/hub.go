package webcode

import (
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"sync"
	"time"
	"game"
	"github.com/teris-io/shortid"
	"encoding/json"
)

var (
	MasterShutdown bool = false
	sidSync sync.Mutex
)

func generateID() (string, error) {
	sidSync.Lock()
	defer sidSync.Unlock()
	return shortid.GetDefault().Generate()
}

// HubLoop waits for either a shutdown signal or a new connection from a websocket through conchan.
// When a new *Conn is sent, LobbyLoop creates a new goroutine on lobby and adds one to wg.
// LobbyLoop calls wg.Done() on exit.
func HubLoop(conchan <-chan *websocket.Conn, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	initWebcode(wg)
	for {
		select {
		case conn := <-conchan:
			
			if pid := newPlayer(); pid != "" {
				go lobby(newAsyncWS(conn, wg), pid, wg)
			} else {
				conn.WriteJSON(SendMessage{"err", "maxcon"})
				conn.Close()
			}

		case <-time.After(time.Second*5):
			if MasterShutdown {
				return
			}
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
func lobby(async *AsyncWS, pid string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
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
	
	if !becomePlayer(async, pid) {
		async.O <- SendMessage{"err", "Unable to become player"}
		async.close()
		return
	}

	async.O <- SendMessage{"ready", game.DefaultSettingsMsg}

	for ;; {
		select {
		case msg, ok := <-async.I:
			if !ok {
				return
			}
			p := getPlayer(async, pid)
			switch msg.Type {
			case "chat":
				var d ChatMessage
				err := json.Unmarshal([]byte(msg.Data), &d)
				
				if err != nil {
					log.Println("Failed to parse chat msg")
					break
				}
				
				if c := p.getChat(d.ChatID); c != nil {
					select {
					case c.Broadcast <- ChatMessage{p.Options.Name, d.Text, p.Options.Color, d.ChatID, false}:

					case <-time.After(time.Second*1):
						log.Println("Dropping chat msg (overloaded?)")
					}
				} else {
					log.Println("User attempted to access invalid chat %s\n", d.ChatID)
				}
			case "options":
				var d game.UOptions
				err := json.Unmarshal([]byte(msg.Data), &d)
				
				if err != nil {
					log.Println("Failed to parse player options")
					break
				}
				p.Options = d

			case "ready":
				p.addChat("global")
			default:
				log.Printf("Not Implimented: %s\n", msg.Type)
			}
			log.Println(msg)
		case <-time.After(time.Second*1):
			if MasterShutdown {
				return
			}
		}
	}
}
