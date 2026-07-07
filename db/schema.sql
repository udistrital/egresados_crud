-- =============================================================
-- SGA Beneficios Egresados — DDL PostgreSQL v4 (VERSIÓN FINAL)
-- Schema: beneficios_egresados
-- =============================================================
-- CAMBIOS RESPECTO A v3.1 (Sustentación — correcciones del profesor):
--
--   [C-5] Radicado como objeto de serialización nativo de PostgreSQL:
--     · Se ELIMINA la tabla secuencia_radicado.
--     · Se reemplaza por una SEQUENCE nativa (seq_radicado_beneficio) más
--       una función fn_siguiente_radicado() que arma el BNF-YYYY-NNNNNN.
--     · El campo solicitud_beneficio.radicado usa esa función como DEFAULT.
--     · Esto REVIERTE la decisión C-4a previa (que descartaba SEQUENCE).
--       Motivo: el profesor pidió usar el objeto de serialización que ya
--       existe en PostgreSQL en lugar de un contador manual.
--
--   [C-6] "Virtualización" de la tabla parametro (sacarla del esquema):
--     · El profesor observó que parametro se volvía un hub con demasiadas
--       FK convergiendo (relaciones que "pasaban todas por parametro").
--     · Se ELIMINAN todas las FK cross-schema REFERENCES parametro.parametro.
--     · Las columnas *_id de estado/categoría/sector quedan como
--       REFERENCIAS LÓGICAS (INTEGER planos). La validación de que el id
--       pertenezca al TipoParametro correcto se hace en el MID (capa de
--       negocio), no en la BD. parametro desaparece del diagrama del esquema.
--     · La semilla institucional pasa a un APÉNDICE claramente marcado como
--       NO perteneciente al schema beneficios_egresados.
--
--   [C-7] Relación EXCLUYENTE usuario ↔ egresado (arco exclusivo / subtipos
--         disjuntos):
--     · Un usuario es EGRESADO **o** representante de EMPRESA, nunca ambos.
--     · Se implementa con el patrón declarativo de subtipos disjuntos:
--         - usuario.tipo_usuario pasa a ser un DISCRIMINADOR LOCAL estable
--           VARCHAR(3) CHECK IN ('EGR','EMP','ADM') (antes era tipo_usuario_id
--           hacia parametro; ahora es local porque una FK compuesta necesita
--           un valor comparable en tiempo de DDL).
--         - usuario declara UNIQUE (id, tipo_usuario).
--         - egresado y usuario_empresa llevan una columna tipo_usuario fijada
--           por CHECK ('EGR' y 'EMP' respectivamente) y una FK COMPUESTA
--           (usuario_id, tipo_usuario) -> usuario(id, tipo_usuario).
--       Como un usuario tiene UN solo tipo_usuario, no puede satisfacer a la
--       vez la FK de egresado (exige EGR) y la de usuario_empresa (exige EMP):
--       la exclusividad queda garantizada en DDL puro, sin triggers.
--
-- CAMBIOS HEREDADOS (v2/v3/v3.1):
--   · Catálogos locales eliminados → referencias a parámetros institucionales.
--   · historial_solicitud unificada; estado vigente = último registro.
--   · VARCHAR con longitud; auditoría (activo/fecha_creacion/fecha_modificacion)
--     en todas las tablas.
-- =============================================================

CREATE SCHEMA IF NOT EXISTS beneficios_egresados;
SET search_path TO beneficios_egresados;

-- =============================================================
-- NOTA SOBRE PARÁMETROS INSTITUCIONALES (C-6 — virtualización)
-- Los catálogos de estado/categoría/sector viven en el servicio
-- institucional de parámetros (schema/servicio "parametro" de la OATI).
-- En este esquema NO se declara ninguna FK hacia ese servicio: las
-- columnas *_id son REFERENCIAS LÓGICAS (virtuales). El MID valida que
-- cada id pertenezca al TipoParametro correcto antes de persistir.
-- Ver el apéndice al final con los TipoParametro/Parametro a aprovisionar.
-- =============================================================


-- -------------------------------------------------------------
-- Usuarios y perfiles
-- -------------------------------------------------------------

-- tipo_usuario: discriminador LOCAL estable del subtipo de usuario.
--   Valores: 'EGR' (Egresado), 'EMP' (Empresa), 'ADM' (Administrador).
--   Es local (no referencia parametro) porque participa en el arco
--   exclusivo egresado/usuario_empresa vía FK compuesta (C-7).

