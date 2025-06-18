package db

import (
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func (o optionFunc) apply(c *tracer) {
	o(c)
}

// WithAttributes specifies additional attributes to be added to the span.
func WithAttributes(attrs ...attribute.KeyValue) Option {
	return optionFunc(func(cfg *tracer) {
		cfg.attrs = append(cfg.attrs, attrs...)
	})
}

func (f metricsOptionFunc) applyMetricsOptions(o *metricsOptions) {
	f(o)
}

// WithMeterProvider sets meter provider.
func WithMeterProvider(p metric.MeterProvider) MetricsOption {
	return struct {
		metricsOptionFunc
	}{
		metricsOptionFunc: func(o *metricsOptions) {
			o.meterProvider = p
		},
	}
}

// WithMinimumReadDBMetricsInterval sets the minimum interval between calls to db.Stats(). Negative values are ignored.
func WithMinimumReadDBMetricsInterval(interval time.Duration) MetricsOption {
	return metricsOptionFunc(func(o *metricsOptions) {
		o.minimumReadDBStatsInterval = interval
	})
}
