package routes

import (
	"github.com/gofiber/fiber/v2"
	"multicatalogo-backend/controllers"
	seguridad "multicatalogo-backend/middleware"
)

// SetupRoutes reúne las direcciones disponibles en la API.
func SetupRoutes(app *fiber.App, api *controllers.API, secret string) {
	app.Get("/api/health", api.HealthCheck)
	app.Post("/api/login", api.Login)
	app.Post("/api/register", api.Register)
	app.Get("/api/productos", api.GetProductos)
	app.Get("/api/productos/:id", api.GetProductoPorID)
	protegida := app.Group("/api", seguridad.RequireAuth(secret))
	protegida.Get("/red", api.GetRed)
	protegida.Post("/checkout", api.Checkout)
	protegida.Get("/pedidos", api.MisPedidos)
	protegida.Get("/comisiones", api.MisComisiones)
	admin := protegida.Group("/admin", seguridad.RequireRole("admin"))
	admin.Get("/productos", api.AdminProductos)
	admin.Post("/productos", api.CrearProducto)
	admin.Put("/productos/:id", api.ActualizarProducto)
	admin.Delete("/productos/:id", api.DesactivarProducto)
	admin.Get("/pedidos", api.AdminPedidos)
	admin.Patch("/pedidos/:id/estado", api.CambiarEstadoPedido)
	admin.Post("/pedidos/:id/cancelar", api.CancelarPedido)
	admin.Get("/comisiones", api.AdminComisiones)
	admin.Post("/liquidaciones", api.LiquidarComisiones)
	admin.Get("/reportes/general", api.ReporteGeneral)
}
