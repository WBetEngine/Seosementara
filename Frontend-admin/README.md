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
├── static/
│   ├── css/admin.css
│   └── js/admin-shell.js
```

## Fitur prototype

- Sidebar gelap + menu berkelompok (Ringkasan, Domain, Konten, SEO, Plugins, Settings)
- Tiga dashboard: Admin, Domain, Global (SA)
- Drawer universal `#app-drawer` — domain (5 tab), post, user, CF token, shortlink
- Responsif: sidebar drawer di mobile, panel kanan full width
- Settings: subnav + list + drawer edit

## Pixel Facebook (kode lama)

Template Go lama masih di `templates/pixel/facebook/` — akan diintegrasikan ke shell baru.

## Deploy Pages (free)

Output folder = `public`. Hubungkan repo ke Cloudflare Pages; route `seosementara.org/admin/*`.
