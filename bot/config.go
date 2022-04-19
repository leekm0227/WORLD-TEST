package main

import (
	"sync"
	"github.com/gorilla/websocket"
)

const (
	ACT_MOVE ActionType = iota
	ACT_ATTACK
	ACT_MAX
)

type Bot struct {
	conn    *websocket.Conn
	player  *Player
	receive chan Message
}

type Test struct {
	sync.Mutex
	uids []string
}

type ActionType int

type Message struct {
	MessageType MessageType `json:"messageType"`
	Payload     interface{} `json:"payload"`
}

type Player struct {
	Id string `json:"id"`
	Hp int    `json:"hp"`
	X  int    `json:"x"`
	Y  int    `json:"y"`
}

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

	X_MIN            = 0
	X_MAX            = 20
	Y_MIN            = 0
	Y_MAX            = 20
)