# Taller 8 — MultiCatálogo PWA

Base acumulativa del Taller 7: registro, catálogo, pedidos, comisiones y administración.
Agrega instalación, interfaz sin conexión, catálogo público guardado durante 24 horas,
carrito local por usuario y checkout con reintentos sin duplicar pedidos.

## Levantar el sistema

Antes de comenzar, abrir Docker Desktop y comprobar que Node.js y npm estén instalados.
Los siguientes comandos parten desde la raíz del repositorio.

En una primera terminal, levantar el backend, PostgreSQL y Redis:

```powershell
cd Taller8_PWA\backend
docker compose -p multicatalogo-pwa up --build -d
docker compose -p multicatalogo-pwa ps
```

En una segunda terminal, instalar las dependencias, compilar la PWA y servirla:

```powershell
cd Taller8_PWA\frontend
npm ci
npm run build
npm run preview -- --host 0.0.0.0 --port 4178
```

Mantener abierta la segunda terminal y visitar http://localhost:4178/catalogo.
El frontend redirige las peticiones `/api` al backend que escucha en el puerto 3008.
La base usa el puerto 5438 y Redis el 6388; el volumen `multicatalogo-pwa-pgdata`
es independiente de los talleres anteriores.

Las cuentas de demostración son:

- Administrador: `admin@upse.edu.ec` / `123456`
- Cliente: `cliente@upse.edu.ec` / `123456`

Para abrir el sistema desde un teléfono conectado a la misma red Wi-Fi, ejecutar `ipconfig`,
buscar la dirección IPv4 del equipo y visitar `http://DIRECCION_IP:4178/catalogo` desde el
teléfono. El Firewall de Windows debe permitir la conexión al puerto 4178. Esta dirección HTTP
sirve para revisar la aplicación; instalar la PWA y usar su modo sin conexión en el teléfono
requiere publicarla con HTTPS.

Si el teléfono no abre la página, ejecutar una vez PowerShell como administrador y permitir el
puerto solamente en redes privadas:

```powershell
New-NetFirewallRule -DisplayName "MultiCatalogo PWA 4178" -Direction Inbound -Protocol TCP -LocalPort 4178 -Action Allow -Profile Private -RemoteAddress LocalSubnet
```

Para detener el sistema, usar `Ctrl+C` en la terminal del frontend y luego ejecutar:

```powershell
cd Taller8_PWA\backend
docker compose -p multicatalogo-pwa down
```

Este último comando detiene los contenedores sin borrar la información de PostgreSQL.
Para desarrollo se puede usar `npm run dev`; el service worker debe comprobarse con
`npm run build` y `npm run preview`.

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

## Publicar en Render

La configuración `render.yaml` crea un servicio web, PostgreSQL y una instancia Key Value usada
como caché. El archivo `Dockerfile.render` compila React y Go en una sola imagen. La página y la
API comparten el mismo dominio, por lo que no hay que copiar direcciones en `VITE_API_URL` ni
configurar CORS entre dos servicios.

Antes de comenzar, fusionar esta rama en `main` y subirla a GitHub. Después:

1. Crear una cuenta en https://render.com y conectar la cuenta de GitHub.
2. En el panel de Render elegir **New → Blueprint** y conectar el repositorio
   `gualter302/Repositorio_FrameworkProgramacionWEB`.
3. Seleccionar la rama `main` y escribir `Taller8_PWA/render.yaml` en **Blueprint Path**.
4. Render solicitará `ADMIN_PASSWORD`. Escribir una contraseña privada de al menos 12 caracteres;
   reemplazará `123456` únicamente para la cuenta `admin@upse.edu.ec` publicada.
5. Revisar los tres recursos y elegir **Apply**. El primer despliegue puede tardar varios minutos
   porque compila el frontend, el backend y crea la base.
6. Abrir el servicio `multicatalogo-pwa-gualter302`. Primero comprobar `/api/health` y después
   abrir `/login`. Las tablas, migraciones y datos iniciales se crean automáticamente.
7. Iniciar sesión con `admin@upse.edu.ec` y la contraseña elegida en el paso 4. La tarjeta con
   `123456` solo aparece en compilaciones locales y queda oculta en Render.
8. Probar catálogo, registro, carrito, pedido y administración. En un teléfono abrir la dirección
   HTTPS de Render, visitar el catálogo una vez e instalar la PWA desde el navegador.

Los siguientes cambios enviados a `main` se despliegan automáticamente. Si falla una nueva
versión, revisar **Events** y **Logs** del servicio. El control de salud es `/api/health` y Render
no reemplaza la versión anterior hasta que la nueva responde correctamente.

### Límites del plan gratuito

- El servicio web se suspende después de 15 minutos sin visitas; la primera apertura puede tardar
  cerca de un minuto.
- PostgreSQL gratuito tiene 1 GB, no incluye copias de seguridad y vence 30 días después de crearse.
- Key Value gratuito puede perder su contenido al reiniciarse. En este proyecto no se pierden datos
  importantes porque solo contiene caché; PostgreSQL sigue siendo la fuente principal.
- Para mantener datos reales se debe actualizar PostgreSQL a un plan de pago antes del vencimiento
  o exportar la información y crear otra base.

### Alternativa para un servidor propio

El `compose.production.yml` y Caddy continúan disponibles para un VPS con dominio propio. Desde
`Taller8_PWA`, definir `DOMAIN`, `DB_PASSWORD` y `JWT_SECRET` y ejecutar:

```powershell
docker compose -f compose.production.yml up --build -d
```

En esa alternativa el dominio debe apuntar al servidor y permitir los puertos 80 y 443. Para
regresar de versión se despliega la imagen anterior conservando PostgreSQL; nunca se eliminan los
volúmenes para actualizar.

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
