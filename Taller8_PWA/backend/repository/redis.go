package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"multicatalogo-backend/models"
)

const (
	claveProductos  = "productos:all"
	tiempoProductos = 5 * time.Minute
	tiempoRed       = 10 * time.Minute
)

// RedisRepository mantiene separada la caché de las consultas de PostgreSQL.
type RedisRepository struct {
	cliente *redis.Client
}

func NewRedisRepository(cliente *redis.Client) *RedisRepository {
	return &RedisRepository{cliente: cliente}
}

func (r *RedisRepository) ObtenerProductos(ctx context.Context) ([]models.Producto, bool) {
	productos := make([]models.Producto, 0)
	return productos, r.cargar(ctx, claveProductos, &productos)
}

func (r *RedisRepository) GuardarProductos(ctx context.Context, productos []models.Producto) {
	r.guardar(ctx, claveProductos, productos, tiempoProductos)
}

func (r *RedisRepository) InvalidarProductos(ctx context.Context) error {
	return r.cliente.Del(ctx, claveProductos).Err()
}

func (r *RedisRepository) ObtenerRed(ctx context.Context, emailUsuario string) ([]*models.Referido, bool) {
	red := make([]*models.Referido, 0)
	return red, r.cargar(ctx, claveRed(emailUsuario), &red)
}

func (r *RedisRepository) GuardarRed(ctx context.Context, emailUsuario string, red []*models.Referido) {
	r.guardar(ctx, claveRed(emailUsuario), red, tiempoRed)
}

func (r *RedisRepository) InvalidarTodasLasRedes(ctx context.Context) error {
	iterador := r.cliente.Scan(ctx, 0, "red:tree:*", 100).Iterator()
	claves := make([]string, 0)
	for iterador.Next(ctx) {
		claves = append(claves, iterador.Val())
	}
	if err := iterador.Err(); err != nil {
		return err
	}
	if len(claves) == 0 {
		return nil
	}
	return r.cliente.Del(ctx, claves...).Err()
}

func claveRed(emailUsuario string) string {
	return fmt.Sprintf("red:tree:%s", emailUsuario)
}

func (r *RedisRepository) cargar(ctx context.Context, clave string, destino any) bool {
	datos, err := r.cliente.Get(ctx, clave).Bytes()
	if err != nil {
		return false
	}
	if err := json.Unmarshal(datos, destino); err != nil {
		_ = r.cliente.Del(ctx, clave).Err()
		return false
	}
	return true
}

func (r *RedisRepository) guardar(ctx context.Context, clave string, valor any, duracion time.Duration) {
	datos, err := json.Marshal(valor)
	if err != nil {
		return
	}
	_ = r.cliente.Set(ctx, clave, datos, duracion).Err()
}
