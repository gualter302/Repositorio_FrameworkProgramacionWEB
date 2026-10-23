-- Sustituye los artículos de demostración por un catálogo propio sin borrar el historial de pedidos.
UPDATE productos
SET activo = false
WHERE id BETWEEN 1 AND 4
  AND nombre IN ('Serum Revitalizante', 'Crema Hidratante Pro', 'Tónico Purificante', 'Mascarilla Nocturna');

INSERT INTO productos (id, nombre, descripcion, precio, categoria, img, galeria, stock, activo) VALUES
    (101, 'Auriculares Orbit', 'Auriculares inalámbricos con estuche compacto y sonido equilibrado.', 39.90, 'Tecnología', '/products/images/auriculares-orbit.webp', ARRAY['/products/images/auriculares-orbit.webp'], 28, true),
    (102, 'Teclado Nova', 'Teclado mecánico compacto con iluminación suave y conexión USB-C.', 54.50, 'Tecnología', '/products/images/teclado-nova.webp', ARRAY['/products/images/teclado-nova.webp'], 16, true),
    (103, 'Mouse Pulse', 'Mouse ergonómico de precisión con seis botones y sensibilidad ajustable.', 24.75, 'Tecnología', '/products/images/mouse-pulse.webp', ARRAY['/products/images/mouse-pulse.webp'], 35, true),
    (104, 'Parlante Echo Mini', 'Parlante Bluetooth portátil con doce horas de autonomía.', 44.00, 'Tecnología', '/products/images/parlante-echo.webp', ARRAY['/products/images/parlante-echo.webp'], 20, true),
    (105, 'Batería Volt 10000', 'Batería externa de carga rápida con dos puertos USB.', 31.90, 'Tecnología', '/products/images/bateria-volt.webp', ARRAY['/products/images/bateria-volt.webp'], 24, true),
    (106, 'Lámpara Aura', 'Lámpara de escritorio regulable con tres tonos de luz.', 27.40, 'Hogar', '/products/images/lampara-aura.webp', ARRAY['/products/images/lampara-aura.webp'], 18, true),
    (107, 'Botella Atlas', 'Botella térmica de acero de 750 ml para bebidas frías o calientes.', 22.80, 'Hogar', '/products/images/botella-atlas.webp', ARRAY['/products/images/botella-atlas.webp'], 32, true),
    (108, 'Mochila Urban', 'Mochila resistente al agua con espacio acolchado para portátil.', 48.60, 'Accesorios', '/products/images/mochila-urban.webp', ARRAY['/products/images/mochila-urban.webp'], 14, true),
    (109, 'Organizador Zen', 'Organizador modular para mantener ordenado el escritorio.', 18.25, 'Hogar', '/products/images/organizador-zen.webp', ARRAY['/products/images/organizador-zen.webp'], 30, true),
    (110, 'Difusor Bruma', 'Difusor ultrasónico con luz ambiental y apagado automático.', 35.70, 'Hogar', '/products/images/difusor-bruma.webp', ARRAY['/products/images/difusor-bruma.webp'], 19, true),
    (111, 'Banda Active', 'Banda elástica de resistencia para entrenamientos y movilidad.', 14.90, 'Bienestar', '/products/images/banda-active.webp', ARRAY['/products/images/banda-active.webp'], 42, true),
    (112, 'Mat Yoga Flow', 'Tapete antideslizante y ligero para yoga o ejercicio en casa.', 29.95, 'Bienestar', '/products/images/mat-yoga.webp', ARRAY['/products/images/mat-yoga.webp'], 21, true),
    (113, 'Kit Facial Lumen', 'Set diario de limpieza e hidratación para todo tipo de piel.', 41.30, 'Bienestar', '/products/images/kit-lumen.webp', ARRAY['/products/images/kit-lumen.webp'], 17, true),
    (114, 'Protector Solar Nube', 'Protector facial ligero FPS 50 de acabado transparente.', 19.80, 'Bienestar', '/products/images/protector-nube.webp', ARRAY['/products/images/protector-nube.webp'], 38, true),
    (115, 'Reloj Vector', 'Reloj minimalista con correa ajustable y resistencia a salpicaduras.', 59.00, 'Accesorios', '/products/images/reloj-vector.webp', ARRAY['/products/images/reloj-vector.webp'], 12, true),
    (116, 'Lentes Prisma', 'Lentes de sol polarizados con marco liviano y estuche protector.', 33.45, 'Accesorios', '/products/images/lentes-prisma.webp', ARRAY['/products/images/lentes-prisma.webp'], 23, true),
    (117, 'Libreta Terra', 'Libreta de tapa dura con 160 páginas y separador de tela.', 12.60, 'Papelería', '/products/images/libreta-terra.webp', ARRAY['/products/images/libreta-terra.webp'], 45, true),
    (118, 'Taza Andina', 'Taza de cerámica de 350 ml con acabado mate.', 16.20, 'Hogar', '/products/images/taza-andina.webp', ARRAY['/products/images/taza-andina.webp'], 27, true),
    (119, 'Hub Nexus USB-C', 'Adaptador multipuerto con HDMI, USB 3.0 y lector de tarjetas.', 46.90, 'Tecnología', '/products/images/hub-nexus.webp', ARRAY['/products/images/hub-nexus.webp'], 15, true),
    (120, 'Soporte Flex', 'Soporte plegable de aluminio para teléfono o tableta.', 21.50, 'Accesorios', '/products/images/soporte-flex.webp', ARRAY['/products/images/soporte-flex.webp'], 34, true)
ON CONFLICT (id) DO UPDATE SET
    nombre = EXCLUDED.nombre,
    descripcion = EXCLUDED.descripcion,
    precio = EXCLUDED.precio,
    categoria = EXCLUDED.categoria,
    img = EXCLUDED.img,
    galeria = EXCLUDED.galeria,
    stock = EXCLUDED.stock,
    activo = true;

SELECT setval('productos_id_seq', (SELECT MAX(id) FROM productos));