CREATE TABLE usuario (
    id                  SERIAL          NOT NULL,
    documento           VARCHAR(20),                -- NULL para empresas self-signup (no tienen cédula); egresados sí lo traen
    nombre              VARCHAR(200)    NOT NULL,
    correo              VARCHAR(150)    NOT NULL,
    tipo_usuario        VARCHAR(3)      NOT NULL,   -- discriminador local (C-7)
    id_externo          VARCHAR(50),
    sistema_origen      VARCHAR(20)     NOT NULL,
    ultimo_acceso       TIMESTAMP,
    activo              BOOLEAN         NOT NULL DEFAULT TRUE,
    fecha_creacion      TIMESTAMP       NOT NULL DEFAULT NOW(),
    fecha_modificacion  TIMESTAMP       NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_usuario PRIMARY KEY (id),
    -- documento único cuando existe; múltiples NULL permitidos (empresas). En Postgres
    -- UNIQUE trata los NULL como distintos, así que las empresas sin documento no colisionan.
    CONSTRAINT uq_documento_usuario UNIQUE (documento),
    -- Identidad de empresa self-signup: se keya por el sub de WSO2 en id_externo.
    CONSTRAINT uq_usuario_id_externo UNIQUE (sistema_origen, id_externo),
    -- Clave compuesta requerida por las FK del arco exclusivo (C-7):
    CONSTRAINT uq_usuario_id_tipo UNIQUE (id, tipo_usuario),
    CONSTRAINT ck_tipo_usuario CHECK (tipo_usuario IN ('EGR','EMP','ADM')),
    CONSTRAINT ck_sistema_origen_usuario CHECK (sistema_origen IN ('SGA','AGORA','LOCAL'))
);
COMMENT ON TABLE usuario IS
    'Identidad local creada vía JIT provisioning al autenticarse contra SGA (egresados) '
    'o Ágora (empresas). tipo_usuario es el discriminador local del subtipo y, junto con '
    'la PK, ancla el arco exclusivo egresado/usuario_empresa (C-7).';
COMMENT ON COLUMN usuario.tipo_usuario IS
    'Discriminador LOCAL del subtipo de usuario: EGR (Egresado), EMP (Empresa), '
    'ADM (Administrador). Es local (no FK a parametro) porque las FK compuestas de '
    'egresado y usuario_empresa lo usan para garantizar exclusividad en DDL puro (C-7). '
    'Conceptualmente equivale al TipoParametro TIPO_USUARIO, pero como solo tiene 3 '
    'valores estructurales estables se modela como CHECK local en vez de referencia virtual.';

CREATE INDEX idx_usuario_tipo_usuario ON usuario(tipo_usuario);
-- (sistema_origen, id_externo) ya es único por la constraint uq_usuario_id_externo,
-- que sirve además como índice de búsqueda para el lookup del JIT de empresa.


-- -------------------------------------------------------------
-- Egresados  (subtipo EXCLUSIVO de usuario — C-7)
-- -------------------------------------------------------------

-- Datos académicos vienen del SGA; esta tabla es espejo local mínimo.
-- La columna tipo_usuario está fijada a 'EGR' y la FK compuesta hacia
-- usuario(id, tipo_usuario) impide que un usuario que sea EMP/ADM tenga
-- perfil de egresado, y a la vez impide que un egresado aparezca como
-- representante de empresa (ver usuario_empresa).

CREATE TABLE egresado (
    id                      SERIAL          NOT NULL,
    usuario_id              INTEGER         NOT NULL,
    tipo_usuario            VARCHAR(3)      NOT NULL DEFAULT 'EGR',  -- fijado por CHECK (C-7)
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
    CONSTRAINT ck_egresado_tipo_usuario CHECK (tipo_usuario = 'EGR'),
    -- FK compuesta: solo enlaza usuarios cuyo tipo_usuario sea EGR (C-7)
    CONSTRAINT fk_egresado_usuario
        FOREIGN KEY (usuario_id, tipo_usuario) REFERENCES usuario(id, tipo_usuario)
);
COMMENT ON TABLE egresado IS
    'Perfil de egresado (1:1 con usuario tipo EGR). Datos académicos sincronizados desde SGA. '
    'Subtipo EXCLUSIVO de usuario (C-7): la FK compuesta (usuario_id, tipo_usuario) hacia '
    'usuario(id, tipo_usuario), con tipo_usuario fijado a EGR por CHECK, impide que un usuario '
    'sea simultáneamente egresado y representante de empresa.';
COMMENT ON COLUMN egresado.tipo_usuario IS
    'Columna discriminadora fijada a ''EGR'' (CHECK). Forma parte de la FK compuesta hacia '
    'usuario(id, tipo_usuario) que materializa el arco exclusivo egresado/usuario_empresa.';


-- -------------------------------------------------------------
-- Empresas
-- -------------------------------------------------------------

-- sector_economico_id / estado_empresa_id: REFERENCIAS LÓGICAS (virtuales)
-- a parámetros institucionales (C-6). NO hay FK; el MID valida el tipo.
--   sector_economico_id  → TipoParametro SECTOR_ECONOMICO
--   estado_empresa_id    → TipoParametro ESTADO_EMPRESA
--                          (ciclo de vida LOCAL: ACTIVA, SUSPENDIDA — las empresas
--                           llegan ya aprobadas desde Ágora, sin aprobación interna)

