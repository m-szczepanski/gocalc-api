package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/m-szczepanski/gocalc-api/internal/models"
)

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(10, 5)

	if rl == nil {
		t.Fatal("expected non-nil rate limiter")
	}

	if rl.rate != 10 {
		t.Errorf("expected rate 10, got %f", rl.rate)
	}

	if rl.burst != 5 {
		t.Errorf("expected burst 5, got %d", rl.burst)
	}

	if rl.limiters == nil {
		t.Error("expected non-nil limiters map")
	}
}

func TestRateLimiter_GetLimiter(t *testing.T) {
	rl := NewRateLimiter(10, 5)

	limiter1 := rl.getLimiter("192.168.1.1")
	if limiter1 == nil {
		t.Fatal("expected non-nil limiter")
	}

	limiter2 := rl.getLimiter("192.168.1.1")
	if limiter1 != limiter2 {
		t.Error("expected same limiter for same IP")
	}

	limiter3 := rl.getLimiter("192.168.1.2")
	if limiter1 == limiter3 {
		t.Error("expected different limiters for different IPs")
	}
}

func TestRateLimiter_Middleware_AllowsRequests(t *testing.T) {
	rl := NewRateLimiter(10, 5)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	middleware := rl.Middleware(handler)

	for i := range 5 {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		ctx := context.WithValue(req.Context(), RequestIDKey, "test-123")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d: expected status 200, got %d", i+1, w.Code)
		}
	}
}

func TestRateLimiter_Middleware_BlocksExcessRequests(t *testing.T) {
	rl := NewRateLimiter(1, 2)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	middleware := rl.Middleware(handler)

	for i := range 2 {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		ctx := context.WithValue(req.Context(), RequestIDKey, "test-123")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d: expected status 200, got %d", i+1, w.Code)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	ctx := context.WithValue(req.Context(), RequestIDKey, "test-123")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected status 429, got %d", w.Code)
	}

	// Verify headers (rate is 1 req/sec = 60 req/min)
	if w.Header().Get("X-RateLimit-Limit") != "60" {
		t.Errorf("expected X-RateLimit-Limit header: 60, got %s", w.Header().Get("X-RateLimit-Limit"))
	}

	if w.Header().Get("Retry-After") != "60" {
		t.Errorf("expected Retry-After header: 60, got %s", w.Header().Get("Retry-After"))
	}

	// Verify error response
	var errResp models.APIErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}

	if errResp.Code != "RATE_LIMIT_EXCEEDED" {
		t.Errorf("expected error code RATE_LIMIT_EXCEEDED, got %s", errResp.Code)
	}
}

func TestRateLimiter_Middleware_DifferentIPs(t *testing.T) {
	rl := NewRateLimiter(1, 1)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	middleware := rl.Middleware(handler)

	// Request from IP1
	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	ctx1 := context.WithValue(req1.Context(), RequestIDKey, "test-123")
	req1 = req1.WithContext(ctx1)
	w1 := httptest.NewRecorder()
	middleware.ServeHTTP(w1, req1)

	// Request from IP2
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.RemoteAddr = "192.168.1.2:12345"
	ctx2 := context.WithValue(req2.Context(), RequestIDKey, "test-456")
	req2 = req2.WithContext(ctx2)
	w2 := httptest.NewRecorder()
	middleware.ServeHTTP(w2, req2)

	if w1.Code != http.StatusOK {
		t.Errorf("IP1: expected status 200, got %d", w1.Code)
	}

	if w2.Code != http.StatusOK {
		t.Errorf("IP2: expected status 200, got %d", w2.Code)
	}
}

func TestRateLimiter_Middleware_RecoveryAfterTime(t *testing.T) {
	rl := NewRateLimiter(10, 1)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	middleware := rl.Middleware(handler)

	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	ctx1 := context.WithValue(req1.Context(), RequestIDKey, "test-123")
	req1 = req1.WithContext(ctx1)
	w1 := httptest.NewRecorder()
	middleware.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("first request: expected status 200, got %d", w1.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.RemoteAddr = "192.168.1.1:12345"
	ctx2 := context.WithValue(req2.Context(), RequestIDKey, "test-456")
	req2 = req2.WithContext(ctx2)
	w2 := httptest.NewRecorder()
	middleware.ServeHTTP(w2, req2)

	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("second request: expected status 429, got %d", w2.Code)
	}

	time.Sleep(150 * time.Millisecond)

	req3 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req3.RemoteAddr = "192.168.1.1:12345"
	ctx3 := context.WithValue(req3.Context(), RequestIDKey, "test-789")
	req3 = req3.WithContext(ctx3)
	w3 := httptest.NewRecorder()
	middleware.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Errorf("third request: expected status 200, got %d", w3.Code)
	}
}

func TestExtractIP(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		expectedIP string
	}{
		{
			name:       "from RemoteAddr with port",
			remoteAddr: "192.168.1.1:12345",
			expectedIP: "192.168.1.1",
		},
		{
			name:       "from RemoteAddr without port",
			remoteAddr: "192.168.1.1",
			expectedIP: "192.168.1.1",
		},
		{
			name:       "IPv6 with port",
			remoteAddr: "[2001:db8::1]:8080",
			expectedIP: "2001:db8::1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.RemoteAddr = tt.remoteAddr

			ip := extractIP(req)
			if ip != tt.expectedIP {
				t.Errorf("expected IP %s, got %s", tt.expectedIP, ip)
			}
		})
	}
}
