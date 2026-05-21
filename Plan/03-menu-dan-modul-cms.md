# 03 — Menu dan Modul CMS (Admin Panel)

Dokumen ini mendefinisikan struktur menu **Admin Panel**. Setiap item menu memetakan ke modul backend dan hak akses. Detail teknis UI ada di [05-admin-panel-htmx.md](./05-admin-panel-htmx.md).

## 1. Struktur Navigasi Utama

```
Dashboard
├── Ringkasan global
├── Aktivitas terbaru
└── Peringatan sistem

Situs (Domain Portfolio)
├── Daftar domain saya
├── Domain dibagikan ke saya
├── Tambah domain baru
├── Detail domain
│   ├── Pengaturan per domain
│   ├── Berbagi akses (share ke admin lain)
│   ├── Undangan menunggu persetujuan (owner)
│   └── DNS & catatan operasi
└── (Super Admin) Semua domain

Konten
├── Post
│   ├── Semua post
│   ├── Tambah baru
│   ├── Draft
│   ├── Terjadwal
│   └── Sampah
├── Halaman (Page)
├── Kategori
├── Tag
└── Template blok (opsional fase 2)

Media
├── Perpustakaan media
├── Upload
└── Pengaturan optimasi gambar

SEO
├── Meta global per situs
├── Meta per konten (bulk editor)
├── Sitemap & robots
├── Redirect manager
└── Schema / structured data

Pengguna & Akses (ringkas — kelola penuh di Setup Backend)
├── Profil saya / sesi saya
└── Log aktivitas (jika berhak)

Tools / Shortlink → [19-modul-url-shortlink.md](./19-modul-url-shortlink.md)
├── Daftar shortlink (auto + manual)
├── Buat shortlink manual (form di admin — bukan otomatis)
└── Analitik klik + Cloudflare

Pixel → [20-pixel-admin-facebook-tiktok-gads.md](./20-pixel-admin-facebook-tiktok-gads.md)
├── /admin/pixel/              → Ringkasan kanal
├── /admin/pixel/facebook/     → Setup, events, domains, analytics
├── /admin/pixel/tiktok/       → Setup, events, domains, analytics
└── /admin/pixel/gads/         → Setup, events, domains, analytics (Google Ads)

Operasi Massal
├── Batch publish / unpublish
├── Bulk update meta SEO
├── Import / export konten
└── Sinkronisasi (jika ada integrasi eksternal)

Jobs / Antrian
├── Job berjalan
├── Riwayat job
└── Gagal & retry

Laporan
├── Statistik konten per situs
├── Status publish
└── Ringkasan error API

Setup (Super Admin / role dengan permission setup.*)
├── Backend → [13](./13-setup-backend-dan-sistem.md) **PUSAT**
│   ├── Ringkasan & health
│   ├── RBAC (peran admin, pengguna, permission sistem)
│   ├── Autentikasi & login aman → [12](./12-autentikasi-dan-login-aman.md)
│   ├── Rate limit (aplikasi + selaras Cloudflare)
│   ├── Operasional (DB, worker, cache, maintenance)
│   ├── Media & storage
│   └── API / webhook / Turnstile
├── Cloudflare → [15](./15-setup-cloudflare-integrasi.md)
├── Host / Subdomain + meta host → [09], [14]
├── Meta global → [14](./14-setup-meta-dan-seo.md)
└── Notifikasi platform

Bantuan
├── Dokumentasi internal
└── Versi sistem
```

## 2. Deskripsi Menu per Modul

### 2.1 Dashboard

| Submenu | Fungsi |
|---------|--------|
| Ringkasan global | Jumlah situs, post published/draft, job aktif |
| Aktivitas terbaru | 20 event terakhir (publish, edit, login) |
| Peringatan sistem | Disk hampir penuh, job gagal, API health |

**Query:** agregat ter-cache (transient 5 menit), bukan hitung penuh tabel tanpa filter.

---

### 2.2 Situs (Domain Portfolio — Ribuan)

