package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration.
type Config struct {
	Server    ServerConfig
	RateLimit RateLimitConfig
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	RequestTimeout  time.Duration
}

// RateLimitConfig holds rate limiting configuration.
type RateLimitConfig struct {
	RequestsPerMinute float64
	Burst             int
}

// Load reads configuration from environment variables with sensible defaults.
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "8080"),
			ReadTimeout:     getDurationEnv("READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getDurationEnv("WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:     getDurationEnv("IDLE_TIMEOUT", 120*time.Second),
			ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT", 15*time.Second),
			RequestTimeout:  getDurationEnv("REQUEST_TIMEOUT", 30*time.Second),
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: getFloatEnv("RATE_LIMIT_RPM", 100.0),
			Burst:             getIntEnv("RATE_LIMIT_BURST", 20),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// validate checks that configuration values are sensible.
func (c *Config) validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("PORT cannot be empty")
	}

	if c.Server.ReadTimeout <= 0 {
		return fmt.Errorf("READ_TIMEOUT must be positive")
	}

	if c.Server.WriteTimeout <= 0 {
		return fmt.Errorf("WRITE_TIMEOUT must be positive")
	}

	if c.Server.IdleTimeout <= 0 {
		return fmt.Errorf("IDLE_TIMEOUT must be positive")
	}

	if c.Server.ShutdownTimeout <= 0 {
		return fmt.Errorf("SHUTDOWN_TIMEOUT must be positive")
	}

	if c.Server.RequestTimeout <= 0 {
		return fmt.Errorf("REQUEST_TIMEOUT must be positive")
	}

	if c.RateLimit.RequestsPerMinute <= 0 {
		return fmt.Errorf("RATE_LIMIT_RPM must be positive")
	}

	if c.RateLimit.Burst <= 0 {
		return fmt.Errorf("RATE_LIMIT_BURST must be positive")
	}

	return nil
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnv retrieves an integer environment variable or returns a default value.
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getFloatEnv retrieves a float environment variable or returns a default value.
func getFloatEnv(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

// getDurationEnv retrieves a duration environment variable or returns a default value.
// Accepts values like "10s", "2m", "1h" etc.
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
