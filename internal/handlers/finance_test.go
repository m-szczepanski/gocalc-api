package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/m-szczepanski/gocalc-api/internal/middleware"
	"github.com/m-szczepanski/gocalc-api/internal/models"
)

func TestVATHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           *models.VATRequest
		expectedStatus int
		expectedVAT    float64
		expectedNet    float64
		expectedGross  float64
		expectError    bool
	}{
		{
			name:           "add 23% VAT to 100",
			method:         http.MethodPost,
			body:           &models.VATRequest{Amount: 100, Rate: 23, Inclusive: false},
			expectedStatus: http.StatusOK,
			expectedVAT:    23,
			expectedNet:    100,
			expectedGross:  123,
			expectError:    false,
		},
		{
			name:           "extract 23% VAT from 123",
			method:         http.MethodPost,
			body:           &models.VATRequest{Amount: 123, Rate: 23, Inclusive: true},
			expectedStatus: http.StatusOK,
			expectedVAT:    23,
			expectedNet:    100,
			expectedGross:  123,
			expectError:    false,
		},
		{
			name:           "add 20% VAT to 50",
			method:         http.MethodPost,
			body:           &models.VATRequest{Amount: 50, Rate: 20, Inclusive: false},
			expectedStatus: http.StatusOK,
			expectedVAT:    10,
			expectedNet:    50,
			expectedGross:  60,
			expectError:    false,
		},
		{
			name:           "zero VAT rate",
			method:         http.MethodPost,
			body:           &models.VATRequest{Amount: 100, Rate: 0, Inclusive: false},
			expectedStatus: http.StatusOK,
			expectedVAT:    0,
			expectedNet:    100,
			expectedGross:  100,
			expectError:    false,
		},
		{
			name:           "negative amount",
			method:         http.MethodPost,
			body:           &models.VATRequest{Amount: -100, Rate: 23, Inclusive: false},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "negative rate",
			method:         http.MethodPost,
			body:           &models.VATRequest{Amount: 100, Rate: -5, Inclusive: false},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "NaN amount",
			method:         http.MethodPost,
			body:           &models.VATRequest{Amount: math.NaN(), Rate: 23, Inclusive: false},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "method not allowed (GET)",
			method:         http.MethodGet,
			body:           &models.VATRequest{Amount: 100, Rate: 23, Inclusive: false},
			expectedStatus: http.StatusMethodNotAllowed,
			expectError:    true,
		},
		{
			name:           "invalid JSON body",
			method:         http.MethodPost,
			body:           nil,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != nil {
				body, _ := json.Marshal(tt.body)
				req = httptest.NewRequest(tt.method, "/api/finance/vat", bytes.NewReader(body))
			} else {
				req = httptest.NewRequest(tt.method, "/api/finance/vat", bytes.NewReader([]byte("invalid json")))
			}
			ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-123")
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			VATHandler(w, req)

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
					return
				}

				vatAmount, _ := data["vat_amount"].(float64)
				netAmount, _ := data["net_amount"].(float64)
				grossAmount, _ := data["gross_amount"].(float64)

				if !almostEqual(vatAmount, tt.expectedVAT, 0.01) {
					t.Errorf("vat_amount = %v, want %v", vatAmount, tt.expectedVAT)
				}
				if !almostEqual(netAmount, tt.expectedNet, 0.01) {
					t.Errorf("net_amount = %v, want %v", netAmount, tt.expectedNet)
				}
				if !almostEqual(grossAmount, tt.expectedGross, 0.01) {
					t.Errorf("gross_amount = %v, want %v", grossAmount, tt.expectedGross)
				}
			}
		})
	}
}

