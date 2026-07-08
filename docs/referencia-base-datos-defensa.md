# Referencia — Defensa de la base de datos (material de sustentación, schema v4)

> **Estado del documento (2026-07-08):** antes archivo `contexto_db.md`. Vigente
> para la sustentación (el schema v4 es el actual). La spec resumida es
> `specs/base-datos/spec.md`; el DDL es `db/schema.sql` (este repo).

> **Cómo usar este documento (modo sustentación):** está pensado para Ctrl+F. Si te
> preguntan por una tabla, busca su nombre (p. ej. `## TABLA: solicitud_beneficio`) y
> encontrarás, en este orden:
> 1. **Para qué sirve** (funcionalidad).
> 2. **Por qué es necesaria** (la defensa: por qué existe y no se eliminó/delegó).
> 3. **Campos** (cada uno con su propósito y justificación).
> 4. **Constraints e índices** (qué garantizan).
> 5. **Preguntas probables del revisor → respuesta lista.**
>
> Esquema: `beneficios_egresados` (PostgreSQL). Archivo DDL: `db/schema.sql`.
> OJO: este documento describe el corte v4 del 2026-06-18; el DDL actual añade
> además `documento_solicitud` + comprobante de aprobación (2026-07-05), `usuario.documento`
> nullable + `uq_usuario_id_externo` (2026-07-02) y la semilla con ids institucionales.
> Submódulo del SGA — Universidad Distrital Francisco José de Caldas.
> 9 tablas locales + 1 vista + 1 secuencia nativa + 1 función. Los catálogos
> (estado/categoría/sector) son **referencias lógicas** al servicio institucional de
> parámetros (sin FK declarada).

---

# PARTE A — CONTEXTO GENERAL

## A.1 Historia del esquema y qué pidió cada revisión

El esquema pasó por varias revisiones. Entender la evolución es clave porque **casi todas las
preguntas son "¿esta tabla/relación es necesaria y por qué se modeló así?"**.

- **v1** (original): 7 tablas de catálogo locales (`tipo_usuario`, `estado_empresa`,
  `estado_beneficio`, `estado_solicitud`, `categoria_beneficio`, `sector_economico`,
  `parametro_sistema`) y el estado de la solicitud guardado en un campo.
- **v2**: correcciones de fondo del primer revisor (C-1 a C-4): eliminó los catálogos
  locales, unificó el estado en el historial, conservó `usuario` y `secuencia_radicado`.
- **v3 / v3.1**: checklist de forma OATI (`VARCHAR` con longitud, auditoría en todas las
  tablas) y aclaraciones de diseño del DBA.
- **v4 (actual)**: correcciones de la **sustentación con el profesor** (C-5, C-6, C-7).
  Son **cambios estructurales**, no solo comentarios. Es el corazón de esta defensa:
  - **C-5** — el radicado pasa de tabla-contador a **secuencia nativa de PostgreSQL**.
  - **C-6** — se **"virtualiza" `parametro`**: se quitan todas las FK hacia él para que
    deje de ser un hub de relaciones; las columnas de catálogo quedan como referencias
    lógicas validadas en el MID.
  - **C-7** — se hace **excluyente** la relación usuario↔egresado (un usuario es egresado
    XOR empresa), con subtipos disjuntos en DDL puro.

## A.2 Las correcciones de fondo — el corazón de la defensa

### Correcciones del primer revisor (C-1 a C-4) — heredadas

| Código  | Qué pedía                                                                              | Cómo respondió el esquema                                                                                                       |
|---------|----------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------|
| **C-1** | "No reinventes catálogos; cárgalos del servicio de parámetros."                        | Se eliminaron las 7 tablas paramétricas; los catálogos viven en el servicio institucional. (En v4 además se virtualizó, ver C-6.) |
| **C-2** | "No dupliques datos que ya existen en SGA (egresado) y Ágora (empresa)."               | `egresado` y `empresa` guardan **solo lo mínimo + un id externo**; el resto se consulta on-demand.                              |
| **C-3** | "Valida si `usuario` es pertinente."                                                   | Sí lo es: identidad local de quién usa *este* módulo (JIT provisioning). Aprobado. **No tocar.**                               |
| **C-4** | "¿`secuencia_radicado` es necesaria? El historial podría ser la única fuente de estado."| El estado se sacó de `solicitud_beneficio` (única fuente = `historial_solicitud`). La secuencia se conservaba… **(revisado en C-5).** |

### Correcciones de la sustentación / profesor (C-5 a C-7) — NUEVAS en v4

| Código  | Qué pidió el profesor                                                                          | Cómo respondió v4                                                                                                                                                          |
|---------|------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **C-5** | "Dejen la tabla de radicación como un **objeto de serialización que ya existe en PostgreSQL**." | Se eliminó la tabla `secuencia_radicado` y se reemplazó por una **`SEQUENCE` nativa** (`seq_radicado_beneficio`) + la función `fn_siguiente_radicado()`, usada como `DEFAULT` del campo `radicado`. |
| **C-6** | "**Virtualicen** la tabla `parametro` para sacarla del esquema; hay demasiadas conexiones cíclicas pasando por ella." | Se quitaron **todas** las FK cross-schema a `parametro.parametro`. Las columnas de catálogo son ahora **referencias lógicas** (INTEGER planos) validadas en el MID. `parametro` desaparece del diagrama. |
| **C-7** | "Hagan **excluyente** la relación usuario↔egresado."                                            | Se modeló como **subtipos disjuntos** en DDL puro: discriminador local `tipo_usuario` + `UNIQUE(id, tipo_usuario)` en `usuario`, y **FK compuestas** desde `egresado` (EGR) y `usuario_empresa` (EMP). |

## A.3 El "antes y después" de C-6 (virtualización de parámetros)

El SGA no permite catálogos locales por módulo: los estados/tipos/categorías viven en un
**servicio institucional de parámetros** con jerarquía `AreaTipo → TipoParametro → Parametro`.

- **v2/v3 (antes):** cada columna `*_id` de catálogo era **FK cross-schema** a
  `parametro.parametro(id)`. Resultado: `parametro` aparecía en el diagrama como un **hub**
  con muchas flechas convergiendo (de `usuario`, `empresa`, `beneficio`,
  `historial_solicitud`). El profesor lo señaló como "conexiones cíclicas" que pasaban todas
  por la misma tabla.
- **v4 (después — C-6):** se **eliminan esas FK**. Las columnas `*_id` siguen existiendo pero
  son **referencias lógicas (virtuales)**: un INTEGER que apunta al id de un `Parametro`, sin
  constraint declarada. La **validación de que el id pertenece al `TipoParametro` correcto se
  hace en el MID** (capa de negocio), con el patrón idiomático del SGA
  (`parametro?query=TipoParametroId__Id:{tipo},Id:{id}`). `parametro` ya **no aparece** en el
  diagrama del esquema.

