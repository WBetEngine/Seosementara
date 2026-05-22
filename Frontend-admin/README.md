# Frontend Admin — UI Prototype

Panel admin **Seosementara CMS** (Plan/27) — **HTML + HTMX + CSS statis**. Backend Go disambungkan besok via `/api/admin/*`.

## Jalankan lokal

```bash
cd Frontend-admin
npm run dev
```

Buka:

- Login (prototype): http://localhost:8788/admin/login.html
- Dashboard: http://localhost:8788/admin/index.html

Atau Cloudflare Pages dev:

```bash
npm run pages:dev
```

## Struktur

```
public/
├── admin/
│   ├── index.html          # Shell: sidebar, topbar, #main, #app-drawer
│   ├── login.html
│   └── _partials/          # Fragment HTMX (mock data)
│       ├── pixel/          # Port templates/pixel/facebook (7 tab)
│       ├── domain-*.html   # Tab domain
│       └── seo-*, pages, … # Semua modul menu
├── static/
│   ├── css/admin.css
│   └── js/admin-shell.js
```

## Fitur prototype (v2)

- Sidebar profesional + **ikon** per menu (SVG inline)
- **Domain:** satu menu → `/admin/domain` dengan **tab panel** (saya / dibagikan / tambah / semua SA)
- **Plugins terpisah:** Shortlink dan Pixel Hub (siap plugin tambahan)
- Tiga dashboard: Admin, Domain, Global (SA)
- Drawer universal `#app-drawer` — domain (5 tab), post, user, CF token, shortlink
- Topbar terang, breadcrumb, stat cards
- Responsif: sidebar drawer di mobile
- Settings: subnav + list + drawer edit

## Pixel Facebook — dipakai di prototype

| Sumber Go (backend) | Prototype statis (Workers) |
|---------------------|----------------------------|
| `templates/pixel/facebook/*.html` | `public/admin/_partials/pixel/*.html` |

Menu **Pixel Hub** → 7 tab (Overview … Analytics) = port konten template ke shell admin v2.  
Saat Backend aktif, handler tetap bisa render template Go; UI final satu shell (`index.html`).

Lihat `templates/pixel/facebook/README.md`.

## Deploy ke Cloudflare Pages (dari GitHub `main`)

**Tidak ada build** untuk prototype ini — deploy langsung folder `public/`.

Push ke `main` **otomatis deploy hanya setelah** proyek Pages di-connect ke repo GitHub (sekali).

Panduan lengkap: **[DEPLOY-CLOUDFLARE-PAGES.md](./DEPLOY-CLOUDFLARE-PAGES.md)**

Ringkas setting Cloudflare:

| Setting | Nilai |
|---------|--------|
| Root directory | `Frontend-admin` |
| Build command | *(kosong)* |
| Build output | `public` |
| Branch | `main` |
