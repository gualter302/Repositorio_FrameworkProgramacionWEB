# Lista de requisitos y avance

## Estado general

- **Hito actual:** Taller 8 PWA publicado en Render y en validación desde dispositivos móviles.
- **Última tarea completada:** catálogo con fotografías, cuentas de acceso rápido y reportes adaptables sin desplazamiento innecesario.
- **Siguiente tarea pendiente:** desplegar los cambios y confirmar la reapertura sin conexión en Android/iOS.
- **Base recomendada:** `Taller8_PWA/`.
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
- [x] Adaptar Sidebar y Navbar al rol autenticado.
- [x] Mantener la interfaz adaptable y accesible en computadora y celular.

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

## Comisiones y liquidaciones

- [x] Convertir cada comisión en un movimiento auditable del pedido.
- [x] Mostrar comisiones propias con sus estados.
- [x] Crear liquidaciones sin permitir pagos duplicados.

## Reportes y operación

- [x] Calcular ventas, pedidos, ticket promedio y comisiones desde PostgreSQL.
- [x] Mostrar productos y referidos principales en el Dashboard administrativo.
- [x] Mostrar comisiones pendientes de pago.

## Plan de conversión a PWA

Objetivo: convertir MultiCatálogo en una aplicación instalable, con interfaz y catálogo
consultado disponibles sin conexión. Las compras y operaciones administrativas requieren
confirmación del servidor. Base: `Taller7_Registro_Referidos/`, conservando sus funciones.

Las fases se ejecutan en orden y se marcan únicamente después de superar su verificación.
Implementación acumulativa en `Taller8_PWA/`, en la rama `plan-pwa` por petición del usuario.
Las casillas pendientes corresponden a verificaciones de despliegue, dispositivos y actualización.

### PWA Fase 1 — Preparación y alcance

- [x] Crear una carpeta acumulativa para el taller PWA desde la última versión funcional y registrarla en `PROJECT_MAP.md`.
- [ ] Comprobar login, catálogo, carrito, checkout y permisos antes de introducir el service worker.
- [x] Definir qué funciona sin conexión: interfaz, catálogo previamente consultado y carrito local.
- [x] Definir como operaciones en línea: login, registro, checkout, cambios de productos/pedidos y liquidaciones.
- [x] Preparar `VITE_API_URL` para desarrollo y producción, evitando `localhost` en la versión publicada.
- [ ] Revisar navegación móvil y formularios a 360–375 px, teclado y botones táctiles.
- [ ] **Verificación:** base compila, flujos principales pasan y el alcance sin conexión queda documentado en el README del taller.

### PWA Fase 2 — Identidad e instalación

- [x] Integrar la generación del manifest y service worker en la compilación de Vite, con versiones compatibles de las dependencias.
- [x] Definir `id`, nombre, nombre corto, idioma español, `start_url`, `scope`, colores y apertura `standalone`.
- [x] Crear iconos de 192 y 512 px, una variante maskable y un icono para dispositivos Apple.
- [x] Actualizar título, descripción, idioma HTML y color del navegador con la identidad de MultiCatálogo.
- [x] Añadir una acción de instalación cuando el navegador la permita e instrucciones para instalación manual en iOS.
- [ ] **Verificación:** manifest sin errores, iconos visibles y apertura de la aplicación instalada en modo independiente en navegadores compatibles.

### PWA Fase 3 — Interfaz y catálogo sin conexión

- [x] Guardar los archivos versionados de la interfaz y ofrecer una respuesta de navegación para las rutas de React.
- [x] Guardar solamente consultas públicas de catálogo: intentar obtener datos actuales y usar la última copia cuando falle la conexión.
- [x] Limitar duración y tamaño de la caché de catálogo e imágenes; mostrar una imagen de respaldo cuando falte una imagen externa.
- [x] Mostrar aviso de desconexión, fecha de la última consulta y aclarar que precio y existencias se validan al comprar.
- [x] Excluir del almacenamiento del service worker las respuestas de sesión, red privada, pedidos, comisiones, reportes y toda escritura.
- [x] Ofrecer un estado claro si la aplicación todavía no tiene un catálogo guardado.
- [ ] **Verificación:** tras una visita en línea, cerrar y abrir sin conexión permite navegar por la interfaz y el catálogo guardado; las rutas privadas no revelan respuestas cacheadas.

