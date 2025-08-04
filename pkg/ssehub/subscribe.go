package ssehub

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

const (
	sendBuffer = 64
)

func (h *Hub) Subscribe(ctx context.Context, room string, clientID string, w http.ResponseWriter) (*Client, error) {
	// Setup SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, errors.New("streaming is not supported")
	}

	ctx, cancel := context.WithCancel(ctx)
	client := &Client{
		id:         clientID,
		room:       room,
		ctx:        ctx,
		cancel:     cancel,
		sendCh:     make(chan Message, sendBuffer), // buffered to absorb burst
		connected:  time.Now(),
		lastActive: time.Now(),
	}

	// Register
	h.mu.Lock()
	if _, exists := h.rooms[room]; !exists {
		h.rooms[room] = make(map[string]*Client)
	}
	h.rooms[room][clientID] = client
	h.mu.Unlock()

	log.Printf("[SSE] Subcribed %s, rooms size: %d", client, h.RoomConnections(room))

	// Start writer goroutine
	go client.writerLoop(w, flusher)

	// Start a heartbeat to keep connection alive / detect silent drop
	go client.heartbeat()

	// Wait for context done to unregister
	go func() {
		<-ctx.Done()
		h.unregister(room, clientID)
	}()

	return client, nil
}

func (h *Hub) unregister(room string, clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clients, exists := h.rooms[room]
	if !exists {
		return
	}

	if _, exists := clients[clientID]; exists {
		delete(clients, clientID)
	}

	if len(clients) == 0 {
		delete(h.rooms, room)
	}
}
