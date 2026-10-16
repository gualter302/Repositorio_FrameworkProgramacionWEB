-- Permite actualizar una base de la Fase 4 sin borrar pedidos ni productos.
CREATE EXTENSION IF NOT EXISTS pgcrypto;

ALTER TABLE referidos ADD COLUMN IF NOT EXISTS codigo TEXT;
UPDATE referidos SET codigo = 'REF' || LPAD(id::TEXT, 6, '0') WHERE codigo IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS referidos_codigo_unique ON referidos(codigo);
ALTER TABLE referidos ALTER COLUMN codigo SET NOT NULL;

ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS nombre TEXT;
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS activo BOOLEAN NOT NULL DEFAULT true;
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS creado_en TIMESTAMPTZ NOT NULL DEFAULT NOW();
UPDATE usuarios SET nombre = SPLIT_PART(email, '@', 1) WHERE nombre IS NULL;
ALTER TABLE usuarios ALTER COLUMN nombre SET NOT NULL;

UPDATE usuarios
SET password = crypt(password, gen_salt('bf'))
WHERE password NOT LIKE '$2%';

UPDATE referidos SET codigo = 'ADMIN' WHERE id = 1;
UPDATE referidos SET codigo = 'ANA' WHERE id = 2;
UPDATE referidos SET codigo = 'LUIS' WHERE id = 3;
UPDATE referidos SET codigo = 'MARIA' WHERE id = 4;
UPDATE referidos SET codigo = 'CARLOS' WHERE id = 5;
UPDATE referidos SET codigo = 'DEMO' WHERE id = 6;
