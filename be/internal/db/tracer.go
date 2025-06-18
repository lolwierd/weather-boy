package db

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var RowsAffectedKey = attribute.Key("pgx.rows_affected")

func NewpgxTracer(opts ...Option) *tracer {
	cfg := &tracer{
		attrs: []attribute.KeyValue{
			semconv.DBSystemPostgreSQL,
		},
	}

	for _, opt := range opts {
		opt.apply(cfg)
	}

	return &tracer{
		tracer: otel.GetTracerProvider().Tracer(otelName),
		attrs:  cfg.attrs,
	}
}
