package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	port := "3333"
	http.HandleFunc("/ws", socketHandler)

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

var upgrader = websocket.Upgrader{
	// ReadBufferSize:  1024,
	// WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	var req Message
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrader.Upgrade: %+v", err)
		return
	}

	go func() {
		for {
			if err := conn.ReadJSON(&req); err != nil {
				return
			}

			conn.WriteJSON(Message{})
		}
	}()
}
