# Spec — Base de datos (schema `beneficios_egresados`, v4)

> **Última actualización:** 2026-07-08 · **DDL vigente:** `db/schema.sql`
> (este repo). La defensa detallada tabla-por-tabla (formato sustentación)
> está en `docs/referencia-base-datos-defensa.md`.

## Objetivo

Modelar la persistencia local del módulo con el mínimo de datos propios:
identidad aprovisionada (JIT), beneficios, solicitudes y su trazabilidad.
Todo lo institucional (catálogos, datos de egresado/empresa) se referencia,
no se duplica (correcciones C-1, C-2, C-6 del revisor).

## Alcance

**In scope:** 11 tablas locales + 1 vista + 1 secuencia nativa + 1 función.
**Out of scope:** catálogos (servicio de parámetros), datos académicos del
egresado y datos de empresa (on-demand de SGA/Ágora), user store (WSO2).

## Repos involucrados

- `egresados_crud` — dueño del schema (`db/schema.sql`) y de los modelos ORM.
- `sga_mid_beneficios_egresados` — valida pertenencia de los ids de parámetro (C-6) y aplica reglas.

## Inventario

| Objeto | Propósito |
|---|---|
| `usuario` | Identidad local JIT. `documento` **NULLABLE** (empresas self-signup sin cédula); UNIQUE `(sistema_origen, id_externo)` con `id_externo` = `sub` WSO2 |
| `egresado` | Perfil 1:1 con usuario EGR; `codigo_institucional` UNIQUE; `fecha_grado` local (sin fuente institucional); programa/facultad NO se almacenan (C-2a) |
| `empresa` | Espejo mínimo de Ágora: `agora_id_externo` + `nit` UNIQUE + estado de ciclo de vida local (`estado_empresa_id`, id de parámetro) |
| `usuario_empresa` | N:M usuario↔empresa (multiempresa real); UNIQUE `(usuario_id, empresa_id)`; `es_principal` |
| `beneficio` | Publicación de la empresa; `cupos_total`/`cupos_disponibles` con CHECKs; categoría/estado = ids de parámetro planos |
| `solicitud_beneficio` | Solicitud del egresado; `radicado` UNIQUE con CHECK `^BNF-\d{4}-\d{6}$`; `datos_complementarios` JSONB; **sin campo de estado** (C-4b) |
| `historial_solicitud` | Única fuente del estado (C-4b): INSERT por transición; `estado_anterior_id`/`estado_nuevo_id` = ids de parámetro (`*int` nullable) |
| `mensaje_solicitud` | Hilo empresa↔egresado |
| `documento_requerido_beneficio` | Documentos que la empresa exige por beneficio (se definen al publicar) |
| `documento_solicitud` | Vínculo solicitud↔documento del gestor documental (`enlace_gestor_documental` = uid Nuxeo) |
| `bitacora_acceso_pii` | Log inmutable de accesos a datos personales (Ley 1581 / RNF-002a); sin borrado lógico |
| `v_solicitud_estado_vigente` (vista) + `idx_historial_solicitud_vigente` | Estado vigente masivo sin N+1 |
| `fn_siguiente_radicado()` + secuencia nativa | Genera el radicado al INSERT (C-5: reemplazó a la tabla-contador `secuencia_radicado`) |

## Requisitos (decisiones de diseño vigentes)

1. **C-1/C-6 — catálogos virtualizados:** las columnas de catálogo son ids planos de `parametro` institucional, sin FK local; la validación de pertenencia es del MID.
2. **C-2 — no duplicar:** `egresado` y `empresa` guardan solo el mínimo + id externo; el resto se consulta on-demand.
3. **C-3 — `usuario` se conserva siempre** (identidad JIT; validado por el revisor).
4. **C-4b — historial como única fuente de estado:** nunca reintroducir un campo de estado en `solicitud_beneficio`; los cambios de estado son INSERT en `historial_solicitud`.
5. **C-5 — radicado por función de BD:** `fn_siguiente_radicado()` en el INSERT; el MID no envía radicado.
6. **C-7 — subtipos disjuntos:** un usuario es egresado XOR empresa (DDL puro).
7. **Convenciones OATI:** borrado lógico (`activo`), auditoría (`fecha_creacion`/`fecha_modificacion`), tablas/columnas comentadas, sin DELETE físico.
8. **Identidad sin documento:** `usuario.documento` NULLABLE + `uq_usuario_id_externo` (migración aplicada a la BD viva el 2026-07-02).

## Criterios de aceptación

1. `db/schema.sql` aplica limpio sobre una BD vacía y deja la semilla local alineada con los ids institucionales (7199+).
2. Insertar una solicitud sin radicado genera `BNF-YYYY-NNNNNN` correlativo (verificado: `BNF-2026-000003`).
3. Dos usuarios con documento NULL coexisten (Postgres permite múltiples NULL en UNIQUE); dos usuarios con el mismo `(sistema_origen, id_externo)` no.
4. `cupos_disponibles` nunca queda < 0 ni > `cupos_total` (CHECKs + UPDATE atómico del CRUD).

## Casos borde

- Fila de empresa soft-deleted (`activo=false`) sigue ocupando el NIT (`uq_nit_empresa` sin condición): el JIT busca sin filtrar `Activo:true` y **reactiva** en vez de chocar con la restricción (fix 2026-07-08).
- `fecha_grado` NULL es válido (no hay fuente institucional); el riesgo Beego zero-time→año-1 está descartado (inserta NULL).
- FK cross-schema hacia `parametro`: eliminadas (C-6); si el servicio vive en otra BD en producción, no hay nada que romper.
