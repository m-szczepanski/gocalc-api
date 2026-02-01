package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/m-szczepanski/gocalc-api/internal/config"
	"github.com/m-szczepanski/gocalc-api/internal/handlers"
	"github.com/m-szczepanski/gocalc-api/internal/middleware"
)

func main() {
	// Load configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/ready", handlers.ReadinessHandler)

	mux.HandleFunc("/api/math/add", handlers.AddHandler)
	mux.HandleFunc("/api/math/subtract", handlers.SubtractHandler)
	mux.HandleFunc("/api/math/multiply", handlers.MultiplyHandler)
	mux.HandleFunc("/api/math/divide", handlers.DivideHandler)

	mux.HandleFunc("/api/finance/vat", handlers.VATHandler)
	mux.HandleFunc("/api/finance/compound-interest", handlers.CompoundInterestHandler)
	mux.HandleFunc("/api/finance/loan-payment", handlers.LoanPaymentHandler)

	mux.HandleFunc("/api/utils/bmi", handlers.BMIHandler)
	mux.HandleFunc("/api/utils/unit-conversion", handlers.UnitConversionHandler)

	// Configure rate limiter with requests per second (RPM / 60)
	rps := cfg.RateLimit.RequestsPerMinute / 60.0
	rateLimiter := middleware.NewRateLimiter(rps, cfg.RateLimit.Burst)

	var handler http.Handler = mux
	handler = middleware.RequestIDMiddleware(handler)
	handler = rateLimiter.Middleware(handler)
	handler = middleware.TimeoutMiddleware(cfg.Server.RequestTimeout)(handler)
	handler = middleware.LoggingMiddleware(handler)
	handler = middleware.NewErrorHandler(handler)

	// Configure server with timeouts from config
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		log.Printf("API running on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down server...")

	// Create a context with timeout from config
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
