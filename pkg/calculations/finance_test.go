package calculations

import (
	"math"
	"testing"
)

func TestCalculateVAT(t *testing.T) {
	tests := []struct {
		name          string
		amount        float64
		rate          float64
		inclusive     bool
		expectedVAT   float64
		expectedNet   float64
		expectedGross float64
		expectError   bool
	}{
		{
			name:          "add 23% VAT to 100",
			amount:        100.0,
			rate:          23.0,
			inclusive:     false,
			expectedVAT:   23.0,
			expectedNet:   100.0,
			expectedGross: 123.0,
			expectError:   false,
		},
		{
			name:          "add 20% VAT to 50",
			amount:        50.0,
			rate:          20.0,
			inclusive:     false,
			expectedVAT:   10.0,
			expectedNet:   50.0,
			expectedGross: 60.0,
			expectError:   false,
		},
		{
			name:          "add 0% VAT to 100",
			amount:        100.0,
			rate:          0.0,
			inclusive:     false,
			expectedVAT:   0.0,
			expectedNet:   100.0,
			expectedGross: 100.0,
			expectError:   false,
		},
		{
			name:          "extract 23% VAT from 123",
			amount:        123.0,
			rate:          23.0,
			inclusive:     true,
			expectedVAT:   23.0,
			expectedNet:   100.0,
			expectedGross: 123.0,
			expectError:   false,
		},
		{
			name:          "extract 20% VAT from 60",
			amount:        60.0,
			rate:          20.0,
			inclusive:     true,
			expectedVAT:   10.0,
			expectedNet:   50.0,
			expectedGross: 60.0,
			expectError:   false,
		},
		{
			name:          "extract 0% VAT from 100",
			amount:        100.0,
			rate:          0.0,
			inclusive:     true,
			expectedVAT:   0.0,
			expectedNet:   100.0,
			expectedGross: 100.0,
			expectError:   false,
		},
		{
			name:          "add 19% VAT to 99.99",
			amount:        99.99,
			rate:          19.0,
			inclusive:     false,
			expectedVAT:   18.9981,
			expectedNet:   99.99,
			expectedGross: 118.9881,
			expectError:   false,
		},
		{
			name:          "extract 5.5% VAT from 105.50",
			amount:        105.50,
			rate:          5.5,
			inclusive:     true,
			expectedVAT:   5.499999999999993, // Expected floating-point precision
			expectedNet:   100.00000000000001,
			expectedGross: 105.50,
			expectError:   false,
		},
		{
			name:          "zero amount",
			amount:        0.0,
			rate:          23.0,
			inclusive:     false,
			expectedVAT:   0.0,
			expectedNet:   0.0,
			expectedGross: 0.0,
			expectError:   false,
		},
		{
			name:        "negative amount",
			amount:      -100.0,
			rate:        23.0,
			inclusive:   false,
			expectError: true,
		},
		{
			name:        "negative rate",
			amount:      100.0,
			rate:        -5.0,
			inclusive:   false,
			expectError: true,
		},
		{
			name:          "large amount with VAT",
			amount:        1000000.0,
			rate:          23.0,
			inclusive:     false,
			expectedVAT:   230000.0,
			expectedNet:   1000000.0,
			expectedGross: 1230000.0,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vatAmount, netAmount, grossAmount, err := CalculateVAT(tt.amount, tt.rate, tt.inclusive)

			if (err != nil) != tt.expectError {
				t.Errorf("CalculateVAT() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if tt.expectError {
				return
			}

			if !almostEqual(vatAmount, tt.expectedVAT, 0.0001) {
				t.Errorf("CalculateVAT() vatAmount = %v, want %v", vatAmount, tt.expectedVAT)
			}
			if !almostEqual(netAmount, tt.expectedNet, 0.0001) {
				t.Errorf("CalculateVAT() netAmount = %v, want %v", netAmount, tt.expectedNet)
			}
			if !almostEqual(grossAmount, tt.expectedGross, 0.0001) {
				t.Errorf("CalculateVAT() grossAmount = %v, want %v", grossAmount, tt.expectedGross)
			}
		})
	}
}

