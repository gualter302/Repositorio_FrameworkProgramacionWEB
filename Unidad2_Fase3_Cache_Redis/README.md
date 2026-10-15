# Unidad 2 · Fase 3 — Caché con Redis

Esta carpeta contiene el mult catálogo completo de la Fase 1, la persistencia en
PostgreSQL de la Fase 2 y la caché con Redis de la Fase 3. Redis permite responder
más rápido las consultas repetidas sin cambiar el formato que recibe el frontend.

## Qué incluye

- PostgreSQL 16 como fuente principal de usuarios, productos y referidos.
- Redis para guardar el catálogo durante 5 minutos.
- Redis para guardar cada árbol de referidos durante 10 minutos.
- Encabezado `X-Cache` con `MISS` en la primera consulta y `HIT` en las siguientes.
- Docker Compose para iniciar PostgreSQL, Redis y el backend en conjunto.
- El frontend React completo heredado de las fases anteriores.

Si Redis no contiene la clave, el backend consulta PostgreSQL y guarda una copia temporal.
Cuando la clave existe, la respuesta sale de Redis y se evita repetir la consulta a la base.

## Requisitos

- Docker Desktop en ejecución.
- Go 1.26 o superior, solamente si ejecutarás el backend fuera de Docker.
- Node.js y npm para ejecutar el frontend.

## Iniciar todo con Docker

Desde la carpeta `backend`:

```bash
docker compose up --build -d
```

Este comando inicia los tres servicios. El backend espera a que PostgreSQL y Redis estén
saludables antes de arrancar y queda disponible en `http://localhost:3000`.

Para detenerlos sin borrar los datos de PostgreSQL:

```bash
docker compose down
```

## Ejecutar el backend de forma local

1. Copia `.env.example` como `.env`.
2. Inicia solamente las dependencias con `docker compose up -d db redis`.
3. Ejecuta `go run .` dentro de `backend`.

## Iniciar el frontend

Desde la carpeta `frontend`:

```bash
npm install
npm run dev
```

El mult catálogo queda disponible en `http://localhost:5173` y consume la API del
backend en `http://localhost:3000`.

## Endpoints con caché

- `GET /api/productos`: usa la clave `productos:all` durante 5 minutos.
- `GET /api/red`: usa `red:tree:admin@upse.edu.ec` durante 10 minutos.
- `GET /api/red?email=correo@ejemplo.com`: separa la copia temporal por correo.

Los demás endpoints conservan el comportamiento de la Fase 2.

## Comprobar la diferencia

Haz dos consultas seguidas y revisa el encabezado `X-Cache`:

```bash
curl.exe -i http://localhost:3000/api/productos
curl.exe -i http://localhost:3000/api/productos
```

La primera respuesta debe indicar `X-Cache: MISS`; la segunda, `X-Cache: HIT`. También
puedes confirmar el tiempo restante de la clave:

```bash
docker exec multicatalogo-fase3-redis redis-cli TTL productos:all
```

## Pruebas

Desde `backend`:

```bash
go test ./...
```

Las pruebas que no requieren servicios externos siempre se ejecutan. Las pruebas heredadas
que consultan datos se omiten cuando PostgreSQL no fue conectado desde el proceso de prueba.

Para comprobar el frontend, ejecuta `npm run lint` y `npm run build` desde `frontend`.
