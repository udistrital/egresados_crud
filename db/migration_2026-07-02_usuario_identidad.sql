-- Migración 2026-07-02: identidad de usuario para el JIT de empresa. Idempotente.
--
-- Los usuarios de empresa (self-signup en WSO2) NO tienen documento en ninguna
-- fuente institucional (userinfo ni userRol, verificado 2026-07-01), así que el
-- JIT los inserta con documento = NULL y su llave de identidad local pasa a ser
-- (sistema_origen, id_externo) con el sub de WSO2. El schema.sql ya lo trae;
-- este script actualiza las BD creadas con la versión anterior (documento
-- NOT NULL), donde el provision de empresa falla con:
--   pq: el valor nulo en la columna «documento» ... viola la restricción not-null

SET search_path = beneficios_egresados;

-- 1. documento nullable (UNIQUE se conserva: Postgres trata los NULL como
--    distintos, así que múltiples empresas sin documento no colisionan).
ALTER TABLE usuario ALTER COLUMN documento DROP NOT NULL;

COMMENT ON COLUMN usuario.documento IS
    'NULL para empresas self-signup (no tienen cédula); egresados sí lo traen.';

-- 2. Identidad externa única: (sistema_origen, id_externo) = sub de WSO2.
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'uq_usuario_id_externo'
          AND conrelid = 'beneficios_egresados.usuario'::regclass
    ) THEN
        ALTER TABLE usuario
            ADD CONSTRAINT uq_usuario_id_externo UNIQUE (sistema_origen, id_externo);
    END IF;
END $$;
