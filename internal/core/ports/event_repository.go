package ports

import (
	"context"

	"github.com/erwin-lovecraft/pistol/internal/core/domain"
)

type EventRepository interface {
	Save(ctx context.Context, roomID string, ev *domain.Event) error

	List(ctx context.Context, roomID string, page, size int) ([]domain.Event, bool, error)
}
