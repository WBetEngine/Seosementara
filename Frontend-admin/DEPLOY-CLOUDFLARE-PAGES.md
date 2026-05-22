# Deploy Admin UI → Cloudflare Pages (dari GitHub `main`)

Prototype admin **tidak punya langkah build** (bukan npm/Vite). Yang di-deploy = folder `public/` apa adanya.

Push ke `main` **hanya otomatis deploy** jika proyek Pages sudah **Connect to Git** di dashboard Cloudflare (sekali setup).

---

## 1. Apa yang terjadi saat “build”?

| Tahap | Admin UI sekarang |
|-------|-------------------|
| **Build** | Tidak mengompilasi — tidak ada `npm run build` |
| **Output** | Isi `Frontend-admin/public/` langsung jadi situs |
| **Trigger** | Push/commit ke branch production (`main`) |
| **Hosting** | Cloudflare Pages (free) |

Setara perintah lokal:

```bash
# Tidak perlu build; yang di-upload = folder public
ls Frontend-admin/public/admin/index.html
```

Nanti jika pakai Tailwind/ bundler, barulah **Build command** diisi (mis. `npm ci && npm run build`).

---

## 2. Setup sekali di Cloudflare Dashboard

1. Login [Cloudflare Dashboard](https://dash.cloudflare.com) → **Workers & Pages** → **Create** → **Pages** → **Connect to Git**.
2. Authorize GitHub → pilih repo **`WBetEngine/Seosementara`** (atau nama repo Anda).
3. **Create project** — isi seperti tabel:

| Field | Nilai |
|-------|--------|
| Project name | `seosementara-admin` (bebas) |
| Production branch | `main` |
| **Root directory** | `Frontend-admin` |
| Framework preset | **None** |
| **Build command** | *(kosongkan)* atau `exit 0` |
| **Build output directory** | `public` |

4. **Environment variables** (opsional untuk prototype):

| Variable | Contoh |
|----------|--------|
| `PRIMARY_DOMAIN` | `seosementara.org` |
| `API_BASE_URL` | `https://seosementara.org` |

5. **Save and Deploy** — Cloudflare clone repo, “build” (hampir instan), publish.

6. **Custom domains** → Add `seosementara.org` (atau subdomain dulu untuk uji).

Setelah ini: **setiap push ke `main`** yang mengubah file di `Frontend-admin/` memicu deploy baru (lihat tab **Deployments**).

---

## 3. Cek deploy berhasil

| Cek | URL |
|-----|-----|
| Pages default | `https://<project-name>.pages.dev/admin/index.html` |
| Login | `.../admin/login.html` |
| Custom domain | `https://seosementara.org/admin/index.html` |

Di dashboard: **Deployments** → status **Success** → view log (build ~ beberapa detik).

---

## 4. Routing `/admin/` (file `_redirects`)

File `public/_redirects` sudah disertakan agar:

- `/admin` → `/admin/index.html`
- `/admin/` → `/admin/index.html`

---

## 5. Backend belum di Pages

| Path | Dilayani oleh |
|------|----------------|
| `/admin/*`, `/static/*` | **Pages** (file di `public/`) |
| `/api/*` | **Tunnel** → Go di mini CPU (besok) |

Prototype memakai mock `/_partials/*.html` di Pages; setelah backend siap, HTMX mengarah ke `/api/admin/...` + Tunnel.

---

## 6. Alternatif: deploy manual (tanpa Git di CF)

```bash
cd Frontend-admin
npx wrangler pages deploy public --project-name=seosementara-admin
```

Perlu `wrangler login` + API token. Untuk produksi tetap disarankan **Git integration** (Plan/15 §7.4).

---

## 7. Troubleshooting

| Masalah | Solusi |
|---------|--------|
| Push `main` tapi tidak deploy | Belum connect Git di Pages, atau root directory salah |
| 404 di `/admin` | Pastikan `_redirects` ada; buka `/admin/index.html` |
| Build gagal “command not found” | Kosongkan build command; preset = None |
| HTMX 404 ke API | Normal — backend belum; mock masih di `_partials/` |

---

## 8. GitHub Actions (opsional, belum wajib)

Plan/16 menyebut workflow CI terpisah. Fase 1 cukup **Cloudflare Connect Git**. Actions hanya jika ingin test/lint sebelum deploy.

**Versi:** 2026-05-22 — selaras Plan/15, Plan/16, Plan/27
