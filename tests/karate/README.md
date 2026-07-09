# Pruebas de API con Karate — `sga_crud_beneficios_egresados`

Suite de pruebas funcionales del contrato HTTP del CRUD usando
[Karate](https://github.com/karatelabs/karate) (v1.4.1, Java 11+). Prueba el
servicio de **caja negra** contra PostgreSQL real — sin mocks: el CRUD no tiene
dependencias externas además de la base de datos.

## Cobertura

| Feature | Qué valida |
|---|---|
| `01-getall-contrato` | El contrato de listado compartido por todos los GetAll (scaffold de terceros_crud, regla 10): `query` (con dot-notation a relaciones), `fields`, `sortby`/`order`, `limit`/`offset`, idioma `[{}]` para lista vacía y 400 por parámetros inconsistentes |
| `02-beneficio-cupos` | RN-002b/c: descuento/devolución **atómica** de cupos con guard en el UPDATE (no baja de 0, no supera `cupos_total`) y borrado lógico |
| `03-solicitud-radicado-historial` | C-5: radicado `BNF-YYYY-NNNNNN` generado por la BD (`fn_siguiente_radicado`); C-4b: `historial_solicitud` como única fuente de estado (`/vigente` = último registro, bitácora en orden descendente, 404 sin historial) |
| `04-crud-basico-usuario` | Ciclo alta→consulta→actualización→borrado lógico, UNIQUE de `documento` y exclusión de inactivos con `Activo:true` |

## Cómo ejecutar

**Prerrequisitos:** Go, Java 11+, Maven y PostgreSQL corriendo.

```powershell
cd tests\karate
.\run_pruebas.ps1
```

El script usa una **BD exclusiva de pruebas** (`beneficios_egresados_pruebas`),
que crea con `db/schema.sql` si no existe y re-siembra con `db/seed_pruebas.sql`
en cada corrida — la BD de desarrollo (`beneficios_egresados`) **no se toca**.
Luego compila y levanta el CRUD contra esa BD y corre `mvn test`. Reporte HTML
en `target/karate-reports/karate-summary.html`.

Si el CRUD ya está corriendo en `:8080`, basta con `mvn test`
(URL sobreescribible con `mvn test "-Dcrud.url=http://otro:8080/v1"`).

## Notas

- Los features corren en orden fijo y un solo hilo (comparten la BD).
- Los ids de parámetros institucionales usados por los datos de prueba
  (7199+, C-1) son los reales del servicio de parámetros (migración
  2026-07-07), definidos en `karate-config.js`.
- La suite del MID (`sga_mid_beneficios_egresados/tests/karate`) cubre las
  reglas de negocio de punta a punta; esta cubre el contrato bajo el MID.
