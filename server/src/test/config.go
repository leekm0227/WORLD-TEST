package test

import (
	"AAA/src/server"

	"github.com/gorilla/websocket"
)

const (
	MOVE ActionType = iota
	ATTACK
	MAX
)

type Dummy struct {
	conn   *websocket.Conn
	player *server.Player
}

type ActionType int
