# sga_crud_beneficios_egresados

API CRUD del submódulo **Beneficios para Egresados** del Sistema de Gestión Académica
(SGA) de la Universidad Distrital Francisco José de Caldas. Expone las entidades
propias del módulo sobre PostgreSQL (schema `beneficios_egresados`) siguiendo el
contrato estándar de los `*_crud` institucionales.

Este servicio no contiene lógica de negocio: la orquestación y las reglas viven en
[`egresados_service`](https://github.com/udistrital/egresados_service).
El micro-frontend es
[`egresados_cliente`](https://github.com/udistrital/egresados_cliente).

## Especificaciones técnicas

- **Go** 1.22 · **Beego** v2.2 (`server/web` + `client/orm`)
- **PostgreSQL** — schema propio `beneficios_egresados` (DDL en [`db/schema.sql`](db/schema.sql))
- Borrado lógico con `activo` y auditoría con `fecha_creacion`/`fecha_modificacion` en todas las tablas (excepto `bitacora_acceso_pii`, que es log inmutable)

### Decisiones de diseño relevantes

- **Sin catálogos locales (C-1):** `estado_*`, `categoria_beneficio`,
  `sector_economico` y `parametro_sistema` viven en el servicio institucional de
  parámetros de la OATI (creados el 2026-07-07). Los campos que los referencian son
  ids planos (`*int` los nullable). `tipo_usuario` es un discriminador local (C-7).
  El apéndice al final de `db/schema.sql` documenta qué se aprovisionó (comentado,
  ya ejecutado — no volver a correr).
- **Historial como única fuente de estado (C-4b):** `solicitud_beneficio` no tiene
  campo de estado; el estado vigente es el último registro de `historial_solicitud`
  (endpoint `/vigente`). Los cambios de estado son INSERT en el historial.
- **Radicados (C-5):** los genera la base de datos al insertar la solicitud
  (`fn_siguiente_radicado()` sobre secuencia nativa, DEFAULT de la columna).
  Formato: `BNF-YYYY-NNNNNN`.

## Variables de entorno

| Variable | Default | Descripción |
|---|---|---|
| `BENEFICIOS_EGRESADOS_CRUD_DB_USER` | `postgres` | Usuario de PostgreSQL |
| `BENEFICIOS_EGRESADOS_CRUD_DB_PASSWORD` | _(vacío)_ | Contraseña de PostgreSQL |
| `BENEFICIOS_EGRESADOS_CRUD_DB_HOST` | `127.0.0.1` | Host de PostgreSQL |
| `BENEFICIOS_EGRESADOS_CRUD_DB_PORT` | `5432` | Puerto de PostgreSQL |
| `BENEFICIOS_EGRESADOS_CRUD_DB_NAME` | `beneficios_egresados` | Base de datos |
| `BENEFICIOS_EGRESADOS_CRUD_DB_SCHEMA` | `beneficios_egresados` | `search_path` |
| `BENEFICIOS_EGRESADOS_CRUD_PORT` | `8080` | Puerto HTTP del servicio |
| `BENEFICIOS_EGRESADOS_CRUD_RUNMODE` | `dev` | Modo de ejecución de Beego |

## Ejecución

```bash
# crear la base de datos y el schema
psql -U postgres -f db/schema.sql

# levantar el servicio
export BENEFICIOS_EGRESADOS_CRUD_DB_PASSWORD=...
go run .
```

## Endpoints

Cada entidad expone `GET /` (listado), `POST /`, `GET /:id`, `PUT /:id` y
`DELETE /:id` (borrado lógico) bajo `/v1`:

`usuario` · `egresado` · `empresa` · `usuario_empresa` · `beneficio` ·
`solicitud_beneficio` · `historial_solicitud` · `mensaje_solicitud` ·
`documento_requerido_beneficio` · `documento_solicitud` ·
`bitacora_acceso_pii` (solo GET/POST, log inmutable)

Rutas especiales:

```
POST /v1/beneficio/:id/cupo/descontar                  → descuento atómico de cupo (RN-002b)
POST /v1/beneficio/:id/cupo/devolver                   → devolución atómica de cupo (RN-002c)
GET  /v1/historial_solicitud/solicitud/:id             → bitácora de la solicitud (desc)
GET  /v1/historial_solicitud/solicitud/:id/vigente     → estado vigente (C-4b)
```

### Contrato de listado (GET de colección)

Mismo contrato que los `*_crud` institucionales (variante de `terceros_crud`):

| Parámetro | Ejemplo | Notas |
|---|---|---|
| `query` | `query=Egresado.Id:1,Activo:true` | `k:v` separados por coma; dot-notation para relaciones; acepta sufijos `__in` (valores con `\|`), `__icontainsall`, `isnull` y los operadores nativos del ORM (`__gte`, `__icontains`, ...) |
| `fields` | `fields=Id,Titulo` | Nombres de campo Go; recorta columnas |
| `sortby` / `order` | `sortby=FechaCreacion&order=desc` | `order` único o uno por campo |
| `limit` | `limit=0` | Default `10`; `0` = sin límite |
| `offset` | `offset=20` | Default `0` |

La lista vacía responde `200` con `[{}]` (idioma estándar del SGA). Los GetAll **no**
filtran `Activo` implícitamente: el consumidor debe pasarlo en `query`.

## Documentación (SDD)

- `specs/base-datos/` — spec del schema (decisiones C-1…C-7, constraints).
- `specs/api-crud/` — spec y tareas de esta API (contrato de listado, rutas especiales).
- `docs/referencia-base-datos-defensa.md` — justificación tabla por tabla del modelo.
- Las especificaciones **transversales** (visión general, autenticación, parámetros)
  viven en `specs/system/` del repo [`egresados_service`](https://github.com/udistrital/egresados_service).

## Contexto

Desarrollado en el marco de la pasantía de Ingeniería de Sistemas (2026) para la
Oficina Asesora de Sistemas (OAS) / OATI. Lineamientos: APIs separadas CRUD/MID,
plantillas `udistrital/plantilla_api_crud`, autenticación WSO2 validada en el MID.
