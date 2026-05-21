# 04 — Backend Golang (Mini CPU)

## 1. Mengapa Golang

| Kriteria | Golang di mini CPU |
|----------|-------------------|
| Memori | Binary tunggal, footprint rendah vs runtime berat |
| Concurrency | Goroutine + worker pool cocok untuk job batch |
| Deploy | Satu binary cross-compile, mudah di systemd |
| Performa | Cukup untuk API I/O-bound + DB |

## 2. Struktur Proyek (Usulan)

```
Backend/
├── cmd/
│   ├── api/          # HTTP server
│   └── worker/       # Job processor (bisa digabung: -mode worker)
├── internal/
│   ├── config/
│   ├── domain/       # Entity & business rules
│   ├── handler/      # HTTP handlers
│   ├── repository/   # DB access
│   ├── service/      # Use cases
│   ├── middleware/   # Auth, CORS, logging, rate limit
│   └── worker/       # Job handlers
├── migrations/
├── go.mod
└── README.md
```

## 3. Library yang Direkomendasikan

| Kebutuhan | Library (contoh) |
|-----------|------------------|
| HTTP router | `chi` atau std `net/http` ServeMux (Go 1.22+) |
| Config | `env` / file YAML |
| DB | `pgx` (PostgreSQL) atau `modernc.org/sqlite` |
| Migrasi | `goose` atau `golang-migrate` |
| Validasi | `go-playground/validator` |
| Log | `slog` (stdlib) |
| Auth JWT/session | `golang-jwt` + secure cookie |
| Test | `testing` + `testify` |

## 4. Pola Arsitektur

**Clean-ish layering:** Handler → Service → Repository.

- **Handler:** parse request, status code, tidak ada SQL
- **Service:** aturan bisnis, transaksi, enqueue job
- **Repository:** query terarah, selalu parameterized

## 5. Aturan Query (Wajib — Skala Massal)

| Larangan | Alternatif |
|----------|------------|
| `SELECT *` tanpa `LIMIT` pada list | Pagination: `limit` + `cursor` atau `offset` max 100 |
| Load semua post satu situs | Filter + page; export via job |
| N+1 query di loop | `JOIN` atau batch `IN (...)` |
| Count penuh tabel tiap request | Cache count 5–15 menit atau perkiraan |
| Operasi DELETE tanpa `WHERE` spesifik | Selalu `site_id` + `id` / status filter |

### Contoh pagination (konsep)

```go
// ListPosts(siteID, status, cursor, limit int) ([]Post, nextCursor, error)
// Index DB: (site_id, status, updated_at DESC)
```

## 6. Job Queue & Worker

Operasi berat **tidak** di HTTP handler utama:

1. Admin trigger → `POST /api/admin/jobs` → insert row `jobs` status `pending`
2. Worker poll (atau NOTIFY) → proses chunk 50 item
3. Update progress → admin poll `GET /api/admin/jobs/{id}`

| Job type | Contoh |
|----------|--------|
| `bulk_publish` | Publish 500 draft |
| `bulk_seo_update` | Update meta field |
| `regenerate_sitemap` | XML per situs |
| `purge_cache` | Invalidate CDN/cache keys |

**Concurrency:** max 2–4 job paralel di mini CPU (configurable).

## 7. Autentikasi & Otorisasi

| Tipe | Mekanisme |
|------|-----------|
| Admin API | Session cookie HttpOnly **atau** JWT short-lived + refresh |
| Public API | Tidak ada auth untuk read; API key untuk form tertentu |
| RBAC | Middleware cek `role` + `site_id` scope |

Password: `bcrypt` atau `argon2id`.

## 8. Media & File

- Upload: `multipart` → validasi ukuran/MIME → simpan disk atau R2
- Serving: URL signed atau public path via reverse proxy
- Thumbnail: generate async di worker (jangan block upload response)

## 9. Cache

| Lapisan | Strategi |
|---------|----------|
| In-memory | LRU kecil untuk settings hot |
| Redis (opsional) | Response fragment key `site:slug:post` |
| HTTP | `Cache-Control` untuk public read endpoint |

Invalidasi: pada `publish`, `update`, `delete` → hapus key terkait situs.

## 10. Health & Observabilitas

```
GET /health        → { "status": "ok", "db": "ok" }
GET /health/ready  → cek DB + disk space
```

Metric opsional: Prometheus endpoint `/metrics` (ringan).

## 11. Konfigurasi Lingkungan

```env
APP_ENV=production
HTTP_ADDR=:8080
DATABASE_URL=postgres://...
CORS_ADMIN_ORIGIN=https://admin.example.pages.dev
CORS_PUBLIC_ORIGIN=https://*.pages.dev
WORKER_CONCURRENCY=2
MAX_UPLOAD_MB=10
```

## 12. Systemd (Mini CPU)

Dua unit: `seosementara-api.service`, `seosementara-worker.service` — restart on failure, memory limit opsional.

## 13. Dokumen Terkait

- API kontrak → [07-api-dan-integrasi.md](./07-api-dan-integrasi.md)
- Menu admin → [03-menu-dan-modul-cms.md](./03-menu-dan-modul-cms.md)
- Infrastruktur → [02-arsitektur-dan-infrastruktur.md](./02-arsitektur-dan-infrastruktur.md)
