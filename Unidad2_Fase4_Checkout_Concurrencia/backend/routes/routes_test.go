package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"multicatalogo-backend/controllers"
	"multicatalogo-backend/models"
	"multicatalogo-backend/repository"
)

type datosFalsos struct {
	errorPedido error
	comisiones  chan int64
}

func (d *datosFalsos) BuscarUsuarioPorCredenciales(_ context.Context, email, password string) (models.Usuario, bool, error) {
	if email == "admin@upse.edu.ec" && password == "123456" {
		return models.Usuario{ID: 1, Email: email, Rol: "admin"}, true, nil
	}
	return models.Usuario{}, false, nil
}

func (d *datosFalsos) BuscarUsuarioPorEmail(_ context.Context, email string) (models.Usuario, bool, error) {
	return models.Usuario{ID: 2, Email: email, Rol: "cliente"}, true, nil
}

func (d *datosFalsos) ObtenerProductos(context.Context) ([]models.Producto, error) {
	return []models.Producto{{ID: 1, Nombre: "Serum", Precio: 45, Stock: 10}}, nil
}

func (d *datosFalsos) ObtenerProductoPorID(_ context.Context, id int) (models.Producto, bool, error) {
	if id != 1 {
		return models.Producto{}, false, nil
	}
	return models.Producto{ID: 1, Nombre: "Serum", Precio: 45, Stock: 10}, true, nil
}

func (d *datosFalsos) ObtenerRed(context.Context) ([]*models.Referido, error) {
	return []*models.Referido{{ID: 1, Nombre: "Administrador", Hijos: []*models.Referido{}}}, nil
}

func (d *datosFalsos) CrearPedido(_ context.Context, req models.CheckoutRequest, _ models.Usuario) (models.PedidoConfirmado, error) {
	if d.errorPedido != nil {
		return models.PedidoConfirmado{}, d.errorPedido
	}
	return models.PedidoConfirmado{
		ID: 1, Numero: "MC-000001", Total: 45, Estado: "confirmado", Fecha: time.Now(),
		Cliente: req.Cliente,
		Items:   []models.PedidoItemConfirmado{{ProductoID: 1, Nombre: "Serum", PrecioUnitario: 45, Cantidad: 1, Subtotal: 45}},
	}, nil
}

func (d *datosFalsos) CalcularComisiones(_ context.Context, pedidoID int64, _ int, _ float64) error {
	if d.comisiones != nil {
		d.comisiones <- pedidoID
	}
	return nil
}

type cacheFalsa struct {
	redInvalidada chan string
}

func (*cacheFalsa) ObtenerProductos(context.Context) ([]models.Producto, bool) { return nil, false }
func (*cacheFalsa) GuardarProductos(context.Context, []models.Producto)        {}
func (*cacheFalsa) InvalidarProductos(context.Context) error                   { return nil }
func (*cacheFalsa) ObtenerRed(context.Context, string) ([]*models.Referido, bool) {
	return nil, false
}
func (*cacheFalsa) GuardarRed(context.Context, string, []*models.Referido) {}
func (c *cacheFalsa) InvalidarRed(_ context.Context, email string) error {
	if c.redInvalidada != nil {
		c.redInvalidada <- email
	}
	return nil
}

func prepararApp(datos *datosFalsos, cache *cacheFalsa) *fiber.App {
	app := fiber.New()
	api := controllers.NewAPI(datos, cache, log.New(io.Discard, "", 0))
	SetupRoutes(app, api)
	return app
}

func TestHealthCheck(t *testing.T) {
	respuesta, err := prepararApp(&datosFalsos{}, &cacheFalsa{}).Test(httptest.NewRequest("GET", "/api/health", nil))
	if err != nil || respuesta.StatusCode != fiber.StatusOK {
		t.Fatalf("se esperaba estado 200, respuesta=%v error=%v", respuesta, err)
	}
}

func TestLoginRechazaCamposVacios(t *testing.T) {
	peticion := httptest.NewRequest("POST", "/api/login", bytes.NewBufferString(`{"email":"","password":""}`))
	peticion.Header.Set("Content-Type", "application/json")
	respuesta, err := prepararApp(&datosFalsos{}, &cacheFalsa{}).Test(peticion)
	if err != nil || respuesta.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("se esperaba estado 400, respuesta=%v error=%v", respuesta, err)
	}

	var apiError models.APIError
	if err := json.NewDecoder(respuesta.Body).Decode(&apiError); err != nil || apiError.Message != "Faltan campos obligatorios" {
		t.Fatalf("el error recibido no es el esperado: %+v, error=%v", apiError, err)
	}
}

