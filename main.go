package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"

	webcode "cshift.net/webcards/webcode"
	card "cshift.net/webcards/card"
)

var up = websocket.Upgrader{}

var conchan = make(chan *websocket.Conn)

// Packs is all the packs able to be read when the program starts.

func upgrade(r http.ResponseWriter, c *http.Request) {
	client, err := up.Upgrade(r, c, nil)

	if err != nil {
		log.Println(err.Error())
		return
	}

	conchan <- client
}

func main() {
	s := flag.String("packdir", "packs", "Set the directory to search for extra card packs")
	p := flag.Int("port", 4040, "Port for the websocket server to run on")
	h := flag.String("host", "127.0.0.1", "Set the host to listen on")
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

	http.HandleFunc("/", upgrade)

	if *sec {
		log.Println(http.ListenAndServeTLS(*h+":"+port, *sec_c, *sec_k, nil))
	} else {
		log.Println(http.ListenAndServe(*h+":"+port, nil))
	}

	webcode.MasterShutdown = true
	webcode.WebWG.Wait()
}
