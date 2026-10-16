package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"multicatalogo-backend/controllers"
	"multicatalogo-backend/models"
	"multicatalogo-backend/repository"
)

type datosFalsos struct {
	errorPedido   error
	errorRegistro error
	registro      models.RegisterRequest
	comisiones    chan int64
}

func (d *datosFalsos) BuscarUsuarioPorEmail(_ context.Context, email string) (models.Usuario, bool, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.MinCost)
	return models.Usuario{ID: 2, Nombre: "Cliente", Email: email, Password: string(hash), Rol: "cliente", Activo: true, ReferidoID: 2}, true, nil
}

func (d *datosFalsos) RegistrarUsuario(_ context.Context, req models.RegisterRequest, hash string) (models.Usuario, error) {
	d.registro = req
	if d.errorRegistro != nil {
		return models.Usuario{}, d.errorRegistro
	}
	return models.Usuario{ID: 3, Nombre: req.Nombre, Email: req.Email, Password: hash, Rol: "cliente", Activo: true, ReferidoID: 3}, nil
}

func (d *datosFalsos) ObtenerProductos(context.Context) ([]models.Producto, error) {
	return []models.Producto{{ID: 1, Nombre: "Serum", Precio: 45, Stock: 10}}, nil
}

func (d *datosFalsos) ObtenerTodosLosProductos(context.Context) ([]models.Producto, error) {
	return []models.Producto{{ID: 1, Nombre: "Serum", Precio: 45, Stock: 10, Activo: true}}, nil
}

func (d *datosFalsos) ObtenerProductoPorID(_ context.Context, id int) (models.Producto, bool, error) {
	if id != 1 {
		return models.Producto{}, false, nil
	}
	return models.Producto{ID: 1, Nombre: "Serum", Precio: 45, Stock: 10}, true, nil
}

func (d *datosFalsos) CrearProducto(_ context.Context, producto models.Producto) (models.Producto, error) {
	producto.ID = 2
	producto.Activo = true
	return producto, nil
}

func (d *datosFalsos) ActualizarProducto(_ context.Context, id int, producto models.Producto) (models.Producto, error) {
	producto.ID = id
	producto.Activo = true
	return producto, nil
}

func (d *datosFalsos) DesactivarProducto(context.Context, int) error {
	return nil
}

func (d *datosFalsos) ObtenerPedidosUsuario(context.Context, string) ([]models.PedidoResumen, error) {
	return []models.PedidoResumen{}, nil
}

func (d *datosFalsos) ObtenerTodosLosPedidos(context.Context, string, int, int) ([]models.PedidoResumen, int, error) {
	return []models.PedidoResumen{}, 0, nil
}

func (d *datosFalsos) ActualizarEstadoPedido(context.Context, int64, string) error {
	return nil
}

func (d *datosFalsos) CancelarPedido(context.Context, int64) error {
	return nil
}

func (d *datosFalsos) ObtenerComisionesUsuario(context.Context, string) ([]models.Comision, error) {
	return []models.Comision{}, nil
}

func (d *datosFalsos) ObtenerTodasLasComisiones(context.Context, string) ([]models.Comision, error) {
	return []models.Comision{}, nil
}

func (d *datosFalsos) LiquidarComisiones(context.Context, []int64) (models.Liquidacion, error) {
	return models.Liquidacion{ID: 1, Monto: 10}, nil
}

func (d *datosFalsos) ObtenerReporteGeneral(context.Context) (models.ReporteGeneral, error) {
	return models.ReporteGeneral{}, nil
}

func (d *datosFalsos) ObtenerRed(context.Context, string) ([]*models.Referido, error) {
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
func (c *cacheFalsa) InvalidarTodasLasRedes(context.Context) error {
	if c.redInvalidada != nil {
		c.redInvalidada <- "todas"
	}
	return nil
}

func prepararApp(datos *datosFalsos, cache *cacheFalsa) *fiber.App {
	app := fiber.New()
	api := controllers.NewAPI(datos, cache, log.New(io.Discard, "", 0), "clave-de-pruebas-con-mas-de-32-caracteres")
	SetupRoutes(app, api, "clave-de-pruebas-con-mas-de-32-caracteres")
	return app
}

func autorizar(peticion *http.Request) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "2", "email": "cliente@upse.edu.ec", "rol": "cliente", "exp": time.Now().Add(time.Hour).Unix(),
	})
	firmado, _ := token.SignedString([]byte("clave-de-pruebas-con-mas-de-32-caracteres"))
	peticion.Header.Set("Authorization", "Bearer "+firmado)
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
	if datos["email"] != "admin@upse.edu.ec" || datos["rol"] != "cliente" || datos["token"] == "" {
		t.Fatalf("la respuesta cambió el contrato del frontend: %+v", datos)
	}
}

func TestRegistroCreaSesionYConservaCodigoReferido(t *testing.T) {
	datos := &datosFalsos{}
	cuerpo := `{"nombre":"Nueva Persona","email":"nueva@upse.edu.ec","password":"secreto","ref":"ana"}`
	peticion := httptest.NewRequest("POST", "/api/register", bytes.NewBufferString(cuerpo))
	peticion.Header.Set("Content-Type", "application/json")
	respuesta, err := prepararApp(datos, &cacheFalsa{}).Test(peticion)
	if err != nil || respuesta.StatusCode != fiber.StatusCreated {
		t.Fatalf("se esperaba estado 201, respuesta=%v error=%v", respuesta, err)
	}
	if datos.registro.Ref != "ANA" || datos.registro.Email != "nueva@upse.edu.ec" {
		t.Fatalf("el registro no normalizó los datos: %+v", datos.registro)
	}
	var sesion models.AuthResponse
	_ = json.NewDecoder(respuesta.Body).Decode(&sesion)
	if sesion.Token == "" || sesion.Rol != "cliente" {
		t.Fatalf("no se creó una sesión válida: %+v", sesion)
	}
}

func TestRegistroRechazaCorreoDuplicado(t *testing.T) {
	datos := &datosFalsos{errorRegistro: repository.ErrEmailRegistrado}
	cuerpo := `{"nombre":"Cliente","email":"cliente@upse.edu.ec","password":"123456"}`
	peticion := httptest.NewRequest("POST", "/api/register", bytes.NewBufferString(cuerpo))
	peticion.Header.Set("Content-Type", "application/json")
	respuesta, err := prepararApp(datos, &cacheFalsa{}).Test(peticion)
	if err != nil || respuesta.StatusCode != fiber.StatusConflict {
		t.Fatalf("se esperaba estado 409, respuesta=%v error=%v", respuesta, err)
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
	autorizar(peticion)

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
		if email != "todas" {
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
	autorizar(peticion)
	respuesta, err := prepararApp(datos, &cacheFalsa{}).Test(peticion)
	if err != nil || respuesta.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("se esperaba estado 400, respuesta=%v error=%v", respuesta, err)
	}
}

func TestCheckoutRequiereSesion(t *testing.T) {
	peticion := httptest.NewRequest("POST", "/api/checkout", bytes.NewBufferString(`{}`))
	peticion.Header.Set("Content-Type", "application/json")
	respuesta, err := prepararApp(&datosFalsos{}, &cacheFalsa{}).Test(peticion)
	if err != nil || respuesta.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("se esperaba estado 401, respuesta=%v error=%v", respuesta, err)
	}
}
