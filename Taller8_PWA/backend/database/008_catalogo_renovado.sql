-- Sustituye los artículos de demostración por un catálogo propio sin borrar el historial de pedidos.
UPDATE productos
SET activo = false
WHERE id BETWEEN 1 AND 4
  AND nombre IN ('Serum Revitalizante', 'Crema Hidratante Pro', 'Tónico Purificante', 'Mascarilla Nocturna');

INSERT INTO productos (id, nombre, descripcion, precio, categoria, img, galeria, stock, activo) VALUES
    (101, 'Auriculares Orbit', 'Auriculares inalámbricos con estuche compacto y sonido equilibrado.', 39.90, 'Tecnología', '/products/auriculares-orbit.svg', ARRAY['/products/auriculares-orbit.svg'], 28, true),
    (102, 'Teclado Nova', 'Teclado mecánico compacto con iluminación suave y conexión USB-C.', 54.50, 'Tecnología', '/products/teclado-nova.svg', ARRAY['/products/teclado-nova.svg'], 16, true),
    (103, 'Mouse Pulse', 'Mouse ergonómico de precisión con seis botones y sensibilidad ajustable.', 24.75, 'Tecnología', '/products/mouse-pulse.svg', ARRAY['/products/mouse-pulse.svg'], 35, true),
    (104, 'Parlante Echo Mini', 'Parlante Bluetooth portátil con doce horas de autonomía.', 44.00, 'Tecnología', '/products/parlante-echo.svg', ARRAY['/products/parlante-echo.svg'], 20, true),
    (105, 'Batería Volt 10000', 'Batería externa de carga rápida con dos puertos USB.', 31.90, 'Tecnología', '/products/bateria-volt.svg', ARRAY['/products/bateria-volt.svg'], 24, true),
    (106, 'Lámpara Aura', 'Lámpara de escritorio regulable con tres tonos de luz.', 27.40, 'Hogar', '/products/lampara-aura.svg', ARRAY['/products/lampara-aura.svg'], 18, true),
    (107, 'Botella Atlas', 'Botella térmica de acero de 750 ml para bebidas frías o calientes.', 22.80, 'Hogar', '/products/botella-atlas.svg', ARRAY['/products/botella-atlas.svg'], 32, true),
    (108, 'Mochila Urban', 'Mochila resistente al agua con espacio acolchado para portátil.', 48.60, 'Accesorios', '/products/mochila-urban.svg', ARRAY['/products/mochila-urban.svg'], 14, true),
    (109, 'Organizador Zen', 'Organizador modular para mantener ordenado el escritorio.', 18.25, 'Hogar', '/products/organizador-zen.svg', ARRAY['/products/organizador-zen.svg'], 30, true),
    (110, 'Difusor Bruma', 'Difusor ultrasónico con luz ambiental y apagado automático.', 35.70, 'Hogar', '/products/difusor-bruma.svg', ARRAY['/products/difusor-bruma.svg'], 19, true),
    (111, 'Banda Active', 'Banda elástica de resistencia para entrenamientos y movilidad.', 14.90, 'Bienestar', '/products/banda-active.svg', ARRAY['/products/banda-active.svg'], 42, true),
    (112, 'Mat Yoga Flow', 'Tapete antideslizante y ligero para yoga o ejercicio en casa.', 29.95, 'Bienestar', '/products/mat-yoga.svg', ARRAY['/products/mat-yoga.svg'], 21, true),
    (113, 'Kit Facial Lumen', 'Set diario de limpieza e hidratación para todo tipo de piel.', 41.30, 'Bienestar', '/products/kit-lumen.svg', ARRAY['/products/kit-lumen.svg'], 17, true),
    (114, 'Protector Solar Nube', 'Protector facial ligero FPS 50 de acabado transparente.', 19.80, 'Bienestar', '/products/protector-nube.svg', ARRAY['/products/protector-nube.svg'], 38, true),
    (115, 'Reloj Vector', 'Reloj minimalista con correa ajustable y resistencia a salpicaduras.', 59.00, 'Accesorios', '/products/reloj-vector.svg', ARRAY['/products/reloj-vector.svg'], 12, true),
    (116, 'Lentes Prisma', 'Lentes de sol polarizados con marco liviano y estuche protector.', 33.45, 'Accesorios', '/products/lentes-prisma.svg', ARRAY['/products/lentes-prisma.svg'], 23, true),
    (117, 'Libreta Terra', 'Libreta de tapa dura con 160 páginas y separador de tela.', 12.60, 'Papelería', '/products/libreta-terra.svg', ARRAY['/products/libreta-terra.svg'], 45, true),
    (118, 'Taza Andina', 'Taza de cerámica de 350 ml con acabado mate.', 16.20, 'Hogar', '/products/taza-andina.svg', ARRAY['/products/taza-andina.svg'], 27, true),
    (119, 'Hub Nexus USB-C', 'Adaptador multipuerto con HDMI, USB 3.0 y lector de tarjetas.', 46.90, 'Tecnología', '/products/hub-nexus.svg', ARRAY['/products/hub-nexus.svg'], 15, true),
    (120, 'Soporte Flex', 'Soporte plegable de aluminio para teléfono o tableta.', 21.50, 'Accesorios', '/products/soporte-flex.svg', ARRAY['/products/soporte-flex.svg'], 34, true)
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
