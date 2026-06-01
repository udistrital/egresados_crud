-- =============================================================
-- SGA Beneficios Egresados — DDL PostgreSQL
-- Schema: beneficios_egresados
-- =============================================================

CREATE SCHEMA IF NOT EXISTS beneficios_egresados;

SET search_path TO beneficios_egresados;

-- -------------------------------------------------------------
-- Catálogos independientes
-- -------------------------------------------------------------

CREATE TABLE tipo_usuario (
    id                  SERIAL          NOT NULL,
    nombre              VARCHAR(100)    NOT NULL,
    descripcion         VARCHAR(500),
    codigo_abreviacion  VARCHAR(50),
    numero_orden        NUMERIC(5,2),
    activo              BOOLEAN         NOT NULL DEFAULT TRUE,
    fecha_creacion      TIMESTAMP       NOT NULL DEFAULT NOW(),
    fecha_modificacion  TIMESTAMP       NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_tipo_usuario PRIMARY KEY (id),
    CONSTRAINT uq_codigo_abreviacion_tipo_usuario UNIQUE (codigo_abreviacion)
);
COMMENT ON TABLE tipo_usuario IS 'Catálogo de tipos de usuario del sistema (EGRESADO, EMPRESA, ADMINISTRADOR).';

CREATE TABLE estado_empresa (
    id                  SERIAL          NOT NULL,
    nombre              VARCHAR(100)    NOT NULL,
    descripcion         VARCHAR(500),
    codigo_abreviacion  VARCHAR(50),
    numero_orden        NUMERIC(5,2),
    activo              BOOLEAN         NOT NULL DEFAULT TRUE,
    fecha_creacion      TIMESTAMP       NOT NULL DEFAULT NOW(),
    fecha_modificacion  TIMESTAMP       NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_estado_empresa PRIMARY KEY (id),
    CONSTRAINT uq_codigo_abreviacion_estado_empresa UNIQUE (codigo_abreviacion)
);
COMMENT ON TABLE estado_empresa IS 'Estados del ciclo de vida de una empresa aliada (EN_REVISION, APROBADA, RECHAZADA, SUSPENDIDA).';

CREATE TABLE estado_beneficio (
    id                  SERIAL          NOT NULL,
    nombre              VARCHAR(100)    NOT NULL,
    descripcion         VARCHAR(500),
    codigo_abreviacion  VARCHAR(50),
    numero_orden        NUMERIC(5,2),
    activo              BOOLEAN         NOT NULL DEFAULT TRUE,
    fecha_creacion      TIMESTAMP       NOT NULL DEFAULT NOW(),
    fecha_modificacion  TIMESTAMP       NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_estado_beneficio PRIMARY KEY (id),
    CONSTRAINT uq_codigo_abreviacion_estado_beneficio UNIQUE (codigo_abreviacion)
);
COMMENT ON TABLE estado_beneficio IS 'Estados de un beneficio publicado: BORRADOR, PUBLICADO, AGOTADO, VENCIDO, RETIRADO.';

CREATE TABLE estado_solicitud (
    id                  SERIAL          NOT NULL,
    nombre              VARCHAR(100)    NOT NULL,
    descripcion         VARCHAR(500),
    codigo_abreviacion  VARCHAR(50),
    numero_orden        NUMERIC(5,2),
    activo              BOOLEAN         NOT NULL DEFAULT TRUE,
    fecha_creacion      TIMESTAMP       NOT NULL DEFAULT NOW(),
    fecha_modificacion  TIMESTAMP       NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_estado_solicitud PRIMARY KEY (id),
    CONSTRAINT uq_codigo_abreviacion_estado_solicitud UNIQUE (codigo_abreviacion)
);
COMMENT ON TABLE estado_solicitud IS 'Estados de una solicitud de beneficio: PENDIENTE, EN_REVISION, REQUIERE_INFO, APROBADA, RECHAZADA, CANCELADA.';

