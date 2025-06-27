package router

import (
	"github.com/lolwierd/weatherboy/be/internal/handlers"
)

func RegisterRoutes() {
	App.Get("/health", handlers.Health)

	v1 := App.Group("/v1")
	v1.Get("/risk/:loc", handlers.GetRisk)
	v1.Get("/bulletin/:loc", handlers.GetBulletin)
	v1.Get("/nowcast/:loc", handlers.GetNowcast)
}
