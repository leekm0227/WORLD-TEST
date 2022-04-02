package test

import (
	"AAA/src/server/world"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	MOVE ActionType = iota
	ATTACK
	MAX
)

type Bot struct {
	conn    *websocket.Conn
	player  *world.Player
	receive chan world.Message
}

type Test struct {
	sync.Mutex
	uids []string
}

type ActionType int
