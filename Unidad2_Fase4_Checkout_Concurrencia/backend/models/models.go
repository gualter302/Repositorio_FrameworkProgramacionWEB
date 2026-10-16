package models

import "time"

// LoginRequest representa los datos que llegan desde el formulario de inicio de sesión.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Usuario guarda los datos que el repositorio necesita para validar una cuenta.
// Password lleva json:"-" para que nunca viaje al frontend.
type Usuario struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Rol      string `json:"rol"`
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
}

// Referido es un nodo de la red multinivel. Hijos se llena al reconstruir el árbol.
type Referido struct {
	ID         int         `json:"id"`
	Nombre     string      `json:"nombre"`
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

// APIError mantiene la misma forma en todas las respuestas de error.
// Message se publica como "error" porque el frontend ya lee esa propiedad.
type APIError struct {
	Status  int         `json:"status"`
	Message string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}
