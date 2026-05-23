# Pasang GitHub self-hosted runner (sekali). Panduan: mini-pc/DEPLOY.md
$ErrorActionPreference = "Stop"
$runnerDir = "C:\actions-runner"
$repoUrl = "https://github.com/WBetEngine/Seosementara"
$ver = "2.334.0"

Write-Host "Registration token: $repoUrl/settings/actions/runners/new" -ForegroundColor Yellow
if (-not (Test-Path $runnerDir)) { New-Item -ItemType Directory -Path $runnerDir -Force | Out-Null }
Set-Location $runnerDir
if (-not (Test-Path ".\config.cmd")) {
  Invoke-WebRequest -Uri "https://github.com/actions/runner/releases/download/v$ver/actions-runner-win-x64-$ver.zip" -OutFile runner.zip
  Expand-Archive runner.zip -DestinationPath . -Force
  Remove-Item runner.zip
}
$regToken = Read-Host "Registration token"
.\config.cmd --url $repoUrl --token $regToken --name "mini-pc-seosementara" --work _work --unattended --replace
.\install.cmd
Write-Host "Runner online. Lanjut setup lewat admin Workers URL." -ForegroundColor Green
