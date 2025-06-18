package otelware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func getTracerMeter(cfg config) (tracer oteltrace.Tracer, meter metric.Meter) {
	if cfg.TracerProvider == nil {
		cfg.TracerProvider = otel.GetTracerProvider()
	}
	tracer = cfg.TracerProvider.Tracer(
		instrumentationName,
	)

	if cfg.MeterProvider == nil {
		cfg.MeterProvider = otel.GetMeterProvider()
	}
	meter = cfg.MeterProvider.Meter(
		instrumentationName,
	)
	return
}

func getTracerKey(cfg config) (tracerKey string) {
	if cfg.TracerKey == "" {
		cfg.TracerKey = "excloud-tracer"
	}
	tracerKey = cfg.TracerKey
	return
}

func getIP(c *fiber.Ctx) string {
	if c.Get("Cf-Connecting-Ip") != "" {
		return utils.CopyString(c.Get("Cf-Connecting-Ip"))
	} else if c.Get("X-Forwarded-For") != "" {
		return utils.CopyString(c.Get("X-Forwarded-For"))
	} else {
		return c.IP()
	}
}
