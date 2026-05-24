# 05 — Admin Panel (HTMX di `/admin/`)

> URL produksi: `https://seosementara.org/admin/` — model domain: [09](./09-model-domain-host-dan-subdomain.md)  
> **First boot infra:** GitHub Pages onboarding — **bukan** halaman di admin CF Pages ([29](./29-frontend-admin-dan-onboarding.md), [28](./28-platform-github-workers.md))

## 1. Peran

Panel admin adalah antarmuka **banyak pekerja** yang mengelola:

- **Ribuan domain portfolio** (situs native CMS — bukan WordPress)
- Hanya domain **milik sendiri** + yang **di-share** (kecuali Super Admin)
- Konten, SEO, media, job batch per domain
- **Host & subdomain produk** (`/admin/settings/host`)

Logika bisnis di backend Go — HTMX hanya memanggil endpoint **sama origin** (`/api/admin/*`).

## 2. Stack UI

| Komponen | Pilihan |
|----------|---------|
| Markup | HTML5 + partial templates |
| Interaktivitas | **HTMX** |
| Styling | CSS ringan |
| URL base | `/admin/` (prefix wajib) |
| Hosting | **Cloudflare Pages** — env dari Setup Cloudflare [15](./15-setup-cloudflare-integrasi.md) |
| Sumber repo | Folder `Frontend-Ui-Admin/` |

## 3. Routing Admin

| Path | Halaman |
|------|---------|
| `/admin/login` | Login pekerja |
| `/admin/` | Dashboard |
| `/admin/sites` | Domain milik saya + dibagikan (paginated) |
| `/admin/sites/{id}/sharing` | Berbagi akses + daftar undangan pending (owner) |
| `/admin/sites/{id}/transfer-owner` | Transfer ownership (**Super Admin**) |
| `/admin/notifications` | Notifikasi (undangan co-admin, transfer, dll.) |
| `/admin/settings/host` | Subdomain produk (**Super Admin**) |
| `/admin/content/posts` | Konten domain aktif |
| `/admin/users` | Manajemen pekerja |
| `/admin/pixel/` | Pixel Hub overview — [20](./20-pixel-admin-facebook-tiktok-gads.md) |
| `/admin/pixel/events/` | Event catalog (single source → 3 platform) |
| `/admin/pixel/facebook` | Ruang kerja kolaborasi Meta (CAPI, diagnostics) |
| `/admin/pixel/tiktok` | Ruang kerja TikTok |
| `/admin/pixel/gads` | Ruang kerja Google Ads |

Semua link internal memakai prefix `/admin/` — hindari path absolut tanpa prefix.

## 4. Pola HTMX (Sama Origin)

```html
<div id="site-list"
     hx-get="/api/admin/managed-domains?page=1&limit=50"
     hx-trigger="load"
     hx-swap="innerHTML">
</div>
```

Keuntungan same-origin:

- Cookie session tanpa masalah CORS cross-domain
- Tidak perlu `API_BASE_URL` ke host lain

### Site switcher (domain portfolio)

```html
<select name="managed_domain_id"
        hx-post="/api/admin/session/active-domain"
        hx-swap="none">
  <!-- opsi dari search, bukan load 1000 option sekaligus -->
</select>
```

Gunakan **combobox search** (ketik → `hx-get` autocomplete), bukan `<select>` 1000 item.

## 5. Skala Ribuan Domain

| UI | Pola |
|----|------|
| Daftar domain | Pagination + filter + indexed search |
| Bulk action | Pilih filter → konfirmasi → job ID → poll progress |
| Dashboard | Angka agregat dari cache — bukan `COUNT(*)` tiap load |
| Berbagi domain | Owner: langsung aktif; Co-admin: form + status pending |
| Notifikasi owner | Badge + list; tombol setujui/tolak (HTMX) |
| Transfer owner | Form Super Admin; owner lama kehilangan akses sepenuhnya |

## 6. Kepemilikan & Banyak Pekerja

| Kebutuhan | Implementasi |
|-----------|--------------|
| Login simultan | Session per user |
| Isolasi | API filter `owner OR domain_shares` |
| Share | Halaman `/admin/sites/{id}/sharing` |
| Super Admin | Menu "Semua domain" + Setup Host |
| Audit | Log share, ubah owner, edit konten |

## 7. Layout

> **Desain lengkap:** tiga dashboard (Global / Admin / Domain), navigasi berkelompok, Setup submenu, responsif mobile — [27-admin-panel-desain-ui-navigasi.md](./27-admin-panel-desain-ui-navigasi.md). **Bukan** toko online (tanpa Cart/Toko).

```
https://seosementara.org/admin/
┌─────────────────────────────────────────────┐
│ Topbar: logo, site switcher (portfolio), user│
├──────────┬──────────────────────────────────┤
│ Sidebar  │ #main (HTMX target)              │
│ Ringkasan│  (3 dashboard — lihat Plan/27)   │
│ Domain   │                                  │
│ Konten   │                                  │
│ Plugins  │  Shortlink, Pixel Hub            │
│ Settings │  Read/Edit/Write sistem          │
└──────────┴──────────────────────────────────┘
```

## 8. Autentikasi

```http
POST /api/admin/auth/login
→ Set-Cookie session (HttpOnly, Secure, SameSite=Lax)
```

Redirect setelah login: `HX-Redirect: /admin/`

Middleware: semua `/admin/*` kecuali login → cek session.

## 9. Struktur Folder

```
Frontend-Ui-Admin/
└── public/
    ├── admin/              # halaman CMS (HTML statis / HTMX)
    ├── assets/css/
    ├── assets/js/
    ├── partials/           # sidebar, topbar
    └── mock-api/           # fase mock — hapus saat apiMode=live
```

Onboarding wizard **tidak** ada di folder ini — lihat `Frontend-Onboarding/` ([29](./29-frontend-admin-dan-onboarding.md)).

## 10. Dokumen Terkait

- Bootstrap & onboarding → [29-frontend-admin-dan-onboarding.md](./29-frontend-admin-dan-onboarding.md)
- Kontrak HTMX & komponen → [17-kontrak-htmx-dan-komponen-ui.md](./17-kontrak-htmx-dan-komponen-ui.md)
- Menu lengkap → [03-menu-dan-modul-cms.md](./03-menu-dan-modul-cms.md)
- Model domain → [09](./09-model-domain-host-dan-subdomain.md)
- API → [07-api-dan-integrasi.md](./07-api-dan-integrasi.md)
