package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lolwierd/weatherboy/be/internal/utils"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// TraceBatchStart is called at the beginning of SendBatch calls. The returned
// context is used for the rest of the call and will be passed to
// TraceBatchQuery and TraceBatchEnd.
func (t *tracer) TraceBatchStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchStartData) context.Context {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx
	}

	var size int
	if b := data.Batch; b != nil {
		size = b.Len()
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(t.attrs...),
		trace.WithAttributes(attribute.Int("pgx.batch.size", size)),
	}

	if conn != nil {
		opts = append(opts, connectionAttributesFromConfig(conn.Config())...)
	}

	ctx, _ = t.tracer.Start(ctx, "batch start", opts...)

	return ctx
}

// TraceBatchQuery is called at the after each query in a batch.
func (t *tracer) TraceBatchQuery(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchQueryData) {
	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(t.attrs...),
	}

	if conn != nil {
		opts = append(opts, connectionAttributesFromConfig(conn.Config())...)
	}

	opts = append(opts, trace.WithAttributes(semconv.DBStatementKey.String(data.SQL)))

	spanName := "query " + sqlOperationName(data.SQL)

	_, span := t.tracer.Start(ctx, spanName, opts...)
	utils.RecordError(span, data.Err)

	span.End()
}

// TraceBatchEnd is called at the end of SendBatch calls.
func (t *tracer) TraceBatchEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceBatchEndData) {
	span := trace.SpanFromContext(ctx)
	utils.RecordError(span, data.Err)

	span.End()
}
