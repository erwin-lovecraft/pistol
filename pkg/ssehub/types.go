package ssehub

import (
	"sync"
)

type Event struct {
	Event   string
	Payload string
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