CREATE TABLE empresa (
    id                      SERIAL          NOT NULL,
    nit                     VARCHAR(20)     NOT NULL,
    razon_social            VARCHAR(200)    NOT NULL,
    agora_id_externo        VARCHAR(50),
    sector_economico_id     INTEGER,                  -- ref. lógica → parametro (C-6)
    estado_empresa_id       INTEGER         NOT NULL, -- ref. lógica → parametro (C-6)
    sitio_web               VARCHAR(255),
    correo_contacto         VARCHAR(150),
    telefono_contacto       VARCHAR(20),
    direccion               VARCHAR(255),
    activo                  BOOLEAN         NOT NULL DEFAULT TRUE,
    fecha_creacion          TIMESTAMP       NOT NULL DEFAULT NOW(),
    fecha_modificacion      TIMESTAMP       NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_empresa PRIMARY KEY (id),
    CONSTRAINT uq_nit_empresa UNIQUE (nit)
    -- Nota: no hay aprobación interna. Las empresas llegan ya validadas desde Ágora
    -- (autenticación WSO2 + JIT), así que se eliminaron usuario_aprobador_id y
    -- fecha_aprobacion. estado_empresa_id queda como ciclo de vida LOCAL (p. ej. SUSPENDIDA).
);
COMMENT ON TABLE empresa IS
    'Empresa aliada. Espejo local de Ágora + estado de ciclo de vida propio del módulo. '
    'sector_economico_id y estado_empresa_id son referencias lógicas a parámetros '
    'institucionales (C-6): sin FK, validadas en el MID.';
COMMENT ON COLUMN empresa.sector_economico_id IS
    'Referencia LÓGICA (virtual) a un Parametro del TipoParametro SECTOR_ECONOMICO. '
    'Sin FK (C-6); el MID valida tipo y existencia. En producción se prefiere el CIIU de Ágora.';
COMMENT ON COLUMN empresa.estado_empresa_id IS
    'Referencia LÓGICA (virtual) a un Parametro del TipoParametro ESTADO_EMPRESA. Sin FK (C-6); '
    'estado de ciclo de vida LOCAL del módulo. Como las empresas llegan ya aprobadas desde Ágora, '
    'no hay estados de aprobación interna (EN_REVISION/RECHAZADA); el uso típico es ACTIVA/SUSPENDIDA.';

CREATE INDEX idx_empresa_estado_empresa   ON empresa(estado_empresa_id);
CREATE INDEX idx_empresa_sector_economico ON empresa(sector_economico_id);


-- -------------------------------------------------------------
-- usuario_empresa  (subtipo EXCLUSIVO de usuario — C-7)
-- -------------------------------------------------------------

-- Relación N:M entre usuarios tipo EMP y empresas. La columna tipo_usuario
-- está fijada a 'EMP' y la FK compuesta impide vincular aquí a un egresado.

CREATE TABLE usuario_empresa (
    id           SERIAL      NOT NULL,
    usuario_id   INTEGER     NOT NULL,
    tipo_usuario VARCHAR(3)  NOT NULL DEFAULT 'EMP',  -- fijado por CHECK (C-7)
    empresa_id   INTEGER     NOT NULL,
    cargo        VARCHAR(100),
    es_principal BOOLEAN     NOT NULL DEFAULT FALSE,
    activo       BOOLEAN     NOT NULL DEFAULT TRUE,
    fecha_creacion      TIMESTAMP NOT NULL DEFAULT NOW(),
    fecha_modificacion  TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_usuario_empresa PRIMARY KEY (id),
    CONSTRAINT uq_usuario_empresa UNIQUE (usuario_id, empresa_id),
    CONSTRAINT ck_usuario_empresa_tipo_usuario CHECK (tipo_usuario = 'EMP'),
    -- FK compuesta: solo enlaza usuarios cuyo tipo_usuario sea EMP (C-7)
    CONSTRAINT fk_usuario_empresa_usuario
        FOREIGN KEY (usuario_id, tipo_usuario) REFERENCES usuario(id, tipo_usuario),
    CONSTRAINT fk_usuario_empresa_empresa
        FOREIGN KEY (empresa_id) REFERENCES empresa(id)
);
COMMENT ON TABLE usuario_empresa IS
    'Relación N:M entre usuarios (tipo EMP) y empresas. Subtipo EXCLUSIVO de usuario (C-7): '
    'la FK compuesta (usuario_id, tipo_usuario) con tipo_usuario fijado a EMP impide que un '
    'egresado (EGR) opere como representante de empresa. Lógica de asignación validada con Ágora.';
COMMENT ON COLUMN usuario_empresa.tipo_usuario IS
    'Columna discriminadora fijada a ''EMP'' (CHECK). Parte de la FK compuesta hacia '
    'usuario(id, tipo_usuario) que materializa el arco exclusivo egresado/usuario_empresa.';

CREATE INDEX idx_usuario_empresa_empresa ON usuario_empresa(empresa_id);


-- -------------------------------------------------------------
-- Beneficios
-- -------------------------------------------------------------

