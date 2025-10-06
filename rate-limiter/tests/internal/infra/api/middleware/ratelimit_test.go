package middleware

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/deal-machine/go-expert/rate-limiter/configs"
	"github.com/deal-machine/go-expert/rate-limiter/internal/errors"
	"github.com/deal-machine/go-expert/rate-limiter/internal/infra/api/routers"
	"github.com/deal-machine/go-expert/rate-limiter/internal/infra/persistence/repository"
	"github.com/deal-machine/go-expert/rate-limiter/internal/strategy"
	"github.com/stretchr/testify/assert"
)

func TestApiEndpoint(t *testing.T) {
	ctx := context.Background()
	t.Cleanup(func() {
		redisClient.FlushDB(ctx)
	})

	repo := repository.NewRedisRepository(redisClient)
	strategy := strategy.NewRateLimitStrategy(repo)
	rateLimitVariables := configs.RateLimitVariables{
		ApiKeyLimit:      250,
		ApiKeyExpiration: time.Second,
		IpLimit:          1000,
		IpExpiration:     time.Second,
	}
	r := routers.Routers(strategy, rateLimitVariables)
	server := httptest.NewServer(r)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL+"/api", nil)
	assert.Nil(t, err)
	req.Header.Set("X-API-KEY", "api_key_test")

	nRequests := 500
	results := make(chan int, nRequests)
	for range nRequests {
		go func() {
			resp, err := http.DefaultClient.Do(req)
			assert.Nil(t, err)
			defer resp.Body.Close()
			results <- resp.StatusCode
		}()
	}

	var allowedCount, blockedCount int
	for range nRequests {
		if <-results == 200 {
			allowedCount++
		} else {
			blockedCount++
		}
	}
	assert.Equal(t, allowedCount, 249)
	assert.Equal(t, blockedCount, 251)
}

func TestApiKeyEndpoint(t *testing.T) {
	ctx := context.Background()
	t.Cleanup(func() {
		redisClient.FlushDB(ctx)
	})

	repo := repository.NewRedisRepository(redisClient)
	strategy := strategy.NewRateLimitStrategy(repo)
	rateLimitVariables := configs.RateLimitVariables{
		ApiKeyLimit:      100,
		ApiKeyExpiration: time.Second,
		IpLimit:          100,
		IpExpiration:     time.Second,
	}
	r := routers.Routers(strategy, rateLimitVariables)
	server := httptest.NewServer(r)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL+"/apiKey", nil)
	assert.Nil(t, err)
	req.Header.Set("X-API-KEY", "api_key_test")

	nRequests := 200
	results := make(chan bool, nRequests)
	for i := 0; i <= nRequests; i++ {
		go func() {
			resp, err := http.DefaultClient.Do(req)
			assert.Nil(t, err)
			defer resp.Body.Close()
			results <- http.StatusOK == resp.StatusCode
		}()
	}

	var allowedCount, blockedCount int
	for i := 0; i <= nRequests; i++ {
		if <-results {
			allowedCount++
		} else {
			blockedCount++
		}
	}
	assert.LessOrEqual(t, allowedCount, 100)
	assert.GreaterOrEqual(t, blockedCount, 100)

	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	assert.Equal(t, errors.ErrTooManyRequests.Error(), string(body))

	time.Sleep(2 * time.Second)

	resp, err = http.DefaultClient.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello World!", string(body))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestIpEndpoint(t *testing.T) {
	ctx := context.Background()
	t.Cleanup(func() {
		redisClient.FlushDB(ctx)
	})

	repo := repository.NewRedisRepository(redisClient)
	strategy := strategy.NewRateLimitStrategy(repo)
	rateLimitVariables := configs.RateLimitVariables{
		ApiKeyLimit:      100,
		ApiKeyExpiration: time.Second,
		IpLimit:          100,
		IpExpiration:     time.Second,
	}
	r := routers.Routers(strategy, rateLimitVariables)
	server := httptest.NewServer(r)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL+"/ip", nil)
	assert.Nil(t, err)
	for range 99 {
		resp, err := http.DefaultClient.Do(req)
		assert.Nil(t, err)
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)
		assert.Equal(t, "Hello World!", string(body))
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}

	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	assert.Equal(t, errors.ErrTooManyRequests.Error(), string(body))
}
