# SIGO Dental

Sistema integral de gestión odontológica y expediente clínico digital. La aplicación reúne recepción, atención clínica, planes de tratamiento, caja, inventario, reportes y auditoría en una interfaz responsive.

## Funciones implementadas

- Autenticación con contraseñas bcrypt, JWT de 15 minutos, segundo factor y permisos por rol.
- Ficha única del paciente, anamnesis y alertas clínicas visibles.
- Agenda por profesional y sillón, estados de atención y bloqueo de superposiciones.
- Odontograma FDI de 32 dientes, superficies y catálogo de hallazgos.
- Periodontograma con seis puntos de sondaje, margen, furca, movilidad y sangrado.
- Notas SOAP firmadas, inmutables y corregibles únicamente mediante adendas.
- Consentimientos por procedimiento con firma y código de validación.
- Presupuestos por fases, sesiones, honorarios, laboratorio y descuentos.
- Aprobación de planes, facturación, pagos y cierre ciego de caja.
- Inventario, alertas de stock y descuento automático al completar procedimientos.
- Liquidaciones profesionales, reportes y auditoría de acciones sensibles.
- Cola transaccional preparada para Email, SMS y WhatsApp.
- Persistencia PostgreSQL, contrato OpenAPI y despliegue Docker/Render.

## Ejecución completa con Docker

Requiere Docker Desktop:

```powershell
docker compose up --build
```

Abra `http://localhost:3000`. PostgreSQL y la aplicación se levantan juntos y los datos sobreviven reinicios mediante el volumen `sigo_postgres`.

## Ejecución para desarrollo

Backend:

```powershell
cd backend
go run .
```

Frontend, en otra terminal:

```powershell
cd frontend
npm install
npm run dev
```

Abra `http://localhost:5173`. Sin `DATABASE_URL`, el backend usa memoria para pruebas rápidas. Para trabajar con PostgreSQL, copie las variables de `.env.example` a su entorno.

## Cuentas de demostración

Contraseña común: `Odonto2026!`.

| Perfil | Correo | Segundo factor local |
|---|---|---|
| Administrador | `admin@sigodental.ec` | `123456` |
| Recepción | `recepcion@sigodental.ec` | No requerido |
| Odontólogo | `odontologo@sigodental.ec` | `123456` |
| Auditor | `auditor@sigodental.ec` | `123456` |

El código local de consentimiento es `654321`. En Render ambos códigos se solicitan como secretos y deben cambiarse.

## Pruebas y calidad

```powershell
cd backend
go test ./...

cd ..\frontend
npm run lint
npm run build
```

## Publicación en Render

El repositorio contiene un Blueprint `render.yaml` y un `Dockerfile` de varias etapas. Render construirá React, compilará Go, creará PostgreSQL y conectará `DATABASE_URL` automáticamente.

1. Suba este repositorio a GitHub con la carpeta `Proyecto_Odontologia/` incluida.
2. En Render seleccione **New > Blueprint**, conecte el repositorio y use `Proyecto_Odontologia/render.yaml` como ruta del Blueprint.
3. Confirme el servicio `sigo-dental` y la base `sigo-dental-db`.
4. Cuando Render lo solicite, asigne valores distintos a `MFA_DEMO_CODE` y `CONSENT_OTP`.
5. Espere a que `/api/health` informe `status: ok`.
6. Ingrese con las cuentas de demostración y cambie sus credenciales antes de usar datos reales.

El Blueprint deja PostgreSQL accesible solo por la red privada de Render. `JWT_SECRET` y `MFA_SECRET` se generan automáticamente y no se guardan en Git.

## Integraciones externas pendientes de credenciales

La emisión fiscal real, el cobro bancario y los envíos por WhatsApp/SMS/Email no se pueden activar sin elegir proveedores y registrar cuentas. La aplicación no finge esos envíos: guarda facturas con estado `Pendiente de integración` y mensajes en la bandeja de salida para procesarlos cuando existan las claves correspondientes.

Antes de producción deben definirse:

- Autoridad tributaria y formato fiscal aplicable.
- Proveedor de pagos y credenciales de prueba/producción.
- Meta Cloud API, Twilio u otro gateway de mensajería.
- Política institucional de respaldos, retención y recuperación.
- Dominio, correo transaccional y responsables de privacidad.

## Estructura

- `frontend/`: React, TypeScript y Vite.
- `backend/`: API Go, seguridad, persistencia y reglas de negocio.
- `database/init.sql`: modelo relacional de referencia.
- `docs/openapi.yaml`: contrato OpenAPI 3.
- `Dockerfile`, `docker-compose.yml`, `render.yaml`: ejecución y despliegue.
- `REQUIREMENTS_CHECKLIST.md`: trazabilidad requisito por requisito.
- `PROGRESO.md`: estado, límites y siguiente trabajo externo.

