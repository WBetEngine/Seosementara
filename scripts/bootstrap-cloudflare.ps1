# Bootstrap Cloudflare via API — baca dari file .env di C:\Seosementara.
param(
  [string]$ApiBase = "https://api.apidevel.org"
)

$ErrorActionPreference = "Stop"
$Root = if ($PSScriptRoot) { Split-Path $PSScriptRoot -Parent } else { "C:\Seosementara" }

. "$PSScriptRoot\load-dotenv.ps1"
Import-DotEnv -Path (Join-Path $Root ".env")

function Require-Env([string]$Name) {
  $v = [Environment]::GetEnvironmentVariable($Name)
  if (-not $v) { throw "Missing di .env: $Name" }
  return $v
}

$token = Require-Env "SUPER_ADMIN_TOKEN"
$cfKey = Require-Env "CLOUDFLARE_API_KEY"
$cfEmail = Require-Env "CLOUDFLARE_ACCOUNT_EMAIL"
$cfAccount = Require-Env "CLOUDFLARE_ACCOUNT_ID"

$headers = @{
  Authorization = "Bearer $token"
  "Content-Type" = "application/json"
}

Write-Host "1. Save Cloudflare credentials..." -ForegroundColor Cyan
$body = @{
  auth_type      = "global_api_key"
  global_api_key = $cfKey
  account_email  = $cfEmail
  account_id     = $cfAccount
} | ConvertTo-Json

Invoke-RestMethod -Uri "$ApiBase/api/admin/setup/cloudflare/credentials?test=1" -Method Put -Headers $headers -Body $body | Out-Null

Write-Host "2. Update domain env..." -ForegroundColor Cyan
$envBody = @{
  PRIMARY_DOMAIN = "apidevel.org"
  APEX_URL       = "https://apidevel.org"
  API_BASE_URL   = "https://api.apidevel.org"
  ENVIRONMENT    = "production"
} | ConvertTo-Json

Invoke-RestMethod -Uri "$ApiBase/api/admin/setup/cloudflare/domain-env" -Method Put -Headers $headers -Body $envBody | Out-Null

Write-Host "3. Refresh tunnel status..." -ForegroundColor Cyan
Invoke-RestMethod -Uri "$ApiBase/api/admin/setup/cloudflare/tunnel/status" -Method Post -Headers $headers | Out-Null

Write-Host "Cloudflare bootstrap OK." -ForegroundColor Green
