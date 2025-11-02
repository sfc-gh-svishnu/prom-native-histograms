# Prometheus Native Histograms Demo

This project demonstrates how to use Prometheus native histograms in a Go application and scrape them using Prometheus server.

## What are Native Histograms?

Native histograms are a new feature in Prometheus that provides more efficient and accurate histogram representation compared to classic histograms. They use a dynamic bucket system that adapts to the data distribution, reducing storage requirements and improving query performance.

## Project Structure

```
.
├── main.go              # Go application with native histogram metrics
├── go.mod               # Go module dependencies
├── prometheus.yml       # Prometheus server configuration
├── docker-compose.yml   # Docker Compose setup for Prometheus
└── README.md           # This file
```

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
go run main.go
```

The application will start on port 8080 and expose the following endpoints:

- `http://localhost:8080/` - Info page
- `http://localhost:8080/api` - Sample API endpoint that generates metrics
- `http://localhost:8080/metrics` - Prometheus metrics endpoint
- `http://localhost:8080/health` - Health check endpoint

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

## Native Histogram Metrics

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

## Stopping the Services

### Stop the Go application:
Press `Ctrl+C` in the terminal running the Go app.

### Stop Prometheus:
```bash
docker-compose down
```

To also remove the stored Prometheus data:
```bash
docker-compose down -v
```

## Key Configuration Details

### Native Histogram Configuration in Go

The native histogram is configured with:

```go
requestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
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
- For the Go app: Change the port in `main.go`
- For Prometheus: Change the port mapping in `docker-compose.yml`

### Prometheus Can't Reach the App

If you're on Linux, change `host.docker.internal` to `172.17.0.1` (Docker bridge IP) in `prometheus.yml`.

### No Native Histogram Data

Ensure:
1. The Go app is running and accessible at `http://localhost:8080/metrics`
2. Prometheus is started with `--enable-feature=native-histograms`
3. The scrape is successful (check Targets page in Prometheus UI)

## Learn More

- [Prometheus Native Histograms Documentation](https://prometheus.io/docs/concepts/metric_types/#histogram)
- [Prometheus Client Golang](https://github.com/prometheus/client_golang)
- [Native Histograms in Prometheus](https://prometheus.io/docs/prometheus/latest/feature_flags/#native-histograms)

## License

MIT

