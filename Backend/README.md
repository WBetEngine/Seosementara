# Backend — Seosementara CMS

## Pixel Facebook Pro

Modul pertama dari Pixel Hub. Menjalankan Conversions API, first-party collect, dan admin UI.

```bash
openssl rand -base64 32   # PIXEL_ENCRYPTION_KEY
export PIXEL_ENCRYPTION_KEY="..."
export DATABASE_URL="postgres://user:pass@localhost:5432/seosementara?sslmode=disable"  # opsional
export ADMIN_TEMPLATES_DIR="../Frontend-admin/templates"
export STATIC_DIR="../Frontend-admin/static"

go run ./cmd/api
```

Migrasi: jalankan `migrations/001_pixel_hub.up.sql` pada PostgreSQL.

## Struktur

- `cmd/api` — HTTP server + worker dispatch inline
- `internal/pixel/facebook` — Meta CAPI client
- `internal/pixel/service` — use cases
- `internal/pixel/store` — postgres / memory
