# Lista de requisitos y avance

## Estado general

- **Hito actual:** Taller 7, registro y red de referidos terminado y verificado.
- **Última tarea completada:** CRUD administrativo de productos con permisos y desactivación lógica.
- **Siguiente tarea pendiente:** libro de comisiones y liquidaciones.
- **Base recomendada:** `Taller7_Registro_Referidos/`.
- **Fuentes disponibles:** guía práctica, planificación, análisis funcional y README de cada taller.
- **Fuentes faltantes:** `Planificacion_Unidad1_Tema5.md` y `Silabo_Framework_Programacion_Web.md`.

## Hitos completados

- [x] Crear el entorno React con TypeScript, Vite y Tailwind CSS.
- [x] Consumir PokeAPI y mostrar los primeros 151 Pokémon.
- [x] Agregar filtros y navegación adaptable a la Pokédex.
- [x] Crear el sidebar colapsable y su versión para celular.
- [x] Conectar el frontend React con el backend Go y Fiber.
- [x] Centralizar las llamadas del frontend en `src/services/api.ts`.
- [x] Consumir `POST /api/login` y `GET /api/productos` desde React.

## Taller 5 Fase 1 Backend robusto

- [x] Separar rutas, controladores, modelos y repositorio.
- [x] Crear una respuesta de error común con estado, mensaje y detalles opcionales.
- [x] Validar que correo y contraseña no estén vacíos.
- [x] Responder 401 cuando las credenciales sean incorrectas.
- [x] Mantener el formato que el frontend ya utiliza para login y productos.
- [x] Agregar `GET /api/productos/:id`.
- [x] Responder 400 cuando el ID del producto no sea numérico o sea menor que uno.
- [x] Responder 404 cuando el producto solicitado no exista.
- [x] Agregar `GET /api/health`.
- [x] Crear pruebas para health, login y consulta de productos.
- [x] Crear y robustecer `GET /api/red` cuando se implemente la red multinivel.

## Taller 6 Ahorro de tokens

- [x] Crear `PROJECT_MAP.md` con la estructura y archivos importantes.
- [x] Mantener `PROJECT_MAP.md` dentro del límite de 100 líneas.
- [x] Crear `REQUIREMENTS_CHECKLIST.md` con el estado real del proyecto.
- [x] Crear `AGENTS.md` para limitar lecturas y pruebas innecesarias en Codex.
- [x] Agregar reglas de uso puntual para Playwright y memoria.
- [x] Documentar la configuración MCP opcional en `CODEX_SETUP.md`.
- [x] Evitar código fuente y trazas completas dentro de los archivos de orientación.

## Frontend de Tema 5

- [x] Hacer que el login distinga los roles `admin` y `cliente`.
- [x] Devolver el rol desde `POST /api/login` sin romper el formato actual.
- [x] Guardar correo, rol y token en `AuthContext`.
- [ ] Proteger Dashboard y Mi Red para que solo acceda el administrador.
- [ ] Crear la tienda a pantalla completa en `/tienda`.
- [ ] Crear datos de productos con descripción, categoría y galería.
- [ ] Separar la obtención de productos en un servicio reutilizable.
- [ ] Agregar búsqueda por texto y filtro por categoría en `/catalogo`.
- [ ] Reflejar la categoría seleccionada en la dirección de la página.
- [ ] Crear `/producto/:id` para mostrar el detalle de un producto.
- [ ] Agregar galería y lightbox con cierre por botón y tecla Escape.
- [ ] Mostrar un mensaje claro cuando el producto no exista.
- [ ] Guardar el carrito en `localStorage` usando una clave por usuario.
- [ ] Permitir aumentar, disminuir y eliminar cantidades del carrito.
- [ ] Mantener separados los carritos del administrador y del cliente.
- [x] Crear el flujo carrito, checkout y confirmación.
- [x] Vaciar el carrito únicamente después de confirmar el pedido.
- [x] Crear la estructura de la red de referidos con tres niveles.
- [ ] Calcular comisiones de 10 %, 5 % y 2 % según el nivel.
- [x] Mostrar el árbol de referidos y sus totales en `/mi-red`.
- [x] Calcular ventas, referidos y comisiones en el Dashboard.
- [ ] Adaptar Sidebar y Navbar al rol autenticado.
- [ ] Mantener la interfaz adaptable y accesible en computadora y celular.