func TestCalculateCompoundInterest(t *testing.T) {
	tests := []struct {
		name              string
		principal         float64
		rate              float64
		time              float64
		compoundFrequency int
		expectedFinal     float64
		expectedInterest  float64
		expectError       bool
	}{
		{
			name:              "1000 at 5% for 10 years, monthly compounding",
			principal:         1000.0,
			rate:              5.0,
			time:              10.0,
			compoundFrequency: 12,
			expectedFinal:     1647.01,
			expectedInterest:  647.01,
			expectError:       false,
		},
		{
			name:              "5000 at 3% for 5 years, annual compounding",
			principal:         5000.0,
			rate:              3.0,
			time:              5.0,
			compoundFrequency: 1,
			expectedFinal:     5796.37,
			expectedInterest:  796.37,
			expectError:       false,
		},
		{
			name:              "2000 at 7.5% for 3 years, quarterly compounding",
			principal:         2000.0,
			rate:              7.5,
			time:              3.0,
			compoundFrequency: 4,
			expectedFinal:     2499.43,
			expectedInterest:  499.43,
			expectError:       false,
		},
		{
			name:              "zero interest rate",
			principal:         1000.0,
			rate:              0.0,
			time:              5.0,
			compoundFrequency: 12,
			expectedFinal:     1000.0,
			expectedInterest:  0.0,
			expectError:       false,
		},
		{
			name:              "zero time period",
			principal:         1000.0,
			rate:              5.0,
			time:              0.0,
			compoundFrequency: 12,
			expectedFinal:     1000.0,
			expectedInterest:  0.0,
			expectError:       false,
		},
		{
			name:              "zero principal",
			principal:         0.0,
			rate:              5.0,
			time:              10.0,
			compoundFrequency: 12,
			expectedFinal:     0.0,
			expectedInterest:  0.0,
			expectError:       false,
		},
		{
			name:              "daily compounding",
			principal:         10000.0,
			rate:              4.5,
			time:              2.0,
			compoundFrequency: 365,
			expectedFinal:     10941.68,
			expectedInterest:  941.68,
			expectError:       false,
		},
		{
			name:              "negative principal",
			principal:         -1000.0,
			rate:              5.0,
			time:              10.0,
			compoundFrequency: 12,
			expectError:       true,
		},
		{
			name:              "negative rate",
			principal:         1000.0,
			rate:              -5.0,
			time:              10.0,
			compoundFrequency: 12,
			expectError:       true,
		},
		{
			name:              "negative time",
			principal:         1000.0,
			rate:              5.0,
			time:              -10.0,
			compoundFrequency: 12,
			expectError:       true,
		},
		{
			name:              "zero compound frequency",
			principal:         1000.0,
			rate:              5.0,
			time:              10.0,
			compoundFrequency: 0,
			expectError:       true,
		},
		{
			name:              "negative compound frequency",
			principal:         1000.0,
			rate:              5.0,
			time:              10.0,
			compoundFrequency: -12,
			expectError:       true,
		},
		{
			name:              "large principal",
			principal:         1000000.0,
			rate:              2.5,
			time:              20.0,
			compoundFrequency: 12,
			expectedFinal:     1647863.98,
			expectedInterest:  647863.98,
			expectError:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finalAmount, interestEarned, err := CalculateCompoundInterest(
				tt.principal,
				tt.rate,
				tt.time,
				tt.compoundFrequency,
			)

			if (err != nil) != tt.expectError {
				t.Errorf("CalculateCompoundInterest() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if tt.expectError {
				return
			}

			if !almostEqual(finalAmount, tt.expectedFinal, 0.01) {
				t.Errorf("CalculateCompoundInterest() finalAmount = %v, want %v", finalAmount, tt.expectedFinal)
			}
			if !almostEqual(interestEarned, tt.expectedInterest, 0.01) {
				t.Errorf("CalculateCompoundInterest() interestEarned = %v, want %v", interestEarned, tt.expectedInterest)
			}
		})
	}
}

