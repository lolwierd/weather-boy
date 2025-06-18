package router

import (
	"github.com/lolwierd/weatherboy/be/internal/handlers"
)

func RegisterRoutes() {
	App.Get("/health", handlers.Health)
}
