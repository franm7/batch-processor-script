# OpenTelemetry Batch Processor Delay Test

This project demonstrates and reproduces timing drift accumulation in the OpenTelemetry batch processor's timeout-based flushing mechanism.

## Problem Description

The OTEL batch processor's timeout mechanism accumulates delays over time, causing batches to drift from their intended schedule. Each batch arrives slightly later than expected, and these delays compound, eventually causing significant timing misalignment.

## What This Test Does

- Generates high-volume metrics data to stress the batch processor
- Uses a 60-second batch timeout to trigger time-based flushing
- Measures actual batch arrival times vs expected intervals
- Demonstrates progressive drift accumulation over time

## Project Structure

```
batch-delay-test/
├── README.md
├── Makefile
├── go.mod
├── collector-config.yaml    # OTEL collector configuration
├── generator.go            # Metric data generator
└── run-test.sh            # Test orchestration script
```

## How to Run

### Prerequisites
- Go 1.21+
- Docker

### Setup
```bash
# Clone and setup
git clone <repo-url>
cd batch-delay-test

# Install dependencies
go mod tidy

# Build generator
make build
```

### Run Test
```bash
# Terminal 1: Start OTEL collector
docker run -p 4317:4317 -p 4318:4318 \
  -v $(pwd)/collector-config.yaml:/etc/otelcol-contrib/config.yaml \
  otel/opentelemetry-collector-contrib:latest

# Terminal 2: Start data generator
./generator

# Let run for a few minutes to observe drift accumulation
```

## Reading Results

Look for log entries in the Docker terminal like:
```
2025-08-25T17:36:54.292Z  info  Metrics {"data points": 154607}
2025-08-25T17:37:54.328Z  info  Metrics {"data points": 170493}
2025-08-25T17:38:54.389Z  info  Metrics {"data points": 165821}
```

### Analyzing Drift

Calculate the interval between consecutive batch logs:
- **Expected**: Exactly 60.000 seconds between batches
- **Actual**: Measure time differences between timestamps
- **Drift**: Difference from expected 60s interval

Example drift progression:
```
Batch 1: 17:36:54.292Z (baseline)
Batch 2: 17:37:54.328Z (60.036s = 36ms drift)
Batch 3: 17:38:54.389Z (60.061s = 97ms total drift)
Batch 4: 17:39:54.467Z (60.078s = 175ms total drift)
```

## Configuration

### Batch Processor Settings
```yaml
processors:
  batch:
    timeout: 60s                    # Target flush interval
    send_batch_size: 1000000000000  # Very large to force timeout-based flushing
```

### Data Generator Settings
```go
ticker := time.NewTicker(1 * time.Millisecond)  // Generate every 1ms
for i := 0; i < 10000000; i++ {                 // 10M metrics per tick
    // Generate counter, gauge, histogram metrics
}
```


## Cleanup

```bash
# Stop all processes
make clean

# Or manually
pkill -f "generator\|otelcol"
docker stop $(docker ps -q --filter ancestor=otel/opentelemetry-collector-contrib)
```