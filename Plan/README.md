# Dokumentasi Perencanaan — Seosementara CMS

Dokumen ini menguraikan perencanaan sistem CMS untuk proyek **Seosementara**. Pembahasan sengaja dipecah ke beberapa file agar mudah dikembangkan, direview, dan diubah tanpa satu file raksasa.

## Cara Membaca

Baca berurutan jika Anda baru memulai. Jika fokus ke satu lapisan, buka file yang relevan saja.

| No. | File | Isi |
|-----|------|-----|
| 01 | [01-visi-dan-gambaran-sistem-cms.md](./01-visi-dan-gambaran-sistem-cms.md) | Visi produk, ruang lingkup CMS, prinsip desain |
| 02 | [02-arsitektur-dan-infrastruktur.md](./02-arsitektur-dan-infrastruktur.md) | Arsitektur tiga lapisan, mini CPU, Cloudflare Pages |
| 03 | [03-menu-dan-modul-cms.md](./03-menu-dan-modul-cms.md) | Daftar menu admin, modul, hak akses per peran |
| 04 | [04-backend-golang.md](./04-backend-golang.md) | Backend Go: struktur, API, database, performa |
| 05 | [05-admin-panel-htmx.md](./05-admin-panel-htmx.md) | Panel admin HTMX di Cloudflare Pages |
| 06 | [06-frontend-users-htmx.md](./06-frontend-users-htmx.md) | Situs pengguna/customer HTMX di Cloudflare Pages |
| 07 | [07-api-dan-integrasi.md](./07-api-dan-integrasi.md) | Kontrak API, auth, CORS, cache |
| 08 | [08-roadmap-implementasi.md](./08-roadmap-implementasi.md) | Fase implementasi dan prioritas |

## Ringkasan Stack

| Lapisan | Teknologi | Hosting |
|---------|-----------|---------|
| Backend API | **Golang** | Mini CPU (self-hosted) |
| Admin Panel | **HTMX** + HTML/CSS | **Cloudflare Pages** |
| Frontend Customer | **HTMX** + HTML/CSS | **Cloudflare Pages** |

## Konteks Skala (Prinsip Tetap)

Proyek ini ditargetkan untuk operasi massal (banyak domain, volume konten besar). Semua keputusan arsitektur harus mengutamakan:

- Query database terarah dan terbatas
- Operasi berat di-batch, bukan loop tanpa batas
- Cache dengan invalidasi yang jelas
- Hindari timeout dan beban berlebihan di server terbatas (mini CPU)

Detail teknis performa ada di file **04** dan **07**.

## Status Dokumen

| Versi | Tanggal | Catatan |
|-------|---------|---------|
| 0.1 | 2026-05-21 | Draft awal struktur perencanaan |
