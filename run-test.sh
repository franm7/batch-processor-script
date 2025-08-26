#!/bin/bash

# run-batch-delay-test.sh

set -e

echo "Starting batch processor delay test..."

# Start OTEL collector in background
echo "Starting OTEL collector..."
otelcol --config=collector-config.yaml &
COLLECTOR_PID=$!

# Wait for collector to start
sleep 2

# Start data generator in background
echo "Starting data generator..."
go run generator.go &
GENERATOR_PID=$!

# Run test for specified duration (default 10 minutes)
DURATION=${1:-600}
echo "Running test for ${DURATION} seconds..."

# Wait for test duration
sleep $DURATION

# Stop processes
echo "Stopping test..."
kill $GENERATOR_PID $COLLECTOR_PID

echo "Test complete. Check logs for batch timing data."