| Columna (tabla)                                  | Referencia lógica al TipoParametro | Valores                                                              |
|--------------------------------------------------|------------------------------------|----------------------------------------------------------------------|
| `empresa.estado_empresa_id`                      | `ESTADO_EMPRESA`                   | ACTIVA, SUSPENDIDA (ciclo de vida local; llegan aprobadas de Ágora)  |
| `empresa.sector_economico_id`                    | `SECTOR_ECONOMICO`                 | TEC, SAL, EDU, IND, COM, FIN, CON, ALI, CON2, OTR                     |
| `beneficio.categoria_beneficio_id`               | `CATEGORIA_BENEFICIO`              | EDUCACION, SALUD, RECREACION, EMPLEO, DESCUENTOS, OTRO                |
| `beneficio.estado_beneficio_id`                  | `ESTADO_BENEFICIO`                 | BORRADOR, PUBLICADO, AGOTADO, VENCIDO, RETIRADO                       |
| `historial_solicitud.estado_anterior_id`/`_nuevo_id` | `ESTADO_SOLICITUD`             | PENDIENTE, EN_REVISION, REQUIERE_INFO, APROBADA, RECHAZADA, CANCELADA |

> **`tipo_usuario` ya NO está en esta tabla.** Por C-7 dejó de ser una referencia a
> parámetros: ahora es un **discriminador local** (`VARCHAR(3)` con CHECK) en la tabla
> `usuario`. Ver A.4.

> **Defensa rápida de C-6:** "No perdemos integridad: los `*_id` siguen siendo referencias a
> parámetros institucionales, solo que la **validación de pertenencia al tipo se delega al
> MID** (que igual debía consultar el servicio para mostrar la etiqueta). A cambio, el modelo
> de datos deja de depender de un schema externo por FK y `parametro` ya no es un hub que
> acopla todo. Es el patrón que el propio SGA usa en producción cuando el servicio de
> parámetros vive en otra base."

## A.4 C-7 en detalle — la relación EXCLUYENTE (subtipos disjuntos)

**Qué se pide:** un `usuario` es **egresado O representante de empresa, nunca ambos**.

**Cómo se garantiza en DDL puro (sin triggers):** patrón clásico de *subtipos disjuntos por
discriminador + FK compuesta*:

1. `usuario` tiene un discriminador **local** `tipo_usuario VARCHAR(3)` con
   `CHECK (tipo_usuario IN ('EGR','EMP','ADM'))`.
2. `usuario` declara `UNIQUE (id, tipo_usuario)` (clave compuesta a la que apuntarán los
   subtipos).
3. `egresado` tiene su propia columna `tipo_usuario` **fijada a `'EGR'`** por
   `CHECK (tipo_usuario = 'EGR')`, y una **FK compuesta**
   `(usuario_id, tipo_usuario) → usuario(id, tipo_usuario)`.
4. `usuario_empresa` hace lo mismo pero **fijada a `'EMP'`**.

**Por qué eso da exclusividad:** un `usuario` tiene **un solo** valor de `tipo_usuario`. Para
insertar un `egresado` la FK exige que ese usuario sea `EGR`; para insertarlo en
`usuario_empresa` la FK exige que sea `EMP`. Como no puede ser las dos cosas a la vez, **no
puede tener perfil de egresado y de empresa simultáneamente**. La exclusividad queda
garantizada por el motor, no por código de aplicación.

> **¿Por qué `tipo_usuario` es local y no una referencia virtual a parámetros (como los demás
> catálogos)?** Porque participa en **constraints estructurales**: una FK compuesta necesita
> un valor **estable y comparable en tiempo de DDL**. El id de un `Parametro` se resuelve en
> runtime por `codigo_abreviacion` y no es estable entre entornos. `TIPO_USUARIO` tiene solo
> 3 valores estructurales fijos (EGR/EMP/ADM), así que modelarlo como `CHECK` local es más
> robusto **y** coherente con sacar `parametro` del esquema (C-6). Los demás catálogos son
> meramente descriptivos (no gobiernan la estructura), por eso sí quedan como referencias
> lógicas.

## A.5 C-5 en detalle — radicado con secuencia nativa

**Antes (v2/v3):** tabla `secuencia_radicado` (una fila por año) con `ultimo_numero`,
incrementada con `SELECT ... FOR UPDATE` dentro de la transacción de negocio.

**Ahora (v4):** dos objetos nativos de PostgreSQL:

- `CREATE SEQUENCE seq_radicado_beneficio` — el "objeto de serialización" que pidió el
  profesor. `nextval()` es **atómico**: garantiza unicidad bajo concurrencia **sin** bloqueos
  explícitos.
- `fn_siguiente_radicado()` — función que arma el formato
  **`BNF-YYYY-NNNNNN`** = `'BNF-' || año actual || '-' || LPAD(nextval, 6)`.
- `solicitud_beneficio.radicado` usa `DEFAULT beneficios_egresados.fn_siguiente_radicado()`:
  el radicado se genera solo al insertar (o el MID puede llamar la función explícitamente).

**Reinicio anual:** la secuencia es continua; si se quiere que `NNNNNN` reinicie cada año, es
una **tarea operativa** (`ALTER SEQUENCE seq_radicado_beneficio RESTART WITH 1;` el 1 de
enero), no un cambio estructural. El formato sigue mostrando el año, así que la trazabilidad
por año se mantiene aunque no se reinicie.

## A.6 Convenciones de modelado OATI/UD que cumple TODO el esquema

- PK siempre `SERIAL` llamada `id`.
- `VARCHAR` **siempre con longitud** explícita. `TEXT` solo para contenido libre
  (descripción, condiciones, justificación, mensaje).
- Campos de auditoría en **todas** las tablas: `activo`, `fecha_creacion`,
  `fecha_modificacion` (ver A.7).
- Constraints nombrados con prefijo: `pk_`, `fk_`, `uq_`, `ck_`.
- Índice sobre toda FK real y sobre campos de búsqueda frecuente (incluidas las referencias
  lógicas `*_id`, que se filtran seguido).
- `COMMENT ON TABLE` en todas; `COMMENT ON COLUMN` en columnas sensibles/no obvias.

## A.7 Campos de auditoría comunes (NO se repiten en cada tabla abajo)

| Campo                | Tipo                               | Para qué sirve                | Por qué es necesario                                                                       |
|----------------------|------------------------------------|-------------------------------|--------------------------------------------------------------------------------------------|
| `activo`             | `BOOLEAN NOT NULL DEFAULT TRUE`    | Borrado lógico (soft delete). | Nada se borra físicamente; conserva trazabilidad. Lineamiento OATI.                        |
| `fecha_creacion`     | `TIMESTAMP NOT NULL DEFAULT NOW()` | Cuándo se creó el registro.   | Auditoría y orden cronológico. Lo asigna el servidor.                                      |
| `fecha_modificacion` | `TIMESTAMP NOT NULL DEFAULT NOW()` | Última modificación.          | Auditoría de cambios. Lo actualiza la capa de aplicación.                                  |

