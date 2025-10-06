package api

import (
	"time"

	"github.com/deal-machine/go-expert/rate-limiter/internal/protocols"
)

type RateLimiter struct {
	Strategy   protocols.IRateLimitStrategy
	Limit      int
	Expiration time.Duration
}

func NewRateLimiter(
	strategy protocols.IRateLimitStrategy,
	limit int,
	expiration time.Duration,
) *RateLimiter {
	return &RateLimiter{
		Strategy:   strategy,
		Limit:      limit,
		Expiration: expiration,
	}
}
