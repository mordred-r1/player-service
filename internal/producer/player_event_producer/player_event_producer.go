package playereventproducer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"

	"github.com/mordred-r1/player-service/internal/models"
)

type PlayerEventProducer struct {
	writer kafkaWriter
	topic  string
}

// kafkaWriter is a minimal interface subset of *kafka.Writer used by the producer.
// Having this interface makes the producer testable with a mock writer.
type kafkaWriter interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

// NewPlayerEventProducer creates a producer that writes to the provided brokers and topic.
func NewPlayerEventProducer(brokers []string, topic string) *PlayerEventProducer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.Hash{},
		// reasonable defaults
		Async: false,
	}
	return &PlayerEventProducer{writer: w, topic: topic}
}

// Produce sends the PlayerEvent to Kafka as a JSON-encoded message. The event Name is used as the message key.
func (p *PlayerEventProducer) Produce(ctx context.Context, event *models.PlayerEvent) error {
	if event == nil {
		return errors.New("event is nil")
	}
	b, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "marshal event")
	}

	msg := kafka.Message{
		Key:   []byte(event.ID),
		Value: b,
		Time:  time.Now(),
	}

	// WriteMessage is blocking and returns when message is accepted by broker
	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return errors.Wrap(err, "write message to kafka")
	}
	return nil
}

// Close closes the underlying writer.
func (p *PlayerEventProducer) Close(ctx context.Context) error {
	if p.writer == nil {
		return nil
	}
	return p.writer.Close()
}
