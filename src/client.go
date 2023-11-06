package src

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var ugrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	uuid string
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) read() {
	defer func() {
		c.hub.remove <- c
	}()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			c.hub.remove <- c
			break
		}

		switch string(msg) {
		case "get_files":
			// TODO: handle errors
			list, _ := GetFiles()
			names := make([]string, len(list))

			for i, item := range list {
				names[i] = item.Name()
			}
			b, _ := json.Marshal(names)

			// could implement broadcasting here
			c.send <- b
		}
	}
}

func (c *Client) write() {
	defer func() {
		c.conn.Close()
	}()
	for {
		message, ok := <-c.send
		fmt.Println(string(message))
		if !ok {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)

		n := len(c.send)
		for i := 0; i < n; i++ {
			w.Write([]byte("\n"))
			w.Write(<-c.send)
		}

		if err := w.Close(); err != nil {
			return
		}
	}
}

func NewClient(h *Hub, id string, w http.ResponseWriter, r *http.Request) error {
	conn, err := ugrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	c := &Client{
		uuid: id,
		hub:  h,
		conn: conn,
		send: make(chan []byte),
	}
	h.add <- c

	go c.read()
	go c.write()

	return nil
}
