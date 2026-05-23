# Deploy mini PC

Setup **sekali** di mini PC. Setelah itu, push ke GitHub = deploy otomatis.

## Ringkasan

| Apa | Di mana |
|-----|---------|
| Secret (DB, token, Cloudflare) | File **`C:\Seosementara\.env`** (manual, sekali) |
| Kode & migrasi | GitHub → sync otomatis via runner |
| Image API | GHCR (pull otomatis) |
| Admin UI | Cloudflare Workers (GitHub Actions terpisah) |

```text
push main → Build image (GHCR) → Deploy Mini PC (runner)
                                      ↓
                              docker compose + .env
```

---

## Langkah 1 — Buat folder & file `.env`

PowerShell di mini PC:

```powershell
mkdir C:\Seosementara -Force
mkdir C:\Seosementara\scripts -Force
mkdir C:\Seosementara\Backend\migrations -Force
notepad C:\Seosementara\.env
```

Isi `.env` (contoh — ganti dengan nilai Anda):

```env
DB_PASSWORD=password-db-anda
MASTER_ENCRYPTION_KEY=base64-32-byte
SUPER_ADMIN_TOKEN=token-panjang-random

CLOUDFLARE_API_KEY=global-api-key-cloudflare
CLOUDFLARE_ACCOUNT_ID=account-id
CLOUDFLARE_ACCOUNT_EMAIL=email-cloudflare
CLOUDFLARE_TUNNEL_ID=opsional
```

Template lengkap ada di repo: `mini-pc/env.example`

> **Penting:** `.env` tidak pernah di-commit ke GitHub. Hanya ada di mini PC.

---

## Langkah 2 — Pasang GitHub runner (sekali)

PowerShell **Administrator**:

### 2a. Download runner

```powershell
mkdir C:\actions-runner
cd C:\actions-runner

Invoke-WebRequest -Uri https://github.com/actions/runner/releases/download/v2.334.0/actions-runner-win-x64-2.334.0.zip -OutFile runner.zip
Expand-Archive runner.zip -DestinationPath . -Force
Remove-Item runner.zip
```

### 2b. Ambil registration token

Buka: https://github.com/WBetEngine/Seosementara/settings/actions/runners/new

Pilih **Windows x64**, salin token dari baris `config.cmd` (kadaluarsa ~1 jam).

### 2c. Configure & install service

Ganti `TOKEN_DARI_GITHUB` dengan token tadi:

```powershell
cd C:\actions-runner
.\config.cmd --url https://github.com/WBetEngine/Seosementara --token TOKEN_DARI_GITHUB --name mini-pc-seosementara --work _work --unattended --replace
.\install.cmd
```

### 2d. Cek runner online

```powershell
Get-Service actions.runner.*
```

Status **Running**. Di GitHub → **Settings → Actions → Runners** → `mini-pc-seosementara` **Idle** (hijau).

---

## Langkah 3 — Deploy pertama

GitHub → **Actions** → **Deploy Mini PC** → **Run workflow**

Runner akan:
1. Sync `docker-compose.prod.yml`, migrasi, dan script ke `C:\Seosementara`
2. Baca secret dari **`.env`** Anda
3. `docker compose pull` + `up -d`

Cek health:

```powershell
curl http://localhost:8080/health
```

---

## Setelah setup

Edit kode di GitHub → push `main` → deploy otomatis (tanpa sentuh mini PC).

---

## File di mini PC

```
C:\Seosementara\
  .env                          ← Anda buat manual (secret)
  docker-compose.prod.yml       ← sync otomatis
  env.example                   ← referensi (sync otomatis)
  Backend\migrations\           ← sync otomatis
  scripts\
    load-dotenv.ps1
    mini-pc-remote-deploy.ps1
    bootstrap-cloudflare.ps1
```

---

## GitHub Secrets (hanya untuk Admin UI)

Secret di GitHub **tidak** dipakai deploy mini PC lagi. Masih dipakai workflow **Deploy Admin UI**:

| Secret | Fungsi |
|--------|--------|
| `SUPER_ADMIN_TOKEN` | Harus **sama** dengan di `.env` mini PC |
| `CLOUDFLARE_API_KEY` | Deploy Workers |
| `CLOUDFLARE_ACCOUNT_ID` | Deploy Workers |
| `CLOUDFLARE_ACCOUNT_EMAIL` | Deploy Workers |

---

## Troubleshooting

| Masalah | Solusi |
|---------|--------|
| `.env belum ada` | Buat `C:\Seosementara\.env` (Langkah 1) |
| Runner offline | `Start-Service actions.runner.*` |
| 401 di admin | Samakan `SUPER_ADMIN_TOKEN` di `.env` dan GitHub Secret |
| Docker error | Buka Docker Desktop |

---

## Prasyarat (sudah)

- Docker Desktop
- cloudflared tunnel → `api.apidevel.org`
