# Frontend Admin — Seosementara CMS

Panel admin **HTML + HTMX + CSS statis**. Cloudflare Workers serve `public/`; API backend di `https://api.apidevel.org`.

## Deploy

```bash
cd Frontend-admin
npm run deploy   # manual
```

Production: GitHub Actions **Deploy Admin UI** (generate `admin-config.js` dari Secret, lalu `wrangler deploy`).

## Dev lokal

```bash
npm run dev
# http://localhost:8788/admin/index.html
```

## Struktur

```
public/
├── admin/
│   ├── index.html
│   └── _partials/       # Fragment HTMX
├── static/js/
│   ├── htmx.min.js
│   ├── admin-shell.js
│   ├── admin-config.example.js
│   └── admin-config.js  # token kosong di repo; CI isi saat deploy
```

## Mini PC / backend

Lihat `mini-pc/DEPLOY.md`.

## Tunnel

Lihat `SETUP-TUNNEL-APIDEVEL.md`.
