package src

type Hub struct {
	clients map[*Client]bool
	// add, remove
	add    chan *Client
	remove chan *Client
}

func NewHub() *Hub {
	newHub := &Hub{
		clients: map[*Client]bool{},
		add:     make(chan *Client),
		remove:  make(chan *Client),
	}

	go newHub.run()

	return newHub
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.add:
			h.clients[client] = true
		case client := <-h.remove:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				client.conn.Close()
			}
		}
	}
}
