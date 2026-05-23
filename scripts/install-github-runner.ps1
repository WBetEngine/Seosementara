# Pasang GitHub self-hosted runner (Windows) — alternatif otomatis.
# Panduan lengkap manual: mini-pc/DEPLOY.md
$ErrorActionPreference = "Stop"

$runnerDir = "C:\actions-runner"
$repoUrl = "https://github.com/WBetEngine/Seosementara"
$runnerVersion = "2.334.0"
$runnerZip = "actions-runner-win-x64-$runnerVersion.zip"
$runnerUrl = "https://github.com/actions/runner/releases/download/v$runnerVersion/$runnerZip"

Write-Host @"

=== GitHub Self-Hosted Runner ===
1. Buka: $repoUrl/settings/actions/runners/new
2. Pilih Windows x64
3. Salin registration token (kadaluarsa ~1 jam)

"@ -ForegroundColor Yellow

if (-not (Test-Path $runnerDir)) {
  New-Item -ItemType Directory -Path $runnerDir -Force | Out-Null
}
Set-Location $runnerDir

if (-not (Test-Path ".\config.cmd")) {
  Write-Host "Download runner v$runnerVersion..." -ForegroundColor Cyan
  Invoke-WebRequest -Uri $runnerUrl -OutFile $runnerZip
  Expand-Archive $runnerZip -DestinationPath . -Force
  Remove-Item $runnerZip
}

$regToken = Read-Host "Registration token dari GitHub"

.\config.cmd --url $repoUrl --token $regToken --name "mini-pc-seosementara" --work _work --unattended --replace

Write-Host "Install as Windows service..." -ForegroundColor Cyan
.\install.cmd

Write-Host @"
Selesai. Runner online di GitHub → Settings → Actions → Runners.
Pastikan C:\Seosementara\.env sudah dibuat sebelum Deploy Mini PC.
"@ -ForegroundColor Green
