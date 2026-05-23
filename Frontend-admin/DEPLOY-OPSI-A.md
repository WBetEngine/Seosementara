# Opsi A — GitHub Actions + GHCR + mini PC

Panduan implementasi **permanen & bersih**: kode di GitHub, mini PC hanya Docker + `.env`.

## Alur

1. Push ke `main` → workflow `.github/workflows/deploy-backend.yml`
2. GitHub build image (Go + partials admin dibundle)
3. Push ke `ghcr.io/wbetengine/seosementara-api:latest` (+ tag commit SHA)
4. Mini PC: `docker compose pull` (otomatis via SSH/webhook atau manual)

Admin UI tetap di **Cloudflare Workers** — tidak ikut image backend.

---

## Bagian 1 — GitHub (sekali)

### Package GHCR

Setelah workflow pertama sukses:

1. Repo → **Packages** → `seosementara-api`
2. **Package settings** → link ke repo `Seosementara` (jika private)
3. Untuk repo private, mini PC wajib `docker login ghcr.io` (PAT scope `read:packages`)

### Deploy otomatis ke mini PC (pilih satu trigger)

#### A) SSH (disarankan jika mini PC bisa di-SSH)

**Variables** (Settings → Actions → Variables):

| Nama | Nilai |
|------|--------|
| `MINI_PC_DEPLOY` | `true` |

**Secrets**:

| Nama | Contoh |
|------|--------|
| `DEPLOY_HOST` | IP publik / hostname mini PC |
| `DEPLOY_USER` | `Administrator` atau user SSH |
| `DEPLOY_SSH_KEY` | Private key OpenSSH (pasangan key di mini PC) |
| `DEPLOY_SSH_PORT` | `22` (opsional) |
| `DEPLOY_PATH` | `C:/Seosementara` |

**Windows — OpenSSH Server:**

1. Settings → Apps → Optional features → **OpenSSH Server**
2. `C:\ProgramData\ssh\administrators_authorized_keys` untuk key deploy
3. Firewall: allow port 22 (atau port custom)

Di mini PC, pastikan ada `docker-compose.prod.yml` + `.env` + `Backend/migrations/`.

#### B) Webhook (mini PC di belakang NAT, tanpa SSH inbound)

**Variable:** `MINI_PC_DEPLOY` = `webhook`

**Secrets:** `DEPLOY_WEBHOOK_URL`, `DEPLOY_WEBHOOK_TOKEN`

Jalankan listener di mini PC (contoh PowerShell terjadwal atau Task Scheduler) yang memanggil `scripts/mini-pc-pull.ps1` saat webhook diterima — implementasi listener sesuai keamanan Anda (IIS, nginx, atau skrip lokal).

#### C) Manual pull (paling sederhana)

Biarkan `MINI_PC_DEPLOY` kosong. Setelah setiap push ke `main`, di mini PC:

```powershell
cd C:\Seosementara
.\scripts\mini-pc-pull.ps1
```

---

## Bagian 2 — Mini PC (sekali)

Lihat `mini-pc/README.md`. Ringkas:

```
C:\Seosementara\
  docker-compose.prod.yml
  .env
  Backend\migrations\   ← hanya folder SQL
```

**Tidak** perlu clone/ZIP seluruh repo.

```powershell
docker login ghcr.io
cd C:\Seosementara
docker compose -f docker-compose.prod.yml up -d
```

---

## Bagian 3 — Develop lokal (opsional)

Di laptop dengan source lengkap:

```bash
docker compose up -d --build
```

Ini memakai `docker-compose.yml` (build lokal + volume partials untuk iterasi cepat).

---

## Troubleshooting

| Masalah | Solusi |
|---------|--------|
| `pull access denied` | `docker login ghcr.io` + PAT `read:packages` |
| Image tidak update | Cek Actions hijau; lalu `pull` di mini PC |
| `failed to read dockerfile` | Jangan pakai compose dev tanpa source; pakai **prod** compose |
| Health gagal | `docker compose -f docker-compose.prod.yml logs api` |

---

## Ringkas

| Komponen | Lokasi |
|----------|--------|
| Kode | GitHub |
| Build | GitHub Actions |
| Image | GHCR |
| Runtime | Mini PC (`docker-compose.prod.yml` + `.env`) |
| Admin UI | Cloudflare Workers |
