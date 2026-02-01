package config

import (
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		want    *Config
		wantErr bool
	}{
		{
			name:    "default configuration",
			envVars: map[string]string{},
			want: &Config{
				Server: ServerConfig{
					Port:            "8080",
					ReadTimeout:     10 * time.Second,
					WriteTimeout:    10 * time.Second,
					IdleTimeout:     120 * time.Second,
					ShutdownTimeout: 15 * time.Second,
					RequestTimeout:  30 * time.Second,
				},
				RateLimit: RateLimitConfig{
					RequestsPerMinute: 100.0,
					Burst:             20,
				},
			},
			wantErr: false,
		},
		{
			name: "custom port",
			envVars: map[string]string{
				"PORT": "3000",
			},
			want: &Config{
				Server: ServerConfig{
					Port:            "3000",
					ReadTimeout:     10 * time.Second,
					WriteTimeout:    10 * time.Second,
					IdleTimeout:     120 * time.Second,
					ShutdownTimeout: 15 * time.Second,
					RequestTimeout:  30 * time.Second,
				},
				RateLimit: RateLimitConfig{
					RequestsPerMinute: 100.0,
					Burst:             20,
				},
			},
			wantErr: false,
		},
		{
			name: "custom timeouts",
			envVars: map[string]string{
				"READ_TIMEOUT":     "5s",
				"WRITE_TIMEOUT":    "15s",
				"IDLE_TIMEOUT":     "60s",
				"SHUTDOWN_TIMEOUT": "30s",
				"REQUEST_TIMEOUT":  "45s",
			},
			want: &Config{
				Server: ServerConfig{
					Port:            "8080",
					ReadTimeout:     5 * time.Second,
					WriteTimeout:    15 * time.Second,
					IdleTimeout:     60 * time.Second,
					ShutdownTimeout: 30 * time.Second,
					RequestTimeout:  45 * time.Second,
				},
				RateLimit: RateLimitConfig{
					RequestsPerMinute: 100.0,
					Burst:             20,
				},
			},
			wantErr: false,
		},
		{
			name: "custom rate limit",
			envVars: map[string]string{
				"RATE_LIMIT_RPM":   "200",
				"RATE_LIMIT_BURST": "50",
			},
			want: &Config{
				Server: ServerConfig{
					Port:            "8080",
					ReadTimeout:     10 * time.Second,
					WriteTimeout:    10 * time.Second,
					IdleTimeout:     120 * time.Second,
					ShutdownTimeout: 15 * time.Second,
					RequestTimeout:  30 * time.Second,
				},
				RateLimit: RateLimitConfig{
					RequestsPerMinute: 200.0,
					Burst:             50,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid timeout format falls back to default",
			envVars: map[string]string{
				"READ_TIMEOUT": "invalid",
			},
			want: &Config{
				Server: ServerConfig{
					Port:            "8080",
					ReadTimeout:     10 * time.Second,
					WriteTimeout:    10 * time.Second,
					IdleTimeout:     120 * time.Second,
					ShutdownTimeout: 15 * time.Second,
					RequestTimeout:  30 * time.Second,
				},
				RateLimit: RateLimitConfig{
					RequestsPerMinute: 100.0,
					Burst:             20,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid rate limit falls back to default",
			envVars: map[string]string{
				"RATE_LIMIT_RPM":   "not-a-number",
				"RATE_LIMIT_BURST": "also-invalid",
			},
			want: &Config{
				Server: ServerConfig{
					Port:            "8080",
					ReadTimeout:     10 * time.Second,
					WriteTimeout:    10 * time.Second,
					IdleTimeout:     120 * time.Second,
					ShutdownTimeout: 15 * time.Second,
					RequestTimeout:  30 * time.Second,
				},
				RateLimit: RateLimitConfig{
					RequestsPerMinute: 100.0,
					Burst:             20,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set test environment variables (auto-restored after test)
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			got, err := Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if got.Server.Port != tt.want.Server.Port {
					t.Errorf("Port = %v, want %v", got.Server.Port, tt.want.Server.Port)
				}
				if got.Server.ReadTimeout != tt.want.Server.ReadTimeout {
					t.Errorf("ReadTimeout = %v, want %v", got.Server.ReadTimeout, tt.want.Server.ReadTimeout)
				}
				if got.Server.WriteTimeout != tt.want.Server.WriteTimeout {
					t.Errorf("WriteTimeout = %v, want %v", got.Server.WriteTimeout, tt.want.Server.WriteTimeout)
				}
				if got.Server.IdleTimeout != tt.want.Server.IdleTimeout {
					t.Errorf("IdleTimeout = %v, want %v", got.Server.IdleTimeout, tt.want.Server.IdleTimeout)
				}
				if got.Server.ShutdownTimeout != tt.want.Server.ShutdownTimeout {
					t.Errorf("ShutdownTimeout = %v, want %v", got.Server.ShutdownTimeout, tt.want.Server.ShutdownTimeout)
				}
				if got.Server.RequestTimeout != tt.want.Server.RequestTimeout {
					t.Errorf("RequestTimeout = %v, want %v", got.Server.RequestTimeout, tt.want.Server.RequestTimeout)
				}
				if got.RateLimit.RequestsPerMinute != tt.want.RateLimit.RequestsPerMinute {
					t.Errorf("RequestsPerMinute = %v, want %v", got.RateLimit.RequestsPerMinute, tt.want.RateLimit.RequestsPerMinute)
				}
				if got.RateLimit.Burst != tt.want.RateLimit.Burst {
					t.Errorf("Burst = %v, want %v", got.RateLimit.Burst, tt.want.RateLimit.Burst)
				}
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid configuration",
			config: &Config{
				Server: ServerConfig{
					Port:            "8080",
					ReadTimeout:     10 * time.Second,
					WriteTimeout:    10 * time.Second,
					IdleTimeout:     120 * time.Second,
					ShutdownTimeout: 15 * time.Second,
					RequestTimeout:  30 * time.Second,
				},
				RateLimit: RateLimitConfig{
					RequestsPerMinute: 100.0,
					Burst:             20,
				},
			},
			wantErr: false,
		},
		{
			name: "empty port",
			config: &Config{
				Server: ServerConfig{
					Port:            "",
					ReadTimeout:     10 * time.Second,
					WriteTimeout:    10 * time.Second,
					IdleTimeout:     120 * time.Second,
					ShutdownTimeout: 15 * time.Second,
					RequestTimeout:  30 * time.Second,
				},
				RateLimit: RateLimitConfig{
					RequestsPerMinute: 100.0,
					Burst:             20,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid port - non-numeric",
			config: &Config{
				Server: ServerConfig{
					Port:            "abc",
					ReadTimeout:     10 * time.Second,
					WriteTimeout:    10 * time.Second,
					IdleTimeout:     120 * time.Second,
					ShutdownTimeout: 15 * time.Second,
					RequestTimeout:  30 * time.Second,
				},
				RateLimit: RateLimitConfig{
					RequestsPerMinute: 100.0,
					Burst:             20,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid port - too low",
			config: &Config{
				Server: ServerConfig{
					Port:            "0",
					ReadTimeout:     10 * time.Second,
					WriteTimeout:    10 * time.Second,
					IdleTimeout:     120 * time.Second,
					ShutdownTimeout: 15 * time.Second,
					RequestTimeout:  30 * time.Second,
				},
				RateLimit: RateLimitConfig{
					RequestsPerMinute: 100.0,
					Burst:             20,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid port - too high",
			config: &Config{
				Server: ServerConfig{
					Port:            "65536",
					ReadTimeout:     10 * time.Second,
					WriteTimeout:    10 * time.Second,
					IdleTimeout:     120 * time.Second,
					ShutdownTimeout: 15 * time.Second,
					RequestTimeout:  30 * time.Second,
				},
				RateLimit: RateLimitConfig{
					RequestsPerMinute: 100.0,
					Burst:             20,
				},
			},
			wantErr: true,
		},
		{
			name: "negative timeout",
			config: &Config{
				Server: ServerConfig{
					Port:            "8080",
					ReadTimeout:     -1 * time.Second,
					WriteTimeout:    10 * time.Second,
					IdleTimeout:     120 * time.Second,
					ShutdownTimeout: 15 * time.Second,
					RequestTimeout:  30 * time.Second,
				},
				RateLimit: RateLimitConfig{
					RequestsPerMinute: 100.0,
					Burst:             20,
				},
			},
			wantErr: true,
		},
		{
			name: "zero rate limit",
			config: &Config{
				Server: ServerConfig{
					Port:            "8080",
					ReadTimeout:     10 * time.Second,
					WriteTimeout:    10 * time.Second,
					IdleTimeout:     120 * time.Second,
					ShutdownTimeout: 15 * time.Second,
					RequestTimeout:  30 * time.Second,
				},
				RateLimit: RateLimitConfig{
					RequestsPerMinute: 0,
					Burst:             20,
				},
			},
			wantErr: true,
		},
		{
			name: "negative burst",
			config: &Config{
				Server: ServerConfig{
					Port:            "8080",
					ReadTimeout:     10 * time.Second,
					WriteTimeout:    10 * time.Second,
					IdleTimeout:     120 * time.Second,
					ShutdownTimeout: 15 * time.Second,
					RequestTimeout:  30 * time.Second,
				},
				RateLimit: RateLimitConfig{
					RequestsPerMinute: 100.0,
					Burst:             -1,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
