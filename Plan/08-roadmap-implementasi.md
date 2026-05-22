# 08 — Roadmap Implementasi

Roadmap ini memecah pembangunan CMS menjadi fase yang dapat dikirim secara bertahap. Estimasi waktu kalender sengaja dihindari — fokus pada **urutan dependensi** dan **risiko**.

## Fase 0 — Fondasi (Prasyarat)

| Item | Output | Folder |
|------|--------|--------|
| Repo layout | Backend, Frontend-admin, Frontend-Users, Plan | root |
| Dokumentasi | File Plan 01–20 | `Plan/` |
| Model domain | [09](./09-model-domain-host-dan-subdomain.md) disepakati | `Plan/` |
| Lingkungan | Local + staging + prod — [16](./16-deploy-dan-lingkungan.md) | ops |
| Infrastruktur | Go di mini CPU + Cloudflare DNS wildcard + Tunnel | ops |
| Routing | `Host` + `/admin/` + `/` + subdomain | Backend |

**Selesai jika:** local & staging deploy OK; smoke test [16](./16-deploy-dan-lingkungan.md) §11 lulus.

---

## Fase 1 — MVP Backend + Auth

| Item | Detail |
|------|--------|
| Migrasi DB | PostgreSQL schema v1 — lihat [10-database-postgresql.md](./10-database-postgresql.md) |
| Auth admin | Login, session, logout |
| CRUD managed-domains | Ownership + list scoped per user |
| domain_shares + invitations | Share langsung (owner) + pending (co-admin) + approve |
| notifications | Notifikasi undangan & transfer ownership |
| transfer-owner | Endpoint Super Admin |
| CRUD hosts | Super Admin — subdomain dinamis |
| Router Host+Path | `/admin/`, apex, subdomain |
| CRUD post (draft/publish) | Pagination wajib |
| Public read API | home + post by slug |
| RBAC + permission checklist | [11](./11-rbac-dan-permission-share.md) |
| Login aman | [12](./12-autentikasi-dan-login-aman.md) |

**Risiko:** CORS cross-origin — selesaikan dengan subdomain konsisten di fase 0.

---

## Fase 2 — Admin Panel HTMX

| Item | Detail |
|------|--------|
| Layout sidebar + site switcher | Menu sesuai [03](./03-menu-dan-modul-cms.md) |
| Login page | HTMX → API |
| Daftar & edit post | Partial swap |
| Media upload | Progress + validasi |
| SEO sidebar per post | Title, description, slug |

**Selesai jika:** operator bisa login, pilih situs, tulis dan publish artikel dari Pages.

---

## Fase 3 — Frontend Publik HTMX + Kontrak UI

| Item | Detail |
|------|--------|
| Kontrak HTMX | [17](./17-kontrak-htmx-dan-komponen-ui.md) |
| Apex `seosementara.org/` | Beranda, halaman statis |
| Modul URL shortlink (MVP) | [19](./19-modul-url-shortlink.md) — auto domain + redirect |
| Admin Setup → Host | CRUD host + `template_id` |
| Cache publik | Cloudflare + invalidasi |

**Selesai jika:** checklist [17](./17-kontrak-htmx-dan-komponen-ui.md) §13; `url.*` redirect jalan; share domain OK.

---

## Fase 4 — Setup Cloudflare, Meta, Backend

| Item | Detail |
|------|--------|
| Cloudflare API + Tunnel + Pages | [15](./15-setup-cloudflare-integrasi.md) |
| Domain .env + sync Pages | `PRIMARY_DOMAIN`, `API_BASE_URL` |
| Setup backend operasional | [13](./13-setup-backend-dan-sistem.md) |
| Meta global + host + domain | [14](./14-setup-meta-dan-seo.md) |

---

## Fase 5 — Operasi Massal & Hardening

| Item | Detail |
|------|--------|
| SEO global per situs | Settings |
| Bulk SEO editor | HTMX + job queue |
| Sitemap XML | Generate per situs |
| Redirect manager | 301 list |
| Worker jobs | Progress UI di admin |

