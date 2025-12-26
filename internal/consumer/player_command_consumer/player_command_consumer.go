package playercommandconsumer

import (
	"context"

	"github.com/mordred-r1/player-service/internal/models"
)

type playerCommandProcessor interface {
	Handle(ctx context.Context, playerCommand *models.PlayerCommand) error
}

type PlayerCommandConsumer struct {
	playerCommandProcessor playerCommandProcessor
	kafkaBroker            []string
	topicName              string
}

func NewPlayerCommandConsumer(playerCommandProcessor playerCommandProcessor, kafkaBroker []string, topicName string) *PlayerCommandConsumer {
	return &PlayerCommandConsumer{
		playerCommandProcessor: playerCommandProcessor,
		kafkaBroker:            kafkaBroker,
		topicName:              topicName,
	}
}
