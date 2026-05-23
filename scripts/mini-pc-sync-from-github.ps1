# Unduh file runtime dari GitHub main (tanpa clone repo penuh).
param(
  [string]$Root = "C:\Seosementara",
  [string]$Repo = "WBetEngine/Seosementara",
  [string]$Branch = "main"
)

$ErrorActionPreference = "Stop"
$base = "https://raw.githubusercontent.com/$Repo/$Branch"

$files = @(
  @{ Rel = "docker-compose.prod.yml"; Local = "docker-compose.prod.yml" },
  @{ Rel = "scripts/mini-pc-pull.ps1"; Local = "scripts\mini-pc-pull.ps1" },
  @{ Rel = "scripts/apply-migrations.ps1"; Local = "scripts\apply-migrations.ps1" },
  @{ Rel = "scripts/bootstrap-cloudflare.ps1"; Local = "scripts\bootstrap-cloudflare.ps1" },
  @{ Rel = "scripts/mini-pc-remote-deploy.ps1"; Local = "scripts\mini-pc-remote-deploy.ps1" },
  @{ Rel = "Backend/migrations/001_pixel_hub.up.sql"; Local = "Backend\migrations\001_pixel_hub.up.sql" },
  @{ Rel = "Backend/migrations/002_cloudflare_setup.up.sql"; Local = "Backend\migrations\002_cloudflare_setup.up.sql" },
  @{ Rel = "Backend/migrations/003_apidevel_domain.up.sql"; Local = "Backend\migrations\003_apidevel_domain.up.sql" }
)

foreach ($f in $files) {
  $dest = Join-Path $Root $f.Local
  $dir = Split-Path $dest -Parent
  if (-not (Test-Path $dir)) { New-Item -ItemType Directory -Path $dir -Force | Out-Null }
  Write-Host "GET $($f.Rel)" -ForegroundColor Cyan
  Invoke-WebRequest -Uri "$base/$($f.Rel)" -OutFile $dest -UseBasicParsing
}

Write-Host "OK — file dari GitHub main sudah di $Root" -ForegroundColor Green
Write-Host "Catatan: .env TIDAK di GitHub — diisi via GitHub Secret MINI_PC_DOTENV (Actions) atau manual."
