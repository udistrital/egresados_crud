-- =============================================================
-- seed_estados_beneficio.sql — Beneficios de prueba en TODOS los estados
-- (para verificar la UI: agotado / vencido / retirado / borrador).
-- SOLO DESARROLLO — no ejecutar en producción.
-- Ejecutar sobre la BD viva ya sembrada (empresa id 2 = AMB GROUP SAS,
-- usuario creador id 2). Ids INSTITUCIONALES de parámetro (2026-07-07),
-- los mismos del servicio real y de la semilla local del MID:
--   ESTADO_BENEFICIO: BORRADOR=7201, PUBLICADO=7202, AGOTADO=7203,
--                     VENCIDO=7204, RETIRADO=7205
--   CATEGORIA_BENEFICIO: EDUCACION=7212, SALUD=7213, RECREACION=7214,
--                        EMPLEO=7215, DESCUENTOS=7216, OTRO=7217
-- =============================================================

INSERT INTO beneficios_egresados.beneficio
  (empresa_id, categoria_beneficio_id, estado_beneficio_id, titulo, descripcion, condiciones,
   fecha_inicio, fecha_fin, cupos_total, cupos_disponibles, fecha_publicacion, usuario_creador_id)
VALUES
  -- PUBLICADOS vigentes pero AGOTADOS (cupos_disponibles = 0) → "Sin cupos" en la UI
  (2, 7216, 7202, 'Bono 50% en certificación de contratación pública',
   'Cofinanciación del 50% del valor de la certificación en contratación estatal con entidad acreditada. Incluye material de estudio y un intento de examen.',
   E'Ser egresado UD con carné vigente\nNo haber recibido este bono en periodos anteriores\nInscribirse antes del cierre de la convocatoria',
   '2026-06-20', '2026-09-15', 15, 0, NOW(), 2),
  (2, 7215, 7202, 'Mentoría ejecutiva 1:1 con la gerencia',
   'Programa de 4 sesiones individuales de mentoría profesional con el equipo directivo de AMB GROUP, orientado a egresados en transición de carrera.',
   E'Contar con mínimo 1 año de experiencia profesional\nDiligenciar el formulario de objetivos de carrera\nCompromiso de asistencia a las 4 sesiones',
   '2026-06-25', '2026-08-20', 5, 0, NOW(), 2),
  -- VENCIDO (fecha_fin pasada) → NO debe aparecer en el catálogo del egresado
  (2, 7213, 7204, 'Jornada de vacunación empresarial 2026-I',
   'Jornada de vacunación gratuita (influenza y tétano) en las instalaciones de la empresa para egresados UD y sus familias.',
   'Presentar documento de identidad y carné de egresado.',
   '2026-03-01', '2026-06-15', 40, 12, '2026-03-01', 2),
  -- RETIRADO por la empresa (vigente en fechas) → NO debe aparecer
  (2, 7214, 7205, 'Pases dobles a feria empresarial',
   'Entrada doble a la feria de proveedores del sector construcción en Corferias.',
   'Registro previo con correo institucional.',
   '2026-06-01', '2026-10-30', 8, 8, '2026-06-01', 2),
  -- BORRADOR (sin publicar) → NO debe aparecer
  (2, 7212, 7201, 'Beca completa bootcamp de análisis de datos',
   'Beca del 100% para bootcamp intensivo de análisis de datos (120 horas, modalidad mixta).',
   'Convocatoria en preparación.',
   '2026-08-01', '2026-11-30', 10, 10, NULL, 2);
