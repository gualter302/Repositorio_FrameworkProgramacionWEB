package controllers

import (
	"errors"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"multicatalogo-backend/models"
	"multicatalogo-backend/repository"
)

type sessionClaims struct {
	Email string `json:"email"`
	Rol   string `json:"rol"`
	jwt.RegisteredClaims
}

// Login valida la contraseña protegida y crea una sesión con duración limitada.
func (api *API) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return responderError(c, fiber.StatusBadRequest, "Cuerpo de petición inválido", nil)
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || strings.TrimSpace(req.Password) == "" {
		return responderError(c, fiber.StatusBadRequest, "Faltan campos obligatorios", nil)
	}

	usuario, encontrado, err := api.datos.BuscarUsuarioPorEmail(c.Context(), email)
	if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudo validar las credenciales", nil)
	}
	if !encontrado || bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(req.Password)) != nil {
		return responderError(c, fiber.StatusUnauthorized, "Credenciales incorrectas", nil)
	}
	if !usuario.Activo {
		return responderError(c, fiber.StatusForbidden, "La cuenta está desactivada", nil)
	}

	return api.responderSesion(c, fiber.StatusOK, usuario)
}

// Register crea la cuenta y su nodo de referido dentro de una sola transacción.
func (api *API) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return responderError(c, fiber.StatusBadRequest, "Cuerpo de petición inválido", nil)
	}

	req.Nombre = strings.TrimSpace(req.Nombre)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Ref = strings.ToUpper(strings.TrimSpace(req.Ref))
	detalles := fiber.Map{}
	if req.Nombre == "" {
		detalles["nombre"] = "El nombre es obligatorio"
	}
	direccion, err := mail.ParseAddress(req.Email)
	if err != nil || direccion.Address != req.Email {
		detalles["email"] = "Escribe un correo válido"
	}
	if len(req.Password) < 6 {
		detalles["password"] = "La contraseña debe tener al menos 6 caracteres"
	}
	if len(detalles) > 0 {
		return responderError(c, fiber.StatusBadRequest, "Revisa los datos del registro", detalles)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudo proteger la contraseña", nil)
	}
	usuario, err := api.datos.RegistrarUsuario(c.Context(), req, string(hash))
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEmailRegistrado):
			return responderError(c, fiber.StatusConflict, "El correo ya está registrado", nil)
		case errors.Is(err, repository.ErrCodigoReferido):
			return responderError(c, fiber.StatusBadRequest, "El código de referido no existe", nil)
		default:
			return responderError(c, fiber.StatusInternalServerError, "No se pudo crear la cuenta", nil)
		}
	}
	_ = api.cache.InvalidarTodasLasRedes(c.Context())

	return api.responderSesion(c, fiber.StatusCreated, usuario)
}

func (api *API) responderSesion(c *fiber.Ctx, estado int, usuario models.Usuario) error {
	ahora := time.Now()
	claims := sessionClaims{
		Email: usuario.Email,
		Rol:   usuario.Rol,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(usuario.ID),
			IssuedAt:  jwt.NewNumericDate(ahora),
			ExpiresAt: jwt.NewNumericDate(ahora.Add(2 * time.Hour)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(api.secret))
	if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudo iniciar la sesión", nil)
	}
	return c.Status(estado).JSON(models.AuthResponse{
		Token: token, Email: usuario.Email, Nombre: usuario.Nombre, Rol: usuario.Rol,
	})
}
