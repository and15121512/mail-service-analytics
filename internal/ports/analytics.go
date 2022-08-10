package ports

import (
	"context"

	"gitlab.com/sukharnikov.aa/mail-service-analytics/internal/domain/models"
)

type Analytics interface {
	StoreEvent(context.Context, *models.Event) error
	GetReport(context.Context, string) (*models.Report, error)

	ValidateAuth(context.Context, *models.TokenPair) (*models.AuthResult, error)
}
