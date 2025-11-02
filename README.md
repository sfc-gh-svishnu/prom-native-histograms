# Prometheus Native Histograms Demo

This project demonstrates how to use Prometheus native histograms in a Go application with a clean, modular architecture. The application exposes metrics that are scraped by a Prometheus server running in Docker.

## What are Native Histograms?

Native histograms are a new feature in Prometheus that provides more efficient and accurate histogram representation compared to classic histograms. They use a dynamic bucket system that adapts to the data distribution, reducing storage requirements and improving query performance.

## Project Structure

```
.
├── cmd/                             # Application entry points
│   └── prom-native-histograms/      # Main application
│       └── main.go                  # Application entry point
├── internal/                        # Private application packages
│   ├── config/                      # Configuration management
│   │   └── config.go
│   ├── metrics/                     # Prometheus metrics definitions
│   │   └── metrics.go
│   ├── handlers/                    # HTTP request handlers
│   │   └── handlers.go
│   └── worker/                      # Background workers
│       └── worker.go
├── go.mod                           # Go module dependencies
├── go.sum                           # Dependency checksums
├── prometheus.yml                   # Prometheus server configuration
├── docker-compose.yml               # Docker Compose setup
└── README.md                        # This file
```

### Package Overview

- **`cmd/prom-native-histograms/main.go`**: Entry point that wires up all components and starts the HTTP server with graceful shutdown
- **`internal/config`**: Centralized configuration with sensible defaults
- **`internal/metrics`**: Prometheus metrics definitions and registration
- **`internal/handlers`**: HTTP handlers for API, health check, and info endpoints
- **`internal/worker`**: Background workers for continuous metrics generation

