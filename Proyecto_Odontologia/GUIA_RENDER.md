# Guía de publicación en Render

Esta guía publica la aplicación completa desde la carpeta `Proyecto_Odontologia/`: frontend
React, API Go y PostgreSQL. El archivo `render.yaml` crea y conecta los recursos automáticamente.

## Antes de publicar

1. Fusiona la rama terminada con `main` y comprueba que la carpeta `Proyecto_Odontologia/`
   aparezca en GitHub.
2. Conserva `Proyecto_Odontologia/render.yaml` en el repositorio.
3. No subas contraseñas ni secretos a Git.

## Crear el despliegue

1. Inicia sesión en Render y abre **New > Blueprint**.
2. Conecta GitHub y selecciona el repositorio
   `gualter302/Repositorio_FrameworkProgramacionWEB`.
3. Selecciona la rama `main`.
4. En **Blueprint Path** escribe `Proyecto_Odontologia/render.yaml`.
5. Revisa que Render muestre estos dos recursos:
   - servicio web `sigo-dental`;
   - base PostgreSQL `sigo-dental-db`.
6. Cuando Render solicite los valores que no se sincronizan, configura:
   - `MFA_DEMO_CODE`: `123456` para la demostración académica;
   - `CONSENT_OTP`: `654321` para la demostración académica.
7. Selecciona **Deploy Blueprint** y espera a que ambos recursos estén disponibles.

Los dos códigos anteriores sirven únicamente para una revisión académica. Cámbialos antes de
usar información real. `JWT_SECRET` y `MFA_SECRET` se generan dentro de Render.

## Comprobar el resultado

1. Abre la URL pública que Render asigna a `sigo-dental`.
2. Comprueba `https://TU-SERVICIO.onrender.com/api/health`. Debe responder con estado `ok`.
3. Ingresa con las cinco cuentas indicadas en el README y confirma sus menús.
4. Crea un paciente de prueba y actualiza la página para confirmar que PostgreSQL conserva los
   datos.

La primera apertura de un servicio gratuito puede tardar mientras Render lo reactiva. La base
gratuita tiene 1 GB, expira después de 30 días y no incluye copias de seguridad; es apropiada
para la presentación, pero no para datos clínicos reales.

## Actualizaciones posteriores

Después de fusionar nuevos cambios en `main`, Render vuelve a construir el servicio web. Si
cambias `render.yaml`, abre el Blueprint en Render y sincronízalo para aplicar la configuración.
Los campos con `sync: false` no vuelven a pedir sus valores en una sincronización: modifícalos
desde **Environment** cuando sea necesario.

## Si el despliegue falla

- Confirma que **Blueprint Path** sea `Proyecto_Odontologia/render.yaml`.
- Revisa que el servicio y la base estén en la misma región (`oregon`).
- Abre los registros de compilación del servicio y busca el primer error.
- Comprueba que `DATABASE_URL` esté enlazada a `sigo-dental-db` y no escrita manualmente.
- Si `/api/health` funciona pero la página no abre, vuelve a desplegar el último commit desde el
  panel del servicio.
