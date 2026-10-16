package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// GetRed devuelve la red de referidos reconstruida como un árbol multinivel.
func (api *API) GetRed(c *fiber.Ctx) error {
	emailUsuario, _ := c.Locals("usuarioEmail").(string)
	emailUsuario = strings.ToLower(strings.TrimSpace(emailUsuario))
	if emailUsuario == "" {
		return responderError(c, fiber.StatusUnauthorized, "Debes iniciar sesión", nil)
	}

	if red, encontrada := api.cache.ObtenerRed(c.Context(), emailUsuario); encontrada {
		c.Set("X-Cache", estadoCache(true))
		return c.JSON(red)
	}

	red, err := api.datos.ObtenerRed(c.Context(), emailUsuario)
	if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudo obtener la red de referidos", nil)
	}
	api.cache.GuardarRed(c.Context(), emailUsuario, red)
	c.Set("X-Cache", estadoCache(false))
	return c.JSON(red)
}
