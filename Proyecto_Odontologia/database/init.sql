CREATE TABLE roles (id SMALLSERIAL PRIMARY KEY, nombre VARCHAR(40) UNIQUE NOT NULL);
CREATE TABLE usuarios (id BIGSERIAL PRIMARY KEY, rol_id SMALLINT NOT NULL REFERENCES roles(id), nombre VARCHAR(120) NOT NULL, email VARCHAR(180) UNIQUE NOT NULL, password_hash TEXT NOT NULL, mfa_habilitado BOOLEAN NOT NULL DEFAULT FALSE, activo BOOLEAN NOT NULL DEFAULT TRUE, creado_en TIMESTAMPTZ NOT NULL DEFAULT NOW());
CREATE TABLE pacientes (id BIGSERIAL PRIMARY KEY, identificacion VARCHAR(30) UNIQUE NOT NULL, nombres VARCHAR(100) NOT NULL, apellidos VARCHAR(100) NOT NULL, fecha_nacimiento DATE NOT NULL, telefono VARCHAR(30) NOT NULL, email VARCHAR(180), direccion TEXT, contacto_emergencia TEXT, tutor_legal TEXT, seguro_convenio TEXT, creado_en TIMESTAMPTZ NOT NULL DEFAULT NOW(), actualizado_en TIMESTAMPTZ NOT NULL DEFAULT NOW());
CREATE TABLE anamnesis (id BIGSERIAL PRIMARY KEY, paciente_id BIGINT UNIQUE NOT NULL REFERENCES pacientes(id), enfermedades TEXT[] NOT NULL DEFAULT '{}', alergias TEXT[] NOT NULL DEFAULT '{}', antecedentes_quirurgicos TEXT, anticoagulado BOOLEAN NOT NULL DEFAULT FALSE, actualizado_en TIMESTAMPTZ NOT NULL DEFAULT NOW());
CREATE TABLE profesionales (id BIGSERIAL PRIMARY KEY, usuario_id BIGINT UNIQUE NOT NULL REFERENCES usuarios(id), especialidad VARCHAR(100), registro_profesional VARCHAR(80) UNIQUE NOT NULL);
CREATE TABLE sillones (id SMALLSERIAL PRIMARY KEY, nombre VARCHAR(60) UNIQUE NOT NULL, activo BOOLEAN NOT NULL DEFAULT TRUE);
CREATE TABLE citas (id BIGSERIAL PRIMARY KEY, paciente_id BIGINT NOT NULL REFERENCES pacientes(id), profesional_id BIGINT NOT NULL REFERENCES profesionales(id), sillon_id SMALLINT NOT NULL REFERENCES sillones(id), inicio TIMESTAMPTZ NOT NULL, fin TIMESTAMPTZ NOT NULL, estado VARCHAR(30) NOT NULL, motivo TEXT NOT NULL, precio_estimado NUMERIC(12,2) NOT NULL DEFAULT 0, CHECK (fin > inicio));
CREATE EXTENSION IF NOT EXISTS btree_gist;
ALTER TABLE citas ADD CONSTRAINT citas_sillon_sin_superposicion EXCLUDE USING gist (sillon_id WITH =, tstzrange(inicio, fin, '[)') WITH &&) WHERE (estado <> 'Cancelado');
CREATE TABLE auditoria (id BIGSERIAL PRIMARY KEY, usuario_id BIGINT REFERENCES usuarios(id), accion VARCHAR(20) NOT NULL, entidad VARCHAR(80) NOT NULL, entidad_id VARCHAR(80), fecha TIMESTAMPTZ NOT NULL DEFAULT NOW(), ip INET, datos JSONB NOT NULL DEFAULT '{}');
INSERT INTO roles(nombre) VALUES ('administrador'),('recepcionista'),('odontologo'),('asistente'),('auditor');
INSERT INTO sillones(nombre) VALUES ('Sillón 1'),('Sillón 2'),('Sillón 3');

