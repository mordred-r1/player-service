package playercommandprocessor

import (
	"context"

	"github.com/mordred-r1/player-service/internal/models"
)

type playerService interface {
	// Get(ctx context.Context, id string) (*models.PlayerState, error)
	Update(ctx context.Context, player *models.PlayerState) error
}

type PlayerCommandProcessor struct {
	playerService playerService
}

func NewPlayerCommandProcessor(playerService playerService) *PlayerCommandProcessor {
	return &PlayerCommandProcessor{
		playerService: playerService,
	}
}
