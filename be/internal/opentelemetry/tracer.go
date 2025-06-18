package opentelemetry

import (
	"context"
	"time"

	"github.com/lolwierd/weatherboy/be/internal/constants"
	"github.com/lolwierd/weatherboy/be/internal/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var Tracer = otel.Tracer(constants.SERVICE_NAME)

func setupTracer(conn *grpc.ClientConn) (*sdktrace.TracerProvider, error) {
	ctx := context.Background()
	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceName(constants.SERVICE_NAME),
		),
	)
	if err != nil {
		logger.Error.Printf("failed to create otel resource: %s", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		logger.Error.Printf("failed to create trace exporter: %s", err)
		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tracerProvider, nil
}
