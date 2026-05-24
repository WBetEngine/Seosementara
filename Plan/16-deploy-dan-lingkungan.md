# 16 — Deploy & Lingkungan (Dev / Staging / Production)

> Menyatukan alur rilis Backend, Admin Workers, Tunnel, PostgreSQL.  
> **Production mini PC:** tanpa file `.env` — lihat [28-platform-github-workers](./28-platform-github-workers.md).

## Implementasi v2 (production mini PC)

| Langkah | Tool | Folder |
|---------|------|--------|
| Onboarding first boot | GitHub Pages + Workers Platform API | `Frontend-Onboarding/` — [28](./28-platform-github-workers.md) |
| Build API image | `deploy-backend.yml` → GHCR | `Backend/` |
| Deploy Docker | `deploy-mini-pc.yml` → self-hosted runner | secrets inject, tanpa `.env` file |
| Deploy admin UI | `deploy-admin.yml` → Cloudflare Pages | `Frontend-Ui-Admin/public/` |
| Deploy publik UI | `deploy-public.yml` → Cloudflare Pages | `Frontend-Publik/public/` |
| Onboarding GitHub Pages | `pages-onboarding.yml` | `Frontend-Onboarding/public/` |

Secrets infra (`DB_PASSWORD`, `MASTER_ENCRYPTION_KEY`) diisi dari onboarding / Settings admin, bukan commit ke repo.

---

## 1. Tujuan (asli)

| Tujuan | Keterangan |
|--------|------------|
| Tiga lingkungan jelas | **Local**, **Staging**, **Production** — tidak campur data/secret |
| Deploy dapat diulang | Script / CI sama setiap rilis |
| Mini CPU aman | Rolling restart tanpa buka port publik |
| UI di edge | Pages deploy terpisah dari binary Go |
| Konfigurasi | GitHub Secrets + `domain_env_config` + Settings admin per lingkungan |

---

## 2. Ringkasan Lingkungan

| | **Local** | **Staging** | **Production** |
|--|-----------|-------------|----------------|
| **Tujuan** | Dev pekerja | Uji sebelum prod | Live |
| **URL apex** | `localhost:8080` | `staging.seosementara.org` | `seosementara.org` |
| **Backend** | `go run` / binary lokal | Mini CPU (bisa mesin kedua) | Mini CPU utama |
| **PostgreSQL** | Docker / lokal | Instance staging | Instance prod |
| **Tunnel** | Opsional (`cloudflared` dev) | Tunnel `sse-staging` | Tunnel `sse-production` |
| **Pages** | `wrangler pages dev` | Project `*-staging` | Project `*-production` |
| **CF Zone** | — / tunnel only | Subdomain staging | Zone produksi |
| **Data** | Seed / sintetis | Anonim / subset | Real |

```mermaid
flowchart TB
  subgraph local [Local Dev]
    LGo[go run api]
    LDB[(Postgres Docker)]
    LPages[wrangler pages dev]
  end
  subgraph staging [Staging]
    SGo[Go binary]
    STun[Tunnel staging]
    SPages[Pages staging]
    SDB[(Postgres staging)]
  end
  subgraph prod [Production]
    PGo[Go binary]
    PTun[Tunnel prod]
    PPages[Pages prod]
    PDB[(Postgres prod)]
  end
  Dev --> local
  CI --> staging
  CI -->|manual approve| prod
```

---

## 3. Artefak Deploy (Per Komponen)

| Komponen | Artefak | Target (Production) | Tool |
|----------|---------|---------------------|------|
| **Backend API** | Docker image `sse-api` (GHCR) | Mini PC container | GitHub Actions + runner |
| **Worker jobs** | Image `sse-worker` atau sidecar | Mini PC container | docker compose |
| **Migrasi DB** | `migrations/*.sql` | PostgreSQL container | `goose` / `migrate` |
| **cloudflared** | Config + tunnel token | Mini PC | Onboarding [28] + Settings [15] |
| **Onboarding UI** | Static HTML | GitHub Pages | push `main` |
| **Admin UI** | `Frontend-Ui-Admin/public/` | Cloudflare Pages | GitHub Actions + `wrangler` |
| **Publik UI** | `Frontend-Publik/public/` | Cloudflare Pages | GitHub Actions + `wrangler` |

**Bukan satu deploy monolith** — pipeline terpisah; urutan bootstrap: [28](./28-platform-github-workers.md) §4.

