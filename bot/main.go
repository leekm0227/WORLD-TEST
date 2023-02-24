package main

import (
	"encoding/json"
	"flag"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var wg sync.WaitGroup

// go run . -test=10
func main() {
	port := flag.String("port", "8888", "port number")
	dummySize := flag.Int("test", 0, "dummy size, max: 500")
	flag.Parse()

	if *dummySize > 500 {
		*dummySize = 500
	}

	if *dummySize > 0 {
		Run(*port, *dummySize)
	}
}

var test Test = Test{
	uids: make([]string, 0),
}

func Run(port string, size int) {
	wg.Add(size)
	for i := 0; i < size; i++ {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
		go newBot().run(port)
	}
	wg.Wait()
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
		playable: true,
		receive:  make(chan Message, 100),
	}
}

func (bot *Bot) run(port string) {
	for {
		var err error
		bot.conn, _, err = websocket.DefaultDialer.Dial("ws://localhost:"+port+"/world", nil)
		if err != nil {
			log.Println("conn fail")
			time.Sleep(2 * time.Second)
			continue
		}

		break
	}

	defer func() {
		wg.Done()
		bot.close()
	}()
	go bot.action()
	go bot.handler()

	for {
		var message Message
		if err := bot.conn.ReadJSON(&message); err != nil {
			log.Printf("read err: %s", err)
			return
		}

		bot.receive <- message
	}

}

func (bot *Bot) close() {
	bot.conn.Close()
	test.removeUid(bot.player.Id)
}

func (bot *Bot) action() {
	for {
		if bot.player == nil || !bot.playable {
			continue
		}
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)+500))

		switch ActionType(rand.Intn(int(ACT_MAX))) {
		case ACT_ATTACK:
			target := test.getRandomUid()
			if target != "" && target != bot.player.Id { // 자살방지
				bot.conn.WriteJSON(map[string]interface{}{
					"messageType": ATTACK,
					"payload": map[string]interface{}{
						"id": target,
						"tx": time.Now().UnixMilli(),
					},
				})

				// log.Printf("%s attacked -> %s", bot.player.Id, target)
			}
		case ACT_MOVE:
			if bot.player != nil {
				switch Direction(rand.Intn(4)) {
				case UP:
					bot.player.Y--
					if bot.player.Y < Y_MIN {
						bot.player.Y = Y_MIN
					}
				case DOWN:
					bot.player.Y++
					if bot.player.Y > Y_MAX {
						bot.player.Y = Y_MAX
					}
				case LEFT:
					bot.player.X--
					if bot.player.X < X_MIN {
						bot.player.X = X_MIN
					}
				case RIGHT:
					bot.player.X++
					if bot.player.X > X_MAX {
						bot.player.X = X_MAX
					}
				}

				bot.conn.WriteJSON(map[string]interface{}{
					"messageType": MOVE,
					"payload": map[string]interface{}{
						"id": bot.player.Id,
						"x":  bot.player.X,
						"y":  bot.player.Y,
						"tx": time.Now().UnixMilli(),
					},
				})

				// log.Printf("%s, moved", bot.player.Id)
			}
		}
	}
}

func (bot *Bot) handler() {
	for {
		message := <-bot.receive
		switch message.MessageType {
		case INIT:
			data, err := json.Marshal(message.Payload.(map[string]interface{})["player"])
			if err == nil {
				var player Player
				json.Unmarshal(data, &player)
				bot.player = &player
				test.addUid(player.Id)
			}
		case DIE:
			data, err := json.Marshal(message.Payload.(map[string]interface{})["id"])
			if err == nil {
				var uid string
				json.Unmarshal(data, &uid)
				if bot.player.Id == uid {
					bot.playable = false
					// log.Printf("died: %s", uid)
				}
			}
		default:
			data, err := json.Marshal(message.Payload.(map[string]interface{})["tx"])
			if err == nil {
				var tx float64
				json.Unmarshal(data, &tx)
				if tx > 0 {
					elapsed := math.Round(float64(time.Now().UnixMilli()) - tx)
					log.Printf("[elapsed: %0.0fms, %d]", elapsed, message.MessageType)
					// log.Printf("[elapsed: %0.0fms, %d]\tpayload: %+v", elapsed, message.MessageType, nil)
				}
			}
		}
	}
}
