package middleware

import (
	"github.com/chandra-shekhar/internal-transfers/internal/server"
)

type Middlewares struct {
	Global          *GlobalMiddlewares
	ContextEnhancer *ContextEnhancer
	RateLimit       *RateLimitMiddleware
}

func NewMiddlewares(s *server.Server) *Middlewares {
	return &Middlewares{
		Global:          NewGlobalMiddlewares(s),
		ContextEnhancer: NewContextEnhancer(s),
		RateLimit:       NewRateLimitMiddleware(s),
	}
}
