package timestamplogger

import (
    "context"
    "time"
    "go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/consumer"
    "go.opentelemetry.io/collector/component"
    "go.uber.org/zap"
)

type timestampLoggerProcessor struct {
    logger *zap.Logger
    startTime time.Time
    batchCount int
}

func newTimestampLoggerProcessor(logger *zap.Logger) *timestampLoggerProcessor {
    return &timestampLoggerProcessor{
        logger: logger,
        startTime: time.Now(),
        batchCount: 0,
    }
}

func (p *timestampLoggerProcessor) Capabilities() consumer.Capabilities {
    return consumer.Capabilities{MutatesData: false}
}

func (p *timestampLoggerProcessor) Shutdown(context.Context) error {
    return nil
}
func (p *timestampLoggerProcessor) Start(context.Context, component.Host) error {
    return nil
}



func (p *timestampLoggerProcessor) ConsumeMetrics(_ context.Context, md pmetric.Metrics) (error) {
    now := time.Now()
    p.batchCount++
    
    // Calculate expected time (batch_count * 60s from start)
    expectedTime := p.startTime.Add(time.Duration(p.batchCount) * 60 * time.Second)
    drift := now.Sub(expectedTime)
    
    p.logger.Info("Batch received",
        zap.Time("timestamp", now),
        zap.Int("batch_number", p.batchCount),
        zap.Duration("drift_ms", drift),
        zap.Int("metric_count", md.MetricCount()),
    )
    
    return nil
}
