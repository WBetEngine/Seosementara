# Frontend-Ui-Admin — Admin Panel HTMX

Panel admin **Seosementara CMS**. Fase 0: UI statis tanpa backend (mock API).

## Kondisi bootstrap

Lihat [Plan/29-bootstrap-admin-ui-pertama-kali.md](../Plan/29-bootstrap-admin-ui-pertama-kali.md):

- Belum ada database
- Tidak pakai `.env`
- Belum ada Cloudflare Tunnel

Semua form disiapkan untuk pengisian dari UI; data demo via mock partials + localStorage.

## Preview lokal

```bash
npx --yes serve public -p 3000
```

Buka: http://localhost:3000/admin/login.html

Login demo: email apa saja + password apa saja → redirect dashboard.

## Struktur

```
public/
├── admin/           # Halaman utama
├── assets/          # CSS, JS
├── partials/        # Sidebar, topbar (HTMX)
└── mock-api/        # Response HTML statis (ganti /api/admin nanti)
```

## Mode API

Edit `public/assets/js/config.js`:

- `mock` — fase 0 (default)
- `live` — setelah backend Go siap
