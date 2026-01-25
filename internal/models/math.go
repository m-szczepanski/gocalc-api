package models

import "time"

type TimeProvider interface {
	Now() time.Time
}

type RealTimeProvider struct{}

func (RealTimeProvider) Now() time.Time {
	return time.Now().UTC()
}

var defaultTimeProvider TimeProvider = RealTimeProvider{}

func SetTimeProvider(tp TimeProvider) {
	defaultTimeProvider = tp
}

type MathRequest struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

type MathResponse struct {
	Result float64 `json:"result"`
}

type APIErrorResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Details   string `json:"details,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp"`
}

func NewAPIErrorResponse(code, message, details, requestID string) *APIErrorResponse {
	return &APIErrorResponse{
		Code:      code,
		Message:   message,
		Details:   details,
		RequestID: requestID,
		Timestamp: defaultTimeProvider.Now().Format(time.RFC3339),
	}
}

type SuccessResponse struct {
	Data      interface{} `json:"data"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp string      `json:"timestamp"`
}

func NewSuccessResponse(data interface{}, requestID string) *SuccessResponse {
	return &SuccessResponse{
		Data:      data,
		RequestID: requestID,
		Timestamp: defaultTimeProvider.Now().Format(time.RFC3339),
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}
