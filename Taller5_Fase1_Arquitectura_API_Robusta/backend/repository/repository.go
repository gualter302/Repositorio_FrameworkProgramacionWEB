package repository

import "multicatalogo-backend/models"

// Por ahora los datos viven en memoria. Más adelante esta parte puede conectarse a una base de datos.
var usuarios = []models.Usuario{
	{Email: "admin@upse.edu.ec", Password: "123456", Token: "fake-jwt-token-123"},
}

var productos = []models.Producto{
	{ID: 1, Nombre: "Serum Revitalizante", Precio: 45.00, Img: "https://picsum.photos/seed/serum/150"},
	{ID: 2, Nombre: "Crema Hidratante Pro", Precio: 32.50, Img: "https://picsum.photos/seed/crema/150"},
	{ID: 3, Nombre: "Tónico Purificante", Precio: 28.00, Img: "https://picsum.photos/seed/tonico/150"},
	{ID: 4, Nombre: "Mascarilla Nocturna", Precio: 50.00, Img: "https://picsum.photos/seed/mascarilla/150"},
}

func BuscarUsuarioPorCredenciales(email, password string) (models.Usuario, bool) {
	for _, usuario := range usuarios {
		if usuario.Email == email && usuario.Password == password {
			return usuario, true
		}
	}

	return models.Usuario{}, false
}

func ObtenerProductos() []models.Producto {
	return productos
}

func ObtenerProductoPorID(id int) (models.Producto, bool) {
	for _, producto := range productos {
		if producto.ID == id {
			return producto, true
		}
	}

	return models.Producto{}, false
}
