package main

func newChannel() *Channel {
	channel := &Channel{
		join:      make(chan *Client),
		leave:     make(chan *Client),
		broadcast: make(chan Message),
		clients:   make(map[*Client]bool),
	}

	go func() {
		for {
			select {
			case client := <-channel.join:
				channel.clients[client] = true
			case client := <-channel.leave:
				delete(channel.clients, client)
			case message := <-channel.broadcast:
				for client := range channel.clients {
					client.send <- message
				}
			}
		}
	}()

	return channel
}

func leave(client *Client) {
	channel.leave <- client
	// TODO: broadcast leave message
}
