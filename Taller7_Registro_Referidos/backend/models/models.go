package models

import "time"

// LoginRequest representa los datos que llegan desde el formulario de inicio de sesión.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest contiene los datos para crear una cuenta y enlazarla a la red.
type RegisterRequest struct {
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Ref      string `json:"ref"`
}

// AuthResponse es la sesión que comparten el login y el registro.
type AuthResponse struct {
	Token  string `json:"token"`
	Email  string `json:"email"`
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

// Usuario guarda los datos que el repositorio necesita para validar una cuenta.
// Password lleva json:"-" para que nunca viaje al frontend.
type Usuario struct {
	ID         int    `json:"id"`
	Nombre     string `json:"nombre"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	Rol        string `json:"rol"`
	Activo     bool   `json:"activo"`
	ReferidoID int    `json:"referidoId"`
}

// Producto representa un artículo del catálogo.
type Producto struct {
	ID          int      `json:"id"`
	Nombre      string   `json:"nombre"`
	Descripcion string   `json:"descripcion"`
	Precio      float64  `json:"precio"`
	Categoria   string   `json:"categoria"`
	Img         string   `json:"img"`
	Galeria     []string `json:"galeria"`
	Stock       int      `json:"stock"`
	Activo      bool     `json:"activo"`
}

// Referido es un nodo de la red multinivel. Hijos se llena al reconstruir el árbol.
type Referido struct {
	ID         int         `json:"id"`
	Nombre     string      `json:"nombre"`
	Codigo     string      `json:"codigo"`
	Nivel      int         `json:"nivel"`
	Ventas     float64     `json:"ventas"`
	Comisiones float64     `json:"comisiones"`
	ParentID   *int        `json:"parentId,omitempty"`
	Hijos      []*Referido `json:"hijos"`
}

// ClienteCheckout contiene los datos necesarios para entregar el pedido.
type ClienteCheckout struct {
	Nombre    string `json:"nombre"`
	Email     string `json:"email"`
	Direccion string `json:"direccion"`
	Ciudad    string `json:"ciudad"`
	Telefono  string `json:"telefono"`
}

// CheckoutItemRequest identifica el producto y la cantidad solicitada.
type CheckoutItemRequest struct {
	ProductoID int `json:"productoId"`
	Cantidad   int `json:"cantidad"`
}

// CheckoutRequest es la compra enviada desde el frontend.
type CheckoutRequest struct {
	UsuarioEmail string                `json:"usuarioEmail"`
	Cliente      ClienteCheckout       `json:"cliente"`
	MetodoPago   string                `json:"metodoPago"`
	Items        []CheckoutItemRequest `json:"items"`
}

// PedidoItemConfirmado muestra el precio realmente usado por el servidor.
type PedidoItemConfirmado struct {
	ProductoID     int     `json:"productoId"`
	Nombre         string  `json:"nombre"`
	Cantidad       int     `json:"cantidad"`
	PrecioUnitario float64 `json:"precio"`
	Subtotal       float64 `json:"subtotal"`
}

// PedidoConfirmado es la respuesta que recibe la pantalla de confirmación.
type PedidoConfirmado struct {
	ID      int64                  `json:"id"`
	Numero  string                 `json:"numero"`
	Total   float64                `json:"total"`
	Estado  string                 `json:"estado"`
	Fecha   time.Time              `json:"fecha"`
	Cliente ClienteCheckout        `json:"cliente"`
	Items   []PedidoItemConfirmado `json:"items"`
}

type PedidoResumen struct {
	PedidoConfirmado
	UsuarioEmail string `json:"usuarioEmail"`
}

type Comision struct {
	ID           int64      `json:"id"`
	PedidoID     int64      `json:"pedidoId"`
	NumeroPedido string     `json:"numeroPedido"`
	ReferidoID   int        `json:"referidoId"`
	Referido     string     `json:"referido"`
	Nivel        int        `json:"nivel"`
	Porcentaje   float64    `json:"porcentaje"`
	Monto        float64    `json:"monto"`
	Estado       string     `json:"estado"`
	Fecha        time.Time  `json:"fecha"`
	PagadaEn     *time.Time `json:"pagadaEn,omitempty"`
}

type Liquidacion struct {
	ID         int64     `json:"id"`
	ReferidoID int       `json:"referidoId"`
	Monto      float64   `json:"monto"`
	Fecha      time.Time `json:"fecha"`
}

type ProductoReporte struct {
	Nombre   string  `json:"nombre"`
	Cantidad int     `json:"cantidad"`
	Ventas   float64 `json:"ventas"`
}

type ReferidoReporte struct {
	Nombre     string  `json:"nombre"`
	Ventas     float64 `json:"ventas"`
	Comisiones float64 `json:"comisiones"`
}

type ReporteGeneral struct {
	TotalVentas          float64           `json:"totalVentas"`
	TotalPedidos         int               `json:"totalPedidos"`
	TicketPromedio       float64           `json:"ticketPromedio"`
	TotalComisiones      float64           `json:"totalComisiones"`
	PendientePagar       float64           `json:"pendientePagar"`
	ProductosPrincipales []ProductoReporte `json:"productosPrincipales"`
	ReferidosPrincipales []ReferidoReporte `json:"referidosPrincipales"`
}

// APIError mantiene la misma forma en todas las respuestas de error.
// Message se publica como "error" porque el frontend ya lee esa propiedad.
type APIError struct {
	Status  int         `json:"status"`
	Message string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}
