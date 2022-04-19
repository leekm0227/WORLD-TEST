package world

import "github.com/gorilla/websocket"

type MessageType int
type Direction int

const (
	JOIN   MessageType = 0
	LEAVE  MessageType = 1
	INIT   MessageType = 2
	MOVE   MessageType = 3
	ATTACK MessageType = 4

	UP    Direction = 0
	DOWN  Direction = 1
	LEFT  Direction = 2
	RIGHT Direction = 3
)

const (
	MAX_SIZE         = 100
	CHAN_BUFFER_SIZE = 1000
	X_MIN            = 0
	X_MAX            = 20
	Y_MIN            = 0
	Y_MAX            = 20
)

type Message struct {
	MessageType MessageType `json:"messageType"`
	Payload     interface{} `json:"payload"`
}

type World struct {
	wid       string
	join      chan *Client
	leave     chan *Client
	move      chan Message
	attack    chan Message
	clientMap map[*Client]bool
	playerMap map[string]Player
}

type Player struct {
	Id string `json:"id"`
	Hp int    `json:"hp"`
	X  int    `json:"x"`
	Y  int    `json:"y"`
}

type Client struct {
	uuid string
	conn *websocket.Conn
}
