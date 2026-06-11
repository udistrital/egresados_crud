-- =============================================================
-- SGA Beneficios Egresados — DDL PostgreSQL v2
-- Schema: beneficios_egresados
-- =============================================================
-- CAMBIOS RESPECTO A v1:
--   · Tablas paramétricas eliminadas (tipo_usuario, estado_empresa,
--     estado_beneficio, estado_solicitud, categoria_beneficio,
--     sector_economico, parametro_sistema) → sustituidas por FK
--     a parametro.parametro (schema compartido, mismo clúster PG).
--   · historial_estado_solicitud eliminada; el estado vigente de
--     una solicitud es el ÚLTIMO registro en historial_solicitud
--     (renombrada). solicitud_beneficio ya NO lleva estado_id propio.
--   · Semilla migrada: los INSERT ahora apuntan a parametro/tipo_parametro.
-- =============================================================

CREATE SCHEMA IF NOT EXISTS beneficios_egresados;
SET search_path TO beneficios_egresados;

-- =============================================================
-- NOTA CROSS-SCHEMA
-- Las tablas parametro.parametro y parametro.tipo_parametro
-- viven en el mismo clúster PostgreSQL pero en schema distinto.
-- Las FK cross-schema se declaran con nombre completo de tabla.
-- Si se prefiere evitar FK cross-schema, reemplazar los REFERENCES
-- por un CHECK o validar en capa de aplicación.
-- =============================================================


-- -------------------------------------------------------------
-- Usuarios y perfiles
-- -------------------------------------------------------------

-- tipo_usuario_id → parametro.parametro.id
--   tipo_parametro.codigo_abreviacion = 'TIPO_USUARIO'
--   Valores semilla: EGR, EMP, ADM

CREATE TABLE usuario (
    id                  SERIAL          NOT NULL,
    documento           VARCHAR(20)     NOT NULL,
    nombre              VARCHAR(200)    NOT NULL,
    correo              VARCHAR(150)    NOT NULL,
    tipo_usuario_id     INTEGER         NOT NULL,   -- FK → parametro.parametro
    id_externo          VARCHAR(50),
    sistema_origen      VARCHAR(20)     NOT NULL,
    ultimo_acceso       TIMESTAMP,
    activo              BOOLEAN         NOT NULL DEFAULT TRUE,
    fecha_creacion      TIMESTAMP       NOT NULL DEFAULT NOW(),
    fecha_modificacion  TIMESTAMP       NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_usuario PRIMARY KEY (id),
    CONSTRAINT uq_documento_usuario UNIQUE (documento),
    CONSTRAINT ck_sistema_origen_usuario CHECK (sistema_origen IN ('SGA','AGORA','LOCAL')),
    CONSTRAINT fk_usuario_tipo_usuario
        FOREIGN KEY (tipo_usuario_id) REFERENCES parametro.parametro(id)
);
COMMENT ON TABLE usuario IS
    'Identidad local creada vía JIT provisioning al autenticarse contra SGA (egresados) '
    'o Ágora (empresas). tipo_usuario_id referencia parametro con tipo TIPO_USUARIO.';
COMMENT ON COLUMN usuario.tipo_usuario_id IS
    'Referencia a parametro.parametro; tipo_parametro.codigo_abreviacion = ''TIPO_USUARIO''. '
    'Valores esperados: EGR (Egresado), EMP (Empresa), ADM (Administrador).';

CREATE INDEX idx_usuario_tipo_usuario ON usuario(tipo_usuario_id);
CREATE INDEX idx_usuario_id_externo   ON usuario(sistema_origen, id_externo);


-- -------------------------------------------------------------
-- Egresados
-- -------------------------------------------------------------

-- Datos académicos vienen desde el Sistema de Gestión Académica (SGA).
-- Esta tabla actúa como espejo local para no depender en tiempo real del SGA.

