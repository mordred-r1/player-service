package bootstrap

import (
	"github.com/mordred-r1/player-service/config"
	rediscache "github.com/mordred-r1/player-service/internal/cache/redis_cache"
)

// InitRedis initializes and returns a Redis cache based on config.
func InitRedis(cfg *config.Config) *rediscache.RedisCache {
	if cfg == nil {
		return nil
	}
	return rediscache.New(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
}
