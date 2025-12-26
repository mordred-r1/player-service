package bootstrap

import (
	"context"

	"github.com/mordred-r1/player-service/config"
	rediscache "github.com/mordred-r1/player-service/internal/cache/redis_cache"
	playereventproducer "github.com/mordred-r1/player-service/internal/producer/player_event_producer"
	playerservice "github.com/mordred-r1/player-service/internal/services/playerService"
	"github.com/mordred-r1/player-service/internal/storage/pgstorage"
)

func InitPlayerService(storage *pgstorage.PGstorage, cfg *config.Config, playerEventProducer *playereventproducer.PlayerEventProducer, redisCache *rediscache.RedisCache) *playerservice.PlayerService {
	return playerservice.NewPlayerService(context.Background(), storage, playerEventProducer, redisCache)
}
