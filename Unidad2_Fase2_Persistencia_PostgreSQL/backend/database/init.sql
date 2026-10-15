-- init.sql: estructura relacional y datos semilla de MultiCatálogo (Fase 2).
-- Docker ejecuta este archivo automáticamente la primera vez que se crea el contenedor.

CREATE TABLE IF NOT EXISTS usuarios (
    id       SERIAL PRIMARY KEY,
    email    TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,           -- Texto plano temporal; se cifrará en una fase posterior.
    rol      TEXT NOT NULL DEFAULT 'cliente'
);

CREATE TABLE IF NOT EXISTS productos (
    id          SERIAL PRIMARY KEY,
    nombre      TEXT NOT NULL,
    descripcion TEXT NOT NULL DEFAULT '',
    precio      DOUBLE PRECISION NOT NULL,
    categoria   TEXT NOT NULL DEFAULT '',
    img         TEXT NOT NULL DEFAULT '',
    galeria     TEXT[] NOT NULL DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS referidos (
    id        SERIAL PRIMARY KEY,
    nombre    TEXT NOT NULL,
    nivel     INT NOT NULL,
    ventas    DOUBLE PRECISION NOT NULL DEFAULT 0,
    parent_id INT REFERENCES referidos(id)   -- Árbol: cada referido apunta a quien lo invitó.
);

-- Usuarios: mismo administrador que existía en memoria.
INSERT INTO usuarios (email, password, rol) VALUES
    ('admin@upse.edu.ec', '123456', 'admin')
ON CONFLICT (email) DO NOTHING;

-- Productos: los mismos cuatro artículos, ampliados con descripción, categoría y galería.
INSERT INTO productos (id, nombre, descripcion, precio, categoria, img, galeria) VALUES
    (1, 'Serum Revitalizante', 'Serum facial con vitamina C para iluminar la piel.', 45.00, 'Cuidado facial', 'https://picsum.photos/seed/serum/150',
        ARRAY['https://picsum.photos/seed/serum1/600', 'https://picsum.photos/seed/serum2/600']),
    (2, 'Crema Hidratante Pro', 'Crema de hidratación profunda para uso diario.', 32.50, 'Hidratación', 'https://picsum.photos/seed/crema/150',
        ARRAY['https://picsum.photos/seed/crema1/600', 'https://picsum.photos/seed/crema2/600']),
    (3, 'Tónico Purificante', 'Tónico que limpia y equilibra la piel grasa.', 28.00, 'Limpieza', 'https://picsum.photos/seed/tonico/150',
        ARRAY['https://picsum.photos/seed/tonico1/600']),
    (4, 'Mascarilla Nocturna', 'Mascarilla reparadora que actúa durante la noche.', 50.00, 'Tratamiento', 'https://picsum.photos/seed/mascarilla/150',
        ARRAY['https://picsum.photos/seed/mascarilla1/600', 'https://picsum.photos/seed/mascarilla2/600'])
ON CONFLICT (id) DO NOTHING;

-- Como se insertaron ids explícitos, adelantamos la secuencia para los próximos productos.
SELECT setval('productos_id_seq', (SELECT MAX(id) FROM productos));

-- Red de referidos: raíz administradora con dos niveles descendientes.
INSERT INTO referidos (id, nombre, nivel, ventas, parent_id) VALUES
    (1, 'Administrador', 1, 0,    NULL),
    (2, 'Ana Torres',    2, 1200, 1),
    (3, 'Luis Mera',     2, 800,  1),
    (4, 'María Yagual',  3, 400,  2),
    (5, 'Carlos Lindao', 3, 300,  3)
ON CONFLICT (id) DO NOTHING;

SELECT setval('referidos_id_seq', (SELECT MAX(id) FROM referidos));
