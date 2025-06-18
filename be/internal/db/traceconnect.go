package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lolwierd/weatherboy/be/internal/utils"
	"go.opentelemetry.io/otel/trace"
)

// TraceConnectStart is called at the beginning of Connect and ConnectConfig
// calls. The returned context is used for the rest of the call and will be
// passed to TraceConnectEnd.
func (t *tracer) TraceConnectStart(ctx context.Context, data pgx.TraceConnectStartData) context.Context {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(t.attrs...),
	}

	if data.ConnConfig != nil {
		opts = append(opts, connectionAttributesFromConfig(data.ConnConfig)...)
	}

	ctx, _ = t.tracer.Start(ctx, "connect", opts...)

	return ctx
}

// TraceConnectEnd is called at the end of Connect and ConnectConfig calls.
func (t *tracer) TraceConnectEnd(ctx context.Context, data pgx.TraceConnectEndData) {
	span := trace.SpanFromContext(ctx)
	utils.RecordError(span, data.Err)

	span.End()
}
