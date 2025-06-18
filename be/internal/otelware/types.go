package otelware

import (
	"github.com/gofiber/fiber/v2"
	otelmetric "go.opentelemetry.io/otel/metric"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName = "github.com/excloud-in/otelware"

	metricNameHttpServerDuration       = "http.server.duration"
	metricNameHttpServerRequestSize    = "http.server.request.size"
	metricNameHttpServerResponseSize   = "http.server.response.size"
	metricNameHttpServerActiveRequests = "http.server.active_requests"

	// Unit constants for deprecated metric units
	unitDimensionless = "1"
	unitBytes         = "By"
	unitMilliseconds  = "ms"
)

type config struct {
	Next           func(*fiber.Ctx) bool
	TracerProvider oteltrace.TracerProvider
	TracerKey      string
	MeterProvider  otelmetric.MeterProvider
	Port           *int
	ServerName     *string
}

// Option specifies instrumentation configuration options.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)
