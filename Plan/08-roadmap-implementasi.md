# 08 — Roadmap Implementasi

Roadmap ini memecah pembangunan CMS menjadi fase yang dapat dikirim secara bertahap. Estimasi waktu kalender sengaja dihindari — fokus pada **urutan dependensi** dan **risiko**.

## Fase 0 — Fondasi (Prasyarat)

| Item | Output | Folder |
|------|--------|--------|
| Repo layout | Backend, Frontend-admin, Frontend-Users, Plan | root |
| Dokumentasi | File Plan 01–09 | `Plan/` |
| Model domain | [09](./09-model-domain-host-dan-subdomain.md) disepakati | `Plan/` |
| Infrastruktur | Go di mini CPU + Cloudflare DNS wildcard + Tunnel | ops |
| Routing | `Host` + `/admin/` + `/` + subdomain | Backend |

**Selesai jika:** `seosementara.org/health` OK, `/admin/login` tampil, satu subdomain contoh respon.

---

## Fase 1 — MVP Backend + Auth

| Item | Detail |
|------|--------|
| Migrasi DB | sites, users, posts, pages, taxonomies |
| Auth admin | Login, session, logout |
| CRUD managed-domains | Ownership + list scoped per user |
| domain_shares | Share / unshare ke admin lain |
| CRUD hosts | Super Admin — subdomain dinamis |
| Router Host+Path | `/admin/`, apex, subdomain |
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

## Fase 3 — Frontend Publik HTMX

| Item | Detail |
|------|--------|
| Apex `seosementara.org/` | Beranda, halaman statis |
| Subdomain contoh | Minimal 2: `url.`, `cdn.` atau `bola.` |
| Admin Setup → Host | CRUD host + pilih template |
| Cache publik | Cloudflare + invalidasi on publish |

**Selesai jika:** pengunjung akses apex + satu subdomain; Super Admin bisa tambah host; pekerja share domain ke rekan.

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
| | | |

---

## Indeks Dokumen Plan

Kembali ke [README.md](./README.md) untuk daftar lengkap semua file perencanaan.