## A.8 Mapa de relaciones (visión rápida — v4)

```
   (servicio institucional de parámetros — FUERA del esquema; SIN FK, solo refs lógicas *_id)
        estado_empresa · sector_economico · categoria_beneficio · estado_beneficio · estado_solicitud

usuario ─┬─1:1─< egresado            ┐
         │        (FK compuesta EGR) ├─ ARCO EXCLUYENTE (C-7): un usuario es EGR xor EMP
         ├─N:M─< usuario_empresa >── empresa   (FK compuesta EMP)
         ├──creador───> beneficio  >── empresa
         ├──autor─────> historial_solicitud
         ├──autor─────> mensaje_solicitud  (chat REQUIERE_INFO)
         └──actor─────> bitacora_acceso_pii

egresado ──< solicitud_beneficio >── beneficio
                  │
                  └──< historial_solicitud   (estado vigente = último registro)

seq_radicado_beneficio + fn_siguiente_radicado()  (secuencia nativa → DEFAULT de radicado, C-5)
v_solicitud_estado_vigente  (vista: último estado por solicitud)
```

> Nota: ya **no** hay un nodo `parametro` con flechas convergiendo (eso era lo que el
> profesor marcaba como cíclico). Las únicas FK que quedan son **locales** (entre tablas del
> propio esquema).

---

# PARTE B — DETALLE TABLA POR TABLA

---

## TABLA: usuario

**Dominio:** identidad · **Correcciones asociadas:** C-3 (aprobada) + C-7 (discriminador del arco exclusivo)

### Para qué sirve

Registra **quién se ha autenticado en este módulo** (JIT provisioning): la primera vez que
una persona entra (egresado vía SGA, o representante de empresa vía Ágora), se inserta aquí a
partir del token WSO2.

### Por qué es necesaria (defensa)

- El revisor preguntó si era pertinente y **concluyó que sí** (C-3).
- SGA y Ágora saben quién existe en *sus* sistemas, pero no **quién usa este módulo ni cuándo**.
- Es el ancla referencial de casi todo: aprobaciones, autoría de beneficios, cambios de
  estado, mensajes y bitácora PII apuntan aquí.
- En v4 su columna `tipo_usuario` además **gobierna el arco exclusivo** egresado/empresa (C-7).

### Campos

| Campo            | Tipo                                  | Para qué sirve                                              | Por qué es necesario                                                                 |
|------------------|---------------------------------------|-------------------------------------------------------------|--------------------------------------------------------------------------------------|
| `id`             | SERIAL PK                             | Identificador interno.                                      | Llave estable para todas las FK locales.                                             |
| `documento`      | VARCHAR(20) NOT NULL **UNIQUE**       | Cédula/identificación (del token).                          | Llave natural; el UNIQUE evita duplicar identidades.                                 |
| `nombre`         | VARCHAR(200) NOT NULL                 | Nombre para mostrar.                                        | Evita ir a SGA/Ágora solo para mostrar el nombre.                                    |
| `correo`         | VARCHAR(150) NOT NULL                 | Correo.                                                     | Notificaciones y contacto.                                                           |
| `tipo_usuario`   | VARCHAR(3) NOT NULL **CHECK EGR/EMP/ADM** | **Discriminador local** del subtipo.                    | **C-7**: local y estable para la FK compuesta del arco exclusivo. (Antes era `tipo_usuario_id` → parametro.) |
| `id_externo`     | VARCHAR(50)                           | ID en el sistema de origen.                                 | Re-vincular con el tercero del SGA o el proveedor de Ágora.                          |
| `sistema_origen` | VARCHAR(20) NOT NULL                  | `SGA`/`AGORA`/`LOCAL`.                                      | De dónde vino la identidad; el CHECK lo restringe.                                   |
| `ultimo_acceso`  | TIMESTAMP                             | Último login.                                               | Métrica de uso y soporte.                                                            |
| *auditoría*      |                                       | `activo`, `fecha_creacion`, `fecha_modificacion` (A.7).     |                                                                                      |

### Constraints e índices

- `uq_documento_usuario (documento)`: una identidad por documento.
- `uq_usuario_id_tipo (id, tipo_usuario)`: **clave compuesta** destino de las FK del arco
  exclusivo (C-7). Trivialmente única (id ya es PK), pero necesaria para poder referenciar la
  pareja `(id, tipo_usuario)`.
- `ck_tipo_usuario`: solo `EGR`/`EMP`/`ADM`.
- `ck_sistema_origen_usuario`: solo `SGA`/`AGORA`/`LOCAL`.
- Índices: `tipo_usuario`; `(sistema_origen, id_externo)`.

### Preguntas probables → respuesta

- **"¿Por qué `tipo_usuario` es local y no referencia a parámetros como los demás catálogos?"**
  Porque participa en una **FK compuesta** (el arco exclusivo, C-7) que necesita un valor
  estable en tiempo de DDL; el id de un parámetro se resuelve en runtime y no es estable entre
  entornos. Además, sacarlo de parámetros es coherente con la virtualización (C-6). Tiene solo
  3 valores estructurales fijos: un CHECK local es lo correcto.
- **"¿No duplica esto al usuario de WSO2/SGA/Ágora?"** No: no guardamos credenciales ni el
  perfil completo, solo la identidad mínima para operar y referenciar acciones.
- **"¿Por qué no el documento como PK?"** PK numérica estable + documento como llave natural
  única; FK livianas y el documento corregible sin romper relaciones.

---

## TABLA: egresado

**Dominio:** identidad · **Correcciones asociadas:** C-2a (no duplicar SGA) + C-7 (subtipo exclusivo)

### Para qué sirve

Perfil del egresado dentro del módulo. Relación **1:1 con `usuario`** tipo EGR. Guarda lo
mínimo; los datos académicos completos vienen del SGA on-demand.

### Por qué es necesaria (defensa)

- `solicitud_beneficio` necesita **a quién apuntar** con integridad referencial local.
- Minimización (C-2a): solo lo imprescindible. `programa_academico`/`facultad` son espejo
  opcional (fuente real = SGA); `fecha_grado` **sí** se guarda porque se verificó que **no hay
  fuente institucional**.
- Es un **subtipo exclusivo** de `usuario` (C-7): no se puede crear un egresado sobre un
  usuario que no sea EGR, ni un egresado puede aparecer también como representante de empresa.

### Campos