-- categoria_beneficio_id / estado_beneficio_id: REFERENCIAS LÓGICAS (C-6).
--   categoria_beneficio_id → TipoParametro CATEGORIA_BENEFICIO
--   estado_beneficio_id    → TipoParametro ESTADO_BENEFICIO
--                            (BORRADOR, PUBLICADO, AGOTADO, VENCIDO, RETIRADO)

CREATE TABLE beneficio (
    id                      SERIAL       NOT NULL,
    empresa_id              INTEGER      NOT NULL,
    categoria_beneficio_id  INTEGER      NOT NULL, -- ref. lógica → parametro (C-6)
    estado_beneficio_id     INTEGER      NOT NULL, -- ref. lógica → parametro (C-6)
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
    -- FK LOCALES reales (a tablas del propio esquema):
    CONSTRAINT fk_beneficio_empresa
        FOREIGN KEY (empresa_id)         REFERENCES empresa(id),
    CONSTRAINT fk_beneficio_usuario_creador
        FOREIGN KEY (usuario_creador_id) REFERENCES usuario(id)
);
COMMENT ON TABLE beneficio IS
    'Beneficio publicado por una empresa aliada para egresados. '
    'categoria_beneficio_id y estado_beneficio_id son referencias lógicas a parámetros '
    'institucionales (C-6): sin FK, validadas en el MID.';
COMMENT ON COLUMN beneficio.categoria_beneficio_id IS
    'Referencia LÓGICA (virtual) a un Parametro del TipoParametro CATEGORIA_BENEFICIO. Sin FK (C-6).';
COMMENT ON COLUMN beneficio.estado_beneficio_id IS
    'Referencia LÓGICA (virtual) a un Parametro del TipoParametro ESTADO_BENEFICIO '
    '(BORRADOR, PUBLICADO, AGOTADO, VENCIDO, RETIRADO). Sin FK (C-6).';
COMMENT ON COLUMN beneficio.usuario_creador_id IS
    'FK LOCAL de trazabilidad/auditoría: qué usuario (tipo EMP o ADM) creó el beneficio. '
    'Es distinto del camino usuario → empresa → beneficio (pertenencia organizacional). '
    'Ambos caminos son intencionales: uno es auditoría de autoría, el otro jerarquía de negocio.';

CREATE INDEX idx_beneficio_estado    ON beneficio(estado_beneficio_id);
CREATE INDEX idx_beneficio_categoria ON beneficio(categoria_beneficio_id);
CREATE INDEX idx_beneficio_empresa   ON beneficio(empresa_id);
CREATE INDEX idx_beneficio_vigencia  ON beneficio(fecha_fin, cupos_disponibles);


-- -------------------------------------------------------------
-- Radicado: SEQUENCE nativa + función (C-5)
-- (Reemplaza la antigua tabla secuencia_radicado)
-- -------------------------------------------------------------

-- Objeto de serialización nativo de PostgreSQL. La unicidad y la
-- concurrencia las garantiza el propio motor (nextval es atómico),
-- sin necesidad de una tabla-contador con SELECT FOR UPDATE.

CREATE SEQUENCE seq_radicado_beneficio
    AS INTEGER
    START WITH 1
    INCREMENT BY 1
    MINVALUE 1
    NO MAXVALUE
    CACHE 1;
COMMENT ON SEQUENCE seq_radicado_beneficio IS
    'Secuencia nativa que alimenta el consecutivo del radicado (C-5). '
    'Reemplaza la tabla secuencia_radicado. nextval() es atómico: garantiza unicidad bajo '
    'concurrencia sin SELECT FOR UPDATE. Reinicio anual = decisión operativa: '
    'ALTER SEQUENCE seq_radicado_beneficio RESTART WITH 1; (cron 1-ene si se desea numeración por año).';

CREATE OR REPLACE FUNCTION fn_siguiente_radicado()
    RETURNS VARCHAR
    LANGUAGE sql
AS $$
    SELECT 'BNF-' || TO_CHAR(CURRENT_DATE, 'YYYY') || '-' ||
           LPAD(nextval('beneficios_egresados.seq_radicado_beneficio')::TEXT, 6, '0');
$$;
COMMENT ON FUNCTION fn_siguiente_radicado() IS
    'Devuelve el siguiente radicado con formato BNF-YYYY-NNNNNN usando la secuencia nativa '
    'seq_radicado_beneficio (C-5). Se usa como DEFAULT de solicitud_beneficio.radicado y puede '
    'invocarse desde el MID con SELECT beneficios_egresados.fn_siguiente_radicado().';


-- -------------------------------------------------------------
-- Solicitudes de beneficio
-- -------------------------------------------------------------

-- El estado vigente es el ÚLTIMO registro en historial_solicitud
-- (ORDER BY fecha_cambio DESC LIMIT 1). NO existe estado_solicitud_id aquí.
-- El radicado se genera con la SEQUENCE nativa vía DEFAULT (C-5).

