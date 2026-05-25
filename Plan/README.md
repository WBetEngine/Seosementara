# Dokumentasi Perencanaan — Seosementara CMS

Dokumen ini menguraikan perencanaan sistem CMS untuk proyek **Seosementara**. Pembahasan sengaja dipecah ke beberapa file agar mudah dikembangkan, direview, dan diubah tanpa satu file raksasa.

## Cara Membaca

Baca berurutan jika Anda baru memulai. Jika fokus ke satu lapisan, buka file yang relevan saja.

| No. | File | Isi |
|-----|------|-----|
| 01 | [01-visi-dan-gambaran-sistem-cms.md](./01-visi-dan-gambaran-sistem-cms.md) | Visi produk, ruang lingkup CMS, prinsip desain |
| 02 | [02-arsitektur-dan-infrastruktur.md](./02-arsitektur-dan-infrastruktur.md) | Satu domain produk, mini PC Docker, Cloudflare, routing Host+Path |
| 28 | [28-platform-github-workers.md](./28-platform-github-workers.md) | **Bootstrap:** GitHub Pages onboarding, Workers Platform API, Docker, tanpa `.env` |
| 29 | [29-frontend-admin-dan-onboarding.md](./29-frontend-admin-dan-onboarding.md) | **Frontend:** pemisahan Onboarding (GH Pages) vs Admin (CF Pages) |
| 03 | [03-menu-dan-modul-cms.md](./03-menu-dan-modul-cms.md) | Menu admin, Settings → Host, skala ribuan domain & pekerja |
| 04 | [04-backend-golang.md](./04-backend-golang.md) | Backend Go: struktur, API, database, performa |
| 05 | [05-admin-panel-htmx.md](./05-admin-panel-htmx.md) | Panel admin HTMX di `/admin/` (Cloudflare Pages) |
| 06 | [06-frontend-users-htmx.md](./06-frontend-users-htmx.md) | Frontend publik apex + subdomain HTMX (`Frontend-Publik/`) |
| 07 | [07-api-dan-integrasi.md](./07-api-dan-integrasi.md) | Kontrak API, auth, CORS, cache |
| 08 | [08-roadmap-implementasi.md](./08-roadmap-implementasi.md) | Fase implementasi dan prioritas |
| 09 | [09-model-domain-host-dan-subdomain.md](./09-model-domain-host-dan-subdomain.md) | **Model domain:** apex, `/admin/`, subdomain, ribuan domain portfolio |
| 10 | [10-database-postgresql.md](./10-database-postgresql.md) | **PostgreSQL:** schema, index, skenario, dampak performa |
| 11 | [11-rbac-dan-permission-share.md](./11-rbac-dan-permission-share.md) | **RBAC** + share read/edit + checklist permission |
| 12 | [12-autentikasi-dan-login-aman.md](./12-autentikasi-dan-login-aman.md) | Login aman, session, rate limit, CSRF |
| 13 | [13-setup-backend-dan-sistem.md](./13-setup-backend-dan-sistem.md) | **Settings Backend:** RBAC, auth, rate limit, operasional |
| 14 | [14-setup-meta-dan-seo.md](./14-setup-meta-dan-seo.md) | Meta: global, subdomain, domain, halaman |
| 15 | [15-setup-cloudflare-integrasi.md](./15-setup-cloudflare-integrasi.md) | **Cloudflare:** Tunnel, DNS, Workers Secrets, domain di DB |
| 16 | [16-deploy-dan-lingkungan.md](./16-deploy-dan-lingkungan.md) | **Deploy:** staging/prod, CI/CD Docker, rollback |
| 17 | [17-kontrak-htmx-dan-komponen-ui.md](./17-kontrak-htmx-dan-komponen-ui.md) | **HTMX:** header, error, swap, komponen admin & publik |
| 18 | [18-bisnis-subdomain-dan-modul.md](./18-bisnis-subdomain-dan-modul.md) | **Subdomain:** bola, cdn, url, ads, comments, review |
| 19 | [19-modul-url-shortlink.md](./19-modul-url-shortlink.md) | **Shortlink:** auto per domain, manual, tracking + CF |
| 20 | [20-pixel-admin-facebook-tiktok-gads.md](./20-pixel-admin-facebook-tiktok-gads.md) | **Pixel Hub:** kolaborasi FB/TikTok/GAds, CAPI, first-party, event catalog |
| 21 | [21-pixel-facebook-pro.md](./21-pixel-facebook-pro.md) | **Pixel Facebook Pro:** fitur profesional, kolaborasi Meta, data per tab |
| 22 | [22-pixel-protokol-komunikasi-dan-data.md](./22-pixel-protokol-komunikasi-dan-data.md) | **Protokol & data lengkap:** canonical event, pipeline, schema DB |
| 23 | [23-meta-conversions-api-kedalaman.md](./23-meta-conversions-api-kedalaman.md) | **Meta CAPI:** auth, payload, EMQ, dedup, hybrid, multi-pixel |
| 24 | [24-meta-akun-bm-pixel-dan-optimasi-iklan.md](./24-meta-akun-bm-pixel-dan-optimasi-iklan.md) | **Realita BM/personal/Fanpage**, SOP BM putus, fitur optimasi biaya |
| 25 | [25-pixel-data-lengkap-emq.md](./25-pixel-data-lengkap-emq.md) | **Data selaras Meta 2026:** tabel resmi `user_data`, hash, wajib website |
| 26 | [26-meta-sumber-resmi-pixel-capi.md](./26-meta-sumber-resmi-pixel-capi.md) | **Peta 3 URL resmi Meta:** standard events, Pixel dev, CAPI business |
| 27 | [27-admin-panel-desain-ui-navigasi.md](./27-admin-panel-desain-ui-navigasi.md) | **Desain admin UI:** 3 dashboard, domain drawer, Plugins, **Settings**, SEO domain-panel, responsif |