CREATE TABLE egresado (
    id                      SERIAL          NOT NULL,
    usuario_id              INTEGER         NOT NULL,
    codigo_institucional    VARCHAR(20)     NOT NULL,
    programa_academico      VARCHAR(150),
    facultad                VARCHAR(150),
    fecha_grado             DATE,
    telefono_contacto       VARCHAR(20),
    activo                  BOOLEAN         NOT NULL DEFAULT TRUE,
    fecha_creacion          TIMESTAMP       NOT NULL DEFAULT NOW(),
    fecha_modificacion      TIMESTAMP       NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_egresado PRIMARY KEY (id),
    CONSTRAINT uq_usuario_id_egresado UNIQUE (usuario_id),
    CONSTRAINT uq_codigo_institucional_egresado UNIQUE (codigo_institucional),
    CONSTRAINT fk_egresado_usuario FOREIGN KEY (usuario_id) REFERENCES usuario(id)
);
COMMENT ON TABLE egresado IS
    'Perfil de egresado (1:1 con usuario tipo EGR). '
    'Datos académicos sincronizados desde SGA.';


-- -------------------------------------------------------------
-- Empresas
-- -------------------------------------------------------------

-- sector_economico_id → parametro.parametro.id
--   tipo_parametro.codigo_abreviacion = 'SECTOR_ECONOMICO'
-- estado_empresa_id   → parametro.parametro.id
--   tipo_parametro.codigo_abreviacion = 'ESTADO_EMPRESA'
--   Valores semilla: EN_REVISION, APROBADA, RECHAZADA, SUSPENDIDA

CREATE TABLE empresa (
    id                      SERIAL          NOT NULL,
    nit                     VARCHAR(20)     NOT NULL,
    razon_social            VARCHAR(200)    NOT NULL,
    agora_id_externo        VARCHAR(50),
    sector_economico_id     INTEGER,                -- FK → parametro.parametro
    estado_empresa_id       INTEGER         NOT NULL, -- FK → parametro.parametro
    sitio_web               VARCHAR(255),
    correo_contacto         VARCHAR(150),
    telefono_contacto       VARCHAR(20),
    direccion               VARCHAR(255),
    fecha_aprobacion        TIMESTAMP,
    usuario_aprobador_id    INTEGER,
    activo                  BOOLEAN         NOT NULL DEFAULT TRUE,
    fecha_creacion          TIMESTAMP       NOT NULL DEFAULT NOW(),
    fecha_modificacion      TIMESTAMP       NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_empresa PRIMARY KEY (id),
    CONSTRAINT uq_nit_empresa UNIQUE (nit),
    CONSTRAINT fk_empresa_sector_economico
        FOREIGN KEY (sector_economico_id)  REFERENCES parametro.parametro(id),
    CONSTRAINT fk_empresa_estado_empresa
        FOREIGN KEY (estado_empresa_id)    REFERENCES parametro.parametro(id),
    CONSTRAINT fk_empresa_usuario_aprobador
        FOREIGN KEY (usuario_aprobador_id) REFERENCES usuario(id)
);
COMMENT ON TABLE empresa IS
    'Empresa aliada. Espejo local de Ágora + estado de ciclo de vida propio del módulo. '
    'sector_economico_id y estado_empresa_id referencian parametro.parametro.';
COMMENT ON COLUMN empresa.sector_economico_id IS
    'Referencia a parametro.parametro; tipo_parametro.codigo_abreviacion = ''SECTOR_ECONOMICO''.';
COMMENT ON COLUMN empresa.estado_empresa_id IS
    'Referencia a parametro.parametro; tipo_parametro.codigo_abreviacion = ''ESTADO_EMPRESA''. '
    'Valores esperados: EN_REVISION, APROBADA, RECHAZADA, SUSPENDIDA.';

CREATE INDEX idx_empresa_estado_empresa   ON empresa(estado_empresa_id);
CREATE INDEX idx_empresa_sector_economico ON empresa(sector_economico_id);


CREATE TABLE usuario_empresa (
    id           SERIAL      NOT NULL,
    usuario_id   INTEGER     NOT NULL,
    empresa_id   INTEGER     NOT NULL,
    cargo        VARCHAR(100),
    es_principal BOOLEAN     NOT NULL DEFAULT FALSE,
    activo       BOOLEAN     NOT NULL DEFAULT TRUE,
    fecha_creacion      TIMESTAMP NOT NULL DEFAULT NOW(),
    fecha_modificacion  TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_usuario_empresa PRIMARY KEY (id),
    CONSTRAINT uq_usuario_empresa UNIQUE (usuario_id, empresa_id),
    CONSTRAINT fk_usuario_empresa_usuario FOREIGN KEY (usuario_id) REFERENCES usuario(id),
    CONSTRAINT fk_usuario_empresa_empresa FOREIGN KEY (empresa_id) REFERENCES empresa(id)
);
COMMENT ON TABLE usuario_empresa IS
    'Relación N:M entre usuarios (tipo EMP) y empresas. Lógica de asignación validada con Ágora.';

