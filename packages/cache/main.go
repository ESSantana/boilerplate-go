package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/application-ellas/ella-backend/packages/cache/interfaces"
	"github.com/redis/go-redis/v9"
)

type cacheManager struct {
	client *redis.Client
}

func NewCacheManager() interfaces.CacheManager {
	fmt.Println("Connecting to Redis database...")
	fmt.Println("REDIS_HOST:", os.Getenv("REDIS_HOST"))
	fmt.Println("REDIS_USER:", os.Getenv("REDIS_USER"))
	fmt.Println("REDIS_PASSWORD:", os.Getenv("REDIS_PASSWORD"))
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PASSWORD"),
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