**Prinsip:** semua bulk > 50 item via job, bukan synchronous loop.

---

## Fase 6 — Hardening & Skala

| Item | Detail |
|------|--------|
| Rate limiting | Public + admin |
| Audit log | Aktivitas admin |
| Backup otomatis | DB + media |
| Monitoring | Health, disk, queue depth |
| Load test mini CPU | Identifikasi bottleneck |

---

## Fase 6b — Pixel Hub (kolaborasi FB, TikTok, GAds)

| Item | Detail |
|------|--------|
| Doc | [20](./20-pixel-admin-facebook-tiktok-gads.md) — menggantikan peran Stape/GTM SS/CAPIG secara native |
| MVP | `pelacak.*` first-party + `sseo-track.js` + `/collect` + CAPI Facebook + worker `pixel_dispatch` |
| Fase 2 | Event catalog + fan-out TikTok & GAds |
| Fase 3 | Mass deploy ribuan domain + privacy hash + consent |
| Fase 4 | Diagnostics, hybrid dedup, analytics sync |

---

## Fase 7 — Modul Subdomain Lanjutan

| Modul | Doc |
|-------|-----|
| CDN, Comments | [18](./18-bisnis-subdomain-dan-modul.md) §5, §8 |
| Bola, Ads, Reviews | [18](./18-bisnis-subdomain-dan-modul.md) §4, §7, §9 |

---

## Fase 8 — Peningkatan (Backlog)

- 2FA admin
- Import/export CSV besar (chunked)
- Integrasi eksternal / sinkronisasi pihak ketiga
- Multi-bahasa konten
- R2 media offload penuh
- JSON-LD schema lanjutan

---

## Matriks Deliverable per Folder

| Fase | Backend (Go) | Frontend-admin | Frontend-Users |
|------|:------------:|:--------------:|:--------------:|
| 1 | ✓✓✓ | — | — |
| 2 | ✓ (partial HTML) | ✓✓✓ | — |
| 3 | ✓ (public API) | ✓ | ✓✓✓ |
| 4 | ✓✓ (worker) | ✓✓ | ✓ |
| 5 | ✓✓ | ✓ | ✓ |

---

## Keputusan yang Harus Diambil Sebelum Fase 1

| # | Keputusan | Opsi |
|---|-----------|------|
| 1 | Database | PostgreSQL vs SQLite |
| 2 | Media storage | Lokal vs Cloudflare R2 |
| 3 | Admin response format | HTML partial only vs JSON hybrid |
| 4 | Domain strategy | Tunnel hostname vs IP |
| 5 | Multi-tenant model | `site_id` column global |

Catat keputusan di bagian bawah file ini setelah final.

---

## Catatan Keputusan

| Tanggal | Keputusan | Catatan |
|---------|-----------|---------|
| 2026-05-21 | Backend: **Golang** | Mini CPU |
| 2026-05-21 | Admin & Users UI: **HTMX** | Cloudflare Pages |
| 2026-05-21 | Admin di **`/admin/`** pada domain produk | Bukan per-domain portfolio |
| 2026-05-21 | Frontend publik = **apex + subdomain** | Setup di `/admin/setup/host` |
| 2026-05-21 | Portfolio = **ribuan domain di DB** | Banyak pekerja, pagination wajib |
| 2026-05-21 | **Bukan WordPress** | Situs native CMS |
| 2026-05-21 | **Ownership + share** | Pekerja hanya domain sendiri + dibagikan |
| 2026-05-21 | **Subdomain dinamis** | Hanya Super Admin di Setup → Host |
| 2026-05-21 | Co-admin undang → **owner setujui** | Notifikasi + `domain_share_invitations` |
| 2026-05-21 | **SA transfer ownership** | API + audit + notifikasi |
| | | |

---

## Indeks Dokumen Plan

Kembali ke [README.md](./README.md) untuk daftar lengkap semua file perencanaan.
