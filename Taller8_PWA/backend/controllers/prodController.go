package controllers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"multicatalogo-backend/models"
	"multicatalogo-backend/repository"
)

// GetProductos devuelve todos los artículos disponibles.
func (api *API) GetProductos(c *fiber.Ctx) error {
	if productos, encontrado := api.cache.ObtenerProductos(c.Context()); encontrado {
		c.Set("X-Cache", estadoCache(true))
		return c.JSON(productos)
	}

	productos, err := api.datos.ObtenerProductos(c.Context())
	if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudo obtener el catálogo", nil)
	}
	api.cache.GuardarProductos(c.Context(), productos)
	c.Set("X-Cache", estadoCache(false))
	return c.JSON(productos)
}

func estadoCache(desdeCache bool) string {
	if desdeCache {
		return "HIT"
	}
	return "MISS"
}

// GetProductoPorID busca un producto usando el número recibido en la dirección.
func (api *API) GetProductoPorID(c *fiber.Ctx) error {
	idTexto := c.Params("id")
	id, err := strconv.Atoi(idTexto)
	if err != nil || id <= 0 {
		return responderError(c, fiber.StatusBadRequest, "El id del producto debe ser un número entero positivo", fiber.Map{"id": idTexto})
	}

	producto, encontrado, err := api.datos.ObtenerProductoPorID(c.Context(), id)
	if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudo obtener el producto", nil)
	}
	if !encontrado {
		return responderError(c, fiber.StatusNotFound, "Producto no encontrado", fiber.Map{"id": id})
	}

	return c.JSON(producto)
}

func validarProducto(producto *models.Producto) string {
	producto.Nombre = strings.TrimSpace(producto.Nombre)
	producto.Descripcion = strings.TrimSpace(producto.Descripcion)
	producto.Categoria = strings.TrimSpace(producto.Categoria)
	producto.Img = strings.TrimSpace(producto.Img)
	if producto.Nombre == "" {
		return "El nombre del producto es obligatorio"
	}
	if producto.Precio < 0 {
		return "El precio no puede ser negativo"
	}
	if producto.Stock < 0 {
		return "El stock no puede ser negativo"
	}
	return ""
}

func (api *API) AdminProductos(c *fiber.Ctx) error {
	productos, err := api.datos.ObtenerTodosLosProductos(c.Context())
	if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudo obtener el catálogo administrativo", nil)
	}
	return c.JSON(productos)
}

func (api *API) CrearProducto(c *fiber.Ctx) error {
	var producto models.Producto
	if err := c.BodyParser(&producto); err != nil {
		return responderError(c, fiber.StatusBadRequest, "Cuerpo de petición inválido", nil)
	}
	if mensaje := validarProducto(&producto); mensaje != "" {
		return responderError(c, fiber.StatusBadRequest, mensaje, nil)
	}
	creado, err := api.datos.CrearProducto(c.Context(), producto)
	if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudo crear el producto", nil)
	}
	if err := api.cache.InvalidarProductos(c.Context()); err != nil {
		api.logger.Printf("no se pudo invalidar la caché de productos: %v", err)
	}
	return c.Status(fiber.StatusCreated).JSON(creado)
}

func (api *API) ActualizarProducto(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return responderError(c, fiber.StatusBadRequest, "El id del producto debe ser un número entero positivo", nil)
	}
	var producto models.Producto
	if err := c.BodyParser(&producto); err != nil {
		return responderError(c, fiber.StatusBadRequest, "Cuerpo de petición inválido", nil)
	}
	if mensaje := validarProducto(&producto); mensaje != "" {
		return responderError(c, fiber.StatusBadRequest, mensaje, nil)
	}
	actualizado, err := api.datos.ActualizarProducto(c.Context(), id, producto)
	if errors.Is(err, repository.ErrProductoNoEncontrado) {
		return responderError(c, fiber.StatusNotFound, "Producto no encontrado", nil)
	}
	if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudo actualizar el producto", nil)
	}
	if err := api.cache.InvalidarProductos(c.Context()); err != nil {
		api.logger.Printf("no se pudo invalidar la caché de productos: %v", err)
	}
	return c.JSON(actualizado)
}

func (api *API) DesactivarProducto(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return responderError(c, fiber.StatusBadRequest, "El id del producto debe ser un número entero positivo", nil)
	}
	err = api.datos.DesactivarProducto(c.Context(), id)
	if errors.Is(err, repository.ErrProductoNoEncontrado) {
		return responderError(c, fiber.StatusNotFound, "Producto no encontrado", nil)
	}
	if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudo desactivar el producto", nil)
	}
	if err := api.cache.InvalidarProductos(c.Context()); err != nil {
		api.logger.Printf("no se pudo invalidar la caché de productos: %v", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
