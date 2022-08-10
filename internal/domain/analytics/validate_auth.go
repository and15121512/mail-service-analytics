package analytics

import (
	"context"

	"gitlab.com/sukharnikov.aa/mail-service-analytics/internal/domain/models"
)

func (s *Service) ValidateAuth(ctx context.Context, tokenpair *models.TokenPair) (*models.AuthResult, error) {
	return s.ac.ValidateAuth(ctx, tokenpair)
}
