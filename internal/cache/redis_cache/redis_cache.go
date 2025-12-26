package rediscache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/mordred-r1/player-service/internal/models"
	"github.com/redis/go-redis/v9"
)

// RedisCache provides simple get/set helpers for PlayerState objects.
type RedisCache struct {
	rdb *redis.Client
}

// New creates a new RedisCache from explicit parameters. It does not perform network checks.
func New(addr, password string, db int) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisCache{rdb: rdb}
}

// Close closes the underlying redis client.
func (c *RedisCache) Close(ctx context.Context) error {
	if c == nil || c.rdb == nil {
		return nil
	}
	return c.rdb.Close()
}

// GetPlayer returns PlayerState from cache or redis.Nil error if missing.
func (c *RedisCache) GetPlayer(ctx context.Context, id string) (*models.PlayerState, error) {
	if c == nil || c.rdb == nil {
		return nil, redis.Nil
	}
	key := "player:" + id
	val, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var ps models.PlayerState
	if err := json.Unmarshal([]byte(val), &ps); err != nil {
		return nil, err
	}
	return &ps, nil
}

// SetPlayer stores PlayerState in cache with ttl.
func (c *RedisCache) SetPlayer(ctx context.Context, p *models.PlayerState, ttl time.Duration) error {
	if c == nil || c.rdb == nil {
		return nil
	}
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	key := "player:" + p.ID
	return c.rdb.Set(ctx, key, string(b), ttl).Err()
}

// DeletePlayer removes player from cache.
func (c *RedisCache) DeletePlayer(ctx context.Context, id string) error {
	if c == nil || c.rdb == nil {
		return nil
	}
	key := "player:" + id
	return c.rdb.Del(ctx, key).Err()
}
