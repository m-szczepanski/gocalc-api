package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	apierrors "github.com/m-szczepanski/gocalc-api/internal/errors"
)

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
	cleanup  time.Duration
}

func NewRateLimiter(rps float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(rps),
		burst:    burst,
		cleanup:  5 * time.Minute,
	}

	go rl.cleanupRoutine()

	return rl
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = limiter
	}

	return limiter
}

func (rl *RateLimiter) cleanupRoutine() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		_ = len(rl.limiters)
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r)
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			requestID := ExtractRequestID(r.Context())
			slog.Warn("rate limit exceeded",
				"request_id", requestID,
				"ip", ip,
				"method", r.Method,
				"path", r.URL.Path,
			)

			w.Header().Set("X-RateLimit-Limit", "100")
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("Retry-After", "60")

			apiErr := apierrors.RateLimitExceeded("rate limit exceeded, please try again later")
			writeErrorResponse(w, apiErr, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func extractIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if ip, _, err := net.SplitHostPort(xff); err == nil {
			return ip
		}
		return xff
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return ip
	}

	return r.RemoteAddr
}
