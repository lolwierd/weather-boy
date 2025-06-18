package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lolwierd/weatherboy/be/internal/utils"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// TracePrepareStart is called at the beginning of Prepare calls. The returned
// context is used for the rest of the call and will be passed to
// TracePrepareEnd.
func (t *tracer) TracePrepareStart(ctx context.Context, conn *pgx.Conn, data pgx.TracePrepareStartData) context.Context {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(t.attrs...),
	}
	if conn != nil {
		opts = append(opts, connectionAttributesFromConfig(conn.Config())...)
	}
	opts = append(opts, trace.WithAttributes(semconv.DBStatementKey.String(data.SQL)))

	spanName := "prepare " + sqlOperationName(data.SQL)

	ctx, _ = t.tracer.Start(ctx, spanName, opts...)
	return ctx
}

// TracePrepareEnd is called at the end of Prepare calls.
func (t *tracer) TracePrepareEnd(ctx context.Context, _ *pgx.Conn, data pgx.TracePrepareEndData) {
	span := trace.SpanFromContext(ctx)
	utils.RecordError(span, data.Err)

	span.End()
}
