-- Migración 2026-07-05: documentos requeridos por beneficio y comprobante de aprobación
-- Delta entre el schema.sql previo (1763eae) y el nuevo (a7db242). Idempotente.

SET search_path = beneficios_egresados;

-- 1. Comprobante opcional al aprobar (historial_solicitud)
ALTER TABLE historial_solicitud
    ADD COLUMN IF NOT EXISTS nombre_archivo_comprobante VARCHAR(300),
    ADD COLUMN IF NOT EXISTS enlace_comprobante         VARCHAR(100);

COMMENT ON COLUMN historial_solicitud.enlace_comprobante IS
    'uid/"Enlace" devuelto por gestor_documental_mid (IdTipoDocumento=167), OPCIONAL. Solo se usa en '
    'la transición a APROBADA: la empresa puede adjuntar un comprobante (p. ej. cupón, certificado) '
    'al aprobar la solicitud. NULL en el resto de transiciones.';

-- 2. Documentos que la empresa exige al publicar un beneficio (RF-005)
CREATE TABLE IF NOT EXISTS documento_requerido_beneficio (
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

CREATE INDEX IF NOT EXISTS idx_documento_requerido_beneficio_beneficio
    ON documento_requerido_beneficio(beneficio_id);

-- 3. PDFs subidos por el egresado por solicitud
CREATE TABLE IF NOT EXISTS documento_solicitud (
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

CREATE INDEX IF NOT EXISTS idx_documento_solicitud_solicitud ON documento_solicitud(solicitud_beneficio_id);
CREATE INDEX IF NOT EXISTS idx_documento_solicitud_requerido ON documento_solicitud(documento_requerido_id);
