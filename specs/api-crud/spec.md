# Spec — API CRUD (`egresados_crud`)

> **Última actualización:** 2026-07-08 · **Estado:** implementado y en uso por el
> MID. Deriva del `BACKEND_SPEC.md` original, actualizado al contrato real.

## Objetivo

Exponer acceso CRUD puro (sin lógica de negocio) al schema `beneficios_egresados`
de PostgreSQL, con el contrato de consulta idiomático del SGA, para consumo
exclusivo del MID.

## Alcance

**In scope:** endpoints CRUD por entidad, DSL de consulta de los GetAll,
endpoints atómicos de cupo, endpoints de historial/estado vigente.
**Out of scope:** reglas de negocio (MID), validación de auth (hoy confía en el
gateway/MID; recibe el Bearer de forma uniforme para un futuro filtro),
resolución de catálogos (ids planos).

## Repos involucrados

- `egresados_crud` (este) — Go + Beego + Beego ORM + PostgreSQL.
- Consumidor único: `sga_mid_beneficios_egresados`. El frontend nunca lo llama.

## Requisitos

1. **CRUD por entidad** bajo `/v1`: `usuario`, `egresado`, `empresa`, `usuario-empresa`, `beneficio`, `solicitud-beneficio`, `historial-solicitud`, `mensaje-solicitud`, `documento-solicitud`, `bitacora-acceso-pii`.
2. **DSL de consulta en todos los GetAll** (contrato de `terceros_crud`, la variante más completa del SGA), centralizado en `models/getall_query.go` + `controllers/getall_params.go` (no copy-paste por entidad):
   - `query=` con dot-notation (`.`→`__`), `__in` con `|`, `__icontainsall`, `isnull`, y operadores nativos del ORM (`__gte`, `__icontains`, …).
   - `fields`, `sortby`, `order`, `limit` (default 10; `limit=0` = todos), `offset`.
   - Los GetAll **no** fuerzan `Activo:true`: el caller lo pasa en `query` (convención SGA).
3. **Lista vacía responde `[{}]`** (idioma SGA). OJO: en RunMode=dev Beego lo pretty-printa; los consumidores deben normalizar por JSON compactado, no por comparación literal.
4. **Cupos atómicos (RN-002b/c):**
   - `POST /v1/beneficio/:id/cupo/descontar` — `UPDATE … WHERE cupos_disponibles > 0`.
   - `POST /v1/beneficio/:id/cupo/devolver` — `UPDATE … WHERE cupos_disponibles < cupos_total`.
   - Sin race conditions; devuelven error si no hay cupo que mover.
5. **Historial / estado vigente (C-4b):**
   - `GET /v1/historial-solicitud/solicitud/:solicitud_id` — bitácora, más reciente primero.
   - `GET /v1/historial-solicitud/solicitud/:solicitud_id/vigente` — estado actual.
6. **Ids de catálogo planos:** campos que referencian parámetros institucionales son `int` (`*int` si nullable, para serializar NULL). Relaciones locales sí son objetos ORM (`{id}`).
7. **PUT reemplaza la fila completa** (`o.Update` sin lista de columnas): el caller debe leer el objeto y mandar el row entero.

## Contrato de integración (con el MID)

- Envelope: respuestas Beego directas (objeto o array JSON); errores con status HTTP correcto — los POST/PUT fallidos deben devolver status ≠ 2xx (bug corregido 2026-07-02: antes se tragaban en silencio).
- El MID normaliza `[{}]` → `[]` (`normalizarListaVacia` con `json.Compact`) y valida status en GET/POST/PUT.
- Variables de entorno: `EGRESADOS_CRUD_PG{USER,PASS,HOST,PORT,DB,SCHEMA}`, `EGRESADOS_CRUD_RUN_MODE`, `EGRESADOS_CRUD_HTTP_PORT` (dev: 8080), `PARAMETER_STORE`.

## Criterios de aceptación

1. `GET /v1/beneficio?query=EstadoBeneficioId:21,Activo:true&limit=0` filtra de verdad en SQL (no en memoria).
2. Descontar cupo con `cupos_disponibles=0` responde error y no modifica la fila; 2 descuentos concurrentes sobre cupo 1 dejan exactamente 0.
3. `GET …/vigente` devuelve el último estado del historial para cualquier solicitud con al menos un registro.
4. POST con FK inválida responde status de error (no 200 con id 0).

## Casos borde

- Paginación: `limit` default 10 puede sorprender — los listados completos del MID pasan `limit=0` explícito.
- La restricción de privacidad/columnas NO va en el CRUD (convención SGA): la minimización RNF-002b es responsabilidad del MID.
