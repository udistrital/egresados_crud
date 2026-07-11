# =============================================================================
# run_pruebas.ps1 — Orquesta la suite Karate del CRUD de Beneficios Egresados
# =============================================================================
# 1. Verifica/arranca PostgreSQL y re-siembra la BD (db/seed_pruebas.sql).
# 2. Compila y levanta el CRUD (:8080).
# 3. Ejecuta `mvn test` (reporte en target/karate-reports/).
# 4. Detiene el servicio.
#
# Prerrequisitos: Go, Java 11+, Maven, PostgreSQL 16 con la BD
# `beneficios_egresados` creada con db/schema.sql.
#
# Uso:  .\run_pruebas.ps1            (desde tests/karate)
#       .\run_pruebas.ps1 -NoReseed
# =============================================================================
param(
    [switch]$NoReseed,
    # Vacío = auto-detectar (PATH y luego C:\Program Files\PostgreSQL\<ver>\bin)
    [string]$PsqlPath   = '',
    [string]$DbUser     = $(if ($env:EGRESADOS_CRUD_PGUSER) { $env:EGRESADOS_CRUD_PGUSER } else { 'postgres' }),
    [string]$DbPassword = $(if ($env:EGRESADOS_CRUD_PGPASS) { $env:EGRESADOS_CRUD_PGPASS } else { '12345' }),
    # BD EXCLUSIVA de pruebas: la suite trunca/siembra datos, por eso NUNCA se
    # apunta a la BD de desarrollo (beneficios_egresados). Se crea sola si falta.
    [string]$DbName     = 'beneficios_egresados_pruebas'
)
$ErrorActionPreference = 'Stop'

$raizKarate = $PSScriptRoot
$raizCrud   = (Resolve-Path (Join-Path $raizKarate '..\..')).Path

if (-not $PsqlPath) {
    $cmd = Get-Command psql -ErrorAction SilentlyContinue
    if ($cmd) { $PsqlPath = $cmd.Source }
    else {
        $PsqlPath = Get-ChildItem 'C:\Program Files\PostgreSQL\*\bin\psql.exe' -ErrorAction SilentlyContinue |
            Sort-Object { $v = 0; [int]::TryParse($_.Directory.Parent.Name, [ref]$v) | Out-Null; $v } -Descending |
            Select-Object -First 1 -ExpandProperty FullName
    }
    if (-not $PsqlPath) {
        throw 'No se encontró psql.exe (ni en el PATH ni en C:\Program Files\PostgreSQL\<ver>\bin). Instala PostgreSQL 16+ o indica la ruta con -PsqlPath.'
    }
    Write-Host "Usando psql: $PsqlPath"
}

function Esperar-Puerto([int]$puerto, [string]$nombre) {
    foreach ($i in 1..60) {
        if ((Test-NetConnection 127.0.0.1 -Port $puerto -WarningAction SilentlyContinue).TcpTestSucceeded) { return }
        Start-Sleep -Milliseconds 500
    }
    throw "$nombre no respondió en el puerto $puerto tras 30s"
}

# ── 1. PostgreSQL (si el puerto ya responde, no se toca el servicio) ──────────
if (-not (Test-NetConnection 127.0.0.1 -Port 5432 -WarningAction SilentlyContinue).TcpTestSucceeded) {
    $svc = Get-Service | Where-Object { $_.Name -match 'postgres' } | Select-Object -First 1
    if ($svc) {
        Write-Host "Arrancando servicio $($svc.Name)..."
        try { Start-Service $svc.Name -ErrorAction Stop } catch { throw "No se pudo arrancar PostgreSQL (¿ejecutar como administrador, o arrancarlo con pg_ctl?): $_" }
    }
}
Esperar-Puerto 5432 'PostgreSQL'

$env:PGPASSWORD = $DbPassword
$existe = & $PsqlPath -U $DbUser -h 127.0.0.1 -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$DbName'"
if ($existe -ne '1') {
    Write-Host "Creando la BD de pruebas $DbName (schema.sql)..."
    & $PsqlPath -U $DbUser -h 127.0.0.1 -d postgres -c "CREATE DATABASE $DbName" | Out-Null
    if ($LASTEXITCODE -ne 0) { throw 'No se pudo crear la BD de pruebas' }
    & $PsqlPath -U $DbUser -h 127.0.0.1 -d $DbName -v ON_ERROR_STOP=1 -f (Join-Path $raizCrud 'db\schema.sql') | Out-Null
    if ($LASTEXITCODE -ne 0) { throw 'Falló la aplicación de db/schema.sql sobre la BD de pruebas' }
}

if (-not $NoReseed) {
    Write-Host "Re-sembrando $DbName con db/seed_pruebas.sql..."
    & $PsqlPath -U $DbUser -h 127.0.0.1 -d $DbName -v ON_ERROR_STOP=1 -f (Join-Path $raizCrud 'db\seed_pruebas.sql') | Out-Null
    if ($LASTEXITCODE -ne 0) { throw 'Falló la siembra de la BD de pruebas (¿credenciales?)' }
}

# ── 2. Compilar y levantar el CRUD ────────────────────────────────────────────
$binDir = Join-Path $raizKarate 'target\bin'
New-Item -ItemType Directory -Force $binDir | Out-Null

Write-Host 'Compilando el CRUD...'
Push-Location $raizCrud
go build -o (Join-Path $binDir 'crud_pruebas.exe') .
if ($LASTEXITCODE -ne 0) { Pop-Location; throw 'No compiló el CRUD' }
Pop-Location

$proc = $null
try {
    Write-Host "Levantando CRUD (:8080) contra la BD $DbName..."
    # Nombres de env estandarizados por la universidad (conf/app.conf, 2026-07-10):
    # EGRESADOS_CRUD_* — sin default quemado, hay que setear TODAS.
    $env:EGRESADOS_CRUD_HTTP_PORT = '8080'
    $env:EGRESADOS_CRUD_RUN_MODE  = 'dev'
    $env:EGRESADOS_CRUD_PGUSER    = $DbUser
    $env:EGRESADOS_CRUD_PGPASS    = $DbPassword
    $env:EGRESADOS_CRUD_PGHOST    = '127.0.0.1'
    $env:EGRESADOS_CRUD_PGPORT    = '5432'
    $env:EGRESADOS_CRUD_PGDB      = $DbName
    $env:EGRESADOS_CRUD_PGSCHEMA  = 'beneficios_egresados'
    $proc = Start-Process (Join-Path $binDir 'crud_pruebas.exe') -WorkingDirectory $raizCrud -PassThru -WindowStyle Hidden
    Esperar-Puerto 8080 'CRUD'

    # ── 3. Suite Karate ───────────────────────────────────────────────────────
    Write-Host 'Ejecutando la suite Karate (mvn test)...'
    Push-Location $raizKarate
    mvn --% -q test
    $exit = $LASTEXITCODE
    Pop-Location

    if ($exit -eq 0) { Write-Host "`n✔ Suite Karate del CRUD: TODO EN VERDE" -ForegroundColor Green }
    else { Write-Host "`n✘ Suite Karate del CRUD: hay fallos. Reporte: tests\karate\target\karate-reports\karate-summary.html" -ForegroundColor Red }
    exit $exit
}
finally {
    # ── 4. Apagar el servicio ─────────────────────────────────────────────────
    if ($proc -and -not $proc.HasExited) { Stop-Process -Id $proc.Id -Force -ErrorAction SilentlyContinue }
}
