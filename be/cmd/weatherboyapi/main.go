package main

import (
	"context"
	"flag"

	"github.com/lolwierd/weatherboy/be/internal/fetch"
	"github.com/lolwierd/weatherboy/be/internal/healthcheck"
	"github.com/lolwierd/weatherboy/be/internal/logger"
	"github.com/lolwierd/weatherboy/be/internal/router"
	"github.com/lolwierd/weatherboy/be/internal/shutdown"
)

var runMode = flag.String("run", "server", "run mode")

func main() {
	flag.Parse()

	switch *runMode {
	case "server":
		logger.Info.Println("Starting API server mode.")
		healthcheck.Healthcheck()
		router.StartServer()
		shutdown.GracefulStop()
	case "fetch_bulletin_once":
		if err := fetch.FetchBulletinOnce(context.Background()); err != nil {
			logger.Error.Println("fetch bulletin:", err)
		}
	default:
		logger.Error.Println("unknown run mode", *runMode)
	}
}
