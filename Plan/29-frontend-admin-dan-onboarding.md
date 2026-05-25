# 29 — Frontend: Admin (Cloudflare Pages) & Onboarding (GitHub Pages)

> **Status:** Implementasi aktif.  
> **Bootstrap infra:** [28-platform-github-workers.md](./28-platform-github-workers.md)  
> **Desain admin:** [27-admin-panel-desain-ui-navigasi.md](./27-admin-panel-desain-ui-navigasi.md)

---

## 1. Dua Aplikasi Frontend Terpisah

| | **Onboarding** | **Admin CMS** |
|--|------------------|---------------|
| **Folder** | `Frontend-Onboarding/` (UI `public/` + Worker `platform-worker/`) | `Frontend-Ui-Admin/` |
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
├── public/                     # wizard UI (GitHub Pages)
│   ├── index.html
│   └── assets/js/onboarding.js
└── platform-worker/            # Cloudflare Worker API (bukan folder terpisah di root)

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

### Fase A (selesai): UI + Platform API

- Layout mobile-friendly, validasi realtime, ikon info di label
- **Platform Worker** `Frontend-Onboarding/platform-worker/` — endpoint `/admin/api/platform/*` (nyata)
- Onboarding memanggil CF API, GitHub API, GitHub Actions (bukan demo toast)

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
- Urutan wizard: **(1) CF Worker** `CLOUDFLARE_API_TOKEN` + `CLOUDFLARE_ACCOUNT_ID` → **(2) GitHub PAT** (secrets + deploy worker) → **(3) CF Zone/domain** → infra
- Setiap step: tombol **Test** sebelum **Lanjut** (langkah 3+ butuh Platform API terhubung)
- Selesai step 9: redirect ke admin CF Pages

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
| Simpan PAT/secret permanen di localStorage | Session edge + GitHub Secrets API |
| Toast "demo" tanpa panggilan API | Semua Test memanggil Platform Worker |

---

## 8. Checklist

- [x] Kerangka `Frontend-Ui-Admin/` (mock)
- [x] Buat `Frontend-Onboarding/` + workflow GitHub Pages
- [x] Pindahkan wizard ke onboarding; `bootstrap.html` → redirect ke GH Pages
- [x] Banner admin + link onboarding; `?from=onboarding` → sessionStorage
- [x] Workers Platform API — `Frontend-Onboarding/platform-worker/` + workflows bootstrap
- [ ] Backend Go API lengkap + `apiMode: live` di admin
- [ ] Admin hapus `mock-api/` setelah Go API hidup

---

## 9. Dokumen Terkait

| Plan | Isi |
|------|-----|
| [28](./28-platform-github-workers.md) | Alur bootstrap lengkap |
| [17](./17-kontrak-htmx-dan-komponen-ui.md) | Kontrak HTMX admin |
| [08](./08-roadmap-implementasi.md) | Urutan fase |

**Versi:** 2.0 — Mei 2026