Modul ini mengelola **domain yang dioperasikan** (ribuan), bukan hostname UI produk.

| Submenu | Fungsi |
|---------|--------|
| Daftar domain saya | Pagination, search — hanya `owner_user_id = saya` |
| Domain dibagikan | Domain yang user lain share ke saya (`domain_shares`) |
| Tambah domain | Buat `managed_domain` + **auto shortlink** `url.../{slug}` [19](./19-modul-url-shortlink.md) |
| Detail → Berbagi akses | Preset read only / edit / co-admin + **checklist permission** → [11](./11-rbac-dan-permission-share.md) |
| Undangan pending | Owner: setujui / tolak undangan dari co-admin |
| Transfer ownership | **Super Admin** — pindah pemilik domain |
| Pengaturan per domain | SEO default, status, catatan — **bukan WordPress** |
| Semua domain | Hanya **Super Admin** — list global |

**Skala:** **1000+ domain** total di sistem; setiap pekerja hanya memuat subset milik + shared.

**Kepemilikan:** pekerja **tidak** melihat domain pekerja lain kecuali di-share. Lihat [09](./09-model-domain-host-dan-subdomain.md) §7.

---

### 2.2b Setup → Host (Domain Produk & Subdomain)

| Submenu | Fungsi |
|---------|--------|
| Daftar host | `seosementara.org`, `bola.seosementara.org`, … |
| Tambah host | Hostname baru + template (subdomain dinamis) |
| Edit / ganti | Ubah template, hostname, maintenance, nonaktifkan |
| Mapping template | Pilih UI HTMX untuk host tersebut |
| Panduan DNS | Wildcard `*.seosementara.org` → origin |

**URL admin:** `/admin/setup/host`

Subdomain **bisa ditambah dan diganti** sewaktu-waktu — keputusan **Super Admin** saja.

**Peran:** **Super Admin eksklusif** — pekerja biasa tidak punya menu ini.

---

### 2.3 Konten

#### Post

| Submenu | Fungsi |
|---------|--------|
| Semua post | List dengan filter `site_id`, status, tanggal |
| Tambah baru | Editor + sidebar SEO |
| Draft / Terjadwal / Sampah | View filter status |

#### Halaman (Page)

Mirip post, tanpa kategori blog (hierarki parent opsional).

#### Kategori & Tag

CRUD taxonomy per situs; hindari load semua term sekaligus — tree lazy-load jika banyak.

---

### 2.4 Media

| Submenu | Fungsi |
|---------|--------|
| Perpustakaan | Grid/list paginated per `site_id` |
| Upload | Chunk upload untuk file besar; validasi MIME |
| Optimasi | Kualitas WebP, max dimension — setting per situs |

---

### 2.5 SEO

| Submenu | Fungsi |
|---------|--------|
| Meta global | Title suffix, default description, OG image |
| Bulk meta editor | Spreadsheet-like HTMX; update batch via job |
| Sitemap & robots | Generate XML per situs; invalidate on publish |
| Redirect | 301/302 list, match path |
| Schema | JSON-LD template per tipe konten |

---

### 2.6 Pengguna & Akses

| Submenu | Fungsi |
|---------|--------|
| Pengguna admin | CRUD user, reset password, 2FA (fase 2) |
| Peran | Super Admin, Site Manager, Editor, SEO Specialist, Viewer |
| Log aktivitas | Audit trail: siapa mengubah apa |

---

### 2.7 Operasi Massal

| Submenu | Fungsi |
|---------|--------|
| Batch publish | Pilih filter → enqueue job |
| Bulk SEO | Update field meta terpilih |
| Import / export | CSV/JSON terbatas; tidak unbounded parse |
| Sinkronisasi | Integrasi eksternal masa depan (opsional) |

**Wajib:** semua operasi > N item (mis. 50) masuk antrian job, bukan loop sinkron di request HTTP.

---

### 2.8 Jobs / Antrian

| Submenu | Fungsi |
|---------|--------|
| Berjalan | Progress bar, ETA perkiraan |
| Riwayat | 30 hari terakhir, paginated |
| Gagal & retry | Detail error per chunk, tombol retry |