| Campo                  | Tipo                                  | Para qué sirve                                    | Por qué es necesario                                                     |
|------------------------|---------------------------------------|---------------------------------------------------|--------------------------------------------------------------------------|
| `id`                   | SERIAL PK                             | Identificador.                                    | FK destino de `solicitud_beneficio`.                                     |
| `usuario_id`           | INTEGER NOT NULL **UNIQUE**           | Parte 1 de la FK compuesta a `usuario`.           | El UNIQUE garantiza el 1:1 usuario↔egresado.                            |
| `tipo_usuario`         | VARCHAR(3) NOT NULL DEFAULT 'EGR' **CHECK = 'EGR'** | Parte 2 de la FK compuesta; fija el subtipo. | **C-7**: obliga a que el usuario referenciado sea EGR → exclusividad.    |
| `codigo_institucional` | VARCHAR(20) NOT NULL **UNIQUE**       | Código estudiantil/carné.                         | Referencia cruzada con el SGA.                                           |
| `programa_academico`   | VARCHAR(150) (nullable)               | Programa cursado.                                 | Espejo del SGA (fuente real allá); útil para mostrar sin llamada remota. |
| `facultad`             | VARCHAR(150) (nullable)               | Facultad.                                         | Espejo del SGA.                                                          |
| `fecha_grado`          | DATE (nullable)                       | Fecha de grado.                                   | **Se almacena local**: no hay fuente institucional (verificado).         |
| `telefono_contacto`    | VARCHAR(20) (nullable)                | Teléfono.                                         | Contacto directo.                                                        |
| *auditoría*            |                                       | `activo`, `fecha_creacion`, `fecha_modificacion`. |                                                                          |

### Constraints e índices

- `uq_usuario_id_egresado`: asegura 1:1 con `usuario`.
- `uq_codigo_institucional_egresado`: un código por egresado.
- `ck_egresado_tipo_usuario`: `tipo_usuario = 'EGR'`.
- `fk_egresado_usuario`: **FK compuesta** `(usuario_id, tipo_usuario) → usuario(id, tipo_usuario)`.

### Preguntas probables → respuesta

- **"¿Cómo se garantiza que un egresado NO sea también empresa?"** (C-7) Por la FK compuesta:
  esta tabla solo enlaza usuarios cuyo `tipo_usuario` sea EGR, y `usuario_empresa` solo enlaza
  EMP. Como un usuario tiene un único `tipo_usuario`, no puede estar en ambas. Exclusividad en
  DDL puro, sin triggers.
- **"¿No viola C-2 guardar programa/facultad?"** Son espejo opcional (nullable); la fuente es
  el SGA, consultado on-demand. Se guardan solo como caché de presentación.
- **"¿Por qué `fecha_grado` sí se guarda?"** Tras búsqueda exhaustiva se confirmó que ni
  `sga_mid` ni `sga_cliente` exponen la fecha de grado; no hay de dónde traerla.

---

## TABLA: empresa

**Dominio:** identidad · **Correcciones asociadas:** C-2b (no duplicar Ágora) + C-6 (refs lógicas)

### Para qué sirve

**Empresa aliada** que publica beneficios. Espejo local mínimo de una empresa de Ágora **más**
el estado de su ciclo de vida *dentro de este módulo*.

### Por qué es necesaria (defensa)

- `beneficio` debe pertenecer a una empresa con integridad referencial local.
- No duplica Ágora: los datos descriptivos se traen on-demand con `agora_id_externo`
  (= `id_proveedor`). Lo único *propio* es el **estado del ciclo de vida en el módulo**, que
  Ágora no conoce.

### Campos

| Campo                  | Tipo                            | Para qué sirve                                    | Por qué es necesario                                                            |
|------------------------|---------------------------------|---------------------------------------------------|---------------------------------------------------------------------------------|
| `id`                   | SERIAL PK                       | Identificador.                                    | FK destino de `beneficio` y `usuario_empresa`.                                  |
| `nit`                  | VARCHAR(20) NOT NULL **UNIQUE** | NIT.                                              | Llave natural; el UNIQUE evita duplicados.                                      |
| `razon_social`         | VARCHAR(200) NOT NULL           | Nombre legal.                                     | Espejo de Ágora para mostrar sin llamada remota.                                |
| `agora_id_externo`     | VARCHAR(50) (nullable)          | `id_proveedor` en Ágora.                          | **Llave para traer datos completos on-demand**; evita duplicar Ágora.           |
| `sector_economico_id`  | INTEGER (nullable)              | **Ref. lógica** → parametro `SECTOR_ECONOMICO`.   | **C-6**: sin FK; el MID valida. En prod se prefiere el CIIU de Ágora.           |
| `estado_empresa_id`    | INTEGER NOT NULL                | **Ref. lógica** → parametro `ESTADO_EMPRESA`.     | **C-6**: ciclo de vida LOCAL (ACTIVA/SUSPENDIDA), sin FK. No hay aprobación interna. |
| `sitio_web`            | VARCHAR(255)                    | Sitio web.                                        | Contacto.                                                                       |
| `correo_contacto`      | VARCHAR(150)                    | Correo.                                           | Contacto operativo.                                                             |
| `telefono_contacto`    | VARCHAR(20)                     | Teléfono.                                         | Contacto operativo.                                                             |
| `direccion`            | VARCHAR(255)                    | Dirección.                                        | Contacto.                                                                       |
| *auditoría*            |                                 | `activo`, `fecha_creacion`, `fecha_modificacion`. |                                                                                 |

> **OJO — campos eliminados a propósito:** v4 (post-sustentación) tenía `fecha_aprobacion` y
> `usuario_aprobador_id`. **Se eliminaron** porque las empresas **llegan ya aprobadas desde
> Ágora** (autenticación WSO2 + JIT): no existe un flujo de aprobación interno en el módulo, así
> que no hay "quién aprobó" ni "cuándo se aprobó" que registrar. Esto además **elimina uno de
> los caminos usuario↔empresa**, simplificando el modelo (menos conexiones, en línea con C-6).

### Constraints e índices

- `uq_nit_empresa`: una empresa por NIT.
- **Sin FK a `usuario`** (se quitó `usuario_aprobador_id`). `sector_economico_id` y
  `estado_empresa_id` son **referencias lógicas sin FK** (C-6).
- Índices: `estado_empresa_id`, `sector_economico_id`.

### Preguntas probables → respuesta

- **"¿No falta registrar quién aprobó la empresa?"** No: las empresas **llegan ya aprobadas
  desde Ágora**; en este módulo no hay aprobación interna, por eso se eliminaron
  `usuario_aprobador_id` y `fecha_aprobacion`. El único estado propio que tiene sentido es el
  ciclo de vida local (ACTIVA/SUSPENDIDA), por si hay que suspender una empresa aquí.
- **"Entonces ¿cuántas relaciones usuario↔empresa quedan?"** Una sola: `usuario_empresa`
  (N:M operativa, quién opera la empresa). La de auditoría administrativa se eliminó.
- **"¿Por qué `estado_empresa_id` no tiene FK?"** (C-6) Porque virtualizamos `parametro`: es
  una referencia lógica al servicio institucional y el MID valida el tipo. Así `parametro` deja
  de ser un hub acoplado al esquema.
- **"¿`sector_economico` no debería venir de Ágora (CIIU)?"** Sí, en producción se usa
  `proveedor_actividad_ciiu`; la referencia local queda como respaldo/clasificación simple.

---

## TABLA: usuario_empresa

