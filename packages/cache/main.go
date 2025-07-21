package cache

import (
	"context"
	"time"

	"github.com/ESSantana/boilerplate-backend/internal/config"
	"github.com/ESSantana/boilerplate-backend/packages/cache/interfaces"
	"github.com/redis/go-redis/v9"
)

type cacheManager struct {
	client *redis.Client
}

func NewCacheManager(cfg *config.Config) interfaces.CacheManager {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Username: cfg.Redis.User,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	return &cacheManager{
		client: client,
	}
}

func (cache *cacheManager) CacheHealthCheck() error {
	return cache.client.Ping(context.Background()).Err()
}

func (cache *cacheManager) SetFlagWithExpiration(ctx context.Context, key string, value bool, expiration time.Duration) error {
	err := cache.client.Set(ctx, key, value, expiration).Err()
	return err
}

func (cache *cacheManager) GetFlag(ctx context.Context, key string) (bool, error) {
	val, err := cache.client.Get(ctx, key).Bool()
	return val, err
}

func (cache *cacheManager) SetStringWithExpiration(ctx context.Context, key, value string, expiration time.Duration) error {
	err := cache.client.Set(ctx, key, value, expiration).Err()
	return err
}

func (cache *cacheManager) GetString(ctx context.Context, key string) (string, error) {
	val, err := cache.client.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return val, err
}
