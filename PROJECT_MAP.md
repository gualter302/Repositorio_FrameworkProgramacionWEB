# Mapa compacto del proyecto

Consulta este archivo antes de buscar código. Las rutas que aparecen aquí corresponden a la
estructura real del repositorio y reemplazan los nombres generales `Unidad1_Frontend` y
`Unidad2_Backend` usados en las guías de clase.

## Tecnologías

- Frontend: React 19, TypeScript, Vite, Tailwind CSS y React Router.
- Backend: Go con Fiber, respuestas JSON y CORS para el frontend local.
- Pruebas: `go test` para las rutas del backend; ESLint y compilación de Vite para el frontend.
- Puertos locales: frontend `5173`; backend `3000`.
- Datos actuales: PostgreSQL como fuente principal y Redis como caché desde las fases 2 y 3.

## Archivos de orientación

- `README.md`: índice de talleres y forma general de ejecutar el proyecto.
- `PROJECT_MAP.md`: ubicación rápida de carpetas y archivos importantes.
- `REQUIREMENTS_CHECKLIST.md`: avance confirmado y siguientes tareas.
- `ANALISIS_FUNCIONAL_U2.md`: funciones de negocio que todavía faltan.
- `PLAN_IMPLEMENTACION_U2.md`: orden de los próximos talleres de back-office.
- `AGENTS.md`: instrucciones que Codex carga automáticamente para este repositorio.
- `CODEX_SETUP.md`: configuración opcional de Playwright y memoria para Codex.
- `CODEX_CONFIG_EXAMPLE.toml`: ejemplo inactivo para la configuración personal de Codex.

## Taller 1 y 2

- `Taller1y2_Entorno_y_API_Pokemon/`: aplicación React que consume PokeAPI.
- `Taller1y2_Entorno_y_API_Pokemon/src/App.tsx`: rutas principales.
- `Taller1y2_Entorno_y_API_Pokemon/src/pages/Pokemons.tsx`: listado y filtros.
- `Taller1y2_Entorno_y_API_Pokemon/src/components/`: layout, menú y tarjetas.
- `Taller1y2_Entorno_y_API_Pokemon/package.json`: comandos y dependencias.

## Taller 3

- `Taller3_Sidebar/`: base MultiCatálogo con navegación adaptable.
- `Taller3_Sidebar/src/App.tsx`: rutas protegidas de la aplicación.
- `Taller3_Sidebar/src/components/Sidebar.tsx`: menú lateral.
- `Taller3_Sidebar/src/components/Navbar.tsx`: barra superior y botón del menú.
- `Taller3_Sidebar/src/context/`: autenticación, carrito y estado del menú.
- `Taller3_Sidebar/src/icons/`: iconos reutilizables.

## Taller 4

- `Taller4_Conexion_Backend_Frontend/`: primera integración completa React y Go.
- `Taller4_Conexion_Backend_Frontend/frontend/src/services/api.ts`: único acceso a la API.
- `Taller4_Conexion_Backend_Frontend/frontend/src/components/Login.tsx`: inicio de sesión.
- `Taller4_Conexion_Backend_Frontend/frontend/src/components/Catalogo.tsx`: catálogo remoto.
- `Taller4_Conexion_Backend_Frontend/backend/main.go`: inicio del servidor y CORS.
- `Taller4_Conexion_Backend_Frontend/backend/routes/routes.go`: rutas del backend.
- `Taller4_Conexion_Backend_Frontend/backend/controllers/`: login y productos.
- `Taller4_Conexion_Backend_Frontend/backend/models/models.go`: datos compartidos.

## Taller 5 y backend actual

