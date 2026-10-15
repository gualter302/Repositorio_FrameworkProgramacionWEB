# Unidad 2 · Fase 4 — Checkout y concurrencia

Esta carpeta conserva todo lo construido anteriormente: el mult catálogo de la Fase 1,
PostgreSQL de la Fase 2 y Redis de la Fase 3. La Fase 4 agrega un checkout real, control de
existencias y cálculo de comisiones en segundo plano.

## Cómo funciona

Al confirmar una compra, el backend toma los precios desde PostgreSQL y abre una transacción.
Dentro de ella registra el pedido, guarda cada producto y descuenta sus existencias. Si una parte
falla, PostgreSQL deshace todo; por eso no quedan pedidos incompletos ni stock mal descontado.

Cuando la compra ya fue confirmada, Go responde al frontend y calcula aparte las comisiones de
los tres niveles superiores: 5 %, 3 % y 2 %. Al terminar elimina la copia antigua de la red en
Redis para que Dashboard y Mi Red reciban los valores nuevos en la siguiente consulta.

Los controladores reciben PostgreSQL y Redis al iniciar la aplicación. Esto evita dependencias
globales y permite probar las rutas con versiones falsas, sin necesitar contenedores.

## Requisitos

- Docker Desktop en ejecución.
- Node.js y npm para el frontend.
- Go solamente si se ejecutará el backend fuera de Docker.

## Iniciar el proyecto

Desde `backend`:

```bash
docker compose up --build -d
```

Desde `frontend`:

```bash
npm install
npm run dev
```

Abre `http://localhost:5173`. Para hacer una compra que genere tres comisiones, inicia sesión con:

- Correo: `cliente@upse.edu.ec`
- Contraseña: `123456`

La cuenta administradora sigue disponible como `admin@upse.edu.ec` con la misma contraseña.

## Endpoints principales

- `POST /api/login`: valida la cuenta.
- `GET /api/productos`: devuelve el catálogo y usa Redis.
- `GET /api/red?email=...`: devuelve la red y sus comisiones.
- `POST /api/checkout`: registra la compra de forma segura.
- `GET /api/health`: confirma que el backend está activo.

Ejemplo de checkout:

```json
{
  "usuarioEmail": "cliente@upse.edu.ec",
  "cliente": {
    "nombre": "Cliente Demo",
    "email": "cliente@upse.edu.ec",
    "direccion": "Barrio Central",
    "ciudad": "La Libertad",
    "telefono": "0999999999"
  },
  "metodoPago": "Efectivo contra entrega",
  "items": [{ "productoId": 1, "cantidad": 1 }]
}
```

## Comprobar la reversión

Envía un carrito con un producto válido y luego otro con id `99999`. La API responde 400 y
PostgreSQL no conserva la cabecera, los detalles ni el descuento del primer producto. Esta prueba
permite demostrar la reversión sin apagar contenedores ni dañar datos.

## Pruebas

En `backend`:

```bash
go test ./...
```

En `frontend`:

```bash
npm run lint
npm run build
```

Para detener los servicios sin borrar los datos:

```bash
docker compose down
```
