package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"multicatalogo-backend/models"
	"multicatalogo-backend/repository"
)

// Login valida los datos recibidos antes de consultar las cuentas guardadas.
func Login(c *fiber.Ctx) error {
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

	usuario, encontrado := repository.BuscarUsuarioPorCredenciales(strings.TrimSpace(req.Email), req.Password)
	if !encontrado {
		return responderError(c, fiber.StatusUnauthorized, "Credenciales incorrectas", nil)
	}

	// La respuesta conserva los campos que el frontend ya utiliza.
	return c.JSON(fiber.Map{"token": usuario.Token, "email": usuario.Email})
}
