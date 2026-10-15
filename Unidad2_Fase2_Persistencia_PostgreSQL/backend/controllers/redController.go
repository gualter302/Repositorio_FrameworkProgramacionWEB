package controllers

import (
	"github.com/gofiber/fiber/v2"
	"multicatalogo-backend/repository"
)

// GetRed devuelve la red de referidos reconstruida como un árbol multinivel.
func GetRed(c *fiber.Ctx) error {
	red, err := repository.ObtenerRed(c.Context())
	if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudo obtener la red de referidos", nil)
	}
	return c.JSON(red)
}
