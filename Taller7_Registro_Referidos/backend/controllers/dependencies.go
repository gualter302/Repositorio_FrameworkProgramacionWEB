package controllers

import (
	"context"
	"log"

	"multicatalogo-backend/models"
)

// DataRepository define lo que los controladores necesitan de PostgreSQL.
type DataRepository interface {
	BuscarUsuarioPorEmail(context.Context, string) (models.Usuario, bool, error)
	RegistrarUsuario(context.Context, models.RegisterRequest, string) (models.Usuario, error)
	ObtenerProductos(context.Context) ([]models.Producto, error)
	ObtenerTodosLosProductos(context.Context) ([]models.Producto, error)
	ObtenerProductoPorID(context.Context, int) (models.Producto, bool, error)
	CrearProducto(context.Context, models.Producto) (models.Producto, error)
	ActualizarProducto(context.Context, int, models.Producto) (models.Producto, error)
	DesactivarProducto(context.Context, int) error
	ObtenerPedidosUsuario(context.Context, string) ([]models.PedidoResumen, error)
	ObtenerTodosLosPedidos(context.Context, string, int, int) ([]models.PedidoResumen, int, error)
	ActualizarEstadoPedido(context.Context, int64, string) error
	CancelarPedido(context.Context, int64) error
	ObtenerComisionesUsuario(context.Context, string) ([]models.Comision, error)
	ObtenerTodasLasComisiones(context.Context, string) ([]models.Comision, error)
	LiquidarComisiones(context.Context, []int64) (models.Liquidacion, error)
	ObtenerReporteGeneral(context.Context) (models.ReporteGeneral, error)
	ObtenerRed(context.Context, string) ([]*models.Referido, error)
	CrearPedido(context.Context, models.CheckoutRequest, models.Usuario) (models.PedidoConfirmado, error)
	CalcularComisiones(context.Context, int64, int, float64) error
}

// CacheRepository permite sustituir Redis por una implementación falsa en las pruebas.
type CacheRepository interface {
	ObtenerProductos(context.Context) ([]models.Producto, bool)
	GuardarProductos(context.Context, []models.Producto)
	InvalidarProductos(context.Context) error
	ObtenerRed(context.Context, string) ([]*models.Referido, bool)
	GuardarRed(context.Context, string, []*models.Referido)
	InvalidarRed(context.Context, string) error
	InvalidarTodasLasRedes(context.Context) error
}

// API recibe sus dependencias al iniciar la aplicación y no usa variables globales.
type API struct {
	datos  DataRepository
	cache  CacheRepository
	logger *log.Logger
	secret string
}

func NewAPI(datos DataRepository, cache CacheRepository, logger *log.Logger, secret string) *API {
	return &API{datos: datos, cache: cache, logger: logger, secret: secret}
}
