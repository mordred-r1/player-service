package roomeventprocessor

import (
	"context"

	"github.com/mordred-r1/player-service/internal/models"
)

type playerService interface {
	// Get(ctx context.Context, id string) (*models.PlayerState, error)
	Create(ctx context.Context, player *models.PlayerState) error
	Delete(ctx context.Context, id string) error
}

type RoomEventProcessor struct {
	playerService playerService
}

func NewRoomEventProcessor(playerService playerService) *RoomEventProcessor {
	return &RoomEventProcessor{
		playerService: playerService,
	}
}
