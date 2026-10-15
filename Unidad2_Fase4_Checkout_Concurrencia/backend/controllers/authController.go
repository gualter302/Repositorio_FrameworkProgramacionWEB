package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"multicatalogo-backend/models"
)

// tokenTemporal se entrega tras un inicio de sesión correcto. En una fase posterior
// se reemplazará por un JWT real; por ahora conserva el valor que el frontend ya usaba.
const tokenTemporal = "fake-jwt-token-123"

// Login valida los datos recibidos antes de consultar las cuentas en PostgreSQL.
func (api *API) Login(c *fiber.Ctx) error {
	var req models.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return responderError(c, fiber.StatusBadRequest, "Cuerpo de petición inválido", nil)
	}

	camposFaltantes := fiber.Map{}
	if strings.TrimSpace(req.Email) == "" {
		camposFaltantes["email"] = "El correo es obligatorio"
	}
	if strings.TrimSpace(req.Password) == "" {
		camposFaltantes["password"] = "La contraseña es obligatoria"
	}
	if len(camposFaltantes) > 0 {
		return responderError(c, fiber.StatusBadRequest, "Faltan campos obligatorios", camposFaltantes)
	}

	usuario, encontrado, err := api.datos.BuscarUsuarioPorCredenciales(c.Context(), strings.ToLower(strings.TrimSpace(req.Email)), req.Password)
	if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudo validar las credenciales", nil)
	}
	if !encontrado {
		return responderError(c, fiber.StatusUnauthorized, "Credenciales incorrectas", nil)
	}

	// La respuesta conserva los campos que el frontend ya utiliza (token y email).
	return c.JSON(fiber.Map{"token": tokenTemporal, "email": usuario.Email})
}