## Ringkasan Stack

| Lapisan | Teknologi | Folder repo | Hosting |
|---------|-----------|-------------|---------|
| **Onboarding (first boot)** | HTML + HTMX | `Frontend-Onboarding/` | **GitHub Pages** |
| Database | **PostgreSQL** | `Backend/` (migrasi) | Mini PC Docker |
| Backend API | **Golang** | `Backend/` | Mini PC Docker; `seosementara.org/api/*` via Tunnel |
| Admin Panel | **HTMX** | `Frontend-Ui-Admin/` | **Cloudflare Pages** — `seosementara.org/admin/*` |
| Frontend publik | **HTMX** | `Frontend-Publik/` | **Cloudflare Pages** — apex + subdomain |
| Platform API | **Cloudflare Workers** | Workers project | `/admin/api/platform/*` (bootstrap) |
| Domain portfolio | Data di DB | — | Ribuan domain **dikelola** di admin — bukan hostname frontend terpisah |
| Edge | Cloudflare | — | Pages, Tunnel, DNS — Settings admin [15] |

**First boot:** operator buka **GitHub Pages onboarding** → Workers Platform API → deploy infra → redirect ke admin CF Pages ([28](./28-platform-github-workers.md)).

## Konteks Skala (Prinsip Tetap)

Proyek ini ditargetkan untuk operasi massal (banyak domain, volume konten besar). Semua keputusan arsitektur harus mengutamakan:

- Query database terarah dan terbatas
- Operasi berat di-batch, bukan loop tanpa batas
- Cache dengan invalidasi yang jelas
- Hindari timeout dan beban berlebihan di server terbatas (mini CPU)

Detail teknis: performa **04**, **07**, **10**; keamanan & hak akses **11**, **12**.

## Status Dokumen

| Versi | Tanggal | Catatan |
|-------|---------|---------|
| 0.1 | 2026-05-21 | Draft awal struktur perencanaan |
| 0.2 | 2026-05-21 | Revisi model domain: `/admin/`, subdomain, ribuan domain portfolio |
| 0.3 | 2026-05-21 | Native CMS (bukan WP), ownership + share, subdomain oleh Super Admin |
| 0.4 | 2026-05-21 | Rencana database PostgreSQL (schema, index, skenario) |
| 0.5 | 2026-05-21 | Co-admin undang dengan persetujuan owner; SA transfer ownership |
| 0.6 | 2026-05-21 | Transfer ownership: owner lama **tanpa akses** |
| 0.7 | 2026-05-21 | RBAC, login aman, setup backend, setup meta (11–14) |
| 0.8 | 2026-05-21 | Integrasi Cloudflare dari admin: Tunnel, Pages, API key |
| 0.9 | 2026-05-21 | Semua setting backend di admin: RBAC, auth, rate limit selaras CF |
| 1.0 | 2026-05-21 | Deploy & lingkungan dev/staging/prod (16) |
| 1.1 | 2026-05-21 | Kontrak HTMX (17) + bisnis subdomain (18) |
| 1.2 | 2026-05-21 | Modul URL shortlink auto/manual + analitik CF (19) |
| 1.3 | 2026-05-21 | Admin Pixel: Facebook, TikTok, Google Ads (20) |
| 1.4 | 2026-05-21 | Pixel Hub: first-party, CAPI, event catalog (20) |
| 1.5 | 2026-05-21 | Rencana implementasi awal Pixel Facebook (21) |
| 1.6 | 2026-05-22 | Spec Pro Facebook + protokol komunikasi & data lengkap (21, 22) |
| 1.7 | 2026-05-22 | Kedalaman Meta Conversions API — CAPI (23) |
| 1.8 | 2026-05-22 | CAPI: hybrid vs server_first + multi-pixel ribuan domain (23) |
| 1.9 | 2026-05-22 | Akun BM putus, pixel tidak transferable, optimasi iklan (24) |
| 2.0 | 2026-05-22 | SOP S0–S3, tabel siapa punya pixel_account, insiden BM putus (24) |
| 2.1 | 2026-05-22 | Data pixel lengkap vs IP-only, tier EMQ, enricher (25) |
| 2.2 | 2026-05-22 | Plan/25 v2: verifikasi dokumentasi resmi Meta CAPI 2026 |
| 2.3 | 2026-05-22 | Plan/26: peta URL Help + Meta Pixel + CAPI business (26) |
| 2.4 | 2026-05-22 | Plan/27: desain admin — 3 dashboard, nav berkelompok, tanpa Cart/Toko (27) |
| 2.5 | 2026-05-22 | Plan/27 v1.1: Plugins, Settings, domain drawer, SEO domain-panel, tanpa jobs (27) |
| 2.7 | 2026-05-24 | Plan/29 v1: bootstrap UI-first admin HTMX mock |
| **3.0** | **2026-05-24** | **Arsitektur final:** onboarding GitHub Pages, folder `Frontend-*`, Docker-only mini PC, Plan/28 v2, hapus Plan/29 bootstrap lama |
