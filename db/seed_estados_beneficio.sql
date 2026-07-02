-- =============================================================
-- seed_estados_beneficio.sql — Beneficios de prueba en TODOS los estados
-- (para verificar la UI: agotado / vencido / retirado / borrador).
-- Ejecutar sobre la BD viva ya sembrada (empresa id 2 = AMB GROUP SAS,
-- usuario creador id 2). Ids de estado del fallback local del MID:
-- BORRADOR=20, PUBLICADO=21, AGOTADO=22, VENCIDO=23, RETIRADO=24.
-- =============================================================

INSERT INTO beneficios_egresados.beneficio
  (empresa_id, categoria_beneficio_id, estado_beneficio_id, titulo, descripcion, condiciones,
   fecha_inicio, fecha_fin, cupos_total, cupos_disponibles, fecha_publicacion, usuario_creador_id)
VALUES
  -- PUBLICADOS vigentes pero AGOTADOS (cupos_disponibles = 0) → "Sin cupos" en la UI
  (2, 44, 21, 'Bono 50% en certificación de contratación pública',
   'Cofinanciación del 50% del valor de la certificación en contratación estatal con entidad acreditada. Incluye material de estudio y un intento de examen.',
   E'Ser egresado UD con carné vigente\nNo haber recibido este bono en periodos anteriores\nInscribirse antes del cierre de la convocatoria',
   '2026-06-20', '2026-09-15', 15, 0, NOW(), 2),
  (2, 43, 21, 'Mentoría ejecutiva 1:1 con la gerencia',
   'Programa de 4 sesiones individuales de mentoría profesional con el equipo directivo de AMB GROUP, orientado a egresados en transición de carrera.',
   E'Contar con mínimo 1 año de experiencia profesional\nDiligenciar el formulario de objetivos de carrera\nCompromiso de asistencia a las 4 sesiones',
   '2026-06-25', '2026-08-20', 5, 0, NOW(), 2),
  -- VENCIDO (fecha_fin pasada) → NO debe aparecer en el catálogo del egresado
  (2, 41, 23, 'Jornada de vacunación empresarial 2026-I',
   'Jornada de vacunación gratuita (influenza y tétano) en las instalaciones de la empresa para egresados UD y sus familias.',
   'Presentar documento de identidad y carné de egresado.',
   '2026-03-01', '2026-06-15', 40, 12, '2026-03-01', 2),
  -- RETIRADO por la empresa (vigente en fechas) → NO debe aparecer
  (2, 42, 24, 'Pases dobles a feria empresarial',
   'Entrada doble a la feria de proveedores del sector construcción en Corferias.',
   'Registro previo con correo institucional.',
   '2026-06-01', '2026-10-30', 8, 8, '2026-06-01', 2),
  -- BORRADOR (sin publicar) → NO debe aparecer
  (2, 40, 20, 'Beca completa bootcamp de análisis de datos',
   'Beca del 100% para bootcamp intensivo de análisis de datos (120 horas, modalidad mixta).',
   'Convocatoria en preparación.',
   '2026-08-01', '2026-11-30', 10, 10, NULL, 2);
