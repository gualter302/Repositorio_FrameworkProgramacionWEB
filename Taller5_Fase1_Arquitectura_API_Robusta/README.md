# Taller 5 - Fase 1: arquitectura y API robusta

Este taller parte de la conexión entre el frontend React y el backend Go realizada en el Taller 4.
El frontend se conserva igual; los cambios se concentran en el backend para responder de forma más
clara cuando una petición tiene datos incompletos o incorrectos.

## Estructura

```text
Taller5_Fase1_Arquitectura_API_Robusta/
├── frontend/       React + Vite (puerto 5173)
└── backend/        Go + Fiber (puerto 3000)
    ├── controllers/  recibe las peticiones y prepara las respuestas
    ├── models/       define la forma de los datos
    ├── repository/   guarda temporalmente usuarios y productos
    └── routes/       registra las direcciones de la API
```

## Qué se hizo

1. Se agregó `APIError` para que los errores siempre incluyan un estado y un mensaje. Los detalles
   son opcionales y sirven para indicar, por ejemplo, qué campo está vacío.
2. El inicio de sesión comprueba que el correo y la contraseña tengan contenido antes de buscar la
   cuenta. Si falta alguno, responde con estado 400.
3. Se agregó `GET /api/productos/:id`. Si el identificador no es un número entero positivo,
   responde con estado 400; si el producto no existe, responde con estado 404.
4. Se agregó `GET /api/health` para comprobar que el servidor está disponible sin iniciar sesión.
5. Los datos fijos se movieron a `repository/`, dejando a los controladores concentrados en recibir
   y responder peticiones.

La propiedad `Message` de `APIError` se envía como `error` en el JSON. Esto mantiene funcionando al
frontend del Taller 4, que ya usa esa propiedad para mostrar los mensajes en pantalla.

## Endpoints

| Método | Dirección | Uso |
|---|---|---|
| `GET` | `/api/health` | Confirma que el backend está activo |
| `POST` | `/api/login` | Valida el correo y la contraseña |
| `GET` | `/api/productos` | Devuelve todo el catálogo |
| `GET` | `/api/productos/:id` | Devuelve un producto por su número |

Ejemplo de error:

```json
{
  "status": 400,
  "error": "Faltan campos obligatorios",
  "details": {
    "email": "El correo es obligatorio"
  }
}
```

## Cómo ejecutarlo

Backend:

```bash
cd Taller5_Fase1_Arquitectura_API_Robusta/backend
go run .
```

Frontend, en otra terminal:

```bash
cd Taller5_Fase1_Arquitectura_API_Robusta/frontend
npm install
npm run dev
```

Abre `http://localhost:5173` e inicia sesión con `admin@upse.edu.ec` y la contraseña `123456`.

## Verificación

Las pruebas del backend revisan el estado del servidor, la validación del login y la consulta de
productos:

```bash
cd Taller5_Fase1_Arquitectura_API_Robusta/backend
go test ./...
```
