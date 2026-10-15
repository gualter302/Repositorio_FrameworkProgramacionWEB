package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"multicatalogo-backend/repository"
)

// GetRed devuelve la red de referidos reconstruida como un árbol multinivel.
func GetRed(c *fiber.Ctx) error {
	emailUsuario := strings.ToLower(strings.TrimSpace(c.Query("email")))
	if emailUsuario == "" {
		emailUsuario = "admin@upse.edu.ec"
	}

	red, desdeCache, err := repository.ObtenerRed(c.Context(), emailUsuario)
	if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudo obtener la red de referidos", nil)
	}
	c.Set("X-Cache", estadoCache(desdeCache))
	return c.JSON(red)
}
