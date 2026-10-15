package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"multicatalogo-backend/models"
)

var (
	ErrUsuarioNoEncontrado  = errors.New("usuario no encontrado")
	ErrProductoNoEncontrado = errors.New("producto no encontrado")
)

// StockInsuficienteError permite explicar qué producto no puede completar la compra.
type StockInsuficienteError struct {
	ProductoID int
	Disponible int
	Solicitado int
}

func (e *StockInsuficienteError) Error() string {
	return fmt.Sprintf("stock insuficiente para el producto %d", e.ProductoID)
}

// PostgresRepository concentra todas las operaciones que dependen de PostgreSQL.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) BuscarUsuarioPorCredenciales(ctx context.Context, email, password string) (models.Usuario, bool, error) {
	var usuario models.Usuario
	err := r.pool.QueryRow(ctx,
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

func (r *PostgresRepository) BuscarUsuarioPorEmail(ctx context.Context, email string) (models.Usuario, bool, error) {
	var usuario models.Usuario
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password, rol FROM usuarios WHERE email = $1`, email,
	).Scan(&usuario.ID, &usuario.Email, &usuario.Password, &usuario.Rol)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Usuario{}, false, nil
	}
	if err != nil {
		return models.Usuario{}, false, err
	}
	return usuario, true, nil
}

func (r *PostgresRepository) ObtenerProductos(ctx context.Context) ([]models.Producto, error) {
	filas, err := r.pool.Query(ctx,
		`SELECT id, nombre, descripcion, precio, categoria, img, galeria, stock
		 FROM productos ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	productos := make([]models.Producto, 0)
	for filas.Next() {
		var producto models.Producto
		if err := filas.Scan(
			&producto.ID,
			&producto.Nombre,
			&producto.Descripcion,
			&producto.Precio,
			&producto.Categoria,
			&producto.Img,
			&producto.Galeria,
			&producto.Stock,
		); err != nil {
			return nil, err
		}
		productos = append(productos, producto)
	}
	return productos, filas.Err()
}

func (r *PostgresRepository) ObtenerProductoPorID(ctx context.Context, id int) (models.Producto, bool, error) {
	var producto models.Producto
	err := r.pool.QueryRow(ctx,
		`SELECT id, nombre, descripcion, precio, categoria, img, galeria, stock
		 FROM productos WHERE id = $1`, id,
	).Scan(
		&producto.ID,
		&producto.Nombre,
		&producto.Descripcion,
		&producto.Precio,
		&producto.Categoria,
		&producto.Img,
		&producto.Galeria,
		&producto.Stock,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Producto{}, false, nil
	}
	if err != nil {
		return models.Producto{}, false, err
	}
	return producto, true, nil
}

func (r *PostgresRepository) ObtenerRed(ctx context.Context) ([]*models.Referido, error) {
	filas, err := r.pool.Query(ctx,
		`SELECT id, nombre, nivel, ventas, comisiones, parent_id
		 FROM referidos ORDER BY nivel, id`)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	porID := make(map[int]*models.Referido)
	orden := make([]*models.Referido, 0)
	for filas.Next() {
		referido := &models.Referido{Hijos: make([]*models.Referido, 0)}
		if err := filas.Scan(
			&referido.ID,
			&referido.Nombre,
			&referido.Nivel,
			&referido.Ventas,
			&referido.Comisiones,
			&referido.ParentID,
		); err != nil {
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
			raices = append(raices, referido)
		}
	}
	return raices, nil
}

// CrearPedido guarda cabecera, detalles y descuento de stock dentro de una sola transacción.
func (r *PostgresRepository) CrearPedido(ctx context.Context, req models.CheckoutRequest, usuario models.Usuario) (models.PedidoConfirmado, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return models.PedidoConfirmado{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	pedido := models.PedidoConfirmado{
		Estado:  "confirmado",
		Cliente: req.Cliente,
		Items:   make([]models.PedidoItemConfirmado, 0, len(req.Items)),
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO pedidos (
			usuario_id, cliente_nombre, cliente_email, direccion, ciudad, telefono,
			metodo_pago, estado, total
		) VALUES ($1, $2, $3, $4, $5, $6, $7, 'procesando', 0)
		RETURNING id, fecha`,
		usuario.ID,
		req.Cliente.Nombre,
		req.Cliente.Email,
		req.Cliente.Direccion,
		req.Cliente.Ciudad,
		req.Cliente.Telefono,
		req.MetodoPago,
	).Scan(&pedido.ID, &pedido.Fecha)
	if err != nil {
		return models.PedidoConfirmado{}, err
	}

	for _, itemSolicitado := range req.Items {
		var nombre string
		var precio float64
		var stock int
		err := tx.QueryRow(ctx,
			`SELECT nombre, precio, stock FROM productos WHERE id = $1 FOR UPDATE`,
			itemSolicitado.ProductoID,
		).Scan(&nombre, &precio, &stock)
		if errors.Is(err, pgx.ErrNoRows) {
			return models.PedidoConfirmado{}, fmt.Errorf("%w: %d", ErrProductoNoEncontrado, itemSolicitado.ProductoID)
		}
		if err != nil {
			return models.PedidoConfirmado{}, err
		}
		if stock < itemSolicitado.Cantidad {
			return models.PedidoConfirmado{}, &StockInsuficienteError{
				ProductoID: itemSolicitado.ProductoID,
				Disponible: stock,
				Solicitado: itemSolicitado.Cantidad,
			}
		}

		subtotal := precio * float64(itemSolicitado.Cantidad)
		if _, err := tx.Exec(ctx,
			`UPDATE productos SET stock = stock - $2 WHERE id = $1`,
			itemSolicitado.ProductoID, itemSolicitado.Cantidad,
		); err != nil {
			return models.PedidoConfirmado{}, err
		}
		if _, err := tx.Exec(ctx,
			`INSERT INTO pedido_items (
				pedido_id, producto_id, nombre, precio_unitario, cantidad, subtotal
			) VALUES ($1, $2, $3, $4, $5, $6)`,
			pedido.ID,
			itemSolicitado.ProductoID,
			nombre,
			precio,
			itemSolicitado.Cantidad,
			subtotal,
		); err != nil {
			return models.PedidoConfirmado{}, err
		}

		pedido.Total += subtotal
		pedido.Items = append(pedido.Items, models.PedidoItemConfirmado{
			ProductoID:     itemSolicitado.ProductoID,
			Nombre:         nombre,
			Cantidad:       itemSolicitado.Cantidad,
			PrecioUnitario: precio,
			Subtotal:       subtotal,
		})
	}

	pedido.Numero = fmt.Sprintf("MC-%06d", pedido.ID)
	if _, err := tx.Exec(ctx,
		`UPDATE pedidos SET numero = $1, total = $2, estado = 'confirmado' WHERE id = $3`,
		pedido.Numero, pedido.Total, pedido.ID,
	); err != nil {
		return models.PedidoConfirmado{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.PedidoConfirmado{}, err
	}
	return pedido, nil
}

// CalcularComisiones reparte la compra entre los tres ascendentes del comprador.
func (r *PostgresRepository) CalcularComisiones(ctx context.Context, pedidoID int64, usuarioID int, total float64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	filas, err := tx.Query(ctx, `
		WITH RECURSIVE ascendentes AS (
			SELECT padre.id, padre.parent_id, 1 AS profundidad
			FROM usuarios usuario
			JOIN referidos actual ON actual.id = usuario.referido_id
			JOIN referidos padre ON padre.id = actual.parent_id
			WHERE usuario.id = $1

			UNION ALL

			SELECT padre.id, padre.parent_id, ascendentes.profundidad + 1
			FROM ascendentes
			JOIN referidos padre ON padre.id = ascendentes.parent_id
			WHERE ascendentes.profundidad < 3
		)
		SELECT id, profundidad FROM ascendentes ORDER BY profundidad`, usuarioID)
	if err != nil {
		return err
	}

	type ascendente struct {
		id          int
		profundidad int
	}
	ascendentes := make([]ascendente, 0, 3)
	for filas.Next() {
		var item ascendente
		if err := filas.Scan(&item.id, &item.profundidad); err != nil {
			filas.Close()
			return err
		}
		ascendentes = append(ascendentes, item)
	}
	if err := filas.Err(); err != nil {
		filas.Close()
		return err
	}
	filas.Close()

	porcentajes := map[int]float64{1: 0.05, 2: 0.03, 3: 0.02}
	for _, referido := range ascendentes {
		porcentaje := porcentajes[referido.profundidad]
		monto := total * porcentaje
		resultado, err := tx.Exec(ctx,
			`INSERT INTO comisiones (pedido_id, referido_id, nivel, porcentaje, monto)
			 VALUES ($1, $2, $3, $4, $5)
			 ON CONFLICT (pedido_id, referido_id) DO NOTHING`,
			pedidoID, referido.id, referido.profundidad, porcentaje, monto,
		)
		if err != nil {
			return err
		}
		if resultado.RowsAffected() == 0 {
			continue
		}
		if _, err := tx.Exec(ctx,
			`UPDATE referidos
			 SET ventas = ventas + $1, comisiones = comisiones + $2
			 WHERE id = $3`,
			total, monto, referido.id,
		); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
