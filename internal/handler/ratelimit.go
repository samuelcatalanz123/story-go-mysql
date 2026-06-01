package handler

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"story-go-mysql/internal/cache"
)

// RateLimit returns a middleware that allows at most `limit` requests per
// client IP within `window`. It counts hits in Redis with INCR (+EXPIRE on the
// first hit). Over the limit, it replies 429 with a Retry-After header.
//
// With cache.Noop (no Redis) Incr returns 0, so limiting is effectively off —
// the app keeps working.
func RateLimit(c cache.Cache, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := fmt.Sprintf("rl:%s:%s", r.URL.Path, clientIP(r))
			count, err := c.Incr(r.Context(), key, window)
			if err == nil && count > int64(limit) {
				w.Header().Set("Retry-After", strconv.Itoa(int(window.Seconds())))
				http.Error(w, `{"error":"too many requests, slow down"}`, http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// clientIP extracts the client's IP from RemoteAddr (host:port).
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
