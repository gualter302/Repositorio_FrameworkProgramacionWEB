package routes

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"multicatalogo-backend/models"
	"multicatalogo-backend/repository"
)

// omitirSinBaseDeDatos salta las pruebas que necesitan PostgreSQL cuando no hay conexión.
func omitirSinBaseDeDatos(t *testing.T) {
	if !repository.EstaConectado() {
		t.Skip("se omite: requiere una conexión a PostgreSQL")
	}
}

func prepararApp() *fiber.App {
	app := fiber.New()
	SetupRoutes(app)
	return app
}

func TestHealthCheck(t *testing.T) {
	respuesta, err := prepararApp().Test(httptest.NewRequest("GET", "/api/health", nil))
	if err != nil {
		t.Fatalf("no se pudo consultar el estado del servidor: %v", err)
	}
	if respuesta.StatusCode != fiber.StatusOK {
		t.Fatalf("se esperaba estado 200 y llegó %d", respuesta.StatusCode)
	}
}

func TestLoginRechazaCamposVacios(t *testing.T) {
	cuerpo := bytes.NewBufferString(`{"email":"","password":""}`)
	peticion := httptest.NewRequest("POST", "/api/login", cuerpo)
	peticion.Header.Set("Content-Type", "application/json")

	respuesta, err := prepararApp().Test(peticion)
	if err != nil {
		t.Fatalf("no se pudo probar el inicio de sesión: %v", err)
	}
	if respuesta.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("se esperaba estado 400 y llegó %d", respuesta.StatusCode)
	}

	var apiError models.APIError
	if err := json.NewDecoder(respuesta.Body).Decode(&apiError); err != nil {
		t.Fatalf("la respuesta no tiene el formato esperado: %v", err)
	}
	if apiError.Status != fiber.StatusBadRequest || apiError.Message != "Faltan campos obligatorios" {
		t.Fatalf("el error recibido está incompleto: %+v", apiError)
	}
}

func TestLoginConservaRespuestaDelFrontend(t *testing.T) {
	omitirSinBaseDeDatos(t)
	cuerpo := bytes.NewBufferString(`{"email":"admin@upse.edu.ec","password":"123456"}`)
	peticion := httptest.NewRequest("POST", "/api/login", cuerpo)
	peticion.Header.Set("Content-Type", "application/json")

	respuesta, err := prepararApp().Test(peticion)
	if err != nil {
		t.Fatalf("no se pudo probar el inicio de sesión: %v", err)
	}
	if respuesta.StatusCode != fiber.StatusOK {
		t.Fatalf("se esperaba estado 200 y llegó %d", respuesta.StatusCode)
	}

	var datos map[string]string
	if err := json.NewDecoder(respuesta.Body).Decode(&datos); err != nil {
		t.Fatalf("la respuesta no tiene el formato esperado: %v", err)
	}
	if datos["email"] != "admin@upse.edu.ec" || datos["token"] == "" {
		t.Fatalf("la respuesta cambió el contrato del frontend: %+v", datos)
	}
}

func TestProductoRechazaIDNoNumerico(t *testing.T) {
	respuesta, err := prepararApp().Test(httptest.NewRequest("GET", "/api/productos/abc", nil))
	if err != nil {
		t.Fatalf("no se pudo consultar el producto: %v", err)
	}
	if respuesta.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("se esperaba estado 400 y llegó %d", respuesta.StatusCode)
	}
}

func TestProductoPorID(t *testing.T) {
	omitirSinBaseDeDatos(t)
	respuesta, err := prepararApp().Test(httptest.NewRequest("GET", "/api/productos/1", nil))
	if err != nil {
		t.Fatalf("no se pudo consultar el producto: %v", err)
	}
	if respuesta.StatusCode != fiber.StatusOK {
		t.Fatalf("se esperaba estado 200 y llegó %d", respuesta.StatusCode)
	}
}
