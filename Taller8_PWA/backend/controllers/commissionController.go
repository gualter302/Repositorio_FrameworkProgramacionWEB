package controllers

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"multicatalogo-backend/repository"
)

func (api *API) MisComisiones(c *fiber.Ctx) error {
	email, _ := c.Locals("usuarioEmail").(string)
	comisiones, err := api.datos.ObtenerComisionesUsuario(c.Context(), email)
	if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudieron obtener tus comisiones", nil)
	}
	return c.JSON(comisiones)
}

func (api *API) AdminComisiones(c *fiber.Ctx) error {
	estado := strings.ToLower(strings.TrimSpace(c.Query("estado")))
	if estado != "" && estado != "pendiente" && estado != "pagada" {
		return responderError(c, fiber.StatusBadRequest, "El estado de comisión no es válido", nil)
	}
	comisiones, err := api.datos.ObtenerTodasLasComisiones(c.Context(), estado)
	if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudieron obtener las comisiones", nil)
	}
	return c.JSON(comisiones)
}

func (api *API) LiquidarComisiones(c *fiber.Ctx) error {
	var cuerpo struct {
		Ids []int64 `json:"comisionIds"`
	}
	if err := c.BodyParser(&cuerpo); err != nil || len(cuerpo.Ids) == 0 {
		return responderError(c, fiber.StatusBadRequest, "Debes seleccionar al menos una comisión", nil)
	}
	for _, id := range cuerpo.Ids {
		if id <= 0 {
			return responderError(c, fiber.StatusBadRequest, "Los identificadores de comisión no son válidos", nil)
		}
	}
	liquidacion, err := api.datos.LiquidarComisiones(c.Context(), cuerpo.Ids)
	if errors.Is(err, repository.ErrComisionNoLiquidable) {
		return responderError(c, fiber.StatusConflict, "Una o más comisiones ya fueron liquidadas o no existen", nil)
	}
	if errors.Is(err, repository.ErrLiquidacionVacia) {
		return responderError(c, fiber.StatusBadRequest, "Debes seleccionar al menos una comisión", nil)
	}
	if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudo crear la liquidación", nil)
	}
	return c.Status(fiber.StatusCreated).JSON(liquidacion)
}
