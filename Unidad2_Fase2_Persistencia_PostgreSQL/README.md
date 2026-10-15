# Unidad 2 · Fase 2 — Persistencia con PostgreSQL

Contiene el mult catálogo completo de la Fase 1 y agrega persistencia con PostgreSQL.
El frontend se mantiene dentro de esta carpeta y conserva el mismo contrato con la API.

## Qué incluye

- `backend/docker-compose.yml`: levanta PostgreSQL 16 y ejecuta el script de creación.
- `backend/database/init.sql`: tablas `usuarios`, `productos` y `referidos` con sus datos semilla.
- `backend/config/db.go`: pool de conexiones con `pgx`, lectura de `.env` y cierre ordenado.
- `backend/repository/repository.go`: todas las consultas (`SELECT`) a la base de datos.
- `backend/controllers/redController.go`: nuevo endpoint `GET /api/red`.
- `frontend/`: aplicación React completa heredada de la Fase 1.

## Requisitos

- Docker Desktop en ejecución.
- Go 1.26 o superior.
- Node.js y npm para ejecutar el frontend.

## Cómo levantarlo

1. Copiar la configuración de ejemplo:

   ```bash
   cp backend/.env.example backend/.env
   ```

2. Levantar PostgreSQL (desde la carpeta `backend`):

   ```bash
   docker compose up -d
   ```

   El archivo `init.sql` se ejecuta automáticamente **solo la primera vez** que se crea el
   volumen. Si cambias el script, recrea el contenedor con `docker compose down -v` y vuelve
   a levantarlo.

3. Iniciar el backend (desde la carpeta `backend`):

   ```bash
   go run .
   ```

   El servidor queda en `http://localhost:3000`.

4. Instalar e iniciar el frontend (desde la carpeta `frontend`):

   ```bash
   npm install
   npm run dev
   ```

   El mult catálogo queda disponible en `http://localhost:5173`.

## Endpoints

- `GET /api/health` — estado del servidor.
- `POST /api/login` — valida credenciales contra la tabla `usuarios` (respuesta: `token`, `email`).
- `GET /api/productos` — catálogo completo.
- `GET /api/productos/:id` — detalle de un producto.
- `GET /api/red` — red de referidos como árbol multinivel.

## Validar la persistencia

1. Conéctate a la base (DBeaver, pgAdmin o terminal) y cambia un precio:

   ```sql
   UPDATE productos SET precio = 99.99 WHERE id = 1;
   ```

2. Recarga el catálogo en el frontend (`/catalogo` o `/tienda`) o consulta
   `http://localhost:3000/api/productos`. El "Serum Revitalizante" mostrará el nuevo precio.

El frontend incluido funciona sin cambios: ya apunta al backend en el puerto `3000`.

## Pruebas

Desde `backend`:

```bash
go test ./...
```

Las pruebas de validación y de estado se ejecutan siempre. Las que necesitan datos de la base
(`TestLoginConservaRespuestaDelFrontend`, `TestProductoPorID`) se omiten automáticamente cuando
no hay una conexión a PostgreSQL disponible.

Para comprobar el frontend, ejecuta `npm run lint` y `npm run build` desde `frontend`.
