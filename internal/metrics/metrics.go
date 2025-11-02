package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Metrics holds all Prometheus metrics for the application
type Metrics struct {
	RequestDuration prometheus.Histogram
	ResponseSize    prometheus.Histogram
	TotalRequests   *prometheus.CounterVec
}

// NewMetrics creates and registers all Prometheus metrics
func NewMetrics() *Metrics {
	m := &Metrics{
		RequestDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:                            "http_request_duration_seconds",
			Help:                            "Duration of HTTP requests in seconds (native histogram)",
			NativeHistogramBucketFactor:     1.1,       // Native histogram bucket growth factor
			NativeHistogramMaxBucketNumber:  100,       // Maximum number of buckets
			NativeHistogramMinResetDuration: time.Hour, // Minimum time between histogram resets
		}),
		ResponseSize: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:                           "http_response_size_bytes",
			Help:                           "Size of HTTP responses in bytes (native histogram)",
			NativeHistogramBucketFactor:    1.1,
			NativeHistogramMaxBucketNumber: 100,
		}),
		TotalRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"endpoint"},
		),
	}

	// Register metrics with Prometheus
	prometheus.MustRegister(m.RequestDuration)
	prometheus.MustRegister(m.ResponseSize)
	prometheus.MustRegister(m.TotalRequests)

	return m
}

