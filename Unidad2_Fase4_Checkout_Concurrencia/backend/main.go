package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"multicatalogo-backend/config"
	"multicatalogo-backend/controllers"
	"multicatalogo-backend/repository"
	"multicatalogo-backend/routes"
)

func main() {
	config.CargarEnv()

	// El contexto se cancela con Ctrl+C para cerrar el servidor y la base de datos ordenadamente.
	ctx, detener := signal.NotifyContext(context.Background(), os.Interrupt)
	defer detener()

	pool, err := config.ConectarBaseDeDatos(ctx)
	if err != nil {
		log.Fatalf("no se pudo conectar a PostgreSQL: %v", err)
	}
	defer pool.Close()

	clienteRedis, err := config.ConectarRedis(ctx)
	if err != nil {
		log.Fatalf("no se pudo conectar a Redis: %v", err)
	}
	defer clienteRedis.Close()

	repositorioDatos := repository.NewPostgresRepository(pool)
	repositorioCache := repository.NewRedisRepository(clienteRedis)
	api := controllers.NewAPI(repositorioDatos, repositorioCache, log.Default())

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

	routes.SetupRoutes(app, api)

	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Fatalf("el servidor se detuvo: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Cerrando servidor, base de datos y caché...")
	_ = app.Shutdown()
}
