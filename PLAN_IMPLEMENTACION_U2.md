# Plan de implementación de Unidad 2

Este plan conserva las ideas nuevas de los archivos entregados, pero refleja solamente
funciones comprobadas en este repositorio. Una casilla se marca cuando sus pruebas pasan.

## Taller 7 — Registro y referidos

- [x] Agregar nombre, estado, fecha y vínculo de referido al usuario.
- [x] Agregar un código único a cada nodo de la red.
- [x] Crear una migración idempotente para una base existente.
- [x] Crear `POST /api/register` con usuario y referido en la misma transacción.
- [x] Proteger contraseñas con bcrypt.
- [x] Rechazar correo duplicado con 409 y referencia inexistente con 400.
- [x] Entregar JWT con correo, rol e identificador del usuario.
- [x] Proteger checkout y Mi Red con la sesión.
- [x] Crear `/registro` y leer el código desde `?ref=`.
- [x] Mostrar y copiar el enlace personal de invitación.
- [x] Invalidar la caché de redes cuando cambia la estructura o las comisiones.

## Siguientes talleres sugeridos

### Administración de productos

- [ ] Crear, editar y desactivar productos con permisos de administrador.
- [ ] Mostrar productos activos en el catálogo e incluir inactivos en el panel admin.
- [ ] Renovar `productos:all` después de cada escritura.

### Consulta y gestión de pedidos

- [x] Permitir que cada cliente consulte solamente sus pedidos.
- [x] Crear listado administrativo con filtros y paginación.
- [x] Gestionar estados y devolver stock al cancelar dentro de una transacción.

### Comisiones y liquidaciones

- [ ] Convertir cada comisión en un movimiento auditable del pedido.
- [ ] Mostrar comisiones propias con sus estados.
- [ ] Crear liquidaciones sin permitir pagos duplicados.

### Reportes y operación

- [ ] Calcular ventas, pedidos, ticket promedio y comisiones desde PostgreSQL.
- [ ] Mostrar productos y referidos principales en el Dashboard.
- [ ] Agregar paginación y búsqueda desde el servidor.
