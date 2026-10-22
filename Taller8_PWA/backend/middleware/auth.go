package middleware

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"multicatalogo-backend/models"
)

type claims struct {
	Email string `json:"email"`
	Rol   string `json:"rol"`
	jwt.RegisteredClaims
}

// RequireAuth comprueba el token y deja la identidad disponible para la ruta.
func RequireAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cabecera := c.Get("Authorization")
		partes := strings.SplitN(cabecera, " ", 2)
		if len(partes) != 2 || !strings.EqualFold(partes[0], "Bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIError{Status: fiber.StatusUnauthorized, Message: "Debes iniciar sesión"})
		}

		datos := &claims{}
		token, err := jwt.ParseWithClaims(partes[1], datos, func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIError{Status: fiber.StatusUnauthorized, Message: "La sesión no es válida o ya terminó"})
		}

		usuarioID, err := strconv.Atoi(datos.Subject)
		if err != nil || usuarioID <= 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIError{Status: fiber.StatusUnauthorized, Message: "La sesión no es válida"})
		}
		c.Locals("usuarioID", usuarioID)
		c.Locals("usuarioEmail", datos.Email)
		c.Locals("usuarioRol", datos.Rol)
		return c.Next()
	}
}

// RequireRole limita una ruta a uno de los roles que ya validó RequireAuth.
func RequireRole(roles ...string) fiber.Handler {
	permitidos := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		permitidos[role] = struct{}{}
	}
	return func(c *fiber.Ctx) error {
		rol, ok := c.Locals("usuarioRol").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(models.APIError{
				Status: fiber.StatusForbidden, Message: "No tienes permisos para esta operación",
			})
		}
		if _, permitido := permitidos[rol]; !permitido {
			return c.Status(fiber.StatusForbidden).JSON(models.APIError{
				Status: fiber.StatusForbidden, Message: "No tienes permisos para esta operación",
			})
		}
		return c.Next()
	}
}
