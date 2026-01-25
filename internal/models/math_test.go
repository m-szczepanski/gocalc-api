package models

import (
	"testing"
	"time"
)

type MockTimeProvider struct {
	mockTime time.Time
}

func (m MockTimeProvider) Now() time.Time {
	return m.mockTime
}

func TestTimeProvider(t *testing.T) {
	originalProvider := defaultTimeProvider
	defer func() { defaultTimeProvider = originalProvider }()

	mockTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	SetTimeProvider(MockTimeProvider{mockTime: mockTime})

	resp := NewAPIErrorResponse("TEST_CODE", "test message", "details", "req-123")
	expectedTimestamp := mockTime.Format(time.RFC3339)
	if resp.Timestamp != expectedTimestamp {
		t.Errorf("expected timestamp %s, got %s", expectedTimestamp, resp.Timestamp)
	}

	successResp := NewSuccessResponse("data", "req-456")
	if successResp.Timestamp != expectedTimestamp {
		t.Errorf("expected timestamp %s, got %s", expectedTimestamp, successResp.Timestamp)
	}
}

func TestRealTimeProvider(t *testing.T) {
	provider := RealTimeProvider{}
	now := provider.Now()

	actualNow := time.Now().UTC()
	diff := actualNow.Sub(now)
	if diff < 0 {
		diff = -diff
	}

	if diff > time.Second {
		t.Errorf("RealTimeProvider.Now() returned time too far from actual time: diff=%v", diff)
	}
}

func TestAPIErrorResponse(t *testing.T) {
	resp := NewAPIErrorResponse("ERROR_CODE", "error message", "error details", "request-123")

	if resp.Code != "ERROR_CODE" {
		t.Errorf("expected code ERROR_CODE, got %s", resp.Code)
	}
	if resp.Message != "error message" {
		t.Errorf("expected message 'error message', got %s", resp.Message)
	}
	if resp.Details != "error details" {
		t.Errorf("expected details 'error details', got %s", resp.Details)
	}
	if resp.RequestID != "request-123" {
		t.Errorf("expected request ID 'request-123', got %s", resp.RequestID)
	}
	if resp.Timestamp == "" {
		t.Error("expected timestamp, got empty")
	}
}

func TestSuccessResponse(t *testing.T) {
	data := map[string]interface{}{"result": 42}
	resp := NewSuccessResponse(data, "request-456")

	if resp.Data == nil {
		t.Error("expected data, got nil")
	}
	if resp.RequestID != "request-456" {
		t.Errorf("expected request ID 'request-456', got %s", resp.RequestID)
	}
	if resp.Timestamp == "" {
		t.Error("expected timestamp, got empty")
	}
}

func TestMathRequest(t *testing.T) {
	req := MathRequest{A: 5.5, B: 3.2}

	if req.A != 5.5 {
		t.Errorf("expected A=5.5, got %f", req.A)
	}
	if req.B != 3.2 {
		t.Errorf("expected B=3.2, got %f", req.B)
	}
}

func TestMathResponse(t *testing.T) {
	resp := MathResponse{Result: 8.7}

	if resp.Result != 8.7 {
		t.Errorf("expected Result=8.7, got %f", resp.Result)
	}
}
