# Aktifkan OpenSSH Server di Windows (sekali, PowerShell Administrator).
# Diperlukan agar MCP / agent bisa SSH ke mini PC via Tailscale.
$ErrorActionPreference = "Stop"

$cap = Get-WindowsCapability -Online | Where-Object { $_.Name -like "OpenSSH.Server*" }
if ($cap.State -ne "Installed") {
  Write-Host "Installing OpenSSH Server..." -ForegroundColor Cyan
  Add-WindowsCapability -Online -Name $cap.Name
}

Start-Service sshd -ErrorAction SilentlyContinue
Set-Service -Name sshd -StartupType Automatic
New-NetFirewallRule -Name "OpenSSH-Server-In-TCP" -DisplayName "OpenSSH Server (sshd)" -Enabled True -Direction Inbound -Protocol TCP -Action Allow -LocalPort 22 -ErrorAction SilentlyContinue

Write-Host "OpenSSH Server: Running" -ForegroundColor Green
Get-Service sshd | Format-List Name, Status, StartType
Write-Host "Test dari laptop (Tailscale): ssh seosementara@100.100.17.92 hostname" -ForegroundColor Yellow
