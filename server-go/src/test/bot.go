package test

import (
	"AAA/src/server/world"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

var test Test = Test{
	uids: make([]string, 0),
}

func Run(port string, size int) {
	for i := 0; i < size; i++ {
		go newBot().run(port)
	}
}

func (test *Test) addUid(uid string) {
	test.Lock()
	defer test.Unlock()
	test.uids = append(test.uids, uid)
}

func (test *Test) removeUid(target string) {
	test.Lock()
	defer test.Unlock()
	for i, uid := range test.uids {
		if uid == target {
			test.uids = append(test.uids[:i], test.uids[i+1:]...)
		}
	}
}

func (test *Test) getRandomUid() string {
	test.Lock()
	defer test.Unlock()

	if len(test.uids) > 0 {
		return test.uids[rand.Intn(len(test.uids))]
	}

	return ""
}

func newBot() *Bot {
	return &Bot{
		receive: make(chan world.Message, 100),
	}
}

func (bot *Bot) run(port string) {
	for {
		var err error
		bot.conn, _, err = websocket.DefaultDialer.Dial("ws://localhost:"+port+"/world", nil)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}

		break
	}

	defer bot.close()
	go bot.handler()

	for {
		var message world.Message
		if err := bot.conn.ReadJSON(&message); err != nil {
			return
		}

		bot.receive <- message
	}

}

func (bot *Bot) close() {
	bot.conn.Close()
	test.removeUid(bot.player.Id)
}

func (bot *Bot) handler() {
	for {
		select {
		case message := <-bot.receive:
			switch message.MessageType {
			case world.INIT:
				data, err := json.Marshal(message.Payload.(map[string]interface{})["player"])
				if err == nil {
					var player world.Player
					json.Unmarshal(data, &player)
					bot.player = &player
					test.addUid(player.Id)
				}
			}
		default:
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)+500))
			if bot.player == nil {
				continue
			}

			switch ActionType(rand.Intn(int(MAX))) {
			case ATTACK:
				target := test.getRandomUid()
				if target != "" && target != bot.player.Id { // 자살방지
					bot.conn.WriteJSON(map[string]interface{}{
						"messageType": world.ATTACK,
						"payload": map[string]interface{}{
							"id": target,
							"tx": time.Now().UnixMilli(),
						},
					})
				}
			case MOVE:
				if bot.player != nil {
					switch world.Direction(rand.Intn(4)) {
					case world.UP:
						bot.player.Y--
						if bot.player.Y < world.Y_MIN {
							bot.player.Y = world.Y_MIN
						}
					case world.DOWN:
						bot.player.Y++
						if bot.player.Y > world.Y_MAX {
							bot.player.Y = world.Y_MAX
						}
					case world.LEFT:
						bot.player.X--
						if bot.player.X < world.X_MIN {
							bot.player.X = world.X_MIN
						}
					case world.RIGHT:
						bot.player.X++
						if bot.player.X > world.X_MAX {
							bot.player.X = world.X_MAX
						}
					}

					bot.conn.WriteJSON(map[string]interface{}{
						"messageType": world.MOVE,
						"payload": map[string]interface{}{
							"id": bot.player.Id,
							"x":  bot.player.X,
							"y":  bot.player.Y,
							"tx": time.Now().UnixMilli(),
						},
					})
				}
			}
		}
	}
}
