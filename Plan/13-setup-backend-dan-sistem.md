# 13 — Setup Backend & Konfigurasi Sistem

> Hanya **Super Admin** — menu di `/admin/setup/*`  
> Bukan pengaturan per domain portfolio (itu di detail domain).

## 1. Ruang Lingkup

| Area | Path admin (contoh) | Fungsi |
|------|----------------------|--------|
| **Host / Subdomain** | `/admin/setup/host` | Hostname produk, template UI — [09](./09-model-domain-host-dan-subdomain.md) |
| **Backend / Sistem** | `/admin/setup/backend` | Konfigurasi operasional platform |
| **Meta global** | `/admin/setup/meta` | Default meta produk — [14](./14-setup-meta-dan-seo.md) |
| **Keamanan** | `/admin/setup/security` | Rate limit, session, kebijakan password |
| **Notifikasi sistem** | `/admin/setup/notifications` | Channel webhook platform |

---

## 2. Menu Setup Backend (`/admin/setup/backend`)

### 2.1 Umum

| Setting | Key | Contoh | Dampak |
|---------|-----|--------|--------|
| Nama platform | `app.name` | Seosementara | UI admin |
| Timezone default | `app.timezone` | Asia/Jakarta | Jadwal publish |
| URL apex | `app.apex_url` | https://seosementara.org | Link email, sitemap |
| Mode maintenance global | `app.maintenance` | false | Semua publik 503 kecuali admin |

### 2.2 Database & performa

| Setting | Key | Dampak |
|---------|-----|--------|
| Pool size hint | `db.pool_size` | Baca di startup Go |
| Query timeout default | `db.query_timeout_ms` | Cegah hung di mini CPU |
| Pagination default | `app.page_size_default` | 50 |
| Pagination max | `app.page_size_max` | 100 |

### 2.3 Worker & job

| Setting | Key | Dampak |
|---------|-----|--------|
| Worker concurrency | `worker.concurrency` | 2–4 di mini CPU |
| Batch size | `worker.batch_size` | 50–200 row per iterasi |
| Job timeout | `worker.job_timeout_sec` | Cegah job menggantung |

### 2.4 Cache

| Setting | Key | Dampak |
|---------|-----|--------|
| Cache enabled | `cache.enabled` | Toggle Redis/memory |
| TTL publik default | `cache.public_ttl_sec` | 60 |
| TTL dashboard stats | `cache.stats_ttl_sec` | 300 |
| Tombol purge all | aksi | Invalidate + Cloudflare API opsional |

### 2.5 Media & storage

| Setting | Key | Dampak |
|---------|-----|--------|
| Max upload MB | `media.max_upload_mb` | 10 |
| Allowed MIME | `media.allowed_mimes` | image/webp,... |
| Storage driver | `media.driver` | `local` \| `r2` |
| R2 bucket / endpoint | `media.r2_*` | Jika pakai Cloudflare R2 |

### 2.6 API & integrasi

| Setting | Key | Dampak |
|---------|-----|--------|
| API rate limit global | `api.rate_limit_per_min` | 300 |
| Webhook signing secret | `webhook.secret` | HMAC outbound |
| Turnstile site key | `turnstile.site_key` | Form publik |

### 2.7 Email (opsional fase 2)

| Setting | Key |
|---------|-----|
| SMTP host / port | `email.smtp_*` |
| From address | `email.from` |

---

## 3. Penyimpanan Konfigurasi

### 3.1 Tabel `system_settings`

```sql
CREATE TABLE system_settings (
  key         TEXT PRIMARY KEY,
  value       JSONB NOT NULL,
  group_name  TEXT NOT NULL,  -- app, worker, cache, media, security
  updated_by  BIGINT REFERENCES users(id),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

| Dampak | |
|--------|--|
| Satu row per key | Update tanpa reload binary |
| JSONB | Tipe fleksibel (number, bool, array) |
| Cache di Go | Load semua settings di startup + refresh tiap 60s atau SET NOTIFY |

**Jangan** simpan secret plain di DB tanpa enkripsi — gunakan:

- Env var untuk secret utama (`DATABASE_URL`, `SESSION_SECRET`)
- `system_settings` hanya untuk non-secret atau secret terenkripsi (AES dengan master key di env)

---

## 4. UI Admin (HTMX)

```
/admin/setup/
├── host/          → subdomain produk
├── backend/       → form grup: Umum, Worker, Cache, Media, API
├── meta/          → meta global produk
├── security/      → rate limit, session TTL, password policy
└── notifications/ → webhook platform
```

Sidebar **Setup** hanya visible untuk `role = super_admin`.

---

## 5. API

| Method | Path | Deskripsi |
|--------|------|-----------|
| GET | `/api/admin/setup/settings` | List by group (masked secrets) |
| PATCH | `/api/admin/setup/settings` | Update batch key-value |
| GET | `/api/admin/setup/health` | DB, disk, queue depth, version |
| POST | `/api/admin/setup/cache/purge` | Purge cache global |

Middleware: **`RequireSuperAdmin`** pada semua `/api/admin/setup/*`.

---

## 6. Skenario & Dampak

| Skenario | Salah konfigurasi | Dampak |
|----------|-------------------|--------|
| `worker.concurrency = 20` di mini CPU | CPU 100%, timeout | Set cap di UI + validasi max 4 |
| `page_size_max = 10000` | OOM pada list | Hard cap 100 di API |
| Maintenance on | Lupa matikan | Publik down — banner di admin |
| Secret di DB plain | DB bocor | Enkripsi / env only |

---

## 7. Relasi dengan Setup Lain

| Dokumen | Isi |
|---------|-----|
| [09](./09-model-domain-host-dan-subdomain.md) | Setup Host |
| [14](./14-setup-meta-dan-seo.md) | Setup Meta |
| [12](./12-autentikasi-dan-login-aman.md) | Setup Security overlap |
| [11](./11-rbac-dan-permission-share.md) | Siapa boleh akses Setup |

---

## 8. Roadmap

| Fase | Item |
|------|------|
| MVP | `system_settings` + UI backend umum + worker + cache |
| Fase 2 | R2 media, email, health dashboard grafik |
| Fase 3 | Import/export settings, staging preview |
