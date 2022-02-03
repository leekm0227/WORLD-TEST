package server

import (
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	MAX_SIZE         = 100
	CHAN_BUFFER_SIZE = 2
)

type MessageType int

const (
	JOIN  MessageType = 0
	LEAVE MessageType = 1
	INIT  MessageType = 2
	MOVE  MessageType = 3
)

type Message struct {
	MessageType MessageType `json:"messageType"`
	Payload     interface{} `json:"payload"`
}

type Channel struct {
	join      chan *Client
	leave     chan *Client
	broadcast chan Message
	clients   map[*Client]bool
}

type Client struct {
	uuid   string
	conn   *websocket.Conn
	action chan Message
}

var channel = Channel{
	join:      make(chan *Client),
	leave:     make(chan *Client),
	broadcast: make(chan Message),
	clients:   make(map[*Client]bool),
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Run(port string) {
	go func() {
		for {
			select {
			case client := <-channel.join:
				channel.clients[client] = true

				x, y := rand.Intn(9)+1, rand.Intn(9)+1
				client.conn.WriteJSON(Message{
					MessageType: INIT,
					Payload: map[string]interface{}{
						"id": client.uuid,
						"x":  x,
						"y":  y,
					},
				})

				message := Message{
					MessageType: JOIN,
					Payload: map[string]interface{}{
						"id": client.uuid,
						"x":  x,
						"y":  y,
					},
				}

				for client := range channel.clients {
					client.conn.WriteJSON(message)
				}
			case client := <-channel.leave:
				delete(channel.clients, client)
				message := Message{
					MessageType: LEAVE,
					Payload:     map[string]string{"id": client.uuid},
				}

				for client := range channel.clients {
					client.conn.WriteJSON(message)
				}
			case message := <-channel.broadcast:
				for client := range channel.clients {
					client.conn.WriteJSON(message)
				}
			}
		}
	}()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var req Message
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("upgrader.Upgrade: %+v", err)
			return
		}

		client := &Client{
			uuid:   strings.Split(uuid.NewString(), "-")[0],
			conn:   conn,
			action: make(chan Message),
		}

		channel.join <- client
		defer func(client *Client) {
			client.conn.Close()
			channel.leave <- client
		}(client)

		for {
			if err := conn.ReadJSON(&req); err != nil {
				log.Println(err)
				return
			}

			switch req.MessageType {
			case MOVE:
				channel.broadcast <- req
			}
		}
	})

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
