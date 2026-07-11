<#
  smoke_test.ps1 — Prueba end-to-end del CRUD (genera beneficio + solicitud).

  Requisitos previos:
    1. PostgreSQL con el esquema v4 aplicado (db/schema.sql) y datos base
       (db/seed_pruebas.sql) cargados.
    2. El CRUD corriendo:  go run main.go   (escucha en :8080 por defecto)

  Uso:
    pwsh scripts/smoke_test.ps1
    pwsh scripts/smoke_test.ps1 -BaseUrl http://localhost:8080/v1

  Qué hace (replica el flujo de CrearSolicitud del MID, pero directo al CRUD):
    POST /beneficio              → publica un beneficio de la empresa demo (id=1)
    POST /solicitud-beneficio    → el egresado demo (id=1) solicita ese beneficio
                                   (el radicado lo genera la SEQUENCE nativa, C-5)
    POST /historial-solicitud    → registro inicial de estado PENDIENTE (C-4b)
    GET  .../vigente             → confirma el estado vigente derivado del historial
#>
param(
  [string]$BaseUrl = "http://localhost:8080/v1"
)

$ErrorActionPreference = "Stop"

# Placeholders de parámetros lógicos (en prod los resuelve el MID vía parametros service)
$CATEGORIA_EDUCACION   = 1
$ESTADO_BEN_PUBLICADO  = 2
$ESTADO_SOL_PENDIENTE  = 1

function Post($path, $body) {
  $json = $body | ConvertTo-Json -Depth 8
  return Invoke-RestMethod -Method Post -Uri "$BaseUrl$path" -ContentType "application/json" -Body $json
}
function Get-($path) { return Invoke-RestMethod -Method Get -Uri "$BaseUrl$path" }

Write-Host "== Smoke test CRUD @ $BaseUrl ==" -ForegroundColor Cyan

# 1) Publicar un beneficio --------------------------------------------------
$beneficio = @{
  empresa                = @{ id = 1 }                  # empresa demo del seed
  categoria_beneficio_id = $CATEGORIA_EDUCACION
  estado_beneficio_id    = $ESTADO_BEN_PUBLICADO
  titulo                 = "Beca de especialización 50% — Demo"
  descripcion            = "Descuento del 50% en programas de posgrado para egresados UD."
  condiciones            = "Egresado activo, promedio >= 3.8, cupos limitados."
  fecha_inicio           = "2026-06-01T00:00:00Z"
  fecha_fin              = "2026-12-31T00:00:00Z"
  cupos_total            = 50
  cupos_disponibles      = 50
  usuario_creador        = @{ id = 1 }                  # representante de la empresa
}
$rBen = Post "/beneficio" $beneficio
$benId = $rBen.id
Write-Host ("[OK] Beneficio publicado  -> id = {0}" -f $benId) -ForegroundColor Green

# 2) Crear la solicitud (radicado lo pone la BD, C-5) -----------------------
$solicitud = @{
  egresado              = @{ id = 1 }                   # egresado demo del seed
  beneficio             = @{ id = $benId }
  datos_complementarios = '{"motivacion":"Interesado en continuar estudios de posgrado"}'
}
$rSol = Post "/solicitud-beneficio" $solicitud
$solId    = $rSol.id
$radicado = $rSol.radicado
Write-Host ("[OK] Solicitud creada     -> id = {0}, radicado = {1}" -f $solId, $radicado) -ForegroundColor Green

# 3) Historial inicial: estado PENDIENTE (única fuente de estado, C-4b) ------
$historial = @{
  solicitud_beneficio = @{ id = $solId }
  estado_nuevo_id     = $ESTADO_SOL_PENDIENTE          # estado_anterior_id = null (inicial)
  usuario             = @{ id = 2 }                     # el propio egresado
  justificacion       = "Creación de la solicitud (estado inicial)."
}
$null = Post "/historial-solicitud" $historial
Write-Host "[OK] Historial inicial     -> PENDIENTE" -ForegroundColor Green

# 4) Verificar estado vigente derivado del historial ------------------------
$vigente = Get- ("/historial-solicitud/solicitud/{0}/vigente" -f $solId)
Write-Host "`n== Estado vigente de la solicitud ==" -ForegroundColor Cyan
$vigente | ConvertTo-Json -Depth 6

Write-Host "`nOK: beneficio y solicitud generados de punta a punta." -ForegroundColor Green
