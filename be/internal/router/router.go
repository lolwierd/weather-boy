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
	v1.Get("/radar/:loc", handlers.GetRadar)
	v1.Get("/riverbasin/:loc", handlers.GetRiverBasin)
	v1.Get("/awsarg/:loc", handlers.GetAWSARG)
}
