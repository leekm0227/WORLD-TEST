package main

import "github.com/gorilla/websocket"

type MessageType int

type Message struct {
	MessageType int64 `json:"messageType"`
	Payload     interface{}
}

type Room struct {
	join      chan *websocket.Conn
	leave     chan *websocket.Conn
	broadcast chan Message
	connMap   map[*websocket.Conn]bool
}
