package playerservice

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/mordred-r1/player-service/internal/models"
)

func (s *PlayerService) Create(ctx context.Context, player *models.PlayerState) error {
	if player == nil {
		return errors.New("player is required")
	}
	if player.ID == "" {
		return errors.New("player ID is required")
	}
	if player.State == "" {
		return errors.New("player state is required")
	}

	if err := s.playerStorage.Create(ctx, player); err != nil {
		return err
	}

	// publish event (log error only)
	if err := s.playerEventProducer.Produce(ctx, &models.PlayerEvent{ID: player.ID, State: player.State}); err != nil {
		log.Printf("failed to produce player created event for %s: %v", player.ID, err)
	}

	// best-effort: populate cache
	if s.cache != nil {
		_ = s.cache.SetPlayer(context.Background(), player, 5*time.Minute)
	}
	return nil
}
