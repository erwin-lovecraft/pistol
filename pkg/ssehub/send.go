package ssehub

import (
	"errors"
)

func (h *Hub) SendToRoom(room string, e Message) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, exists := h.rooms[room]
	if !exists {
		return errors.New("room not found")
	}

	for _, cl := range clients {
		select {
		case cl.sendCh <- e:
		default:
			return errors.New("send channel is full")
		}
	}
	return nil
}

func (h *Hub) Broadcast(e Message) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, clients := range h.rooms {
		for _, cl := range clients {
			select {
			case cl.sendCh <- e:
			default:
				return errors.New("send channel is full")
			}
		}
	}
	return nil
}
