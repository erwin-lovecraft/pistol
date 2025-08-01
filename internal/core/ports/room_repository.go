package ports

import (
	"context"

	"github.com/erwin-lovecraft/pistol/internal/core/domain"
)

type RoomRepository interface {
	SaveRoom(ctx context.Context, room domain.Room) error

	List(ctx context.Context, filter RoomFilter) ([]domain.Room, error)

	DeleteByID(ctx context.Context, id string) error
}

type RoomFilter struct{}
