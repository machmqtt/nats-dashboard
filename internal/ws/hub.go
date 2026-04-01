package ws

import (
	"log/slog"
	"sync"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]bool
	log     *slog.Logger
}

func NewHub(log *slog.Logger) *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		log:     log,
	}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	h.clients[c] = true
	h.mu.Unlock()
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	delete(h.clients, c)
	h.mu.Unlock()
}

// Broadcast sends a message to all clients subscribed to the given environment.
func (h *Hub) Broadcast(env string, msgType string, data any) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	msg := Message{Type: msgType, Env: env, Data: data}
	for c := range h.clients {
		if c.Env() == env {
			select {
			case c.send <- msg:
			default:
				// Drop message if client is slow.
			}
		}
	}
}
