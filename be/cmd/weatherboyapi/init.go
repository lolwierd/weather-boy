package main

import (
	"fmt"
	"os"

	"github.com/lolwierd/weatherboy/be/internal/db"
	"github.com/lolwierd/weatherboy/be/internal/opentelemetry"
)

func init() {
	// Setup OTEL
	opentelemetry.InitOtel()

	// Setup DB
	db.InitDBPool(
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_DB"),
			os.Getenv("POSTGRES_PORT")),
	)
}
