package playerservice

import (
	"context"
	"errors"
	"log"

	"github.com/mordred-r1/player-service/internal/models"
)

func (s *PlayerService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	if err := s.playerStorage.Delete(ctx, id); err != nil {
		return err
	}

	if err := s.playerEventProducer.Produce(ctx, &models.PlayerEvent{ID: id, State: "deleted"}); err != nil {
		log.Printf("failed to produce player deleted event for %s: %v", id, err)
	}
	return nil
}
