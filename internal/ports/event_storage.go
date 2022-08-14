package ports

import (
	"context"

	"gitlab.com/sukharnikov.aa/mail-service-analytics/internal/domain/models"
)

type EventStorage interface {
	CreateEventIfNotExists(ctx context.Context, event *models.Event) error
	CountDoneEvents(ctx context.Context) (int, error)
	CountDeclinedEvents(ctx context.Context) (int, error)
	GetEventsByTaskID(ctx context.Context, taskId string) ([]models.Event, error)
}
