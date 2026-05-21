# 06 — Frontend Customer / Users (HTMX + Cloudflare Pages)

## 1. Peran

Lapisan ini adalah **tampilan publik** yang dilihat pengunjung setiap domain customer. Sama seperti admin, memakai **HTMX** dan di-host di **Cloudflare Pages** — bukan di mini CPU.

## 2. Perbedaan Admin vs Customer

| Aspek | Admin Panel | Frontend Customer |
|-------|-------------|-------------------|
| Pengguna | Operator internal | Pengunjung internet |
| Auth | Wajib login | Umumnya anonymous (read) |
| API prefix | `/api/admin/*` | `/api/public/*` |
| Cache | Sedikit | Agresif (edge + CDN) |
| Tema | Satu UI admin | Per situs / per domain |
| SEO output | Tidak relevan | Sangat relevan (HTML meta, schema) |

## 3. Stack

| Komponen | Pilihan |
|----------|---------|
| Interaktivitas | **HTMX** |
| Markup | HTML + partial server-driven |
| Styling | Per-site CSS theme |
| Hosting | **Cloudflare Pages** |
| Data | API Golang (read-mostly) |

## 4. Pola Halaman

### 4.1 Beranda

```html
<main hx-get="/api/public/sites/{site}/home"
      hx-trigger="load"
      hx-swap="innerHTML">
</main>
```

Backend mengembalikan daftar artikel terbaru, hero, dll. sebagai HTML fragment.

### 4.2 Artikel / halaman tunggal

- URL cantik: `/blog/{slug}` — Pages rewrite ke template + fetch by slug
- Meta SEO di `<head>` di-render server (Go template) atau di-inject saat swap

### 4.3 Arsip & kategori

Pagination HTMX:

```html
<button hx-get="/api/public/posts?page=2"
        hx-target="#listing"
        hx-swap="beforeend">
  Muat lebih
</button>
```

### 4.4 Pencarian (opsional)

`hx-get` dengan query `q=` — backend limit 20 hasil, debounce input.

## 5. Multi-Domain di Cloudflare Pages

| Strategi | Kapan dipakai |
|----------|---------------|
| Satu project, banyak custom domain | Banyak situs, tema sama |
| Beberapa project per kelompok | Tema berbeda total |
| `host` header → resolve `site_id` | API lookup situs by domain |

Flow:

1. Request masuk `https://customer-a.com/artikel/foo`
2. Pages serve shell HTML
3. HTMX call API dengan header `X-Forwarded-Host` atau embed `site_id` di build config per domain

## 6. Cache & Performa Publik

| Lapisan | TTL |
|---------|-----|
| Cloudflare CDN | Cache static asset 1 tahun |
| API response public | `Cache-Control: public, max-age=60` untuk listing; invalidasi on publish |
| HTMX partial | ETag support dari backend |

**Prinsip:** mini CPU tidak diload oleh traffic static — hanya API ringan dan cacheable.

## 7. SEO di Frontend

Wajib di response HTML (bukan hanya setelah JS):

- `<title>`, meta description, canonical
- Open Graph / Twitter Card
- JSON-LD (Article, WebSite) — dari modul SEO CMS
- Sitemap XML di route terpisah (bisa generate statis ke Pages saat publish)

## 8. Form Publik (Terbatas)

Jika ada kontak / newsletter:

- POST ke `/api/public/forms/contact`
- Rate limit + honeypot + Turnstile (Cloudflare)
- Tidak expose admin credentials

## 9. Struktur Folder (Usulan)

```
Frontend-Users/
├── public/
│   ├── index.html          # shell generik
│   ├── css/
│   │   ├── base.css
│   │   └── themes/
│   │       ├── site-a.css
│   │       └── site-b.css
│   └── js/htmx.min.js
├── themes/
│   └── site-a/
│       └── config.json     # site_id, domain
├── _redirects
└── wrangler.toml
```

## 10. Navigasi Publik (Bukan Menu CMS)

Dikonfigurasi per situs di admin → disimpan di API → dirender di header/footer:

- Beranda
- Blog / Artikel
- Kategori populer
- Halaman: Tentang, Kontak, Kebijakan Privasi
- Footer: social links, sitemap link

## 11. Offline & Error

- Fallback statis di Pages jika API down (halaman maintenance)
- HTMX `hx-on::response-error` tampilkan pesan ramah

## 12. Dokumen Terkait

- SEO modul admin → [03-menu-dan-modul-cms.md](./03-menu-dan-modul-cms.md)
- Public API → [07-api-dan-integrasi.md](./07-api-dan-integrasi.md)
- Backend cache invalidation → [04-backend-golang.md](./04-backend-golang.md)
