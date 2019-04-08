package main

import (
	"BountyWarServerG/game"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var addr = flag.String("addr", "192.168.0.146:8082", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer cl(c)

	game.PlayerConnect(c)
	c.SetCloseHandler(func(code int, text string) error {
		game.PlayerDisconnect(c)
		return nil
	})
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Error:", err)
			break
		}
		//log.Println("message:", message)
		game.IncomeMessage(message, c)
	}
}

func cl(c *websocket.Conn) {
	_ = c.Close()
}

func main() {
	go game.Cycle()
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/bountywar", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
