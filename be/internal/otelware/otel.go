package otelware

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func New(opts ...Option) fiber.Handler {
	cfg := config{}
	for _, opt := range opts {
		opt.apply(&cfg)
	}
	tracer, meter := getTracerMeter(cfg)
	tracerKey := getTracerKey(cfg)

	httpServerDuration, err := meter.Float64Histogram(metricNameHttpServerDuration, metric.WithUnit(unitMilliseconds), metric.WithDescription("measures the duration inbound HTTP requests"))
	if err != nil {
		otel.Handle(err)
	}
	httpServerRequestSize, err := meter.Int64Histogram(metricNameHttpServerRequestSize, metric.WithUnit(unitBytes), metric.WithDescription("measures the size of HTTP request messages"))
	if err != nil {
		otel.Handle(err)
	}
	httpServerResponseSize, err := meter.Int64Histogram(metricNameHttpServerResponseSize, metric.WithUnit(unitBytes), metric.WithDescription("measures the size of HTTP response messages"))
	if err != nil {
		otel.Handle(err)
	}
	httpServerActiveRequests, err := meter.Int64UpDownCounter(metricNameHttpServerActiveRequests, metric.WithUnit(unitDimensionless), metric.WithDescription("measures the number of concurrent HTTP requests that are currently in-flight"))
	if err != nil {
		otel.Handle(err)
	}

	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		c.Locals(tracerKey, tracer)

		savedCtx, cancel := context.WithCancel(c.UserContext())

		start := time.Now()

		requestMetricsAttrs := httpServerMetricAttributesFromRequest(c, cfg)
		httpServerActiveRequests.Add(savedCtx, 1, metric.WithAttributes(requestMetricsAttrs...))

		responseMetricAttrs := make([]attribute.KeyValue, len(requestMetricsAttrs))
		copy(responseMetricAttrs, requestMetricsAttrs)

		reqHeader := make(http.Header)
		c.Request().Header.VisitAll(func(k, v []byte) {
			reqHeader.Add(string(k), string(v))
		})
		ctx := context.Background()
		opts := []oteltrace.SpanStartOption{
			oteltrace.WithAttributes(httpServerTraceAttributesFromRequest(c, cfg)...),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}

		spanName := utils.CopyString(c.Path())
		ctx, span := tracer.Start(ctx, spanName, opts...)
		defer span.End()

		c.SetUserContext(ctx)

		if err := c.Next(); err != nil {
			span.RecordError(err)
			_ = c.App().Config().ErrorHandler(c, err)
		}

		responseAttrs := append(
			semconv.HTTPAttributesFromHTTPStatusCode(c.Response().StatusCode()),
			semconv.HTTPRouteKey.String(c.Route().Path),
		)

		requestSize := int64(len(c.Request().Body()))
		responseSize := int64(len(c.Response().Body()))

		defer func() {
			responseMetricAttrs = append(
				responseMetricAttrs,
				responseAttrs...)

			httpServerActiveRequests.Add(savedCtx, -1, metric.WithAttributes(requestMetricsAttrs...))
			httpServerDuration.Record(savedCtx, float64(time.Since(start).Microseconds())/1000, metric.WithAttributes(responseMetricAttrs...))
			httpServerRequestSize.Record(savedCtx, requestSize, metric.WithAttributes(responseMetricAttrs...))
			httpServerResponseSize.Record(savedCtx, responseSize, metric.WithAttributes(responseMetricAttrs...))

			c.SetUserContext(savedCtx)
			cancel()
		}()

		span.SetAttributes(
			append(
				responseAttrs,
				semconv.HTTPResponseContentLengthKey.Int64(responseSize),
			)...)

		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCodeAndSpanKind(c.Response().StatusCode(), oteltrace.SpanKindServer)
		span.SetStatus(spanStatus, spanMessage)

		return nil
	}
}
