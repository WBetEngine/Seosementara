# Pull image terbaru dari GHCR dan restart API (Opsi A — tanpa git pull).
$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot\..
docker compose -f docker-compose.prod.yml pull api
docker compose -f docker-compose.prod.yml up -d api
Invoke-WebRequest -Uri http://localhost:8080/health -UseBasicParsing | Select-Object -ExpandProperty Content
