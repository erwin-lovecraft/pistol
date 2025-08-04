package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/erwin-lovecraft/pistol/internal/adapters/ormmodel"
	"github.com/erwin-lovecraft/pistol/internal/core/domain"
	"github.com/erwin-lovecraft/pistol/internal/core/ports"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ ports.EventRepository = (*eventRepository)(nil)

type eventRepository struct {
	queries *ormmodel.Queries
}

func NewEventRepository(dbPool *pgxpool.Pool) ports.EventRepository {
	return eventRepository{
		queries: ormmodel.New(dbPool),
	}
}

func (repo eventRepository) Save(ctx context.Context, roomID string, ev *domain.Event) error {
	if ev.ID == 0 {
		id, err := sf.NextID()
		if err != nil {
			return fmt.Errorf("generate id: %w", err)
		}
		ev.ID = id
	}

	var (
		err             error
		headerBytes     []byte
		queryParamBytes []byte
	)
	if ev.Header != nil {
		if headerBytes, err = json.Marshal(ev.Header); err != nil {
			return fmt.Errorf("marshal header: %w", err)
		}
	}
	if ev.QueryParams != nil {
		if queryParamBytes, err = json.Marshal(ev.QueryParams); err != nil {
			return fmt.Errorf("marshal query param: %w", err)
		}
	}

	var pgRoomID pgtype.UUID
	if err := pgRoomID.Scan(roomID); err != nil {
		return fmt.Errorf("scan room id: %w", err)
	}

	createdAt, err := repo.queries.SaveEvent(ctx, ormmodel.SaveEventParams{
		ID:          ev.ID,
		Method:      ev.Method,
		Header:      headerBytes,
		QueryParams: queryParamBytes,
		Body:        ev.Body,
		RoomID:      pgRoomID,
	})
	if err != nil {
		return fmt.Errorf("save event: %w", err)
	}
	ev.CreatedAt = createdAt.Time
	return nil
}

func (repo eventRepository) List(ctx context.Context, roomID string, page, size int) ([]domain.Event, bool, error) {
	offset := (page - 1) * size
	limit := size

	var pgRoomID pgtype.UUID
	if err := pgRoomID.Scan(roomID); err != nil {
		return nil, false, fmt.Errorf("scan room id: %w", err)
	}

	models, err := repo.queries.ListEvents(ctx, ormmodel.ListEventsParams{
		RoomID: pgRoomID,
		Offset: int32(offset),
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, false, fmt.Errorf("list events: %w", err)
	}

	events := make([]domain.Event, len(models))
	for idx, model := range models {
		var (
			evHeader  http.Header
			evQueries map[string][]string
		)
		if err := json.Unmarshal(model.Header, &evHeader); err != nil {
			log.Printf("unmarshal header: %v", err)
			continue
		}
		if err := json.Unmarshal(model.QueryParams, &evQueries); err != nil {
			log.Printf("unmarshal query params: %v", err)
			continue
		}

		events[idx] = domain.Event{
			ID:          model.ID,
			Method:      model.Method,
			Body:        model.Body,
			Header:      evHeader,
			QueryParams: evQueries,
			CreatedAt:   model.CreatedAt.Time,
		}
	}

	return events, true, nil
}
