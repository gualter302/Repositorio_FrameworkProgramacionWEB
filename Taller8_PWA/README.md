# Taller 8 — MultiCatálogo PWA

Base acumulativa del Taller 7: registro, catálogo, pedidos, comisiones y administración.
Agrega instalación, interfaz sin conexión, catálogo público guardado durante 24 horas,
carrito local por usuario y checkout con reintentos sin duplicar pedidos.

## Ejecutar

En backend:

```powershell
docker compose -p multicatalogo-pwa up --build -d
```

En frontend:

```powershell
npm ci
npm run build
npm run preview -- --host localhost --port 4178
```

Abrir http://localhost:4178/catalogo. Vite redirige /api al backend en 3008.
La base usa 5438 y Redis 6388; el volumen multicatalogo-pwa-pgdata es independiente.
Las cuentas de demostración conservan admin@upse.edu.ec y cliente@upse.edu.ec / 123456.
Para desarrollo usar npm run dev; comprobar el service worker con build y preview.

## Uso sin conexión

Visitar el catálogo una vez con conexión. Después se puede abrir la interfaz, consultar
esa copia pública y preparar el carrito. La última fecha de consulta aparece en el catálogo.
La copia local se limita a 500 productos y un millón de caracteres.
Las imágenes consultadas se guardan hasta siete días (máximo 40) y hay imagen de respaldo.
Sin una visita previa o con catálogo de más de 24 horas se muestra un error de conexión.

Login, registro, compras, pedidos privados, red, reportes y liquidaciones necesitan internet.
El service worker nunca almacena respuestas de API ni reenvía escrituras automáticamente.
Redis continúa siendo una caché del servidor; el almacenamiento sin conexión vive en el navegador.

El carrito se conserva por correo al cerrar sesión. Una cuenta distinta no recibe ese carrito.
La sesión vence a las dos horas; al entrar nuevamente se recupera el carrito correspondiente.
Los datos locales pertenecen al perfil del navegador: en un equipo compartido conviene borrar
los datos del sitio al terminar. No se cifra el almacenamiento local.

## Compra y respuesta perdida

Antes del primer envío se consulta catálogo actual y se comprueba precio y disponibilidad.
El servidor vuelve a validar stock bajo bloqueo y calcula el total.
Cada intento guarda una clave y el cuerpo original por usuario. Si la conexión se pierde,
confirmar otra vez recupera esa misma compra; también puede consultarse Mis pedidos.
El reintento conserva los datos originales aunque se haya editado el formulario.
Una respuesta definitiva 400/409 libera el intento para corregirlo. Un error de red o servidor
conserva la clave. La clave, huella y respuesta quedan en PostgreSQL dentro de la transacción.
Reutilizar la clave con datos distintos responde 409.

Las comisiones siguen calculándose después del commit como en el Taller 7.
Un cierre del proceso en ese intervalo sigue siendo una limitación heredada; no hay worker persistente.

## Instalar y actualizar

En Edge/Chrome usar Instalar aplicación cuando aparezca o el menú del navegador.
En iPhone usar Safari → Compartir → Añadir a pantalla de inicio.
La nueva versión se anuncia y solo se aplica por decisión del usuario desde catálogo o carrito.
El carrito se guarda antes de actualizar. No se fuerza una recarga durante checkout o administración.

## HTTPS y producción

Hay Dockerfile de frontend y configuración Caddy con proxy /api y soporte para rutas React.
Desde esta carpeta, definir DOMAIN, DB_PASSWORD y JWT_SECRET (al menos 32 caracteres) y ejecutar:

```powershell
docker compose -f compose.production.yml up --build -d
```

El dominio debe apuntar a un servidor con puertos 80/443 accesibles para que Caddy obtenga
un certificado. No se publicó en internet: falta proporcionar dominio y alojamiento.
La instalación en iOS/Android reales y HTTPS público quedan pendientes de esa publicación.
El frontend usa /api; nunca incluir secretos en variables VITE_.
Antes de producción cambiar las contraseñas de demostración y hacer copia de PostgreSQL.

En una base existente aplicar database/007_pwa.sql manualmente: los scripts de inicio solo
se ejecutan al crear el volumen. No borrar volúmenes para actualizar.
Para regresar de versión desplegar la imagen anterior, conservar PostgreSQL y recompilar
el frontend anterior con un nuevo service worker; el carrito usa una clave estable.

## Verificación

- Backend: go test ./...
- Compra real: ejecutar backend/test-idempotencia.ps1 con los contenedores activos.
- Frontend: npm run lint, npm run build y npm run test:pwa (Microsoft Edge instalado).
- Iconos: npm run icons, generados desde public/logo.svg.

Las pruebas de navegador usan API simulada y la compilación real con service worker:
reapertura sin conexión, carrito persistente, manifest, iconos, cachés sin respuestas privadas,
cambio de cuenta, sesión vencida y ancho de 375 px. Otra prueba simula dos versiones reales
del service worker y verifica el aviso, el bloqueo en checkout y la conservación del carrito.
El script de compra usa PostgreSQL real.
Quedan por verificar actualización entre dos despliegues públicos, instalación física Android/iOS,
auditoría completa de accesibilidad y condiciones de red lenta/caída en dispositivos.
