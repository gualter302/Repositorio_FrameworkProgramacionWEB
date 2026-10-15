package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"multicatalogo-backend/routes"
)

func main() {
	app := fiber.New()

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:5173,http://127.0.0.1:5173"
	}

	// CORS permite que el frontend de Vite se comunique con este servidor.
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
