package ssehub

import (
	"sync"
)

type Message struct {
	Event string
	Data  string
	ID    string
	Retry int64
}

type Hub struct {
	rooms map[string]map[string]*Client
	mu    sync.RWMutex
}

func (h *Hub) NewRoom(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.rooms[id]; !ok {
		h.rooms[id] = make(map[string]*Client)
	}
}
