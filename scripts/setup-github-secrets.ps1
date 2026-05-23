# Isi GitHub Secrets — jalankan di laptop (gh auth login dulu).
# winget install GitHub.cli  →  gh auth login
# Salin secrets.local.example → secrets.local, isi nilai, lalu jalankan script ini.

$ErrorActionPreference = "Stop"
$Repo = "WBetEngine/Seosementara"
$localFile = Join-Path $PSScriptRoot "secrets.local"

if (-not (Test-Path $localFile)) {
  Write-Host "Buat file: $localFile" -ForegroundColor Red
  Write-Host "Salin dari secrets.local.example dan isi nilainya."
  exit 1
}

Get-Content $localFile | ForEach-Object {
  $line = $_.Trim()
  if ($line -match '^\s*#' -or $line -eq "") { return }
  if ($line -match '^([^=]+)=(.*)$') {
    $name = $matches[1].Trim()
    $val = $matches[2].Trim().Trim('"')
    gh secret set $name -R $Repo -b $val
    Write-Host "OK $name" -ForegroundColor Green
  }
}

Write-Host "Selesai: https://github.com/$Repo/settings/secrets/actions" -ForegroundColor Cyan
