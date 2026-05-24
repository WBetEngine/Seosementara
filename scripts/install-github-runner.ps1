# Pasang GitHub self-hosted runner (sekali). Panduan: mini-pc/DEPLOY.md
$ErrorActionPreference = "Stop"
$runnerDir = "C:\actions-runner"
$repoUrl = "https://github.com/WBetEngine/Seosementara"
$ver = "2.334.0"

function Ensure-RunnerFiles {
  param([string]$Dir, [string]$Version)
  $need = @("config.cmd", "run.cmd")
  $missing = $need | Where-Object { -not (Test-Path (Join-Path $Dir $_)) }
  if ($missing.Count -eq 0) { return }

  Write-Host "Download runner (file hilang: $($missing -join ', '))..." -ForegroundColor Yellow
  $zip = Join-Path $Dir "actions-runner-win-x64-$Version.zip"
  Invoke-WebRequest -Uri "https://github.com/actions/runner/releases/download/v$Version/actions-runner-win-x64-$Version.zip" -OutFile $zip
  Expand-Archive $zip -DestinationPath $Dir -Force
  Remove-Item $zip -Force

  $nested = Get-ChildItem -Path $Dir -Directory | Where-Object {
    Test-Path (Join-Path $_.FullName "config.cmd")
  } | Select-Object -First 1
  if ($nested) {
    Write-Host "Pindahkan file dari subfolder $($nested.Name)..." -ForegroundColor Yellow
    Get-ChildItem -Path $nested.FullName -Force | Move-Item -Destination $Dir -Force
    Remove-Item $nested.FullName -Recurse -Force
  }

  $stillMissing = $need | Where-Object { -not (Test-Path (Join-Path $Dir $_)) }
  if ($stillMissing.Count -gt 0) {
    throw "Runner package tidak lengkap di $Dir - missing: $($stillMissing -join ', ')"
  }
}

Write-Host "Registration token: $repoUrl/settings/actions/runners/new" -ForegroundColor Yellow
Write-Host "Catatan: runner v2.334+ tidak punya install.cmd - pakai --runasservice saat config." -ForegroundColor Cyan
if (-not (Test-Path $runnerDir)) { New-Item -ItemType Directory -Path $runnerDir -Force | Out-Null }
Ensure-RunnerFiles -Dir $runnerDir -Version $ver

$regToken = Read-Host "Registration token"
Push-Location $runnerDir
try {
  & cmd /c "config.cmd --url $repoUrl --token $regToken --name mini-pc-seosementara --work _work --unattended --replace --runasservice"
  if (-not (Test-Path ".\install.cmd")) {
    Write-Host "Runner v2.334+: service dipasang lewat config --runasservice (tanpa install.cmd)." -ForegroundColor Cyan
  } else {
    & cmd /c install.cmd
  }
} finally {
  Pop-Location
}
Write-Host "Runner online. Cek: $repoUrl/settings/actions/runners" -ForegroundColor Green
