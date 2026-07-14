package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// MessageHandler is a function that processes incoming messages.
type MessageHandler func(ctx context.Context, data []byte) error

// Subscriber handles subscribing to NATS JetStream subjects.
type Subscriber struct {
	nc     *nats.Conn
	js     jetstream.JetStream
	logger *slog.Logger
}

// NewSubscriber creates a new NATS subscriber.
func NewSubscriber(natsURL string, logger *slog.Logger) (*Subscriber, error) {
	nc, err := nats.Connect(natsURL,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(2*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	return &Subscriber{
		nc:     nc,
		js:     js,
		logger: logger,
	}, nil
}

// Subscribe subscribes to a subject with a durable consumer.
func (s *Subscriber) Subscribe(ctx context.Context, subject, consumerName string, handler MessageHandler) error {
	// Create or get existing consumer
	consumer, err := s.js.CreateOrUpdateConsumer(ctx, "cloudcommerce", jetstream.ConsumerConfig{
		Name:          consumerName,
		Durable:       consumerName,
		FilterSubject: subject,
		AckPolicy:     jetstream.AckExplicitPolicy,
		MaxDeliver:    3,
		AckWait:       30 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	// Start consuming messages
	_, err = consumer.Consume(func(msg jetstream.Msg) {
		// Process message
		if err := handler(ctx, msg.Data()); err != nil {
			s.logger.Error("failed to process message",
				"subject", subject,
				"error", err,
			)
			// Negative acknowledgment - will be redelivered
			msg.Nak()
			return
		}

		// Acknowledge successful processing
		msg.Ack()
	})

	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	s.logger.Info("subscribed to subject",
		"subject", subject,
		"consumer", consumerName,
	)

	return nil
}

// SubscribeWithStruct subscribes and unmarshals JSON into a struct.
func (s *Subscriber) SubscribeWithStruct(
	ctx context.Context,
	subject, consumerName string,
	eventType interface{},
	handler func(ctx context.Context, event interface{}) error,
) error {
	return s.Subscribe(ctx, subject, consumerName, func(ctx context.Context, data []byte) error {
		// Unmarshal JSON into event struct
		if err := json.Unmarshal(data, eventType); err != nil {
			return fmt.Errorf("failed to unmarshal event: %w", err)
		}
		return handler(ctx, eventType)
	})
}

// Close closes the NATS connection.
func (s *Subscriber) Close() {
	if s.nc != nil {
		s.nc.Close()
	}
}

// Health checks if the NATS connection is healthy.
func (s *Subscriber) Health() error {
	if s.nc == nil {
		return fmt.Errorf("NATS connection is nil")
	}
	if !s.nc.IsConnected() {
		return fmt.Errorf("NATS is not connected")
	}
	return nil
}
