package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"multicatalogo-backend/models"
)

// pool es la conexión compartida a PostgreSQL. main la inyecta con UsarPool.
var pool *pgxpool.Pool

// UsarPool guarda el pool de conexiones creado al iniciar la aplicación.
func UsarPool(p *pgxpool.Pool) {
	pool = p
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

// ObtenerProductos devuelve todos los artículos del catálogo.
func ObtenerProductos(ctx context.Context) ([]models.Producto, error) {
	filas, err := pool.Query(ctx,
		`SELECT id, nombre, descripcion, precio, categoria, img, galeria
		 FROM productos ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	productos := make([]models.Producto, 0)
	for filas.Next() {
		var p models.Producto
		if err := filas.Scan(&p.ID, &p.Nombre, &p.Descripcion, &p.Precio, &p.Categoria, &p.Img, &p.Galeria); err != nil {
			return nil, err
		}
		productos = append(productos, p)
	}
	return productos, filas.Err()
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
func ObtenerRed(ctx context.Context) ([]*models.Referido, error) {
	filas, err := pool.Query(ctx,
		`SELECT id, nombre, nivel, ventas, parent_id FROM referidos ORDER BY nivel, id`)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	porID := make(map[int]*models.Referido)
	orden := make([]*models.Referido, 0)
	for filas.Next() {
		referido := &models.Referido{Hijos: make([]*models.Referido, 0)}
		if err := filas.Scan(&referido.ID, &referido.Nombre, &referido.Nivel, &referido.Ventas, &referido.ParentID); err != nil {
			return nil, err
		}
		porID[referido.ID] = referido
		orden = append(orden, referido)
	}
	if err := filas.Err(); err != nil {
		return nil, err
	}

	raices := make([]*models.Referido, 0)
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
	return raices, nil
}
