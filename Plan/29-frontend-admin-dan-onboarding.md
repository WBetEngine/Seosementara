# 29 — Frontend: Admin (Cloudflare Pages) & Onboarding (GitHub Pages)

> **Status:** Implementasi aktif.  
> **Bootstrap infra:** [28-platform-github-workers.md](./28-platform-github-workers.md)  
> **Desain admin:** [27-admin-panel-desain-ui-navigasi.md](./27-admin-panel-desain-ui-navigasi.md)

---

## 1. Dua Aplikasi Frontend Terpisah

| | **Onboarding** | **Admin CMS** |
|--|------------------|---------------|
| **Folder** | `Frontend-Onboarding/` | `Frontend-Ui-Admin/` |
| **Hosting** | GitHub Pages | Cloudflare Pages |
| **Deploy** | Otomatis tiap push `main` | Setelah bootstrap (GitHub Actions) |
| **Pengguna** | Operator setup infra (sekali) | Pekerja CMS (sehari-hari) |
| **Backend** | Workers Platform API | Go API via Tunnel (`/api/admin/*`) |
| **Mock data** | Hanya sampai Workers API hidup | Fase 0: mock → live setelah Go API |

---

## 2. Kondisi First Boot

| # | Kondisi |
|---|---------|
| 1 | Belum ada PostgreSQL / data |
| 2 | Tidak ada `.env` di mini PC |
| 3 | Tunnel belum aktif |
| 4 | Admin CF Pages belum deploy (normal) |
| 5 | **GitHub Pages onboarding sudah bisa online** hanya dengan push repo |

Operator **pertama kali** selalu mulai dari **URL GitHub Pages onboarding**, bukan admin Cloudflare.

---

## 3. Struktur Folder

```
Frontend-Onboarding/
└── public/
    ├── index.html              # wizard bootstrap
    ├── assets/css/onboarding.css
    └── assets/js/onboarding.js   # panggil Workers Platform API

Frontend-Ui-Admin/
└── public/
    ├── admin/                    # halaman CMS (tanpa bootstrap wizard)
    ├── assets/
    ├── partials/
    └── mock-api/                 # hapus setelah apiMode=live

Frontend-Publik/
└── public/                     # apex + subdomain UI
```

---

## 4. Admin UI — Fase Implementasi

### Fase A (selesai): Kerangka statis

- Layout mobile-friendly, drawer, sidebar ([27](./27-admin-panel-desain-ui-navigasi.md))
- Mock partials untuk demo HTMX
- **Pindahkan** `/admin/bootstrap.html` → `Frontend-Onboarding` (rencana)

### Fase B: Sambung Go API

```javascript
// Frontend-Ui-Admin/public/assets/js/config.js
SSEO.apiMode = 'live';
SSEO.apiBase = '';  // same-origin via Tunnel
```

- Hapus banner demo & folder `mock-api/`
- Login session nyata ([12](./12-autentikasi-dan-login-aman.md))

### Fase C: Cek bootstrap

Saat buka admin, `GET /admin/api/platform/setup/status`:

- Jika belum selesai → banner + link ke GitHub Pages onboarding
- Jika selesai → dashboard normal

---

## 5. Onboarding UI — Konfigurasi

```javascript
// Frontend-Onboarding/public/assets/js/config.js
SSEO.platformApiBase = 'https://<workers-host>/admin/api/platform';
SSEO.adminUrlAfterComplete = 'https://seosementara.org/admin/login.html';
```

- **Jangan** simpan API Key / password SSH di `localStorage` permanen
- Setiap step: tombol **Test** sebelum **Lanjut**
- Selesai step 8: redirect ke admin CF Pages

---

## 6. Halaman Admin (Cloudflare Pages)

Selaras [27](./27-admin-panel-desain-ui-navigasi.md) — **tanpa** menu bootstrap:

| Grup | Path |
|------|------|
| Ringkasan | `/admin/dashboard*.html` |
| Domain | `/admin/domain/` |
| Konten | `/admin/content/` |
| SEO | `/admin/seo/` |
| Plugins | `/admin/plugins/` |
| Settings | `/admin/settings/` |

Settings → Infra / Mini PC (status runner, tunnel, health) — **bukan** first boot wizard.

---

## 7. Larangan

| Larangan | Alasan |
|----------|--------|
| Onboarding wizard di CF Pages admin | Pemisahan hosting ([28](./28-platform-github-workers.md)) |
| Script generator Python di repo | HTML sudah final; hapus setelah generate |
| `npx serve` sebagai alur produksi | Kode di GitHub; deploy via Pages |
| Bootstrap form simpan secret ke localStorage | Keamanan |

---

## 8. Checklist

- [x] Kerangka `Frontend-Ui-Admin/` (mock)
- [x] Buat `Frontend-Onboarding/` + workflow GitHub Pages
- [x] Pindahkan wizard ke onboarding; `bootstrap.html` → redirect ke GH Pages
- [x] Banner admin + link onboarding; `?from=onboarding` → sessionStorage
- [ ] Workers Platform API ([28](./28-platform-github-workers.md) §5)
- [ ] Backend Go + Docker
- [ ] Admin `apiMode: live`

---

## 9. Dokumen Terkait

| Plan | Isi |
|------|-----|
| [28](./28-platform-github-workers.md) | Alur bootstrap lengkap |
| [17](./17-kontrak-htmx-dan-komponen-ui.md) | Kontrak HTMX admin |
| [08](./08-roadmap-implementasi.md) | Urutan fase |

**Versi:** 2.0 — Mei 2026