**Dominio:** identidad · **Correcciones asociadas:** C-2c (validar contra Ágora) + C-7 (subtipo exclusivo)

### Para qué sirve

Relación **N:M** entre usuarios tipo EMP y empresas: qué personas operan en nombre de qué
empresa, con su cargo, y cuál es el contacto principal.

### Por qué es necesaria (defensa)

- Una empresa puede tener **varios** representantes, y una persona podría representar a más de
  una empresa → tabla puente N:M.
- Ágora (`proveedor_representante_legal`) solo admite **un** representante por proveedor y no
  se vincula con WSO2/terceros → **insuficiente**. La relación operativa vive aquí (JIT).
- Es el **otro lado del arco exclusivo** (C-7): solo enlaza usuarios EMP, nunca egresados.

### Campos

| Campo          | Tipo                                  | Para qué sirve                                    | Por qué es necesario                                     |
|----------------|---------------------------------------|---------------------------------------------------|----------------------------------------------------------|
| `id`           | SERIAL PK                             | Identificador.                                    | Llave de la fila de relación.                            |
| `usuario_id`   | INTEGER NOT NULL                      | Parte 1 de la FK compuesta a `usuario`.           | Lado "persona" de la relación.                           |
| `tipo_usuario` | VARCHAR(3) NOT NULL DEFAULT 'EMP' **CHECK = 'EMP'** | Parte 2 de la FK compuesta; fija el subtipo. | **C-7**: obliga a que el usuario sea EMP → exclusividad. |
| `empresa_id`   | INTEGER NOT NULL                      | FK → `empresa`.                                   | Lado "empresa" de la relación.                           |
| `cargo`        | VARCHAR(100) (nullable)               | Cargo de la persona.                              | Contexto operativo/contractual.                          |
| `es_principal` | BOOLEAN NOT NULL DEFAULT FALSE        | Marca al contacto principal.                      | A quién dirigir comunicaciones cuando hay varios.        |
| *auditoría*    |                                       | `activo`, `fecha_creacion`, `fecha_modificacion`. |                                                          |

### Constraints e índices

- `uq_usuario_empresa (usuario_id, empresa_id)`: evita duplicar el mismo vínculo.
- `ck_usuario_empresa_tipo_usuario`: `tipo_usuario = 'EMP'`.
- `fk_usuario_empresa_usuario`: **FK compuesta** `(usuario_id, tipo_usuario) → usuario(id, tipo_usuario)`.
- `fk_usuario_empresa_empresa → empresa`. Índice: `empresa_id`.

### Preguntas probables → respuesta

- **"¿Cómo impide esto que un egresado opere una empresa?"** (C-7) La FK compuesta solo enlaza
  usuarios cuyo `tipo_usuario` sea EMP; un EGR no puede insertarse aquí. Espejo exacto de la
  restricción en `egresado`.
- **"¿Ágora no resuelve la relación?"** No: `proveedor_representante_legal` permite un solo
  representante y no enlaza con WSO2. Por eso es local.

---

## TABLA: beneficio

**Dominio:** negocio · **Correcciones asociadas:** C-1/C-6 (catálogos como refs lógicas)

### Para qué sirve

**Entidad central**: el beneficio que una empresa publica para los egresados (título,
descripción, condiciones, vigencia y cupos).

### Por qué es necesaria (defensa)

Es el objeto de negocio principal; sin ella no hay catálogo ni solicitudes. La defensa se
centra en que sus catálogos (categoría/estado) **no son tablas locales ni FK**, sino
referencias lógicas a parámetros (C-1 + C-6).

### Campos

| Campo                    | Tipo                    | Para qué sirve                                       | Por qué es necesario                                                    |
|--------------------------|-------------------------|------------------------------------------------------|-------------------------------------------------------------------------|
| `id`                     | SERIAL PK               | Identificador.                                       | FK destino de `solicitud_beneficio`.                                    |
| `empresa_id`             | INTEGER NOT NULL        | **FK LOCAL** → `empresa`.                            | Pertenencia: de qué empresa es el beneficio.                            |
| `categoria_beneficio_id` | INTEGER NOT NULL        | **Ref. lógica** → parametro `CATEGORIA_BENEFICIO`.   | Clasificación sin catálogo local ni FK (C-1/C-6).                       |
| `estado_beneficio_id`    | INTEGER NOT NULL        | **Ref. lógica** → parametro `ESTADO_BENEFICIO`.      | Ciclo de vida del beneficio sin catálogo local ni FK (C-1/C-6).         |
| `titulo`                 | VARCHAR(200) NOT NULL   | Título.                                              | Encabezado en el catálogo.                                              |
| `descripcion`            | TEXT NOT NULL           | Descripción libre.                                   | Contenido sin límite fijo → `TEXT`.                                     |
| `condiciones`            | TEXT NOT NULL           | Condiciones/requisitos.                              | Contenido libre; el MID exige separarla de la descripción (RN-008b).    |
| `fecha_inicio`           | DATE NOT NULL           | Inicio de vigencia.                                  | Define cuándo aplica.                                                   |
| `fecha_fin`              | DATE NOT NULL           | Fin de vigencia.                                     | CHECK `fecha_fin > fecha_inicio` evita rangos inválidos.                |
| `cupos_total`            | INTEGER NOT NULL        | Cupos ofrecidos.                                     | CHECK `>= 1`: no se publica sin cupos.                                  |
| `cupos_disponibles`      | INTEGER NOT NULL        | Cupos restantes.                                     | CHECK `0 <= disp <= total`: integridad del inventario.                  |
| `imagen_url`             | VARCHAR(500) (nullable) | Imagen.                                              | Presentación.                                                           |
| `fecha_publicacion`      | TIMESTAMP (nullable)    | Cuándo pasó a PUBLICADO.                             | Trazabilidad.                                                           |
| `usuario_creador_id`     | INTEGER NOT NULL        | **FK LOCAL** → `usuario`.                            | Auditoría de autoría. Distinto de empresa→beneficio.                    |
| *auditoría*              |                         | `activo`, `fecha_creacion`, `fecha_modificacion`.    |                                                                         |

### Constraints e índices

- CHECKs: `ck_fechas_beneficio`, `ck_cupos_total_beneficio`, `ck_cupos_disponibles_beneficio`.
- **FK reales:** `empresa`, `usuario_creador`. `categoria`/`estado` son refs lógicas sin FK (C-6).
- Índices: `estado_beneficio_id`, `categoria_beneficio_id`, `empresa_id`,
  `(fecha_fin, cupos_disponibles)`.

### Preguntas probables → respuesta

- **"¿Por qué `categoria`/`estado` no tienen FK?"** (C-6) Virtualización: son referencias
  lógicas a parámetros; el MID valida el tipo. Evita acoplar el esquema al servicio externo.
- **"¿Por qué `usuario_creador_id` si ya está `empresa_id`?"** `empresa_id` = jerarquía de
  negocio; `usuario_creador_id` = auditoría (qué persona lo creó). Ambos necesarios.
