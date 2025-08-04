package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/erwin-lovecraft/pistol/internal/core/domain"
	"github.com/erwin-lovecraft/pistol/internal/core/ports"
)

var _ ports.EventRepository = (*InMemoryEventRepository)(nil)

type InMemoryEventRepository struct {
	cache sync.Map
}

func NewInMemoryEventRepository() *InMemoryEventRepository {
	return &InMemoryEventRepository{
		cache: sync.Map{},
	}
}

func (i *InMemoryEventRepository) Save(ctx context.Context, roomID string, ev domain.Event) error {
	data, ok := i.cache.Load(roomID)
	if !ok || data == nil {
		i.cache.Store(roomID, []domain.Event{ev})
		return nil
	}

	events, ok := data.([]domain.Event)
	if !ok {
		return errors.New("invalid events type")
	}

	events = append(events, ev)
	i.cache.Store(roomID, events)
	return nil
}

func (i *InMemoryEventRepository) List(ctx context.Context, roomID string, page, size int) ([]domain.Event, bool, error) {
	data, ok := i.cache.Load(roomID)
	if !ok || data == nil {
		return nil, false, nil
	}

	events, ok := data.([]domain.Event)
	if !ok {
		return nil, false, errors.New("invalid data")
	}

	offset := (page - 1) * size
	hi := offset + size
	if hi > len(events) {
		hi = len(events)
	}

	var rs []domain.Event
	for idx := hi - 1; idx >= offset; idx-- {
		rs = append(rs, events[idx])
	}

	return rs, offset+size < len(events), nil
}
