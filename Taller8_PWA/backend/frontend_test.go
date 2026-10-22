package main

import (
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestConfigurarFrontendSirveRutasYArchivos(t *testing.T) {
	directorio := filepath.Join("testdata", "frontend")

	app := fiber.New()
	defer func() { _ = app.Shutdown() }()
	if err := configurarFrontend(app, directorio); err != nil {
		t.Fatal(err)
	}

	respuesta, err := app.Test(httptest.NewRequest("GET", "/catalogo", nil))
	if err != nil || respuesta.StatusCode != fiber.StatusOK {
		t.Fatalf("la ruta de React no respondió correctamente: estado=%d error=%v", respuesta.StatusCode, err)
	}
	if cache := respuesta.Header.Get(fiber.HeaderCacheControl); cache != "no-cache" {
		t.Fatalf("se esperaba no-cache para la interfaz, se recibió %q", cache)
	}
	_ = respuesta.Body.Close()

	respuesta, err = app.Test(httptest.NewRequest("GET", "/assets/app.js", nil))
	if err != nil || respuesta.StatusCode != fiber.StatusOK {
		t.Fatalf("el archivo versionado no respondió correctamente: estado=%d error=%v", respuesta.StatusCode, err)
	}
	if cache := respuesta.Header.Get(fiber.HeaderCacheControl); !strings.Contains(cache, "immutable") {
		t.Fatalf("se esperaba caché inmutable para assets, se recibió %q", cache)
	}
	_ = respuesta.Body.Close()

	respuesta, err = app.Test(httptest.NewRequest("GET", "/api/no-existe", nil))
	if err != nil || respuesta.StatusCode != fiber.StatusNotFound {
		t.Fatalf("una ruta de API desconocida debe responder 404: estado=%d error=%v", respuesta.StatusCode, err)
	}
	_ = respuesta.Body.Close()
}
