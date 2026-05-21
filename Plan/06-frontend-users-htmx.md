# 06 — Frontend Publik (HTMX) — Domain Produk

> **Bukan** frontend per domain portfolio. Lihat [09-model-domain-host-dan-subdomain.md](./09-model-domain-host-dan-subdomain.md).

## 1. Peran (Revisi)

**Frontend customer** = antarmuka publik yang dilayani dari **domain backend** `seosementara.org` dan **subdomain**-nya — yang dilihat pengunjung internet.

| Yang ini BUKAN | Yang ini ADALAH |
|----------------|-----------------|
| UI terpisah untuk setiap dari ribuan domain portfolio | UI produk di `https://seosementara.org/` |
| Deploy hostname `toko-abc.com` di CMS ini | Kelola `toko-abc.com` sebagai **data** di admin |
| Satu UI per domain portfolio | Beberapa **subdomain layanan produk** dengan tampilan berbeda |

## 2. Dua Jenis Tampilan Publik

### 2.1 Apex — `seosementara.org`

Situs utama produk: beranda, halaman marketing, dokumentasi, blog produk, dll.

| URL contoh | Isi |
|------------|-----|
| `/` | Beranda |
| `/blog/{slug}` | Artikel produk |
| `/tentang` | Halaman statis |

### 2.2 Subdomain — layanan terpisah

Setiap subdomain punya **UI HTMX sendiri** (layout, menu, fungsi):

| Host (contoh) | Modul | Spesifikasi |
|---------------|-------|-------------|
| `bola.seosementara.org` | Bola | [18](./18-bisnis-subdomain-dan-modul.md) §4 |
| `cdn.seosementara.org` | CDN | [18](./18-bisnis-subdomain-dan-modul.md) §5 |
| `url.seosementara.org` | URL Short | [18](./18-bisnis-subdomain-dan-modul.md) §6 |
| `ads.seosementara.org` | Ads | [18](./18-bisnis-subdomain-dan-modul.md) §7 |
| `comments.seosementara.org` | Comments | [18](./18-bisnis-subdomain-dan-modul.md) §8 |
| `review.seosementara.org` | Reviews | [18](./18-bisnis-subdomain-dan-modul.md) §9 |

Daftar subdomain **bukan hardcode** — didaftarkan di admin:

`https://seosementara.org/admin/setup/host`

## 3. Stack

| Komponen | Pilihan |
|----------|---------|
| Interaktivitas | **HTMX** |
| Render | Go `html/template` + partial swap |
| Hosting | **Cloudflare Pages** + API via Tunnel [15](./15-setup-cloudflare-integrasi.md) |
| Sumber file | Repo `Frontend-Users/` (+ subfolder per subdomain opsional) |

## 4. Routing (Backend)

```go
// Host: bola.seosementara.org, Path: /match/123
hostCfg := db.GetHostByHostname(host)
if hostCfg == nil { return 404 }
return renderTemplate(hostCfg.TemplateID, path)
```

Template ID contoh: `apex_default`, `subdomain_bola`, `subdomain_cdn`.

## 5. Pola HTMX (Apex)

```html
<!-- seosementara.org -->
<main hx-get="/api/public/home"
      hx-trigger="load"
      hx-swap="innerHTML">
</main>
```

Sama origin → **tanpa CORS** untuk request ke `/api/public/*`.

## 6. Pola HTMX (Subdomain)

```html
<!-- bola.seosementara.org -->
<section hx-get="/api/public/bola/fixtures"
         hx-trigger="load">
</section>
```

Namespace API publik bisa diprefix per layanan: `/api/public/bola/...`, `/api/public/url/...`.

## 7. Hubungan dengan Ribuan Domain Portfolio

Domain portfolio (mis. `toko-abc.com`) dikelola sebagai **situs native CMS** di `/admin/` — **bukan** WordPress. Tampilan publik domain portfolio (jika nanti ada) terpisah dari dokumen ini; file ini fokus ke **brand produk** `seosementara.org` dan **subdomain layanan** (bola, cdn, url, …).

## 8. Struktur Folder (Usulan)

```
Frontend-Users/
├── apex/
│   ├── layouts/
│   ├── pages/
│   └── partials/
├── subdomains/
│   ├── bola/
│   ├── cdn/
│   ├── url/
│   ├── ads/
│   ├── comments/
│   └── review/
├── static/
│   ├── css/
│   └── js/htmx.min.js
└── README.md
```

## 9. Cache & SEO

| Host | SEO |
|------|-----|
| Apex | Meta produk, OG, sitemap `seosementara.org/sitemap.xml` |
| Subdomain | Meta per layanan; sitemap opsional per host |

Cache agresif di Cloudflare untuk GET publik; invalidasi saat publish dari admin.

## 10. Dokumen Terkait

- Model domain → [09](./09-model-domain-host-dan-subdomain.md)
- Setup host → [03](./03-menu-dan-modul-cms.md)
- API publik → [07](./07-api-dan-integrasi.md)
