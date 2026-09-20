# Taller 4 - Conexión del Frontend (React) con el Backend (Go)

El **frontend** es el proyecto del Taller 3 (con el sidebar terminado). El **backend** es la API REST en Go con Fiber (`multicatalogo-backend`).

```
Taller4_Conexion_Backend_Frontend/
├── frontend/   -> React + Vite (puerto 5173)
│   └── src/services/api.ts   <- ÚNICO archivo que habla con el backend
└── backend/    -> Go + Fiber (puerto 3000)
    ├── main.go               -> arranca el servidor y configura CORS
    ├── routes/routes.go      -> define las rutas /api/login y /api/productos
    ├── controllers/          -> la lógica de cada ruta
    └── models/models.go      -> las estructuras JSON (LoginRequest, Producto)
```

> La carpeta `backend` va **al lado** del frontend (no dentro de `src`): `src` es código que se
> ejecuta en el navegador; el backend se ejecuta en el servidor y es otro programa (Go).

## Cómo se conectan (idea general)

Son **dos programas independientes** que se hablan por HTTP enviándose JSON:

```
Navegador (React :5173)  --- fetch("http://localhost:3000/api/...") --->  Go/Fiber (:3000)
                         <----------------- respuesta JSON --------------
```

| Acción en el frontend | Petición | Lo atiende en Go |
|---|---|---|
| Iniciar sesión (`Login.tsx`) | `POST /api/login` con `{email, password}` | `controllers.Login` |
| Ver el catálogo (`Catalogo.tsx`) | `GET /api/productos` | `controllers.GetProductos` |

Pasos que se hicieron:

1. **Backend encendido** en `:3000` con las rutas `/api/login` y `/api/productos`.
2. **CORS**: el navegador bloquea las peticiones entre puertos distintos (5173 → 3000) salvo que el backend lo autorice. En `main.go` se autorizan `http://localhost:5173` y `http://127.0.0.1:5173` (configurable con la variable `ALLOWED_ORIGINS`).
3. **`src/services/api.ts`** centraliza las llamadas con `fetch` y la URL base (`VITE_API_URL`, por defecto `http://localhost:3000/api`).
4. **`Login.tsx`** dejó de comparar contra un usuario quemado: ahora llama a `loginRequest()` y, si el backend responde 200, guarda el correo y el **token** en `AuthContext`. Si responde 401 muestra el error del backend.
5. **`Catalogo.tsx`** dejó de tener el arreglo de productos escrito a mano: los pide con `getProductos()` dentro de un `useEffect` y muestra "Cargando..." o el error si falla.

## Cómo ejecutarlo (2 terminales)

**Terminal 1 - Backend** (requiere tener [Go](https://go.dev/dl/) instalado):

```bash
cd backend
go run .
```

**Terminal 2 - Frontend** (requiere Node.js):

```bash
cd frontend
npm install
npm run dev
```

Abrir `http://localhost:5173` e ingresar con `admin@upse.edu.ec` / `123456`.

### Si el frontend corre en otra dirección (ej. la IP del laboratorio)

Backend, autorizando esa dirección:

```bash
# Linux / macOS / Git Bash
ALLOWED_ORIGINS="http://172.17.82.108:5173" go run .
# PowerShell
$env:ALLOWED_ORIGINS="http://172.17.82.108:5173"; go run .
```

Frontend: copiar `.env.example` a `.env` y poner `VITE_API_URL=http://<IP-del-backend>:3000/api`.

## Problemas típicos

- **"No se pudo conectar con el servidor"**: el backend está apagado o el puerto/URL no coincide.
- **Error de CORS en la consola del navegador**: el origen del frontend no está en `ALLOWED_ORIGINS`.
- **Las rutas de la API no existen (404)**: recuerda que empiezan por `/api`.

## Nota

El backend valida un único usuario fijo y devuelve un token de ejemplo (`fake-jwt-token-123`); es suficiente para el taller. Un proyecto real usaría base de datos, contraseñas cifradas y un JWT verdadero.
