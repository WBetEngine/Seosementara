# Deploy Docker di mini PC — env dari GitHub Secrets (process), bukan file .env.
$ErrorActionPreference = "Stop"
$Root = if ($PSScriptRoot) { Split-Path $PSScriptRoot -Parent } else { "C:\Seosementara" }
Set-Location $Root

foreach ($var in @("DB_PASSWORD", "MASTER_ENCRYPTION_KEY")) {
  if (-not [Environment]::GetEnvironmentVariable($var)) {
    throw "$var belum diset — isi lewat admin Settings → Infra & GitHub"
  }
}

Write-Host "=== Seosementara deploy (GitHub Secrets inject) ===" -ForegroundColor Cyan
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d --force-recreate
Start-Sleep -Seconds 8
docker compose -f docker-compose.prod.yml ps
$health = curl.exe -sf http://localhost:8080/health 2>$null
Write-Host "health: $health" -ForegroundColor Green
