package calculations

import (
	"fmt"
	"math"
)

// CalculateVAT calculates VAT tax for a given amount and rate.
// If inclusive is true, it extracts VAT from the amount (amount already includes VAT).
// If inclusive is false, it adds VAT to the amount (amount is net).
//
// Formula for VAT-exclusive (add VAT):
//
//	VAT Amount = Amount * (Rate / 100)
//	Gross Amount = Amount + VAT Amount
//
// Formula for VAT-inclusive (extract VAT):
//
//	Net Amount = Amount / (1 + Rate / 100)
//	VAT Amount = Amount - Net Amount
//
// Precision: Uses float64 arithmetic which may have rounding errors for very large
// amounts or many decimal places. For financial applications requiring exact decimal
// precision, consider using a decimal library.
//
// Returns: (vatAmount, netAmount, grossAmount, error)
func CalculateVAT(amount, rate float64, inclusive bool) (float64, float64, float64, error) {
	if amount < 0 {
		return 0, 0, 0, fmt.Errorf("amount cannot be negative")
	}
	if rate < 0 {
		return 0, 0, 0, fmt.Errorf("rate cannot be negative")
	}

	var vatAmount, netAmount, grossAmount float64

	if inclusive {
		netAmount = amount / (1 + rate/100)
		vatAmount = amount - netAmount
		grossAmount = amount
	} else {
		netAmount = amount
		vatAmount = amount * (rate / 100)
		grossAmount = amount + vatAmount
	}

	return vatAmount, netAmount, grossAmount, nil
}

// CalculateCompoundInterest calculates compound interest for an investment.
//
// Formula: A = P * (1 + r/n)^(n*t)
// Where:
//
//	A = Final amount
//	P = Principal (initial investment)
//	r = Annual interest rate (as decimal, e.g., 0.05 for 5%)
//	n = Number of times interest is compounded per year
//	t = Time in years
//
// Assumptions:
// - Interest is compounded at regular intervals (discrete compounding)
// - No additional deposits or withdrawals during the period
// - Rate remains constant throughout the investment period
//
// Precision: Uses float64 arithmetic. For very long time periods or high compound
// frequencies, floating-point errors may accumulate. Results are rounded to 2 decimal
// places for practical financial use.
//
// Returns: (finalAmount, interestEarned, error)
func CalculateCompoundInterest(principal, rate, time float64, compoundFrequency int) (float64, float64, error) {
	if principal < 0 {
		return 0, 0, fmt.Errorf("principal cannot be negative")
	}
	if rate < 0 {
		return 0, 0, fmt.Errorf("rate cannot be negative")
	}
	if time < 0 {
		return 0, 0, fmt.Errorf("time cannot be negative")
	}
	if compoundFrequency <= 0 {
		return 0, 0, fmt.Errorf("compound frequency must be positive")
	}

	rateDecimal := rate / 100

	// A = P * (1 + r/n)^(n*t)
	n := float64(compoundFrequency)
	exponent := n * time
	base := 1 + (rateDecimal / n)
	finalAmount := principal * math.Pow(base, exponent)
	interestEarned := finalAmount - principal

	finalAmount = math.Round(finalAmount*100) / 100
	interestEarned = math.Round(interestEarned*100) / 100

	return finalAmount, interestEarned, nil
}

// CalculateLoanPayment calculates the periodic payment amount for a loan using
// the standard amortization formula.
//
// Formula: M = P * [r(1+r)^n] / [(1+r)^n - 1]
// Where:
//
//	M = Payment amount per period
//	P = Principal loan amount
//	r = Interest rate per period (annual rate / payments per year / 100)
//	n = Total number of payments (years * payments per year)
//
// Assumptions:
// - Fixed interest rate throughout the loan term
// - Equal payment amounts (amortized loan)
// - Payments made at regular intervals
// - No additional fees, insurance, or prepayments
//
// Special case: If interest rate is 0, payment = principal / total payments
//
// Precision: Uses float64 arithmetic. Monthly payment calculations should be
// accurate for typical loan amounts and terms. For exact amortization schedules,
// consider using a decimal library.
//
// Returns: (paymentAmount, totalPayment, totalInterest, error)
func CalculateLoanPayment(principal, annualRate, years float64, paymentsPerYear int) (float64, float64, float64, error) {
	if principal < 0 {
		return 0, 0, 0, fmt.Errorf("principal cannot be negative")
	}
	if annualRate < 0 {
		return 0, 0, 0, fmt.Errorf("annual rate cannot be negative")
	}
	if years <= 0 {
		return 0, 0, 0, fmt.Errorf("years must be positive")
	}
	if paymentsPerYear <= 0 {
		return 0, 0, 0, fmt.Errorf("payments per year must be positive")
	}

	totalPayments := years * float64(paymentsPerYear)

	// Special case: zero interest rate
	if annualRate == 0 {
		paymentAmount := principal / totalPayments
		paymentAmount = math.Round(paymentAmount*100) / 100
		// For a 0% loan, totalPayment should exactly match the principal (to cents),
		// and totalInterest must be 0, regardless of per-period rounding.
		totalPayment := math.Round(principal*100) / 100
		totalInterest := 0.0
		return paymentAmount, totalPayment, totalInterest, nil
	}

	ratePerPeriod := (annualRate / 100) / float64(paymentsPerYear)

	// M = P * [r(1+r)^n] / [(1+r)^n - 1]
	onePlusR := 1 + ratePerPeriod
	powerN := math.Pow(onePlusR, totalPayments)
	paymentAmount := principal * (ratePerPeriod * powerN) / (powerN - 1)

	paymentAmount = math.Round(paymentAmount*100) / 100

	totalPayment := paymentAmount * totalPayments
	totalInterest := totalPayment - principal

	totalPayment = math.Round(totalPayment*100) / 100
	totalInterest = math.Round(totalInterest*100) / 100

	return paymentAmount, totalPayment, totalInterest, nil
}
