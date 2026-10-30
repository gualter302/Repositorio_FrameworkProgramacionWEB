# Mapa compacto del proyecto

## Orientación

- `README.md`: instalación, cuentas, pruebas y publicación.
- `PROGRESO.md`: avance por fase y dependencias externas.
- `REQUIREMENTS_CHECKLIST.md`: trazabilidad del SRS.

## Frontend

- `frontend/src/App.tsx`: sesión, navegación, pacientes y agenda.
- `frontend/src/api.ts`: contratos y único acceso a la API.
- `frontend/src/modules/ClinicalPage.tsx`: odontograma, periodoncia, SOAP, consentimientos y planes.
- `frontend/src/modules/FinancePage.tsx`: presupuestos, facturas y caja.
- `frontend/src/modules/InventoryPage.tsx`: existencias y ajustes.
- `frontend/src/modules/ReportsPage.tsx`: indicadores y auditoría.
- `frontend/src/styles.css`: diseño responsive.

## Backend

- `backend/main.go`: inicio del servidor.
- `backend/server.go`: rutas base, seguridad, pacientes y agenda.
- `backend/auth.go`: firma JWT y segundo factor.
- `backend/store.go`: memoria para pruebas y PostgreSQL para despliegue.
- `backend/handlers_clinical.go`: expediente y planes.
- `backend/handlers_finance.go`: inventario, facturación, caja, honorarios y reportes.
- `backend/models.go`: contratos persistidos.
- `backend/server_test.go`: pruebas de reglas críticas.

## Infraestructura y documentación

- `database/init.sql`: esquema relacional de referencia.
- `docs/openapi.yaml`: API OpenAPI 3.
- `Dockerfile`: imagen única de frontend y backend.
- `docker-compose.yml`: entorno local con PostgreSQL.
- `render.yaml`: Blueprint para Render.
- `.env.example`: variables locales.

