package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, addr string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
