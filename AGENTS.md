# Instrucciones de Codex para este repositorio

## Orientación

- Usa `PROJECT_MAP.md` cuando necesites localizar una carpeta o archivo que no conoces.
- Usa `REQUIREMENTS_CHECKLIST.md` cuando la tarea trate sobre el avance, el siguiente hito o un
  requisito del proyecto.
- No leas estos dos documentos por obligación para cambios pequeños y claramente localizados.
- Evita explorar `node_modules`, `dist`, dependencias descargadas y archivos generados.
- Busca primero una ruta o nombre concreto antes de ampliar la búsqueda.

## Cambios

- Lee únicamente los archivos relacionados con la tarea actual.
- Modifica solamente las líneas o funciones necesarias.
- Conserva la organización y el estilo que ya usa cada taller.
- No cambies talleres anteriores cuando el trabajo corresponda a una carpeta nueva.
- Usa una rama independiente para cada taller.
- No agregues coautores a los commits; el autor del repositorio es `gualter302`.

## Seguimiento

- Cuando completes un requisito del proyecto, marca su casilla en `REQUIREMENTS_CHECKLIST.md`.
- Actualiza el Estado general únicamente cuando cambie el hito, la última tarea o el siguiente paso.
- Actualiza `PROJECT_MAP.md` cuando se agregue una carpeta o archivo importante para navegar el
  proyecto; no registres archivos temporales.

## Pruebas y errores

- Ejecuta primero las pruebas directamente relacionadas con el cambio.
- Si una prueba falla, repite esa prueba mientras corriges el problema.
- Ejecuta la suite completa al final cuando el cambio afecte el funcionamiento de la aplicación.
- Si el mismo error continúa después de dos intentos, explica la causa probable antes de cambiar
  más archivos.
- Para comprobar una API o una respuesta JSON, prefiere una petición HTTP directa.

## Uso de herramientas

- En pruebas de navegador, consulta texto visible o selectores concretos y evita pedir todo el HTML.
- Usa capturas de pantalla solamente cuando sea necesario revisar el diseño visual.
- No guardes en memoria código completo, trazas de error, conversaciones ni historiales de trabajo.
- Si guardas información persistente, limita el contenido a puertos, variables importantes y la
  solución final de problemas repetitivos del entorno.

## Entrega

- Explica los cambios con lenguaje natural y directo.
- Entrega el título y la descripción sugeridos para el PR, pero no crees ni publiques el PR.
- No hagas `push` salvo que el usuario lo solicite expresamente.