CREATE INDEX idx_usuario_empresa_empresa ON usuario_empresa(empresa_id);


-- -------------------------------------------------------------
-- Beneficios
-- -------------------------------------------------------------

-- categoria_beneficio_id → parametro.parametro.id
--   tipo_parametro.codigo_abreviacion = 'CATEGORIA_BENEFICIO'
--   Valores semilla: Educación, Salud, Recreación, Empleo, Descuentos, etc.
-- estado_beneficio_id    → parametro.parametro.id
--   tipo_parametro.codigo_abreviacion = 'ESTADO_BENEFICIO'
--   Valores semilla: BORRADOR, PUBLICADO, AGOTADO, VENCIDO, RETIRADO

CREATE TABLE beneficio (
    id                      SERIAL       NOT NULL,
    empresa_id              INTEGER      NOT NULL,
    categoria_beneficio_id  INTEGER      NOT NULL, -- FK → parametro.parametro
    estado_beneficio_id     INTEGER      NOT NULL, -- FK → parametro.parametro
    titulo                  VARCHAR(200) NOT NULL,
    descripcion             TEXT         NOT NULL,
    condiciones             TEXT         NOT NULL,
    fecha_inicio            DATE         NOT NULL,
    fecha_fin               DATE         NOT NULL,
    cupos_total             INTEGER      NOT NULL,
    cupos_disponibles       INTEGER      NOT NULL,
    imagen_url              VARCHAR(500),
    fecha_publicacion       TIMESTAMP,
    usuario_creador_id      INTEGER      NOT NULL,
    activo                  BOOLEAN      NOT NULL DEFAULT TRUE,
    fecha_creacion          TIMESTAMP    NOT NULL DEFAULT NOW(),
    fecha_modificacion      TIMESTAMP    NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_beneficio PRIMARY KEY (id),
    CONSTRAINT ck_fechas_beneficio
        CHECK (fecha_fin > fecha_inicio),
    CONSTRAINT ck_cupos_total_beneficio
        CHECK (cupos_total >= 1),
    CONSTRAINT ck_cupos_disponibles_beneficio
        CHECK (cupos_disponibles >= 0 AND cupos_disponibles <= cupos_total),
    CONSTRAINT fk_beneficio_empresa
        FOREIGN KEY (empresa_id)             REFERENCES empresa(id),
    CONSTRAINT fk_beneficio_categoria
        FOREIGN KEY (categoria_beneficio_id) REFERENCES parametro.parametro(id),
    CONSTRAINT fk_beneficio_estado
        FOREIGN KEY (estado_beneficio_id)    REFERENCES parametro.parametro(id),
    CONSTRAINT fk_beneficio_usuario_creador
        FOREIGN KEY (usuario_creador_id)     REFERENCES usuario(id)
);
COMMENT ON TABLE beneficio IS
    'Beneficio publicado por una empresa aliada para egresados. '
    'categoria_beneficio_id y estado_beneficio_id referencian parametro.parametro.';
COMMENT ON COLUMN beneficio.categoria_beneficio_id IS
    'Referencia a parametro.parametro; tipo_parametro.codigo_abreviacion = ''CATEGORIA_BENEFICIO''.';
COMMENT ON COLUMN beneficio.estado_beneficio_id IS
    'Referencia a parametro.parametro; tipo_parametro.codigo_abreviacion = ''ESTADO_BENEFICIO''. '
    'Valores esperados: BORRADOR, PUBLICADO, AGOTADO, VENCIDO, RETIRADO.';

