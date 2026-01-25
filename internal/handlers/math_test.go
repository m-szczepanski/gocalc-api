package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/m-szczepanski/gocalc-api/internal/models"
)

func TestAddHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           models.MathRequest
		expectedStatus int
		expectedResult float64
	}{
		{
			name:           "valid addition",
			method:         http.MethodPost,
			body:           models.MathRequest{A: 5, B: 3},
			expectedStatus: http.StatusOK,
			expectedResult: 8,
		},
		{
			name:           "negative numbers",
			method:         http.MethodPost,
			body:           models.MathRequest{A: -5, B: 3},
			expectedStatus: http.StatusOK,
			expectedResult: -2,
		},
		{
			name:           "decimals",
			method:         http.MethodPost,
			body:           models.MathRequest{A: 1.5, B: 2.5},
			expectedStatus: http.StatusOK,
			expectedResult: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(tt.method, "/api/math/add", bytes.NewReader(body))
			w := httptest.NewRecorder()

			AddHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.expectedStatus)
			}

			var resp models.MathResponse
			json.NewDecoder(w.Body).Decode(&resp)
			if resp.Result != tt.expectedResult {
				t.Errorf("result = %v, want %v", resp.Result, tt.expectedResult)
			}
		})
	}
}

func TestSubtractHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           models.MathRequest
		expectedStatus int
		expectedResult float64
	}{
		{
			name:           "valid subtraction",
			method:         http.MethodPost,
			body:           models.MathRequest{A: 5, B: 3},
			expectedStatus: http.StatusOK,
			expectedResult: 2,
		},
		{
			name:           "negative result",
			method:         http.MethodPost,
			body:           models.MathRequest{A: 3, B: 5},
			expectedStatus: http.StatusOK,
			expectedResult: -2,
		},
		{
			name:           "decimals",
			method:         http.MethodPost,
			body:           models.MathRequest{A: 5.5, B: 2.5},
			expectedStatus: http.StatusOK,
			expectedResult: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(tt.method, "/api/math/subtract", bytes.NewReader(body))
			w := httptest.NewRecorder()

			SubtractHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.expectedStatus)
			}

			var resp models.MathResponse
			json.NewDecoder(w.Body).Decode(&resp)
			if resp.Result != tt.expectedResult {
				t.Errorf("result = %v, want %v", resp.Result, tt.expectedResult)
			}
		})
	}
}

func TestMultiplyHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           models.MathRequest
		expectedStatus int
		expectedResult float64
	}{
		{
			name:           "valid multiplication",
			method:         http.MethodPost,
			body:           models.MathRequest{A: 5, B: 3},
			expectedStatus: http.StatusOK,
			expectedResult: 15,
		},
		{
			name:           "negative numbers",
			method:         http.MethodPost,
			body:           models.MathRequest{A: -5, B: 3},
			expectedStatus: http.StatusOK,
			expectedResult: -15,
		},
		{
			name:           "zero",
			method:         http.MethodPost,
			body:           models.MathRequest{A: 0, B: 5},
			expectedStatus: http.StatusOK,
			expectedResult: 0,
		},
		{
			name:           "decimals",
			method:         http.MethodPost,
			body:           models.MathRequest{A: 2.5, B: 4},
			expectedStatus: http.StatusOK,
			expectedResult: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(tt.method, "/api/math/multiply", bytes.NewReader(body))
			w := httptest.NewRecorder()

			MultiplyHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.expectedStatus)
			}

			var resp models.MathResponse
			json.NewDecoder(w.Body).Decode(&resp)
			if resp.Result != tt.expectedResult {
				t.Errorf("result = %v, want %v", resp.Result, tt.expectedResult)
			}
		})
	}
}

func TestDivideHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           models.MathRequest
		expectedStatus int
		expectedResult float64
		checkError     bool
	}{
		{
			name:           "valid division",
			method:         http.MethodPost,
			body:           models.MathRequest{A: 6, B: 2},
			expectedStatus: http.StatusOK,
			expectedResult: 3,
			checkError:     false,
		},
		{
			name:           "division by zero",
			method:         http.MethodPost,
			body:           models.MathRequest{A: 5, B: 0},
			expectedStatus: http.StatusBadRequest,
			expectedResult: 0,
			checkError:     true,
		},
		{
			name:           "decimals",
			method:         http.MethodPost,
			body:           models.MathRequest{A: 5, B: 2},
			expectedStatus: http.StatusOK,
			expectedResult: 2.5,
			checkError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(tt.method, "/api/math/divide", bytes.NewReader(body))
			w := httptest.NewRecorder()

			DivideHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.expectedStatus)
			}

			if tt.checkError {
				var resp models.ErrorResponse
				json.NewDecoder(w.Body).Decode(&resp)
				if resp.Error == "" {
					t.Errorf("expected error message, got empty")
				}
			} else {
				var resp models.MathResponse
				json.NewDecoder(w.Body).Decode(&resp)
				if resp.Result != tt.expectedResult {
					t.Errorf("result = %v, want %v", resp.Result, tt.expectedResult)
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
			w := httptest.NewRecorder()

			handler(w, req)

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("method %s: status = %d, want %d", method, w.Code, http.StatusMethodNotAllowed)
			}
		}
	}
}

func TestInvalidJSON(t *testing.T) {
	invalidJSON := io.NopCloser(bytes.NewReader([]byte("invalid json")))
	req := httptest.NewRequest(http.MethodPost, "/api/math/add", invalidJSON)
	w := httptest.NewRecorder()

	AddHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var resp models.ErrorResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Error == "" {
		t.Errorf("expected error message, got empty")
	}
}
