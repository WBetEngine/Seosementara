# 08 — Roadmap Implementasi

Roadmap ini memecah pembangunan CMS menjadi fase yang dapat dikirim secara bertahap. Estimasi waktu kalender sengaja dihindari — fokus pada **urutan dependensi** dan **risiko**.

## Fase 0 — Fondasi (Prasyarat)

| Item | Output | Folder |
|------|--------|--------|
| Repo layout | Backend, Frontend-admin, Frontend-Users, Plan | root |
| Dokumentasi | File Plan 01–08 | `Plan/` |
| Infrastruktur mini CPU | Go binary, DB, reverse proxy, Tunnel | ops |
| Cloudflare Pages | Dua project (admin + users), env vars | CF dashboard |

**Selesai jika:** health check API OK, Pages deploy hello-world memanggil API.

---

## Fase 1 — MVP Backend + Auth

| Item | Detail |
|------|--------|
| Migrasi DB | sites, users, posts, pages, taxonomies |
| Auth admin | Login, session, logout |
| CRUD situs | List/create/update |
| CRUD post (draft/publish) | Pagination wajib |
| Public read API | home + post by slug |
| RBAC dasar | Super Admin, Editor |

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

## Fase 3 — Frontend Customer HTMX

| Item | Detail |
|------|--------|
| Shell HTML per tema | Beranda, single post, arsip |
| Resolve site by host | Multi-domain |
| Meta SEO di HTML | Title, OG, canonical |
| Cache header public API | TTL + invalidasi on publish |

**Selesai jika:** pengunjung membaca artikel published dengan performa baik di edge.

---

## Fase 4 — SEO & Operasi Massal

| Item | Detail |
|------|--------|
| SEO global per situs | Settings |
| Bulk SEO editor | HTMX + job queue |
| Sitemap XML | Generate per situs |
| Redirect manager | 301 list |
| Worker jobs | Progress UI di admin |

**Prinsip:** semua bulk > 50 item via job, bukan synchronous loop.

---

## Fase 5 — Hardening & Skala

| Item | Detail |
|------|--------|
| Rate limiting | Public + admin |
| Audit log | Aktivitas admin |
| Backup otomatis | DB + media |
| Monitoring | Health, disk, queue depth |
| Load test mini CPU | Identifikasi bottleneck |

---

## Fase 6 — Peningkatan (Backlog)

- 2FA admin
- Import/export CSV besar (chunked)
- Integrasi WordPress / sinkronisasi
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
| | | |

---

## Indeks Dokumen Plan

Kembali ke [README.md](./README.md) untuk daftar lengkap semua file perencanaan.
