# Taller 5 - Lección Tema 5: interfaces a pantalla completa y catálogo de venta multinivel

Parte del **Taller 4** (frontend React conectado al backend Go/Fiber) y agrega lo pedido en la
*Guía de Laboratorio — Práctica 03 (Unidad 1, Tema 5)*: roles `admin` / `cliente`, tienda a
pantalla completa, catálogo con filtros, detalle con lightbox, carrito persistente por usuario,
checkout, red multinivel y dashboard con KPIs calculados.

```
Taller5_Leccion_Catalogo_Multinivel/
├── backend/    -> Go + Fiber (puerto 3000). /api/login ahora devuelve el rol
└── frontend/   -> React + Vite + Tailwind (puerto 5173)
    └── src/
        ├── data/        -> productos.ts (mock del catálogo) y red.ts (árbol MLM + funciones puras)
        ├── services/    -> api.ts (backend real) y productosService.ts (mock asíncrono)
        ├── utils/       -> rutas.ts (ruta de inicio de cada rol)
        ├── context/     -> AuthContext (usuario + rol), CartContext (carrito por usuario), SidebarContext
        └── components/  -> Storefront, Catalogo, DetalleProducto, Carrito, Checkout, Confirmacion, MiRed, Dashboard...
```

## Entregables

| # | Entregable | Ruta / vista |
|---|---|---|
| F1 | Storefront a pantalla completa (hero + secciones) | `/tienda` |
| F2 | Ruta dinámica de detalle | `/producto/:id` |
| F3 | Búsqueda + filtro por categoría (`useMemo`, `?categoria=`) | `/catalogo` |
| F4 | Galería y lightbox (cierra con `Escape`) | `DetalleProducto` |
| F5 | Carrito persistente por usuario + checkout + confirmación | `/carrito`, `/checkout`, `/confirmacion` |
| F6 | Árbol multinivel recursivo con comisiones | `/mi-red` |
| F7 | Dashboard con KPIs calculados | `/` |
| F8 | Login con rol y navegación según el rol | `Login`, `Sidebar`, `Navbar` |

## Rutas y roles

| Ruta | Acceso |
|---|---|
| `/login` | público |
| `/` (Dashboard), `/mi-red` | solo admin (un cliente es enviado a `/tienda`) |
| `/tienda`, `/catalogo`, `/producto/:id`, `/carrito`, `/checkout`, `/confirmacion` | admin y cliente |

## Cómo ejecutarlo (2 terminales)

**Terminal 1 - Backend** (requiere Go):

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

Abrir `http://localhost:5173`.

| Rol | Correo | Contraseña |
|---|---|---|
| Admin | `admin@upse.edu.ec` | `123456` |
| Cliente | `cliente@upse.edu.ec` | `123456` |

Si el frontend corre en otra dirección (IP del laboratorio), ver el README del Taller 4:
`ALLOWED_ORIGINS` en el backend y `VITE_API_URL` en `frontend/.env`.

## Diferencias con la guía (adaptaciones a este proyecto)

- **Login** sigue usando `src/services/api.ts` (`loginRequest`) del Taller 4 en lugar de un `fetch` directo, y el `AuthContext` guarda también el **token**.
- **`src/utils/rutas.ts`** centraliza la ruta de inicio de cada rol (admin → `/`, cliente → `/tienda`).
- **Sidebar / Navbar** conservan el diseño del Taller 3 (`SidebarContext`, iconos SVG, `NavLink`); se agregó el icono `TiendaIcon`, el filtro por rol y la insignia de rol.
- **Catálogo**: la categoría se lee directamente de la URL (`?categoria=`), así el filtro no se desincroniza al navegar desde el Sidebar o el Storefront.
- **Storefront**: el hero ocupa el alto visible (`100dvh` menos el Navbar) y anula el padding del `<main>` para verse de borde a borde.
- **Carrito**: mantiene el diseño responsivo del Taller 3 en celular.

## Verificación rápida

1. Admin → cae en el Dashboard: **$4,090** ventas, **7** referidos, **$319.50** comisiones, nivel **Plata**.
2. Cliente → cae en `/tienda`; el Sidebar solo muestra Tienda y Catálogo; si escribe `/mi-red` vuelve a `/tienda`.
3. Añadir al carrito, recargar: el carrito sigue ahí. Cerrar sesión y entrar con la otra cuenta: cada una tiene su propio carrito.
4. Carrito → Checkout → Confirmar: aparece el número de pedido y el carrito queda vacío.

## Nota

El catálogo del Tema 5 usa datos mock (`productosService.ts`) porque necesita descripción, categoría y
galería, que la API todavía no entrega. En la Unidad 2 solo se cambiará el cuerpo de
`getProductos()` / `getProductoById()` por `fetch` a la API (PostgreSQL), sin tocar los componentes.
