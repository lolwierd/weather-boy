package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lolwierd/weatherboy/be/internal/healthcheck"
)

func Health(c *fiber.Ctx) error {
	if !healthcheck.IsHealthy {
		c.SendStatus(fiber.StatusTeapot)
		c.SendString("Shuting Down")
	} else {
		c.SendStatus(fiber.StatusOK)
		c.SendString("OK")
	}

	return nil
}