func TestCalculateLoanPayment(t *testing.T) {
	tests := []struct {
		name             string
		principal        float64
		annualRate       float64
		years            float64
		paymentsPerYear  int
		expectedPayment  float64
		expectedTotal    float64
		expectedInterest float64
		expectError      bool
	}{
		{
			name:             "30-year mortgage at 4.5%",
			principal:        300000.0,
			annualRate:       4.5,
			years:            30.0,
			paymentsPerYear:  12,
			expectedPayment:  1520.06,
			expectedTotal:    547221.60,
			expectedInterest: 247221.60,
			expectError:      false,
		},
		{
			name:             "car loan 5 years at 6%",
			principal:        25000.0,
			annualRate:       6.0,
			years:            5.0,
			paymentsPerYear:  12,
			expectedPayment:  483.32,
			expectedTotal:    28999.20,
			expectedInterest: 3999.20,
			expectError:      false,
		},
		{
			name:             "personal loan 3 years at 9%",
			principal:        10000.0,
			annualRate:       9.0,
			years:            3.0,
			paymentsPerYear:  12,
			expectedPayment:  318.00,
			expectedTotal:    11448.00,
			expectedInterest: 1448.00,
			expectError:      false,
		},
		{
			name:             "zero interest loan",
			principal:        10000.0,
			annualRate:       0.0,
			years:            5.0,
			paymentsPerYear:  12,
			expectedPayment:  166.67,
			expectedTotal:    10000.0,
			expectedInterest: 0.0,
			expectError:      false,
		},
		{
			name:             "quarterly payments",
			principal:        50000.0,
			annualRate:       5.0,
			years:            10.0,
			paymentsPerYear:  4,
			expectedPayment:  1596.07,
			expectedTotal:    63842.80,
			expectedInterest: 13842.80,
			expectError:      false,
		},
		{
			name:             "short-term loan 1 year",
			principal:        5000.0,
			annualRate:       12.0,
			years:            1.0,
			paymentsPerYear:  12,
			expectedPayment:  444.24,
			expectedTotal:    5330.88,
			expectedInterest: 330.88,
			expectError:      false,
		},
		{
			name:            "negative principal",
			principal:       -10000.0,
			annualRate:      5.0,
			years:           5.0,
			paymentsPerYear: 12,
			expectError:     true,
		},
		{
			name:            "negative rate",
			principal:       10000.0,
			annualRate:      -5.0,
			years:           5.0,
			paymentsPerYear: 12,
			expectError:     true,
		},
		{
			name:            "zero years",
			principal:       10000.0,
			annualRate:      5.0,
			years:           0.0,
			paymentsPerYear: 12,
			expectError:     true,
		},
		{
			name:            "negative years",
			principal:       10000.0,
			annualRate:      5.0,
			years:           -5.0,
			paymentsPerYear: 12,
			expectError:     true,
		},
		{
			name:            "zero payments per year",
			principal:       10000.0,
			annualRate:      5.0,
			years:           5.0,
			paymentsPerYear: 0,
			expectError:     true,
		},
		{
			name:            "negative payments per year",
			principal:       10000.0,
			annualRate:      5.0,
			years:           5.0,
			paymentsPerYear: -12,
			expectError:     true,
		},
		{
			name:             "large mortgage",
			principal:        1000000.0,
			annualRate:       3.5,
			years:            30.0,
			paymentsPerYear:  12,
			expectedPayment:  4490.45,
			expectedTotal:    1616562.00,
			expectedInterest: 616562.00,
			expectError:      false,
		},
		{
			name:             "small loan with decimals",
			principal:        1500.50,
			annualRate:       7.25,
			years:            2.0,
			paymentsPerYear:  12,
			expectedPayment:  67.35,
			expectedTotal:    1616.40,
			expectedInterest: 115.90,
			expectError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paymentAmount, totalPayment, totalInterest, err := CalculateLoanPayment(
				tt.principal,
				tt.annualRate,
				tt.years,
				tt.paymentsPerYear,
			)

			if (err != nil) != tt.expectError {
				t.Errorf("CalculateLoanPayment() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if tt.expectError {
				return
			}

			if !almostEqual(paymentAmount, tt.expectedPayment, 0.01) {
				t.Errorf("CalculateLoanPayment() paymentAmount = %v, want %v", paymentAmount, tt.expectedPayment)
			}
			if !almostEqual(totalPayment, tt.expectedTotal, 0.50) {
				t.Errorf("CalculateLoanPayment() totalPayment = %v, want %v", totalPayment, tt.expectedTotal)
			}
			if !almostEqual(totalInterest, tt.expectedInterest, 0.50) {
				t.Errorf("CalculateLoanPayment() totalInterest = %v, want %v", totalInterest, tt.expectedInterest)
			}
		})
	}
}

func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
