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

### Environment `production` (404 = belum dibuat)

URL ini **404** sampai environment `production` ada di repo:

https://github.com/WBetEngine/Seosementara/settings/environments/production

**Buat environment (pilih salah satu):**

| Cara | Kapan |
|------|--------|
| **Actions → Ensure Production Environment → Run workflow** | Sekarang — tidak butuh secret Cloudflare |
| **Actions → Deploy Admin UI** (sukses) | Otomatis membuat environment |
| **Bootstrap Platform** di admin | Menulis secret ke environment |

Daftar environment: https://github.com/WBetEngine/Seosementara/settings/environments

Secret infra mini PC (`DB_PASSWORD`, dll.) baru muncul di environment **setelah Bootstrap + Infra** di admin — bukan isi manual di halaman 404.

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

Runner **v2.334+** tidak punya `install.cmd` — service dipasang lewat `config.cmd --runasservice`.
Jika sudah config tanpa service, jalankan `C:\actions-runner\run.cmd` (biarkan window terbuka)
atau config ulang dengan token baru + `--runasservice`.

**GHCR pull gagal `denied`:** buat package public di  
https://github.com/orgs/WBetEngine/packages/container/seosementara-api/settings  
atau tambah permission **Packages: Read** pada PAT Bootstrap (`PLATFORM_GITHUB_PAT`).

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
