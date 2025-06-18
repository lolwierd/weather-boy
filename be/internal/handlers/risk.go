package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lolwierd/weatherboy/be/internal/logger"
	"github.com/lolwierd/weatherboy/be/internal/score"
)

func GetRisk(c *fiber.Ctx) error {
	loc := c.Params("loc")
	res, err := score.RiskLevel(c.Context(), loc)
	if err != nil {
		logger.Error.Println("risk level:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}
