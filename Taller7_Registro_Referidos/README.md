# Taller 7 — Registro y red de referidos

Este taller parte de la Fase 4 completa y agrega la puerta de entrada del negocio:
una persona puede crear su cuenta, usar el código de quien la invitó y obtener su
propio lugar dentro de la red multinivel.

## Qué se agregó

- `POST /api/register` para crear la cuenta y el nodo de red en una transacción.
- Formulario público en `/registro`, incluyendo `?ref=CODIGO`.
- Contraseñas protegidas con bcrypt; ya no se guardan como texto normal.
- JWT real con dos horas de duración.
- Checkout y Mi Red protegidos: la identidad sale del token, no del cuerpo enviado.
- Una red propia por usuario y enlace de invitación con botón para copiar.
- Migración `004_registro_referidos.sql` para bases creadas en la fase anterior.
- Invalidación de las redes guardadas en Redis cuando se registra alguien o cambian comisiones.
- CRUD de productos para administradores, con desactivación lógica e invalidación del catálogo en Redis.
- Historial de pedidos para clientes y panel administrativo con búsqueda, paginación, estados y cancelación transaccional.
- Libro auditable de comisiones, consulta personal y liquidaciones protegidas contra duplicados.
- Reporte administrativo de ventas, pedidos, ticket promedio, comisiones y productos principales.

## Iniciar el backend

Con Docker Desktop encendido, entra en `backend` y ejecuta:

```bash
docker compose up --build -d
```

Los servicios quedan disponibles en los puertos 3000, 5432 y 6379. El volumen del
Taller 7 es independiente y no modifica los datos de la Fase 4.

Si el volumen ya existía antes de esta versión, aplica una vez la migración:

```bash
docker compose exec db psql -U multicatalogo -d multicatalogo -f /docker-entrypoint-initdb.d/005_productos_activos.sql
```

## Iniciar el frontend

En `frontend`:

```bash
npm install
npm run dev
```

Abre `http://localhost:5173/registro?ref=ANA` para probar una invitación.

## Cuentas iniciales

- Administrador: `admin@upse.edu.ec` / `123456`
- Cliente: `cliente@upse.edu.ec` / `123456`

Los códigos iniciales incluyen `ADMIN`, `ANA`, `LUIS`, `MARIA`, `CARLOS` y `DEMO`.

## Respuestas esperadas

- Registro correcto: 201 y una sesión lista para usar.
- Correo repetido: 409.
- Código de referido inexistente: 400.
- Checkout o Mi Red sin token: 401.
- Cuenta desactivada: 403 al iniciar sesión.
- Cliente intentando administrar productos: 403.
- Cliente consultando `/api/pedidos`: solo recibe sus propios pedidos.
- Liquidación repetida de una comisión: 409.
- Cliente intentando consultar reportes administrativos: 403.

## Pruebas

Backend:

```bash
go test ./...
```

Frontend:

```bash
npm run lint
npm run build
```