- `Taller5_Fase1_Arquitectura_API_Robusta/`: base recomendada para continuar el proyecto.
- `Taller5_Fase1_Arquitectura_API_Robusta/frontend/`: copia funcional del frontend del Taller 4.
- `Taller5_Fase1_Arquitectura_API_Robusta/frontend/src/services/api.ts`: conexión con Go.
- `Taller5_Fase1_Arquitectura_API_Robusta/backend/main.go`: servidor en el puerto 3000.
- `Taller5_Fase1_Arquitectura_API_Robusta/backend/routes/routes.go`: health, login y productos.
- `Taller5_Fase1_Arquitectura_API_Robusta/backend/controllers/authController.go`: validación del login.
- `Taller5_Fase1_Arquitectura_API_Robusta/backend/controllers/prodController.go`: lista y detalle.
- `Taller5_Fase1_Arquitectura_API_Robusta/backend/controllers/healthController.go`: estado del servidor.
- `Taller5_Fase1_Arquitectura_API_Robusta/backend/controllers/errors.go`: formato común de errores.
- `Taller5_Fase1_Arquitectura_API_Robusta/backend/models/models.go`: modelos y `APIError`.
- `Taller5_Fase1_Arquitectura_API_Robusta/backend/repository/repository.go`: datos en memoria.
- `Taller5_Fase1_Arquitectura_API_Robusta/backend/routes/routes_test.go`: pruebas de rutas.

## Unidad 2 Fases 2, 3 y 4

- `Unidad2_Fase2_Persistencia_PostgreSQL/`: Fase 1 más persistencia en PostgreSQL.
- `Unidad2_Fase3_Cache_Redis/`: Fases 1 y 2 más caché de catálogo y red.
- `Unidad2_Fase4_Checkout_Concurrencia/`: Fases 1, 2 y 3 más checkout transaccional.
- `Unidad2_Fase4_Checkout_Concurrencia/backend/repository/postgres.go`: pedido, stock y comisiones.
- `Unidad2_Fase4_Checkout_Concurrencia/backend/controllers/dependencies.go`: dependencias comprobables.
- `Unidad2_Fase4_Checkout_Concurrencia/frontend/src/components/Checkout.tsx`: compra conectada a la API.
- `Unidad2_Fase4_Checkout_Concurrencia/README.md`: ejecución y prueba segura de reversión.

## Taller 7

- `Taller7_Registro_Referidos/`: Fase 4 más registro, JWT y red propia.
- `Taller7_Registro_Referidos/backend/database/004_registro_referidos.sql`: migración acumulativa.
- `Taller7_Registro_Referidos/backend/middleware/auth.go`: validación de sesiones.
- `Taller7_Registro_Referidos/backend/controllers/prodController.go`: catálogo público y CRUD administrativo.
- `Taller7_Registro_Referidos/backend/database/005_productos_activos.sql`: migración de desactivación lógica.
- `Taller7_Registro_Referidos/frontend/src/components/Registro.tsx`: alta con código de invitación.
- `Taller7_Registro_Referidos/frontend/src/components/AdminProductos.tsx`: gestión de productos para administradores.
- `Taller7_Registro_Referidos/frontend/src/components/Pedidos.tsx`: historial de pedidos del cliente.
- `Taller7_Registro_Referidos/frontend/src/components/AdminPedidos.tsx`: gestión administrativa de pedidos.
- `Taller7_Registro_Referidos/frontend/src/components/Comisiones.tsx`: comisiones del usuario autenticado.
- `Taller7_Registro_Referidos/frontend/src/components/AdminComisiones.tsx`: libro y liquidación administrativa.
- `Taller7_Registro_Referidos/frontend/src/components/AdminReportes.tsx`: ventas, pedidos, ticket y rankings.
- `Taller7_Registro_Referidos/backend/controllers/reportController.go`: endpoint de reportes administrativos.
- `Taller7_Registro_Referidos/backend/database/006_comisiones_liquidaciones.sql`: estados y liquidaciones de comisiones.

## Comandos puntuales

- Backend actual: entrar en `Taller7_Registro_Referidos/backend` y ejecutar `docker compose up --build -d`.
- Pruebas backend: en esa misma carpeta ejecutar `go test ./...`.
- Frontend actual: entrar en `Taller7_Registro_Referidos/frontend` y ejecutar `npm run dev`.
- Revisión frontend: ejecutar `npm run lint` y `npm run build` solo cuando cambie el frontend.

## Documentación usada para los requisitos

- Disponible fuera del repositorio: `Guia_Practica_Tema5_leccion.md`.
- Disponible fuera del repositorio: `Planificación.docx`, con las fases backend de la Unidad 2.
- No encontrados: `Planificacion_Unidad1_Tema5.md` y `Silabo_Framework_Programacion_Web.md`.
