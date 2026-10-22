package controllers

import "github.com/gofiber/fiber/v2"

func (api *API) ReporteGeneral(c *fiber.Ctx) error {
	reporte, err := api.datos.ObtenerReporteGeneral(c.Context())
	if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudo generar el reporte", nil)
	}
	return c.JSON(reporte)
}
