# Platform Worker — API bootstrap (production)

Cloudflare Worker: `/admin/api/platform/*`

Dipanggil oleh onboarding GitHub Pages (bukan mode demo).

## Setup sekali

```bash
cd platform-worker
npm install
npx wrangler kv namespace create SETUP_KV
# Salin id ke wrangler.toml → ganti PLACEHOLDER_KV_ID
```

## GitHub Secrets (repo)

| Secret | Untuk |
|--------|--------|
| `CLOUDFLARE_API_TOKEN` | Deploy worker + workflows Pages |
| `CLOUDFLARE_ACCOUNT_ID` | Deploy worker |
| `PLATFORM_KV_ID` | (opsional) ganti PLACEHOLDER di wrangler.toml via CI |

## Deploy

Push ke `main` → workflow **Deploy Platform Worker**  
Menulis `Frontend-Onboarding/public/assets/js/platform-api-url.js` dengan URL worker.

## Uji lokal

```bash
npm run dev
# Buka onboarding dengan ?api=http://localhost:8787
```

## Endpoint

Lihat [Plan/28-platform-github-workers.md](../Plan/28-platform-github-workers.md) §5.
