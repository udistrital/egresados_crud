# Tasks — API CRUD

> Estado al 2026-07-08.

## Completadas

1. [x] Scaffold Beego + modelos de las 9 tablas (schema v4).
2. [x] DSL de consulta en todos los GetAll (contrato terceros_crud, centralizado). — 2026-06-10
3. [x] Endpoints atómicos de cupo (descontar/devolver, RN-002b/c). — 2026-07-01
4. [x] Endpoints de historial + estado vigente (C-4b). — 2026-06-09
5. [x] Radicado generado por la BD (`fn_siguiente_radicado`, C-5); ruta de secuencia eliminada.
6. [x] `usuario.documento` nullable + `uq_usuario_id_externo` (identidad sin cédula). — 2026-07-01/02
7. [x] Migración de la BD dev a los ids institucionales de parámetros (7199+). — 2026-07-07

## Pendientes

1. [ ] Filtro de seguridad JWT propio (hoy confía en gateway/MID; el Bearer ya llega de forma uniforme). Baja prioridad — decisión de despliegue.
2. [ ] Pruebas automatizadas de `getall_query` (DSL de filtros) con BD de prueba.
3. [ ] Aplicar `db/schema.sql` completo a una BD limpia (la viva se migró con ALTER a mano) para validar reproducibilidad.