This structure follows the [golang-standards/project-layout](https://github.com/golang-standards/project-layout) convention, making it easy to add multiple binaries if needed in the future.

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose

## Setup and Running

### 1. Install Go Dependencies

```bash
go mod download
```

### 2. Run the Go Application

```bash
go run cmd/prom-native-histograms/main.go
```

The application will start on port 8080 and expose the following endpoints:

- `http://localhost:8080/` - Info page
- `http://localhost:8080/api` - Sample API endpoint that generates metrics
- `http://localhost:8080/metrics` - Prometheus metrics endpoint
- `http://localhost:8080/health` - Health check endpoint

To stop the server gracefully, press `Ctrl+C`.

### 3. Start Prometheus Server

In a separate terminal, start the Prometheus server using Docker Compose:

```bash
docker-compose up -d
```

This will:
- Start Prometheus on port 9090
- Enable native histogram support with the `--enable-feature=native-histograms` flag
- Configure Prometheus to scrape metrics from your Go application every 5 seconds

### 4. Access Prometheus UI

Open your browser and navigate to:

```
http://localhost:9090
```

### 5. Query Native Histograms

In the Prometheus UI, try these queries:

#### View the histogram:
```promql
http_request_duration_seconds
```

#### Calculate quantiles (e.g., 95th percentile):
```promql
histogram_quantile(0.95, http_request_duration_seconds)
```

#### View response size histogram:
```promql
http_response_size_bytes
```

#### Calculate average response time:
```promql
histogram_avg(http_request_duration_seconds)
```

#### View request rate:
```promql
rate(http_requests_total[1m])
```

## Architecture Details

### Modular Design

The application follows Go best practices with a clean, modular architecture:

1. **Separation of Concerns**: Each package has a single, well-defined responsibility
2. **Dependency Injection**: Dependencies are explicitly passed to constructors
3. **Graceful Shutdown**: The server handles SIGINT/SIGTERM signals gracefully
4. **Context-based Cancellation**: Background workers respect context cancellation
5. **Configuration Management**: Centralized configuration with default values

### Native Histogram Metrics

The application exposes two native histograms:

1. **`http_request_duration_seconds`** - Tracks the duration of HTTP requests
   - Uses native histogram bucket factor of 1.1
   - Maximum of 100 buckets
   - Automatically generates data every second

2. **`http_response_size_bytes`** - Tracks the size of HTTP responses
   - Uses native histogram bucket factor of 1.1
   - Maximum of 100 buckets
   - Simulates varying response sizes

3. **`http_requests_total`** - Counter tracking total requests by endpoint

### Configuration

Default configuration values are set in `internal/config/config.go`:

```go
ServerAddress:             ":8080"
ReadTimeout:               15 * time.Second
WriteTimeout:              15 * time.Second
MetricsGenerationInterval: 1 * time.Second
MinRequestDuration:        10 * time.Millisecond
MaxRequestDuration:        500 * time.Millisecond
MinResponseSize:           1024      // 1KB
MaxResponseSize:           102400    // 100KB
```

You can easily modify these values or extend the configuration to read from environment variables or config files.

## Testing the Setup

### Generate Traffic

You can generate traffic to see the metrics in action:

```bash
# Generate some API requests
for i in {1..100}; do curl http://localhost:8080/api; done
```

### Verify Metrics

Check the raw metrics:

```bash
curl http://localhost:8080/metrics | grep http_request_duration_seconds
```

You should see native histogram buckets in the output.

### Check Prometheus Targets

Verify Prometheus is successfully scraping:

```bash
curl -s http://localhost:9090/api/v1/targets | python3 -m json.tool
```

## Building for Production

### Build a Binary

```bash
go build -o prom-native-histograms ./cmd/prom-native-histograms
./prom-native-histograms
```

Or install it globally:

```bash
go install ./cmd/prom-native-histograms
prom-native-histograms
```

### Build with Docker

Create a `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o prom-native-histograms ./cmd/prom-native-histograms

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/prom-native-histograms .
EXPOSE 8080
CMD ["./prom-native-histograms"]
```

Build and run:

```bash
docker build -t prom-native-histograms .
docker run -p 8080:8080 prom-native-histograms
```

## Stopping the Services

### Stop the Go application:
Press `Ctrl+C` in the terminal running the Go app (graceful shutdown is automatic).

### Stop Prometheus:
```bash
docker-compose down
```

To also remove the stored Prometheus data:
```bash
docker-compose down -v
```

## Development

### Running Tests

```bash
go test ./...
```

### Running Linters

```bash
golangci-lint run
```

### Code Organization Best Practices

1. **Use internal/ for private packages**: Prevents external imports
2. **Keep packages focused**: Each package has one responsibility
3. **Minimize dependencies**: Only import what you need
4. **Use interfaces**: Makes testing and mocking easier
5. **Document exports**: Add comments for all exported types and functions

## Key Configuration Details

### Native Histogram Configuration

The native histogram is configured in `internal/metrics/metrics.go`:

```go
prometheus.NewHistogram(prometheus.HistogramOpts{
    Name:                        "http_request_duration_seconds",
    Help:                        "Duration of HTTP requests in seconds",
    NativeHistogramBucketFactor: 1.1,  // Controls bucket resolution
    NativeHistogramMaxBucketNumber: 100, // Maximum number of buckets
})
```

### Prometheus Configuration

In `prometheus.yml`, native histogram support is enabled by:

1. Using the `--enable-feature=native-histograms` flag in Docker Compose
2. Configuring scrape protocols to include newer formats that support native histograms

## Troubleshooting

### Port Already in Use

If port 8080 or 9090 is already in use, you can modify:
- For the Go app: Change `ServerAddress` in `internal/config/config.go`
- For Prometheus: Change the port mapping in `docker-compose.yml`

### Prometheus Can't Reach the App

If you're on Linux, change `host.docker.internal` to `172.17.0.1` (Docker bridge IP) in `prometheus.yml`.

### No Native Histogram Data

Ensure:
1. The Go app is running and accessible at `http://localhost:8080/metrics`
2. Prometheus is started with `--enable-feature=native-histograms`
3. The scrape is successful (check Targets page in Prometheus UI at http://localhost:9090/targets)

### Import Path Issues

Make sure your `go.mod` file has the correct module name. If you clone this repo, you may need to update import paths.

## Learn More

- [Prometheus Native Histograms Documentation](https://prometheus.io/docs/concepts/metric_types/#histogram)
- [Prometheus Client Golang](https://github.com/prometheus/client_golang)
- [Native Histograms in Prometheus](https://prometheus.io/docs/prometheus/latest/feature_flags/#native-histograms)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

## License

MIT
