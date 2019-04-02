package main

import (
	"BountyWarServerG/game"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var addr = flag.String("addr", "192.168.111.164:8082", "http service address")

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
			log.Println("Error:")
			break
		}
		game.IncomeMessage(message, c)
		/*
			log.Printf("recv type:%d, mess:%s\n", mt, message)
			err = c.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}*/
	}
}

func cl(c *websocket.Conn) {
	_ = c.Close()
}

func main() {
	//fmt.Println(game.KeyGo, " ", game.KeyBack, " ", game.KeyLeft, " ", game.KeyRight, " ", game.KeyFire)
	go game.Cycle()
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/bountywar", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
