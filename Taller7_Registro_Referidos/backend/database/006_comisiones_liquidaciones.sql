-- Agrega estados auditables y liquidaciones sin borrar el historial de comisiones.
ALTER TABLE comisiones ADD COLUMN IF NOT EXISTS estado TEXT NOT NULL DEFAULT 'pendiente';
ALTER TABLE comisiones ADD COLUMN IF NOT EXISTS pagada_en TIMESTAMPTZ;

CREATE TABLE IF NOT EXISTS liquidaciones (
    id          BIGSERIAL PRIMARY KEY,
    referido_id INT NOT NULL REFERENCES referidos(id),
    monto       DOUBLE PRECISION NOT NULL CHECK (monto > 0),
    fecha        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (id, referido_id)
);

CREATE TABLE IF NOT EXISTS liquidacion_comisiones (
    liquidacion_id BIGINT NOT NULL REFERENCES liquidaciones(id) ON DELETE CASCADE,
    comision_id    BIGINT NOT NULL UNIQUE REFERENCES comisiones(id),
    PRIMARY KEY (liquidacion_id, comision_id)
);
