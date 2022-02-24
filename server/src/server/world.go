package server

import (
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var world = World{
	wid:       "world-0",
	join:      make(chan *Client),
	leave:     make(chan *Client),
	move:      make(chan Message),
	attack:    make(chan Message),
	clientMap: make(map[*Client]bool),
	playerMap: make(map[string]Player),
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Run(port string) {
	go runWorld()
	http.HandleFunc("/ws", wsHandler)
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
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

	world.join <- client
	defer func(client *Client) {
		world.leave <- client
		client.conn.Close()
	}(client)

	for {
		if err := conn.ReadJSON(&req); err != nil {
			// log.Println(err)
			return
		}

		switch req.MessageType {
		case ATTACK:
			world.attack <- req
		case MOVE:
			world.move <- req
		}
	}
}

func runWorld() {
	for {
		select {
		case client := <-world.join:
			player := Player{
				Id: client.uuid,
				Hp: 10,
				X:  rand.Intn(X_MAX),
				Y:  rand.Intn(Y_MAX),
			}

			client.conn.WriteJSON(Message{
				MessageType: INIT,
				Payload: map[string]interface{}{
					"player":  player,
					"players": world.playerMap,
				},
			})

			for c := range world.clientMap {
				c.conn.WriteJSON(Message{
					MessageType: JOIN,
					Payload:     player,
				})
			}

			world.playerMap[client.uuid] = player
			world.clientMap[client] = true
		case client := <-world.leave:
			delete(world.clientMap, client)
			delete(world.playerMap, client.uuid)

			message := Message{
				MessageType: LEAVE,
				Payload:     map[string]string{"id": client.uuid},
			}

			for client := range world.clientMap {
				client.conn.WriteJSON(message)
			}
		case message := <-world.move:
			payload := message.Payload.(map[string]interface{})
			uuid := payload["id"].(string)
			x := int(payload["x"].(float64))
			y := int(payload["y"].(float64))

			if player, ok := world.playerMap[uuid]; ok {
				player.X = x
				player.Y = y
				world.playerMap[uuid] = player

				for client := range world.clientMap {
					client.conn.WriteJSON(message)
				}
			}
		case message := <-world.attack:
			payload := message.Payload.(map[string]interface{})
			uuid := payload["id"].(string)

			if player, ok := world.playerMap[uuid]; ok {
				player.Hp = player.Hp - 1

				if player.Hp > 0 {
					world.playerMap[uuid] = player
					message.Payload.(map[string]interface{})["hp"] = player.Hp
					for client := range world.clientMap {
						client.conn.WriteJSON(message)
					}
				} else {
					for client := range world.clientMap {
						if client.uuid == uuid {
							client.conn.Close()
							delete(world.clientMap, client)
							delete(world.playerMap, client.uuid)
							continue
						}

						client.conn.WriteJSON(Message{
							MessageType: LEAVE,
							Payload:     map[string]string{"id": client.uuid},
						})
					}
				}
			}
		}
	}
}
