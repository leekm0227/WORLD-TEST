package main

import (
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func newClient(conn *websocket.Conn) *Client {
	client := &Client{
		uuid:   strings.Split(uuid.NewString(), "-")[0],
		conn:   conn,
		action: make(chan Message),
		send:   make(chan Message),
	}

	go func() {
		for {
			select {
			case message := <-client.action:
				println(message)
			case message := <-client.send:
				client.conn.WriteJSON(message)
			}
		}
	}()

	return client
}
