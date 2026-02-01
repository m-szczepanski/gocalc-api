package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	HealthHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var response struct {
		Data HealthResponse `json:"data"`
	}

	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Data.Status != "healthy" {
		t.Errorf("expected status 'healthy', got %s", response.Data.Status)
	}

	if response.Data.Version == "" {
		t.Error("expected version to be set")
	}

	if response.Data.Uptime == "" {
		t.Error("expected uptime to be set")
	}

	if response.Data.Timestamp == "" {
		t.Error("expected timestamp to be set")
	}

	if response.Data.Checks == nil {
		t.Error("expected checks to be set")
	}

	if _, ok := response.Data.Checks["memory"]; !ok {
		t.Error("expected memory check")
	}

	if _, ok := response.Data.Checks["goroutines"]; !ok {
		t.Error("expected goroutines check")
	}
}

func TestReadinessHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()

	ReadinessHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var response struct {
		Data ReadinessResponse `json:"data"`
	}

	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !response.Data.Ready {
		t.Error("expected ready to be true")
	}

	if response.Data.Status != "ready" {
		t.Errorf("expected status 'ready', got %s", response.Data.Status)
	}
}

func TestCheckMemory(t *testing.T) {
	result := checkMemory()
	if result != "ok" && result != "warning" {
		t.Errorf("expected 'ok' or 'warning', got %s", result)
	}
}

func TestCheckGoroutines(t *testing.T) {
	result := checkGoroutines()
	if result != "ok" && result != "warning" {
		t.Errorf("expected 'ok' or 'warning', got %s", result)
	}
}

func TestFormatMemoryCheck(t *testing.T) {
	tests := []struct {
		name  string
		alloc uint64
		sys   uint64
		want  string
	}{
		{
			name:  "normal usage",
			alloc: 100,
			sys:   500,
			want:  "ok",
		},
		{
			name:  "high alloc",
			alloc: 600,
			sys:   500,
			want:  "warning",
		},
		{
			name:  "high sys",
			alloc: 100,
			sys:   1100,
			want:  "warning",
		},
		{
			name:  "both high",
			alloc: 600,
			sys:   1100,
			want:  "warning",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatMemoryCheck(tt.alloc, tt.sys)
			if got != tt.want {
				t.Errorf("formatMemoryCheck(%d, %d) = %s, want %s", tt.alloc, tt.sys, got, tt.want)
			}
		})
	}
}