CREATE INDEX idx_beneficio_estado    ON beneficio(estado_beneficio_id);
CREATE INDEX idx_beneficio_categoria ON beneficio(categoria_beneficio_id);
CREATE INDEX idx_beneficio_empresa   ON beneficio(empresa_id);
CREATE INDEX idx_beneficio_vigencia  ON beneficio(fecha_fin, cupos_disponibles);


-- -------------------------------------------------------------
-- Secuencia de radicado
-- (Se conserva por buenas prácticas: garantía de unicidad y
--  control transaccional con SELECT FOR UPDATE)
-- -------------------------------------------------------------

CREATE TABLE secuencia_radicado (
    id              SERIAL    NOT NULL,
    anio            INTEGER   NOT NULL,
    ultimo_numero   INTEGER   NOT NULL DEFAULT 0,
    activo          BOOLEAN   NOT NULL DEFAULT TRUE,
    fecha_creacion      TIMESTAMP NOT NULL DEFAULT NOW(),
    fecha_modificacion  TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_secuencia_radicado  PRIMARY KEY (id),
    CONSTRAINT uq_anio_secuencia_radicado UNIQUE (anio)
);
COMMENT ON TABLE secuencia_radicado IS
    'Contador para generar radicados BNF-YYYY-NNNNNN. '
    'Usar SELECT FOR UPDATE al asignar número para evitar duplicados en concurrencia.';


-- -------------------------------------------------------------
-- Solicitudes de beneficio
-- -------------------------------------------------------------

-- El estado vigente de una solicitud es el ÚLTIMO registro
-- en historial_solicitud (ORDER BY fecha_cambio DESC LIMIT 1).
-- Se elimina estado_solicitud_id de esta tabla para evitar
-- redundancia y posible inconsistencia.

CREATE TABLE solicitud_beneficio (
    id                    SERIAL      NOT NULL,
    radicado              VARCHAR(20) NOT NULL,
    egresado_id           INTEGER     NOT NULL,
    beneficio_id          INTEGER     NOT NULL,
    datos_complementarios JSONB,
    fecha_solicitud       TIMESTAMP   NOT NULL DEFAULT NOW(),
    activo                BOOLEAN     NOT NULL DEFAULT TRUE,
    fecha_creacion        TIMESTAMP   NOT NULL DEFAULT NOW(),
    fecha_modificacion    TIMESTAMP   NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_solicitud_beneficio PRIMARY KEY (id),
    CONSTRAINT uq_radicado_solicitud_beneficio UNIQUE (radicado),
    CONSTRAINT ck_radicado_solicitud_beneficio
        CHECK (radicado ~ '^BNF-[0-9]{4}-[0-9]{6}$'),
    CONSTRAINT fk_solicitud_beneficio_egresado
        FOREIGN KEY (egresado_id)  REFERENCES egresado(id),
    CONSTRAINT fk_solicitud_beneficio_beneficio
        FOREIGN KEY (beneficio_id) REFERENCES beneficio(id)
);
COMMENT ON TABLE solicitud_beneficio IS
    'Solicitud de un egresado sobre un beneficio. '
    'El estado vigente se obtiene del último registro en historial_solicitud.';

CREATE INDEX idx_solicitud_beneficio_egresado  ON solicitud_beneficio(egresado_id);
CREATE INDEX idx_solicitud_beneficio_beneficio ON solicitud_beneficio(beneficio_id);
CREATE INDEX idx_solicitud_beneficio_fecha     ON solicitud_beneficio(fecha_solicitud);


-- -------------------------------------------------------------
-- Historial de estado de solicitud
-- (Tabla unificada: el último registro es el estado vigente)
-- -------------------------------------------------------------

-- estado_anterior_id / estado_nuevo_id → parametro.parametro.id
--   tipo_parametro.codigo_abreviacion = 'ESTADO_SOLICITUD'
--   Valores semilla: PENDIENTE, EN_REVISION, REQUIERE_INFO,
--                    APROBADA, RECHAZADA, CANCELADA

