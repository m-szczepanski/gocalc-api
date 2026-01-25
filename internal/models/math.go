package models

type MathRequest struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

type MathResponse struct {
	Result float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
