package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger := zerolog.Ctx(r.Context())
		logger.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Dur("duration", time.Since(start)).
			Msg("HTTP request")
	})
}