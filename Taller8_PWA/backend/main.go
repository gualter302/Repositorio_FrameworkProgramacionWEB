package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"multicatalogo-backend/config"
	"multicatalogo-backend/controllers"
	"multicatalogo-backend/database"
	"multicatalogo-backend/repository"
	"multicatalogo-backend/routes"
)

func main() {
	config.CargarEnv()

	// El contexto se cancela con Ctrl+C para cerrar el servidor y la base de datos ordenadamente.
	ctx, detener := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer detener()

	pool, err := config.ConectarBaseDeDatos(ctx)
	if err != nil {
		log.Fatalf("no se pudo conectar a PostgreSQL: %v", err)
	}
	defer pool.Close()
	migracionesCtx, cancelarMigraciones := context.WithTimeout(ctx, 2*time.Minute)
	if err := database.AplicarMigraciones(migracionesCtx, pool); err != nil {
		cancelarMigraciones()
		log.Fatalf("no se pudo preparar PostgreSQL: %v", err)
	}
	if err := database.ConfigurarClaveAdministrador(migracionesCtx, pool, os.Getenv("ADMIN_PASSWORD")); err != nil {
		cancelarMigraciones()
		log.Fatalf("no se pudo configurar la cuenta administradora: %v", err)
	}
	cancelarMigraciones()

	clienteRedis, err := config.ConectarRedis(ctx)
	var repositorioCache controllers.CacheRepository
	if err != nil {
		log.Printf("Redis no está disponible; la API continuará consultando PostgreSQL: %v", err)
		repositorioCache = repository.NewNoopCache()
	} else {
		defer clienteRedis.Close()
		repositorioCache = repository.NewRedisRepository(clienteRedis)
	}

	repositorioDatos := repository.NewPostgresRepository(pool)
	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET debe tener al menos 32 caracteres")
	}
	api := controllers.NewAPI(repositorioDatos, repositorioCache, log.Default(), jwtSecret)

	app := fiber.New()

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:5173,http://127.0.0.1:5173"
	}

	// CORS permite que el frontend de Vite se comunique con este servidor.
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, Idempotency-Key",
	}))

	routes.SetupRoutes(app, api, jwtSecret)
	if err := configurarFrontend(app, os.Getenv("STATIC_DIR")); err != nil {
		log.Fatalf("no se pudo preparar el frontend: %v", err)
	}

	puerto := os.Getenv("PORT")
	if puerto == "" {
		puerto = "3000"
	}
	direccion := "0.0.0.0:" + puerto
	erroresServidor := make(chan error, 1)
	go func() {
		log.Printf("Servidor disponible en http://%s", direccion)
		erroresServidor <- app.Listen(direccion)
	}()

	select {
	case <-ctx.Done():
		log.Println("Cerrando servidor, base de datos y caché...")
		if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
			log.Printf("el servidor no pudo cerrar todas las conexiones a tiempo: %v", err)
		}
	case err := <-erroresServidor:
		if err != nil {
			log.Fatalf("el servidor se detuvo: %v", err)
		}
	}
}