- **"¿`descripcion`/`condiciones` por qué `TEXT`?"** Contenido libre sin longitud fija; el
  estándar reserva `VARCHAR(n)` para campos acotados.
- **"¿Cupos en concurrencia?"** El CHECK garantiza el dato; la atomicidad del descuento se
  maneja en aplicación (pendiente endpoint dedicado, RN-002b/c).

---

## OBJETO: seq_radicado_beneficio + fn_siguiente_radicado()

**Dominio:** soporte · **Corrección asociada:** C-5 (¡la más preguntada en esta sustentación!)

### Para qué sirve

Generar **radicados únicos** con formato **`BNF-YYYY-NNNNNN`**. Sustituyen a la antigua tabla
`secuencia_radicado`.

- `seq_radicado_beneficio`: **secuencia nativa** de PostgreSQL (el "objeto de serialización"
  que pidió el profesor). `nextval()` entrega consecutivos de forma **atómica**.
- `fn_siguiente_radicado()`: función que arma el radicado
  `'BNF-' || año || '-' || LPAD(nextval, 6)`.
- Se usa como `DEFAULT` del campo `solicitud_beneficio.radicado`.

### Por qué es así (defensa)

- **El profesor pidió explícitamente** usar el objeto de serialización nativo de PostgreSQL en
  vez de una tabla-contador manual. Una `SEQUENCE` es justamente eso: el mecanismo nativo,
  atómico y libre de bloqueos para generar números únicos crecientes.
- **Concurrencia:** `nextval()` nunca entrega el mismo número a dos transacciones, **sin**
  necesidad de `SELECT FOR UPDATE`. La unicidad final del radicado la refuerzan el `UNIQUE` y
  el `CHECK` de formato en `solicitud_beneficio`.
- **Menos superficie:** se eliminan una tabla, su semilla y la ruta CRUD
  `POST /siguiente/:anio`. El radicado se genera solo (DEFAULT) o con
  `SELECT beneficios_egresados.fn_siguiente_radicado();`.

### Preguntas probables → respuesta

- **"¿Y el reinicio por año? Una SEQUENCE no reinicia sola."** Correcto: es una **decisión
  operativa**, no estructural. Si se quiere `NNNNNN` por año, un job el 1 de enero ejecuta
  `ALTER SEQUENCE seq_radicado_beneficio RESTART WITH 1;`. El año ya va en el formato, así que
  la trazabilidad anual se mantiene aunque el consecutivo sea continuo.
- **"¿Qué pasa si pasa de 999999?"** El `CHECK` de 6 dígitos fallaría; con reinicio anual no se
  alcanza. Si se previera ese volumen, se amplía el `LPAD`/regex a 7 dígitos.
- **"¿Por qué antes era tabla y ahora secuencia?"** Antes (C-4a) se argumentó control
  transaccional fino; el profesor pidió simplificar al objeto nativo (C-5). La secuencia cubre
  unicidad y concurrencia con menos código y es el idioma de PostgreSQL.
- **"Si dos solicitudes fallan/rollback, ¿se 'pierden' números?"** Sí, una SEQUENCE puede dejar
  huecos (es el comportamiento estándar y aceptable: garantiza unicidad, no contigüidad
  perfecta). Para un radicado eso es irrelevante.

---

## TABLA: solicitud_beneficio

**Dominio:** negocio · **Correcciones asociadas:** C-4b (estado fuera) + C-5 (radicado por secuencia)

### Para qué sirve

Registra que un **egresado solicita un beneficio**. Es la transacción central del módulo.

### Por qué es necesaria (defensa)

Es el hecho de negocio (quién pidió qué y cuándo). Lo defendible es **lo que NO tiene**: ya
**no** guarda `estado_solicitud_id` (C-4b), y el `radicado` se genera con la **secuencia
nativa** (C-5), no con una tabla-contador.

### Campos

| Campo                   | Tipo                                                | Para qué sirve                                    | Por qué es necesario                                                                |
|-------------------------|-----------------------------------------------------|---------------------------------------------------|-------------------------------------------------------------------------------------|
| `id`                    | SERIAL PK                                           | Identificador.                                    | FK destino de historial y mensajes.                                                 |
| `radicado`              | VARCHAR(20) NOT NULL **UNIQUE** DEFAULT `fn_siguiente_radicado()` | Radicado `BNF-YYYY-NNNNNN`.          | Identificador oficial; lo genera la **secuencia nativa** (C-5). CHECK valida formato. |
| `egresado_id`           | INTEGER NOT NULL                                    | FK → `egresado`.                                  | Quién solicita.                                                                     |
| `beneficio_id`          | INTEGER NOT NULL                                    | FK → `beneficio`.                                 | Qué solicita.                                                                       |
| `datos_complementarios` | JSONB (nullable)                                    | Datos extra del formulario.                       | Cada beneficio puede pedir campos distintos → JSON flexible.                        |
| `fecha_solicitud`       | TIMESTAMP NOT NULL DEFAULT NOW()                    | Cuándo se radicó.                                 | Orden y métricas.                                                                   |
| *auditoría*             |                                                     | `activo`, `fecha_creacion`, `fecha_modificacion`. |                                                                                     |

> **OJO — campo eliminado a propósito:** v1 tenía `estado_solicitud_id` aquí. **Se eliminó**:
> el estado vigente sale del último registro de `historial_solicitud` (C-4b).

### Constraints e índices

- `uq_radicado_solicitud_beneficio` + `ck_radicado_solicitud_beneficio` (regex
  `^BNF-[0-9]{4}-[0-9]{6}$`).
- FKs: `egresado`, `beneficio`.
- Índices: `egresado_id`, `beneficio_id`, `fecha_solicitud`.

### Preguntas probables → respuesta

- **"¿Cómo se genera el radicado?"** (C-5) Por `DEFAULT` con la función
  `fn_siguiente_radicado()`, que usa la secuencia nativa. El MID también puede invocarla
  explícitamente si necesita el valor antes del INSERT.
- **"¿Dónde está el estado?"** No se almacena aquí. Es el `estado_nuevo_id` del registro más
  reciente en `historial_solicitud` (vista `v_solicitud_estado_vigente`). (C-4b.)
- **"¿Por qué `datos_complementarios` en JSONB?"** Los requisitos varían por beneficio; un
  JSONB evita crear columnas/tablas por cada formulario.

---

## TABLA: historial_solicitud

**Dominio:** negocio · **Correcciones asociadas:** C-4b (única fuente de estado) + C-6 (refs lógicas)

### Para qué sirve

Bitácora de **todas las transiciones de estado** de cada solicitud. Cada cambio es un INSERT.
**El estado vigente = el registro con mayor `fecha_cambio`.**

### Por qué es necesaria (defensa)

- El revisor pidió que "el último estado sea el actual": esta tabla es la **única fuente de
  verdad** del estado (C-4b).
