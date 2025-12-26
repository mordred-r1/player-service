package roomeventconsumer

import (
	"context"

	"github.com/mordred-r1/player-service/internal/models"
)

type roomEventProcessor interface {
	Handle(ctx context.Context, roomEvent *models.RoomEvent) error
}

type RoomEventConsumer struct {
	roomEventProcessor roomEventProcessor
	kafkaBroker        []string
	topicName          string
}

func NewRoomEventConsumer(roomEventProcessor roomEventProcessor, kafkaBroker []string, topicName string) *RoomEventConsumer {
	return &RoomEventConsumer{
		roomEventProcessor: roomEventProcessor,
		kafkaBroker:        kafkaBroker,
		topicName:          topicName,
	}
}
