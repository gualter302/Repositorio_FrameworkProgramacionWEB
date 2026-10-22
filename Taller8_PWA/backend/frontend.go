package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"multicatalogo-backend/models"
)

// configurarFrontend sirve la compilación de React desde el mismo dominio que la API.
func configurarFrontend(app *fiber.App, directorio string) error {
	if strings.TrimSpace(directorio) == "" {
		return nil
	}

	absoluto, err := filepath.Abs(directorio)
	if err != nil {
		return fmt.Errorf("ruta del frontend inválida: %w", err)
	}
	indice := filepath.Join(absoluto, "index.html")
	if _, err := os.Stat(indice); err != nil {
		return fmt.Errorf("no se encontró la compilación del frontend en %s: %w", indice, err)
	}

	app.Use(func(c *fiber.Ctx) error {
		ruta := c.Path()
		switch {
		case ruta == "/sw.js", ruta == "/manifest.webmanifest", strings.HasPrefix(ruta, "/workbox-"):
			c.Set(fiber.HeaderCacheControl, "no-cache")
		case strings.HasPrefix(ruta, "/assets/"):
			c.Set(fiber.HeaderCacheControl, "public, max-age=31536000, immutable")
		default:
			c.Set(fiber.HeaderCacheControl, "no-cache")
		}
		return c.Next()
	})

	app.Static("/", absoluto, fiber.Static{Index: "index.html", ByteRange: true})
	app.Use(func(c *fiber.Ctx) error {
		if c.Path() == "/api" || strings.HasPrefix(c.Path(), "/api/") {
			return c.Status(fiber.StatusNotFound).JSON(models.APIError{
				Status: fiber.StatusNotFound, Message: "Ruta de API no encontrada",
			})
		}
		return c.SendFile(indice)
	})

	return nil
}
