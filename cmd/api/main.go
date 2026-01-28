package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/m-szczepanski/gocalc-api/internal/handlers"
	"github.com/m-szczepanski/gocalc-api/internal/middleware"
)

const (
	port            = ":8080"
	shutdownTimeout = 15 * time.Second
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.HealthHandler)

	mux.HandleFunc("/api/math/add", handlers.AddHandler)
	mux.HandleFunc("/api/math/subtract", handlers.SubtractHandler)
	mux.HandleFunc("/api/math/multiply", handlers.MultiplyHandler)
	mux.HandleFunc("/api/math/divide", handlers.DivideHandler)

	mux.HandleFunc("/api/finance/vat", handlers.VATHandler)
	mux.HandleFunc("/api/finance/compound-interest", handlers.CompoundInterestHandler)
	mux.HandleFunc("/api/finance/loan-payment", handlers.LoanPaymentHandler)

	mux.HandleFunc("/api/utils/bmi", handlers.BMIHandler)
	mux.HandleFunc("/api/utils/unit-conversion", handlers.UnitConversionHandler)

	var handler http.Handler = mux
	handler = middleware.RequestIDMiddleware(handler)
	handler = middleware.LoggingMiddleware(handler)
	handler = middleware.NewErrorHandler(handler)

	server := &http.Server{
		Addr:         port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("API running on %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down server...")

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
