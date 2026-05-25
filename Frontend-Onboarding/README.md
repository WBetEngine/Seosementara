# Frontend-Onboarding — Setup Infra Pertama Kali

Wizard **production** — memanggil **Platform Worker API** (bukan demo).

| | |
|--|--|
| **URL** | https://wbetengine.github.io/Seosementara/ |
| **API** | `https://sse-platform.<akun>.workers.dev/admin/api/platform/*` |
| **Deploy UI** | push `main` → folder `docs/` atau GitHub Actions |

## Prasyarat (urutan)

1. **GitHub Secrets** di repo: `CLOUDFLARE_API_TOKEN`, `CLOUDFLARE_ACCOUNT_ID`
2. **KV namespace** + deploy worker — lihat [platform-worker/README.md](../platform-worker/README.md)
3. Workflow **Deploy Platform Worker** sukses → file `platform-api-url.js` terisi URL worker
4. **Pages** → branch `main`, folder **`/docs`**

## Uji lokal

```bash
# Terminal 1 — worker
cd platform-worker && npm install && npm run dev

# Terminal 2 — UI
npx --yes serve Frontend-Onboarding/public -p 3080
# Buka http://localhost:3080/?api=http://localhost:8787
```

## Alur operator

1. Validasi PAT → session edge
2. Test + simpan Cloudflare → CF API nyata + GitHub Secrets
3. Test SSH → GitHub Actions `bootstrap-ssh-test`
4. Register runner → SSH install runner di mini PC
5. Buat tunnel → Cloudflare API + workflow install cloudflared
6. Simpan DB secrets → GitHub Secrets API
7. Deploy → trigger backend + admin Pages + public Pages
8. Buka admin CF Pages

Dokumen: [Plan/28](../Plan/28-platform-github-workers.md), [Plan/29](../Plan/29-frontend-admin-dan-onboarding.md).
