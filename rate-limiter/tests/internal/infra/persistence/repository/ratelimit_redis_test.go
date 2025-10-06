package repository

import (
	"context"
	"testing"
	"time"

	"github.com/deal-machine/go-expert/rate-limiter/internal/infra/persistence/repository"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitRedisRepository_IncrementAndGet(t *testing.T) {
	ctx := context.Background()
	t.Cleanup(func() {
		redisClient.FlushDB(ctx)
	})

	repo := repository.NewRedisRepository(redisClient)

	// loop increment key
	for i := 1; i <= 10; i++ {
		incr, err := repo.IncrementAndGet(ctx, "key_test")
		assert.Nil(t, err)
		assert.Equal(t, int64(i), incr)
	}

	// total increment value
	incr, err := repo.IncrementAndGet(ctx, "key_test")
	assert.Nil(t, err)
	assert.Equal(t, int64(11), incr)

	// another key
	incr, err = repo.IncrementAndGet(ctx, "another_key_test")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), incr)
}
func TestRateLimitRedisRepository_DefineExpiration(t *testing.T) {
	ctx := context.Background()
	t.Cleanup(func() {
		redisClient.FlushDB(ctx)
	})

	repo := repository.NewRedisRepository(redisClient)

	incr, err := repo.IncrementAndGet(ctx, "key")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), incr)
	incr, err = repo.IncrementAndGet(ctx, "key")
	assert.Nil(t, err)
	assert.Equal(t, int64(2), incr)

	err = repo.DefineExpiration(ctx, "key", time.Second)
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)

	incr, err = repo.IncrementAndGet(ctx, "key")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), incr)
}
