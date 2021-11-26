package main

import "github.com/gorilla/websocket"

type MessageType int

const (
	MOVE MessageType = 0
)

type Message struct {
	MessageType MessageType `json:"messageType"`
	Payload     interface{} `json:"payload"`
}

type Channel struct {
	join      chan *Client
	leave     chan *Client
	broadcast chan Message
	clients   map[*Client]bool
}

type Client struct {
	uuid   string
	conn   *websocket.Conn
	action chan Message
	send   chan Message
}
