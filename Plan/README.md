# Dokumentasi Perencanaan — Seosementara CMS

Dokumen ini menguraikan perencanaan sistem CMS untuk proyek **Seosementara**. Pembahasan sengaja dipecah ke beberapa file agar mudah dikembangkan, direview, dan diubah tanpa satu file raksasa.

## Cara Membaca

Baca berurutan jika Anda baru memulai. Jika fokus ke satu lapisan, buka file yang relevan saja.

| No. | File | Isi |
|-----|------|-----|
| 01 | [01-visi-dan-gambaran-sistem-cms.md](./01-visi-dan-gambaran-sistem-cms.md) | Visi produk, ruang lingkup CMS, prinsip desain |
| 02 | [02-arsitektur-dan-infrastruktur.md](./02-arsitektur-dan-infrastruktur.md) | Satu domain produk, mini CPU, Cloudflare, routing Host+Path |
| 03 | [03-menu-dan-modul-cms.md](./03-menu-dan-modul-cms.md) | Menu admin, Setup → Host, skala ribuan domain & pekerja |
| 04 | [04-backend-golang.md](./04-backend-golang.md) | Backend Go: struktur, API, database, performa |
| 05 | [05-admin-panel-htmx.md](./05-admin-panel-htmx.md) | Panel admin HTMX di `/admin/` |
| 06 | [06-frontend-users-htmx.md](./06-frontend-users-htmx.md) | Frontend publik apex + subdomain HTMX |
| 07 | [07-api-dan-integrasi.md](./07-api-dan-integrasi.md) | Kontrak API, auth, CORS, cache |
| 08 | [08-roadmap-implementasi.md](./08-roadmap-implementasi.md) | Fase implementasi dan prioritas |
| 09 | [09-model-domain-host-dan-subdomain.md](./09-model-domain-host-dan-subdomain.md) | **Model domain:** apex, `/admin/`, subdomain, ribuan domain portfolio |
| 10 | [10-database-postgresql.md](./10-database-postgresql.md) | **PostgreSQL:** schema, index, skenario, dampak performa |

## Ringkasan Stack

| Lapisan | Teknologi | URL / Hosting |
|---------|-----------|---------------|
| Database | **PostgreSQL** | Mini CPU; pool via PgBouncer |
| Backend API | **Golang** | Mini CPU; `seosementara.org/api/*` |
| Admin Panel | **HTMX** | `seosementara.org/admin/*` (sama origin) |
| Frontend publik | **HTMX** | `seosementara.org/` + subdomain (`bola.`, `cdn.`, …) |
| Domain portfolio | Data di DB | Ribuan domain **dikelola** di admin — bukan hostname frontend terpisah |
| Edge | Cloudflare | DNS wildcard, SSL, cache, Tunnel ke mini CPU |

## Konteks Skala (Prinsip Tetap)

Proyek ini ditargetkan untuk operasi massal (banyak domain, volume konten besar). Semua keputusan arsitektur harus mengutamakan:

- Query database terarah dan terbatas
- Operasi berat di-batch, bukan loop tanpa batas
- Cache dengan invalidasi yang jelas
- Hindari timeout dan beban berlebihan di server terbatas (mini CPU)

Detail teknis performa ada di file **04**, **07**, dan **10**.

## Status Dokumen

| Versi | Tanggal | Catatan |
|-------|---------|---------|
| 0.1 | 2026-05-21 | Draft awal struktur perencanaan |
| 0.2 | 2026-05-21 | Revisi model domain: `/admin/`, subdomain, ribuan domain portfolio |
| 0.3 | 2026-05-21 | Native CMS (bukan WP), ownership + share, subdomain oleh Super Admin |
| 0.4 | 2026-05-21 | Rencana database PostgreSQL (schema, index, skenario) |
