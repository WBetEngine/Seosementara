# Platform Worker (`sse-platform`)

API bootstrap onboarding — di-host di **Cloudflare Workers**, bukan GitHub.

Parent folder: `Frontend-Onboarding/` (satu tempat dengan UI wizard di `../public/`).

## Setup lokal

```bash
cd Frontend-Onboarding/platform-worker
npm install
export CLOUDFLARE_API_TOKEN=...
export CLOUDFLARE_ACCOUNT_ID=...
bash scripts/ensure-setup-kv.sh
npm run dev
```

## Deploy

GitHub Actions: workflow **Deploy Platform Worker** (path `Frontend-Onboarding/platform-worker`).

Endpoint base: `https://sse-platform.<subdomain>.workers.dev/admin/api/platform`
