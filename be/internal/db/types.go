package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	otelName            = "github.com/lolwierd/weatherboy/be/internal/db"
	sqlOperationUnknown = "UNKNOWN"
	// Unit constants for deprecated metric units
	unitDimensionless = "1"
	unitBytes         = "By"
	unitMilliseconds  = "ms"

	// defaultMinimumReadDBStatsInterval is the default minimum interval between calls to db.Stats().
	defaultMinimumReadDBStatsInterval = time.Second

	pgxPoolAcquireCount            = "pgxpool_acquires"
	pgxpoolAcquireDuration         = "pgxpool_acquire_duration"
	pgxpoolAcquiredConns           = "pgxpool_acquired_conns"
	pgxpoolCancelledAcquires       = "pgxpool_canceled_acquires"
	pgxpoolConstructingConns       = "pgxpool_constructing_conns"
	pgxpoolEmptyAcquire            = "pgxpool_empty_acquire"
	pgxpoolIdleConns               = "pgxpool_idle_conns"
	pgxpoolMaxConns                = "pgxpool_max_conns"
	pgxpoolMaxIdleDestroyCount     = "pgxpool_max_idle_destroys"
	pgxpoolMaxLifetimeDestroyCount = "pgxpool_max_lifetime_destroys"
	pgxpoolNewConnsCount           = "pgxpool_new_conns"
	pgxpoolTotalConns              = "pgxpool_total_conns"
)

// ConnPool abstracts the pgxpool.Pool methods used in the repository.
type ConnPool interface {
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
	Ping(context.Context) error
	Close()
}

type Driver struct {
	ConnPool ConnPool
}

type tracer struct {
	tracer trace.Tracer
	attrs  []attribute.KeyValue
}

// Option specifies instrumentation configuration options.
type Option interface {
	apply(*tracer)
}

type optionFunc func(*tracer)

// MetricsOption allows for managing otelsql configuration using functional options.
type MetricsOption interface {
	applyMetricsOptions(o *metricsOptions)
}

type metricsOptions struct {
	// meterProvider sets the metric.MeterProvider. If nil, the global Provider will be used.
	meterProvider metric.MeterProvider

	// minimumReadDBStatsInterval sets the minimum interval between calls to db.Stats(). Negative values are ignored.
	minimumReadDBStatsInterval time.Duration

	// defaultAttributes will be set to each metrics as default.
	defaultAttributes []attribute.KeyValue
}

type metricsOptionFunc func(o *metricsOptions)