> **Legacy (dev lokal):** binary ke `/opt/seosementara/bin/` + systemd — bukan alur production.

---

## 4. Variabel Lingkungan

### 4.1 Server mini PC — **bukan** file `.env` di Git (production)

| Variable | Sumber |
|----------|--------|
| `DB_PASSWORD` | GitHub Secret ← admin **Infra & GitHub** |
| `MASTER_ENCRYPTION_KEY` | GitHub Secret ← admin |
| `SUPER_ADMIN_TOKEN` | GitHub Secret ← admin (opsional) |
| Inject ke container | Self-hosted runner → `docker compose` env |

Dev lokal boleh pakai env process atau file lokal — tidak di-commit.

### 4.1b Legacy (deprecated): file `/etc/seosementara/env`

| Variable | Local | Staging | Production |
|----------|-------|---------|------------|
| `APP_ENV` | `local` | `staging` | `production` |
| `HTTP_ADDR` | `127.0.0.1:8080` | `127.0.0.1:8080` | `127.0.0.1:8080` |
| `DATABASE_URL` | local DSN | staging DSN | prod DSN |
| `SESSION_SECRET` | dev-only | staging secret | prod secret (kuat) |
| `MASTER_ENCRYPTION_KEY` | dev-only | staging | prod |
| `LOG_LEVEL` | `debug` | `info` | `info` |

### 4.2 Domain utama (DB `domain_env_config` + sync Pages — [15])

| Key | Staging contoh | Production contoh |
|-----|----------------|-------------------|
| `PRIMARY_DOMAIN` | `staging.seosementara.org` | `seosementara.org` |
| `APEX_URL` | `https://staging.seosementara.org` | `https://seosementara.org` |
| `API_BASE_URL` | sama dengan APEX_URL | sama |
| `ENVIRONMENT` | `staging` | `production` |

Diset lewat **`/admin/settings/cloudflare/domain-utama`** per lingkungan (staging admin terpisah atau flag di DB).

### 4.3 Cloudflare (DB terenkripsi — [15])

| Resource | Staging | Production |
|----------|---------|------------|
| API Token | Token scoped staging (disarankan) | Token scoped prod |
| Tunnel name | `sse-staging` | `sse-production` |
| Pages project admin | `seosementara-admin-staging` | `seosementara-admin` |
| Pages project public | `seosementara-public-staging` | `seosementara-public` |

**Jangan** pakai tunnel/production token yang sama.

---

## 5. Struktur Repo & Branch

| Branch | Lingkungan | Deploy otomatis |
|--------|------------|-----------------|
| `main` | Production | Ya (dengan approval manual) |
| `staging` | Staging | Ya (setiap push) |
| `feature/*` | — | CI test saja |

```
.github/workflows/
├── ci.yml              → test + lint Go, validate SQL
├── deploy-staging.yml  → push staging branch
└── deploy-production.yml → workflow_dispatch / tag v*
```

---

## 6. Urutan Deploy (Runbook)

### 6.1 Production — urutan wajib (post-bootstrap)

```text
1. Backup PostgreSQL (pg_dump)
2. Migrasi DB (goose up) — backward-compatible only
3. Pull & restart Docker containers (API + worker)
4. Cek /health & /health/ready via Tunnel
5. Deploy Pages (admin + public) — jika UI berubah
6. Sync env Pages dari Settings admin (jika domain env berubah)
7. Apply Tunnel routes (jika berubah) — biasanya jarang
8. Purge cache Cloudflare (opsional)
9. Smoke test: login admin, list domain, 1 API publik
```

| Langkah gagal | Rollback |
|---------------|----------|
| Migrasi | `goose down 1` + restore dump |
| Container API | Redeploy image tag sebelumnya |
| Pages | Redeploy deployment sebelumnya di dashboard CF |

### 6.2 Staging

Sama seperti production, tanpa backup formal wajib — bisa reset DB dari seed.

### 6.3 Local (opsional — hanya dev pekerja)

```bash
docker compose up -d postgres
cd Backend && go run ./cmd/api
cd Frontend-Ui-Admin && npx wrangler pages dev public
```

**Bukan** alur onboarding produksi. Operator first boot memakai GitHub Pages ([28](./28-platform-github-workers.md)).

---

## 7. Backend Go — Detail Deploy (Mini PC Docker)

### 7.1 Build image (CI)

```bash
cd Backend
docker build -t ghcr.io/<org>/sse-api:${GIT_SHA} .
docker push ghcr.io/<org>/sse-api:${GIT_SHA}
```

