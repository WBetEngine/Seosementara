# 09 — Model Domain, Host, dan Subdomain

Dokumen ini merangkum keputusan arsitektur dari diskusi desain. File lain (01–08) harus selaras dengan model di sini.

## 1. Konsep Inti (Wajib Dipahami)

Ada **dua kelompok "domain"** yang berbeda:

| Kelompok | Contoh | Fungsi |
|----------|--------|--------|
| **Domain produk (sistem)** | `seosementara.org`, `bola.seosementara.org` | Tempat UI admin, UI publik produk, dan layanan subdomain berjalan |
| **Domain yang dikelola (portfolio)** | Ribuan domain customer / proyek SEO | **Data** di database — dikelola lewat admin, **bukan** hostname terpisah untuk frontend HTMX |

**Frontend customer** = tampilan publik **domain backend** (`seosementara.org` dan subdomain-nya), **bukan** satu frontend HTMX per domain portfolio.

**Admin panel** = path **`/admin/`** di domain utama, mis. `https://seosementara.org/admin/`.

## 2. Peta URL

### 2.1 Domain utama — `seosementara.org`

| Path / area | Pengguna | Contoh URL |
|-------------|----------|------------|
| `/` | Pengunjung produk | `https://seosementara.org/` |
| `/blog/...` | Konten publik (jika ada) | `https://seosementara.org/blog/artikel` |
| `/admin/` | Pekerja / operator internal | `https://seosementara.org/admin/` |
| `/admin/login` | Login pekerja | `https://seosementara.org/admin/login` |
| `/api/admin/*` | API admin (sama origin) | `https://seosementara.org/api/admin/posts` |
| `/api/public/*` | API publik | `https://seosementara.org/api/public/home` |

### 2.2 Subdomain — layanan berbeda, UI berbeda

Setiap subdomain punya **tampilan HTMX sendiri** (template, menu, fungsi), tetap di ekosistem `seosementara.org`:

| Subdomain (contoh) | Peran (draft) |
|--------------------|---------------|
| `bola.seosementara.org` | Modul/layanan Bola — UI khusus |
| `cdn.seosementara.org` | CDN / aset / delivery |
| `url.seosementara.org` | Short URL / redirect |
| `ads.seosementara.org` | Iklan / kampanye |
| `comments.seosementara.org` | Komentar |
| `review.seosementara.org` | Ulasan |
| *(tambah lainnya)* | Didefinisikan lewat admin |

Subdomain **bukan** domain portfolio ribuan — itu entri terpisah di modul **Situs / Domain**.

### 2.3 Domain portfolio (ribuan)

- Disimpan sebagai record: `managed_domains` (nama bebas di DB)
- Dikelola di admin: SEO, konten, batch, status — tanpa deploy frontend baru per domain
- Pekerja memilih **domain mana** yang sedang dikerjakan (site switcher / filter)

## 3. Diagram Arsitektur

```mermaid
flowchart TB
  subgraph workers [Banyak Pekerja]
    W1[Pekerja A]
    W2[Pekerja B]
  end
  subgraph visitors [Pengunjung Publik]
    V1[Pengunjung seosementara.org]
    V2[Pengunjung bola.*]
  end
  subgraph cf [Cloudflare - DNS SSL Cache]
    Wildcard["*.seosementara.org + seosementara.org"]
  end
  subgraph mini [Mini CPU - Go]
    Router[Router: Host + Path]
    AdminUI["/admin/* HTMX"]
    PublicUI["/ HTMX - root domain"]
    SubUI["Subdomain templates"]
    API["/api/*"]
    DB[(Database)]
  end
  W1 --> Wildcard
  W2 --> Wildcard
  V1 --> Wildcard
  V2 --> Wildcard
  Wildcard --> Router
  Router --> AdminUI
  Router --> PublicUI
  Router --> SubUI
  Router --> API
  AdminUI --> API
  PublicUI --> API
  SubUI --> API
  API --> DB
```

## 4. Konfigurasi Host di Admin

Semua subdomain dan binding host dikelola di:

**`https://seosementara.org/admin/setup/host`**

