-- Agrega información de muestra para recorrer pedidos, comisiones y reportes.
-- El registro de versión evita duplicar datos cuando PostgreSQL y el backend ejecutan el archivo.
CREATE TABLE IF NOT EXISTS demo_seed_history (
    version TEXT PRIMARY KEY,
    aplicada_en TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM demo_seed_history WHERE version = '009_datos_demostracion') THEN
        RETURN;
    END IF;

    INSERT INTO demo_seed_history (version) VALUES ('009_datos_demostracion');

    UPDATE usuarios
    SET nombre = 'Administrador', rol = 'admin', activo = true, referido_id = 1
    WHERE email = 'admin@upse.edu.ec';

    UPDATE usuarios
    SET nombre = 'Cliente Demo', password = crypt('123456', gen_salt('bf')),
        rol = 'cliente', activo = true, referido_id = 6
    WHERE email = 'cliente@upse.edu.ec';

    INSERT INTO pedidos (
        numero, usuario_id, cliente_nombre, cliente_email, direccion, ciudad,
        telefono, metodo_pago, estado, total, fecha
    )
    SELECT datos.numero, usuario.id, 'Cliente Demo', 'cliente@upse.edu.ec', datos.direccion,
           datos.ciudad, '0991234567', datos.metodo_pago, datos.estado, datos.total, datos.fecha
    FROM usuarios usuario
    CROSS JOIN (VALUES
        ('DEMO-2026-001', 'Av. Malecón 120', 'Salinas', 'Tarjeta', 'entregado', 65.10::DOUBLE PRECISION, NOW() - INTERVAL '30 days'),
        ('DEMO-2026-002', 'Barrio Central, casa 8', 'La Libertad', 'Transferencia', 'entregado', 71.40::DOUBLE PRECISION, NOW() - INTERVAL '16 days'),
        ('DEMO-2026-003', 'Av. Eleodoro Solórzano 45', 'La Libertad', 'Tarjeta', 'enviado', 79.25::DOUBLE PRECISION, NOW() - INTERVAL '9 days'),
        ('DEMO-2026-004', 'Ciudadela Universitaria', 'Santa Elena', 'Efectivo', 'preparando', 59.75::DOUBLE PRECISION, NOW() - INTERVAL '4 days'),
        ('DEMO-2026-005', 'Sector La Milina', 'Salinas', 'Transferencia', 'confirmado', 92.45::DOUBLE PRECISION, NOW() - INTERVAL '1 day')
    ) AS datos(numero, direccion, ciudad, metodo_pago, estado, total, fecha)
    WHERE usuario.email = 'cliente@upse.edu.ec'
    ON CONFLICT (numero) DO NOTHING;

    INSERT INTO pedido_items (pedido_id, producto_id, nombre, precio_unitario, cantidad, subtotal)
    SELECT pedido.id, datos.producto_id, datos.nombre, datos.precio, datos.cantidad,
           datos.precio * datos.cantidad
    FROM pedidos pedido
    JOIN (VALUES
        ('DEMO-2026-001', 101, 'Auriculares Orbit', 39.90::DOUBLE PRECISION, 1),
        ('DEMO-2026-001', 117, 'Libreta Terra', 12.60::DOUBLE PRECISION, 2),
        ('DEMO-2026-002', 108, 'Mochila Urban', 48.60::DOUBLE PRECISION, 1),
        ('DEMO-2026-002', 107, 'Botella Atlas', 22.80::DOUBLE PRECISION, 1),
        ('DEMO-2026-003', 102, 'Teclado Nova', 54.50::DOUBLE PRECISION, 1),
        ('DEMO-2026-003', 103, 'Mouse Pulse', 24.75::DOUBLE PRECISION, 1),
        ('DEMO-2026-004', 112, 'Mat Yoga Flow', 29.95::DOUBLE PRECISION, 1),
        ('DEMO-2026-004', 111, 'Banda Active', 14.90::DOUBLE PRECISION, 2),
        ('DEMO-2026-005', 115, 'Reloj Vector', 59.00::DOUBLE PRECISION, 1),
        ('DEMO-2026-005', 116, 'Lentes Prisma', 33.45::DOUBLE PRECISION, 1)
    ) AS datos(numero, producto_id, nombre, precio, cantidad)
      ON pedido.numero = datos.numero;

    UPDATE productos producto
    SET stock = GREATEST(0, producto.stock - datos.cantidad)
    FROM (VALUES
        (101, 1), (117, 2), (108, 1), (107, 1), (102, 1),
        (103, 1), (112, 1), (111, 2), (115, 1), (116, 1)
    ) AS datos(producto_id, cantidad)
    WHERE producto.id = datos.producto_id;

    WITH nuevas AS (
        INSERT INTO comisiones (
            pedido_id, referido_id, nivel, porcentaje, monto, fecha, estado, pagada_en
        )
        SELECT pedido.id, nivel.referido_id, nivel.nivel, nivel.porcentaje,
               ROUND((pedido.total * nivel.porcentaje)::NUMERIC, 2)::DOUBLE PRECISION,
               pedido.fecha,
               CASE WHEN pedido.numero IN ('DEMO-2026-001', 'DEMO-2026-002') THEN 'pagada' ELSE 'pendiente' END,
               CASE WHEN pedido.numero IN ('DEMO-2026-001', 'DEMO-2026-002') THEN pedido.fecha + INTERVAL '2 days' END
        FROM pedidos pedido
        CROSS JOIN (VALUES
            (4, 1, 0.05::DOUBLE PRECISION),
            (2, 2, 0.03::DOUBLE PRECISION),
            (1, 3, 0.02::DOUBLE PRECISION)
        ) AS nivel(referido_id, nivel, porcentaje)
        WHERE pedido.numero IN (
            'DEMO-2026-001', 'DEMO-2026-002', 'DEMO-2026-003',
            'DEMO-2026-004', 'DEMO-2026-005'
        )
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

    INSERT INTO liquidaciones (id, referido_id, monto, fecha)
    SELECT datos.id, datos.referido_id, SUM(comision.monto), NOW() - INTERVAL '14 days'
    FROM (VALUES (9101::BIGINT, 1), (9102::BIGINT, 2), (9103::BIGINT, 4)) AS datos(id, referido_id)
    JOIN comisiones comision ON comision.referido_id = datos.referido_id
    JOIN pedidos pedido ON pedido.id = comision.pedido_id
    WHERE pedido.numero IN ('DEMO-2026-001', 'DEMO-2026-002')
    GROUP BY datos.id, datos.referido_id
    ON CONFLICT (id, referido_id) DO NOTHING;

    INSERT INTO liquidacion_comisiones (liquidacion_id, comision_id)
    SELECT CASE comision.referido_id WHEN 1 THEN 9101 WHEN 2 THEN 9102 WHEN 4 THEN 9103 END,
           comision.id
    FROM comisiones comision
    JOIN pedidos pedido ON pedido.id = comision.pedido_id
    WHERE pedido.numero IN ('DEMO-2026-001', 'DEMO-2026-002')
      AND comision.referido_id IN (1, 2, 4)
    ON CONFLICT (comision_id) DO NOTHING;

    PERFORM setval('pedidos_id_seq', (SELECT MAX(id) FROM pedidos));
    PERFORM setval('pedido_items_id_seq', (SELECT MAX(id) FROM pedido_items));
    PERFORM setval('comisiones_id_seq', (SELECT MAX(id) FROM comisiones));
    PERFORM setval('liquidaciones_id_seq', (SELECT MAX(id) FROM liquidaciones));
END $$;
