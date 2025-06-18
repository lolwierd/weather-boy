package router

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/lolwierd/weatherboy/be/internal/logger"
	"github.com/lolwierd/weatherboy/be/internal/otelware"
)

var App *fiber.App

// @title						Excloud Billing API
// @version					1.0
// @description				API for Excloud billing and cost management.
// @servers.url				https://billing.excloud.in
// @externalDocs.description	Docs for excloud
// @externalDocs.url			https://docs.excloud.in
func StartServer() {
	App = fiber.New(fiber.Config{
		AppName:               "weatherboy/be",
		ReadTimeout:           5 * time.Second,
		WriteTimeout:          5 * time.Second,
		IdleTimeout:           30 * time.Second,
		DisableStartupMessage: true,
	})

	// Middleware
	App.Use(recover.New())
	App.Use(cors.New(cors.Config{
		AllowHeaders: "*",
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	App.Use(
		otelware.New(
			otelware.WithServerName("weatherboy/be"),
			otelware.WithServerName("weatherboy/be-tracer"),
		),
	)

	RegisterRoutes()

	go func() {
		logger.Info.Println("Starting server on " + os.Getenv("LISTEN_ADDR"))
		err := App.Listen(os.Getenv("LISTEN_ADDR"))
		if err != nil {
			logger.Error.Fatalln("Could not start API Server: ", err)
		}
	}()
}
