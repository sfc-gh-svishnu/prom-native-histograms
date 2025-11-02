package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Create a native histogram with custom buckets
	// NativeHistogramBucketFactor controls the resolution of the histogram
	requestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:                            "http_request_duration_seconds",
		Help:                            "Duration of HTTP requests in seconds (native histogram)",
		NativeHistogramBucketFactor:     1.1,       // Native histogram bucket growth factor
		NativeHistogramMaxBucketNumber:  100,       // Maximum number of buckets
		NativeHistogramMinResetDuration: time.Hour, // Minimum time between histogram resets
	})

	// Another native histogram for response sizes
	responseSize = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:                           "http_response_size_bytes",
		Help:                           "Size of HTTP responses in bytes (native histogram)",
		NativeHistogramBucketFactor:    1.1,
		NativeHistogramMaxBucketNumber: 100,
	})

	// A counter for total requests
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"endpoint"},
	)
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(responseSize)
	prometheus.MustRegister(totalRequests)
}

// Handler that simulates some work and records metrics
func apiHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Increment request counter
	totalRequests.WithLabelValues("/api").Inc()

	// Simulate some work with random duration (10ms to 500ms)
	sleepDuration := time.Duration(10+rand.Intn(490)) * time.Millisecond
	time.Sleep(sleepDuration)

	// Record request duration
	duration := time.Since(start).Seconds()
	requestDuration.Observe(duration)

	// Simulate response size (1KB to 100KB)
	size := float64(1024 + rand.Intn(99*1024))
	responseSize.Observe(size)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Background worker to generate some continuous metrics
func generateMetrics() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Simulate background operations with varying durations
		duration := float64(rand.Intn(1000)) / 1000.0 // 0 to 1 second
		requestDuration.Observe(duration)

		// Simulate varying response sizes
		size := float64(rand.Intn(50000))
		responseSize.Observe(size)
	}
}

func main() {
	// Start background metrics generator
	go generateMetrics()

	// Set up HTTP handlers
	http.HandleFunc("/api", apiHandler)
	http.Handle("/metrics", promhttp.Handler())

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("healthy"))
	})

	// Root endpoint with info
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	})

	log.Println("Starting server on :8080")
	log.Println("Metrics available at http://localhost:8080/metrics")
	log.Println("API endpoint at http://localhost:8080/api")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