CREATE TABLE historial_solicitud (
    id                      SERIAL    NOT NULL,
    solicitud_beneficio_id  INTEGER   NOT NULL,
    estado_anterior_id      INTEGER,             -- FK → parametro.parametro (NULL en estado inicial)
    estado_nuevo_id         INTEGER   NOT NULL,  -- FK → parametro.parametro
    usuario_id              INTEGER   NOT NULL,
    justificacion           TEXT,
    fecha_cambio            TIMESTAMP NOT NULL DEFAULT NOW(),
    activo                  BOOLEAN   NOT NULL DEFAULT TRUE,
    fecha_creacion          TIMESTAMP NOT NULL DEFAULT NOW(),
    fecha_modificacion      TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_historial_solicitud PRIMARY KEY (id),
    CONSTRAINT fk_historial_solicitud_solicitud
        FOREIGN KEY (solicitud_beneficio_id) REFERENCES solicitud_beneficio(id),
    CONSTRAINT fk_historial_estado_anterior
        FOREIGN KEY (estado_anterior_id) REFERENCES parametro.parametro(id),
    CONSTRAINT fk_historial_estado_nuevo
        FOREIGN KEY (estado_nuevo_id)    REFERENCES parametro.parametro(id),
    CONSTRAINT fk_historial_usuario
        FOREIGN KEY (usuario_id) REFERENCES usuario(id)
);
COMMENT ON TABLE historial_solicitud IS
    'Bitácora de transiciones de estado de cada solicitud. '
    'El estado vigente de una solicitud es el registro con mayor fecha_cambio. '
    'estado_anterior_id es NULL en el primer registro (creación). '
    'estado_anterior_id y estado_nuevo_id referencian parametro.parametro '
    'con tipo_parametro.codigo_abreviacion = ''ESTADO_SOLICITUD''.';
COMMENT ON COLUMN historial_solicitud.estado_anterior_id IS
    'NULL en la creación de la solicitud (no hay estado previo).';
COMMENT ON COLUMN historial_solicitud.estado_nuevo_id IS
    'Referencia a parametro.parametro; tipo_parametro.codigo_abreviacion = ''ESTADO_SOLICITUD''.';

CREATE INDEX idx_historial_solicitud_solicitud ON historial_solicitud(solicitud_beneficio_id);
CREATE INDEX idx_historial_solicitud_fecha     ON historial_solicitud(fecha_cambio);
-- Índice para obtener el estado vigente eficientemente
CREATE INDEX idx_historial_solicitud_vigente
    ON historial_solicitud(solicitud_beneficio_id, fecha_cambio DESC);


-- -------------------------------------------------------------
-- Mensajes de solicitud
-- -------------------------------------------------------------

CREATE TABLE mensaje_solicitud (
    id                      SERIAL    NOT NULL,
    solicitud_beneficio_id  INTEGER   NOT NULL,
    usuario_id              INTEGER   NOT NULL,
    mensaje                 TEXT      NOT NULL,
    fecha_envio             TIMESTAMP NOT NULL DEFAULT NOW(),
    activo                  BOOLEAN   NOT NULL DEFAULT TRUE,
    fecha_creacion          TIMESTAMP NOT NULL DEFAULT NOW(),
    fecha_modificacion      TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_mensaje_solicitud PRIMARY KEY (id),
    CONSTRAINT fk_mensaje_solicitud_solicitud
        FOREIGN KEY (solicitud_beneficio_id) REFERENCES solicitud_beneficio(id),
    CONSTRAINT fk_mensaje_solicitud_usuario
        FOREIGN KEY (usuario_id)             REFERENCES usuario(id)
);
COMMENT ON TABLE mensaje_solicitud IS
    'Intercambio empresa ↔ egresado cuando una solicitud está en estado REQUIERE_INFO.';

CREATE INDEX idx_mensaje_solicitud_solicitud ON mensaje_solicitud(solicitud_beneficio_id, fecha_envio);


-- -------------------------------------------------------------
-- Bitácora PII (registro inmutable, sin FK de borrado)
-- -------------------------------------------------------------

