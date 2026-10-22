package repository

import (
	"context"

	"multicatalogo-backend/models"
)

// NoopCache permite que la API continúe usando PostgreSQL si Redis está temporalmente fuera de servicio.
type NoopCache struct{}

func NewNoopCache() *NoopCache { return &NoopCache{} }

func (*NoopCache) ObtenerProductos(context.Context) ([]models.Producto, bool) { return nil, false }
func (*NoopCache) GuardarProductos(context.Context, []models.Producto)        {}
func (*NoopCache) InvalidarProductos(context.Context) error                   { return nil }
func (*NoopCache) ObtenerRed(context.Context, string) ([]*models.Referido, bool) {
	return nil, false
}
func (*NoopCache) GuardarRed(context.Context, string, []*models.Referido) {}
func (*NoopCache) InvalidarRed(context.Context, string) error             { return nil }
func (*NoopCache) InvalidarTodasLasRedes(context.Context) error           { return nil }
