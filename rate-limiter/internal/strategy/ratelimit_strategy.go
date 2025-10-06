package strategy

import (
	"context"
	"fmt"
	"log"
	"time"

	internalErrors "github.com/deal-machine/go-expert/rate-limiter/internal/errors"
	"github.com/deal-machine/go-expert/rate-limiter/internal/protocols"
)

type RateLimitStrategy struct {
	repository protocols.IRateLimitRepository
	ctx        context.Context
}

func NewRateLimitStrategy(r protocols.IRateLimitRepository) *RateLimitStrategy {
	return &RateLimitStrategy{
		repository: r,
		ctx:        context.Background(),
	}
}

func (rls *RateLimitStrategy) CheckLimit(key string, limit int, ex time.Duration) error {
	key = fmt.Sprintf("rate:%s", key)
	log.Println("ratelimiting key: ", key)
	count, err := rls.repository.IncrementAndGet(rls.ctx, key)
	if err != nil {
		log.Println("error on increment and get", err)
		return err
	}
	isFirstRequestForKey := count == 1
	if isFirstRequestForKey {
		err = rls.repository.DefineExpiration(rls.ctx, key, ex)
		if err != nil {
			log.Println("error on define expiration", err)
			return err
		}
	}
	isLimitReached := count >= int64(limit)
	if isLimitReached {
		log.Println("limit reached count", count)
		return internalErrors.ErrTooManyRequests
	}
	return nil
}
