# Estado del proyecto

Última actualización: 30 de octubre de 2026.

## Fase 1 - Core administrativo y turnos

- [x] Pacientes, filiación, anamnesis y alertas.
- [x] Agenda diaria, profesionales, sillones y estados de cita.
- [x] Bloqueo de superposición por sillón.
- [x] Dashboard y navegación por roles.

## Fase 2 - Core clínico

- [x] Odontograma FDI por dientes, superficies y condiciones.
- [x] Periodontograma de seis puntos.
- [x] Notas SOAP firmadas e inmutables con adendas.
- [x] Consentimientos firmados y validados con código.
- [x] Planes por fases, sesiones y costos.

## Fase 3 - Finanzas, inventario y reportes

- [x] Aprobación total de presupuestos y transición a facturación.
- [x] Comprobantes locales y formas de pago.
- [x] Cierre ciego de caja.
- [x] Inventario, mínimos y descarga automática de insumos.
- [x] Liquidaciones profesionales.
- [x] Reporte operativo/financiero y auditoría.
- [x] Bandeja de salida para notificaciones multicanal.

## Fase 4 - Seguridad, pruebas y despliegue

- [x] bcrypt, JWT, RBAC y segundo factor.
- [x] Expiración/cierre por 15 minutos.
- [x] Cabeceras de seguridad y CORS configurable.
- [x] Persistencia PostgreSQL.
- [x] Pruebas automatizadas de reglas críticas.
- [x] OpenAPI 3.
- [x] Docker, Docker Compose y Blueprint de Render.
- [x] Health check y frontend servido por el mismo proceso.

## Límites que requieren decisiones o cuentas externas

- [ ] Timbrado fiscal real: falta definir autoridad, ambiente y credenciales.
- [ ] Pasarela bancaria real: falta escoger proveedor y registrar comercio.
- [ ] WhatsApp/SMS/Email real: falta proveedor, remitente y plantillas aprobadas.
- [ ] Firma biométrica certificada: actualmente se registra firma electrónica simple y OTP.
- [ ] PACS/DICOM: fuera del alcance de la fase solicitada por el SRS.
- [ ] UAT formal con usuarios de una clínica y actas de aceptación.
- [ ] Escaneo DAST/SAST en la plataforma CI elegida.
- [ ] Prueba de carga externa de 150 usuarios y validación del SLA en infraestructura contratada.

## Siguiente paso recomendado

1. Revisar la aplicación local con los cuatro perfiles.
2. Crear el repositorio remoto y subir la rama principal cuando el propietario lo autorice.
3. Importar `render.yaml` como Blueprint.
4. Realizar UAT y registrar incidencias reales.
5. Elegir proveedores externos antes de activar datos clínicos reales.

