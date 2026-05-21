# 05 — Admin Panel (HTMX di `/admin/`)

> URL produksi: `https://seosementara.org/admin/` — model domain: [09](./09-model-domain-host-dan-subdomain.md)

## 1. Peran

Panel admin adalah antarmuka **banyak pekerja** yang mengelola:

- **Ribuan domain portfolio** (record di database)
- Konten, SEO, media, job batch per domain
- **Host & subdomain produk** (`/admin/setup/host`)

Logika bisnis di backend Go — HTMX hanya memanggil endpoint **sama origin** (`/api/admin/*`).

## 2. Stack UI

| Komponen | Pilihan |
|----------|---------|
| Markup | HTML5 + partial templates |
| Interaktivitas | **HTMX** |
| Styling | CSS ringan |
| URL base | `/admin/` (prefix wajib) |
| Hosting | Dilayani **origin mini CPU** (bukan hostname terpisah) |
| Sumber repo | Folder `Frontend-admin/` |

## 3. Routing Admin

| Path | Halaman |
|------|---------|
| `/admin/login` | Login pekerja |
| `/admin/` | Dashboard |
| `/admin/sites` | Daftar domain portfolio (paginated) |
| `/admin/posts` | Konten domain aktif |
| `/admin/setup/host` | Konfigurasi host & subdomain |
| `/admin/users` | Manajemen pekerja |

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
| Assign pekerja | Admin assign subset domain ke user |

## 6. Skala Banyak Pekerja

| Kebutuhan | Implementasi |
|-----------|--------------|
| Login simultan | Session per user, Redis/store session |
| RBAC | Middleware + hide menu sidebar |
| Audit | Log siapa ubah domain X |
| Scope | User hanya lihat domain yang di-assign |

## 7. Layout

```
https://seosementara.org/admin/
┌─────────────────────────────────────────────┐
│ Topbar: logo, site switcher (portfolio), user│
├──────────┬──────────────────────────────────┤
│ Sidebar  │ #main (HTMX target)              │
│ Dashboard│                                  │
│ Situs    │                                  │
│ Konten   │                                  │
│ Setup    │                                  │
│  └ Host  │                                  │
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
Frontend-admin/
├── templates/
│   ├── layouts/admin.html
│   ├── pages/
│   │   ├── dashboard.html
│   │   ├── sites/
│   │   └── setup/
│   │       └── host.html
│   └── partials/
├── static/css/
└── static/js/htmx.min.js
```

## 10. Dokumen Terkait

- Menu lengkap → [03-menu-dan-modul-cms.md](./03-menu-dan-modul-cms.md)
- Model domain → [09](./09-model-domain-host-dan-subdomain.md)
- API → [07-api-dan-integrasi.md](./07-api-dan-integrasi.md)
