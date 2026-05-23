# Deploy mini PC — GitHub pusat, tanpa file `.env`

## Arsitektur

```text
GitHub (kode + GHCR + Secrets)
        │
        ├─► Workers admin (setup UI)
        │     ├─ Infra → GitHub Secrets → workflow Deploy Mini PC
        │     └─ Cloudflare → Workers Secrets
        │
        └─► Self-hosted runner (mini PC)
              └─ docker compose (env inject dari GitHub Secrets)
```

| Setup | Di mana | Penyimpanan |
|-------|---------|-------------|
| DB password, encryption key | Admin → **Settings → Infra & GitHub** | GitHub Secrets → Docker inject |
| Global API Key Cloudflare | Admin → **Settings → Cloudflare → Koneksi** | Workers Secrets |
| Domain, tunnel, DNS | Admin → tab Cloudflare (Go API) | PostgreSQL |
| Kode & image | GitHub `main` | GHCR |

**Tidak ada file `.env` di mini PC.**

---

## Bootstrap sekali (mini PC)

### 1. Runner + Docker + cloudflared

```powershell
# Administrator — lihat scripts/install-github-runner.ps1
cd C:\Seosementara
git clone https://github.com/WBetEngine/Seosementara.git .
.\scripts\install-github-runner.ps1
```

### 2. GitHub Environment `production`

Repo → **Settings → Environments → production** — tambah secrets bootstrap:

| Secret | Fungsi |
|--------|--------|
| `GITHUB_SETUP_TOKEN` | PAT: repo secrets write + actions write (→ Worker via deploy-admin) |
| `CLOUDFLARE_API_KEY` | Deploy Wrangler (Global API Key) |
| `CLOUDFLARE_ACCOUNT_EMAIL` | Deploy Wrangler |
| `CLOUDFLARE_ACCOUNT_ID` | Deploy Wrangler |

Secrets **DB / encryption** diisi lewat admin (langkah 3), bukan manual di GitHub.

### 3. Setup lewat admin Workers

1. Buka `https://seosementara.seosementara3.workers.dev/admin/settings/backend/infra`  
   → isi `DB_PASSWORD`, `MASTER_ENCRYPTION_KEY` → **Simpan & deploy mini PC**

2. Buka **Settings → Cloudflare → Koneksi**  
   → Global API Key + email + Account ID → **Simpan ke Workers**

3. Tab Tunnel / Domain / DNS → konfigurasi via Go API (`api.apidevel.org`)

---

## Update rutin

Push `main` → build GHCR → runner deploy otomatis.

Manual: **Actions → Deploy Mini PC → Run workflow**

---

## File di mini PC

```
C:\Seosementara\
  docker-compose.prod.yml
  Backend\migrations\
  scripts\mini-pc-deploy.ps1
```

---

## Dev lokal (opsional)

Developers boleh pakai `docker compose up` dengan env lokal — bukan alur production.
