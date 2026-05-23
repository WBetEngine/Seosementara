# Muat file .env ke environment process (PowerShell tidak baca .env otomatis).
function Import-DotEnv {
  param(
    [Parameter(Mandatory = $true)]
    [string]$Path
  )

  if (-not (Test-Path $Path)) {
    throw ".env tidak ditemukan: $Path`nSalin mini-pc/env.example ke .env dan isi nilainya."
  }

  Get-Content $Path | ForEach-Object {
    $line = $_.Trim()
    if ($line -match '^\s*#' -or $line -eq "") { return }
    if ($line -match '^([^=]+)=(.*)$') {
      $name = $matches[1].Trim()
      $val = $matches[2].Trim().Trim('"').Trim("'")
      [Environment]::SetEnvironmentVariable($name, $val, "Process")
    }
  }
}
