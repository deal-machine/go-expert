package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimitRedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RateLimitRedisRepository {
	return &RateLimitRedisRepository{
		client: client,
	}
}

func (r *RateLimitRedisRepository) IncrementAndGet(ctx context.Context, key string) (int64, error) {
	// Incrementa contador e pega valor atual
	val, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (r *RateLimitRedisRepository) DefineExpiration(ctx context.Context, key string, expiration time.Duration) error {
	// Primeira requisição, define expiração da janela
	return r.client.Expire(ctx, key, expiration).Err()
}
