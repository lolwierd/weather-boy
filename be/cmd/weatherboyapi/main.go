package main

import (
	"github.com/lolwierd/weatherboy/be/internal/healthcheck"
	"github.com/lolwierd/weatherboy/be/internal/logger"
	"github.com/lolwierd/weatherboy/be/internal/router"
	"github.com/lolwierd/weatherboy/be/internal/shutdown"
)

func main() {
	logger.Info.Println("Starting API server mode.")
	healthcheck.Healthcheck()
	router.StartServer()
	shutdown.GracefulStop()
}
