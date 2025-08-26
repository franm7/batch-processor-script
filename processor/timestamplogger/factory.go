// processor/timestamplogger/factory.go
package timestamplogger

import (
    "context"
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/consumer"
    "go.opentelemetry.io/collector/processor"
)

func NewFactory() processor.Factory {
    return processor.NewFactory(
        "timestamplogger",
        func() component.Config { return &struct{}{} },
        processor.WithMetrics(createMetricsProcessor, component.StabilityLevelDevelopment),
    )
}

func createMetricsProcessor(_ context.Context, set processor.CreateSettings, _ component.Config, nextConsumer consumer.Metrics) (processor.Metrics, error) {
    return newTimestampLoggerProcessor(set.Logger), nil
}
