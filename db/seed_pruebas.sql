-- =============================================================
-- seed_pruebas.sql — Datos base para probar el módulo (NO producción)
-- =============================================================
-- Crea los registros prerrequisito (usuario / empresa / egresado) que las FK
-- locales exigen antes de poder generar beneficios y solicitudes.
--
-- Ejecutar DESPUÉS de aplicar db/schema.sql (v4) sobre la BD beneficios_egresados:
--     psql -d beneficios_egresados -f db/schema.sql
--     psql -d beneficios_egresados -f db/seed_pruebas.sql
--
-- Nota (C-6): los campos *_id de estado/categoría/sector son REFERENCIAS LÓGICAS
-- a parámetros institucionales; aquí se usan valores PLACEHOLDER (no hay FK en v4).
-- En producción esos ids los resuelve el MID contra el servicio de parámetros.
--   estado_empresa_id   = 1  (~ ACTIVA)
--   sector_economico_id = 1  (~ TEC)
-- =============================================================

SET search_path TO beneficios_egresados;

-- Limpieza previa (orden inverso a las FK) + reinicio de identidades
TRUNCATE historial_solicitud, mensaje_solicitud, solicitud_beneficio, beneficio,
         usuario_empresa, egresado, empresa, usuario, bitacora_acceso_pii
         RESTART IDENTITY CASCADE;

-- ── Usuarios (tipo_usuario es el discriminador local C-7: EGR | EMP | ADM) ──
INSERT INTO usuario (id, documento, nombre, correo, tipo_usuario, id_externo, sistema_origen, activo) VALUES
  (1, '900111222',  'Representante Empresa Demo', 'rep@empresademo.com',                 'EMP', 'AGORA-1', 'AGORA', TRUE),
  (2, '1016060113', 'Egresado Demo',              'egresado@correo.udistrital.edu.co',   'EGR', 'SGA-2',   'SGA',   TRUE),
  (3, '79999999',   'Admin Demo',                 'admin@udistrital.edu.co',             'ADM', NULL,      'LOCAL', TRUE);

-- ── Empresa (llega ya aprobada desde Ágora; estado de ciclo de vida local) ──
INSERT INTO empresa (id, nit, razon_social, agora_id_externo, sector_economico_id, estado_empresa_id, correo_contacto, telefono_contacto, activo) VALUES
  (1, '900111222-3', 'Empresa Demo S.A.S.', 'AG-PROV-1', 1, 1, 'contacto@empresademo.com', '6017000000', TRUE);

-- ── Vínculo usuario(EMP) ↔ empresa (C-7: tipo_usuario fijado a 'EMP') ──
INSERT INTO usuario_empresa (id, usuario_id, tipo_usuario, empresa_id, cargo, es_principal, activo) VALUES
  (1, 1, 'EMP', 1, 'Gerente de Talento Humano', TRUE, TRUE);

-- ── Egresado (subtipo EXCLUSIVO de usuario; tipo_usuario fijado a 'EGR') ──
INSERT INTO egresado (id, usuario_id, tipo_usuario, codigo_institucional, programa_academico, facultad, fecha_grado, telefono_contacto, activo) VALUES
  (1, 2, 'EGR', '20201020113', 'Ingeniería de Sistemas', 'Facultad de Ingeniería', DATE '2024-12-01', '3001234567', TRUE);

-- ── Reiniciar las secuencias SERIAL al máximo insertado ──
SELECT setval(pg_get_serial_sequence('beneficios_egresados.usuario',         'id'), (SELECT MAX(id) FROM usuario));
SELECT setval(pg_get_serial_sequence('beneficios_egresados.empresa',         'id'), (SELECT MAX(id) FROM empresa));
SELECT setval(pg_get_serial_sequence('beneficios_egresados.usuario_empresa', 'id'), (SELECT MAX(id) FROM usuario_empresa));
SELECT setval(pg_get_serial_sequence('beneficios_egresados.egresado',        'id'), (SELECT MAX(id) FROM egresado));

-- Verificación rápida
SELECT 'usuario'  AS tabla, COUNT(*) FROM usuario
UNION ALL SELECT 'empresa',         COUNT(*) FROM empresa
UNION ALL SELECT 'usuario_empresa', COUNT(*) FROM usuario_empresa
UNION ALL SELECT 'egresado',        COUNT(*) FROM egresado;
