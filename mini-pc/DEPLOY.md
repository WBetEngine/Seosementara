# Mini PC — Docker saja (tanpa kode / tanpa Git)

Mini PC **bukan** tempat kode. Hanya menjalankan container.

## Yang boleh ada di mini PC

| Komponen | Fungsi |
|----------|--------|
| Docker Desktop | Postgres + API (image dari GHCR) |
| cloudflared | Tunnel → `api.apidevel.org` |
| GitHub runner | Terima deploy dari GitHub Actions |
| `C:\Seosementara\` | **Runtime saja** (bukan repo Git) — disync otomatis oleh CI |

```text
C:\Seosementara\          ← disalin otomatis workflow Deploy Mini PC
  docker-compose.prod.yml
  Backend\migrations\*.sql
  scripts\mini-pc-deploy.ps1

Tidak ada: git, Frontend-admin, Backend source, .env
```

Semua kode & admin UI ada di **GitHub** + **Cloudflare Workers**.

---

## Setup operator (100% lewat admin Workers)

Buka: **https://seosementara.seosementara3.workers.dev/admin/settings/backend/infra**

### 1. Bootstrap Platform

GitHub PAT + Global API Key + email + Account ID → GitHub Environment `production` + Workers Secrets.

### 2. Infra mini PC

DB_PASSWORD + MASTER_ENCRYPTION_KEY → GitHub Environment → runner inject Docker.

### 3. Cloudflare (opsional)

Settings → Cloudflare → Koneksi.

**Tidak perlu** buka GitHub Settings manual.

---

## Setup mini PC (sekali, tanpa Git)

### A. Pasang runner

PowerShell **Administrator** — download script saja (tanpa clone repo):

```powershell
mkdir C:\Seosementara\scripts -Force
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/WBetEngine/Seosementara/main/scripts/install-github-runner.ps1" -OutFile C:\Seosementara\scripts\install-github-runner.ps1
C:\Seosementara\scripts\install-github-runner.ps1
```

Token runner: https://github.com/WBetEngine/Seosementara/settings/actions/runners/new

### B. Docker + cloudflared

Sudah terpasang (prasyarat).

### C. Deploy pertama

Setelah Bootstrap + Infra di admin → workflow **Deploy Mini PC** otomatis:

1. Sync file runtime ke `C:\Seosementara`
2. `docker compose pull` + `up` (secret dari GitHub, bukan `.env`)

Manual: GitHub → **Actions → Deploy Mini PC → Run workflow**

---

## Update rutin

Anda **hanya** edit & push di GitHub. Mini PC tidak disentuh.

```text
push main → GHCR image → Deploy Mini PC (runner) → Docker restart
push Frontend-admin → Deploy Admin UI → Cloudflare Workers
```

---

## Verifikasi

```powershell
curl.exe http://localhost:8080/health
```

GitHub → **Actions** → Deploy Mini PC / Deploy Admin UI = Success.

Panduan arsitektur: [Plan/28-platform-github-workers.md](../Plan/28-platform-github-workers.md)