CREATE TABLE categoria_beneficio (
    id                  SERIAL          NOT NULL,
    nombre              VARCHAR(100)    NOT NULL,
    descripcion         VARCHAR(500),
    codigo_abreviacion  VARCHAR(50),
    numero_orden        NUMERIC(5,2),
    activo              BOOLEAN         NOT NULL DEFAULT TRUE,
    fecha_creacion      TIMESTAMP       NOT NULL DEFAULT NOW(),
    fecha_modificacion  TIMESTAMP       NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_categoria_beneficio PRIMARY KEY (id),
    CONSTRAINT uq_codigo_abreviacion_categoria_beneficio UNIQUE (codigo_abreviacion)
);
COMMENT ON TABLE categoria_beneficio IS 'Categoría de los beneficios (Educación, Salud, Recreación, Empleo, Descuentos, etc.).';

CREATE TABLE sector_economico (
    id                  SERIAL          NOT NULL,
    nombre              VARCHAR(100)    NOT NULL,
    descripcion         VARCHAR(500),
    codigo_abreviacion  VARCHAR(50),
    numero_orden        NUMERIC(5,2),
    activo              BOOLEAN         NOT NULL DEFAULT TRUE,
    fecha_creacion      TIMESTAMP       NOT NULL DEFAULT NOW(),
    fecha_modificacion  TIMESTAMP       NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_sector_economico PRIMARY KEY (id),
    CONSTRAINT uq_codigo_abreviacion_sector_economico UNIQUE (codigo_abreviacion)
);
COMMENT ON TABLE sector_economico IS 'Sectores económicos para clasificación de empresas.';

CREATE TABLE parametro_sistema (
    id                  SERIAL          NOT NULL,
    clave               VARCHAR(100)    NOT NULL,
    valor               VARCHAR(500)    NOT NULL,
    tipo_dato           VARCHAR(20)     NOT NULL,
    descripcion         VARCHAR(500),
    activo              BOOLEAN         NOT NULL DEFAULT TRUE,
    fecha_creacion      TIMESTAMP       NOT NULL DEFAULT NOW(),
    fecha_modificacion  TIMESTAMP       NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_parametro_sistema PRIMARY KEY (id),
    CONSTRAINT uq_clave_parametro_sistema UNIQUE (clave),
    CONSTRAINT ck_tipo_dato_parametro_sistema CHECK (tipo_dato IN ('INTEGER','STRING','BOOLEAN','DECIMAL','JSON'))
);
COMMENT ON TABLE parametro_sistema IS 'Parámetros configurables del sistema (ej. límite de solicitudes activas por egresado).';

-- -------------------------------------------------------------
-- Usuarios y perfiles
-- -------------------------------------------------------------

CREATE TABLE usuario (
    id                  SERIAL          NOT NULL,
    documento           VARCHAR(20)     NOT NULL,
    nombre              VARCHAR(200)    NOT NULL,
    correo              VARCHAR(150)    NOT NULL,
    tipo_usuario_id     INTEGER         NOT NULL,
    id_externo          VARCHAR(50),
    sistema_origen      VARCHAR(20)     NOT NULL,
    ultimo_acceso       TIMESTAMP,
    activo              BOOLEAN         NOT NULL DEFAULT TRUE,
    fecha_creacion      TIMESTAMP       NOT NULL DEFAULT NOW(),
    fecha_modificacion  TIMESTAMP       NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_usuario PRIMARY KEY (id),
    CONSTRAINT uq_documento_usuario UNIQUE (documento),
    CONSTRAINT ck_sistema_origen_usuario CHECK (sistema_origen IN ('SGA','AGORA','LOCAL')),
    CONSTRAINT fk_usuario_tipo_usuario FOREIGN KEY (tipo_usuario_id) REFERENCES tipo_usuario(id)
);
COMMENT ON TABLE usuario IS 'Identidad local creada vía JIT provisioning al autenticarse contra SGA (egresados) o Ágora (empresas).';
CREATE INDEX idx_usuario_tipo_usuario ON usuario(tipo_usuario_id);
CREATE INDEX idx_usuario_id_externo   ON usuario(sistema_origen, id_externo);

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
COMMENT ON TABLE egresado IS 'Perfil de egresado (1:1 con usuario tipo EGRESADO).';

