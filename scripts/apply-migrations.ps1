# Jalankan migrasi SQL ke Postgres (jika API belum auto-migrate).
$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot\..

$files = @(
  "Backend\migrations\001_pixel_hub.up.sql",
  "Backend\migrations\002_cloudflare_setup.up.sql",
  "Backend\migrations\003_apidevel_domain.up.sql"
)

foreach ($f in $files) {
  if (-not (Test-Path $f)) { Write-Warning "Skip missing: $f"; continue }
  Write-Host "Applying $f ..." -ForegroundColor Cyan
  Get-Content $f -Raw | docker compose -f docker-compose.prod.yml exec -T db psql -U seosementara -d seosementara
}

Write-Host "Done. Restart API:" -ForegroundColor Green
Write-Host "  docker compose -f docker-compose.prod.yml restart api"
