CREATE TABLE IF NOT EXISTS compras_idempotentes (
    usuario_id INT NOT NULL REFERENCES usuarios(id),
    clave TEXT NOT NULL,
    huella TEXT NOT NULL,
    respuesta JSONB NOT NULL,
    creado_en TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (usuario_id, clave)
);
