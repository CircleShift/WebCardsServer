package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"fmt"
	"time"

	"github.com/gorilla/websocket"

	webcode "cshift.net/webcards/webcode"
	card "cshift.net/webcards/card"
)

var up = websocket.Upgrader{}

var conchan = make(chan *websocket.Conn)

// Packs is all the packs able to be read when the program starts.

func passConn(r http.ResponseWriter, c *http.Request) {
	client, err := up.Upgrade(r, c, nil)

	if err != nil || webcode.MasterShutdown {
		log.Println(err.Error())
		return
	}

	conchan <- client
}

func handleUsrInput() {
	var i string
	for {
		fmt.Scanln(&i)
		switch i {
		case "exit":
			return
		default:
			log.Printf("Command not recognized %s", i)
		}
	}
}

func main() {
	s := flag.String("packdir", "packs", "Set the directory to search for extra card packs")
	p := flag.Int("port", 4040, "Port for the websocket server to run on")
	h := flag.String("host", "", "Set the host to listen on")
	sec := flag.Bool("tls", false, "Use HTTPS/WSS instead of HTTP/WS")
	sec_c := flag.String("cert", "host.csr", "Cert file for tls")
	sec_k := flag.String("key", "host.key", "Key file for tls")

	flag.Parse()

	up.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	err := card.InitCardPacks(*s)

	if err != nil {
		log.Println(err)
		log.Fatal("Failed to read default.json, exiting.")
	}

	log.Println(card.Packs)

	port := strconv.Itoa(*p)

	log.Println("Starting server on " + *h + ":" + port)

	go webcode.HubLoop(conchan)

	server := &http.Server {
		Addr: *h+":"+port,
		Handler: http.HandlerFunc(passConn),
		ReadTimeout: 10*time.Second,
		WriteTimeout: 10*time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if *sec {
		go server.ListenAndServeTLS(*sec_c, *sec_k)
	} else {
		go server.ListenAndServe()
	}

	defer server.Close()

	handleUsrInput()

	webcode.MasterShutdown = true
	webcode.WebWG.Wait()
}