CREATE TABLE solicitud_beneficio (
    id                    SERIAL      NOT NULL,
    radicado              VARCHAR(20) NOT NULL DEFAULT beneficios_egresados.fn_siguiente_radicado(),
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
    'Solicitud de un egresado sobre un beneficio. El radicado se genera con la secuencia nativa '
    'fn_siguiente_radicado() (C-5). El estado vigente se obtiene del último registro en '
    'historial_solicitud (C-4b): NO se almacena estado aquí.';
COMMENT ON COLUMN solicitud_beneficio.radicado IS
    'Radicado oficial BNF-YYYY-NNNNNN. DEFAULT = fn_siguiente_radicado() (secuencia nativa, C-5). '
    'UNIQUE + CHECK de formato garantizan unicidad e integridad. El consecutivo NNNNNN admite '
    'hasta 999999 por ciclo de la secuencia (reinicio anual operativo si se requiere).';

CREATE INDEX idx_solicitud_beneficio_egresado  ON solicitud_beneficio(egresado_id);
CREATE INDEX idx_solicitud_beneficio_beneficio ON solicitud_beneficio(beneficio_id);
CREATE INDEX idx_solicitud_beneficio_fecha     ON solicitud_beneficio(fecha_solicitud);


-- -------------------------------------------------------------
-- Historial de estado de solicitud
-- (Tabla unificada: el último registro es el estado vigente — C-4b)
-- -------------------------------------------------------------

-- estado_anterior_id / estado_nuevo_id: REFERENCIAS LÓGICAS (C-6) a
--   Parametro del TipoParametro ESTADO_SOLICITUD
--   (PENDIENTE, EN_REVISION, REQUIERE_INFO, APROBADA, RECHAZADA, CANCELADA)

CREATE TABLE historial_solicitud (
    id                      SERIAL    NOT NULL,
    solicitud_beneficio_id  INTEGER   NOT NULL,
    estado_anterior_id      INTEGER,             -- ref. lógica → parametro (NULL en estado inicial)
    estado_nuevo_id         INTEGER   NOT NULL,  -- ref. lógica → parametro
    usuario_id              INTEGER   NOT NULL,
    justificacion           TEXT,
    nombre_archivo_comprobante VARCHAR(300),  -- comprobante OPCIONAL que la empresa adjunta al aprobar
    enlace_comprobante      VARCHAR(100),      -- uid/"Enlace" en gestor_documental_mid (ref. lógica, sin FK)
    fecha_cambio            TIMESTAMP NOT NULL DEFAULT NOW(),
    activo                  BOOLEAN   NOT NULL DEFAULT TRUE,
    fecha_creacion          TIMESTAMP NOT NULL DEFAULT NOW(),
    fecha_modificacion      TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_historial_solicitud PRIMARY KEY (id),
    -- FK LOCALES reales (estado_*_id son referencias lógicas SIN FK — C-6):
    CONSTRAINT fk_historial_solicitud_solicitud
        FOREIGN KEY (solicitud_beneficio_id) REFERENCES solicitud_beneficio(id),
    CONSTRAINT fk_historial_usuario
        FOREIGN KEY (usuario_id) REFERENCES usuario(id)
);
COMMENT ON TABLE historial_solicitud IS
    'Bitácora de transiciones de estado de cada solicitud. El estado vigente es el registro con '
    'mayor fecha_cambio (C-4b). estado_anterior_id/estado_nuevo_id son referencias lógicas (C-6) '
    'a Parametro del TipoParametro ESTADO_SOLICITUD; sin FK, validadas en el MID.';
COMMENT ON COLUMN historial_solicitud.estado_anterior_id IS
    'Referencia LÓGICA (virtual, C-6) al estado ORIGEN de la transición (TipoParametro '
    'ESTADO_SOLICITUD). NULL en el primer registro (creación). Con estado_nuevo_id forma el par '
    'origen→destino; ambos describen el mismo registro de transición, no es redundancia.';
COMMENT ON COLUMN historial_solicitud.estado_nuevo_id IS
    'Referencia LÓGICA (virtual, C-6) al estado DESTINO de la transición; nunca NULL. El estado '
    'vigente se obtiene de este campo en el registro con mayor fecha_cambio (vista '
    'v_solicitud_estado_vigente). usuario_id registra quién ejecutó el cambio (auditoría de acción).';
COMMENT ON COLUMN historial_solicitud.enlace_comprobante IS
    'uid/"Enlace" devuelto por gestor_documental_mid (IdTipoDocumento=167), OPCIONAL. Solo se usa en '
    'la transición a APROBADA: la empresa puede adjuntar un comprobante (p. ej. cupón, certificado) '
    'al aprobar la solicitud. NULL en el resto de transiciones.';

CREATE INDEX idx_historial_solicitud_solicitud ON historial_solicitud(solicitud_beneficio_id);
CREATE INDEX idx_historial_solicitud_fecha     ON historial_solicitud(fecha_cambio);
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
-- Documentos requeridos por beneficio y documentos subidos en solicitud
-- -------------------------------------------------------------

-- documento_requerido_beneficio: qué documentos pide la empresa al publicar
-- el beneficio (p. ej. "Hoja de vida", "Cédula"). El binario NO se guarda aquí:
-- vive en el servicio institucional gestor_documental_mid (Nuxeo); este schema
-- solo referencia el enlace/uid que ese servicio devuelve (ver documento_solicitud).

CREATE TABLE documento_requerido_beneficio (
    id                  SERIAL       NOT NULL,
    beneficio_id        INTEGER      NOT NULL,
    nombre              VARCHAR(200) NOT NULL,
    descripcion         TEXT         NOT NULL,
    activo              BOOLEAN      NOT NULL DEFAULT TRUE,
    fecha_creacion      TIMESTAMP    NOT NULL DEFAULT NOW(),
    fecha_modificacion  TIMESTAMP    NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_documento_requerido_beneficio PRIMARY KEY (id),
    CONSTRAINT fk_documento_requerido_beneficio_beneficio
        FOREIGN KEY (beneficio_id) REFERENCES beneficio(id)
);
COMMENT ON TABLE documento_requerido_beneficio IS
    'Documento que la empresa exige al egresado para postularse a un beneficio '
    '(definido al publicar el beneficio, RF-005). Solo nombre/descripción: el archivo '
    'en sí lo sube el egresado por solicitud (ver documento_solicitud).';

CREATE INDEX idx_documento_requerido_beneficio_beneficio ON documento_requerido_beneficio(beneficio_id);

-- documento_solicitud: el PDF que el egresado subió para cumplir un requisito
-- de una solicitud puntual. enlace_gestor_documental es el uid/"Enlace" que
-- devuelve gestor_documental_mid al subir (referencia LÓGICA a un servicio
-- externo, igual criterio que las referencias a parametro, C-6): no hay FK
-- posible porque Nuxeo/gestor_documental_mid no pertenece a este esquema.

CREATE TABLE documento_solicitud (
    id                          SERIAL       NOT NULL,
    solicitud_beneficio_id      INTEGER      NOT NULL,
    documento_requerido_id      INTEGER      NOT NULL,
    nombre_archivo              VARCHAR(300) NOT NULL,
    enlace_gestor_documental    VARCHAR(100) NOT NULL,
    comentario_empresa          TEXT,
    fecha_comentario            TIMESTAMP,
    activo                      BOOLEAN      NOT NULL DEFAULT TRUE,
    fecha_creacion              TIMESTAMP    NOT NULL DEFAULT NOW(),
    fecha_modificacion          TIMESTAMP    NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_documento_solicitud PRIMARY KEY (id),
    CONSTRAINT fk_documento_solicitud_solicitud
        FOREIGN KEY (solicitud_beneficio_id) REFERENCES solicitud_beneficio(id),
    CONSTRAINT fk_documento_solicitud_requerido
        FOREIGN KEY (documento_requerido_id) REFERENCES documento_requerido_beneficio(id)
);
COMMENT ON TABLE documento_solicitud IS
    'PDF subido por el egresado para cumplir un documento_requerido_beneficio dentro de '
    'una solicitud. El binario vive en gestor_documental_mid (Nuxeo); enlace_gestor_documental '
    'es el uid ("Enlace") que ese servicio devuelve al subir — referencia lógica, sin FK '
    '(el servicio es externo a este esquema).';
COMMENT ON COLUMN documento_solicitud.enlace_gestor_documental IS
    'uid/"Enlace" devuelto por gestor_documental_mid (POST document/upload, '
    'IdTipoDocumento=167). Se usa para consultar (GET document/:uid) o eliminar '
    '(DELETE document/:uid) el archivo en Nuxeo.';
COMMENT ON COLUMN documento_solicitud.comentario_empresa IS
    'Observación de la empresa sobre el documento (p. ej. "no es legible", "falta la '
    'firma"). Campo único: se sobreescribe si la empresa vuelve a comentar; fecha_comentario '
    'registra cuándo se dejó el comentario vigente.';

CREATE INDEX idx_documento_solicitud_solicitud ON documento_solicitud(solicitud_beneficio_id);
CREATE INDEX idx_documento_solicitud_requerido ON documento_solicitud(documento_requerido_id);


-- -------------------------------------------------------------
-- Bitácora PII (registro inmutable, sin FK de borrado)
-- -------------------------------------------------------------

CREATE TABLE bitacora_acceso_pii (
    id                  SERIAL        NOT NULL,
    usuario_id          INTEGER       NOT NULL,
    recurso_tipo        VARCHAR(50)   NOT NULL,
    recurso_id          INTEGER,
    accion              VARCHAR(50)   NOT NULL,
    direccion_ip        VARCHAR(45),
    user_agent          VARCHAR(500),
    detalle             JSONB,
    fecha_evento        TIMESTAMP     NOT NULL DEFAULT NOW(),
    activo              BOOLEAN       NOT NULL DEFAULT TRUE,
    fecha_creacion      TIMESTAMP     NOT NULL DEFAULT NOW(),
    fecha_modificacion  TIMESTAMP     NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_bitacora_acceso_pii PRIMARY KEY (id),
    CONSTRAINT fk_bitacora_acceso_pii_usuario FOREIGN KEY (usuario_id) REFERENCES usuario(id)
);
COMMENT ON TABLE bitacora_acceso_pii IS
    'Bitácora de accesos a datos personales (Ley 1581 Colombia). Retención mínima 6 meses. '
    'El campo activo permite deshabilitar registros sin borrarlos físicamente.';

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
    'Usar en listados para no repetir la lógica ORDER BY / DISTINCT ON.';


-- =============================================================
-- APÉNDICE — APROVISIONAMIENTO EN EL SERVICIO INSTITUCIONAL DE PARÁMETROS
-- =============================================================
-- IMPORTANTE: ESTO **NO** ES PARTE DEL SCHEMA beneficios_egresados.
-- Tras la "virtualización" (C-6), parametro/tipo_parametro NO se referencian
-- por FK desde este esquema. Este bloque queda como GUÍA OPERATIVA de qué
-- TipoParametro/Parametro deben existir en el servicio institucional
-- (parametros_crud de la OATI) para que el MID resuelva las referencias
-- lógicas (*_id) de estado/categoría/sector.
--
-- Nota: TIPO_USUARIO ya NO aparece aquí: pasó a ser discriminador local en
-- la tabla usuario (C-7). Solo se aprovisionan los catálogos que siguen
-- siendo referencias lógicas a parametro.
--
-- Ejecutar SOLO en el servicio/schema "parametro" del clúster, una vez,
-- añadiendo el area_tipo_id real y control de idempotencia (ON CONFLICT).
-- Se deja COMENTADO para evitar ejecución accidental dentro de este DDL.
-- =============================================================

/*
SET search_path TO parametro;

INSERT INTO tipo_parametro (nombre, descripcion, codigo_abreviacion, activo, numero_orden)
VALUES
    ('Estado Empresa',       'Estados del ciclo de vida de una empresa aliada',                       'ESTADO_EMPRESA',      TRUE, 1),
    ('Estado Beneficio',     'Estados de un beneficio publicado',                                     'ESTADO_BENEFICIO',    TRUE, 2),
    ('Estado Solicitud',     'Estados de una solicitud de beneficio',                                 'ESTADO_SOLICITUD',    TRUE, 3),
    ('Categoría Beneficio',  'Categorías de clasificación de beneficios para egresados',              'CATEGORIA_BENEFICIO', TRUE, 4),
    ('Sector Económico',     'Sectores económicos para clasificación de empresas aliadas',            'SECTOR_ECONOMICO',    TRUE, 5),
    ('Parámetro Sistema',    'Parámetros configurables del módulo (límites, paginación, validaciones)','PARAMETRO_SISTEMA',   TRUE, 6);

-- ESTADO_EMPRESA
INSERT INTO parametro (nombre, descripcion, codigo_abreviacion, activo, numero_orden, tipo_parametro_id)
SELECT nombre, descripcion, codigo_abreviacion, activo, numero_orden,
       (SELECT id FROM tipo_parametro WHERE codigo_abreviacion = 'ESTADO_EMPRESA')
FROM (VALUES
    ('Activa',     'Empresa activa en el módulo (llega aprobada desde Ágora)', 'ACTIVA',     TRUE, 1),
    ('Suspendida', 'Empresa suspendida temporalmente del módulo',              'SUSPENDIDA', TRUE, 2)
) AS v(nombre, descripcion, codigo_abreviacion, activo, numero_orden);

-- ESTADO_BENEFICIO
INSERT INTO parametro (nombre, descripcion, codigo_abreviacion, activo, numero_orden, tipo_parametro_id)
SELECT nombre, descripcion, codigo_abreviacion, activo, numero_orden,
       (SELECT id FROM tipo_parametro WHERE codigo_abreviacion = 'ESTADO_BENEFICIO')
FROM (VALUES
    ('Borrador',  'Beneficio en edición, no visible para egresados', 'BORRADOR',  TRUE, 1),
    ('Publicado', 'Beneficio activo y visible en el catálogo',       'PUBLICADO', TRUE, 2),
    ('Agotado',   'Beneficio sin cupos disponibles',                 'AGOTADO',   TRUE, 3),
    ('Vencido',   'Beneficio fuera de su periodo de vigencia',       'VENCIDO',   TRUE, 4),
    ('Retirado',  'Beneficio retirado manualmente',                  'RETIRADO',  TRUE, 5)
) AS v(nombre, descripcion, codigo_abreviacion, activo, numero_orden);

-- ESTADO_SOLICITUD
INSERT INTO parametro (nombre, descripcion, codigo_abreviacion, activo, numero_orden, tipo_parametro_id)
SELECT nombre, descripcion, codigo_abreviacion, activo, numero_orden,
       (SELECT id FROM tipo_parametro WHERE codigo_abreviacion = 'ESTADO_SOLICITUD')
FROM (VALUES
    ('Pendiente',            'Solicitud recibida, sin revisión aún',           'PENDIENTE',     TRUE, 1),
    ('En revisión',          'Solicitud siendo evaluada por la empresa',       'EN_REVISION',   TRUE, 2),
    ('Requiere información', 'Empresa solicita datos adicionales al egresado',  'REQUIERE_INFO', TRUE, 3),
    ('Aprobada',             'Solicitud aprobada por la empresa',              'APROBADA',      TRUE, 4),
    ('Rechazada',            'Solicitud rechazada por la empresa',             'RECHAZADA',     TRUE, 5),
    ('Cancelada',            'Solicitud cancelada por el egresado o el admin', 'CANCELADA',     TRUE, 6)
) AS v(nombre, descripcion, codigo_abreviacion, activo, numero_orden);

-- CATEGORIA_BENEFICIO
INSERT INTO parametro (nombre, descripcion, codigo_abreviacion, activo, numero_orden, tipo_parametro_id)
SELECT nombre, descripcion, codigo_abreviacion, activo, numero_orden,
       (SELECT id FROM tipo_parametro WHERE codigo_abreviacion = 'CATEGORIA_BENEFICIO')
FROM (VALUES
    ('Educación',  'Formación académica y capacitación',                'EDUCACION',  TRUE, 1),
    ('Salud',      'Servicios de salud y bienestar',                    'SALUD',      TRUE, 2),
    ('Recreación', 'Recreativos, culturales y deportivos',              'RECREACION', TRUE, 3),
    ('Empleo',     'Oportunidades laborales y prácticas profesionales', 'EMPLEO',     TRUE, 4),
    ('Descuentos', 'Descuentos en productos y servicios',               'DESCUENTOS', TRUE, 5),
    ('Otro',       'No clasificable en las categorías anteriores',      'OTRO',       TRUE, 6)
) AS v(nombre, descripcion, codigo_abreviacion, activo, numero_orden);

-- SECTOR_ECONOMICO
INSERT INTO parametro (nombre, descripcion, codigo_abreviacion, activo, numero_orden, tipo_parametro_id)
SELECT nombre, descripcion, codigo_abreviacion, activo, numero_orden,
       (SELECT id FROM tipo_parametro WHERE codigo_abreviacion = 'SECTOR_ECONOMICO')
FROM (VALUES
    ('Tecnología e Innovación', 'TIC, software y servicios digitales',          'TEC',  TRUE, 1),
    ('Salud y Farmacéutico',    'Clínicas, hospitales, laboratorios',           'SAL',  TRUE, 2),
    ('Educación',               'Instituciones educativas y e-learning',        'EDU',  TRUE, 3),
    ('Industria y Manufactura', 'Fabricación y producción industrial',          'IND',  TRUE, 4),
    ('Comercio y Retail',       'Comercio mayorista y minorista',               'COM',  TRUE, 5),
    ('Servicios Financieros',   'Bancos, aseguradoras, fintech',                'FIN',  TRUE, 6),
    ('Construcción',            'Constructoras, inmobiliarias, infraestructura', 'CON',  TRUE, 7),
    ('Alimentos y Bebidas',     'Producción y venta de alimentos y bebidas',    'ALI',  TRUE, 8),
    ('Consultoría y Servicios', 'Consultoría, outsourcing y servicios',         'CON2', TRUE, 9),
    ('Otro',                    'Sectores no clasificados',                     'OTR',  TRUE, 10)
) AS v(nombre, descripcion, codigo_abreviacion, activo, numero_orden);

-- PARAMETRO_SISTEMA
INSERT INTO parametro (nombre, descripcion, codigo_abreviacion, activo, numero_orden, tipo_parametro_id)
SELECT nombre, descripcion, codigo_abreviacion, activo, numero_orden,
       (SELECT id FROM tipo_parametro WHERE codigo_abreviacion = 'PARAMETRO_SISTEMA')
FROM (VALUES
    ('Límite solicitudes activas egresado',
     'Máximo de solicitudes activas por egresado (RN-010). Valor: 5',
     'LIMITE_SOLICITUDES_ACTIVAS_EGRESADO', TRUE, 1),
    ('Paginación catálogo por defecto',
     'Tamaño de página por defecto del catálogo. Valor: 20',
     'PAGINACION_CATALOGO_DEFAULT',         TRUE, 2),
    ('Mínimo caracteres justificación rechazo',
     'Longitud mínima de la justificación al rechazar (RN-003). Valor: 20',
     'JUSTIFICACION_RECHAZO_MIN_CARACTERES', TRUE, 3)
) AS v(nombre, descripcion, codigo_abreviacion, activo, numero_orden);

SET search_path TO beneficios_egresados;
*/

-- =============================================================
-- FIN DDL v4
-- =============================================================
