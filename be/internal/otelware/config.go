package otelware

import (
	"github.com/gofiber/fiber/v2"
	otelmetric "go.opentelemetry.io/otel/metric"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// config is used to configure the Fiber middleware.
func (o optionFunc) apply(c *config) {
	o(c)
}

// WithNext takes a function that will be called on every
// request, the middleware will be skipped if returning true
func WithNext(f func(ctx *fiber.Ctx) bool) Option {
	return optionFunc(func(cfg *config) {
		cfg.Next = f
	})
}

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
// If none is specified, the global provider is used.
func WithTracerProvider(provider oteltrace.TracerProvider) Option {
	return optionFunc(func(cfg *config) {
		cfg.TracerProvider = provider
	})
}

// WithTracerKey specifies a tracer key to use with a tracer.
// If none is specified, the excloud-tracer key is used.
func WithTracerKey(trackerKey string) Option {
	return optionFunc(func(cfg *config) {
		cfg.TracerKey = trackerKey
	})
}

// WithMeterProvider specifies a meter provider to use for reporting.
// If none is specified, the global provider is used.
func WithMeterProvider(provider otelmetric.MeterProvider) Option {
	return optionFunc(func(cfg *config) {
		cfg.MeterProvider = provider
	})
}

// WithServerName specifies the value to use when setting the `http.server_name`
// attribute on metrics/spans.
func WithServerName(serverName string) Option {
	return optionFunc(func(cfg *config) {
		cfg.ServerName = &serverName
	})
}

// WithPort specifies the value to use when setting the `net.host.port`
// attribute on metrics/spans. Attribute is "Conditionally Required: If not
// default (`80` for `http`, `443` for `https`).
func WithPort(port int) Option {
	return optionFunc(func(cfg *config) {
		cfg.Port = &port
	})
}
