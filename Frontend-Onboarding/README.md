# Frontend-Onboarding — Setup Infra Pertama Kali

Wizard **first boot** di **GitHub Pages** (bukan Cloudflare Pages admin).

| | |
|--|--|
| **URL produksi** | `https://wbetengine.github.io/Seosementara/` |
| **Deploy** | GitHub Actions `pages-onboarding.yml` pada push `main` |
| **Backend** | Workers Platform API (`/admin/api/platform/*`) — fase berikutnya |

## Prasyarat GitHub

1. Repo → **Settings** → **Pages** → Source: **GitHub Actions**
2. Push folder `Frontend-Onboarding/public/` ke `main`

## Preview lokal

```bash
npx --yes serve Frontend-Onboarding/public -p 3080
```

Buka http://localhost:3080/

## Setelah selesai

Operator diarahkan ke admin Cloudflare Pages (`SSEO.adminUrlAfterComplete` di `config.js`).

Dokumen: [Plan/29-frontend-admin-dan-onboarding.md](../Plan/29-frontend-admin-dan-onboarding.md), [Plan/28](../Plan/28-platform-github-workers.md).
