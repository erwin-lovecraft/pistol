package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/erwin-lovecraft/pistol/internal/core/domain"
	"github.com/erwin-lovecraft/pistol/internal/core/ports"
	"github.com/erwin-lovecraft/pistol/pkg/ssehub"
	"github.com/google/uuid"
)

var (
	uuidFunc = uuid.New
)

type Service interface {
	ListenEvents(ctx context.Context, roomID string, w http.ResponseWriter) (*ssehub.Client, error)

	PushEvent(ctx context.Context, roomID string, event domain.Event) error

	ListEvents(ctx context.Context, roomID string, page int, size int) ([]domain.Event, bool, error)
}

type service struct {
	hub             *ssehub.Hub
	roomRepository  ports.RoomRepository
	eventRepository ports.EventRepository
}

func NewService(roomRepository ports.RoomRepository, eventRepository ports.EventRepository) Service {
	return &service{
		hub:             ssehub.NewHub(),
		eventRepository: eventRepository,
		roomRepository:  roomRepository,
	}
}

func (s *service) CreateRoom(ctx context.Context, name, avatar string) (domain.Room, string, error) {
	id := uuidFunc()
	s.hub.NewRoom(id.String())

	room := domain.Room{
		ID:     id.String(),
		Name:   name,
		Avatar: avatar,
	}
	if err := s.roomRepository.SaveRoom(ctx, room); err != nil {
		return domain.Room{}, "", err
	}

	return room, buildRoomLink(room), nil
}

func (s *service) ListRoom(ctx context.Context) ([]domain.Room, error) {
	return s.roomRepository.List(ctx, ports.RoomFilter{})
}

func buildRoomLink(room domain.Room) string {
	return fmt.Sprintf("/rooms/%s/events", room.ID)
}

func (s *service) ListenEvents(ctx context.Context, roomID string, w http.ResponseWriter) (*ssehub.Client, error) {
	clientID := uuidFunc()

	cl, err := s.hub.Subscribe(ctx, roomID, clientID.String(), w)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to client: %w", err)
	}

	return cl, nil
}

func (s *service) PushEvent(ctx context.Context, roomID string, event domain.Event) error {
	// Sanitize headers
	for k := range event.Header {
		if slices.Contains(secretHeaders, k) {
			event.Header.Del(k)
		}
	}
	for k := range event.QueryParams {
		if slices.Contains(secretQueryParams, k) {
			delete(event.QueryParams, k)
		}
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Persist event
	if err := s.eventRepository.Save(ctx, roomID, &event); err != nil {
		return fmt.Errorf("failed to save event: %w", err)
	}

	return s.hub.SendToRoom(roomID, ssehub.Message{
		Event: "message",
		Data:  string(payload),
	})
}

func (s *service) ListEvents(ctx context.Context, roomID string, page int, size int) ([]domain.Event, bool, error) {
	rs, hasMore, err := s.eventRepository.List(ctx, roomID, page, size)
	if err != nil {
		return nil, false, err
	}

	return rs, hasMore, nil
}