### PWA Fase 4 — Carrito y sesión

- [x] Persistir el carrito por usuario en almacenamiento local y validar su estructura al recuperarlo.
- [x] Mantener separados los carritos al cambiar de cuenta y definir su conservación al cerrar sesión.
- [x] Permitir preparar el carrito sin conexión y conservarlo al reiniciar la aplicación.
- [x] Gestionar token vencido y respuestas 401: solicitar inicio de sesión sin perder el carrito del usuario correspondiente.
- [x] Bloquear la confirmación sin conexión y conservar el carrito ante errores del servidor.
- [x] Revalidar precio, disponibilidad y stock al volver a conectar, mostrando los cambios antes de confirmar.
- [ ] **Verificación:** recarga, cierre, desconexión y cambio de cuenta conservan el carrito correcto sin mezclar información entre usuarios.

### PWA Fase 5 — Compras y actualizaciones seguras

- [x] Mantener compras y liquidaciones fuera de cualquier cola de reenvío automático del service worker.
- [x] Evitar envíos repetidos mientras una operación está en curso y tratar una respuesta perdida como resultado incierto.
- [x] Implementar una clave de idempotencia para checkout, con restricción persistente en PostgreSQL, para que reintentar la misma compra no duplique pedidos ni descuentos de stock.
- [x] Avisar cuando exista una nueva versión y permitir actualizar después de guardar el estado local y terminar operaciones en curso.
- [x] Eliminar cachés antiguas sin borrar el carrito ni forzar recargas durante un formulario o una compra.
- [x] **Verificación:** reintento contra PostgreSQL crea un solo pedido; dos versiones simuladas del service worker conservan el carrito y bloquean actualizar en checkout.

### PWA Fase 6 — Despliegue HTTPS

- [ ] Preparar la compilación de producción y el alojamiento del frontend y la API bajo HTTPS.
- [ ] Configurar rutas de React, acceso a `/api`, CORS cuando corresponda y ausencia de contenido HTTP dentro del sitio HTTPS.
- [x] Configurar el alcance y los encabezados de caché del service worker, manifest y archivos versionados para que se detecten actualizaciones.
- [x] Mantener secretos del backend fuera del paquete frontend y documentar variables, arranque y recuperación de una versión anterior.
- [ ] **Verificación:** desde un teléfono real la versión publicada conecta con la API, se instala y puede abrir su interfaz sin conexión después de la primera visita.

### PWA Fase 7 — Validación y entrega

- [x] Ejecutar lint, compilación y pruebas del backend si fue modificado; revisar manifest, almacenamiento y service worker en las herramientas del navegador.
- [ ] Probar en Chrome/Android, Safari/iOS y un navegador de escritorio, documentando diferencias de instalación.
- [ ] Probar primera visita sin datos, reapertura sin conexión, red lenta, caída de API, token vencido y cambio de usuario.
- [ ] Verificar compra correcta, stock insuficiente, reintento sin duplicados y actualización de la aplicación instalada.
- [ ] Revisar accesibilidad: foco, etiquetas, contraste, textos alternativos y avisos de conexión comprensibles.
- [x] Documentar instalación, límites sin conexión, actualización y resultados; actualizar el mapa y este checklist.
- [ ] **Verificación:** todos los escenarios anteriores pasan con evidencias y las limitaciones conocidas quedan registradas antes de entregar.

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
- [x] Probar consulta de comisiones, liquidación y rechazo de duplicados.
- [x] Probar reporte general y rechazo al acceso de un cliente.
- [ ] Verificar navegación por teclado, textos alternativos y contraste.
- [ ] Actualizar este checklist inmediatamente después de cada prueba exitosa.
