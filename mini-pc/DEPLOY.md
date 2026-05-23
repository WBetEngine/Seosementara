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

## Bootstrap (admin panel)

Buka: **https://seosementara.seosementara3.workers.dev/admin/settings/backend/infra**

Bootstrap admin akan **membuat** Environment `production` otomatis via API (tidak perlu buka URL environment manual).

### Deploy Workers pertama kali (jika admin belum update)

Isi **Repository Secrets** (bukan Environment — link ini selalu ada):

**https://github.com/WBetEngine/Seosementara/settings/secrets/actions**

| Secret | Wajib untuk deploy Workers |
|--------|----------------------------|
| `CLOUDFLARE_API_KEY` | Ya |
| `CLOUDFLARE_ACCOUNT_EMAIL` | Ya |
| `CLOUDFLARE_ACCOUNT_ID` | Ya |
| `SUPER_ADMIN_TOKEN` | Opsional |

Lalu: **Actions → Deploy Admin UI → Run workflow**

### Environment `production` (setelah Bootstrap)

URL hanya ada **setelah** environment dibuat:

1. List: https://github.com/WBetEngine/Seosementara/settings/environments  
2. Atau otomatis saat **Bootstrap Platform** di admin sukses  

Jika `/settings/environments/production` → **404**, environment belum dibuat — normal.

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
