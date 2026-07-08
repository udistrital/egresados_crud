-- =============================================================
-- Migración 2026-07-07 — Remapeo de ids de parámetros:
-- semilla LOCAL del MID → ids REALES del servicio institucional
-- (creados por el equipo el 2026-07-07: area_tipo EGR id=32,
-- tipo_parametro 174-179, parametro 7199-7230).
--
-- CUÁNDO ejecutar: al apagar BENEFICIOS_EGRESADOS_MID_PARAMETROS_LOCAL
-- (el MID pasa a resolver códigos contra el servicio real). Sin este
-- remapeo el catálogo sale vacío y los estados no resuelven, porque
-- las filas de la BD dev guardan los ids de la semilla local.
--
-- Idempotente: los ids viejos (10-62) y nuevos (7199+) no se solapan,
-- así que re-ejecutarla es un no-op.
--
-- NOTA: db/seed_estados_beneficio.sql y db/seed_pruebas.sql ya usan los
-- ids institucionales (actualizados 2026-07-08); esta migración solo hace
-- falta en BDs sembradas ANTES de esa fecha.
-- =============================================================

BEGIN;

SET search_path TO beneficios_egresados;

-- ── ESTADO_EMPRESA (tipo 174) ────────────────────────────────
-- Local: 10=ACTIVA (antes APROBADA — mismo significado: empresa nace
-- operativa, sin flujo de aprobación), 12=SUSPENDIDA.
UPDATE empresa SET estado_empresa_id = CASE estado_empresa_id
    WHEN 10 THEN 7199  -- ACTIVA
    WHEN 11 THEN 7199  -- PENDIENTE (semilla vieja) → ACTIVA: sin flujo de aprobación
    WHEN 12 THEN 7200  -- SUSPENDIDA
    ELSE estado_empresa_id END
WHERE estado_empresa_id IN (10, 11, 12);

-- ── SECTOR_ECONOMICO (tipo 178; columna nullable) ────────────
UPDATE empresa SET sector_economico_id = CASE sector_economico_id
    WHEN 50 THEN 7218  -- TEC
    WHEN 51 THEN 7222  -- COM
    WHEN 52 THEN 7227  -- OTR
    ELSE sector_economico_id END
WHERE sector_economico_id IN (50, 51, 52);

-- ── ESTADO_BENEFICIO (tipo 175) ──────────────────────────────
UPDATE beneficio SET estado_beneficio_id = CASE estado_beneficio_id
    WHEN 20 THEN 7201  -- BORRADOR
    WHEN 21 THEN 7202  -- PUBLICADO
    WHEN 22 THEN 7203  -- AGOTADO
    WHEN 23 THEN 7204  -- VENCIDO
    WHEN 24 THEN 7205  -- RETIRADO
    ELSE estado_beneficio_id END
WHERE estado_beneficio_id IN (20, 21, 22, 23, 24);

-- ── CATEGORIA_BENEFICIO (tipo 177) ───────────────────────────
UPDATE beneficio SET categoria_beneficio_id = CASE categoria_beneficio_id
    WHEN 40 THEN 7212  -- EDUCACION
    WHEN 41 THEN 7213  -- SALUD
    WHEN 42 THEN 7214  -- RECREACION
    WHEN 43 THEN 7215  -- EMPLEO
    WHEN 44 THEN 7216  -- DESCUENTOS
    WHEN 45 THEN 7217  -- OTRO
    ELSE categoria_beneficio_id END
WHERE categoria_beneficio_id IN (40, 41, 42, 43, 44, 45);

-- ── ESTADO_SOLICITUD (tipo 176) — historial, C-4b ────────────
UPDATE historial_solicitud SET estado_nuevo_id = CASE estado_nuevo_id
    WHEN 30 THEN 7206  -- PENDIENTE
    WHEN 31 THEN 7207  -- EN_REVISION
    WHEN 32 THEN 7208  -- REQUIERE_INFO
    WHEN 33 THEN 7209  -- APROBADA
    WHEN 34 THEN 7210  -- RECHAZADA
    WHEN 35 THEN 7211  -- CANCELADA
    ELSE estado_nuevo_id END
WHERE estado_nuevo_id IN (30, 31, 32, 33, 34, 35);

UPDATE historial_solicitud SET estado_anterior_id = CASE estado_anterior_id
    WHEN 30 THEN 7206
    WHEN 31 THEN 7207
    WHEN 32 THEN 7208
    WHEN 33 THEN 7209
    WHEN 34 THEN 7210
    WHEN 35 THEN 7211
    ELSE estado_anterior_id END
WHERE estado_anterior_id IN (30, 31, 32, 33, 34, 35);

COMMIT;

-- ── Verificación: no debe quedar ningún id local (<100) ──────
-- SELECT 'empresa.estado' col, count(*) FROM empresa WHERE estado_empresa_id < 100
-- UNION ALL SELECT 'empresa.sector', count(*) FROM empresa WHERE sector_economico_id < 100
-- UNION ALL SELECT 'beneficio.estado', count(*) FROM beneficio WHERE estado_beneficio_id < 100
-- UNION ALL SELECT 'beneficio.categoria', count(*) FROM beneficio WHERE categoria_beneficio_id < 100
-- UNION ALL SELECT 'historial.nuevo', count(*) FROM historial_solicitud WHERE estado_nuevo_id < 100
-- UNION ALL SELECT 'historial.anterior', count(*) FROM historial_solicitud WHERE estado_anterior_id < 100;
