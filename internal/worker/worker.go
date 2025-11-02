package worker

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/svishnu/prom-native-histograms/internal/config"
	"github.com/svishnu/prom-native-histograms/internal/metrics"
)

// MetricsGenerator generates metrics in the background
type MetricsGenerator struct {
	metrics *metrics.Metrics
	config  *config.Config
}

// NewMetricsGenerator creates a new metrics generator
func NewMetricsGenerator(m *metrics.Metrics, cfg *config.Config) *MetricsGenerator {
	return &MetricsGenerator{
		metrics: m,
		config:  cfg,
	}
}

// Start begins generating metrics in the background
func (mg *MetricsGenerator) Start(ctx context.Context) {
	ticker := time.NewTicker(mg.config.MetricsGenerationInterval)
	defer ticker.Stop()

	log.Println("Background metrics generator started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Background metrics generator stopped")
			return
		case <-ticker.C:
			mg.generateMetrics()
		}
	}
}

func (mg *MetricsGenerator) generateMetrics() {
	// Simulate background operations with varying durations
	duration := float64(rand.Intn(1000)) / 1000.0 // 0 to 1 second
	mg.metrics.RequestDuration.Observe(duration)

	// Simulate varying response sizes
	size := float64(rand.Intn(50000))
	mg.metrics.ResponseSize.Observe(size)
}

