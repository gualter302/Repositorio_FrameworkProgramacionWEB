package routes

import (
	"github.com/gofiber/fiber/v2"
	"multicatalogo-backend/controllers"
)

// SetupRoutes reúne las direcciones disponibles en la API.
func SetupRoutes(app *fiber.App) {
	app.Get("/api/health", controllers.HealthCheck)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/productos", controllers.GetProductos)
	app.Get("/api/productos/:id", controllers.GetProductoPorID)
	app.Get("/api/red", controllers.GetRed)
}
