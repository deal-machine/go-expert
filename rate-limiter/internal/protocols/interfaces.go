package protocols

import (
	"context"
	"time"
)

type IRateLimitRepository interface {
	IncrementAndGet(ctx context.Context, key string) (int64, error)
	DefineExpiration(ctx context.Context, key string, expiration time.Duration) error
}

type IRateLimitStrategy interface {
	CheckLimit(key string, limit int, ex time.Duration) error
}
