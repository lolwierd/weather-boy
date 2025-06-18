package opentelemetry

import (
	"context"
	"os"
	"time"

	"github.com/lolwierd/weatherboy/be/internal/logger"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	MeterProvider  *sdkmetric.MeterProvider
	TracerProvider *sdktace.TracerProvider
	conn           *grpc.ClientConn
)

func InitOtel() {
	var err error
	ctx := context.Background()
	otelcollectorhost := os.Getenv("OTELCOL_HOST") + ":" + os.Getenv("OTELCOL_PORT")
	ctxtimeout, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()
	conn, err = grpc.DialContext(ctxtimeout, otelcollectorhost,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		logger.Error.Printf("failed to create gRPC connection to collector: %s", otelcollectorhost)
	}

	MeterProvider, err = setupMetrics(conn)
	if err != nil {
		logger.Error.Fatalln("Could not bootstrap OTEL MeterProvider")
	}
	TracerProvider, err = setupTracer(conn)
	if err != nil {
		logger.Error.Fatalln("Could not bootstrap OTEL TraceProvider")
	}
}

func IsOtelConnHealthy() bool {
	return (conn.GetState() == connectivity.Idle || conn.GetState() == connectivity.Ready)
}
