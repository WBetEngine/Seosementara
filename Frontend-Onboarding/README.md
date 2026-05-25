# Frontend-Onboarding — satu folder untuk first boot

Semua yang terkait **wizard onboarding** ada di sini (UI + Cloudflare Worker API). Tidak ada salinan terpisah di `docs/` atau `platform-worker/` di root repo.

## Struktur

```
Frontend-Onboarding/
├── public/                 # UI wizard → GitHub Pages
│   ├── index.html
│   └── assets/
│       ├── css/onboarding.css
│       └── js/
│           ├── config.js
│           ├── onboarding.js
│           ├── platform-api.js
│           └── platform-api-url.js   # diisi CI setelah deploy worker
└── platform-worker/        # Cloudflare Worker sse-platform (API bootstrap)
    ├── src/
    ├── wrangler.toml
    ├── package.json
    └── scripts/ensure-setup-kv.sh
```

| Komponen | Platform | Folder |
|----------|----------|--------|
| Wizard UI | **GitHub Pages** | `public/` |
| Platform API | **Cloudflare Workers** | `platform-worker/` |

GitHub Actions (di root `.github/workflows/`) hanya **memanggil** path di atas — bukan folder duplikat.

## Deploy UI (GitHub Pages)

**Disarankan:** Settings → Pages → Source: **GitHub Actions** → workflow `Deploy Onboarding (GitHub Pages)` → artifact `Frontend-Onboarding/public`.

Alternatif lama `docs/` di root **tidak dipakai lagi** — hindari duplikasi.

## Deploy Platform Worker (Cloudflare)

**Penting:** Setiap push ke `Frontend-Onboarding/**` memicu workflow ini. Jika secret belum ada, job deploy **dilewati** (bukan error) — UI Pages tetap bisa hijau sementara Worker belum ada.

### Cara pertama kali (urutan disarankan)

1. Repo → **Settings → Secrets and variables → Actions** → **New repository secret** (bukan *Variables*):
   - `CLOUDFLARE_API_TOKEN` — token dari [Cloudflare API Tokens](https://dash.cloudflare.com/profile/api-tokens) (permission Workers + KV minimal).
   - `CLOUDFLARE_ACCOUNT_ID` — ID akun (sidebar dashboard Cloudflare).
2. Actions → **Deploy Platform Worker** → **Run workflow** (boleh kosongkan input jika secret sudah di-set).
3. Setelah job **deploy** hijau, `platform-api-url.js` ter-update otomatis → refresh onboarding langkah **2c**.

Alternatif tanpa secret: Run workflow dan isi input `cloudflare_api_token` + `cloudflare_account_id` (nilai terlihat di log run — kurang aman).

Setelah Worker hidup, langkah 2 wizard bisa menyimpan `CLOUDFLARE_*` lewat PAT (`initial-setup`) dan memicu deploy berikutnya via `repository_dispatch`.

Manual:

```bash
cd Frontend-Onboarding/platform-worker
npm install
export CLOUDFLARE_API_TOKEN=...
export CLOUDFLARE_ACCOUNT_ID=...
bash scripts/ensure-setup-kv.sh
npm run dev   # http://localhost:8787
```

UI lokal:

```bash
npx --yes serve Frontend-Onboarding/public -p 3080
# http://localhost:3080/?api=http://localhost:8787
```

## Urutan wizard (ringkas)

1. Kredensial Cloudflare → 2. GitHub PAT + deploy worker + URL `*.workers.dev` → 3. Zone → 4–9 infra.

Lihat `Plan/29-frontend-admin-dan-onboarding.md`.
