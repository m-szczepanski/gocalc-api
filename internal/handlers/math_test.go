package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/m-szczepanski/gocalc-api/internal/middleware"
	"github.com/m-szczepanski/gocalc-api/internal/models"
)

func TestAddHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           *models.MathRequest
		expectedStatus int
		expectedResult float64
		expectError    bool
	}{
		{
			name:           "valid addition",
			method:         http.MethodPost,
			body:           &models.MathRequest{A: 5, B: 3},
			expectedStatus: http.StatusOK,
			expectedResult: 8,
			expectError:    false,
		},
		{
			name:           "negative numbers",
			method:         http.MethodPost,
			body:           &models.MathRequest{A: -5, B: 3},
			expectedStatus: http.StatusOK,
			expectedResult: -2,
			expectError:    false,
		},
		{
			name:           "decimals",
			method:         http.MethodPost,
			body:           &models.MathRequest{A: 1.5, B: 2.5},
			expectedStatus: http.StatusOK,
			expectedResult: 4,
			expectError:    false,
		},
		{
			name:           "method not allowed (GET)",
			method:         http.MethodGet,
			body:           &models.MathRequest{A: 5, B: 3},
			expectedStatus: http.StatusMethodNotAllowed,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != nil {
				body, _ := json.Marshal(tt.body)
				req = httptest.NewRequest(tt.method, "/api/math/add", bytes.NewReader(body))
			} else {
				req = httptest.NewRequest(tt.method, "/api/math/add", nil)
			}
			// Add request ID to context to simulate middleware
			ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-123")
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			AddHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.expectedStatus)
			}

			if tt.expectError {
				var resp models.APIErrorResponse
				json.NewDecoder(w.Body).Decode(&resp)
				if resp.Code == "" {
					t.Errorf("expected error code, got empty")
				}
			} else {
				var resp models.SuccessResponse
				json.NewDecoder(w.Body).Decode(&resp)
				data, ok := resp.Data.(map[string]interface{})
				if !ok {
					t.Errorf("expected data in response")
				} else if result, ok := data["result"].(float64); !ok || result != tt.expectedResult {
					t.Errorf("result = %v, want %v", result, tt.expectedResult)
				}
			}
		})
	}
}

func TestSubtractHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           *models.MathRequest
		expectedStatus int
		expectedResult float64
		expectError    bool
	}{
		{
			name:           "valid subtraction",
			method:         http.MethodPost,
			body:           &models.MathRequest{A: 5, B: 3},
			expectedStatus: http.StatusOK,
			expectedResult: 2,
			expectError:    false,
		},
		{
			name:           "negative result",
			method:         http.MethodPost,
			body:           &models.MathRequest{A: 3, B: 5},
			expectedStatus: http.StatusOK,
			expectedResult: -2,
			expectError:    false,
		},
		{
			name:           "decimals",
			method:         http.MethodPost,
			body:           &models.MathRequest{A: 5.5, B: 2.5},
			expectedStatus: http.StatusOK,
			expectedResult: 3,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(tt.method, "/api/math/subtract", bytes.NewReader(body))
			ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-123")
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			SubtractHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.expectedStatus)
			}

			if !tt.expectError {
				var resp models.SuccessResponse
				json.NewDecoder(w.Body).Decode(&resp)
				data, ok := resp.Data.(map[string]interface{})
				if !ok {
					t.Errorf("expected data in response")
				} else if result, ok := data["result"].(float64); !ok || result != tt.expectedResult {
					t.Errorf("result = %v, want %v", result, tt.expectedResult)
				}
			}
		})
	}
}

func TestMultiplyHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           *models.MathRequest
		expectedStatus int
		expectedResult float64
		expectError    bool
	}{
		{
			name:           "valid multiplication",
			method:         http.MethodPost,
			body:           &models.MathRequest{A: 5, B: 3},
			expectedStatus: http.StatusOK,
			expectedResult: 15,
			expectError:    false,
		},
		{
			name:           "negative numbers",
			method:         http.MethodPost,
			body:           &models.MathRequest{A: -5, B: 3},
			expectedStatus: http.StatusOK,
			expectedResult: -15,
			expectError:    false,
		},
		{
			name:           "zero",
			method:         http.MethodPost,
			body:           &models.MathRequest{A: 0, B: 5},
			expectedStatus: http.StatusOK,
			expectedResult: 0,
			expectError:    false,
		},
		{
			name:           "decimals",
			method:         http.MethodPost,
			body:           &models.MathRequest{A: 2.5, B: 4},
			expectedStatus: http.StatusOK,
			expectedResult: 10,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(tt.method, "/api/math/multiply", bytes.NewReader(body))
			ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-123")
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			MultiplyHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.expectedStatus)
			}

			if !tt.expectError {
				var resp models.SuccessResponse
				json.NewDecoder(w.Body).Decode(&resp)
				data, ok := resp.Data.(map[string]interface{})
				if !ok {
					t.Errorf("expected data in response")
				} else if result, ok := data["result"].(float64); !ok || result != tt.expectedResult {
					t.Errorf("result = %v, want %v", result, tt.expectedResult)
				}
			}
		})
	}
}

func TestDivideHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           *models.MathRequest
		expectedStatus int
		expectedResult float64
		expectError    bool
		expectedCode   string
	}{
		{
			name:           "valid division",
			method:         http.MethodPost,
			body:           &models.MathRequest{A: 6, B: 2},
			expectedStatus: http.StatusOK,
			expectedResult: 3,
			expectError:    false,
		},
		{
			name:           "division by zero",
			method:         http.MethodPost,
			body:           &models.MathRequest{A: 5, B: 0},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
			expectedCode:   "DIVISION_BY_ZERO",
		},
		{
			name:           "decimals",
			method:         http.MethodPost,
			body:           &models.MathRequest{A: 5, B: 2},
			expectedStatus: http.StatusOK,
			expectedResult: 2.5,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(tt.method, "/api/math/divide", bytes.NewReader(body))
			ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-123")
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			DivideHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.expectedStatus)
			}

			if tt.expectError {
				var resp models.APIErrorResponse
				json.NewDecoder(w.Body).Decode(&resp)
				if resp.Code != tt.expectedCode {
					t.Errorf("error code = %s, want %s", resp.Code, tt.expectedCode)
				}
			} else {
				var resp models.SuccessResponse
				json.NewDecoder(w.Body).Decode(&resp)
				data, ok := resp.Data.(map[string]interface{})
				if !ok {
					t.Errorf("expected data in response")
				} else if result, ok := data["result"].(float64); !ok || result != tt.expectedResult {
					t.Errorf("result = %v, want %v", result, tt.expectedResult)
				}
			}
		})
	}
}

func TestMethodNotAllowed(t *testing.T) {
	handlers := []func(http.ResponseWriter, *http.Request){
		AddHandler,
		SubtractHandler,
		MultiplyHandler,
		DivideHandler,
	}

	methods := []string{http.MethodGet, http.MethodPut, http.MethodDelete}

	for _, handler := range handlers {
		for _, method := range methods {
			req := httptest.NewRequest(method, "/api/math/add", nil)
			ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-123")
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler(w, req)

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("method %s: status = %d, want %d", method, w.Code, http.StatusMethodNotAllowed)
			}

			var resp models.APIErrorResponse
			json.NewDecoder(w.Body).Decode(&resp)
			if resp.Code != "METHOD_NOT_ALLOWED" {
				t.Errorf("error code = %s, want METHOD_NOT_ALLOWED", resp.Code)
			}
		}
	}
}

func TestInvalidJSON(t *testing.T) {
	invalidJSON := io.NopCloser(bytes.NewReader([]byte("invalid json")))
	req := httptest.NewRequest(http.MethodPost, "/api/math/add", invalidJSON)
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-123")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	AddHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var resp models.APIErrorResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Code != "INVALID_INPUT" {
		t.Errorf("error code = %s, want INVALID_INPUT", resp.Code)
	}
}

// TestErrorCodeStructure verifies that error responses have the correct structure.
func TestErrorCodeStructure(t *testing.T) {
	tests := []struct {
		name         string
		handler      func(http.ResponseWriter, *http.Request)
		body         *models.MathRequest
		expectedCode string
	}{
		{
			name:         "add with invalid JSON",
			handler:      AddHandler,
			body:         nil,
			expectedCode: "INVALID_INPUT",
		},
		{
			name:         "divide by zero",
			handler:      DivideHandler,
			body:         &models.MathRequest{A: 5, B: 0},
			expectedCode: "DIVISION_BY_ZERO",
		},
		{
			name:         "method not allowed",
			handler:      AddHandler,
			body:         &models.MathRequest{A: 5, B: 3},
			expectedCode: "METHOD_NOT_ALLOWED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			switch tt.name {
			case "add with invalid JSON":
				req = httptest.NewRequest(http.MethodPost, "/api/math/add", io.NopCloser(bytes.NewReader([]byte("bad"))))
			case "method not allowed":
				body, _ := json.Marshal(tt.body)
				req = httptest.NewRequest(http.MethodGet, "/api/math/add", bytes.NewReader(body))
			default:
				body, _ := json.Marshal(tt.body)
				req = httptest.NewRequest(http.MethodPost, "/api/math/divide", bytes.NewReader(body))
			}

			ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-123")
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			tt.handler(w, req)

			var resp models.APIErrorResponse
			json.NewDecoder(w.Body).Decode(&resp)

			if resp.Code != tt.expectedCode {
				t.Errorf("code = %s, want %s", resp.Code, tt.expectedCode)
			}
			if resp.Message == "" {
				t.Errorf("message is empty")
			}
			if resp.Timestamp == "" {
				t.Errorf("timestamp is empty")
			}
			if resp.RequestID != "test-123" {
				t.Errorf("request_id = %s, want test-123", resp.RequestID)
			}
		})
	}
}

// TestSuccessResponseStructure verifies that success responses have the correct structure.
func TestSuccessResponseStructure(t *testing.T) {
	body, _ := json.Marshal(models.MathRequest{A: 5, B: 3})
	req := httptest.NewRequest(http.MethodPost, "/api/math/add", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-456")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	AddHandler(w, req)

	var resp models.SuccessResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Data == nil {
		t.Errorf("data is nil")
	}
	if resp.Timestamp == "" {
		t.Errorf("timestamp is empty")
	}
	if resp.RequestID != "test-456" {
		t.Errorf("request_id = %s, want test-456", resp.RequestID)
	}
}
