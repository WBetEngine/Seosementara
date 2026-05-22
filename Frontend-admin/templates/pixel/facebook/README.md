# Template Pixel Facebook (Go)

Template HTML untuk **Backend Go** (`AdminPixelFacebook` handler).

## Dua jalur UI

| Deploy | File | Cara render |
|--------|------|-------------|
| **Cloudflare Workers** (statis) | `public/admin/_partials/pixel/*.html` | HTMX → `#main` di shell `index.html` |
| **Backend + Tunnel** | `templates/pixel/facebook/*.html` | `html/template` full page + `pixel.css` |

Konten **sama** (Overview, Setup, Connection, Domains, Diagnostics, Events, Analytics).  
Prototype statis memakai desain admin v2; template Go memakai layout gelap `pixel_head` sampai backend digabung ke shell universal.

## Integrasi ke shell (fase berikutnya)

1. Pecah `pixel_head` / `pixel_foot` — hanya body masuk `#main` atau `#pixel-tab-panel`
2. Handler render fragment HTMX, bukan halaman penuh
3. `pixel.css` opsional; prefer token `admin.css`

## API terkait

Lihat Plan/22 § API admin pixel Facebook.
