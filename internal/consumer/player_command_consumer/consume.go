package playercommandconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/mordred-r1/player-service/internal/models"
	"github.com/segmentio/kafka-go"
)

func (c *PlayerCommandConsumer) Consume(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           c.kafkaBroker,
		GroupID:           "PlayerCommandConsumer_group",
		Topic:             c.topicName,
		HeartbeatInterval: 3 * time.Second,
		SessionTimeout:    30 * time.Second,
	})
	defer r.Close()

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			slog.Error("PlayerCommandConsumer.consume error", "error", err.Error())
		} else {
			fmt.Println("Received message:", string(msg.Value))
		}
		var playerCommand *models.PlayerCommand
		err = json.Unmarshal(msg.Value, &playerCommand)
		if err != nil {
			slog.Error("parce", "error", err)
			continue
		}
		err = c.playerCommandProcessor.Handle(ctx, playerCommand)
		if err != nil {
			slog.Error("Handle", "error", err)
		}
	}

}
