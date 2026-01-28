package models

import (
	"encoding/json"
	"testing"
)

func TestVATRequestJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected VATRequest
		wantErr  bool
	}{
		{
			name:  "valid VAT request exclusive",
			input: `{"amount": 100, "rate": 23, "inclusive": false}`,
			expected: VATRequest{
				Amount:    100,
				Rate:      23,
				Inclusive: false,
			},
			wantErr: false,
		},
		{
			name:  "valid VAT request inclusive",
			input: `{"amount": 123, "rate": 23, "inclusive": true}`,
			expected: VATRequest{
				Amount:    123,
				Rate:      23,
				Inclusive: true,
			},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   `{"amount": 100, "rate": }`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req VATRequest
			err := json.Unmarshal([]byte(tt.input), &req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if req.Amount != tt.expected.Amount {
					t.Errorf("Amount = %v, want %v", req.Amount, tt.expected.Amount)
				}
				if req.Rate != tt.expected.Rate {
					t.Errorf("Rate = %v, want %v", req.Rate, tt.expected.Rate)
				}
				if req.Inclusive != tt.expected.Inclusive {
					t.Errorf("Inclusive = %v, want %v", req.Inclusive, tt.expected.Inclusive)
				}
			}
		})
	}
}

func TestVATResponseJSON(t *testing.T) {
	response := VATResponse{
		VATAmount:   23,
		NetAmount:   100,
		GrossAmount: 123,
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var decoded VATResponse
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if decoded.VATAmount != response.VATAmount {
		t.Errorf("VATAmount = %v, want %v", decoded.VATAmount, response.VATAmount)
	}
	if decoded.NetAmount != response.NetAmount {
		t.Errorf("NetAmount = %v, want %v", decoded.NetAmount, response.NetAmount)
	}
	if decoded.GrossAmount != response.GrossAmount {
		t.Errorf("GrossAmount = %v, want %v", decoded.GrossAmount, response.GrossAmount)
	}
}

func TestCompoundInterestRequestJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected CompoundInterestRequest
		wantErr  bool
	}{
		{
			name:  "valid compound interest request",
			input: `{"principal": 1000, "rate": 5, "time": 10, "compound_frequency": 12}`,
			expected: CompoundInterestRequest{
				Principal:         1000,
				Rate:              5,
				Time:              10,
				CompoundFrequency: 12,
			},
			wantErr: false,
		},
		{
			name:  "valid with decimal values",
			input: `{"principal": 1500.50, "rate": 3.5, "time": 2.5, "compound_frequency": 4}`,
			expected: CompoundInterestRequest{
				Principal:         1500.50,
				Rate:              3.5,
				Time:              2.5,
				CompoundFrequency: 4,
			},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   `{"principal": 1000, "rate": }`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req CompoundInterestRequest
			err := json.Unmarshal([]byte(tt.input), &req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if req.Principal != tt.expected.Principal {
					t.Errorf("Principal = %v, want %v", req.Principal, tt.expected.Principal)
				}
				if req.Rate != tt.expected.Rate {
					t.Errorf("Rate = %v, want %v", req.Rate, tt.expected.Rate)
				}
				if req.Time != tt.expected.Time {
					t.Errorf("Time = %v, want %v", req.Time, tt.expected.Time)
				}
				if req.CompoundFrequency != tt.expected.CompoundFrequency {
					t.Errorf("CompoundFrequency = %v, want %v", req.CompoundFrequency, tt.expected.CompoundFrequency)
				}
			}
		})
	}
}

func TestCompoundInterestResponseJSON(t *testing.T) {
	response := CompoundInterestResponse{
		FinalAmount:    1647.01,
		InterestEarned: 647.01,
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var decoded CompoundInterestResponse
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if decoded.FinalAmount != response.FinalAmount {
		t.Errorf("FinalAmount = %v, want %v", decoded.FinalAmount, response.FinalAmount)
	}
	if decoded.InterestEarned != response.InterestEarned {
		t.Errorf("InterestEarned = %v, want %v", decoded.InterestEarned, response.InterestEarned)
	}
}

func TestLoanPaymentRequestJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected LoanPaymentRequest
		wantErr  bool
	}{
		{
			name:  "valid loan payment request",
			input: `{"principal": 300000, "annual_rate": 4.5, "years": 30, "payments_per_year": 12}`,
			expected: LoanPaymentRequest{
				Principal:       300000,
				AnnualRate:      4.5,
				Years:           30,
				PaymentsPerYear: 12,
			},
			wantErr: false,
		},
		{
			name:  "valid with decimal years",
			input: `{"principal": 25000, "annual_rate": 6, "years": 5.5, "payments_per_year": 12}`,
			expected: LoanPaymentRequest{
				Principal:       25000,
				AnnualRate:      6,
				Years:           5.5,
				PaymentsPerYear: 12,
			},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   `{"principal": 10000, "annual_rate": }`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req LoanPaymentRequest
			err := json.Unmarshal([]byte(tt.input), &req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if req.Principal != tt.expected.Principal {
					t.Errorf("Principal = %v, want %v", req.Principal, tt.expected.Principal)
				}
				if req.AnnualRate != tt.expected.AnnualRate {
					t.Errorf("AnnualRate = %v, want %v", req.AnnualRate, tt.expected.AnnualRate)
				}
				if req.Years != tt.expected.Years {
					t.Errorf("Years = %v, want %v", req.Years, tt.expected.Years)
				}
				if req.PaymentsPerYear != tt.expected.PaymentsPerYear {
					t.Errorf("PaymentsPerYear = %v, want %v", req.PaymentsPerYear, tt.expected.PaymentsPerYear)
				}
			}
		})
	}
}

func TestLoanPaymentResponseJSON(t *testing.T) {
	response := LoanPaymentResponse{
		PaymentAmount: 1520.06,
		TotalPayment:  547221.60,
		TotalInterest: 247221.60,
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var decoded LoanPaymentResponse
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if decoded.PaymentAmount != response.PaymentAmount {
		t.Errorf("PaymentAmount = %v, want %v", decoded.PaymentAmount, response.PaymentAmount)
	}
	if decoded.TotalPayment != response.TotalPayment {
		t.Errorf("TotalPayment = %v, want %v", decoded.TotalPayment, response.TotalPayment)
	}
	if decoded.TotalInterest != response.TotalInterest {
		t.Errorf("TotalInterest = %v, want %v", decoded.TotalInterest, response.TotalInterest)
	}
}