CREATE TABLE bitacora_acceso_pii (
    id           SERIAL      NOT NULL,
    usuario_id   INTEGER     NOT NULL,
    recurso_tipo VARCHAR(50) NOT NULL,
    recurso_id   INTEGER,
    accion       VARCHAR(50) NOT NULL,
    direccion_ip VARCHAR(45),
    user_agent   VARCHAR(500),
    detalle      JSONB,
    fecha_evento TIMESTAMP   NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_bitacora_acceso_pii PRIMARY KEY (id),
    CONSTRAINT fk_bitacora_acceso_pii_usuario FOREIGN KEY (usuario_id) REFERENCES usuario(id)
);
COMMENT ON TABLE bitacora_acceso_pii IS
    'Bitácora inmutable de accesos a datos personales (Ley 1581 Colombia). '
    'Retención mínima 6 meses. No se aplica borrado lógico (activo) en esta tabla.';

CREATE INDEX idx_bitacora_pii_usuario ON bitacora_acceso_pii(usuario_id);
CREATE INDEX idx_bitacora_pii_recurso ON bitacora_acceso_pii(recurso_tipo, recurso_id);
CREATE INDEX idx_bitacora_pii_fecha   ON bitacora_acceso_pii(fecha_evento);


-- =============================================================
-- VISTA DE APOYO: estado vigente de solicitudes
-- =============================================================

CREATE OR REPLACE VIEW v_solicitud_estado_vigente AS
SELECT DISTINCT ON (hs.solicitud_beneficio_id)
    sb.id                    AS solicitud_id,
    sb.radicado,
    sb.egresado_id,
    sb.beneficio_id,
    hs.estado_nuevo_id       AS estado_actual_id,
    hs.fecha_cambio          AS fecha_ultimo_estado,
    hs.usuario_id            AS usuario_ultimo_cambio,
    hs.justificacion         AS justificacion_ultimo_cambio
FROM solicitud_beneficio sb
JOIN historial_solicitud hs ON hs.solicitud_beneficio_id = sb.id
ORDER BY hs.solicitud_beneficio_id, hs.fecha_cambio DESC;

COMMENT ON VIEW v_solicitud_estado_vigente IS
    'Estado vigente de cada solicitud (último registro en historial_solicitud). '
    'Usar esta vista en consultas de listado para evitar repetir la lógica ORDER BY / DISTINCT ON.';


-- =============================================================
-- DATOS SEMILLA — parametro.tipo_parametro y parametro.parametro
-- =============================================================
-- Ejecutar en el schema "parametro" (mismo clúster).
-- Se asume que area_tipo_id corresponde al área de Egresados;
-- ajustar el valor según el id real en producción.
-- =============================================================

SET search_path TO parametro;

-- -----------------------------------------------------------
-- tipo_parametro: grupos de parámetros para este módulo
-- -----------------------------------------------------------

INSERT INTO tipo_parametro (nombre, descripcion, codigo_abreviacion, activo, numero_orden)
VALUES
    ('Tipo de Usuario',       'Tipos de usuario del módulo de beneficios (Egresado, Empresa, Administrador)',          'TIPO_USUARIO',        TRUE, 1),
    ('Estado Empresa',        'Estados del ciclo de vida de una empresa aliada',                                       'ESTADO_EMPRESA',      TRUE, 2),
    ('Estado Beneficio',      'Estados de un beneficio publicado',                                                     'ESTADO_BENEFICIO',    TRUE, 3),
    ('Estado Solicitud',      'Estados de una solicitud de beneficio',                                                 'ESTADO_SOLICITUD',    TRUE, 4),
    ('Categoría Beneficio',   'Categorías de clasificación de beneficios para egresados',                              'CATEGORIA_BENEFICIO', TRUE, 5),
    ('Sector Económico',      'Sectores económicos para clasificación de empresas aliadas',                            'SECTOR_ECONOMICO',    TRUE, 6),
    ('Parámetro Sistema',     'Parámetros configurables del módulo de beneficios (límites, paginación, validaciones)', 'PARAMETRO_SISTEMA',   TRUE, 7);


-- -----------------------------------------------------------
-- parametro: valores de cada tipo
-- (tipo_parametro_id se resuelve por subconsulta de codigo_abreviacion)
-- -----------------------------------------------------------

-- TIPO_USUARIO
INSERT INTO parametro (nombre, descripcion, codigo_abreviacion, activo, numero_orden, tipo_parametro_id)
SELECT nombre, descripcion, codigo_abreviacion, activo, numero_orden,
       (SELECT id FROM tipo_parametro WHERE codigo_abreviacion = 'TIPO_USUARIO')