- Da trazabilidad completa (quién, cuándo, de qué estado a cuál, por qué). Un solo campo de
  estado en `solicitud_beneficio` no daría historia y podría desincronizarse.

### Campos

| Campo                    | Tipo                             | Para qué sirve                                                | Por qué es necesario                                                 |
|--------------------------|----------------------------------|---------------------------------------------------------------|----------------------------------------------------------------------|
| `id`                     | SERIAL PK                        | Identificador.                                                | Llave del registro de transición.                                    |
| `solicitud_beneficio_id` | INTEGER NOT NULL                 | **FK LOCAL** → `solicitud_beneficio`.                         | A qué solicitud pertenece la transición.                             |
| `estado_anterior_id`     | INTEGER (nullable)               | **Ref. lógica** → parametro `ESTADO_SOLICITUD`. Estado origen. | NULL en el primer registro. Parte del par origen→destino. Sin FK (C-6). |
| `estado_nuevo_id`        | INTEGER NOT NULL                 | **Ref. lógica** → parametro `ESTADO_SOLICITUD`. Estado destino.| Nunca NULL; define el estado vigente. Sin FK (C-6).                  |
| `usuario_id`             | INTEGER NOT NULL                 | **FK LOCAL** → `usuario`.                                     | Auditoría de acción: quién ejecutó el cambio.                        |
| `justificacion`          | TEXT (nullable)                  | Motivo del cambio.                                            | Obligatorio en negocio para rechazos (RN-003).                       |
| `fecha_cambio`           | TIMESTAMP NOT NULL DEFAULT NOW() | Momento de la transición.                                    | Define el orden; el más reciente es el estado vigente.               |
| *auditoría*              |                                  | `activo`, `fecha_creacion`, `fecha_modificacion`.            |                                                                      |

### Constraints e índices

- **FK reales:** `solicitud_beneficio`, `usuario`. `estado_anterior_id`/`estado_nuevo_id` son
  referencias lógicas **sin FK** (C-6).
- Índices: `solicitud_beneficio_id`, `fecha_cambio`, y
  `(solicitud_beneficio_id, fecha_cambio DESC)` para el estado vigente.

### Preguntas probables → respuesta

- **"¿Por qué dos columnas de estado a parámetros y ninguna con FK?"** Son origen y destino de
  la transición (roles semánticos distintos, no redundancia). En v4 además **no llevan FK**
  (C-6): son referencias lógicas validadas en el MID, para que `parametro` no acople el
  esquema.
- **"¿Por qué `estado_anterior_id` permite NULL?"** En el primer registro (creación) no hay
  estado previo.
- **"¿No es más simple un campo de estado en la solicitud?"** Pierde la historia y permite
  inconsistencias; el revisor pidió este modelo (C-4b).
- **"¿Hay un tercer camino usuario↔solicitud aquí?"** Sí: `usuario_id` = quién cambió el
  estado (auditoría), distinto de `egresado`→solicitud (quién la creó).

---

## TABLA: mensaje_solicitud

**Dominio:** negocio

### Para qué sirve

Hilo de **mensajes entre empresa y egresado** cuando una solicitud queda en `REQUIERE_INFO`.

### Por qué es necesaria (defensa)

- El estado `REQUIERE_INFO` necesita un **canal** con trazabilidad; sin esta tabla sería un
  estado muerto.
- Mantiene la comunicación **dentro del expediente** (ligada al radicado), con autoría y fecha.

### Campos

| Campo                    | Tipo                             | Para qué sirve                                    | Por qué es necesario                    |
|--------------------------|----------------------------------|---------------------------------------------------|-----------------------------------------|
| `id`                     | SERIAL PK                        | Identificador.                                    | Llave del mensaje.                      |
| `solicitud_beneficio_id` | INTEGER NOT NULL                 | FK → `solicitud_beneficio`.                       | Ata el mensaje al expediente.           |
| `usuario_id`             | INTEGER NOT NULL                 | FK → `usuario`.                                   | Autor del mensaje (empresa o egresado). |
| `mensaje`                | TEXT NOT NULL                    | Contenido.                                        | Texto libre → `TEXT`.                   |
| `fecha_envio`            | TIMESTAMP NOT NULL DEFAULT NOW() | Cuándo se envió.                                  | Ordena el hilo.                         |
| *auditoría*              |                                  | `activo`, `fecha_creacion`, `fecha_modificacion`. |                                         |

### Constraints e índices

- FKs: `solicitud_beneficio`, `usuario`. Índice: `(solicitud_beneficio_id, fecha_envio)`.

### Preguntas probables → respuesta

- **"¿No basta `justificacion` del historial?"** No: la justificación es puntual del cambio de
  estado; aquí hay un ida-y-vuelta entre dos partes con varios mensajes.

---

## TABLA: bitacora_acceso_pii

**Dominio:** soporte/cumplimiento legal

### Para qué sirve

Registra los **accesos a datos personales** (quién accedió a qué dato, cuándo, desde qué IP).
Requisito de la **Ley 1581 de 2012** (protección de datos personales, Colombia). Retención
mínima 6 meses.

### Por qué es necesaria (defensa)

- El módulo maneja PII de egresados; la ley exige poder demostrar **quién accedió**.
- No es opcional: sin ella el módulo no respondería ante una auditoría de tratamiento de datos.

### Campos

| Campo          | Tipo                             | Para qué sirve                       | Por qué es necesario              |
|----------------|----------------------------------|--------------------------------------|-----------------------------------|
| `id`           | SERIAL PK                        | Identificador.                       | Llave del evento.                 |
| `usuario_id`   | INTEGER NOT NULL                 | FK → `usuario`.                      | Quién accedió.                    |
| `recurso_tipo` | VARCHAR(50) NOT NULL             | Tipo de recurso (p. ej. `egresado`). | Qué clase de dato.                |
| `recurso_id`   | INTEGER (nullable)               | ID del recurso.                      | A qué registro concreto.          |
| `accion`       | VARCHAR(50) NOT NULL             | Acción (lectura, exportación…).      | Qué se hizo con el dato.          |
| `direccion_ip` | VARCHAR(45) (nullable)           | IP de origen.                        | Trazabilidad (45 = soporta IPv6). |
| `user_agent`   | VARCHAR(500) (nullable)          | Cliente/navegador.                   | Contexto del acceso.              |
| `detalle`      | JSONB (nullable)                 | Detalle adicional.                   | Información variable del evento.  |
| `fecha_evento` | TIMESTAMP NOT NULL DEFAULT NOW() | Cuándo ocurrió.                      | Momento exacto del acceso.        |
| *auditoría*    |                                  | `activo`, `fecha_creacion`, `fecha_modificacion`. |                      |

### Constraints e índices

- FK: `usuario`. Índices: `usuario_id`, `(recurso_tipo, recurso_id)`, `fecha_evento`.

### Preguntas probables → respuesta