| Field (konsep) | Keterangan |
|----------------|------------|
| `hostname` | `bola.seosementara.org` atau apex |
| `type` | `apex` \| `subdomain` \| `path_prefix` |
| `template_id` | UI HTMX mana yang dipakai |
| `enabled` | Aktif / maintenance |
| `notes` | Keterangan operator |

Tanpa entri di **Setup → Host**, subdomain tidak dilayani (404 atau halaman default).

## 5. Routing di Backend (Go)

Pseudo-logic:

```go
func route(req) {
  host := normalizeHost(req.Host)   // bola.seosementara.org
  path := req.URL.Path

  switch {
  case host == "seosementara.org" && strings.HasPrefix(path, "/admin/"):
    serveAdminHTMX(host, path)
  case host == "seosementara.org" && strings.HasPrefix(path, "/api/"):
    serveAPI(path)
  case host == "seosementara.org":
    servePublicHTMX(path)           // frontend customer - apex
  default:
    h := lookupHostConfig(host)     // dari DB, diisi via admin/setup/host
    if h == nil { return 404 }
    serveSubdomainHTMX(h, path)
  }
}
```

## 6. Skala: Ribuan Domain + Banyak Pekerja

### 6.1 Ribuan domain (portfolio)

| Tantangan | Solusi |
|-----------|--------|
| List domain lambat | Pagination + search + index `(status, name)` |
| Filter operasi | Wajib pilih / filter domain sebelum bulk job |
| Batch | Job queue per chunk; tidak load 1000 sekaligus |
| Audit | Siapa mengubah domain X — log aktivitas |

### 6.2 Banyak pekerja (concurrent)

| Tantangan | Solusi |
|-----------|--------|
| Tabrakan edit | Optional: lock optimistik `updated_at` / pesan konflik |
| Hak akses | RBAC + scope domain (hanya subset domain) |
| Beban login | Session terpisah; rate limit login |
| Dashboard | Agregat cached — bukan hitung ulang 1000 domain tiap refresh |

## 7. Hosting (Revisi dari Draft Awal)

| Komponen | Revisi |
|----------|--------|
| Backend Go | Mini CPU — **tetap** |
| Admin HTMX | **`seosementara.org/admin/`** — dilayani origin (Go), bukan proyek terpisah per domain |
| Frontend customer | **`seosementara.org/`** + subdomain — dilayani origin (Go) |
| Cloudflare | DNS wildcard, proxy/Tunnel, cache — **bukan** satu Pages per domain portfolio |
| Folder `Frontend-admin/` & `Frontend-Users/` | Sumber template HTML/HTMX di repo; di-build atau di-embed ke binary Go |

Cloudflare Pages masih bisa dipakai untuk **asset statis** (CSS/JS) jika diinginkan, asalkan routing `/admin/` dan subdomain tetap konsisten di DNS (Workers route → origin).

## 8. Perbedaan dengan Asumsi Lama (Catatan Migrasi Plan)

| Asumsi lama (salah) | Model baru (benar) |
|---------------------|-------------------|
| Frontend customer = domain per situs portfolio | Frontend customer = UI domain produk `seosementara.org` |
| Admin di subdomain `admin.` atau Pages terpisah | Admin di path `/admin/` |
| Ribuan hostname frontend | Ribuan **record domain** di DB, satu panel admin |
| Subdomain = customer site | Subdomain = **layanan produk** (bola, cdn, url, …) |

## 9. Pertanyaan Terbuka

- Apakah `www.seosementara.org` redirect ke apex?
- Apakah portfolio domain perlu CNAME ke infrastruktur lain (WP di shared hosting) terpisah dari UI produk?
- Template per subdomain: satu repo folder per subdomain atau config-driven?

Jawaban ditambahkan ke file ini setelah diputuskan.

## 10. Dokumen Terkait

- [02-arsitektur-dan-infrastruktur.md](./02-arsitektur-dan-infrastruktur.md)
- [03-menu-dan-modul-cms.md](./03-menu-dan-modul-cms.md) — menu Setup → Host
- [05-admin-panel-htmx.md](./05-admin-panel-htmx.md)
- [06-frontend-users-htmx.md](./06-frontend-users-htmx.md)
