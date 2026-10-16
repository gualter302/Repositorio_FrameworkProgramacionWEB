package controllers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"multicatalogo-backend/repository"
)

func (api *API) MisPedidos(c *fiber.Ctx) error {
	email, _ := c.Locals("usuarioEmail").(string)
	pedidos, err := api.datos.ObtenerPedidosUsuario(c.Context(), email)
	if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudieron obtener tus pedidos", nil)
	}
	return c.JSON(pedidos)
}

func (api *API) AdminPedidos(c *fiber.Ctx) error {
	pagina, err := strconv.Atoi(c.Query("pagina", "1"))
	if err != nil || pagina < 1 {
		return responderError(c, fiber.StatusBadRequest, "La página debe ser un entero positivo", nil)
	}
	limite, err := strconv.Atoi(c.Query("limite", "20"))
	if err != nil || limite < 1 || limite > 100 {
		return responderError(c, fiber.StatusBadRequest, "El límite debe estar entre 1 y 100", nil)
	}
	pedidos, total, err := api.datos.ObtenerTodosLosPedidos(c.Context(), strings.TrimSpace(c.Query("buscar")), pagina, limite)
	if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudieron obtener los pedidos", nil)
	}
	return c.JSON(fiber.Map{"datos": pedidos, "pagina": pagina, "limite": limite, "total": total})
}

func (api *API) CambiarEstadoPedido(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return responderError(c, fiber.StatusBadRequest, "El id del pedido no es válido", nil)
	}
	var cuerpo struct {
		Estado string `json:"estado"`
	}
	if err := c.BodyParser(&cuerpo); err != nil {
		return responderError(c, fiber.StatusBadRequest, "Cuerpo de petición inválido", nil)
	}
	estado := strings.ToLower(strings.TrimSpace(cuerpo.Estado))
	permitidos := map[string]bool{
		"pendiente": true, "procesando": true, "confirmado": true,
		"preparando": true, "enviado": true, "entregado": true,
	}
	if !permitidos[estado] {
		return responderError(c, fiber.StatusBadRequest, "El estado del pedido no es válido", nil)
	}
	if err := api.datos.ActualizarEstadoPedido(c.Context(), id, estado); errors.Is(err, repository.ErrPedidoNoEncontrado) {
		return responderError(c, fiber.StatusNotFound, "Pedido no encontrado", nil)
	} else if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudo actualizar el pedido", nil)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (api *API) CancelarPedido(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return responderError(c, fiber.StatusBadRequest, "El id del pedido no es válido", nil)
	}
	err = api.datos.CancelarPedido(c.Context(), id)
	if errors.Is(err, repository.ErrPedidoNoEncontrado) {
		return responderError(c, fiber.StatusNotFound, "Pedido no encontrado", nil)
	}
	if errors.Is(err, repository.ErrEstadoPedido) {
		return responderError(c, fiber.StatusConflict, "El pedido no se puede cancelar en su estado actual", nil)
	}
	if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudo cancelar el pedido", nil)
	}
	if err := api.cache.InvalidarProductos(c.Context()); err != nil {
		api.logger.Printf("no se pudo invalidar la caché de productos: %v", err)
	}
	if err := api.cache.InvalidarTodasLasRedes(c.Context()); err != nil {
		api.logger.Printf("no se pudo invalidar la caché de redes: %v", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