-- -------------------------------------------------------------
-- Empresas
-- -------------------------------------------------------------

CREATE TABLE empresa (
    id                  SERIAL          NOT NULL,
    nit                 VARCHAR(20)     NOT NULL,
    razon_social        VARCHAR(200)    NOT NULL,
    agora_id_externo    VARCHAR(50),
    sector_economico_id INTEGER,
    estado_empresa_id   INTEGER         NOT NULL,
    sitio_web           VARCHAR(255),
    correo_contacto     VARCHAR(150),
    telefono_contacto   VARCHAR(20),
    direccion           VARCHAR(255),
    fecha_aprobacion    TIMESTAMP,
    usuario_aprobador_id INTEGER,
    activo              BOOLEAN         NOT NULL DEFAULT TRUE,
    fecha_creacion      TIMESTAMP       NOT NULL DEFAULT NOW(),
    fecha_modificacion  TIMESTAMP       NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_empresa PRIMARY KEY (id),
    CONSTRAINT uq_nit_empresa UNIQUE (nit),
    CONSTRAINT fk_empresa_sector_economico  FOREIGN KEY (sector_economico_id)  REFERENCES sector_economico(id),
    CONSTRAINT fk_empresa_estado_empresa    FOREIGN KEY (estado_empresa_id)    REFERENCES estado_empresa(id),
    CONSTRAINT fk_empresa_usuario_aprobador FOREIGN KEY (usuario_aprobador_id) REFERENCES usuario(id)
);
COMMENT ON TABLE empresa IS 'Empresa aliada. Espejo local de Ágora + estado de ciclo de vida propio del módulo.';
CREATE INDEX idx_empresa_estado_empresa  ON empresa(estado_empresa_id);
CREATE INDEX idx_empresa_sector_economico ON empresa(sector_economico_id);

CREATE TABLE usuario_empresa (
    id          SERIAL      NOT NULL,
    usuario_id  INTEGER     NOT NULL,
    empresa_id  INTEGER     NOT NULL,
    cargo       VARCHAR(100),
    es_principal BOOLEAN    NOT NULL DEFAULT FALSE,
    activo      BOOLEAN     NOT NULL DEFAULT TRUE,
    fecha_creacion      TIMESTAMP NOT NULL DEFAULT NOW(),
    fecha_modificacion  TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_usuario_empresa PRIMARY KEY (id),
    CONSTRAINT uq_usuario_empresa UNIQUE (usuario_id, empresa_id),
    CONSTRAINT fk_usuario_empresa_usuario  FOREIGN KEY (usuario_id)  REFERENCES usuario(id),
    CONSTRAINT fk_usuario_empresa_empresa  FOREIGN KEY (empresa_id)  REFERENCES empresa(id)
);
COMMENT ON TABLE usuario_empresa IS 'Relación N:M entre usuarios (tipo EMPRESA) y empresas.';
CREATE INDEX idx_usuario_empresa_empresa ON usuario_empresa(empresa_id);

-- -------------------------------------------------------------
-- Beneficios y solicitudes
-- -------------------------------------------------------------

