package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"multicatalogo-backend/models"
)

var (
	ErrUsuarioNoEncontrado  = errors.New("usuario no encontrado")
	ErrProductoNoEncontrado = errors.New("producto no encontrado")
	ErrEmailRegistrado      = errors.New("correo ya registrado")
	ErrCodigoReferido       = errors.New("código de referido no encontrado")
	ErrPedidoNoEncontrado   = errors.New("pedido no encontrado")
	ErrEstadoPedido         = errors.New("estado de pedido inválido")
	ErrComisionNoLiquidable = errors.New("la comisión ya fue liquidada o no existe")
	ErrLiquidacionVacia     = errors.New("no hay comisiones para liquidar")
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

func (r *PostgresRepository) BuscarUsuarioPorEmail(ctx context.Context, email string) (models.Usuario, bool, error) {
	var usuario models.Usuario
	err := r.pool.QueryRow(ctx,
		`SELECT id, nombre, email, password, rol, activo, referido_id
		 FROM usuarios WHERE email = $1`, email,
	).Scan(
		&usuario.ID,
		&usuario.Nombre,
		&usuario.Email,
		&usuario.Password,
		&usuario.Rol,
		&usuario.Activo,
		&usuario.ReferidoID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Usuario{}, false, nil
	}
	if err != nil {
		return models.Usuario{}, false, err
	}
	return usuario, true, nil
}

// RegistrarUsuario crea al usuario y su lugar en la red como una sola operación.
func (r *PostgresRepository) RegistrarUsuario(ctx context.Context, req models.RegisterRequest, passwordHash string) (models.Usuario, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return models.Usuario{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var yaExiste bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM usuarios WHERE email = $1)`, req.Email).Scan(&yaExiste); err != nil {
		return models.Usuario{}, err
	}
	if yaExiste {
		return models.Usuario{}, ErrEmailRegistrado
	}

	var padreID int
	var nivelPadre int
	if req.Ref == "" {
		err = tx.QueryRow(ctx,
			`SELECT id, nivel FROM referidos WHERE parent_id IS NULL ORDER BY id LIMIT 1 FOR UPDATE`,
		).Scan(&padreID, &nivelPadre)
	} else {
		err = tx.QueryRow(ctx,
			`SELECT id, nivel FROM referidos WHERE UPPER(codigo) = $1 FOR UPDATE`, req.Ref,
		).Scan(&padreID, &nivelPadre)
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Usuario{}, ErrCodigoReferido
	}
	if err != nil {
		return models.Usuario{}, err
	}

	var referidoID int
	if err := tx.QueryRow(ctx, `SELECT nextval('referidos_id_seq')`).Scan(&referidoID); err != nil {
		return models.Usuario{}, err
	}
	codigo := fmt.Sprintf("REF%06d", referidoID)
	if _, err := tx.Exec(ctx,
		`INSERT INTO referidos (id, nombre, codigo, nivel, ventas, comisiones, parent_id)
		 VALUES ($1, $2, $3, $4, 0, 0, $5)`,
		referidoID, req.Nombre, codigo, nivelPadre+1, padreID,
	); err != nil {
		return models.Usuario{}, err
	}

	usuario := models.Usuario{
		Nombre: req.Nombre, Email: req.Email, Password: passwordHash,
		Rol: "cliente", Activo: true, ReferidoID: referidoID,
	}
	err = tx.QueryRow(ctx,
		`INSERT INTO usuarios (nombre, email, password, rol, activo, referido_id)
		 VALUES ($1, $2, $3, 'cliente', true, $4) RETURNING id`,
		usuario.Nombre, usuario.Email, usuario.Password, usuario.ReferidoID,
	).Scan(&usuario.ID)
	if err != nil {
		var postgresErr *pgconn.PgError
		if errors.As(err, &postgresErr) && postgresErr.Code == "23505" {
			return models.Usuario{}, ErrEmailRegistrado
		}
		return models.Usuario{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Usuario{}, err
	}
	return usuario, nil
}

func (r *PostgresRepository) ObtenerProductos(ctx context.Context) ([]models.Producto, error) {
	filas, err := r.pool.Query(ctx,
		`SELECT id, nombre, descripcion, precio, categoria, img, galeria, stock, activo
		 FROM productos WHERE activo = true ORDER BY id`)
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
			&producto.Activo,
		); err != nil {
			return nil, err
		}
		productos = append(productos, producto)
	}
	return productos, filas.Err()
}

func (r *PostgresRepository) ObtenerTodosLosProductos(ctx context.Context) ([]models.Producto, error) {
	filas, err := r.pool.Query(ctx,
		`SELECT id, nombre, descripcion, precio, categoria, img, galeria, stock, activo
		 FROM productos ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	productos := make([]models.Producto, 0)
	for filas.Next() {
		var producto models.Producto
		if err := filas.Scan(&producto.ID, &producto.Nombre, &producto.Descripcion, &producto.Precio,
			&producto.Categoria, &producto.Img, &producto.Galeria, &producto.Stock, &producto.Activo); err != nil {
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
		 FROM productos WHERE id = $1 AND activo = true`, id,
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

func (r *PostgresRepository) CrearProducto(ctx context.Context, producto models.Producto) (models.Producto, error) {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO productos (nombre, descripcion, precio, categoria, img, galeria, stock, activo)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, true)
		 RETURNING id, activo`,
		producto.Nombre, producto.Descripcion, producto.Precio, producto.Categoria,
		producto.Img, producto.Galeria, producto.Stock,
	).Scan(&producto.ID, &producto.Activo)
	return producto, err
}

func (r *PostgresRepository) ActualizarProducto(ctx context.Context, id int, producto models.Producto) (models.Producto, error) {
	err := r.pool.QueryRow(ctx,
		`UPDATE productos
		 SET nombre = $1, descripcion = $2, precio = $3, categoria = $4,
		     img = $5, galeria = $6, stock = $7
		 WHERE id = $8
		 RETURNING id, activo`,
		producto.Nombre, producto.Descripcion, producto.Precio, producto.Categoria,
		producto.Img, producto.Galeria, producto.Stock, id,
	).Scan(&producto.ID, &producto.Activo)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Producto{}, ErrProductoNoEncontrado
	}
	return producto, err
}

func (r *PostgresRepository) DesactivarProducto(ctx context.Context, id int) error {
	resultado, err := r.pool.Exec(ctx, `UPDATE productos SET activo = false WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if resultado.RowsAffected() == 0 {
		return ErrProductoNoEncontrado
	}
	return nil
}

func (r *PostgresRepository) ObtenerRed(ctx context.Context, emailUsuario string) ([]*models.Referido, error) {
	filas, err := r.pool.Query(ctx,
		`WITH RECURSIVE red_usuario AS (
			SELECT referido.id, referido.nombre, referido.codigo, referido.nivel,
			       referido.ventas, referido.comisiones, referido.parent_id
			FROM usuarios usuario
			JOIN referidos referido ON referido.id = usuario.referido_id
			WHERE usuario.email = $1

			UNION ALL

			SELECT hijo.id, hijo.nombre, hijo.codigo, hijo.nivel,
			       hijo.ventas, hijo.comisiones, hijo.parent_id
			FROM referidos hijo
			JOIN red_usuario padre ON hijo.parent_id = padre.id
		)
		SELECT id, nombre, codigo, nivel, ventas, comisiones, parent_id
		FROM red_usuario ORDER BY nivel, id`, emailUsuario)
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
			&referido.Codigo,
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

func (r *PostgresRepository) ObtenerPedidosUsuario(ctx context.Context, email string) ([]models.PedidoResumen, error) {
	return r.consultarPedidos(ctx, `WHERE u.email = $1`, []any{email}, ` ORDER BY p.fecha DESC`)
}

func (r *PostgresRepository) ObtenerTodosLosPedidos(ctx context.Context, busqueda string, pagina int, limite int) ([]models.PedidoResumen, int, error) {
	if pagina < 1 {
		pagina = 1
	}
	if limite < 1 || limite > 100 {
		limite = 20
	}
	filtro := ""
	args := make([]any, 0, 2)
	if busqueda != "" {
		filtro = `WHERE p.numero ILIKE $1 OR u.email ILIKE $1 OR p.estado ILIKE $1`
		args = append(args, "%"+busqueda+"%")
	}
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM pedidos p JOIN usuarios u ON u.id = p.usuario_id `+filtro, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	args = append(args, limite, (pagina-1)*limite)
	pedidos, err := r.consultarPedidos(ctx, filtro, args, ` ORDER BY p.fecha DESC LIMIT $`+strconv.Itoa(len(args)-1)+` OFFSET $`+strconv.Itoa(len(args)))
	return pedidos, total, err
}

func (r *PostgresRepository) consultarPedidos(ctx context.Context, filtro string, args []any, suffix string) ([]models.PedidoResumen, error) {
	filas, err := r.pool.Query(ctx, `
		SELECT p.id, p.numero, p.total, p.estado, p.fecha,
		       p.cliente_nombre, p.cliente_email, p.direccion, p.ciudad, p.telefono,
		       u.email
		FROM pedidos p JOIN usuarios u ON u.id = p.usuario_id `+filtro+suffix, args...)
	if err != nil {
		return nil, err
	}
	defer filas.Close()
	pedidos := make([]models.PedidoResumen, 0)
	porID := make(map[int64]int)
	for filas.Next() {
		var pedido models.PedidoResumen
		if err := filas.Scan(&pedido.ID, &pedido.Numero, &pedido.Total, &pedido.Estado, &pedido.Fecha,
			&pedido.Cliente.Nombre, &pedido.Cliente.Email, &pedido.Cliente.Direccion,
			&pedido.Cliente.Ciudad, &pedido.Cliente.Telefono, &pedido.UsuarioEmail); err != nil {
			return nil, err
		}
		pedido.Items = make([]models.PedidoItemConfirmado, 0)
		porID[pedido.ID] = len(pedidos)
		pedidos = append(pedidos, pedido)
	}
	if err := filas.Err(); err != nil {
		return nil, err
	}
	if len(pedidos) == 0 {
		return pedidos, nil
	}
	ids := make([]int64, 0, len(pedidos))
	for _, pedido := range pedidos {
		ids = append(ids, pedido.ID)
	}
	items, err := r.pool.Query(ctx, `SELECT pedido_id, producto_id, nombre, cantidad, precio_unitario, subtotal
		FROM pedido_items WHERE pedido_id = ANY($1) ORDER BY pedido_id, id`, ids)
	if err != nil {
		return nil, err
	}
	defer items.Close()
	for items.Next() {
		var pedidoID int64
		var item models.PedidoItemConfirmado
		if err := items.Scan(&pedidoID, &item.ProductoID, &item.Nombre, &item.Cantidad, &item.PrecioUnitario, &item.Subtotal); err != nil {
			return nil, err
		}
		pedidos[porID[pedidoID]].Items = append(pedidos[porID[pedidoID]].Items, item)
	}
	return pedidos, items.Err()
}

func (r *PostgresRepository) ActualizarEstadoPedido(ctx context.Context, id int64, estado string) error {
	resultado, err := r.pool.Exec(ctx, `UPDATE pedidos SET estado = $1 WHERE id = $2`, estado, id)
	if err != nil {
		return err
	}
	if resultado.RowsAffected() == 0 {
		return ErrPedidoNoEncontrado
	}
	return nil
}

func (r *PostgresRepository) CancelarPedido(ctx context.Context, id int64) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var estado string
	if err := tx.QueryRow(ctx, `SELECT estado FROM pedidos WHERE id = $1 FOR UPDATE`, id).Scan(&estado); errors.Is(err, pgx.ErrNoRows) {
		return ErrPedidoNoEncontrado
	} else if err != nil {
		return err
	}
	if estado == "cancelado" {
		return ErrEstadoPedido
	}
	if estado != "procesando" && estado != "confirmado" && estado != "pendiente" {
		return ErrEstadoPedido
	}

	items, err := tx.Query(ctx, `SELECT producto_id, cantidad FROM pedido_items WHERE pedido_id = $1`, id)
	if err != nil {
		return err
	}
	type itemStock struct{ productoID, cantidad int }
	var stocks []itemStock
	for items.Next() {
		var item itemStock
		if err := items.Scan(&item.productoID, &item.cantidad); err != nil {
			items.Close()
			return err
		}
		stocks = append(stocks, item)
	}
	if err := items.Err(); err != nil {
		items.Close()
		return err
	}
	items.Close()
	for _, item := range stocks {
		if _, err := tx.Exec(ctx, `UPDATE productos SET stock = stock + $1 WHERE id = $2`, item.cantidad, item.productoID); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(ctx, `
		UPDATE referidos r SET ventas = r.ventas - c.total_ventas, comisiones = r.comisiones - c.total_comisiones
		FROM (SELECT referido_id, SUM((SELECT total FROM pedidos WHERE id = $1)) AS total_ventas,
		             SUM(monto) AS total_comisiones FROM comisiones WHERE pedido_id = $1 GROUP BY referido_id) c
		WHERE r.id = c.referido_id`, id); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM comisiones WHERE pedido_id = $1`, id); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE pedidos SET estado = 'cancelado' WHERE id = $1`, id); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *PostgresRepository) ObtenerComisionesUsuario(ctx context.Context, email string) ([]models.Comision, error) {
	return r.consultarComisiones(ctx, `
		JOIN usuarios u ON u.referido_id = c.referido_id
		WHERE u.email = $1`, []any{email})
}

func (r *PostgresRepository) ObtenerTodasLasComisiones(ctx context.Context, estado string) ([]models.Comision, error) {
	filtro := ""
	args := []any{}
	if estado != "" {
		filtro = " WHERE c.estado = $1"
		args = append(args, estado)
	}
	return r.consultarComisiones(ctx, filtro, args)
}

func (r *PostgresRepository) consultarComisiones(ctx context.Context, filtro string, args []any) ([]models.Comision, error) {
	filas, err := r.pool.Query(ctx, `
		SELECT c.id, c.pedido_id, p.numero, c.referido_id, r.nombre, c.nivel,
		       c.porcentaje, c.monto, c.estado, c.fecha, c.pagada_en
		FROM comisiones c
		JOIN pedidos p ON p.id = c.pedido_id
		JOIN referidos r ON r.id = c.referido_id`+filtro+` ORDER BY c.fecha DESC`, args...)
	if err != nil {
		return nil, err
	}
	defer filas.Close()
	comisiones := make([]models.Comision, 0)
	for filas.Next() {
		var comision models.Comision
		if err := filas.Scan(&comision.ID, &comision.PedidoID, &comision.NumeroPedido,
			&comision.ReferidoID, &comision.Referido, &comision.Nivel, &comision.Porcentaje,
			&comision.Monto, &comision.Estado, &comision.Fecha, &comision.PagadaEn); err != nil {
			return nil, err
		}
		comisiones = append(comisiones, comision)
	}
	return comisiones, filas.Err()
}

func (r *PostgresRepository) LiquidarComisiones(ctx context.Context, ids []int64) (models.Liquidacion, error) {
	if len(ids) == 0 {
		return models.Liquidacion{}, ErrLiquidacionVacia
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return models.Liquidacion{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var referidoID int
	var monto float64
	if err := tx.QueryRow(ctx, `
		SELECT referido_id, SUM(monto) FROM comisiones
		WHERE id = ANY($1) AND estado = 'pendiente'
		GROUP BY referido_id`, ids).Scan(&referidoID, &monto); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Liquidacion{}, ErrComisionNoLiquidable
		}
		return models.Liquidacion{}, err
	}
	var cantidad int
	if err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM comisiones WHERE id = ANY($1) AND estado = 'pendiente' AND referido_id = $2`, ids, referidoID).Scan(&cantidad); err != nil {
		return models.Liquidacion{}, err
	}
	if cantidad != len(ids) {
		return models.Liquidacion{}, ErrComisionNoLiquidable
	}

	var liquidacion models.Liquidacion
	if err := tx.QueryRow(ctx, `INSERT INTO liquidaciones (referido_id, monto) VALUES ($1, $2) RETURNING id, fecha`, referidoID, monto).Scan(&liquidacion.ID, &liquidacion.Fecha); err != nil {
		return models.Liquidacion{}, err
	}
	for _, id := range ids {
		if _, err := tx.Exec(ctx, `INSERT INTO liquidacion_comisiones (liquidacion_id, comision_id) VALUES ($1, $2)`, liquidacion.ID, id); err != nil {
			return models.Liquidacion{}, ErrComisionNoLiquidable
		}
	}
	if _, err := tx.Exec(ctx, `UPDATE comisiones SET estado = 'pagada', pagada_en = NOW() WHERE id = ANY($1)`, ids); err != nil {
		return models.Liquidacion{}, err
	}
	liquidacion.ReferidoID = referidoID
	liquidacion.Monto = monto
	if err := tx.Commit(ctx); err != nil {
		return models.Liquidacion{}, err
	}
	return liquidacion, nil
}
