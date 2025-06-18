package db

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lolwierd/weatherboy/be/internal/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func recordMetrics(db *pgxpool.Pool, opts ...MetricsOption) error {
	var (
		err error

		acquireCount                         metric.Int64ObservableCounter
		acquireDuration                      metric.Float64ObservableCounter
		acquiredConns                        metric.Int64ObservableUpDownCounter
		cancelledAcquires                    metric.Int64ObservableCounter
		constructingConns                    metric.Int64ObservableUpDownCounter
		emptyAcquires                        metric.Int64ObservableCounter
		idleConns                            metric.Int64ObservableUpDownCounter
		maxConns                             metric.Int64ObservableGauge
		maxIdleDestroyCount                  metric.Int64ObservableCounter
		maxLifetimeDestroyCountifetimeClosed metric.Int64ObservableCounter
		newConnsCount                        metric.Int64ObservableCounter
		totalConns                           metric.Int64ObservableUpDownCounter

		dbStats     *pgxpool.Stat
		lastDBStats time.Time

		// lock prevents a race between batch observer and metric registration.
		lock sync.Mutex
	)

	options := metricsOptions{
		meterProvider:              otel.GetMeterProvider(),
		minimumReadDBStatsInterval: defaultMinimumReadDBStatsInterval,
		defaultAttributes: []attribute.KeyValue{
			semconv.DBSystemPostgreSQL,
		},
	}

	for _, opt := range opts {
		opt.applyMetricsOptions(&options)
	}

	meter := options.meterProvider.Meter(otelName)

	lock.Lock()
	defer lock.Unlock()

	if acquireCount, err = meter.Int64ObservableCounter(
		pgxPoolAcquireCount,
		metric.WithUnit(unitDimensionless),
		metric.WithDescription("Cumulative count of successful acquires from the pool."),
	); err != nil {
		logger.Error.Printf("could not create acquireCount metric: %s", err)
		return err
	}

	if acquireDuration, err = meter.Float64ObservableCounter(
		pgxpoolAcquireDuration,
		metric.WithUnit(unitDimensionless),
		metric.WithDescription("Total duration of all successful acquires from the pool in nanoseconds."),
	); err != nil {
		logger.Error.Printf("could not create acquireDuration metric: %s", err)
		return err
	}

	if acquiredConns, err = meter.Int64ObservableUpDownCounter(
		pgxpoolAcquiredConns,
		metric.WithUnit(unitDimensionless),
		metric.WithDescription("Number of currently acquired connections in the pool."),
	); err != nil {
		logger.Error.Printf("could not create acquireConns metric: %s", err)
		return err
	}

	if cancelledAcquires, err = meter.Int64ObservableCounter(
		pgxpoolCancelledAcquires,
		metric.WithUnit(unitDimensionless),
		metric.WithDescription("Cumulative count of acquires from the pool that were canceled by a context."),
	); err != nil {
		logger.Error.Printf("could not create cancelledAcquires metric: %s", err)
		return err
	}

	if constructingConns, err = meter.Int64ObservableUpDownCounter(
		pgxpoolConstructingConns,
		metric.WithUnit(unitMilliseconds),
		metric.WithDescription("Number of conns with construction in progress in the pool."),
	); err != nil {
		logger.Error.Printf("could not create constructingConns metric: %s", err)
		return err
	}

	if emptyAcquires, err = meter.Int64ObservableCounter(
		pgxpoolEmptyAcquire,
		metric.WithUnit(unitDimensionless),
		metric.WithDescription("Cumulative count of successful acquires from the pool that waited for a resource to be released or constructed because the pool was empty."),
	); err != nil {
		logger.Error.Printf("could not create emptyAcquires metric: %s", err)
		return err
	}

	if idleConns, err = meter.Int64ObservableUpDownCounter(
		pgxpoolIdleConns,
		metric.WithUnit(unitDimensionless),
		metric.WithDescription("Number of currently idle conns in the pool."),
	); err != nil {
		logger.Error.Printf("could not create idleConns metric: %s", err)
		return err
	}

	if maxConns, err = meter.Int64ObservableGauge(
		pgxpoolMaxConns,
		metric.WithUnit(unitDimensionless),
		metric.WithDescription("Maximum size of the pool."),
	); err != nil {
		logger.Error.Printf("could not create maxConns metric: %s", err)
		return err
	}

	if maxIdleDestroyCount, err = meter.Int64ObservableCounter(
		pgxpoolMaxIdleDestroyCount,
		metric.WithUnit(unitDimensionless),
		metric.WithDescription("Cumulative count of connections destroyed because they exceeded MaxConnIdleTime."),
	); err != nil {
		logger.Error.Printf("could not create maxIdleDestroyCount metric: %s", err)
		return err
	}

	if maxLifetimeDestroyCountifetimeClosed, err = meter.Int64ObservableCounter(
		pgxpoolMaxLifetimeDestroyCount,
		metric.WithUnit(unitDimensionless),
		metric.WithDescription("Cumulative count of connections destroyed because they exceeded MaxConnLifetime."),
	); err != nil {
		logger.Error.Printf("could not create maxLifetimeDestroyCountifetimeClosed metric: %s", err)
		return err
	}

	if newConnsCount, err = meter.Int64ObservableCounter(
		pgxpoolNewConnsCount,
		metric.WithUnit(unitDimensionless),
		metric.WithDescription("Cumulative count of new connections opened."),
	); err != nil {
		logger.Error.Printf("could not create newConnsCount metric: %s", err)
		return err
	}

	if totalConns, err = meter.Int64ObservableUpDownCounter(
		pgxpoolTotalConns,
		metric.WithUnit(unitDimensionless),
		metric.WithDescription("Total number of resources currently in the pool. The value is the sum of ConstructingConns, AcquiredConns, and IdleConns."),
	); err != nil {
		logger.Error.Printf("could not create totalConns metric: %s", err)
		return err
	}

	_, err = meter.RegisterCallback(
		func(ctx context.Context, o metric.Observer) error {
			lock.Lock()
			defer lock.Unlock()

			now := time.Now()
			if now.Sub(lastDBStats) >= options.minimumReadDBStatsInterval {
				dbStats = db.Stat()
				lastDBStats = now
			}

			o.ObserveInt64(acquireCount, dbStats.AcquireCount(), metric.WithAttributes(options.defaultAttributes...))
			o.ObserveFloat64(acquireDuration, float64(dbStats.AcquireDuration())/1e6, metric.WithAttributes(options.defaultAttributes...))
			o.ObserveInt64(acquiredConns, int64(dbStats.AcquiredConns()), metric.WithAttributes(options.defaultAttributes...))
			o.ObserveInt64(cancelledAcquires, dbStats.CanceledAcquireCount(), metric.WithAttributes(options.defaultAttributes...))
			o.ObserveInt64(constructingConns, int64(dbStats.ConstructingConns()), metric.WithAttributes(options.defaultAttributes...))
			o.ObserveInt64(emptyAcquires, dbStats.EmptyAcquireCount(), metric.WithAttributes(options.defaultAttributes...))
			o.ObserveInt64(idleConns, int64(dbStats.IdleConns()), metric.WithAttributes(options.defaultAttributes...))
			o.ObserveInt64(maxConns, int64(dbStats.MaxConns()), metric.WithAttributes(options.defaultAttributes...))
			o.ObserveInt64(maxIdleDestroyCount, dbStats.MaxIdleDestroyCount(), metric.WithAttributes(options.defaultAttributes...))
			o.ObserveInt64(maxLifetimeDestroyCountifetimeClosed, dbStats.MaxLifetimeDestroyCount(), metric.WithAttributes(options.defaultAttributes...))
			o.ObserveInt64(newConnsCount, dbStats.NewConnsCount(), metric.WithAttributes(options.defaultAttributes...))
			o.ObserveInt64(totalConns, int64(dbStats.TotalConns()), metric.WithAttributes(options.defaultAttributes...))

			return nil
		},
		acquireCount,
		acquireDuration,
		acquiredConns,
		cancelledAcquires,
		constructingConns,
		emptyAcquires,
		idleConns,
		maxConns,
		maxIdleDestroyCount,
		maxLifetimeDestroyCountifetimeClosed,
		newConnsCount,
		totalConns,
	)

	if err != nil {
		logger.Error.Printf("Could not register callback for DB Metrics: %s", err)
	}
	return err
}
