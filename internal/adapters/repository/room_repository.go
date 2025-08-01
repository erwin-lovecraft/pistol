package repository

import (
	"context"
	"sync"

	"github.com/erwin-lovecraft/pistol/internal/core/domain"
	"github.com/erwin-lovecraft/pistol/internal/core/ports"
)

var _ ports.RoomRepository = (*InMemoryRoomRepository)(nil)

type InMemoryRoomRepository struct {
	cache map[string]domain.Room
	mu    sync.RWMutex
}

func NewInMemoryRoomRepository() *InMemoryRoomRepository {
	return &InMemoryRoomRepository{
		cache: make(map[string]domain.Room),
	}
}

func (i *InMemoryRoomRepository) SaveRoom(ctx context.Context, room domain.Room) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.cache[room.ID] = room
	return nil
}

func (i *InMemoryRoomRepository) List(ctx context.Context, filter ports.RoomFilter) ([]domain.Room, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	var rooms []domain.Room
	for _, room := range i.cache {
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (i *InMemoryRoomRepository) DeleteByID(ctx context.Context, id string) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	delete(i.cache, id)
	return nil
}
