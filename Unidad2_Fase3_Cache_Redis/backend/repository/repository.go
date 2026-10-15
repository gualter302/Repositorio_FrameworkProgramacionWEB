package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"multicatalogo-backend/models"
)

// pool es la conexión compartida a PostgreSQL. main la inyecta con UsarPool.
var pool *pgxpool.Pool
var cache *redis.Client

const (
	claveProductos  = "productos:all"
	tiempoProductos = 5 * time.Minute
	tiempoRed       = 10 * time.Minute
)

// UsarPool guarda el pool de conexiones creado al iniciar la aplicación.
func UsarPool(p *pgxpool.Pool) {
	pool = p
}

// UsarCache guarda la conexión que se usará para responder consultas repetidas.
func UsarCache(cliente *redis.Client) {
	cache = cliente
}

// EstaConectado indica si ya existe un pool disponible. Las pruebas lo usan
// para omitirse cuando no hay base de datos.
func EstaConectado() bool {
	return pool != nil
}

// BuscarUsuarioPorCredenciales consulta la cuenta que coincide con el correo y la contraseña.
func BuscarUsuarioPorCredenciales(ctx context.Context, email, password string) (models.Usuario, bool, error) {
	var usuario models.Usuario
	err := pool.QueryRow(ctx,
		`SELECT id, email, password, rol FROM usuarios WHERE email = $1 AND password = $2`,
		email, password,
	).Scan(&usuario.ID, &usuario.Email, &usuario.Password, &usuario.Rol)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Usuario{}, false, nil
	}
	if err != nil {
		return models.Usuario{}, false, err
	}
	return usuario, true, nil
}

// ObtenerProductos busca primero una copia temporal y consulta PostgreSQL cuando no existe.
func ObtenerProductos(ctx context.Context) ([]models.Producto, bool, error) {
	productos := make([]models.Producto, 0)
	if cargarCache(ctx, claveProductos, &productos) {
		return productos, true, nil
	}

	filas, err := pool.Query(ctx,
		`SELECT id, nombre, descripcion, precio, categoria, img, galeria
		 FROM productos ORDER BY id`)
	if err != nil {
		return nil, false, err
	}
	defer filas.Close()

	for filas.Next() {
		var p models.Producto
		if err := filas.Scan(&p.ID, &p.Nombre, &p.Descripcion, &p.Precio, &p.Categoria, &p.Img, &p.Galeria); err != nil {
			return nil, false, err
		}
		productos = append(productos, p)
	}
	if err := filas.Err(); err != nil {
		return nil, false, err
	}

	guardarCache(ctx, claveProductos, productos, tiempoProductos)
	return productos, false, nil
}

// ObtenerProductoPorID busca un producto por su identificador.
func ObtenerProductoPorID(ctx context.Context, id int) (models.Producto, bool, error) {
	var p models.Producto
	err := pool.QueryRow(ctx,
		`SELECT id, nombre, descripcion, precio, categoria, img, galeria
		 FROM productos WHERE id = $1`, id,
	).Scan(&p.ID, &p.Nombre, &p.Descripcion, &p.Precio, &p.Categoria, &p.Img, &p.Galeria)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Producto{}, false, nil
	}
	if err != nil {
		return models.Producto{}, false, err
	}
	return p, true, nil
}

// ObtenerRed lee la lista plana de referidos y reconstruye el árbol multinivel en memoria.
// Devuelve los nodos raíz (los que no tienen parent_id).
func ObtenerRed(ctx context.Context, emailUsuario string) ([]*models.Referido, bool, error) {
	clave := fmt.Sprintf("red:tree:%s", emailUsuario)
	raices := make([]*models.Referido, 0)
	if cargarCache(ctx, clave, &raices) {
		return raices, true, nil
	}

	filas, err := pool.Query(ctx,
		`SELECT id, nombre, nivel, ventas, parent_id FROM referidos ORDER BY nivel, id`)
	if err != nil {
		return nil, false, err
	}
	defer filas.Close()

	porID := make(map[int]*models.Referido)
	orden := make([]*models.Referido, 0)
	for filas.Next() {
		referido := &models.Referido{Hijos: make([]*models.Referido, 0)}
		if err := filas.Scan(&referido.ID, &referido.Nombre, &referido.Nivel, &referido.Ventas, &referido.ParentID); err != nil {
			return nil, false, err
		}
		porID[referido.ID] = referido
		orden = append(orden, referido)
	}
	if err := filas.Err(); err != nil {
		return nil, false, err
	}

	for _, referido := range orden {
		if referido.ParentID == nil {
			raices = append(raices, referido)
			continue
		}
		if padre, ok := porID[*referido.ParentID]; ok {
			padre.Hijos = append(padre.Hijos, referido)
		} else {
			raices = append(raices, referido) // Padre ausente: lo tratamos como raíz para no perderlo.
		}
	}

	guardarCache(ctx, clave, raices, tiempoRed)
	return raices, false, nil
}

func cargarCache(ctx context.Context, clave string, destino any) bool {
	if cache == nil {
		return false
	}

	datos, err := cache.Get(ctx, clave).Bytes()
	if err != nil {
		return false
	}
	if err := json.Unmarshal(datos, destino); err != nil {
		_ = cache.Del(ctx, clave).Err()
		return false
	}
	return true
}

func guardarCache(ctx context.Context, clave string, valor any, duracion time.Duration) {
	if cache == nil {
		return
	}

	datos, err := json.Marshal(valor)
	if err != nil {
		return
	}
	_ = cache.Set(ctx, clave, datos, duracion).Err()
}
