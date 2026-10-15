package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"multicatalogo-backend/repository"
)

// GetProductos devuelve todos los artículos disponibles.
func GetProductos(c *fiber.Ctx) error {
	return c.JSON(repository.ObtenerProductos())
}

// GetProductoPorID busca un producto usando el número recibido en la dirección.
func GetProductoPorID(c *fiber.Ctx) error {
	idTexto := c.Params("id")
	id, err := strconv.Atoi(idTexto)
	if err != nil || id <= 0 {
		return responderError(
			c,
			fiber.StatusBadRequest,
			"El id del producto debe ser un número entero positivo",
			fiber.Map{"id": idTexto},
		)
	}

	producto, encontrado := repository.ObtenerProductoPorID(id)
	if !encontrado {
		return responderError(c, fiber.StatusNotFound, "Producto no encontrado", fiber.Map{"id": id})
	}

	return c.JSON(producto)
}
