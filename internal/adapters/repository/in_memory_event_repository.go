package repository

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/erwin-lovecraft/pistol/internal/core/domain"
	"github.com/erwin-lovecraft/pistol/internal/core/ports"
)

var (
	timeNowFunc     = time.Now
	cleanupSchedule = 15 * time.Minute
	defaultTTL      = time.Hour
)

var _ ports.EventRepository = (*InMemoryEventRepository)(nil)

type InMemoryEventRepository struct {
	cache sync.Map
}

func NewInMemoryEventRepository() *InMemoryEventRepository {
	repo := InMemoryEventRepository{
		cache: sync.Map{},
	}

	go func() {
		ticker := time.NewTicker(cleanupSchedule)
		for {
			select {
			case <-ticker.C:
				repo.cleanUp(defaultTTL)
			}
		}
	}()

	return &repo
}

func (i *InMemoryEventRepository) Save(ctx context.Context, roomID string, ev *domain.Event) error {
	ev.CreatedAt = timeNowFunc().UTC()

	data, ok := i.cache.Load(roomID)
	if !ok || data == nil {
		i.cache.Store(roomID, []domain.Event{*ev})
		return nil
	}

	events, ok := data.([]domain.Event)
	if !ok {
		return errors.New("invalid events type")
	}

	events = append(events, *ev)
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

func (i *InMemoryEventRepository) cleanUp(ttl time.Duration) {
	log.Printf("[event_repository] cleaning up expired events")
	now := timeNowFunc().UTC()

	i.cache.Range(func(key, value interface{}) bool {
		var retained []domain.Event
		for _, ev := range value.([]domain.Event) {
			if now.After(ev.CreatedAt.Add(ttl)) {
				continue
			}

			retained = append(retained, ev)
		}

		i.cache.Store(key, retained)
		return true
	})
}
