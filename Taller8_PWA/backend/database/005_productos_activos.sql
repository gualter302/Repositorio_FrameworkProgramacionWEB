-- Permite desactivar productos sin romper el historial de pedidos.
ALTER TABLE productos ADD COLUMN IF NOT EXISTS activo BOOLEAN NOT NULL DEFAULT true;