CREATE TABLE beneficio (
    id                      SERIAL      NOT NULL,
    empresa_id              INTEGER     NOT NULL,
    categoria_beneficio_id  INTEGER     NOT NULL,
    estado_beneficio_id     INTEGER     NOT NULL,
    titulo                  VARCHAR(200) NOT NULL,
    descripcion             TEXT        NOT NULL,
    condiciones             TEXT        NOT NULL,
    fecha_inicio            DATE        NOT NULL,
    fecha_fin               DATE        NOT NULL,
    cupos_total             INTEGER     NOT NULL,
    cupos_disponibles       INTEGER     NOT NULL,
    imagen_url              VARCHAR(500),
    fecha_publicacion       TIMESTAMP,
    usuario_creador_id      INTEGER     NOT NULL,
    activo                  BOOLEAN     NOT NULL DEFAULT TRUE,
    fecha_creacion          TIMESTAMP   NOT NULL DEFAULT NOW(),
    fecha_modificacion      TIMESTAMP   NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_beneficio PRIMARY KEY (id),
    CONSTRAINT ck_fechas_beneficio            CHECK (fecha_fin > fecha_inicio),
    CONSTRAINT ck_cupos_total_beneficio       CHECK (cupos_total >= 1),
    CONSTRAINT ck_cupos_disponibles_beneficio CHECK (cupos_disponibles >= 0 AND cupos_disponibles <= cupos_total),
    CONSTRAINT fk_beneficio_empresa           FOREIGN KEY (empresa_id)             REFERENCES empresa(id),
    CONSTRAINT fk_beneficio_categoria         FOREIGN KEY (categoria_beneficio_id) REFERENCES categoria_beneficio(id),
    CONSTRAINT fk_beneficio_estado_beneficio  FOREIGN KEY (estado_beneficio_id)    REFERENCES estado_beneficio(id),
    CONSTRAINT fk_beneficio_usuario_creador   FOREIGN KEY (usuario_creador_id)     REFERENCES usuario(id)
);
COMMENT ON TABLE beneficio IS 'Beneficio publicado por una empresa aliada para egresados.';
CREATE INDEX idx_beneficio_estado_beneficio ON beneficio(estado_beneficio_id);
CREATE INDEX idx_beneficio_categoria        ON beneficio(categoria_beneficio_id);
CREATE INDEX idx_beneficio_empresa          ON beneficio(empresa_id);
CREATE INDEX idx_beneficio_vigencia         ON beneficio(fecha_fin, cupos_disponibles);

CREATE TABLE secuencia_radicado (
    id              SERIAL      NOT NULL,
    anio            INTEGER     NOT NULL,
    ultimo_numero   INTEGER     NOT NULL DEFAULT 0,
    activo          BOOLEAN     NOT NULL DEFAULT TRUE,
    fecha_creacion      TIMESTAMP NOT NULL DEFAULT NOW(),
    fecha_modificacion  TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_secuencia_radicado  PRIMARY KEY (id),
    CONSTRAINT uq_anio_secuencia_radicado UNIQUE (anio)
);
COMMENT ON TABLE secuencia_radicado IS 'Contador para generar radicados BNF-YYYY-NNNNNN. Usar SELECT FOR UPDATE al asignar.';

CREATE TABLE solicitud_beneficio (
    id                  SERIAL      NOT NULL,
    radicado            VARCHAR(20) NOT NULL,
    egresado_id         INTEGER     NOT NULL,
    beneficio_id        INTEGER     NOT NULL,
    estado_solicitud_id INTEGER     NOT NULL,
    datos_complementarios JSONB,
    fecha_solicitud     TIMESTAMP   NOT NULL DEFAULT NOW(),
    activo              BOOLEAN     NOT NULL DEFAULT TRUE,
    fecha_creacion      TIMESTAMP   NOT NULL DEFAULT NOW(),
    fecha_modificacion  TIMESTAMP   NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_solicitud_beneficio       PRIMARY KEY (id),
    CONSTRAINT uq_radicado_solicitud_beneficio UNIQUE (radicado),
    CONSTRAINT ck_radicado_solicitud_beneficio CHECK (radicado ~ '^BNF-[0-9]{4}-[0-9]{6}$'),
    CONSTRAINT fk_solicitud_beneficio_egresado       FOREIGN KEY (egresado_id)         REFERENCES egresado(id),
    CONSTRAINT fk_solicitud_beneficio_beneficio      FOREIGN KEY (beneficio_id)         REFERENCES beneficio(id),
    CONSTRAINT fk_solicitud_beneficio_estado         FOREIGN KEY (estado_solicitud_id)  REFERENCES estado_solicitud(id)
);
COMMENT ON TABLE solicitud_beneficio IS 'Solicitud de un egresado sobre un beneficio.';
CREATE INDEX idx_solicitud_beneficio_egresado              ON solicitud_beneficio(egresado_id);
CREATE INDEX idx_solicitud_beneficio_beneficio             ON solicitud_beneficio(beneficio_id);
CREATE INDEX idx_solicitud_beneficio_estado                ON solicitud_beneficio(estado_solicitud_id);
CREATE INDEX idx_solicitud_beneficio_fecha                 ON solicitud_beneficio(fecha_solicitud);
CREATE INDEX idx_solicitud_beneficio_egresado_beneficio_estado ON solicitud_beneficio(egresado_id, beneficio_id, estado_solicitud_id);

