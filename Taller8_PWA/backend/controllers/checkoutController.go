package controllers

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"multicatalogo-backend/models"
	"multicatalogo-backend/repository"
)

// Checkout valida la compra y delega el guardado transaccional a PostgreSQL.
func (api *API) Checkout(c *fiber.Ctx) error {
	var req models.CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return responderError(c, fiber.StatusBadRequest, "Cuerpo de petición inválido", nil)
	}
	req.ClaveCompra = strings.TrimSpace(c.Get("Idempotency-Key"))
	if len(req.ClaveCompra) < 16 || len(req.ClaveCompra) > 100 {
		return responderError(c, fiber.StatusBadRequest, "La compra necesita una clave válida para evitar duplicados", nil)
	}

	emailSesion, _ := c.Locals("usuarioEmail").(string)
	req.UsuarioEmail = strings.ToLower(strings.TrimSpace(emailSesion))
	detalles := validarCheckout(req)
	if len(detalles) > 0 {
		return responderError(c, fiber.StatusBadRequest, "La compra contiene datos inválidos", detalles)
	}

	usuario, encontrado, err := api.datos.BuscarUsuarioPorEmail(c.Context(), req.UsuarioEmail)
	if err != nil {
		return responderError(c, fiber.StatusInternalServerError, "No se pudo validar el comprador", nil)
	}
	if !encontrado {
		return responderError(c, fiber.StatusBadRequest, "El usuario de la compra no existe", nil)
	}

	pedido, err := api.datos.CrearPedido(c.Context(), req, usuario)
	if err != nil {
		var stockErr *repository.StockInsuficienteError
		switch {
		case errors.Is(err, repository.ErrClaveCompra):
			return responderError(c, fiber.StatusConflict, "Esta clave ya corresponde a otra compra", nil)
		case errors.Is(err, repository.ErrProductoNoEncontrado):
			return responderError(c, fiber.StatusBadRequest, "Uno de los productos no existe", nil)
		case errors.As(err, &stockErr):
			return responderError(c, fiber.StatusConflict, "No hay existencias suficientes", fiber.Map{
				"productoId": stockErr.ProductoID,
				"disponible": stockErr.Disponible,
				"solicitado": stockErr.Solicitado,
			})
		default:
			return responderError(c, fiber.StatusInternalServerError, "No se pudo completar la compra", nil)
		}
	}

	_ = api.cache.InvalidarProductos(c.Context())
	go api.procesarComisiones(pedido, usuario)

	return c.Status(fiber.StatusCreated).JSON(pedido)
}

func validarCheckout(req models.CheckoutRequest) fiber.Map {
	detalles := fiber.Map{}
	if req.UsuarioEmail == "" {
		detalles["usuarioEmail"] = "El correo del usuario es obligatorio"
	}
	if strings.TrimSpace(req.Cliente.Nombre) == "" || strings.TrimSpace(req.Cliente.Email) == "" ||
		strings.TrimSpace(req.Cliente.Direccion) == "" || strings.TrimSpace(req.Cliente.Ciudad) == "" ||
		strings.TrimSpace(req.Cliente.Telefono) == "" {
		detalles["cliente"] = "Completa todos los datos de entrega"
	}
	if strings.TrimSpace(req.MetodoPago) == "" {
		detalles["metodoPago"] = "Selecciona un método de pago"
	}
	if len(req.Items) == 0 {
		detalles["items"] = "El carrito no puede estar vacío"
	}

	productos := make(map[int]bool)
	for _, item := range req.Items {
		if item.ProductoID <= 0 || item.Cantidad <= 0 {
			detalles["items"] = "Cada producto debe tener un id y una cantidad válidos"
			break
		}
		if productos[item.ProductoID] {
			detalles["items"] = "Cada producto debe aparecer una sola vez"
			break
		}
		productos[item.ProductoID] = true
	}
	return detalles
}

func (api *API) procesarComisiones(pedido models.PedidoConfirmado, usuario models.Usuario) {
	ctx, cancelar := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelar()

	if err := api.datos.CalcularComisiones(ctx, pedido.ID, usuario.ID, pedido.Total); err != nil {
		api.logger.Printf("no se pudieron calcular las comisiones del pedido %s: %v", pedido.Numero, err)
		return
	}

	if err := api.cache.InvalidarTodasLasRedes(ctx); err != nil {
		api.logger.Printf("no se pudo renovar la caché de las redes: %v", err)
	}
}