Self-hosted runner di mini PC:

```bash
docker compose pull && docker compose up -d
```

Secrets inject via GitHub Actions env — **tanpa** file `.env` di disk.

### 7.2 Layout production (Docker)

```text
/opt/sse-docker/          # hanya compose + volume mounts
├── docker-compose.yml
├── data/postgres/
└── data/media/
```

> **Legacy dev:** binary di `/opt/seosementara/bin/` + systemd — lihat §7.3 jika masih dipakai lokal.

### 7.3 systemd (legacy / dev binary — opsional)

```ini
# /etc/systemd/system/sse-api.service — HANYA jika tidak pakai Docker
[Service]
ExecStart=/opt/seosementara/bin/sse-api
EnvironmentFile=/etc/seosementara/env
Restart=on-failure
```

```ini
# sse-worker.service — terpisah agar job tidak mati saat API restart
```

| Dampak | |
|--------|--|
| Restart API < 2 detik | Session di DB tetap valid |
| Worker terpisah | Bulk job lanjut saat deploy API |

### 7.4 Migrasi database

| Aturan | Dampak |
|--------|--------|
| Migrasi **forward-only** di prod | Hindari down yang hapus kolom dipakai |
| Index baru | `CREATE INDEX CONCURRENTLY` di migrasi terpisah |
| Seed | Hanya local/staging |

```bash
goose -dir migrations postgres "$DATABASE_URL" up
```

### 7.5 cloudflared

| Langkah | |
|---------|--|
| Token dari Setup admin [15] | |
| `cloudflared service install <token>` | |
| Config routes dari DB → apply via admin | |

Service terpisah: `cloudflared.service` — restart tidak wajib saat deploy Go kecuali port berubah.

---

## 8. Cloudflare Pages — Detail Deploy

### 8.1 Dua proyek × dua lingkungan = 4 project (disarankan)

| Project | Branch CF Pages | Folder |
|---------|-----------------|--------|
| `seosementara-admin-staging` | `staging` | `Frontend-Ui-Admin` |
| `seosementara-admin` | `main` | `Frontend-Ui-Admin` |
| `seosementara-public-staging` | `staging` | `Frontend-Publik` |
| `seosementara-public` | `main` | `Frontend-Publik` |

### 8.2 GitHub Actions (contoh)

```yaml
# deploy-staging.yml — ringkas
jobs:
  pages-admin:
    steps:
      - uses: actions/checkout@v4
      - run: npx wrangler pages deploy Frontend-Ui-Admin/public \
          --project-name=seosementara-admin-staging
    env:
      CLOUDFLARE_API_TOKEN: ${{ secrets.CF_TOKEN_STAGING }}
```

| Secret GitHub | Isi |
|---------------|-----|
| `CF_TOKEN_STAGING` | API token scoped Pages staging |
| `CF_TOKEN_PRODUCTION` | API token scoped Pages prod |
| `SSH_DEPLOY_KEY` | Opsional: rsync binary ke mini CPU |

### 8.3 Deploy dari admin panel (opsional fase 2)

Tombol di [15] §7.3: trigger `wrangler pages deploy` via job worker — token dari `cloudflare_credentials`.

| Pro | Kontra |
|-----|--------|
| Super Admin tidak perlu GitHub | Build di mini CPU lambat — **disarankan tetap CI utama** |

---

## 9. CI Pipeline (`ci.yml`)

| Job | Isi |
|-----|-----|
| `go-test` | `go test ./...`, race detector opsional |
| `go-lint` | `staticcheck` / `golangci-lint` |
| `sql-check` | Validasi migrasi goose |
| `build` | Artefak `sse-api`, `sse-worker` — upload artifact |

Tidak deploy ke prod dari PR — hanya dari `main` + approval.

---

## 10. Routing per Lingkungan

| Path | Local | Staging/Prod |
|------|-------|--------------|
| UI `/admin/*` | Pages dev atau Go | Pages + route |
| UI `/` | Pages dev | Pages |
| `/api/*` | Go :8080 | Tunnel → Go |
| `*.domain` | Hosts file / staging DNS | CF DNS |

**Staging** memakai subdomain `staging.` agar tidak bentrok cookie/session dengan prod.

---

## 11. Smoke Test Pasca-Deploy

