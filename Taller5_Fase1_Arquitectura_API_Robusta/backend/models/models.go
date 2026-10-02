package models

// LoginRequest representa los datos que llegan desde el formulario de inicio de sesión.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Usuario guarda los datos que el repositorio necesita para validar una cuenta.
type Usuario struct {
	Email    string
	Password string
	Token    string
}

// Producto representa un artículo del catálogo.
type Producto struct {
	ID     int     `json:"id"`
	Nombre string  `json:"nombre"`
	Precio float64 `json:"precio"`
	Img    string  `json:"img"`
}

// APIError mantiene la misma forma en todas las respuestas de error.
// Message se publica como "error" porque el frontend ya lee esa propiedad.
type APIError struct {
	Status  int         `json:"status"`
	Message string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}
