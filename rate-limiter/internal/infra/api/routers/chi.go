package routers

import (
	"github.com/deal-machine/go-expert/rate-limiter/configs"
	"github.com/deal-machine/go-expert/rate-limiter/internal/infra/api"
	"github.com/deal-machine/go-expert/rate-limiter/internal/infra/api/handlers"
	"github.com/deal-machine/go-expert/rate-limiter/internal/infra/api/middleware"
	"github.com/deal-machine/go-expert/rate-limiter/internal/protocols"
	"github.com/go-chi/chi/v5"
	go_chi_middleware "github.com/go-chi/chi/v5/middleware"
)

func Routers(s protocols.IRateLimitStrategy, rateLimitVariables configs.RateLimitVariables) *chi.Mux {

	limiterApiKey := api.NewRateLimiter(s, rateLimitVariables.ApiKeyLimit, rateLimitVariables.ApiKeyExpiration)
	limiterIp := api.NewRateLimiter(s, rateLimitVariables.IpLimit, rateLimitVariables.IpExpiration)

	r := chi.NewRouter()

	r.Use(go_chi_middleware.Recoverer)
	r.Use(go_chi_middleware.Logger)

	r.Route("/apiKey", func(r chi.Router) {
		r.Use(middleware.RateLimitMiddleware(limiterApiKey, nil).Handle)
		r.Get("/", handlers.HelloWorldHandler)
	})
	r.Route("/ip", func(r chi.Router) {
		r.Use(middleware.RateLimitMiddleware(nil, limiterIp).Handle)
		r.Get("/", handlers.HelloWorldHandler)
	})
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.RateLimitMiddleware(limiterApiKey, limiterIp).Handle)
		r.Get("/", handlers.HelloWorldHandler)
	})

	return r
}