## Unidad 2 Fase 2 Persistencia con PostgreSQL

- [x] Reemplazar usuarios, productos y red guardados en memoria.
- [x] Diseñar las tablas y relaciones de PostgreSQL.
- [x] Crear una conexión segura entre Go y PostgreSQL.
- [x] Mover las consultas de datos a la capa `repository`.
- [x] Confirmar que los datos se mantienen después de reiniciar el backend.
- [x] Mostrar en la tienda los productos agregados directamente a la base de datos.

## Unidad 2 Fase 3 Caché con Redis

- [x] Agregar Redis al entorno del proyecto.
- [x] Guardar temporalmente el catálogo para evitar consultas repetidas.
- [x] Guardar temporalmente el resultado de la red multinivel.
- [x] Definir cuándo se renueva o elimina cada dato guardado en caché.
- [x] Comparar la primera consulta (`MISS`) con las siguientes (`HIT`).

## Unidad 2 Fase 4 Checkout y concurrencia

- [x] Crear `POST /api/checkout` usando dependencias separadas y comprobables.
- [x] Procesar cada compra dentro de una transacción de base de datos.
- [x] Descontar existencias y registrar el pedido como una sola operación.
- [x] Deshacer todos los cambios cuando el pago o el stock fallen.
- [x] Calcular las comisiones de la red después de una compra correcta.
- [x] Usar concurrencia de Go únicamente donde no comprometa la transacción.
- [x] Actualizar el KPI de comisiones después de confirmar el pedido.
- [x] Mostrar al usuario un error claro cuando el checkout no pueda completarse.

## Taller 7 Registro y referidos

- [x] Crear una cuenta con nombre, correo y contraseña protegida.
- [x] Crear usuario y nodo de referido dentro de una sola transacción.
- [x] Leer un código de invitación desde `/registro?ref=CODIGO`.
- [x] Responder 409 para correos repetidos y 400 para códigos inexistentes.
- [x] Reemplazar el token temporal por un JWT con vencimiento.
- [x] Proteger checkout y red usando la identidad del JWT.
- [x] Mostrar a cada usuario su red y su enlace de invitación.
- [x] Agregar una migración para conservar bases de fases anteriores.
- [x] Mantener una sola copia de cada documento de orientación.

## Administración de productos

- [x] Crear, editar y desactivar productos con permisos de administrador.
- [x] Mostrar productos activos en el catálogo e incluir inactivos en el panel admin.
- [x] Renovar `productos:all` después de cada escritura.
- [x] Proteger la pantalla administrativa según el rol autenticado.

## Consulta y gestión de pedidos

- [x] Permitir que cada cliente consulte solamente sus pedidos.
- [x] Crear listado administrativo con filtros y paginación.
- [x] Gestionar estados y devolver stock al cancelar dentro de una transacción.
- [x] Revertir las comisiones asociadas cuando se cancela un pedido.

## Verificación y entrega

- [x] Ejecutar `go test ./...` en el backend de la Fase 1.
- [x] Ejecutar `npm run lint` después del cambio del frontend de la Fase 4.
- [x] Ejecutar `npm run build` después del cambio del frontend de la Fase 4.
- [ ] Probar manualmente los dos roles cuando sean implementados.
- [x] Probar carrito, checkout correcto y reversión ante un producto inexistente.
- [x] Probar el árbol multinivel y los valores actualizados de comisiones.
- [x] Probar registro correcto, duplicado, referencia inválida y ruta sin token.
- [x] Verificar que el checkout use el usuario del JWT aunque el cuerpo envíe otro correo.
- [x] Probar CRUD de productos, rechazo del cliente y desactivación lógica.
- [x] Probar historial del cliente, listado administrativo y cancelación con devolución de stock.
- [ ] Verificar navegación por teclado, textos alternativos y contraste.
- [ ] Actualizar este checklist inmediatamente después de cada prueba exitosa.
