# Jalankan di mini PC (sekali) — PowerShell Administrator.
# Setelah ini: push ke GitHub main = deploy otomatis, tanpa SSH dari internet.
$ErrorActionPreference = "Stop"

$runnerDir = "C:\actions-runner"
$repoUrl = "https://github.com/WBetEngine/Seosementara"
$token = Read-Host "GitHub PAT (scope: repo + admin:org read for private, atau repo saja)"

Write-Host @"

=== GitHub Self-Hosted Runner (sekali saja) ===
1. Buka: $repoUrl/settings/actions/runners/new
2. Pilih Windows x64
3. Salin token registration (bukan PAT) — atau pakai PAT di bawah

"@ -ForegroundColor Yellow

if (-not (Test-Path $runnerDir)) {
  New-Item -ItemType Directory -Path $runnerDir -Force | Out-Null
}
Set-Location $runnerDir

if (-not (Test-Path ".\config.cmd")) {
  Write-Host "Download runner..." -ForegroundColor Cyan
  Invoke-WebRequest -Uri "https://github.com/actions/runner/releases/download/v2.321.0/actions-runner-win-x64-2.321.0.zip" -OutFile runner.zip
  Expand-Archive runner.zip -DestinationPath . -Force
  Remove-Item runner.zip
}

$regToken = Read-Host "Registration token dari GitHub (Settings → Actions → Runners → New)"

.\config.cmd --url $repoUrl --token $regToken --name "mini-pc-seosementara" --work _work --unattended --replace

Write-Host "Install as Windows service..." -ForegroundColor Cyan
.\install.cmd

Write-Host @"
Selesai. Runner muncul di GitHub → Settings → Actions → Runners (online).
Workflow 'Deploy Mini PC (local)' akan jalan di PC ini otomatis.
"@ -ForegroundColor Green
