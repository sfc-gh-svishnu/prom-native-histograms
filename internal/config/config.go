package config

import (
	"time"
)

// Config holds application configuration
type Config struct {
	// Server configuration
	ServerAddress string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration

	// Metrics generation configuration
	MetricsGenerationInterval time.Duration
	MinRequestDuration        time.Duration
	MaxRequestDuration        time.Duration
	MinResponseSize           int
	MaxResponseSize           int
}

// NewConfig creates a new configuration with default values
func NewConfig() *Config {
	return &Config{
		ServerAddress:             ":8080",
		ReadTimeout:               15 * time.Second,
		WriteTimeout:              15 * time.Second,
		MetricsGenerationInterval: 1 * time.Second,
		MinRequestDuration:        10 * time.Millisecond,
		MaxRequestDuration:        500 * time.Millisecond,
		MinResponseSize:           1024,   // 1KB
		MaxResponseSize:           102400, // 100KB
	}
}
