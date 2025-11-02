package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/svishnu/prom-native-histograms/internal/config"
	"github.com/svishnu/prom-native-histograms/internal/handlers"
	"github.com/svishnu/prom-native-histograms/internal/metrics"
	"github.com/svishnu/prom-native-histograms/internal/worker"
)

func main() {
	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize metrics
	m := metrics.NewMetrics()

	// Initialize handlers
	h := handlers.NewHandler(m, cfg)

	// Start background metrics generator
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	generator := worker.NewMetricsGenerator(m, cfg)
	go generator.Start(ctx)

	// Set up HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.RootHandler)
	mux.HandleFunc("/api", h.APIHandler)
	mux.HandleFunc("/health", h.HealthHandler)
	mux.Handle("/metrics", promhttp.Handler())

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", cfg.ServerAddress)
		log.Printf("Metrics available at http://localhost%s/metrics", cfg.ServerAddress)
		log.Printf("API endpoint at http://localhost%s/api", cfg.ServerAddress)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Cancel context to stop background workers
	cancel()

	// Gracefully shutdown the server with a timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
