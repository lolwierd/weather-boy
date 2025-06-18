package db

import (
	"strings"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// connectionAttributesFromConfig returns a slice of SpanStartOptions that contain
// attributes from the given connection config.
func connectionAttributesFromConfig(config *pgx.ConnConfig) []trace.SpanStartOption {
	if config != nil {
		return []trace.SpanStartOption{
			trace.WithAttributes(attribute.String(string(semconv.NetPeerNameKey), config.Host)),
			trace.WithAttributes(attribute.Int(string(semconv.NetPeerPortKey), int(config.Port))),
			trace.WithAttributes(attribute.String(string(semconv.DBUserKey), config.User)),
		}
	}
	return nil
}

// sqlOperationName attempts to get the first 'word' from a given SQL query, which usually
// is the operation name (e.g. 'SELECT').
func sqlOperationName(stmt string) string {
	parts := strings.Fields(stmt)
	if len(parts) == 0 {
		// Fall back to a fixed value to prevent creating lots of tracing operations
		// differing only by the amount of whitespace in them (in case we'd fall back
		// to the full query or a cut-off version).
		return sqlOperationUnknown
	}
	return strings.ToUpper(parts[0])
}
