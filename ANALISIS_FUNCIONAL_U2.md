# Análisis funcional de Unidad 2

Este documento adapta el análisis entregado al estado real de este repositorio. Las
Fases 1 a 4 ya ofrecen catálogo, PostgreSQL, Redis y checkout transaccional. El sistema
puede vender, pero todavía necesita funciones administrativas para operar diariamente.

## Avance incorporado en el Taller 7

- Registro de clientes con nombre, correo y contraseña protegida.
- Código opcional de referido para que la red pueda crecer.
- JWT real y rutas privadas que usan la identidad de la sesión.
- Red propia y enlace personal de invitación.
- Migración que conserva los datos de las fases anteriores.

## Brechas que continúan abiertas

1. Administración de productos: crear, editar y desactivar productos desde la aplicación.
2. Historial de pedidos: clientes consultan los suyos y administradores gestionan estados. **Implementado en el Taller 7.**
3. Libro de comisiones: conservar cada comisión como movimiento auditable.
4. Liquidaciones: marcar comisiones pendientes como pagadas sin duplicarlas.
5. Reportes: ventas, ticket promedio, productos principales y deuda de comisiones.
6. Operación: paginación y búsqueda desde el servidor para listas grandes.

## Decisiones importantes

- Una escritura debe estar validada y protegida, excepto registro e inicio de sesión.
- Usuario y nodo de referido se crean juntos; si algo falla, no se guarda ninguno.
- El backend toma la identidad desde el JWT y no confía en un correo enviado por el navegador.
- Los pedidos, stock y comisiones de la Fase 4 se conservan como base acumulativa.
- El borrado futuro de productos debe ser lógico para no perder el historial de pedidos.

Lo relacionado con SaaS multicliente, despliegue público, correo real y pasarelas de
pago queda fuera de este taller y corresponde a unidades posteriores.