FROM (VALUES
    ('Egresado',      'Egresado autenticado vía SGA',                           'EGR', TRUE, 1),
    ('Empresa',       'Representante de empresa autenticado vía Ágora',          'EMP', TRUE, 2),
    ('Administrador', 'Administrador del módulo (OATI / Oficina de Egresados)',  'ADM', TRUE, 3)
) AS v(nombre, descripcion, codigo_abreviacion, activo, numero_orden);

-- ESTADO_EMPRESA
INSERT INTO parametro (nombre, descripcion, codigo_abreviacion, activo, numero_orden, tipo_parametro_id)
SELECT nombre, descripcion, codigo_abreviacion, activo, numero_orden,
       (SELECT id FROM tipo_parametro WHERE codigo_abreviacion = 'ESTADO_EMPRESA')
FROM (VALUES
    ('En revisión', 'Empresa en proceso de validación',                        'EN_REVISION', TRUE, 1),
    ('Aprobada',    'Empresa aprobada para publicar beneficios',               'APROBADA',    TRUE, 2),
    ('Rechazada',   'Empresa rechazada en el proceso de validación',           'RECHAZADA',   TRUE, 3),
    ('Suspendida',  'Empresa suspendida temporalmente del módulo',             'SUSPENDIDA',  TRUE, 4)
) AS v(nombre, descripcion, codigo_abreviacion, activo, numero_orden);

-- ESTADO_BENEFICIO
INSERT INTO parametro (nombre, descripcion, codigo_abreviacion, activo, numero_orden, tipo_parametro_id)
SELECT nombre, descripcion, codigo_abreviacion, activo, numero_orden,
       (SELECT id FROM tipo_parametro WHERE codigo_abreviacion = 'ESTADO_BENEFICIO')
FROM (VALUES
    ('Borrador',  'Beneficio en edición, no visible para egresados',           'BORRADOR',  TRUE, 1),
    ('Publicado', 'Beneficio activo y visible en el catálogo',                 'PUBLICADO', TRUE, 2),
    ('Agotado',   'Beneficio sin cupos disponibles',                           'AGOTADO',   TRUE, 3),
    ('Vencido',   'Beneficio fuera de su periodo de vigencia',                 'VENCIDO',   TRUE, 4),
    ('Retirado',  'Beneficio retirado manualmente por la empresa o el admin',  'RETIRADO',  TRUE, 5)
) AS v(nombre, descripcion, codigo_abreviacion, activo, numero_orden);

-- ESTADO_SOLICITUD
INSERT INTO parametro (nombre, descripcion, codigo_abreviacion, activo, numero_orden, tipo_parametro_id)
SELECT nombre, descripcion, codigo_abreviacion, activo, numero_orden,
       (SELECT id FROM tipo_parametro WHERE codigo_abreviacion = 'ESTADO_SOLICITUD')
FROM (VALUES
    ('Pendiente',            'Solicitud recibida, sin revisión aún',                         'PENDIENTE',     TRUE, 1),
    ('En revisión',          'Solicitud siendo evaluada por la empresa',                     'EN_REVISION',   TRUE, 2),
    ('Requiere información', 'Empresa solicita datos adicionales al egresado',               'REQUIERE_INFO', TRUE, 3),
    ('Aprobada',             'Solicitud aprobada por la empresa',                            'APROBADA',      TRUE, 4),
    ('Rechazada',            'Solicitud rechazada por la empresa',                           'RECHAZADA',     TRUE, 5),
    ('Cancelada',            'Solicitud cancelada por el egresado o el administrador',       'CANCELADA',     TRUE, 6)
) AS v(nombre, descripcion, codigo_abreviacion, activo, numero_orden);

-- CATEGORIA_BENEFICIO
INSERT INTO parametro (nombre, descripcion, codigo_abreviacion, activo, numero_orden, tipo_parametro_id)
SELECT nombre, descripcion, codigo_abreviacion, activo, numero_orden,
       (SELECT id FROM tipo_parametro WHERE codigo_abreviacion = 'CATEGORIA_BENEFICIO')
