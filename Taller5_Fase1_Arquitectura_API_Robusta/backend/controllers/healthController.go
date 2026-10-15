package controllers

import "github.com/gofiber/fiber/v2"

// HealthCheck permite comprobar rápidamente que el backend está funcionando.
func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"message": "Servidor Go/Fiber operativo",
	})
}
