# Frontend-Ui-Admin — Admin Panel HTMX

Panel admin **Seosementara CMS** di Cloudflare Pages (`seosementara.org/admin/*`).

## First boot vs operasional

| Fase | URL | Doc |
|------|-----|-----|
| **Setup infra pertama kali** | GitHub Pages onboarding | [Plan/29-frontend-admin-dan-onboarding.md](../Plan/29-frontend-admin-dan-onboarding.md), [Plan/28](../Plan/28-platform-github-workers.md) |
| **Operasional sehari-hari** | Cloudflare Pages admin | Dokumen ini |

Wizard bootstrap ada di **`Frontend-Onboarding/`** (GitHub Pages). File `admin/bootstrap.html` hanya redirect ke URL onboarding.

Saat infra belum selesai, admin menampilkan banner + link ke GitHub Pages onboarding.

## Fase 0 (mock)

- Belum ada database / Tunnel
- `SSEO.apiMode = 'mock'` di `public/assets/js/config.js`
- Data demo via `mock-api/` + partials HTMX

## Preview lokal (opsional dev)

```bash
npx --yes serve public -p 3000
```

Hanya untuk dev pekerja — **bukan** alur produksi atau onboarding.

## Struktur

```
public/
├── admin/           # Halaman CMS (tanpa bootstrap wizard)
├── assets/          # CSS, JS, config.js
├── partials/        # Sidebar, topbar (HTMX)
└── mock-api/        # Hapus saat apiMode=live
```

## Mode API

Edit `public/assets/js/config.js`:

- `mock` — fase 0 (default)
- `live` — setelah backend Go + Tunnel siap (same-origin `/api/admin/*`)