FROM (VALUES
    ('Educación',   'Beneficios relacionados con formación académica y capacitación', 'EDUCACION',   TRUE, 1),
    ('Salud',       'Beneficios en servicios de salud y bienestar',                   'SALUD',       TRUE, 2),
    ('Recreación',  'Beneficios recreativos, culturales y deportivos',                'RECREACION',  TRUE, 3),
    ('Empleo',      'Oportunidades laborales y prácticas profesionales',              'EMPLEO',      TRUE, 4),
    ('Descuentos',  'Descuentos en productos y servicios para egresados',             'DESCUENTOS',  TRUE, 5),
    ('Otro',        'Beneficios que no clasifican en las categorías anteriores',      'OTRO',        TRUE, 6)
) AS v(nombre, descripcion, codigo_abreviacion, activo, numero_orden);

-- SECTOR_ECONOMICO
INSERT INTO parametro (nombre, descripcion, codigo_abreviacion, activo, numero_orden, tipo_parametro_id)
SELECT nombre, descripcion, codigo_abreviacion, activo, numero_orden,
       (SELECT id FROM tipo_parametro WHERE codigo_abreviacion = 'SECTOR_ECONOMICO')
FROM (VALUES
    ('Tecnología e Innovación',  'Empresas del sector TIC, software y servicios digitales',         'TEC',  TRUE, 1),
    ('Salud y Farmacéutico',     'Clínicas, hospitales, laboratorios y farmacéuticas',              'SAL',  TRUE, 2),
    ('Educación',                'Instituciones educativas, academias y plataformas e-learning',    'EDU',  TRUE, 3),
    ('Industria y Manufactura',  'Fabricación, producción industrial y manufactura',                'IND',  TRUE, 4),
    ('Comercio y Retail',        'Comercio al por mayor y menor, retail físico y en línea',         'COM',  TRUE, 5),
    ('Servicios Financieros',    'Bancos, aseguradoras, fintech y servicios financieros',           'FIN',  TRUE, 6),
    ('Construcción',             'Constructoras, inmobiliarias e infraestructura',                  'CON',  TRUE, 7),
    ('Alimentos y Bebidas',      'Producción, distribución y venta de alimentos y bebidas',         'ALI',  TRUE, 8),
    ('Consultoría y Servicios',  'Consultoría profesional, outsourcing y servicios empresariales',  'CON2', TRUE, 9),
    ('Otro',                     'Sectores no clasificados en las categorías anteriores',           'OTR',  TRUE, 10)
) AS v(nombre, descripcion, codigo_abreviacion, activo, numero_orden);

-- PARAMETRO_SISTEMA
-- Parámetros configurables del módulo (antes en tabla parametro_sistema)
INSERT INTO parametro (nombre, descripcion, codigo_abreviacion, activo, numero_orden, tipo_parametro_id)
SELECT nombre, descripcion, codigo_abreviacion, activo, numero_orden,
       (SELECT id FROM tipo_parametro WHERE codigo_abreviacion = 'PARAMETRO_SISTEMA')
FROM (VALUES
    ('Límite solicitudes activas egresado',
     'Máximo de solicitudes en estado activo (Pendiente/En revisión/Requiere info) por egresado (RN-010). Valor: 5',
     'LIMITE_SOLICITUDES_ACTIVAS_EGRESADO', TRUE, 1),
    ('Paginación catálogo por defecto',
     'Tamaño de página por defecto del catálogo de beneficios. Valor: 20',
     'PAGINACION_CATALOGO_DEFAULT',         TRUE, 2),
    ('Mínimo caracteres justificación rechazo',
     'Longitud mínima de la justificación al rechazar una solicitud (RN-003). Valor: 20',
     'JUSTIFICACION_RECHAZO_MIN_CARACTERES', TRUE, 3)
) AS v(nombre, descripcion, codigo_abreviacion, activo, numero_orden);


-- -----------------------------------------------------------
-- Volver al schema del módulo
-- -----------------------------------------------------------
SET search_path TO beneficios_egresados;

-- Semilla operativa: inicializar secuencia del año en curso
INSERT INTO secuencia_radicado (anio, ultimo_numero)
    VALUES (EXTRACT(YEAR FROM NOW())::INTEGER, 0);
