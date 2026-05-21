# 07 — API dan Integrasi Antar Lapisan

## 1. Prinsip API

| Prinsip | Implementasi |
|---------|--------------|
| RESTful | Resource-oriented URL, verb HTTP standar |
| Versioning | Prefix `/api/v1` (opsional di MVP) |
| Format admin | HTML partial **atau** JSON — MVP: partial HTML untuk HTMX admin |
| Format public | HTML partial untuk HTMX customer |
| Error | JSON `{ "error": { "code", "message" } }` atau HTML alert fragment |

## 2. Pembagian Namespace

```
/api/admin/*     → Auth required, RBAC
/api/public/*    → Read mostly, rate limited
/health            → No auth
```

## 3. Autentikasi

### Login admin

```http
POST /api/admin/auth/login
Content-Type: application/json

{ "email": "...", "password": "..." }
```

Response:

```http
Set-Cookie: session=...; HttpOnly; Secure; SameSite=None; Path=/
```

### Request berikutnya

```http
GET /api/admin/posts?site_id=1&limit=20
Cookie: session=...
X-Site-ID: 1
```

### Logout

```http
POST /api/admin/auth/logout
```

## 4. CORS (Cloudflare Pages → Mini CPU)

| Header | Nilai |
|--------|-------|
| `Access-Control-Allow-Origin` | Origin admin/users Pages (whitelist) |
| `Access-Control-Allow-Credentials` | `true` |
| `Access-Control-Allow-Methods` | GET, POST, PUT, PATCH, DELETE, OPTIONS |
| `Access-Control-Allow-Headers` | Content-Type, X-Site-ID, HX-Request |

Preflight `OPTIONS` harus di-handle cepat (< 50ms).

## 5. Endpoint Admin (Ringkas)

### Situs

| Method | Path | Deskripsi |
|--------|------|-----------|
| GET | `/api/admin/sites` | List paginated |
| POST | `/api/admin/sites` | Buat situs |
| GET | `/api/admin/sites/{id}` | Detail |
| PATCH | `/api/admin/sites/{id}` | Update |

### Post

| Method | Path | Deskripsi |
|--------|------|-----------|
| GET | `/api/admin/posts` | List (`site_id`, `status`, `cursor`, `limit`) |
| POST | `/api/admin/posts` | Buat draft |
| GET | `/api/admin/posts/{id}` | Detail + SEO meta |
| PATCH | `/api/admin/posts/{id}` | Update |
| POST | `/api/admin/posts/{id}/publish` | Publish |
| DELETE | `/api/admin/posts/{id}` | Soft delete |

### Media

| Method | Path | Deskripsi |
|--------|------|-----------|
| GET | `/api/admin/media` | List |
| POST | `/api/admin/media/upload` | Multipart upload |
| DELETE | `/api/admin/media/{id}` | Hapus |

### SEO

| Method | Path | Deskripsi |
|--------|------|-----------|
| GET | `/api/admin/seo/site/{site_id}` | Settings |
| PATCH | `/api/admin/seo/site/{site_id}` | Update |
| POST | `/api/admin/seo/bulk` | Enqueue bulk job |

### Jobs

| Method | Path | Deskripsi |
|--------|------|-----------|
| GET | `/api/admin/jobs` | List |
| GET | `/api/admin/jobs/{id}` | Status + progress |
| POST | `/api/admin/jobs/{id}/retry` | Retry failed |

### Dashboard

| Method | Path | Deskripsi |
|--------|------|-----------|
| GET | `/api/admin/dashboard` | HTML partial ringkasan |

## 6. Endpoint Public (Customer)

| Method | Path | Deskripsi | Cache |
|--------|------|-----------|-------|
| GET | `/api/public/sites/by-host` | Resolve site dari Host header | 5m |
| GET | `/api/public/home` | Fragment beranda | 60s |
| GET | `/api/public/posts` | List published | 60s |
| GET | `/api/public/posts/{slug}` | Artikel | 5m |
| GET | `/api/public/pages/{slug}` | Halaman statis | 5m |
| GET | `/api/public/sitemap.xml` | Sitemap | 1h |
| POST | `/api/public/forms/contact` | Form kontak | no cache |

Semua public list **wajib** `limit` default 20, max 50.

## 7. Header HTMX

Backend dapat mendeteksi request HTMX:

| Header | Penggunaan |
|--------|------------|
| `HX-Request: true` | Return partial tanpa layout penuh |
| `HX-Target` | Validasi target yang diizinkan |
| `HX-Trigger` | Logging / analytics |

Response:

| Header | Penggunaan |
|--------|------------|
| `HX-Redirect` | Redirect setelah login |
| `HX-Trigger` | Toast notification client |
| `HX-Retarget` | Ganti elemen swap |

## 8. Rate Limiting

| Kelompok | Limit (contoh) |
|----------|----------------|
| Public read | 120 req/menit per IP |
| Public form POST | 10 req/menit per IP |
| Admin API | 300 req/menit per user |
| Login | 5 percobaan/menit per IP |

Implementasi: middleware token bucket in-memory atau Redis.

## 9. Invalidasi Cache (Alur Publish)

```mermaid
sequenceDiagram
  participant Admin as Admin HTMX
  participant API as Go API
  participant Worker as Worker
  participant CF as Cloudflare

  Admin->>API: POST publish post
  API->>API: Update DB status=published
  API->>Worker: Enqueue purge_cache
  Worker->>CF: Purge URL (optional API)
  API-->>Admin: 200 + fragment updated
```

## 10. Kontrak Error

| HTTP | Code | Arti |
|------|------|------|
| 400 | `validation_error` | Input tidak valid |
| 401 | `unauthorized` | Belum login |
| 403 | `forbidden` | RBAC gagal |
| 404 | `not_found` | Resource tidak ada |
| 409 | `conflict` | Slug duplikat |
| 429 | `rate_limited` | Terlalu banyak request |
| 500 | `internal_error` | Log di server, pesan generik ke client |

## 11. Webhook (Opsional)

Admin pengaturan dapat mendaftarkan webhook:

- `content.published`
- `job.completed`
- `job.failed`

Payload JSON + HMAC signature.

## 12. Dokumen Terkait

- Backend implementasi → [04-backend-golang.md](./04-backend-golang.md)
- Admin HTMX → [05-admin-panel-htmx.md](./05-admin-panel-htmx.md)
- Users HTMX → [06-frontend-users-htmx.md](./06-frontend-users-htmx.md)
