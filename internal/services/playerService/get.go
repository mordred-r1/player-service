package playerservice

import (
	"context"
	"errors"
	"time"

	"github.com/mordred-r1/player-service/internal/models"
)

func (s *PlayerService) Get(ctx context.Context, id string) (*models.PlayerState, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	// Try cache first
	if s.cache != nil {
		ps, err := s.cache.GetPlayer(ctx, id)
		if err == nil {
			return ps, nil
		}
	}

	// Fallback to storage
	ps, err := s.playerStorage.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Best-effort cache population
	if s.cache != nil {
		go func(p *models.PlayerState) {
			_ = s.cache.SetPlayer(context.Background(), p, 5*time.Minute)
		}(ps)
	}

	return ps, nil
}
