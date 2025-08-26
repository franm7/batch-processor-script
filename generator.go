package main

import (
    "context"
    "log"
    "time"
    "github.com/brianvoe/gofakeit/v6"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
    sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func main() {
    ctx := context.Background()
    
    // Setup OTLP exporter
    exporter, err := otlpmetricgrpc.New(ctx,
        otlpmetricgrpc.WithEndpoint("localhost:4317"),
        otlpmetricgrpc.WithInsecure(),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Setup meter provider
    provider := sdkmetric.NewMeterProvider(
    sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter, 
        sdkmetric.WithInterval(1*time.Millisecond))),
	)
    otel.SetMeterProvider(provider)
    
    // Create meter and instruments
    meter := otel.Meter("batch-delay-test")
    counter, _ := meter.Int64Counter("test_counter")
    gauge, _ := meter.Float64Gauge("test_gauge")
    
    log.Println("Starting data generation...")
    
    ticker := time.NewTicker(1 * time.Millisecond)
    defer ticker.Stop()

    for range ticker.C {
        for i := 0; i < 10000000; i++ {
            counter.Add(ctx, int64(gofakeit.Number(1, 100)))
            gauge.Record(ctx, gofakeit.Float64())
            histogram, _ := meter.Float64Histogram("test_histogram")
            histogram.Record(ctx, gofakeit.Float64Range(0, 1000))
        }
    }
}