func TestCompoundInterestHandler(t *testing.T) {
	tests := []struct {
		name             string
		method           string
		body             *models.CompoundInterestRequest
		expectedStatus   int
		expectedFinal    float64
		expectedInterest float64
		expectError      bool
	}{
		{
			name:   "1000 at 5% for 10 years, monthly",
			method: http.MethodPost,
			body: &models.CompoundInterestRequest{
				Principal:         1000,
				Rate:              5,
				Time:              10,
				CompoundFrequency: 12,
			},
			expectedStatus:   http.StatusOK,
			expectedFinal:    1647.01,
			expectedInterest: 647.01,
			expectError:      false,
		},
		{
			name:   "5000 at 3% for 5 years, annual",
			method: http.MethodPost,
			body: &models.CompoundInterestRequest{
				Principal:         5000,
				Rate:              3,
				Time:              5,
				CompoundFrequency: 1,
			},
			expectedStatus:   http.StatusOK,
			expectedFinal:    5796.37,
			expectedInterest: 796.37,
			expectError:      false,
		},
		{
			name:   "zero interest rate",
			method: http.MethodPost,
			body: &models.CompoundInterestRequest{
				Principal:         1000,
				Rate:              0,
				Time:              5,
				CompoundFrequency: 12,
			},
			expectedStatus:   http.StatusOK,
			expectedFinal:    1000,
			expectedInterest: 0,
			expectError:      false,
		},
		{
			name:   "negative principal",
			method: http.MethodPost,
			body: &models.CompoundInterestRequest{
				Principal:         -1000,
				Rate:              5,
				Time:              10,
				CompoundFrequency: 12,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:   "negative rate",
			method: http.MethodPost,
			body: &models.CompoundInterestRequest{
				Principal:         1000,
				Rate:              -5,
				Time:              10,
				CompoundFrequency: 12,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:   "zero compound frequency",
			method: http.MethodPost,
			body: &models.CompoundInterestRequest{
				Principal:         1000,
				Rate:              5,
				Time:              10,
				CompoundFrequency: 0,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:   "NaN principal",
			method: http.MethodPost,
			body: &models.CompoundInterestRequest{
				Principal:         math.NaN(),
				Rate:              5,
				Time:              10,
				CompoundFrequency: 12,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:   "method not allowed (PUT)",
			method: http.MethodPut,
			body: &models.CompoundInterestRequest{
				Principal:         1000,
				Rate:              5,
				Time:              10,
				CompoundFrequency: 12,
			},
			expectedStatus: http.StatusMethodNotAllowed,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != nil {
				body, _ := json.Marshal(tt.body)
				req = httptest.NewRequest(tt.method, "/api/finance/compound-interest", bytes.NewReader(body))
			} else {
				req = httptest.NewRequest(tt.method, "/api/finance/compound-interest", nil)
			}
			ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-123")
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			CompoundInterestHandler(w, req)

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
					return
				}

				finalAmount, _ := data["final_amount"].(float64)
				interestEarned, _ := data["interest_earned"].(float64)

				if !almostEqual(finalAmount, tt.expectedFinal, 0.01) {
					t.Errorf("final_amount = %v, want %v", finalAmount, tt.expectedFinal)
				}
				if !almostEqual(interestEarned, tt.expectedInterest, 0.01) {
					t.Errorf("interest_earned = %v, want %v", interestEarned, tt.expectedInterest)
				}
			}
		})
	}
}

func TestLoanPaymentHandler(t *testing.T) {
	tests := []struct {
		name             string
		method           string
		body             *models.LoanPaymentRequest
		expectedStatus   int
		expectedPayment  float64
		expectedTotal    float64
		expectedInterest float64
		expectError      bool
	}{
		{
			name:   "30-year mortgage at 4.5%",
			method: http.MethodPost,
			body: &models.LoanPaymentRequest{
				Principal:       300000,
				AnnualRate:      4.5,
				Years:           30,
				PaymentsPerYear: 12,
			},
			expectedStatus:   http.StatusOK,
			expectedPayment:  1520.06,
			expectedTotal:    547221.60,
			expectedInterest: 247221.60,
			expectError:      false,
		},
		{
			name:   "car loan 5 years at 6%",
			method: http.MethodPost,
			body: &models.LoanPaymentRequest{
				Principal:       25000,
				AnnualRate:      6,
				Years:           5,
				PaymentsPerYear: 12,
			},
			expectedStatus:   http.StatusOK,
			expectedPayment:  483.32,
			expectedTotal:    28999.20,
			expectedInterest: 3999.20,
			expectError:      false,
		},
		{
			name:   "zero interest loan",
			method: http.MethodPost,
			body: &models.LoanPaymentRequest{
				Principal:       10000,
				AnnualRate:      0,
				Years:           5,
				PaymentsPerYear: 12,
			},
			expectedStatus:   http.StatusOK,
			expectedPayment:  166.67,
			expectedTotal:    10000.20,
			expectedInterest: 0,
			expectError:      false,
		},
		{
			name:   "negative principal",
			method: http.MethodPost,
			body: &models.LoanPaymentRequest{
				Principal:       -10000,
				AnnualRate:      5,
				Years:           5,
				PaymentsPerYear: 12,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:   "negative rate",
			method: http.MethodPost,
			body: &models.LoanPaymentRequest{
				Principal:       10000,
				AnnualRate:      -5,
				Years:           5,
				PaymentsPerYear: 12,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:   "zero years",
			method: http.MethodPost,
			body: &models.LoanPaymentRequest{
				Principal:       10000,
				AnnualRate:      5,
				Years:           0,
				PaymentsPerYear: 12,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:   "zero payments per year",
			method: http.MethodPost,
			body: &models.LoanPaymentRequest{
				Principal:       10000,
				AnnualRate:      5,
				Years:           5,
				PaymentsPerYear: 0,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:   "Inf principal",
			method: http.MethodPost,
			body: &models.LoanPaymentRequest{
				Principal:       math.Inf(1),
				AnnualRate:      5,
				Years:           5,
				PaymentsPerYear: 12,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:   "method not allowed (DELETE)",
			method: http.MethodDelete,
			body: &models.LoanPaymentRequest{
				Principal:       10000,
				AnnualRate:      5,
				Years:           5,
				PaymentsPerYear: 12,
			},
			expectedStatus: http.StatusMethodNotAllowed,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != nil {
				body, _ := json.Marshal(tt.body)
				req = httptest.NewRequest(tt.method, "/api/finance/loan-payment", bytes.NewReader(body))
			} else {
				req = httptest.NewRequest(tt.method, "/api/finance/loan-payment", nil)
			}
			ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-123")
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			LoanPaymentHandler(w, req)

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
					return
				}

				paymentAmount, _ := data["payment_amount"].(float64)
				totalPayment, _ := data["total_payment"].(float64)
				totalInterest, _ := data["total_interest"].(float64)

				if !almostEqual(paymentAmount, tt.expectedPayment, 0.01) {
					t.Errorf("payment_amount = %v, want %v", paymentAmount, tt.expectedPayment)
				}
				if !almostEqual(totalPayment, tt.expectedTotal, 0.50) {
					t.Errorf("total_payment = %v, want %v", totalPayment, tt.expectedTotal)
				}
				if !almostEqual(totalInterest, tt.expectedInterest, 0.50) {
					t.Errorf("total_interest = %v, want %v", totalInterest, tt.expectedInterest)
				}
			}
		})
	}
}

// almostEqual checks if two float64 values are approximately equal within a tolerance
func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
