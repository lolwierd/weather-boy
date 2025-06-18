package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lolwierd/weatherboy/be/internal/utils"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// TraceQueryStart is called at the beginning of Query, QueryRow, and Exec calls.
// The returned context is used for the rest of the call and will be passed to TraceQueryEnd.
func (t *tracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
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

	spanName := "query " + sqlOperationName(data.SQL)

	ctx, _ = t.tracer.Start(ctx, spanName, opts...)

	return ctx
}

// TraceQueryEnd is called at the end of Query, QueryRow, and Exec calls.
func (t *tracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	span := trace.SpanFromContext(ctx)
	utils.RecordError(span, data.Err)

	span.SetAttributes(RowsAffectedKey.Int(int(data.CommandTag.RowsAffected())))

	span.End()
}
