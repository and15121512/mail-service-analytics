package kafka_consumer

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaConsumer struct {
	reader *kafka.Reader
	logger *zap.SugaredLogger
}

func NewConsumer(brokers []string, topic string, groupID string, logger *zap.SugaredLogger) (*KafkaConsumer, error) {
	if len(brokers) == 0 || brokers[0] == "" || topic == "" || groupID == "" {
		logger.Errorf("invalid config parameters for kafka consumer")
		return nil, fmt.Errorf("invalid config parameters for kafka consumer")
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	return &KafkaConsumer{
		reader: reader,
		logger: logger,
	}, nil
}

func (kc *KafkaConsumer) ReadSync(ctx context.Context) (*kafka.Message, error) {
	kafkaMsg, err := kc.reader.FetchMessage(ctx)
	if err != nil {
		kc.logger.Errorf("failed to fetch message from kafka: %s", err.Error())
		return nil, fmt.Errorf("failed to fetch message from kafka: %s", err.Error())
	}
	return &kafkaMsg, nil
}

func (kc *KafkaConsumer) ProcessAndCommitAsync(ctx context.Context, m *kafka.Message, fn func(context.Context, *kafka.Message) error) {
	go func() {
		if fn != nil {
			if err := fn(ctx, m); err != nil {
				kc.logger.Errorf("failed to process message from kafka: %s", err.Error())
				return
			}
		}

		err := kc.reader.CommitMessages(ctx, *m)
		if err != nil {
			kc.logger.Errorf("failed to commit message from kafka: %s", err.Error())
		}
	}()
}
