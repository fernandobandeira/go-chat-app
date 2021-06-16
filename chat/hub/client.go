package hub

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait     = 10 * time.Second
	pongWait      = 60 * time.Second
	pingPeriod    = (pongWait * 9) / 10
	maxMesageSize = 512
)

type client struct {
	hub  *hub
	conn *websocket.Conn
	send chan []byte
}

func NewClient(hub *hub, conn *websocket.Conn) *client {
	c := &client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	c.hub.register <- c
	go c.writePump()
	go c.readPump()

	return c
}

func (c *client) readPump() {
	defer func() {
		c.hub.unregister <- c
		err := c.conn.Close()
		if err != nil {
			log.Printf("error on close connection (readPump): %v", err)
		}
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(message)
		c.hub.broadcast <- message
	}
}

func (c *client) writePump() {
	for {
		message, ok := <-c.send
		if !ok {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)

		if err := w.Close(); err != nil {
			return
		}
	}
}
