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

1. Isi GitHub Secrets: `CLOUDFLARE_API_TOKEN`, `CLOUDFLARE_ACCOUNT_ID` (opsional `PLATFORM_KV_ID`).
2. Actions → **Deploy Platform Worker** → Run workflow.
3. Setelah hijau, `platform-api-url.js` ter-update otomatis.

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
