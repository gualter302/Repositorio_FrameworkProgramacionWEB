package routes

import (
	"github.com/gofiber/fiber/v2"
	"multicatalogo-backend/controllers"
)

// SetupRoutes reúne las direcciones disponibles en la API.
func SetupRoutes(app *fiber.App, api *controllers.API) {
	app.Get("/api/health", api.HealthCheck)
	app.Post("/api/login", api.Login)
	app.Get("/api/productos", api.GetProductos)
	app.Get("/api/productos/:id", api.GetProductoPorID)
	app.Get("/api/red", api.GetRed)
	app.Post("/api/checkout", api.Checkout)
}
