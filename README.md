# egresados_crud
:heavy_check_mark: Check: Repositorio API CRUD del submódulo **Beneficios para Egresados** del Sistema de Gestión Académica (SGA).

API CRUD del submódulo **Beneficios para Egresados** del Sistema de Gestión Académica
(SGA) de la Universidad Distrital Francisco José de Caldas. Expone las entidades
propias del módulo sobre PostgreSQL (schema `beneficios_egresados`) siguiendo el
contrato estándar de los `*_crud` institucionales.

Este servicio no contiene lógica de negocio: la orquestación y las reglas viven en
[`egresados_service`](https://github.com/udistrital/egresados_service).
El micro-frontend es
[`egresados_cliente`](https://github.com/udistrital/egresados_cliente).

## Especificaciones Técnicas

### Tecnologías Implementadas y Versiones
* [Golang 1.22](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/golang.md) (imagen de CI en `golang:1.25`, compatible hacia atrás)
* [BeeGo v2.2](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/beego.md) (`server/web` + `client/orm`)
* [PostgreSQL](https://www.postgresql.org/) — schema propio `beneficios_egresados` (DDL en [`db/schema.sql`](db/schema.sql))
* [Docker](https://docs.docker.com/engine/install/ubuntu/)

Borrado lógico con `activo` y auditoría con `fecha_creacion`/`fecha_modificacion` en
todas las tablas (excepto `bitacora_acceso_pii`, que es log inmutable).

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

### Variables de Entorno

Definidas en [`conf/app.conf`](conf/app.conf) vía `${VAR||default}` (expansión nativa
de Beego). Si una variable no está seteada, `main.go` cae a los defaults de
desarrollo local.

```shell
# Parámetros de la API
EGRESADOS_CRUD_HTTPPORT=8080
EGRESADOS_CRUD_RUNMODE=dev

# Database
EGRESADOS_CRUD_DB_USER=postgres
EGRESADOS_CRUD_DB_PASS=postgres
EGRESADOS_CRUD_DB_URL=localhost
EGRESADOS_CRUD_DB_NAME=beneficios_egresados
EGRESADOS_CRUD_DB_SCHEMA=beneficios_egresados
EGRESADOS_CRUD_DB_PORT=5432

# Institucional
PARAMETER_STORE=
```

### Ejecución del Proyecto
```shell
# 1. Clonar el repositorio
git clone -b develop https://github.com/udistrital/egresados_crud.git

# 2. Moverse a la carpeta del repositorio
cd egresados_crud

# 3. Moverse a la rama **develop**
git pull origin develop && git checkout develop

# 4. Crear la base de datos y el schema
psql -U postgres -f db/schema.sql

# 5. Configurar las variables de entorno
export EGRESADOS_CRUD_RUNMODE=dev
export EGRESADOS_CRUD_DB_PASS=postgres

# 6. Ejecutar el proyecto
bee run
```

### Ejecución Dockerfile
```shell
# El Dockerfile está implementado para el despliegue mediante
# el sistema de integración continua (CI).

# 1. Construir la imagen
docker build -t egresados_crud .

# 2. Ejecutar el contenedor
docker run --name egresados_crud \
  -e EGRESADOS_CRUD_RUNMODE=dev \
  -e EGRESADOS_CRUD_DB_USER=postgres \
  -e EGRESADOS_CRUD_DB_PASS=postgres \
  -e EGRESADOS_CRUD_DB_URL=host.docker.internal \
  -e EGRESADOS_CRUD_DB_PORT=5432 \
  -e EGRESADOS_CRUD_DB_NAME=beneficios_egresados \
  -e EGRESADOS_CRUD_DB_SCHEMA=beneficios_egresados \
  -p 8080:8080 \
  egresados_crud

# 3. Comprobar que el contenedor esté en ejecución
docker ps
```

### Ejecución docker-compose
```shell
# No implementado actualmente.
```

## Estado CI

| Develop | Master | Sonar |
| -- | -- | -- |
| [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/egresados_crud/status.svg?ref=refs/heads/develop)](https://hubci.portaloas.udistrital.edu.co/udistrital/egresados_crud) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/egresados_crud/status.svg?ref=refs/heads/master)](https://hubci.portaloas.udistrital.edu.co/udistrital/egresados_crud) | [![Quality Gate Status](https://sonarqube.portaloas.udistrital.edu.co/api/project_badges/measure?project=egresados_crud&metric=alert_status)](https://sonar.portaloas.udistrital.edu.co/dashboard?id=egresados_crud) |

## Modelo de datos

- DDL: [`db/schema.sql`](db/schema.sql)
- Justificación tabla por tabla: [`docs/referencia-base-datos-defensa.md`](docs/referencia-base-datos-defensa.md)
- Spec del schema (decisiones C-1…C-7, constraints): [`specs/base-datos/spec.md`](specs/base-datos/spec.md)

> No hay un diagrama `.svg` renderizado todavía (a diferencia de otros `*_crud`
> institucionales que sí lo traen bajo `database/`). Si se necesita, se puede generar
> a partir de `db/schema.sql`.

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

Documentación interactiva (Swagger UI, generada con `bee generate docs`):
[`swagger/`](swagger/), servida en `/swagger/` cuando `EnableDocs = true`.

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

## Licencia

This file is part of egresados_crud.

egresados_crud is free software: you can redistribute it and/or modify it under the
terms of the GNU General Public License as published by the Free Software
Foundation, either version 3 of the License, or (at your option) any later version.

egresados_crud is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with
egresados_crud. If not, see https://www.gnu.org/licenses/.
