package controllers

import (
	"github.com/gofiber/fiber/v2"
	"multicatalogo-backend/models"
)

// responderError evita que cada controlador construya los errores de una forma distinta.
func responderError(c *fiber.Ctx, status int, message string, details interface{}) error {
	return c.Status(status).JSON(models.APIError{
		Status:  status,
		Message: message,
		Details: details,
	})
}
