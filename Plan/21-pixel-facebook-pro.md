# 21 — Pixel Facebook Pro (Implementasi)

> Spesifikasi fitur **Pro** untuk kanal pertama. Dasar arsitektur: [20-pixel-admin-facebook-tiktok-gads.md](./20-pixel-admin-facebook-tiktok-gads.md).

## 1. Perbedaan Basic vs Pro

| Aspek | Basic (sebelumnya, hanya Plan) | Pro (implementasi saat ini) |
|-------|----------------------------------|-----------------------------|
| UI | Form setup saja | 7 tab: Overview, Setup, Connection, Domains, Diagnostics, Events, Analytics |
| Pelacakan | Snippet FB langsung | **First-party** `sseo-track.js` → `POST /collect` |
| Pengiriman | Manual / tidak ada | Worker **`pixel_dispatch`** → Meta CAPI |
| Credential | - | AES-GCM (`PIXEL_ENCRYPTION_KEY`) |
| Dedup | - | `event_id` UUID per event |
| EMQ | - | `fbp`, `fbc`, hash `em`/`ph` di payload |
| Test | - | Uji koneksi + test event + Test Event Code |
| Skala domain | - | Assign `pixel_domain_assignments` |
| Diagnostik | - | Pending, failure rate 24j, last error |

## 2. Kode & Path

| Komponen | Lokasi |
|----------|--------|
| API server | `Backend/cmd/api/main.go` |
| CAPI client | `Backend/internal/pixel/facebook/capi.go` |
| Service | `Backend/internal/pixel/service/facebook.go` |
| Store | `Backend/internal/pixel/store/` (postgres + memory) |
| Migrasi | `Backend/migrations/001_pixel_hub.up.sql` |
| Admin UI | `Frontend-admin/templates/pixel/facebook/` |
| Skrip | `Frontend-admin/static/js/sseo-track.js` |

## 3. Endpoint

| Method | Path | Fungsi |
|--------|------|--------|
| POST | `/collect` | Ingest event first-party |
| GET | `/sseo-track.js` | Skrip pelacakan |
| GET | `/admin/pixel/facebook/*` | Halaman HTMX |
| POST | `/api/admin/pixel/facebook/setup` | Simpan Pixel ID + token |
| POST | `/api/admin/pixel/facebook/test-connection` | Uji CAPI |
| POST | `/api/admin/pixel/facebook/test-event` | Kirim event uji |
| GET | `/api/admin/pixel/facebook/diagnostics` | JSON diagnostik |
| GET | `/api/admin/pixel/facebook/events` | Log event |
| POST | `/api/admin/pixel/facebook/domains/assign` | Assign domain |

## 4. Menjalankan (dev)

```bash
# Generate key 32 byte base64
openssl rand -base64 32

cd Backend
export PIXEL_ENCRYPTION_KEY="<base64>"
# opsional: export DATABASE_URL="postgres://..."
go run ./cmd/api
```

Buka: `http://localhost:8080/admin/pixel/facebook/`

Tanpa `DATABASE_URL` → store in-memory (dev).

## 5. Roadmap Pro lanjutan (belum di kode)

| Fitur | Prioritas |
|-------|-----------|
| EMQ score pull dari Meta API | Tinggi |
| Webhook / batch export Events Manager | Sedang |
| Job `pixel_deploy_snippet` mass 3000 domain | Tinggi |
| Consent mode + GDPR banner | Sedang |
| Hybrid dedup browser+CAPI otomatis | Sedang |
| OAuth refresh long-lived token | Sedang |

## 6. Dokumen terkait

- [20](./20-pixel-admin-facebook-tiktok-gads.md) — Pixel Hub umum
- [08](./08-roadmap-implementasi.md) — Fase 6b
