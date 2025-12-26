package roomeventconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/mordred-r1/player-service/internal/models"
	"github.com/segmentio/kafka-go"
)

func (c *RoomEventConsumer) Consume(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           c.kafkaBroker,
		GroupID:           "RoomEventConsumer_group",
		Topic:             c.topicName,
		HeartbeatInterval: 3 * time.Second,
		SessionTimeout:    30 * time.Second,
	})
	defer r.Close()

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			slog.Error("RoomEventConsumer.consume error", "error", err.Error())
		} else {
			fmt.Println("Received message:", string(msg.Value))
		}
		var roomEvent *models.RoomEvent
		err = json.Unmarshal(msg.Value, &roomEvent)
		if err != nil {
			slog.Error("parce", "error", err)
			continue
		}
		err = c.roomEventProcessor.Handle(ctx, roomEvent)
		if err != nil {
			slog.Error("Handle", "error", err)
		}
	}

}