func TestLoginConservaRespuestaDelFrontend(t *testing.T) {
	peticion := httptest.NewRequest("POST", "/api/login", bytes.NewBufferString(`{"email":"admin@upse.edu.ec","password":"123456"}`))
	peticion.Header.Set("Content-Type", "application/json")
	respuesta, err := prepararApp(&datosFalsos{}, &cacheFalsa{}).Test(peticion)
	if err != nil || respuesta.StatusCode != fiber.StatusOK {
		t.Fatalf("se esperaba estado 200, respuesta=%v error=%v", respuesta, err)
	}
	var datos map[string]string
	_ = json.NewDecoder(respuesta.Body).Decode(&datos)
	if datos["email"] != "admin@upse.edu.ec" || datos["token"] == "" {
		t.Fatalf("la respuesta cambió el contrato del frontend: %+v", datos)
	}
}

func TestProductoRechazaIDNoNumerico(t *testing.T) {
	respuesta, err := prepararApp(&datosFalsos{}, &cacheFalsa{}).Test(httptest.NewRequest("GET", "/api/productos/abc", nil))
	if err != nil || respuesta.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("se esperaba estado 400, respuesta=%v error=%v", respuesta, err)
	}
}

func TestProductoPorID(t *testing.T) {
	respuesta, err := prepararApp(&datosFalsos{}, &cacheFalsa{}).Test(httptest.NewRequest("GET", "/api/productos/1", nil))
	if err != nil || respuesta.StatusCode != fiber.StatusOK {
		t.Fatalf("se esperaba estado 200, respuesta=%v error=%v", respuesta, err)
	}
}

func TestCheckoutEjecutaComisionesDespuesDeCrearPedido(t *testing.T) {
	comisiones := make(chan int64, 1)
	redInvalidada := make(chan string, 2)
	datos := &datosFalsos{comisiones: comisiones}
	cache := &cacheFalsa{redInvalidada: redInvalidada}
	cuerpo := `{"usuarioEmail":"cliente@upse.edu.ec","cliente":{"nombre":"Cliente","email":"cliente@upse.edu.ec","direccion":"Centro","ciudad":"La Libertad","telefono":"0999999999"},"metodoPago":"Efectivo","items":[{"productoId":1,"cantidad":1}]}`
	peticion := httptest.NewRequest("POST", "/api/checkout", bytes.NewBufferString(cuerpo))
	peticion.Header.Set("Content-Type", "application/json")

	respuesta, err := prepararApp(datos, cache).Test(peticion)
	if err != nil || respuesta.StatusCode != fiber.StatusCreated {
		t.Fatalf("se esperaba estado 201, respuesta=%v error=%v", respuesta, err)
	}
	select {
	case pedidoID := <-comisiones:
		if pedidoID != 1 {
			t.Fatalf("se procesó un pedido inesperado: %d", pedidoID)
		}
	case <-time.After(time.Second):
		t.Fatal("las comisiones no se ejecutaron en segundo plano")
	}
	select {
	case email := <-redInvalidada:
		if email != "cliente@upse.edu.ec" {
			t.Fatalf("se invalidó una red inesperada: %s", email)
		}
	case <-time.After(time.Second):
		t.Fatal("no se invalidó la caché de la red")
	}
}

func TestCheckoutReportaProductoInexistente(t *testing.T) {
	datos := &datosFalsos{errorPedido: fmt.Errorf("%w: 99999", repository.ErrProductoNoEncontrado)}
	cuerpo := `{"usuarioEmail":"cliente@upse.edu.ec","cliente":{"nombre":"Cliente","email":"cliente@upse.edu.ec","direccion":"Centro","ciudad":"La Libertad","telefono":"0999999999"},"metodoPago":"Efectivo","items":[{"productoId":99999,"cantidad":1}]}`
	peticion := httptest.NewRequest("POST", "/api/checkout", bytes.NewBufferString(cuerpo))
	peticion.Header.Set("Content-Type", "application/json")
	respuesta, err := prepararApp(datos, &cacheFalsa{}).Test(peticion)
	if err != nil || respuesta.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("se esperaba estado 400, respuesta=%v error=%v", respuesta, err)
	}
}
