package middleware

import (
	"context"
	"net/http"

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
		requestID := extractRequestID(r.Context())
		println("[" + requestID + "] " + r.Method + " " + r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func writeErrorResponse(w http.ResponseWriter, apiErr *apierrors.APIError, r *http.Request) {
	requestID := extractRequestID(r.Context())
	statusCode := apiErr.HTTPStatus()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := models.NewAPIErrorResponse(apiErr.Code, apiErr.Message, apiErr.Details, requestID)
	encodeJSON(w, resp)
}

func extractRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		return "unknown"
	}
	return requestID
}

func generateRequestID() string {
	return randomHexString(16)
}

func randomHexString(length int) string {
	const hex = "0123456789abcdef"
	result := make([]byte, length)
	for i := range result {
		result[i] = hex[fastRand()%len(hex)]
	}
	return string(result)
}

var fastRandSeed uint32 = 1

func fastRand() int {
	fastRandSeed = fastRandSeed*1103515245 + 12345
	return int((fastRandSeed / 65536) % 32768)
}
