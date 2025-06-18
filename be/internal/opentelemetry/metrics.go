package opentelemetry

import (
	"context"
	"time"

	"github.com/lolwierd/weatherboy/be/internal/constants"
	"github.com/lolwierd/weatherboy/be/internal/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
)

func setupMetrics(conn *grpc.ClientConn) (*sdkmetric.MeterProvider, error) {
	ctx := context.Background()
	ctxtimeout, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()
	exporter, err := otlpmetricgrpc.New(
		ctxtimeout,
		otlpmetricgrpc.WithGRPCConn(conn),
	)
	if err != nil {
		logger.Error.Printf("failed to create a new metric client")
		return nil, err
	}

	// labels/tags/resources that are common to all metrics.
	attributes := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(constants.SERVICE_NAME),
	)

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(attributes),
		sdkmetric.WithReader(
			// collects and exports metric data every 30 seconds.
			sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(30*time.Second)),
		),
	)

	otel.SetMeterProvider(mp)

	return mp, nil
}
