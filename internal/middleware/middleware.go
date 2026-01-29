package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	apierrors "github.com/m-szczepanski/gocalc-api/internal/errors"
	"github.com/m-szczepanski/gocalc-api/internal/models"
)

type contextKey string

const (
	RequestIDKey contextKey = "request-id"
)

type ErrorHandler struct {
	handler http.Handler
}

func NewErrorHandler(handler http.Handler) *ErrorHandler {
	return &ErrorHandler{handler: handler}
}

func (eh *ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	defer func() {
		if err := recover(); err != nil {
			slog.Error("panic recovered",
				"error", err,
				"stack", string(debug.Stack()),
				"path", r.URL.Path,
				"method", r.Method,
			)
			rw.statusCode = http.StatusInternalServerError
			writeErrorResponse(rw, apierrors.InternalError("internal server error"), r)
		}
	}()

	eh.handler.ServeHTTP(rw, r)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	if !rw.written {
		rw.statusCode = statusCode
		rw.written = true
		rw.ResponseWriter.WriteHeader(statusCode)
	}
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	if !rw.written {
		rw.written = true
	}
	return rw.ResponseWriter.Write(data)
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := generateRequestID()

		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		r = r.WithContext(ctx)

		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := ExtractRequestID(r.Context())
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		slog.Info("request started",
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		slog.Info("request completed",
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.statusCode,
			"duration_ms", duration.Milliseconds(),
			"remote_addr", r.RemoteAddr,
		)
	})
}

func writeErrorResponse(w http.ResponseWriter, apiErr *apierrors.APIError, r *http.Request) {
	requestID := ExtractRequestID(r.Context())
	statusCode := apiErr.HTTPStatus()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := models.NewAPIErrorResponse(apiErr.Code, apiErr.Message, apiErr.Details, requestID)
	if err := encodeJSON(w, resp); err != nil {
		slog.Error("failed to encode error response",
			"error", err,
			"request_id", requestID,
			"status_code", statusCode,
		)
	}
}

func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r = r.WithContext(ctx)

			done := make(chan struct{})
			panicChan := make(chan interface{}, 1)

			go func() {
				defer func() {
					if p := recover(); p != nil {
						panicChan <- p
					}
				}()
				next.ServeHTTP(w, r)
				close(done)
			}()

			select {
			case p := <-panicChan:
				panic(p)
			case <-done:
				return
			case <-ctx.Done():
				requestID := ExtractRequestID(r.Context())
				slog.Warn("request timeout",
					"request_id", requestID,
					"method", r.Method,
					"path", r.URL.Path,
					"timeout", timeout.String(),
				)
				writeErrorResponse(w, apierrors.InternalError("request timeout"), r)
			}
		})
	}
}

func ExtractRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		return "unknown"
	}
	return requestID
}

func generateRequestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("err-%d", fastTimestamp())
	}
	return hex.EncodeToString(bytes)
}

func fastTimestamp() int64 {
	return 0 // Will be replaced by actual timestamp if needed
}
