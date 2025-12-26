package playerservice

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/mordred-r1/player-service/internal/models"
)

func (s *PlayerService) Update(ctx context.Context, player *models.PlayerState) error {
	if player == nil {
		return errors.New("player is required")
	}
	if player.ID == "" {
		return errors.New("player ID is required")
	}
	if player.State == "" {
		return errors.New("player state is required")
	}
	if err := s.playerStorage.Update(ctx, player); err != nil {
		return err
	}
	if err := s.playerEventProducer.Produce(ctx, &models.PlayerEvent{ID: player.ID, State: player.State}); err != nil {
		log.Printf("failed to produce player updated event for %s: %v", player.ID, err)
	}

	// best-effort: update cache
	if s.cache != nil {
		_ = s.cache.SetPlayer(context.Background(), player, 5*time.Minute)
	}
	return nil
}
