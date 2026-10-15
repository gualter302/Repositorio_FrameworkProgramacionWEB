package models

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
}

// Referido es un nodo de la red multinivel. Hijos se llena al reconstruir el árbol.
type Referido struct {
	ID       int         `json:"id"`
	Nombre   string      `json:"nombre"`
	Nivel    int         `json:"nivel"`
	Ventas   float64     `json:"ventas"`
	ParentID *int        `json:"parentId,omitempty"`
	Hijos    []*Referido `json:"hijos"`
}

// APIError mantiene la misma forma en todas las respuestas de error.
// Message se publica como "error" porque el frontend ya lee esa propiedad.
type APIError struct {
	Status  int         `json:"status"`
	Message string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}
