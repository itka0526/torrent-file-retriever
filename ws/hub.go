package ws

import (
	"encoding/json"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// slice header is being copied
// even if the original slice is modified it will reflected in the copy because the underlying array is the same
// func passByValue() []string {
// 	v := []string{"apples", "bananas"}
// 	return v
// }

// this is for if you want to the caller to modify the original slice's content.
// func passByReference() *[]string {
// 	v := []string{"apples", "bananas"}
// 	return &v
// }

func spewFileInfos(files []MyFileInfo) []byte {
	b, err := json.Marshal(files)
	if err != nil {
		return nil
	}
	return b
}

func (h *Hub) Run(files []MyFileInfo) {
	for {
		// Feed the latest information about downloads
		go func() {
			for {
				h.broadcast <- spewFileInfos(files)
				time.Sleep(500 * time.Millisecond)
			}
		}()
		// for client := range h.clients {
		// 	for {
		// 		client.send <- spewFileInfos(files)
		// 		time.Sleep(500 * time.Millisecond)
		// 	}
		// }

		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
