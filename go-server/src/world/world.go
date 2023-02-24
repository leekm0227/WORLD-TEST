package world

import (
	"log"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var world = World{
	wid:       "world-0",
	join:      make(chan *Client, CHAN_BUFFER_SIZE),
	leave:     make(chan *Client, CHAN_BUFFER_SIZE),
	move:      make(chan Message, CHAN_BUFFER_SIZE),
	attack:    make(chan Message, CHAN_BUFFER_SIZE),
	clientMap: make(map[*Client]bool, MAX_SIZE),
	playerMap: make(map[string]Player, MAX_SIZE),
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(c *gin.Context) {
	var req Message
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("upgrader.Upgrade: %+v", err)
		return
	}

	client := &Client{
		uuid: strings.Split(uuid.NewString(), "-")[0],
		conn: conn,
	}

	world.join <- client
	defer func(client *Client) {
		world.leave <- client
		client.conn.Close()
	}(client)

	for {
		if err := conn.ReadJSON(&req); err != nil {
			log.Printf("read err: %s", err)
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

func Run() {
	for {
		select {
		case client := <-world.join:
			join(client)
		case client := <-world.leave:
			leave(client)
		case message := <-world.move:
			move(message)
		case message := <-world.attack:
			attack(message)
		}
	}
}

func join(client *Client) {
	player := Player{
		Id: client.uuid,
		Hp: 5,
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
}

func leave(client *Client) {
	delete(world.clientMap, client)
	delete(world.playerMap, client.uuid)

	message := Message{
		MessageType: LEAVE,
		Payload:     map[string]string{"id": client.uuid},
	}

	for client := range world.clientMap {
		client.conn.WriteJSON(message)
	}
}

func move(message Message) {
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

		// logElapsed("move", payload)
	}
}

func attack(message Message) {
	payload := message.Payload.(map[string]interface{})
	uuid := payload["id"].(string)

	if player, ok := world.playerMap[uuid]; ok {
		player.Hp = player.Hp - 1

		message.Payload.(map[string]interface{})["hp"] = player.Hp
		for client := range world.clientMap {
			client.conn.WriteJSON(message)
		}

		// if player.Hp > 0 {
		// 	world.playerMap[uuid] = player
		// 	message.Payload.(map[string]interface{})["hp"] = player.Hp
		// 	for client := range world.clientMap {
		// 		client.conn.WriteJSON(message)
		// 	}

		// 	// logElapsed("attack", payload)
		// } else {
		// 	delete(world.playerMap, uuid)
		// 	for client := range world.clientMap {
		// 		client.conn.WriteJSON(Message{
		// 			MessageType: DIE,
		// 			Payload:     map[string]string{"id": player.Id},
		// 		})
		// 	}

		// 	// logElapsed("died", payload)
		// }
	}
}

func logElapsed(key string, payload map[string]interface{}) {
	if payload["tx"] != nil {
		elapsed := math.Round(float64(time.Now().UnixMilli()) - payload["tx"].(float64))
		delete(payload, "tx")
		log.Printf("[elapsed: %0.0fms, %s]\tpayload: %+v", elapsed, key, payload)
	}
}
