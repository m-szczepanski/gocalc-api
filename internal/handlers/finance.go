package handlers

import (
	"net/http"

	apierrors "github.com/m-szczepanski/gocalc-api/internal/errors"
	"github.com/m-szczepanski/gocalc-api/internal/models"
	"github.com/m-szczepanski/gocalc-api/internal/validation"
	"github.com/m-szczepanski/gocalc-api/pkg/calculations"
)

func VATHandler(w http.ResponseWriter, r *http.Request) {
	if err := validation.ValidateMethod(r.Method, http.MethodPost); err != nil {
		writeErrorWithDetails(w, r, err)
		return
	}

	var req models.VATRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		writeErrorWithDetails(w, r, apierrors.InvalidInput("invalid request body").WithError(err))
		return
	}

	if err := validation.ValidateVATRequest(&req); err != nil {
		writeErrorWithDetails(w, r, err)
		return
	}

	vatAmount, netAmount, grossAmount, err := calculations.CalculateVAT(req.Amount, req.Rate, req.Inclusive)
	if err != nil {
		writeErrorWithDetails(w, r, apierrors.InternalError("VAT calculation failed").WithError(err))
		return
	}

	response := models.VATResponse{
		VATAmount:   vatAmount,
		NetAmount:   netAmount,
		GrossAmount: grossAmount,
	}

	if err := writeSuccessResponse(w, r, response); err != nil {
		// Error already logged, headers likely already sent
		return
	}
}

func CompoundInterestHandler(w http.ResponseWriter, r *http.Request) {
	if err := validation.ValidateMethod(r.Method, http.MethodPost); err != nil {
		writeErrorWithDetails(w, r, err)
		return
	}

	var req models.CompoundInterestRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		writeErrorWithDetails(w, r, apierrors.InvalidInput("invalid request body").WithError(err))
		return
	}

	if err := validation.ValidateCompoundInterestRequest(&req); err != nil {
		writeErrorWithDetails(w, r, err)
		return
	}

	finalAmount, interestEarned, err := calculations.CalculateCompoundInterest(
		req.Principal,
		req.Rate,
		req.Time,
		req.CompoundFrequency,
	)
	if err != nil {
		writeErrorWithDetails(w, r, apierrors.InternalError("compound interest calculation failed").WithError(err))
		return
	}

	response := models.CompoundInterestResponse{
		FinalAmount:    finalAmount,
		InterestEarned: interestEarned,
	}

	if err := writeSuccessResponse(w, r, response); err != nil {
		// Error already logged, headers likely already sent
		return
	}
}

func LoanPaymentHandler(w http.ResponseWriter, r *http.Request) {
	if err := validation.ValidateMethod(r.Method, http.MethodPost); err != nil {
		writeErrorWithDetails(w, r, err)
		return
	}

	var req models.LoanPaymentRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		writeErrorWithDetails(w, r, apierrors.InvalidInput("invalid request body").WithError(err))
		return
	}

	if err := validation.ValidateLoanPaymentRequest(&req); err != nil {
		writeErrorWithDetails(w, r, err)
		return
	}

	paymentAmount, totalPayment, totalInterest, err := calculations.CalculateLoanPayment(
		req.Principal,
		req.AnnualRate,
		req.Years,
		req.PaymentsPerYear,
	)
	if err != nil {
		writeErrorWithDetails(w, r, apierrors.InternalError("loan payment calculation failed").WithError(err))
		return
	}

	response := models.LoanPaymentResponse{
		PaymentAmount: paymentAmount,
		TotalPayment:  totalPayment,
		TotalInterest: totalInterest,
	}

	if err := writeSuccessResponse(w, r, response); err != nil {
		// Error already logged, headers likely already sent
		return
	}
}
