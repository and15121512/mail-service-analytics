package task

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"gitlab.com/sukharnikov.aa/mail-service-analytics/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-analytics/internal/ports"
	"gitlab.com/sukharnikov.aa/mail-service-analytics/internal/utils"
	"go.uber.org/zap"
)

type Task struct {
	analytics ports.Analytics
	co        Consumer
	logger    *zap.SugaredLogger
	ctx       context.Context
	cancel    context.CancelFunc
}

func New(logger *zap.SugaredLogger, analytics ports.Analytics, co Consumer) *Task {
	return &Task{
		analytics: analytics,
		co:        co,
		logger:    logger,
	}
}

func (t *Task) annotatedLogger(ctx context.Context) *zap.SugaredLogger {
	request_id, _ := ctx.Value(utils.CtxKeyRequestIDGet()).(string)
	method, _ := ctx.Value(utils.CtxKeyMethodGet()).(string)
	url, _ := ctx.Value(utils.CtxKeyURLGet()).(string)

	return t.logger.With(
		"request_id", request_id,
		"method", method,
		"url", url,
	)
}

func (t *Task) Start(ctx context.Context) error {
	t.ctx, t.cancel = context.WithCancel(ctx)
	for {
		msg, err := t.co.ReadSync(t.ctx)
		if err != nil {
			t.logger.Errorf("failed to read message: %s", err.Error())
			return fmt.Errorf("failed to read message: %s", err.Error())
		}
		t.co.ProcessAndCommitAsync(t.ctx, msg, t.StoreEvent)
	}
}

func (t *Task) Stop() {
	t.cancel()
}

func (t *Task) StoreEvent(ctx context.Context, eventMsg *kafka.Message) error {
	logger := t.annotatedLogger(ctx)

	var event models.Event
	err := json.Unmarshal(eventMsg.Value, &event)
	if err != nil {
		logger.Errorf("failed to deserialize event message from kafka: %s", err.Error())
		return fmt.Errorf("failed to deserialize event message from kafka: %s", err.Error())
	}

	err = t.analytics.StoreEvent(ctx, &event)
	if err != nil {
		logger.Errorf("failed to store event")
		return fmt.Errorf("failed to store event")
	}
	return nil
}
