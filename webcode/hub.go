package webcode

import (
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"sync"
	"time"
	"github.com/teris-io/shortid"
	"encoding/json"
)

var (
	MasterShutdown bool = false
	sidSync sync.Mutex
	WebWG sync.WaitGroup
)

func generateID() (string, error) {
	sidSync.Lock()
	defer sidSync.Unlock()
	return shortid.GetDefault().Generate()
}

// HubLoop waits for either a shutdown signal or a new connection from a websocket through conchan.
// When a new *Conn is sent, LobbyLoop creates a new goroutine on lobby and adds one to wg.
// LobbyLoop calls wg.Done() on exit.
func HubLoop(conchan <-chan *websocket.Conn) {
	WebWG.Add(1)
	defer WebWG.Done()
	initWebcode()
	cleanTicker := time.NewTicker(10*time.Second)
	for {
		select {
		case conn := <-conchan:
			
			if pid := newPlayer(); pid != "" {
				go lobby(newAsyncWS(conn), pid)
			} else {
				conn.WriteJSON(SendMessage{"err", "maxcon"})
				conn.Close()
			}

		case <-time.After(time.Second*5):
			if MasterShutdown {
				return
			}
		}

		select {
		case <-cleanTicker.C:
			cleanML()
		default:
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
func lobby(async *AsyncWS, pid string) {
	WebWG.Add(1)
	defer WebWG.Done()

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
		async.trySend(SendMessage{"err", "Unable to become player"})
		async.close()
		return
	} else if !async.isClosed() {
		async.trySend(SendMessage{"ready", DefaultSettingsMsg})
	}

	for ;; {
		select {
		case msg, ok := <-async.I:
			if !ok {
				log.Printf("Player id:%s left.\n", pid)
				return
			}
			p := getPlayer(pid)
			switch msg.Type {
			case "chat":
				var d ChatMessage
				err := json.Unmarshal([]byte(msg.Data), &d)
				
				if err != nil {
					log.Println("Failed to parse chat msg")
					break
				}
				
				if c := p.getChat(d.ChatID); c != nil && !c.shutdown {
					select {
					case c.Broadcast <- ChatMessage{p.Options.Name, d.Text, p.Options.Color, d.ChatID, false}:

					case <-time.After(time.Second*1):
						log.Println("Dropping chat msg (overloaded?)")
					}
				} else {
					log.Println("User attempted to access invalid chat %s\n", d.ChatID)
				}
			case "options":
				var d UOptions
				err := json.Unmarshal([]byte(msg.Data), &d)
				
				if err != nil {
					log.Println("Failed to parse player options")
					break
				}
				p.Options = d
			case "create":
				var d GOptions
				err := json.Unmarshal([]byte(msg.Data), &d)

				if err != nil {
					log.Println("Failed to parse game options msg")
					break
				}

				if gid := newGame(d, pid); gid != "" {
					p.gameID = gid
				} else {
					log.Println("Unable to create new game")
					async.trySend(SendMessage{"err", "Game creation error"})
					break
				}
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
