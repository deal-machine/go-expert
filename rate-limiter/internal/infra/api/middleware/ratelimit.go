package middleware

import (
	"log"
	"net"
	"net/http"
	"strings"

	internalErrors "github.com/deal-machine/go-expert/rate-limiter/internal/errors"
	"github.com/deal-machine/go-expert/rate-limiter/internal/infra/api"
)

type rateLimitMiddleware struct {
	LimiterApiKey *api.RateLimiter
	LimiterIP     *api.RateLimiter
}

func RateLimitMiddleware(limiterApiKey *api.RateLimiter, limiterIP *api.RateLimiter) *rateLimitMiddleware {
	return &rateLimitMiddleware{
		LimiterApiKey: limiterApiKey,
		LimiterIP:     limiterIP,
	}
}

func (rl *rateLimitMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var statusCode int = http.StatusBadRequest

		apiKey := strings.TrimSpace(r.Header.Get("X-API-Key"))
		if apiKey != "" && rl.LimiterApiKey != nil {
			log.Println("apiKey", apiKey)
			err := rl.LimiterApiKey.Strategy.CheckLimit(
				apiKey,
				rl.LimiterApiKey.Limit,
				rl.LimiterApiKey.Expiration,
			)
			if err != nil {
				if err == internalErrors.ErrTooManyRequests {
					statusCode = http.StatusTooManyRequests
				}
				w.WriteHeader(statusCode)
				w.Write([]byte(err.Error()))
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(statusCode)
			w.Write([]byte(err.Error()))
			return
		}
		if rl.LimiterIP != nil {
			err = rl.LimiterIP.Strategy.CheckLimit(host, rl.LimiterIP.Limit, rl.LimiterIP.Expiration)
			if err != nil {
				if err == internalErrors.ErrTooManyRequests {
					statusCode = http.StatusTooManyRequests
				}
				w.WriteHeader(statusCode)
				w.Write([]byte(err.Error()))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
