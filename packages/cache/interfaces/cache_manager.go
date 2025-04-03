package interfaces

import (
	"context"
	"time"
)

type CacheManager interface {
	CacheHealthCheck() error
	SetFlagWithExpiration(ctx context.Context, key string, value bool, expiration time.Duration) error
	GetFlag(ctx context.Context, key string) (bool, error)
	SetStringWithExpiration(ctx context.Context, key, value string, expiration time.Duration) error
	GetString(ctx context.Context, key string) (string, error)
}
