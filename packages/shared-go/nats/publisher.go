package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Publisher handles publishing events to NATS JetStream.
type Publisher struct {
	nc *nats.Conn
	js jetstream.JetStream
}

// NewPublisher creates a new NATS publisher.
func NewPublisher(natsURL string) (*Publisher, error) {
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

	// Automatically bootstrap "cloudcommerce" JetStream stream
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     "cloudcommerce",
		Subjects: []string{"order.*", "payment.*", "inventory.*", "notification.*", "tenant.*"},
	})

	return &Publisher{
		nc: nc,
		js: js,
	}, nil
}

// Publish publishes an event to a subject.
func (p *Publisher) Publish(ctx context.Context, subject string, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	_, err = p.js.Publish(ctx, subject, data)
	if err != nil {
		return fmt.Errorf("failed to publish event to %s: %w", subject, err)
	}

	return nil
}

// PublishWithMetadata publishes an event with additional NATS message options.
func (p *Publisher) PublishWithMetadata(ctx context.Context, subject string, event interface{}, msgID string) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	_, err = p.js.Publish(ctx, subject, data, jetstream.WithMsgID(msgID))
	if err != nil {
		return fmt.Errorf("failed to publish event to %s: %w", subject, err)
	}

	return nil
}

// Close closes the NATS connection.
func (p *Publisher) Close() {
	if p.nc != nil {
		p.nc.Close()
	}
}

// Health checks if the NATS connection is healthy.
func (p *Publisher) Health() error {
	if p.nc == nil {
		return fmt.Errorf("NATS connection is nil")
	}
	if !p.nc.IsConnected() {
		return fmt.Errorf("NATS is not connected")
	}
	return nil
}
