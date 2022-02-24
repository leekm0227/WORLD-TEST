package test

import (
	"AAA/src/server"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

var uids []string = make([]string, 0)

func Run(port string, size int) {
	for i := 0; i < size; i++ {
		go run(port)
	}
}

func run(port string) {
	var dummy Dummy

	for {
		var err error
		dummy.conn, _, err = websocket.DefaultDialer.Dial("ws://localhost:"+port+"/ws", nil)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}

		break
	}

	go dummy.receive()
	defer dummy.conn.Close()

	for {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)+500))

		switch ActionType(rand.Intn(int(MAX))) {
		case ATTACK:
			if len(uids) > 0 {
				uid := uids[server.Direction(rand.Intn(len(uids)))]

				if uid != dummy.player.Id { // 자살방지
					dummy.conn.WriteJSON(map[string]interface{}{
						"messageType": server.ATTACK,
						"payload": map[string]interface{}{
							"id": uid,
						},
					})
				}
			}
		case MOVE:
			if dummy.player != nil {
				switch server.Direction(rand.Intn(4)) {
				case server.UP:
					dummy.player.Y--
					if dummy.player.Y < server.Y_MIN {
						dummy.player.Y = server.Y_MIN
					}
				case server.DOWN:
					dummy.player.Y++
					if dummy.player.Y > server.Y_MAX {
						dummy.player.Y = server.Y_MAX
					}
				case server.LEFT:
					dummy.player.X--
					if dummy.player.X < server.X_MIN {
						dummy.player.X = server.X_MIN
					}
				case server.RIGHT:
					dummy.player.X++
					if dummy.player.X > server.X_MAX {
						dummy.player.X = server.X_MAX
					}
				}

				dummy.conn.WriteJSON(map[string]interface{}{
					"messageType": server.MOVE,
					"payload": map[string]interface{}{
						"id": dummy.player.Id,
						"x":  dummy.player.X,
						"y":  dummy.player.Y,
					},
				})
			}
		}
	}
}

func (dummy *Dummy) receive() {
	var res server.Message

	for {
		if err := dummy.conn.ReadJSON(&res); err != nil {
			for i, uid := range uids {
				if uid == dummy.player.Id {
					uids = append(uids[:i], uids[i+1:]...)
				}
			}

			log.Printf("connection closed, left uids: %d", len(uids))
			return
		}

		switch res.MessageType {
		case server.INIT:
			data, err := json.Marshal(res.Payload.(map[string]interface{})["player"])
			if err == nil {
				var player server.Player
				json.Unmarshal(data, &player)
				uids = append(uids, player.Id)
				dummy.player = &player
			}
		}
	}
}