CREATE TABLE historial_estado_solicitud (
    id                      SERIAL      NOT NULL,
    solicitud_beneficio_id  INTEGER     NOT NULL,
    estado_anterior_id      INTEGER,
    estado_nuevo_id         INTEGER     NOT NULL,
    usuario_id              INTEGER     NOT NULL,
    justificacion           TEXT,
    fecha_cambio            TIMESTAMP   NOT NULL DEFAULT NOW(),
    activo                  BOOLEAN     NOT NULL DEFAULT TRUE,
    fecha_creacion          TIMESTAMP   NOT NULL DEFAULT NOW(),
    fecha_modificacion      TIMESTAMP   NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_historial_estado_solicitud PRIMARY KEY (id),
    CONSTRAINT fk_historial_solicitud_beneficio FOREIGN KEY (solicitud_beneficio_id) REFERENCES solicitud_beneficio(id),
    CONSTRAINT fk_historial_estado_anterior     FOREIGN KEY (estado_anterior_id)     REFERENCES estado_solicitud(id),
    CONSTRAINT fk_historial_estado_nuevo        FOREIGN KEY (estado_nuevo_id)        REFERENCES estado_solicitud(id),
    CONSTRAINT fk_historial_usuario             FOREIGN KEY (usuario_id)             REFERENCES usuario(id)
);
COMMENT ON TABLE historial_estado_solicitud IS 'Bitácora de transiciones de estado de cada solicitud.';
CREATE INDEX idx_historial_estado_solicitud_solicitud ON historial_estado_solicitud(solicitud_beneficio_id);
CREATE INDEX idx_historial_estado_solicitud_fecha     ON historial_estado_solicitud(fecha_cambio);

CREATE TABLE mensaje_solicitud (
    id                      SERIAL      NOT NULL,
    solicitud_beneficio_id  INTEGER     NOT NULL,
    usuario_id              INTEGER     NOT NULL,
    mensaje                 TEXT        NOT NULL,
    fecha_envio             TIMESTAMP   NOT NULL DEFAULT NOW(),
    activo                  BOOLEAN     NOT NULL DEFAULT TRUE,
    fecha_creacion          TIMESTAMP   NOT NULL DEFAULT NOW(),
    fecha_modificacion      TIMESTAMP   NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_mensaje_solicitud PRIMARY KEY (id),
    CONSTRAINT fk_mensaje_solicitud_solicitud FOREIGN KEY (solicitud_beneficio_id) REFERENCES solicitud_beneficio(id),
    CONSTRAINT fk_mensaje_solicitud_usuario   FOREIGN KEY (usuario_id)             REFERENCES usuario(id)
);
COMMENT ON TABLE mensaje_solicitud IS 'Intercambio empresa ↔ egresado cuando una solicitud está en estado REQUIERE_INFO.';
CREATE INDEX idx_mensaje_solicitud_solicitud ON mensaje_solicitud(solicitud_beneficio_id, fecha_envio);