CREATE TABLE odontograma (
  id BIGSERIAL PRIMARY KEY, paciente_id BIGINT NOT NULL REFERENCES pacientes(id), diente SMALLINT NOT NULL,
  superficie VARCHAR(20) NOT NULL, condicion VARCHAR(60) NOT NULL, nomenclatura VARCHAR(20) NOT NULL DEFAULT 'FDI',
  notas TEXT, registrado_por BIGINT REFERENCES usuarios(id), registrado_en TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE periodontograma (
  id BIGSERIAL PRIMARY KEY, paciente_id BIGINT NOT NULL REFERENCES pacientes(id), diente SMALLINT NOT NULL,
  profundidades SMALLINT[6] NOT NULL, margen SMALLINT NOT NULL DEFAULT 0, furca SMALLINT NOT NULL DEFAULT 0,
  movilidad SMALLINT NOT NULL DEFAULT 0, sangrado BOOLEAN NOT NULL DEFAULT FALSE, registrado_en TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE notas_clinicas (
  id BIGSERIAL PRIMARY KEY, paciente_id BIGINT NOT NULL REFERENCES pacientes(id), subjetivo TEXT NOT NULL,
  objetivo TEXT NOT NULL, evaluacion TEXT NOT NULL, plan TEXT NOT NULL, firmado_por BIGINT REFERENCES usuarios(id),
  firmado_en TIMESTAMPTZ NOT NULL DEFAULT NOW(), bloqueado_en TIMESTAMPTZ NOT NULL
);
CREATE TABLE adendas_clinicas (
  id BIGSERIAL PRIMARY KEY, nota_id BIGINT NOT NULL REFERENCES notas_clinicas(id), motivo TEXT NOT NULL,
  contenido TEXT NOT NULL, firmado_por BIGINT REFERENCES usuarios(id), creado_en TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE consentimientos (
  id BIGSERIAL PRIMARY KEY, paciente_id BIGINT NOT NULL REFERENCES pacientes(id), procedimiento TEXT NOT NULL,
  riesgos TEXT NOT NULL, firmante VARCHAR(180) NOT NULL, firma TEXT NOT NULL, otp_validado BOOLEAN NOT NULL,
  firmado_en TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE planes_tratamiento (
  id BIGSERIAL PRIMARY KEY, paciente_id BIGINT NOT NULL REFERENCES pacientes(id), estado VARCHAR(30) NOT NULL,
  total NUMERIC(12,2) NOT NULL DEFAULT 0, creado_en TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE items_tratamiento (
  id BIGSERIAL PRIMARY KEY, plan_id BIGINT NOT NULL REFERENCES planes_tratamiento(id), procedimiento TEXT NOT NULL,
  fase SMALLINT NOT NULL, sesiones SMALLINT NOT NULL, honorarios NUMERIC(12,2) NOT NULL,
  laboratorio NUMERIC(12,2) NOT NULL DEFAULT 0, descuento NUMERIC(12,2) NOT NULL DEFAULT 0, estado VARCHAR(30) NOT NULL
);
CREATE TABLE inventario (
  id BIGSERIAL PRIMARY KEY, nombre VARCHAR(160) UNIQUE NOT NULL, unidad VARCHAR(40) NOT NULL,
  stock NUMERIC(12,3) NOT NULL DEFAULT 0, minimo NUMERIC(12,3) NOT NULL DEFAULT 0,
  costo_unitario NUMERIC(12,2) NOT NULL DEFAULT 0, actualizado_en TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE facturas (
  id BIGSERIAL PRIMARY KEY, numero VARCHAR(60) UNIQUE NOT NULL, paciente_id BIGINT NOT NULL REFERENCES pacientes(id),
  plan_id BIGINT REFERENCES planes_tratamiento(id), subtotal NUMERIC(12,2) NOT NULL, impuesto NUMERIC(12,2) NOT NULL,
  total NUMERIC(12,2) NOT NULL, forma_pago VARCHAR(50) NOT NULL, estado VARCHAR(30) NOT NULL,
  estado_fiscal VARCHAR(50) NOT NULL, creado_en TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE cierres_caja (
  id BIGSERIAL PRIMARY KEY, fecha DATE UNIQUE NOT NULL, efectivo_declarado NUMERIC(12,2) NOT NULL,
  efectivo_esperado NUMERIC(12,2) NOT NULL, diferencia NUMERIC(12,2) NOT NULL,
  cerrado_por BIGINT REFERENCES usuarios(id), creado_en TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE liquidaciones_profesionales (
  id BIGSERIAL PRIMARY KEY, profesional_id BIGINT REFERENCES profesionales(id), periodo VARCHAR(7) NOT NULL,
  monto NUMERIC(12,2) NOT NULL, estado VARCHAR(30) NOT NULL, creado_en TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE bandeja_mensajes (
  id BIGSERIAL PRIMARY KEY, canal VARCHAR(30) NOT NULL, destinatario VARCHAR(180) NOT NULL,
  plantilla VARCHAR(80) NOT NULL, estado VARCHAR(30) NOT NULL DEFAULT 'pending', creado_en TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- La aplicación actual guarda una instantánea transaccional en este registro para operar
-- todos los módulos con una sola escritura atómica. Las tablas normalizadas anteriores
-- documentan el modelo destino para la separación posterior por repositorios.
CREATE TABLE sigo_app_state (
  id SMALLINT PRIMARY KEY CHECK (id = 1), data JSONB NOT NULL, updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

