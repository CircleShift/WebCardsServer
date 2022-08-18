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
			p := getPlayer(pid)
			
			if !ok {
				log.Printf("Player id:%s left.\n", pid)
				if g := getGame(p.gameID); g != nil {
					g.playerLeave(pid)
				}
				p.gameID = ""
				p.chatList = []string{"global"}
				p.options = DefaultUserSettings
				return
			}

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
					case c.Broadcast <- ChatMessage{p.options.Name, d.Text, p.options.Color, d.ChatID, false}:

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
				p.options = d
			case "create":
				var d GOptions
				err := json.Unmarshal([]byte(msg.Data), &d)

				if err != nil {
					log.Println("Failed to parse game options msg")
					break
				} else if g := getGame(p.gameID); g != nil && !g.end {
					g.playerLeave(pid)
				}

				if gid := newGame(d, pid); gid != "" {
					p.gameID = gid
					p.joinGame()
					g := getGame(gid)
					g.playerJoin(pid)
				} else {
					log.Println("Unable to create new game")
					p.noJoinGame("Game creation error")
				}
			case "join":
				var d JoinGameMessage
				err := json.Unmarshal([]byte(msg.Data), &d)

				if err != nil {
					log.Println("Failed to parse join msg")
					break
				}

				g := getGame(d.GameID)

				if g_old := getGame(p.gameID); p.gameID != "" && g_old != nil {
					g_old.playerLeave(pid)
				}

				p.gameID = ""

				if g != nil && g.playerCanJoin(pid, d.Password) {
					p.gameID = d.GameID
					p.joinGame()
					g.playerJoin(pid)
				} else {
					p.noJoinGame("Failed to join")
				}
			case "leave":
				g := getGame(p.gameID)
				p.gameID = ""
				if g != nil {
					g.playerLeave(pid)
				}
				p.leaveGame()
			case "move":
				var d MoveCardMessage
				err := json.Unmarshal([]byte(msg.Data), &d)

				if err != nil {
					log.Println("Failed to parse move msg")
					break
				}

				if g := getGame(p.gameID); g != nil {
					g.tryMove(pid, d)
				}
			case "ready":
				p.addChat("global")
				pub_game_lock.Lock()
				log.Println(pub_game_msg)
				async.trySend(SendMessage{"lobby", SendMessage{"gameList", pub_game_msg}})
				pub_game_lock.Unlock()
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
