package ports

import (
	"context"

	"gitlab.com/sukharnikov.aa/mail-service-analytics/internal/domain/models"
)

type Auth interface {
	ValidateAuth(ctx context.Context, tokenpair *models.TokenPair) (*models.AuthResult, error)
}
