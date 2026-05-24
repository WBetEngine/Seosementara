# Pasang runner sebagai Windows Service (--runasservice). Panduan: mini-pc/DEPLOY.md
# Jalankan PowerShell Administrator. Tutup dulu jendela run.cmd jika masih terbuka.
$ErrorActionPreference = "Stop"
$runnerDir = "C:\actions-runner"
$repoUrl = "https://github.com/WBetEngine/Seosementara"
$runnerName = "mini-pc-seosementara"

function Test-Admin {
  $id = [Security.Principal.WindowsIdentity]::GetCurrent()
  $p = New-Object Security.Principal.WindowsPrincipal($id)
  return $p.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

if (-not (Test-Admin)) {
  throw "Jalankan PowerShell as Administrator (klik kanan -> Run as administrator)"
}

if (-not (Test-Path (Join-Path $runnerDir "config.cmd"))) {
  throw "Runner belum ada di $runnerDir - jalankan install-github-runner.ps1 dulu"
}

Write-Host ""
Write-Host "=== Pasang GitHub Runner sebagai Windows Service ===" -ForegroundColor Cyan
Write-Host "1. Tutup jendela run.cmd (Ctrl+C) jika masih jalan" -ForegroundColor Yellow
Write-Host "2. Token baru (sekali pakai):" -ForegroundColor Yellow
Write-Host "   $repoUrl/settings/actions/runners/new" -ForegroundColor White
Write-Host ""

$regToken = Read-Host "Registration token"
if (-not $regToken) { throw "Token wajib diisi" }

Push-Location $runnerDir
try {
  Write-Host "Config ulang dengan --runasservice ..." -ForegroundColor Cyan
  & cmd /c "config.cmd --url $repoUrl --token $regToken --name $runnerName --work _work --unattended --replace --runasservice"
} finally {
  Pop-Location
}

Start-Sleep -Seconds 3
$svc = Get-Service | Where-Object { $_.Name -like "actions.runner.*" } | Select-Object -First 1
if ($svc) {
  if ($svc.Status -ne "Running") {
    Write-Host "Menyalakan service $($svc.Name) ..." -ForegroundColor Yellow
    Start-Service $svc.Name
  }
  Write-Host "Service: $($svc.Name) - $($svc.Status)" -ForegroundColor Green
} else {
  Write-Host "Service tidak terdeteksi - buka services.msc, cari GitHub Actions Runner" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Cek runner Idle di GitHub:" -ForegroundColor Green
Write-Host "$repoUrl/settings/actions/runners" -ForegroundColor White
Write-Host "Window run.cmd tidak perlu dibuka lagi." -ForegroundColor Green
