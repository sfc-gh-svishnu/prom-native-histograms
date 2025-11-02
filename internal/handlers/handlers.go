package handlers

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/svishnu/prom-native-histograms/internal/config"
	"github.com/svishnu/prom-native-histograms/internal/metrics"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	metrics *metrics.Metrics
	config  *config.Config
}

// NewHandler creates a new handler with dependencies
func NewHandler(m *metrics.Metrics, cfg *config.Config) *Handler {
	return &Handler{
		metrics: m,
		config:  cfg,
	}
}

// APIHandler simulates some work and records metrics
func (h *Handler) APIHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Increment request counter
	h.metrics.TotalRequests.WithLabelValues("/api").Inc()

	// Simulate some work with random duration
	minDuration := h.config.MinRequestDuration.Milliseconds()
	maxDuration := h.config.MaxRequestDuration.Milliseconds()
	sleepDuration := time.Duration(minDuration+rand.Int63n(maxDuration-minDuration)) * time.Millisecond
	time.Sleep(sleepDuration)

	// Record request duration
	duration := time.Since(start).Seconds()
	h.metrics.RequestDuration.Observe(duration)

	// Simulate response size
	size := float64(h.config.MinResponseSize + rand.Intn(h.config.MaxResponseSize-h.config.MinResponseSize))
	h.metrics.ResponseSize.Observe(size)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// HealthHandler handles health check requests
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("healthy"))
}

// RootHandler handles root path requests
func (h *Handler) RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`Prometheus Native Histogram Demo

Endpoints:
- /api - Sample API endpoint that generates metrics
- /metrics - Prometheus metrics endpoint
- /health - Health check endpoint

Native histograms are enabled for:
- http_request_duration_seconds
- http_response_size_bytes
`))
}

