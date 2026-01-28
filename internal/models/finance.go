package models

type VATRequest struct {
	Amount    float64 `json:"amount"`
	Rate      float64 `json:"rate"`      // VAT rate as percentage (e.g., 23 for 23%)
	Inclusive bool    `json:"inclusive"` // true: extract VAT from amount, false: add VAT to amount
}

type VATResponse struct {
	VATAmount   float64 `json:"vat_amount"`
	NetAmount   float64 `json:"net_amount"`
	GrossAmount float64 `json:"gross_amount"`
}

type CompoundInterestRequest struct {
	Principal         float64 `json:"principal"`          // Initial principal amount
	Rate              float64 `json:"rate"`               // Annual interest rate as percentage (e.g., 5 for 5%)
	Time              float64 `json:"time"`               // Time period in years
	CompoundFrequency int     `json:"compound_frequency"` // Number of times interest is compounded per year (e.g., 12 for monthly)
}

type CompoundInterestResponse struct {
	FinalAmount    float64 `json:"final_amount"`
	InterestEarned float64 `json:"interest_earned"`
}

type LoanPaymentRequest struct {
	Principal       float64 `json:"principal"`         // Loan principal amount
	AnnualRate      float64 `json:"annual_rate"`       // Annual interest rate as percentage (e.g., 5 for 5%)
	Years           float64 `json:"years"`             // Loan term in years
	PaymentsPerYear int     `json:"payments_per_year"` // Number of payments per year (e.g., 12 for monthly)
}

type LoanPaymentResponse struct {
	PaymentAmount float64 `json:"payment_amount"`
	TotalPayment  float64 `json:"total_payment"`
	TotalInterest float64 `json:"total_interest"`
}