---

### 2.9 Laporan

| Submenu | Fungsi |
|---------|--------|
| Statistik konten | Count per status per situs (cached) |
| Status publish | Kalender / list publish hari ini |
| Error API | Agregat 4xx/5xx dari log |

---

### 2.10 Pengaturan

| Submenu | Fungsi |
|---------|--------|
| Umum | Timezone, format tanggal |
| API keys | Token untuk integrasi eksternal |
| Cache | TTL default, tombol purge per situs |
| Notifikasi | Webhook Slack/Telegram saat job selesai |
| Maintenance | Mode maintenance per situs atau global |

---

## 3. RBAC (Ringkas)

| Lapisan | Kelola di admin | Dokumen |
|---------|-----------------|---------|
| Sistem (role admin, setup) | `/admin/setup/backend/rbac/` | [13](./13-setup-backend-dan-sistem.md) §3, [11](./11-rbac-dan-permission-share.md) |
| Domain (share checklist) | `/admin/sites/{id}/sharing` | [11](./11-rbac-dan-permission-share.md) |
| Login & keamanan | `/admin/setup/backend/autentikasi/` | [12](./12-autentikasi-dan-login-aman.md) |
| Rate limit | `/admin/setup/backend/rate-limit/` | [13](./13-setup-backend-dan-sistem.md) §5 |

## 4. Menu Frontend Publik (Bukan Menu Admin)

Navigasi pengunjung di **`seosementara.org`** dan subdomain — dikonfigurasi per **host** di Setup → Host, bukan per domain portfolio:

| Host | Navigasi (contoh) |
|------|-------------------|
| Apex | Beranda, Blog, Dokumentasi, Kontak |
| `bola.*` | Jadwal, Liga, Statistik |
| `url.*` | Buat link, Statistik klik |
| `cdn.*` | Browse aset (jika publik) |

Detail di [06-frontend-users-htmx.md](./06-frontend-users-htmx.md) dan [09](./09-model-domain-host-dan-subdomain.md).

## 4b. Kepemilikan & Banyak Pekerja

| Fitur | Keterangan |
|-------|------------|
| Ownership | Setiap domain punya owner; bukan WordPress |
| Isolasi | List/query filter owner + shared only |
| Share owner | Undangan langsung aktif |
| Share co-admin | Undangan **pending** sampai owner **setujui** |
| Notifikasi | Owner dapat alert undangan co-admin |
| Transfer owner | Hanya Super Admin |
| Audit log | Share, approve/reject, transfer ownership |
| Site switcher | Domain milik + dibagikan (bukan pending) |

Detail alur: [09](./09-model-domain-host-dan-subdomain.md) §7.3–7.6.

## 5. Pemetaan Menu → Endpoint API (Ringkas)

| Menu | Prefix API (admin) |
|------|-------------------|
| Dashboard | `GET /api/admin/dashboard` |
| Situs | `/managed-domains`, `/shares`, `/share-invitations` |
| Notifikasi | `/api/admin/notifications` |
| Setup Host | `/api/admin/hosts` (Super Admin) |
| Konten | `/api/admin/posts`, `/pages`, `/taxonomies` |
| Media | `/api/admin/media` |
| SEO | `/api/admin/seo` |
| Pengguna | `/api/admin/users` |
| Jobs | `/api/admin/jobs` |
| Pengaturan | `/api/admin/settings` |

Spesifikasi lengkap di [07-api-dan-integrasi.md](./07-api-dan-integrasi.md).

## 6. Prioritas Implementasi Menu

| Fase | Menu |
|------|------|
| MVP | Dashboard, Situs, Konten (post/page), Media, SEO dasar, Pengguna dasar |
| Fase 2 | Operasi massal, Jobs, Redirect, Schema |
| Fase 3 | Laporan lanjutan, Import/export, 2FA |

Lihat [08-roadmap-implementasi.md](./08-roadmap-implementasi.md).
