package handlers

import (
	"net/http"
	"runtime"
	"time"
)

var startTime = time.Now()

// HealthResponse contains health check information.
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Uptime    string            `json:"uptime"`
	Version   string            `json:"version"`
	Checks    map[string]string `json:"checks"`
}

// ReadinessResponse contains readiness probe information.
type ReadinessResponse struct {
	Ready  bool   `json:"ready"`
	Status string `json:"status"`
}

// HealthHandler returns detailed health information about the service.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime)

	// Perform health checks
	checks := make(map[string]string)
	checks["memory"] = checkMemory()
	checks["goroutines"] = checkGoroutines()

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Uptime:    uptime.String(),
		Version:   "1.0.0",
		Checks:    checks,
	}

	writeSuccessResponse(w, r, response)
}

// ReadinessHandler indicates if the service is ready to accept traffic.
// This is useful for Kubernetes readiness probes.
func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	// Service is always ready if it's running (no external dependencies)
	response := ReadinessResponse{
		Ready:  true,
		Status: "ready",
	}

	writeSuccessResponse(w, r, response)
}

// checkMemory returns memory usage information.
func checkMemory() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Convert to MB for readability
	allocMB := m.Alloc / 1024 / 1024
	sysMB := m.Sys / 1024 / 1024

	return formatMemoryCheck(allocMB, sysMB)
}

// formatMemoryCheck formats memory statistics.
func formatMemoryCheck(alloc, sys uint64) string {
	if alloc > 500 || sys > 1000 {
		return "warning"
	}
	return "ok"
}

// checkGoroutines returns goroutine count status.
func checkGoroutines() string {
	count := runtime.NumGoroutine()
	if count > 1000 {
		return "warning"
	}
	return "ok"
}
