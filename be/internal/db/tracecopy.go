package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lolwierd/weatherboy/be/internal/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// TraceCopyFromStart is called at the beginning of CopyFrom calls. The
// returned context is used for the rest of the call and will be passed to
// TraceCopyFromEnd.
func (t *tracer) TraceCopyFromStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceCopyFromStartData) context.Context {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(t.attrs...),
		trace.WithAttributes(attribute.String("db.table", data.TableName.Sanitize())),
	}

	if conn != nil {
		opts = append(opts, connectionAttributesFromConfig(conn.Config())...)
	}

	ctx, _ = t.tracer.Start(ctx, "copy_from "+data.TableName.Sanitize(), opts...)

	return ctx
}

// TraceCopyFromEnd is called at the end of CopyFrom calls.
func (t *tracer) TraceCopyFromEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceCopyFromEndData) {
	span := trace.SpanFromContext(ctx)
	utils.RecordError(span, data.Err)

	span.SetAttributes(RowsAffectedKey.Int(int(data.CommandTag.RowsAffected())))

	span.End()
}
