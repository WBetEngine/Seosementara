# 05 — Admin Panel (HTMX + Cloudflare Pages)

## 1. Peran

Panel admin adalah antarmuka operator untuk mengelola CMS. **Tidak** ada logika bisnis kritis di browser — HTMX memanggil API Golang di mini CPU dan menukar fragmen HTML.

## 2. Stack UI

| Komponen | Pilihan |
|----------|---------|
| Markup | HTML5 semantic |
| Interaktivitas | **HTMX** (hx-get, hx-post, hx-swap, hx-trigger) |
| Styling | CSS vanilla atau utility ringan (tanpa framework berat) |
| Icons | SVG inline atau sprite |
| Build | Opsional: `templ` (Go) generate HTML saat CI, deploy static ke Pages |
| Hosting | **Cloudflare Pages** |

**Tidak memakai:** React/Vue/Angular penuh — menjaga bundle kecil dan cocok edge.

## 3. Pola HTMX Utama

### 3.1 List dengan pagination

```html
<div id="post-list"
     hx-get="/api/admin/posts?site_id=1&page=1"
     hx-trigger="load"
     hx-swap="innerHTML">
  <!-- skeleton -->
</div>
```

Backend mengembalikan **partial HTML** (bukan JSON) untuk admin, atau JSON + client template — **keputusan:** partial HTML dari Go `html/template` direkomendasikan agar satu sumber template.

### 3.2 Form inline edit

```html
<form hx-post="/api/admin/posts/123"
      hx-target="#post-editor"
      hx-swap="outerHTML">
```

### 3.3 Notifikasi & polling job

```html
<div hx-get="/api/admin/jobs/45/status"
     hx-trigger="every 2s"
     hx-swap="innerHTML">
```

Hentikan polling saat status `completed` (HX-Trigger header atau `htmx:afterSwap`).

## 4. Layout Halaman

```
┌─────────────────────────────────────────────┐
│ Topbar: logo, site switcher, user menu      │
├──────────┬──────────────────────────────────┤
│ Sidebar  │ Main content (HTMX swap target)  │
│ (menu    │                                  │
│  03)     │                                  │
└──────────┴──────────────────────────────────┘
```

- **Site switcher:** set cookie/header `X-Site-ID` untuk semua request berikutnya
- Sidebar: render menu sesuai RBAC (hide item tanpa permission)

## 5. Pemetaan Menu → Halaman

| Menu (lihat 03) | Route Pages (contoh) | Partial target |
|-----------------|----------------------|----------------|
| Dashboard | `/admin/` | `#main` |
| Daftar post | `/admin/posts` | `#main` |
| Edit post | `/admin/posts/{id}` | `#editor` |
| Media | `/admin/media` | `#main` |
| SEO bulk | `/admin/seo/bulk` | `#main` |
| Jobs | `/admin/jobs` | `#main` |

Routing: Cloudflare Pages `_redirects` atau Functions ringan jika perlu.

## 6. Autentikasi di Edge

1. Halaman `/admin/login` — form HTMX POST ke `https://api.../api/admin/auth/login`
2. Backend set `Set-Cookie` (Secure, HttpOnly, SameSite=None untuk cross-site Pages→API)
3. HTMX config global:

```html
<body hx-headers='{"X-Requested-With": "XMLHttpRequest"}' ...>
```

4. Request berikutnya `hx-include` cookie otomatis (same-site policy — pertimbangkan subdomain atau Tunnel sama root domain).

**Catatan:** cross-origin Pages → API memerlukan CORS + credentials; idealnya `admin.domain.com` dan `api.domain.com` under satu registrable domain.

## 7. Struktur Folder (Usulan)

```
Frontend-admin/
├── public/
│   ├── index.html
│   ├── css/
│   └── js/htmx.min.js
├── partials/          # jika pre-render
├── _redirects
└── wrangler.toml      # Pages config
```

## 8. Performa di Cloudflare Pages

| Praktik | Alasan |
|---------|--------|
| Asset CSS/JS cache panjang | Immutable hash filename |
| Fragment kecil | Swap cepat, sedikit HTML |
| Debounce search | `hx-trigger="keyup changed delay:300ms"` |
| Tidak fetch ribuan row | Pagination server-side wajib |

## 9. Aksesibilitas & UX

- Focus management setelah swap (`htmx:focus-scroll`)
- Pesan error dari `HX-Trigger` atau fragment alert
- Konfirmasi destructive: `hx-confirm="Yakin hapus?"`

## 10. Environment Pages

| Variable | Contoh |
|----------|--------|
| `API_BASE_URL` | `https://api.seosementara.example` |
| `APP_NAME` | Seosementara Admin |

Set di Cloudflare dashboard → Settings → Environment variables.

## 11. Dokumen Terkait

- Menu lengkap → [03-menu-dan-modul-cms.md](./03-menu-dan-modul-cms.md)
- API → [07-api-dan-integrasi.md](./07-api-dan-integrasi.md)
- Arsitektur → [02-arsitektur-dan-infrastruktur.md](./02-arsitektur-dan-infrastruktur.md)
