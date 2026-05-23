# Deploy ke mini PC — dipanggil GitHub Actions self-hosted runner.
# Tidak membutuhkan file .env: secret sudah di environment process.
$ErrorActionPreference = "Stop"
$Root = if ($PSScriptRoot) { Split-Path $PSScriptRoot -Parent } else { "C:\Seosementara" }
Set-Location $Root

Write-Host "=== Seosementara deploy ===" -ForegroundColor Cyan

foreach ($var in @("DB_PASSWORD", "MASTER_ENCRYPTION_KEY", "SUPER_ADMIN_TOKEN")) {
  if (-not [Environment]::GetEnvironmentVariable($var)) {
    throw "Environment variable $var belum diset (GitHub Secrets → Actions)"
  }
}

if (-not (Test-Path "$Root\docker-compose.prod.yml")) {
  throw "docker-compose.prod.yml tidak ditemukan"
}

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
$health = curl.exe -sf http://localhost:8080/health 2>$null
Write-Host "localhost:8080/health => $health"

if ($env:CLOUDFLARE_API_KEY -and (Test-Path "$Root\scripts\bootstrap-cloudflare.ps1")) {
  Write-Host "Bootstrap Cloudflare..." -ForegroundColor Cyan
  & "$Root\scripts\bootstrap-cloudflare.ps1" 2>&1 | Out-Host
}

Write-Host "Deploy selesai." -ForegroundColor Green
