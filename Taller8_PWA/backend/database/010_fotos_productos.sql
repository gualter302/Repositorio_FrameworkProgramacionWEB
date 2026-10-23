-- Cambia las ilustraciones iniciales por fotografías locales optimizadas para el catálogo.
UPDATE productos AS producto
SET img = foto.ruta,
    galeria = ARRAY[foto.ruta]
FROM (VALUES
    (101, '/products/images/auriculares-orbit.webp'),
    (102, '/products/images/teclado-nova.webp'),
    (103, '/products/images/mouse-pulse.webp'),
    (104, '/products/images/parlante-echo.webp'),
    (105, '/products/images/bateria-volt.webp'),
    (106, '/products/images/lampara-aura.webp'),
    (107, '/products/images/botella-atlas.webp'),
    (108, '/products/images/mochila-urban.webp'),
    (109, '/products/images/organizador-zen.webp'),
    (110, '/products/images/difusor-bruma.webp'),
    (111, '/products/images/banda-active.webp'),
    (112, '/products/images/mat-yoga.webp'),
    (113, '/products/images/kit-lumen.webp'),
    (114, '/products/images/protector-nube.webp'),
    (115, '/products/images/reloj-vector.webp'),
    (116, '/products/images/lentes-prisma.webp'),
    (117, '/products/images/libreta-terra.webp'),
    (118, '/products/images/taza-andina.webp'),
    (119, '/products/images/hub-nexus.webp'),
    (120, '/products/images/soporte-flex.webp')
) AS foto(id, ruta)
WHERE producto.id = foto.id;
