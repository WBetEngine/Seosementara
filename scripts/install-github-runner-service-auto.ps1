# Pasang runner sebagai Windows Service. Token registrasi otomatis via GitHub API (PAT).
# PowerShell Administrator. Tutup run.cmd dulu (Ctrl+C).
param(
  [string]$GitHubPat = $env:PLATFORM_GITHUB_PAT,
  [string]$RunnerDir = "C:\actions-runner",
  [string]$RepoOwner = "WBetEngine",
  [string]$RepoName = "Seosementara",
  [string]$RunnerName = "mini-pc-seosementara"
)

$ErrorActionPreference = "Stop"
$repoUrl = "https://github.com/$RepoOwner/$RepoName"

function Test-Admin {
  $id = [Security.Principal.WindowsIdentity]::GetCurrent()
  $p = New-Object Security.Principal.WindowsPrincipal($id)
  return $p.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

function Get-RunnerRegistrationToken {
  param([string]$Pat)
  if (-not $Pat) { throw "PAT kosong - set PLATFORM_GITHUB_PAT atau -GitHubPat" }
  $headers = @{
    Authorization = "Bearer $Pat"
    Accept = "application/vnd.github+json"
    "X-GitHub-Api-Version" = "2022-11-28"
    "User-Agent" = "Seosementara-Platform/1.0"
  }
  $uri = "https://api.github.com/repos/$RepoOwner/$RepoName/actions/runners/registration-token"
  try {
    $resp = Invoke-RestMethod -Method Post -Uri $uri -Headers $headers -Body "{}"
  } catch {
    throw "GitHub API registration-token gagal: $($_.Exception.Message) - PAT perlu Administration write"
  }
  if (-not $resp.token) { throw "GitHub API tidak mengembalikan token registrasi" }
  return $resp.token
}

if (-not (Test-Admin)) {
  throw "Jalankan PowerShell as Administrator"
}

if (-not (Test-Path (Join-Path $RunnerDir "config.cmd"))) {
  throw "Runner belum ada di $RunnerDir"
}

Write-Host "=== Pasang GitHub Runner sebagai Windows Service (auto token) ===" -ForegroundColor Cyan

if (-not $GitHubPat) {
  $GitHubPat = Read-Host "GitHub PAT (fine-grained: Administration write)"
}
$regToken = Get-RunnerRegistrationToken -Pat $GitHubPat
Write-Host "Registration token didapat via GitHub API." -ForegroundColor Green

Push-Location $RunnerDir
try {
  & cmd /c "config.cmd --url $repoUrl --token $regToken --name $RunnerName --work _work --unattended --replace --runasservice"
} finally {
  Pop-Location
}

Start-Sleep -Seconds 3
$svc = Get-Service | Where-Object { $_.Name -like "actions.runner.*" } | Select-Object -First 1
if ($svc) {
  if ($svc.Status -ne "Running") { Start-Service $svc.Name }
  Write-Host "Service: $($svc.Name) - $($svc.Status)" -ForegroundColor Green
} else {
  Write-Host "Service tidak terdeteksi - cek services.msc" -ForegroundColor Yellow
}

Write-Host "Runner: $repoUrl/settings/actions/runners" -ForegroundColor Green
