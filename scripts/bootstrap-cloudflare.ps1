# Bootstrap Cloudflare credentials + domain env via API (mini PC / production).
# Baca secret dari C:\Seosementara\.env — jangan commit file ini dengan secret hardcoded.
param(
  [string]$ApiBase = "https://api.apidevel.org",
  [string]$EnvFile = ".env"
)

$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot\..

function Get-EnvValue([string]$Name) {
  $line = Select-String -Path $EnvFile -Pattern "^$Name=" | Select-Object -First 1
  if (-not $line) { throw "Missing $Name in $EnvFile" }
  return ($line.Line -replace "^$Name=", "").Trim('"')
}

$token = Get-EnvValue "SUPER_ADMIN_TOKEN"
$cfKey = Get-EnvValue "CLOUDFLARE_API_KEY"
$cfEmail = Get-EnvValue "CLOUDFLARE_ACCOUNT_EMAIL"
$cfAccount = Get-EnvValue "CLOUDFLARE_ACCOUNT_ID"

$headers = @{
  Authorization = "Bearer $token"
  "Content-Type" = "application/json"
}

Write-Host "1. Save Cloudflare credentials..." -ForegroundColor Cyan
$body = @{
  auth_type     = "global_api_key"
  global_api_key = $cfKey
  account_email = $cfEmail
  account_id    = $cfAccount
} | ConvertTo-Json

$r = Invoke-RestMethod -Uri "$ApiBase/api/admin/setup/cloudflare/credentials?test=1" -Method Put -Headers $headers -Body $body
$r | ConvertTo-Json

Write-Host "2. Update domain env..." -ForegroundColor Cyan
$envBody = @{
  PRIMARY_DOMAIN = "apidevel.org"
  APEX_URL       = "https://apidevel.org"
  API_BASE_URL   = "https://api.apidevel.org"
  ENVIRONMENT    = "production"
} | ConvertTo-Json

Invoke-RestMethod -Uri "$ApiBase/api/admin/setup/cloudflare/domain-env" -Method Put -Headers $headers -Body $envBody | Out-Null

Write-Host "3. Refresh tunnel status..." -ForegroundColor Cyan
Invoke-RestMethod -Uri "$ApiBase/api/admin/setup/cloudflare/tunnel/status" -Method Post -Headers $headers | ConvertTo-Json

Write-Host "OK — buka admin Settings → Cloudflare" -ForegroundColor Green
