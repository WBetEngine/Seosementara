# Backend — Cloudflare Settings (produksi)

Modul **Settings → Cloudflare** sesuai Plan/15.

## Migrasi PostgreSQL

```bash
psql "$DATABASE_URL" -f migrations/002_cloudflare_setup.up.sql
```

Atau via Docker Compose (otomatis dari `docker-entrypoint-initdb.d` pada DB baru).

## Jalankan

```bash
export MASTER_ENCRYPTION_KEY="$(openssl rand -base64 32)"
export SUPER_ADMIN_TOKEN="your-secret"
export DATABASE_URL="postgres://..."
go run ./cmd/api
```

Docker (dari root repo):

```bash
cp mini-pc/env.example .env
# isi MASTER_ENCRYPTION_KEY dan SUPER_ADMIN_TOKEN
docker compose up -d --build
```

## Endpoint

### HTMX HTML (partial)

| GET | Path |
|-----|------|
| Shell | `/api/admin/settings/cloudflare/` |
| Koneksi | `/api/admin/settings/cloudflare/koneksi` |
| Domain | `/api/admin/settings/cloudflare/domain` |
| Tunnel | `/api/admin/settings/cloudflare/tunnel` |
| Pages | `/api/admin/settings/cloudflare/pages` |
| DNS | `/api/admin/settings/cloudflare/dns` |

Header: `Authorization: Bearer $SUPER_ADMIN_TOKEN`

### JSON API

Prefix: `/api/admin/setup/cloudflare/` (sama handler)

| Method | Path | Fungsi |
|--------|------|--------|
| GET/PUT | `/credentials` | Token / Global API Key |
| POST | `/credentials/test` | Uji ke Cloudflare API |
| GET/PUT | `/domain-env` | Variabel domain |
| POST | `/domain-env/sync-pages` | Sync ke Pages env |
| GET | `/tunnel` | Config + routes |
| POST | `/tunnel` | Buat tunnel + install command |
| POST | `/tunnel/routes/apply` | Push ingress ke CF |
| GET | `/tunnel/status` | Cek connector |
| PUT | `/tunnel/route` | Simpan route |
| GET/PUT | `/pages` | Proyek Pages admin |
| POST | `/pages/deploy` | Trigger deployment |
| POST | `/dns/apply` | Buat record DNS |
| GET | `/logs` | Audit API CF |

## Admin UI (Workers) + API

Saat Tunnel mengarahkan `https://seosementara.org/api/*` ke Go:

1. Di `index.html` set `hx-headers` dengan token Super Admin.
2. Ubah `hx-get` Cloudflare ke `/api/admin/settings/cloudflare/...`.

Tanpa Tunnel, UI Workers tetap memakai partial statis di `/_partials/`.
