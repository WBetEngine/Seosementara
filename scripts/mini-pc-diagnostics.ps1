# Mini PC diagnostics — output JSON untuk GitHub Actions / MCP
$ErrorActionPreference = "Continue"
$report = [ordered]@{
  timestamp = (Get-Date).ToUniversalTime().ToString("o")
  hostname  = $env:COMPUTERNAME
  checks    = @()
}

function Add-Check($name, $ok, $detail) {
  $report.checks += @{ name = $name; ok = [bool]$ok; detail = "$detail" }
}

# API health
try {
  $local = curl.exe -sf http://localhost:8080/health 2>$null
  Add-Check "api_localhost" ($local -eq "ok") $local
} catch {
  Add-Check "api_localhost" $false $_.Exception.Message
}

try {
  $remote = curl.exe -sf https://api.apidevel.org/health 2>$null
  Add-Check "api_tunnel" ($remote -eq "ok") $remote
} catch {
  Add-Check "api_tunnel" $false $_.Exception.Message
}

# Docker
try {
  Set-Location C:\Seosementara -ErrorAction Stop
  $psOut = docker compose -f docker-compose.prod.yml ps --format json 2>&1 | Out-String
  Add-Check "docker_compose" ($LASTEXITCODE -eq 0) $psOut.Trim()
} catch {
  Add-Check "docker_compose" $false $_.Exception.Message
}

# Runner service
$runnerSvc = Get-Service | Where-Object { $_.Name -like "actions.runner.*" } | Select-Object -First 1
if ($runnerSvc) {
  Add-Check "runner_service" ($runnerSvc.Status -eq "Running") "$($runnerSvc.Name)=$($runnerSvc.Status)"
} else {
  Add-Check "runner_service" $false "Tidak ada service actions.runner.*"
}

# cloudflared
$cf = Get-Service cloudflared -ErrorAction SilentlyContinue
if ($cf) {
  Add-Check "cloudflared" ($cf.Status -eq "Running") "$($cf.Name)=$($cf.Status)"
} else {
  Add-Check "cloudflared" $false "Service cloudflared tidak ditemukan"
}

$report.ok = -not ($report.checks | Where-Object { -not $_.ok })
$report | ConvertTo-Json -Depth 5
