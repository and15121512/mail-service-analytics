package analytics

import (
	"context"
	"fmt"
	"sort"
	"time"

	"gitlab.com/sukharnikov.aa/mail-service-analytics/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-analytics/internal/ports"
	"gitlab.com/sukharnikov.aa/mail-service-analytics/internal/utils"
	"go.uber.org/zap"
)

type Service struct {
	es     ports.EventStorage
	ac     ports.Auth
	logger *zap.SugaredLogger
}

func New(es ports.EventStorage, ac ports.Auth, logger *zap.SugaredLogger) *Service {
	return &Service{
		es:     es,
		ac:     ac,
		logger: logger,
	}
}

func (s *Service) annotatedLogger(ctx context.Context) *zap.SugaredLogger {
	request_id, _ := ctx.Value(utils.CtxKeyRequestIDGet()).(string)
	method, _ := ctx.Value(utils.CtxKeyMethodGet()).(string)
	url, _ := ctx.Value(utils.CtxKeyURLGet()).(string)

	return s.logger.With(
		"request_id", request_id,
		"method", method,
		"url", url,
	)
}

func (s *Service) StoreEvent(ctx context.Context, event *models.Event) error {
	return s.es.CreateEvent(ctx, event)
}

func (s *Service) GetReport(ctx context.Context, taskId string) (*models.Report, error) {
	logger := s.annotatedLogger(ctx)

	doneCnt, err := s.es.CountDoneEvents(ctx)
	if err != nil {
		logger.Errorf("failed to count done events")
		return &models.Report{}, fmt.Errorf("failed to count done events")
	}
	declinedCnt, err := s.es.CountDeclinedEvents(ctx)
	if err != nil {
		logger.Errorf("failed to count declined events")
		return &models.Report{}, fmt.Errorf("failed to count declined events")
	}

	events, err := s.es.GetEventsByTaskID(ctx, taskId)
	if err != nil {
		logger.Errorf("failed to get all events for task ID %s", taskId)
		return &models.Report{}, fmt.Errorf("failed to get all events for task ID %s", taskId)
	}
	reactionDurations := getReactionDurations(ctx, events)

	return &models.Report{
		DoneCnt:           doneCnt,
		DeclinedCnt:       declinedCnt,
		TaskId:            taskId,
		ReactionDurations: reactionDurations,
	}, nil
}

func getReactionDurations(ctx context.Context, events []models.Event) []time.Duration {
	sort.Slice(events, func(i int, j int) bool {
		return events[i].Time.Before(events[j].Time)
	})

	reactionDurations := make([]time.Duration, len(events)-1)
	for i := range reactionDurations {
		reactionDurations[i] = events[i+1].Time.Sub(events[i].Time)
	}
	return reactionDurations
}
