package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lolwierd/weatherboy/be/internal/logger"
	"github.com/lolwierd/weatherboy/be/internal/repository"
)

func GetAWSARG(c *fiber.Ctx) error {
	loc := c.Params("loc")
	r, err := repository.LatestAWSARG(c.Context(), loc)
	if err != nil {
		logger.Error.Println("aws/arg fetch:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(r)
}
