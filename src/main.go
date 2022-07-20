package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"

	"card"
	"webcode"
)

var up = websocket.Upgrader{}

var conchan = make(chan *websocket.Conn)
var shutdown = make(chan bool)

// Packs is all the packs able to be read when the program starts.

func upgrade(r http.ResponseWriter, c *http.Request) {
	client, err := up.Upgrade(r, c, nil)

	if err != nil {
		log.Println(err.Error())
	}

	conchan <- client
}

func main() {
	s := flag.String("setdir", "sets", "Set the directory to search for cardsets")
	p := flag.Int("port", 4040, "Port for the websocket server to run on")
	h := flag.String("host", "127.0.0.1", "Set the host to listen on")

	flag.Parse()

	err := card.InitCardPacks(*s)

	if err != nil {
		log.Println(err)
		log.Fatal("Failed to read default.json, exiting.")
	}

	log.Println(card.Packs)

	port := strconv.Itoa(*p)

	log.Println("Starting server on " + *h + ":" + port)

	var wg sync.WaitGroup
	wg.Add(1)
	go webcode.HubLoop(conchan, shutdown, &wg)

	http.HandleFunc("/", upgrade)

	log.Println(http.ListenAndServe(*h+":"+port, nil))

	close(shutdown)
	wg.Wait()
}
