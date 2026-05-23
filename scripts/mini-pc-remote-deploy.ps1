# Dijalankan oleh GitHub Actions (SSH) atau manual di mini PC.
$ErrorActionPreference = "Stop"
$Root = if ($PSScriptRoot) { Split-Path $PSScriptRoot -Parent } else { "C:\Seosementara" }
Set-Location $Root

Write-Host "=== Seosementara remote deploy ===" -ForegroundColor Cyan
Write-Host "Root: $Root"

if (-not (Test-Path "$Root\.env")) {
  throw ".env tidak ditemukan di $Root"
}
if (-not (Test-Path "$Root\docker-compose.prod.yml")) {
  throw "docker-compose.prod.yml tidak ditemukan"
}

# Pastikan folder migrasi ada
$migDir = Join-Path $Root "Backend\migrations"
if (-not (Test-Path $migDir)) {
  New-Item -ItemType Directory -Path $migDir -Force | Out-Null
}

Write-Host "Pull image GHCR..." -ForegroundColor Cyan
docker compose -f docker-compose.prod.yml pull

Write-Host "Start stack..." -ForegroundColor Cyan
docker compose -f docker-compose.prod.yml up -d --force-recreate

Start-Sleep -Seconds 8

Write-Host "Health check..." -ForegroundColor Cyan
docker compose -f docker-compose.prod.yml ps
$local = curl.exe -sf http://localhost:8080/health 2>$null
Write-Host "localhost:8080/health => $local"

if (Test-Path "$Root\scripts\bootstrap-cloudflare.ps1") {
  $hasCf = Select-String -Path "$Root\.env" -Pattern "^CLOUDFLARE_API_KEY=.+" -Quiet
  if ($hasCf) {
    Write-Host "Bootstrap Cloudflare..." -ForegroundColor Cyan
    & "$Root\scripts\bootstrap-cloudflare.ps1" -EnvFile "$Root\.env" 2>&1 | Out-Host
  }
}

Write-Host "Deploy selesai." -ForegroundColor Green