-- -------------------------------------------------------------
-- Bitácora PII (sin FK a activo — registro inmutable)
-- -------------------------------------------------------------

CREATE TABLE bitacora_acceso_pii (
    id              SERIAL      NOT NULL,
    usuario_id      INTEGER     NOT NULL,
    recurso_tipo    VARCHAR(50) NOT NULL,
    recurso_id      INTEGER,
    accion          VARCHAR(50) NOT NULL,
    direccion_ip    VARCHAR(45),
    user_agent      VARCHAR(500),
    detalle         JSONB,
    fecha_evento    TIMESTAMP   NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_bitacora_acceso_pii PRIMARY KEY (id),
    CONSTRAINT fk_bitacora_acceso_pii_usuario FOREIGN KEY (usuario_id) REFERENCES usuario(id)
);
COMMENT ON TABLE bitacora_acceso_pii IS 'Bitácora inmutable de accesos a datos personales (Ley 1581). Retención mínima 6 meses.';
CREATE INDEX idx_bitacora_acceso_pii_usuario       ON bitacora_acceso_pii(usuario_id);
CREATE INDEX idx_bitacora_acceso_pii_recurso       ON bitacora_acceso_pii(recurso_tipo, recurso_id);
CREATE INDEX idx_bitacora_acceso_pii_fecha_evento  ON bitacora_acceso_pii(fecha_evento);

-- =============================================================
-- Datos semilla
-- =============================================================

INSERT INTO tipo_usuario (nombre, descripcion, codigo_abreviacion, numero_orden) VALUES
    ('Egresado',       'Egresado autenticado vía SGA',                              'EGR', 1),
    ('Empresa',        'Representante de empresa autenticado vía Ágora',             'EMP', 2),
    ('Administrador',  'Administrador del módulo (OATI / Oficina de Egresados)',     'ADM', 3);

INSERT INTO estado_empresa (nombre, codigo_abreviacion, numero_orden) VALUES
    ('En revisión', 'EN_REVISION', 1),
    ('Aprobada',    'APROBADA',    2),
    ('Rechazada',   'RECHAZADA',   3),
    ('Suspendida',  'SUSPENDIDA',  4);

INSERT INTO estado_beneficio (nombre, codigo_abreviacion, numero_orden) VALUES
    ('Borrador',   'BORRADOR',   1),
    ('Publicado',  'PUBLICADO',  2),
    ('Agotado',    'AGOTADO',    3),
    ('Vencido',    'VENCIDO',    4),
    ('Retirado',   'RETIRADO',   5);

INSERT INTO estado_solicitud (nombre, codigo_abreviacion, numero_orden) VALUES
    ('Pendiente',             'PENDIENTE',      1),
    ('En revisión',           'EN_REVISION',    2),
    ('Requiere información',  'REQUIERE_INFO',  3),
    ('Aprobada',              'APROBADA',       4),
    ('Rechazada',             'RECHAZADA',      5),
    ('Cancelada',             'CANCELADA',      6);

INSERT INTO parametro_sistema (clave, valor, tipo_dato, descripcion) VALUES
    ('LIMITE_SOLICITUDES_ACTIVAS_EGRESADO',  '5',  'INTEGER', 'Máximo de solicitudes en estado activo (Pendiente/En revisión/Requiere info) por egresado (RN-010).'),
    ('PAGINACION_CATALOGO_DEFAULT',          '20', 'INTEGER', 'Tamaño de página por defecto del catálogo de beneficios.'),
    ('JUSTIFICACION_RECHAZO_MIN_CARACTERES', '20', 'INTEGER', 'Longitud mínima de la justificación al rechazar una solicitud (RN-003).');

INSERT INTO secuencia_radicado (anio, ultimo_numero)
    VALUES (EXTRACT(YEAR FROM NOW())::INTEGER, 0);
