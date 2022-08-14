package task

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Consumer interface {
	ReadSync(ctx context.Context) (*kafka.Message, error)
	ProcessAndCommitAsync(ctx context.Context, m *kafka.Message, fn func(context.Context, *kafka.Message) error)
}
