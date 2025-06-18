package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lolwierd/weatherboy/be/internal/logger"
)

var dbDriver *Driver

func GetDBDriver() *Driver {
	return dbDriver
}

// TODO think about cleaning up idle DB connections.
func InitDBPool(dsn string) {
	dbDriver = new(Driver)
	connConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Error.Fatalln("Could not parse DB ConnConfig: ", err)
	}
	connConfig.ConnConfig.Tracer = NewpgxTracer()

	connPool, err := pgxpool.NewWithConfig(context.Background(), connConfig)
	if err != nil {
		logger.Error.Fatalln("Could not create DB Pool: ", err)
	}
	err = connPool.Ping(context.Background())
	if err != nil {
		logger.Error.Fatalln("unable to ping database: ", err)
	}
	if err := recordMetrics(connPool); err != nil {
		logger.Error.Fatalln("unable to record database stats: ", err)
	}
	dbDriver.ConnPool = connPool
}