| # | Cek | Harapan |
|---|-----|---------|
| 1 | `GET /health` | `200` `db: ok` |
| 2 | `GET /health/ready` | disk, tunnel ok |
| 3 | Login admin | Redirect `/admin/` |
| 4 | `GET /api/admin/dashboard` | 200 + HTML partial |
| 5 | Halaman publik apex | 200 |
| 6 | Satu subdomain contoh | 200 |
| 7 | Tunnel connector | healthy di Setup admin |
| 8 | Pages deployment | success di CF |

Otomatisasi: script `scripts/smoke.sh` dengan `curl` + exit code.

---

## 12. Rollback & Darurat

| Keadaan | Tindakan |
|---------|----------|
| Bug kritis API | Rollback binary `.bak` + restart |
| Migrasi rusak | Restore dump + fix forward migration |
| Pages rusak | Rollback deployment CF |
| Tunnel mati | Restart `cloudflared`; cek token |
| Maintenance | `app.maintenance=true` dari Setup backend [13] — tanpa deploy |

**RTO target (internal):** API rollback < 5 menit jika binary `.bak` siap.

---

## 13. Integrasi Setup Admin Panel

| Aksi deploy | Dari mana |
|-------------|-----------|
| Lihat versi terdeploy | `/admin/settings/backend/ringkasan` — `GIT_SHA`, `build_time` |
| Trigger Pages deploy | `/admin/settings/cloudflare/pages` (fase 2) |
| Env domain | `/admin/settings/cloudflare/domain-utama` |
| Maintenance mode | `/admin/settings/backend/operasional` |
| Health tunnel | `/admin/settings/cloudflare/tunnel` |

Backend expose:

```json
GET /api/admin/settings/backend/overview
{
  "version": "abc123",
  "env": "production",
  "uptime_sec": 86400,
  "tunnel_status": "healthy",
  "last_deploy_at": "..."
}
```

---

## 14. Matriks Skenario & Dampak

| # | Skenario | Dampak | Mitigasi |
|---|----------|--------|----------|
| D1 | Deploy API tanpa migrasi | 500 error | CI cek migrasi pending |
| D2 | Migrasi lock tabel lama | Timeout admin | CONCURRENTLY, maintenance window |
| D3 | Pages env salah | HTMX panggil API prod dari staging | Env terpisah per project |
| D4 | Binary salah arch | Cannot execute | CI build matrix match CPU |
| D5 | Dua worker aktif (blue/green salah) | Duplicate job | Satu worker enabled |
| D6 | Secret prod di staging | Kebocoran | Token CF terpisah |
| D7 | Rollback tanpa down migration | Schema mismatch | Hanya forward-compatible migration |
| D8 | cloudflared tidak restart | OK — tunnel independen | |
| D9 | Deploy Jumat sore | Incident weekend | Deploy staging Kamis, prod awal minggu |

---

## 15. Checklist Pertama Kali (Bootstrap Production)

**First boot:** ikuti [28-platform-github-workers.md](./28-platform-github-workers.md) — wizard **GitHub Pages onboarding**, bukan checklist manual SSH.

Ringkasan hasil bootstrap:

1. [ ] Workers Platform API online  
2. [ ] GitHub Secrets terisi (DB, encryption)  
3. [ ] Self-hosted runner terdaftar di mini PC  
4. [ ] Tunnel + routes aktif  
5. [ ] Docker API + PostgreSQL running  
6. [ ] `goose up` migrasi awal  
7. [ ] Super Admin pertama (seed CLI)  
8. [ ] CF Pages admin + publik deployed  
9. [ ] Smoke test §11  
10. [ ] Backup otomatis harian (cron `pg_dump`)  

---

## 16. Roadmap Implementasi Deploy

| Fase | Deliverable |
|------|-------------|
| MVP | Onboarding GH Pages + Docker deploy via Actions |
| Fase 2 | GitHub Actions staging auto, prod manual approve |
| Fase 3 | Smoke otomatis, overview versi di Settings admin |

---

## 17. Dokumen Terkait

- [02-arsitektur-dan-infrastruktur.md](./02-arsitektur-dan-infrastruktur.md)
- [04-backend-golang.md](./04-backend-golang.md)
- [05-admin-panel-htmx.md](./05-admin-panel-htmx.md)
- [06-frontend-users-htmx.md](./06-frontend-users-htmx.md)
- [08-roadmap-implementasi.md](./08-roadmap-implementasi.md)
- [13-setup-backend-dan-sistem.md](./13-setup-backend-dan-sistem.md)
- [15-setup-cloudflare-integrasi.md](./15-setup-cloudflare-integrasi.md)
