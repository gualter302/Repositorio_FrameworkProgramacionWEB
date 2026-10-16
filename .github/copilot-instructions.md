# Instrucciones del repositorio

## Alcance

- El trabajo vigente se realiza principalmente en `Taller7_Registro_Referidos/`.
- Los talleres anteriores son referencias históricas; no modificarlos salvo que la tarea lo pida.
- Consultar `PROJECT_MAP.md` para localizar archivos y `REQUIREMENTS_CHECKLIST.md` para conocer el avance.

## Cambios

- Leer solamente los archivos relacionados con la tarea.
- Evitar `node_modules/`, `dist/`, dependencias descargadas y archivos generados.
- Mantener la separación existente entre frontend React, backend Go, PostgreSQL y Redis.
- No guardar secretos, tokens ni credenciales reales.
- No hacer `push`, crear PR ni modificar ramas remotas sin solicitud expresa.

## Verificación

- Backend: ejecutar `go test ./...` desde `Taller7_Registro_Referidos/backend`.
- Frontend: ejecutar `npm run lint` y `npm run build` desde `Taller7_Registro_Referidos/frontend` cuando corresponda.
- Para API, preferir peticiones HTTP directas y comprobar los códigos de respuesta.
- Actualizar el checklist y el mapa del proyecto cuando se complete una función o se agregue un archivo importante.
