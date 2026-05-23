#Install cloudflared on Windows (mini PC). Run PowerShell as Administrator.
$ErrorActionPreference = "Stop"

$installDir = "C:\Program Files\cloudflared"
$exePath = Join-Path $installDir "cloudflared.exe"

if (-not (Test-Path $installDir)) {
  New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}

Write-Host "Downloading cloudflared..." -ForegroundColor Cyan
$zip = Join-Path $env:TEMP "cloudflared-windows-amd64.exe"
Invoke-WebRequest -Uri "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-windows-amd64.exe" -OutFile $zip -UseBasicParsing
Copy-Item $zip $exePath -Force

Write-Host "Installed: $exePath" -ForegroundColor Green
& $exePath --version

Write-Host @"

Next steps (token from Cloudflare Zero Trust → Tunnels → seosementara-api):
  & `"$exePath`" service install YOUR_TUNNEL_TOKEN
  Start-Service cloudflared
  curl.exe https://api.apidevel.org/health

See: Frontend-admin/SETUP-TUNNEL-APIDEVEL.md
"@ -ForegroundColor Yellow
