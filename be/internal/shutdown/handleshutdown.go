package shutdown

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/lolwierd/weatherboy/be/internal/db"
	"github.com/lolwierd/weatherboy/be/internal/healthcheck"
	"github.com/lolwierd/weatherboy/be/internal/logger"
	"github.com/lolwierd/weatherboy/be/internal/opentelemetry"
	"github.com/lolwierd/weatherboy/be/internal/router"
)

var WG sync.WaitGroup
var IsShuttingDown bool

// This function will handle graceful shutdown of the application in a proper order and blocks until SIGTERM
func GracefulStop() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.Printf("Caught signal %v, shutting down application.", <-sigChan)
	IsShuttingDown = true

	//First become unhealthy and wait 5 seconds to drain.
	healthcheck.IsHealthy = false
	time.Sleep(5 * time.Second)

	//Close fiber connections
	router.App.Shutdown()

	logger.Info.Println("Waiting for all goroutines to finish")
	// Wait for all goroutines to complete.
	WG.Wait()

	//Close All DB connections
	db.GetDBDriver().ConnPool.Close()

	//Close OTEL providers
	if err := opentelemetry.TracerProvider.Shutdown(context.Background()); err != nil {
		log.Printf("Error shutting down tracer provider: %v", err)
	}
	if err := opentelemetry.MeterProvider.Shutdown(context.Background()); err != nil {
		log.Printf("Error shutting down metrics provider: %v", err)
	}
}
