package playerservice

import (
	"context"
	"time"

	"github.com/mordred-r1/player-service/internal/models"
)

type playerStorage interface {
	Create(ctx context.Context, player *models.PlayerState) error
	Get(ctx context.Context, id string) (*models.PlayerState, error)
	Update(ctx context.Context, player *models.PlayerState) error
	Delete(ctx context.Context, id string) error
}

type playerEventProducer interface {
	Produce(ctx context.Context, event *models.PlayerEvent) error
}

type PlayerService struct {
	playerStorage       playerStorage
	playerEventProducer playerEventProducer
	cache               playerCache
}

func NewPlayerService(ctx context.Context, playerStorage playerStorage, playerEventProducer playerEventProducer, cache playerCache) *PlayerService {
	return &PlayerService{
		playerStorage:       playerStorage,
		playerEventProducer: playerEventProducer,
		cache:               cache,
	}
}

// playerCache is used by the service for caching player states.
type playerCache interface {
	GetPlayer(ctx context.Context, id string) (*models.PlayerState, error)
	SetPlayer(ctx context.Context, p *models.PlayerState, ttl time.Duration) error
	DeletePlayer(ctx context.Context, id string) error
}
