# GitHub Actions → Mini PC: setup Secrets (sekali)

GitHub **tidak** menyimpan `.env` di repo. Isi disimpan sebagai **Secret** `MINI_PC_DOTENV`, lalu Actions menulis ke `C:\Seosementara\.env` via SSH/SCP.

## 1. Aktifkan OpenSSH di Windows (mini PC)

PowerShell **Administrator**:

```powershell
Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0
Start-Service sshd
Set-Service -Name sshd -StartupType Automatic
New-NetFirewallRule -Name "OpenSSH" -DisplayName "OpenSSH Server (sshd)" -Enabled True -Direction Inbound -Protocol TCP -Action Allow -LocalPort 22
```

Pastikan login password aktif (`PasswordAuthentication yes` di `C:\ProgramData\ssh\sshd_config`), lalu:

```powershell
Restart-Service sshd
```

## 2. Secrets di GitHub repo

**Settings → Secrets and variables → Actions → New repository secret**

| Secret | Nilai (contoh Anda) |
|--------|---------------------|
| `DEPLOY_HOST` | `100.100.17.92` |
| `DEPLOY_USER` | `seosementara` |
| `DEPLOY_SSH_PASSWORD` | password Windows user |
| `DEPLOY_PATH` | `C:/Seosementara` |
| `DEPLOY_SSH_PORT` | `22` (opsional) |
| `MINI_PC_DOTENV` | **seluruh isi** file `.env` (lihat bawah) |

Opsional (lebih aman dari password): `DEPLOY_SSH_KEY` = private key OpenSSH.

### Isi `MINI_PC_DOTENV` (copy-paste ke Secret)

```env
DB_PASSWORD=...
DATABASE_URL=postgres://seosementara:...@localhost:5432/seosementara?sslmode=disable
MASTER_ENCRYPTION_KEY=...
SUPER_ADMIN_TOKEN=...
CLOUDFLARE_API_KEY=...
CLOUDFLARE_ACCOUNT_ID=...
CLOUDFLARE_ACCOUNT_EMAIL=...
CLOUDFLARE_TUNNEL_ID=...
```

**Penting:** `SUPER_ADMIN_TOKEN` harus sama dengan `admin-config.js` di Workers.

## 3. Variable (opsional)

| Variable | Nilai |
|----------|--------|
| `MINI_PC_DEPLOY` | `true` — agar deploy image + sync setelah build |

## 4. Jalankan sync

**Actions → Sync Mini PC → Run workflow**

Atau otomatis setelah push ke `main` (migrasi / compose berubah) atau setelah **Deploy Backend API** sukses.

## 5. Manual tanpa Actions

```powershell
cd C:\Seosementara
.\scripts\mini-pc-sync-from-github.ps1
# .env tetap manual atau dari Secret yang Anda salin
docker compose -f docker-compose.prod.yml up -d
```

## Alur

```text
GitHub Secret MINI_PC_DOTENV
    → Actions SCP → C:\Seosementara\.env
GitHub main (migrations, compose, scripts)
    → Actions SCP → C:\Seosementara\
    → SSH → mini-pc-remote-deploy.ps1
    → docker compose pull + up
```

**Bukan MCP** — ini **GitHub Actions + SSH/SCP**, standar untuk deploy ke server rumah (Opsi A).
