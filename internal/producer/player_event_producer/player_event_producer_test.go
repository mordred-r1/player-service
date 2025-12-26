package playereventproducer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mordred-r1/player-service/internal/models"
	"github.com/segmentio/kafka-go"
)

// mockWriter implements kafkaWriter for tests.
type mockWriter struct {
	writeErr error
	closed   bool
	msgs     []kafka.Message
}

func (m *mockWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	if m.writeErr != nil {
		return m.writeErr
	}
	m.msgs = append(m.msgs, msgs...)
	return nil
}

func (m *mockWriter) Close() error {
	m.closed = true
	return nil
}

func TestProduce_Success(t *testing.T) {
	m := &mockWriter{}
	p := &PlayerEventProducer{writer: m, topic: "players.events"}
	ctx := context.Background()
	e := &models.PlayerEvent{ID: "123", State: "playing"}
	if err := p.Produce(ctx, e); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(m.msgs) != 1 {
		t.Fatalf("expected 1 message written, got %d", len(m.msgs))
	}
	if string(m.msgs[0].Key) != e.ID {
		t.Fatalf("expected message key %s, got %s", e.ID, string(m.msgs[0].Key))
	}
}

func TestProduce_NilEvent(t *testing.T) {
	m := &mockWriter{}
	p := &PlayerEventProducer{writer: m}
	ctx := context.Background()
	if err := p.Produce(ctx, nil); err == nil {
		t.Fatalf("expected error for nil event")
	}
}

func TestProduce_WriteError(t *testing.T) {
	m := &mockWriter{writeErr: errors.New("write failed")}
	p := &PlayerEventProducer{writer: m}
	ctx := context.Background()
	e := &models.PlayerEvent{ID: "123", State: "stopped"}
	if err := p.Produce(ctx, e); err == nil {
		t.Fatalf("expected error from writer")
	}
}

func TestClose(t *testing.T) {
	m := &mockWriter{}
	p := &PlayerEventProducer{writer: m}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := p.Close(ctx); err != nil {
		t.Fatalf("expected no error on close, got %v", err)
	}
	if !m.closed {
		t.Fatalf("expected writer to be closed")
	}
}
