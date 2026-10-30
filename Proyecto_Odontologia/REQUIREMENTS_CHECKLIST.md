# Trazabilidad de requisitos

## Requisitos funcionales

- [x] RF-01 Registro unificado del paciente.
- [x] RF-02 Anamnesis y alertas médicas.
- [x] RF-03 Odontograma gráfico interactivo.
- [x] RF-04 Periodontograma digital.
- [x] RF-05 Evolución y notas SOAP inmutables.
- [x] RF-06 Consentimiento informado digital.
- [x] RF-07 Presupuestos multinivel.
- [x] RF-08 Aprobación y transición de ítems a estado agendable.
- [x] RF-09 Agendamiento multirrecinto por profesional y sillón.
- [x] RF-10 Estados de cita.
- [~] RF-11 Notificaciones automáticas: cola implementada; envío real pendiente de proveedor.
- [~] RF-12 Punto de venta y facturación: comprobante local implementado; timbrado fiscal pendiente.
- [x] RF-13 Liquidación de honorarios.
- [x] RF-14 Descargo automático de insumos.

## Requisitos no funcionales

- [~] RNF-01 Tiempo de respuesta: aplicación optimizada; requiere medición en Render.
- [~] RNF-02 Disponibilidad: health check configurado; el SLA depende del plan contratado.
- [~] RNF-03 Concurrencia: bloqueos internos seguros; falta prueba externa de 150 usuarios.
- [~] RNF-04 TLS en tránsito mediante Render; el cifrado en reposo debe confirmarse con el plan final contratado.
- [x] RNF-05 RBAC, JWT, MFA y cierre a los 15 minutos.
- [x] RNF-06 Auditoría de operaciones sensibles con usuario, fecha e IP.
- [x] RNF-07 Diseño responsive para escritorio, laptop y tablet.
- [~] RNF-08 Navegadores: compilación estándar; falta matriz manual en cuatro navegadores.

## Reglas de negocio

- [x] RN-01 Historial clínico sin borrado y correcciones mediante adendas.
- [x] RN-02 Cierre ciego de caja.
- [x] RN-03 Prohibición de superponer citas en el mismo sillón.

`[~]` indica que la parte controlable por el código está lista, pero la aceptación final depende de infraestructura, proveedor o prueba externa.
