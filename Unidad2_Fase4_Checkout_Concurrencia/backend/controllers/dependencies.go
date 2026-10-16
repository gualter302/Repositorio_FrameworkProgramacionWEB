package controllers

import (
	"context"
	"log"

	"multicatalogo-backend/models"
)

// DataRepository define lo que los controladores necesitan de PostgreSQL.
type DataRepository interface {
	BuscarUsuarioPorCredenciales(context.Context, string, string) (models.Usuario, bool, error)
	BuscarUsuarioPorEmail(context.Context, string) (models.Usuario, bool, error)
	ObtenerProductos(context.Context) ([]models.Producto, error)
	ObtenerProductoPorID(context.Context, int) (models.Producto, bool, error)
	ObtenerRed(context.Context) ([]*models.Referido, error)
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
}

// API recibe sus dependencias al iniciar la aplicación y no usa variables globales.
type API struct {
	datos  DataRepository
	cache  CacheRepository
	logger *log.Logger
}

func NewAPI(datos DataRepository, cache CacheRepository, logger *log.Logger) *API {
	return &API{datos: datos, cache: cache, logger: logger}
}
