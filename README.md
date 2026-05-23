# Seosementara

CMS untuk operasi massal domain & iklan.

## Deploy (GitHub pusat — tanpa `.env` di mini PC)

| Target | Panduan |
|--------|---------|
| Arsitektur platform | [Plan/28-platform-github-workers.md](Plan/28-platform-github-workers.md) |
| Mini PC (Docker) | [mini-pc/DEPLOY.md](mini-pc/DEPLOY.md) |
| Admin Workers + setup UI | `https://seosementara.seosementara3.workers.dev/admin/` |
| Tunnel API | [Frontend-admin/SETUP-TUNNEL-APIDEVEL.md](Frontend-admin/SETUP-TUNNEL-APIDEVEL.md) |

**Setup operator:** Settings → **Infra & GitHub** (DB, encryption) dan **Cloudflare → Koneksi** (Global API Key) — semua lewat admin Workers.

## Dokumentasi Pixel

- [Plan/20](Plan/20-pixel-admin-facebook-tiktok-gads.md) — Pixel Hub
- [Plan/21](Plan/21-pixel-facebook-pro.md) — Facebook Pro
- [Plan/22](Plan/22-pixel-protokol-komunikasi-dan-data.md) — Protokol data
- [Plan/23](Plan/23-meta-conversions-api-kedalaman.md) — Meta CAPI
