package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/m-szczepanski/gocalc-api/internal/models"
)

func TestRequestIDMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := ExtractRequestID(r.Context())
		if requestID == "" || requestID == "unknown" {
			t.Error("expected request ID in context, got empty or unknown")
		}
		w.WriteHeader(http.StatusOK)
	})

	middleware := RequestIDMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	requestID := w.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Error("expected X-Request-ID header, got empty")
	}

	if len(requestID) != 32 {
		t.Errorf("expected request ID length 32, got %d", len(requestID))
	}
}

func TestRequestIDMiddleware_UniqueIDs(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := RequestIDMiddleware(handler)

	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		middleware.ServeHTTP(w, req)

		requestID := w.Header().Get("X-Request-ID")
		if ids[requestID] {
			t.Errorf("duplicate request ID generated: %s", requestID)
		}
		ids[requestID] = true
	}
}

func TestLoggingMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := LoggingMiddleware(handler)

	req := httptest.NewRequest(http.MethodPost, "/api/math/add", nil)
	ctx := context.WithValue(req.Context(), RequestIDKey, "test-id-123")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestErrorHandler_NoPanic(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	errorHandler := NewErrorHandler(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx := context.WithValue(req.Context(), RequestIDKey, "test-123")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	errorHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestErrorHandler_WithPanic(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	errorHandler := NewErrorHandler(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx := context.WithValue(req.Context(), RequestIDKey, "test-123")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	errorHandler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}

	var resp models.APIErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if resp.Code != "INTERNAL_ERROR" {
		t.Errorf("expected error code INTERNAL_ERROR, got %s", resp.Code)
	}

	if resp.RequestID != "test-123" {
		t.Errorf("expected request ID test-123, got %s", resp.RequestID)
	}
}

func TestErrorHandler_WithPanicString(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("string panic message")
	})

	errorHandler := NewErrorHandler(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx := context.WithValue(req.Context(), RequestIDKey, "panic-test")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	errorHandler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	w := httptest.NewRecorder()
	rw := &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	rw.WriteHeader(http.StatusCreated)

	if rw.statusCode != http.StatusCreated {
		t.Errorf("expected status code 201, got %d", rw.statusCode)
	}

	if !rw.written {
		t.Error("expected written flag to be true")
	}

	rw.WriteHeader(http.StatusBadRequest)
	if rw.statusCode != http.StatusCreated {
		t.Error("status code should not change after first write")
	}
}

func TestResponseWriter_Write(t *testing.T) {
	w := httptest.NewRecorder()
	rw := &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	data := []byte("test data")
	n, err := rw.Write(data)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if n != len(data) {
		t.Errorf("expected %d bytes written, got %d", len(data), n)
	}

	if !rw.written {
		t.Error("expected written flag to be true")
	}
}

func TestExtractRequestID(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "with request ID",
			ctx:      context.WithValue(context.Background(), RequestIDKey, "test-id-123"),
			expected: "test-id-123",
		},
		{
			name:     "without request ID",
			ctx:      context.Background(),
			expected: "unknown",
		},
		{
			name:     "with wrong type",
			ctx:      context.WithValue(context.Background(), RequestIDKey, 12345),
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractRequestID(tt.ctx)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGenerateRequestID(t *testing.T) {
	id := generateRequestID()

	if len(id) != 32 {
		t.Errorf("expected request ID length 32, got %d", len(id))
	}

	for _, c := range id {
		if !strings.ContainsRune("0123456789abcdef", c) {
			t.Errorf("invalid hex character in request ID: %c", c)
		}
	}
}

func TestMiddlewareChain(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := ExtractRequestID(r.Context())
		if requestID == "" || requestID == "unknown" {
			t.Error("request ID should be set by middleware")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	var h http.Handler = handler
	h = RequestIDMiddleware(h)
	h = LoggingMiddleware(h)
	h = NewErrorHandler(h)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Header().Get("X-Request-ID") == "" {
		t.Error("expected X-Request-ID header")
	}
}

func TestLoggingMiddleware_ResponseTracking(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		expectedStatus int
	}{
		{
			name:           "success response",
			statusCode:     http.StatusOK,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "error response",
			statusCode:     http.StatusBadRequest,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "server error",
			statusCode:     http.StatusInternalServerError,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			})

			middleware := LoggingMiddleware(handler)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			ctx := context.WithValue(req.Context(), RequestIDKey, "test-123")
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			middleware.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestTimeoutMiddleware(t *testing.T) {
	t.Run("request completes before timeout", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		})

		middleware := TimeoutMiddleware(100 * time.Millisecond)(handler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), RequestIDKey, "test-123")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("request times out", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case <-time.After(200 * time.Millisecond):
				w.WriteHeader(http.StatusOK)
			case <-r.Context().Done():
				return
			}
		})

		middleware := TimeoutMiddleware(50 * time.Millisecond)(handler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), RequestIDKey, "test-123")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}

		var errResp models.APIErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
			t.Fatalf("failed to unmarshal error response: %v", err)
		}

		if errResp.Code != "INTERNAL_ERROR" {
			t.Errorf("expected error code INTERNAL_ERROR, got %s", errResp.Code)
		}
	})

	t.Run("context cancellation is propagated", func(t *testing.T) {
		contextCancelled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			<-r.Context().Done()
			contextCancelled = true
		})

		middleware := TimeoutMiddleware(50 * time.Millisecond)(handler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), RequestIDKey, "test-123")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		time.Sleep(100 * time.Millisecond)

		if !contextCancelled {
			t.Error("expected context to be cancelled in handler")
		}
	})
}