- **"Es 'inmutable', ¿para qué `activo`/`fecha_modificacion`?"** Conceptualmente es de
  solo-escritura, pero el lineamiento OATI exige auditoría en **todas** las tablas
  (uniformidad). `activo` permite deshabilitar lógicamente sin borrar.
- **"¿Por qué `direccion_ip` VARCHAR(45)?"** Para soportar IPv6 completo.

---

# PARTE C — VISTA Y APROVISIONAMIENTO

## VISTA: v_solicitud_estado_vigente

### Para qué sirve

Devuelve el **estado actual de cada solicitud** sin repetir en cada consulta la lógica
`DISTINCT ON (solicitud) ... ORDER BY fecha_cambio DESC`.

### Por qué es necesaria (defensa)

Como el estado vive en el historial (C-4b), calcular "el último estado" en cada listado sería
repetitivo y propenso a errores. La vista centraliza esa lógica una sola vez.

### Qué devuelve

`solicitud_id`, `radicado`, `egresado_id`, `beneficio_id`, `estado_actual_id`
(= `estado_nuevo_id` del último registro), `fecha_ultimo_estado`, `usuario_ultimo_cambio`,
`justificacion_ultimo_cambio`.

---

## APROVISIONAMIENTO en el servicio institucional de parámetros (NO es parte del esquema)

Tras la virtualización (C-6), `parametro`/`tipo_parametro` **no se referencian por FK** desde
este esquema. El DDL incluye, al final y **comentado**, un script-guía de los
`TipoParametro`/`Parametro` que deben existir en el servicio institucional
(`parametros_crud` de la OATI) para que el MID resuelva las referencias lógicas `*_id`:

- **6 `tipo_parametro`**: `ESTADO_EMPRESA`, `ESTADO_BENEFICIO`, `ESTADO_SOLICITUD`,
  `CATEGORIA_BENEFICIO`, `SECTOR_ECONOMICO`, `PARAMETRO_SISTEMA`.
- Sus valores `parametro` hijos (resolviendo `tipo_parametro_id` por subconsulta de
  `codigo_abreviacion`).

> **`TIPO_USUARIO` ya NO se aprovisiona como parámetro**: por C-7 es un discriminador local
> (`usuario.tipo_usuario`, CHECK EGR/EMP/ADM).

### Pendientes operativos al cargar en producción (NO son fallas de diseño)

1. El `INSERT INTO tipo_parametro` no incluye `area_tipo_id`; en `parametros_crud` suele ser
   FK obligatoria → añadir el área de Egresados real al cargar.
2. La semilla no tiene `ON CONFLICT`; correrla dos veces duplicaría. Ejecutar una sola vez o
   añadir idempotencia.
3. Confirmar con OATI que los `codigo_abreviacion` genéricos no colisionen con tipos
   existentes.

---

# PARTE D — CHULETA RÁPIDA (resumen de una línea por objeto)

| Objeto                               | Es para…                                        | Es necesario / se modeló así porque…                                                                       |
|--------------------------------------|-------------------------------------------------|-------------------------------------------------------------------------------------------------------------|
| `usuario`                            | Identidad local de quién usa el módulo (JIT).   | SGA/Ágora no saben quién opera *este* módulo; `tipo_usuario` (local) ancla el arco exclusivo (C-7).        |
| `egresado`                           | Perfil mínimo del egresado (1:1 con usuario).   | Subtipo **exclusivo** EGR (FK compuesta, C-7); guarda solo lo mínimo + `fecha_grado` (sin fuente externa).  |
| `empresa`                            | Espejo mínimo de empresa + estado en el módulo. | `beneficio` necesita dueño; catálogos como refs lógicas sin FK (C-6); estado propio que Ágora no conoce.    |
| `usuario_empresa`                    | Relación N:M usuario↔empresa.                   | Subtipo **exclusivo** EMP (FK compuesta, C-7); Ágora solo admite 1 representante sin vínculo WSO2.          |
| `beneficio`                          | El beneficio publicado (entidad central).       | Objeto de negocio principal; categoría/estado como refs lógicas sin FK (C-6).                               |
| `seq_radicado_beneficio` + función   | Generar radicados `BNF-YYYY-NNNNNN`.            | **Objeto de serialización nativo** de PostgreSQL (C-5): atómico, sin tabla-contador.                        |
| `solicitud_beneficio`                | El hecho: egresado pide beneficio.              | Transacción central; sin estado embebido (C-4b); radicado por secuencia nativa (C-5).                       |
| `historial_solicitud`                | Bitácora de estados; último = vigente.          | Única fuente de verdad del estado (C-4b); estados como refs lógicas sin FK (C-6).                           |
| `mensaje_solicitud`                  | Chat empresa↔egresado en REQUIERE_INFO.         | El estado REQUIERE_INFO necesita canal con trazabilidad dentro del expediente.                              |
| `bitacora_acceso_pii`                | Auditoría de acceso a datos personales.         | Exigido por la Ley 1581; demuestra quién accedió a PII.                                                      |
| `v_solicitud_estado_vigente` (vista) | Estado actual por solicitud.                    | Centraliza el cálculo del último estado (no es tabla).                                                       |

---

# PARTE E — IMPACTO EN EL CÓDIGO (pendiente de sincronizar tras v4)

Los cambios estructurales de v4 **aún no se reflejan** en los microservicios. Al retomar el
backend hay que sincronizar:

1. **CRUD (`sga_crud_beneficios_egresados`)**:
   - `models/usuario.go`: `tipo_usuario_id int` → `tipo_usuario string` (EGR/EMP/ADM).
   - `models/egresado.go` y `models/usuario_empresa.go`: agregar campo `tipo_usuario` y la FK
     compuesta (en Beego ORM, vía `rel(fk)` no aplica directo; se valida la pareja o se deja la
     restricción a nivel de BD).
   - Eliminar modelo/controlador/rutas de `secuencia_radicado`; el radicado lo da el DEFAULT o
     `SELECT fn_siguiente_radicado()`.
   - `db/schema.sql` del repo está en v2: **reemplazarlo por v4**.
2. **MID (`sga_mid_beneficios_egresados`)**:
   - Ya no llamar `POST /v1/secuencia_radicado/siguiente/:anio` (no existe). El radicado lo
     genera la BD.
   - `parametros_service.go`: la validación de tipo de los `*_id` (estado/categoría/sector)
     ahora es **obligatoria en el MID** (ya no hay FK que respalde). Anclar el tipo en la query
     de desreferencia.
   - JIT provisioning: al crear `usuario`, fijar `tipo_usuario` según el origen del token
     (EGR si viene de SGA, EMP si viene de Ágora) — esto **habilita** el arco exclusivo.
3. **Documentación del proyecto**: reglas C-1/C-4a/C-4b siguen vigentes; añadir C-5/C-6/C-7
   y marcar que `tipo_usuario` dejó de ser parámetro y que `secuencia_radicado` ya no es tabla.
   (Nota: ya propagado a las specs el 2026-07-08.)
