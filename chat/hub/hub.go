package hub

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/fernandomalmeida/go-chat-app/dbgen/dbchat"
)

type hub struct {
	store      dbchat.Store
	clients    map[*client]bool
	broadcast  chan []byte
	register   chan *client
	unregister chan *client
}

func New(chatStore dbchat.Store) *hub {
	return &hub{
		store:      chatStore,
		clients:    make(map[*client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
	}
}

func (h *hub) Run() {
	for {
		select {
		case client := <-h.register:
			log.Print("client registered!")
			h.clients[client] = true
			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
			defer cancel()
			msgs, err := h.store.GetMessages(ctx)
			if err != nil {
				log.Printf("error on get messages: %v", err)
			}
			for _, msg := range msgs {
				msgJson, _ := json.Marshal(msg)
				client.send <- msgJson
			}
		case client := <-h.unregister:
			log.Print("client unregisterd...")
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			log.Printf("broadcast message: %d", message)
			var msg dbchat.AddMessageParams
			json.Unmarshal(message, &msg)
			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
			defer cancel()
			msgAdded, err := h.store.AddMessage(ctx, msg)
			if err != nil {
				log.Printf("error on add message: %v", err)
			}
			msgBytes, _ := json.Marshal(msgAdded)
			for client := range h.clients {
				select {
				case client.send <- msgBytes:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
