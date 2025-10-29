package middleware

import (
	"github.com/chandra-shekhar/internal-transfers/internal/server"
)

type RateLimitMiddleware struct {
	server *server.Server
}

func NewRateLimitMiddleware(s *server.Server) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		server: s,
	}
}

func (r *RateLimitMiddleware) RecordRateLimitHit(endpoint string) {
	// Log rate limit hit
	r.server.Logger.Warn().
		Str("endpoint", endpoint).
		Msg("rate limit hit recorded")
}
