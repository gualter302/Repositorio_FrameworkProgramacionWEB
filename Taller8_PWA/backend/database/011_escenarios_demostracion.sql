-- Completa los escenarios de prueba del administrador y del cliente sin duplicarlos.
CREATE TABLE IF NOT EXISTS demo_seed_history (
    version TEXT PRIMARY KEY,
    aplicada_en TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DO $$
DECLARE
    cliente_ref_id INT;
    sofia_ref_id INT;
    mateo_ref_id INT;
    valentina_ref_id INT;
BEGIN
    IF EXISTS (SELECT 1 FROM demo_seed_history WHERE version = '011_escenarios_demostracion') THEN
        RETURN;
    END IF;

    SELECT id INTO cliente_ref_id FROM referidos WHERE codigo = 'DEMO';
    IF cliente_ref_id IS NULL THEN
        RAISE EXCEPTION 'No existe el referido DEMO requerido para preparar los datos de prueba';
    END IF;

    INSERT INTO referidos (nombre, codigo, nivel, ventas, comisiones, parent_id)
    SELECT 'Sofía Rivera', 'DEMO-SOFIA', nivel + 1, 0, 0, id
    FROM referidos WHERE id = cliente_ref_id
    ON CONFLICT (codigo) DO UPDATE SET nombre = EXCLUDED.nombre, parent_id = EXCLUDED.parent_id
    RETURNING id INTO sofia_ref_id;

    INSERT INTO referidos (nombre, codigo, nivel, ventas, comisiones, parent_id)
    SELECT 'Mateo Salazar', 'DEMO-MATEO', nivel + 1, 0, 0, id
    FROM referidos WHERE id = cliente_ref_id
    ON CONFLICT (codigo) DO UPDATE SET nombre = EXCLUDED.nombre, parent_id = EXCLUDED.parent_id
    RETURNING id INTO mateo_ref_id;

    INSERT INTO referidos (nombre, codigo, nivel, ventas, comisiones, parent_id)
    SELECT 'Valentina Cruz', 'DEMO-VALE', nivel + 1, 0, 0, id
    FROM referidos WHERE id = sofia_ref_id
    ON CONFLICT (codigo) DO UPDATE SET nombre = EXCLUDED.nombre, parent_id = EXCLUDED.parent_id
    RETURNING id INTO valentina_ref_id;

    INSERT INTO usuarios (nombre, email, password, rol, activo, referido_id) VALUES
        ('Sofía Rivera', 'sofia.demo@multicatalogo.local', crypt('123456', gen_salt('bf')), 'cliente', true, sofia_ref_id),
        ('Mateo Salazar', 'mateo.demo@multicatalogo.local', crypt('123456', gen_salt('bf')), 'cliente', true, mateo_ref_id),
        ('Valentina Cruz', 'valentina.demo@multicatalogo.local', crypt('123456', gen_salt('bf')), 'cliente', true, valentina_ref_id)
    ON CONFLICT (email) DO UPDATE SET activo = true, referido_id = EXCLUDED.referido_id;

    INSERT INTO pedidos (
        numero, usuario_id, cliente_nombre, cliente_email, direccion, ciudad,
        telefono, metodo_pago, estado, total, fecha
    )
    SELECT datos.numero, usuario.id, datos.nombre, usuario.email, datos.direccion,
           datos.ciudad, datos.telefono, datos.metodo_pago, datos.estado, datos.total, datos.fecha
    FROM (VALUES
        ('PRUEBA-ADMIN-001', 'admin@upse.edu.ec', 'Administrador', 'Av. Universidad 100', 'Santa Elena', '0990000001', 'Tarjeta', 'entregado', 75.20::DOUBLE PRECISION, NOW() - INTERVAL '20 days'),
        ('PRUEBA-ADMIN-002', 'admin@upse.edu.ec', 'Administrador', 'Av. Universidad 100', 'Santa Elena', '0990000001', 'Transferencia', 'procesando', 68.40::DOUBLE PRECISION, NOW() - INTERVAL '2 days'),
        ('PRUEBA-RED-001', 'sofia.demo@multicatalogo.local', 'Sofía Rivera', 'Barrio Las Acacias', 'Salinas', '0990000011', 'Tarjeta', 'entregado', 62.70::DOUBLE PRECISION, NOW() - INTERVAL '12 days'),
        ('PRUEBA-RED-002', 'mateo.demo@multicatalogo.local', 'Mateo Salazar', 'Sector El Tambo', 'Santa Elena', '0990000012', 'Transferencia', 'enviado', 79.25::DOUBLE PRECISION, NOW() - INTERVAL '6 days'),
        ('PRUEBA-RED-003', 'valentina.demo@multicatalogo.local', 'Valentina Cruz', 'Ciudadela Italiana', 'Salinas', '0990000013', 'Efectivo', 'confirmado', 66.85::DOUBLE PRECISION, NOW() - INTERVAL '1 day')
    ) AS datos(numero, email, nombre, direccion, ciudad, telefono, metodo_pago, estado, total, fecha)
    JOIN usuarios usuario ON usuario.email = datos.email
    ON CONFLICT (numero) DO NOTHING;

    INSERT INTO pedido_items (pedido_id, producto_id, nombre, precio_unitario, cantidad, subtotal)
    SELECT pedido.id, datos.producto_id, datos.nombre, datos.precio, datos.cantidad,
           datos.precio * datos.cantidad
    FROM pedidos pedido
    JOIN (VALUES
        ('PRUEBA-ADMIN-001', 115, 'Reloj Vector', 59.00::DOUBLE PRECISION, 1),
        ('PRUEBA-ADMIN-001', 118, 'Taza Andina', 16.20::DOUBLE PRECISION, 1),
        ('PRUEBA-ADMIN-002', 119, 'Hub Nexus USB-C', 46.90::DOUBLE PRECISION, 1),
        ('PRUEBA-ADMIN-002', 120, 'Soporte Flex', 21.50::DOUBLE PRECISION, 1),
        ('PRUEBA-RED-001', 101, 'Auriculares Orbit', 39.90::DOUBLE PRECISION, 1),
        ('PRUEBA-RED-001', 107, 'Botella Atlas', 22.80::DOUBLE PRECISION, 1),
        ('PRUEBA-RED-002', 102, 'Teclado Nova', 54.50::DOUBLE PRECISION, 1),
        ('PRUEBA-RED-002', 103, 'Mouse Pulse', 24.75::DOUBLE PRECISION, 1),
        ('PRUEBA-RED-003', 108, 'Mochila Urban', 48.60::DOUBLE PRECISION, 1),
        ('PRUEBA-RED-003', 109, 'Organizador Zen', 18.25::DOUBLE PRECISION, 1)
    ) AS datos(numero, producto_id, nombre, precio, cantidad)
      ON pedido.numero = datos.numero;

    UPDATE productos producto
    SET stock = GREATEST(0, producto.stock - datos.cantidad)
    FROM (
        SELECT item.producto_id, SUM(item.cantidad)::INT AS cantidad
        FROM pedido_items item
        JOIN pedidos pedido ON pedido.id = item.pedido_id
        WHERE pedido.numero IN (
            'PRUEBA-ADMIN-001', 'PRUEBA-ADMIN-002', 'PRUEBA-RED-001',
            'PRUEBA-RED-002', 'PRUEBA-RED-003'
        )
        GROUP BY item.producto_id
    ) datos
    WHERE producto.id = datos.producto_id;

    WITH RECURSIVE ascendentes AS (
        SELECT pedido.id AS pedido_id, comprador.parent_id AS referido_id, 1 AS nivel
        FROM pedidos pedido
        JOIN usuarios usuario ON usuario.id = pedido.usuario_id
        JOIN referidos comprador ON comprador.id = usuario.referido_id
        WHERE pedido.numero IN ('PRUEBA-RED-001', 'PRUEBA-RED-002', 'PRUEBA-RED-003')
          AND comprador.parent_id IS NOT NULL

        UNION ALL

        SELECT ascendente.pedido_id, referido.parent_id, ascendente.nivel + 1
        FROM ascendentes ascendente
        JOIN referidos referido ON referido.id = ascendente.referido_id
        WHERE ascendente.nivel < 3 AND referido.parent_id IS NOT NULL
    ), nuevas AS (
        INSERT INTO comisiones (
            pedido_id, referido_id, nivel, porcentaje, monto, fecha, estado, pagada_en
        )
        SELECT pedido.id, ascendente.referido_id, ascendente.nivel,
               CASE ascendente.nivel WHEN 1 THEN 0.05 WHEN 2 THEN 0.03 ELSE 0.02 END,
               ROUND((pedido.total * CASE ascendente.nivel WHEN 1 THEN 0.05 WHEN 2 THEN 0.03 ELSE 0.02 END)::NUMERIC, 2)::DOUBLE PRECISION,
               pedido.fecha,
               CASE WHEN pedido.numero = 'PRUEBA-RED-001' AND referido.codigo = 'DEMO' THEN 'pagada' ELSE 'pendiente' END,
               CASE WHEN pedido.numero = 'PRUEBA-RED-001' AND referido.codigo = 'DEMO' THEN pedido.fecha + INTERVAL '3 days' END
        FROM ascendentes ascendente
        JOIN pedidos pedido ON pedido.id = ascendente.pedido_id
        JOIN referidos referido ON referido.id = ascendente.referido_id
        ON CONFLICT (pedido_id, referido_id) DO NOTHING
        RETURNING pedido_id, referido_id, monto
    ), totales AS (
        SELECT nueva.referido_id, SUM(pedido.total) AS ventas, SUM(nueva.monto) AS comisiones
        FROM nuevas nueva
        JOIN pedidos pedido ON pedido.id = nueva.pedido_id
        GROUP BY nueva.referido_id
    )
    UPDATE referidos referido
    SET ventas = referido.ventas + totales.ventas,
        comisiones = referido.comisiones + totales.comisiones
    FROM totales
    WHERE referido.id = totales.referido_id;

    WITH comision_pagada AS (
        SELECT comision.id, comision.referido_id, comision.monto, comision.pagada_en
        FROM comisiones comision
        JOIN pedidos pedido ON pedido.id = comision.pedido_id
        JOIN referidos referido ON referido.id = comision.referido_id
        WHERE pedido.numero = 'PRUEBA-RED-001' AND referido.codigo = 'DEMO'
    ), nueva_liquidacion AS (
        INSERT INTO liquidaciones (referido_id, monto, fecha)
        SELECT referido_id, monto, pagada_en FROM comision_pagada
        RETURNING id
    )
    INSERT INTO liquidacion_comisiones (liquidacion_id, comision_id)
    SELECT liquidacion.id, comision.id
    FROM nueva_liquidacion liquidacion
    CROSS JOIN comision_pagada comision;

    INSERT INTO demo_seed_history (version) VALUES ('011_escenarios_demostracion');

    PERFORM setval('referidos_id_seq', (SELECT MAX(id) FROM referidos));
    PERFORM setval('usuarios_id_seq', (SELECT MAX(id) FROM usuarios));
    PERFORM setval('pedidos_id_seq', (SELECT MAX(id) FROM pedidos));
    PERFORM setval('pedido_items_id_seq', (SELECT MAX(id) FROM pedido_items));
    PERFORM setval('comisiones_id_seq', (SELECT MAX(id) FROM comisiones));
    PERFORM setval('liquidaciones_id_seq', (SELECT MAX(id) FROM liquidaciones));
END $$;
