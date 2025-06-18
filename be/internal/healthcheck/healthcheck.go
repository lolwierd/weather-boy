package healthcheck

import (
	"context"
	"time"

	"github.com/lolwierd/weatherboy/be/internal/db"
	"github.com/lolwierd/weatherboy/be/internal/logger"
	"github.com/lolwierd/weatherboy/be/internal/opentelemetry"
)

var IsHealthy = false

func Healthcheck() {
	go func() {
		for {
			isAnyoneUnhealthy := false
			if !isDBHealthy() {
				logger.Error.Println("db is unhealthy. setting API unhealthy.")
				IsHealthy = false
				isAnyoneUnhealthy = true
			}
			if !isOTELHealthy() {
				logger.Warn.Println("otel is unhealthy")
			}
			if !isAnyoneUnhealthy {
				IsHealthy = true
			}
			time.Sleep(time.Second * 1)
		}
	}()
}

func isOTELHealthy() bool {
	return opentelemetry.IsOtelConnHealthy()
}

func isDBHealthy() bool {
	err := db.GetDBDriver().ConnPool.Ping(context.Background())
	if err != nil {
		logger.Error.Println("unable to ping database: ", err)
		logger.Error.Println("marking API as unhealthy")
	}
	return err == nil
}
