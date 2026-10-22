-- Estructura y datos iniciales de MultiCatálogo para el Taller 7.
-- Docker ejecuta este archivo cuando crea por primera vez el volumen de PostgreSQL.

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS referidos (
    id          SERIAL PRIMARY KEY,
    nombre      TEXT NOT NULL,
    codigo      TEXT NOT NULL UNIQUE,
    nivel       INT NOT NULL,
    ventas      DOUBLE PRECISION NOT NULL DEFAULT 0,
    comisiones  DOUBLE PRECISION NOT NULL DEFAULT 0,
    parent_id   INT REFERENCES referidos(id)
);

CREATE TABLE IF NOT EXISTS usuarios (
    id          SERIAL PRIMARY KEY,
    nombre      TEXT NOT NULL,
    email       TEXT NOT NULL UNIQUE,
    password    TEXT NOT NULL,
    rol         TEXT NOT NULL DEFAULT 'cliente',
    activo      BOOLEAN NOT NULL DEFAULT true,
    referido_id INT NOT NULL REFERENCES referidos(id),
    creado_en   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS productos (
    id          SERIAL PRIMARY KEY,
    nombre      TEXT NOT NULL,
    descripcion TEXT NOT NULL DEFAULT '',
    precio      DOUBLE PRECISION NOT NULL CHECK (precio >= 0),
    categoria   TEXT NOT NULL DEFAULT '',
    img         TEXT NOT NULL DEFAULT '',
    galeria     TEXT[] NOT NULL DEFAULT '{}',
    stock       INT NOT NULL DEFAULT 0 CHECK (stock >= 0),
    activo      BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE IF NOT EXISTS pedidos (
    id                BIGSERIAL PRIMARY KEY,
    numero            TEXT UNIQUE,
    usuario_id        INT NOT NULL REFERENCES usuarios(id),
    cliente_nombre    TEXT NOT NULL,
    cliente_email     TEXT NOT NULL,
    direccion         TEXT NOT NULL,
    ciudad            TEXT NOT NULL,
    telefono          TEXT NOT NULL,
    metodo_pago       TEXT NOT NULL,
    estado            TEXT NOT NULL DEFAULT 'Confirmado',
    total             DOUBLE PRECISION NOT NULL CHECK (total >= 0),
    fecha             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pedido_items (
    id              BIGSERIAL PRIMARY KEY,
    pedido_id       BIGINT NOT NULL REFERENCES pedidos(id) ON DELETE CASCADE,
    producto_id     INT NOT NULL REFERENCES productos(id),
    nombre          TEXT NOT NULL,
    precio_unitario DOUBLE PRECISION NOT NULL CHECK (precio_unitario >= 0),
    cantidad        INT NOT NULL CHECK (cantidad > 0),
    subtotal        DOUBLE PRECISION NOT NULL CHECK (subtotal >= 0)
);

CREATE TABLE IF NOT EXISTS comisiones (
    id          BIGSERIAL PRIMARY KEY,
    pedido_id   BIGINT NOT NULL REFERENCES pedidos(id) ON DELETE CASCADE,
    referido_id INT NOT NULL REFERENCES referidos(id),
    nivel       INT NOT NULL CHECK (nivel BETWEEN 1 AND 3),
    porcentaje  DOUBLE PRECISION NOT NULL,
    monto       DOUBLE PRECISION NOT NULL CHECK (monto >= 0),
    fecha       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    estado      TEXT NOT NULL DEFAULT 'pendiente',
    pagada_en   TIMESTAMPTZ,
    UNIQUE (pedido_id, referido_id)
);

CREATE TABLE IF NOT EXISTS liquidaciones (
    id          BIGSERIAL PRIMARY KEY,
    referido_id INT NOT NULL REFERENCES referidos(id),
    monto       DOUBLE PRECISION NOT NULL CHECK (monto > 0),
    fecha       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (id, referido_id)
);

CREATE TABLE IF NOT EXISTS liquidacion_comisiones (
    liquidacion_id BIGINT NOT NULL REFERENCES liquidaciones(id) ON DELETE CASCADE,
    comision_id    BIGINT NOT NULL UNIQUE REFERENCES comisiones(id),
    PRIMARY KEY (liquidacion_id, comision_id)
);

-- Tres niveles ascendentes desde el cliente de demostración permiten comprobar las comisiones.
INSERT INTO referidos (id, nombre, codigo, nivel, ventas, comisiones, parent_id) VALUES
    (1, 'Administrador', 'ADMIN', 1, 0,    0, NULL),
    (2, 'Ana Torres',    'ANA',   2, 1200, 0, 1),
    (3, 'Luis Mera',     'LUIS',  2, 800,  0, 1),
    (4, 'María Yagual',  'MARIA', 3, 400,  0, 2),
    (5, 'Carlos Lindao', 'CARLOS',3, 300,  0, 3),
    (6, 'Cliente Demo',  'DEMO',  4, 0,    0, 4)
ON CONFLICT (id) DO NOTHING;

SELECT setval('referidos_id_seq', (SELECT MAX(id) FROM referidos));

INSERT INTO usuarios (nombre, email, password, rol, activo, referido_id) VALUES
    ('Administrador', 'admin@upse.edu.ec', crypt('123456123456', gen_salt('bf')), 'admin', true, 1),
    ('Cliente Demo', 'cliente@upse.edu.ec', crypt('123456', gen_salt('bf')), 'cliente', true, 6)
ON CONFLICT (email) DO NOTHING;

INSERT INTO productos (id, nombre, descripcion, precio, categoria, img, galeria, stock) VALUES
    (1, 'Serum Revitalizante', 'Serum facial con vitamina C para iluminar la piel.', 45.00, 'Cuidado facial', 'https://picsum.photos/seed/serum/150',
        ARRAY['https://picsum.photos/seed/serum1/600', 'https://picsum.photos/seed/serum2/600'], 25),
    (2, 'Crema Hidratante Pro', 'Crema de hidratación profunda para uso diario.', 32.50, 'Hidratación', 'https://picsum.photos/seed/crema/150',
        ARRAY['https://picsum.photos/seed/crema1/600', 'https://picsum.photos/seed/crema2/600'], 30),
    (3, 'Tónico Purificante', 'Tónico que limpia y equilibra la piel grasa.', 28.00, 'Limpieza', 'https://picsum.photos/seed/tonico/150',
        ARRAY['https://picsum.photos/seed/tonico1/600'], 20),
    (4, 'Mascarilla Nocturna', 'Mascarilla reparadora que actúa durante la noche.', 50.00, 'Tratamiento', 'https://picsum.photos/seed/mascarilla/150',
        ARRAY['https://picsum.photos/seed/mascarilla1/600', 'https://picsum.photos/seed/mascarilla2/600'], 15)
ON CONFLICT (id) DO NOTHING;

SELECT setval('productos_id_seq', (SELECT MAX(id) FROM productos));
