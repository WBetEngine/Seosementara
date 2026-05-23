# Isi sekali di GitHub → Settings → Secrets and variables → Actions
# (Browser saja — tidak perlu buka mini PC lagi setelah runner terpasang)

| Secret | Contoh / catatan |
|--------|------------------|
| `DB_PASSWORD` | Password Postgres |
| `MASTER_ENCRYPTION_KEY` | base64 32 byte |
| `SUPER_ADMIN_TOKEN` | Sama dengan admin-config.js |
| `CLOUDFLARE_API_KEY` | Global API Key / token CF |
| `CLOUDFLARE_ACCOUNT_ID` | `c3180cc77d27c46189e672bc4a74ab57` |
| `CLOUDFLARE_ACCOUNT_EMAIL` | Email akun CF |
| `CLOUDFLARE_TUNNEL_ID` | UUID tunnel |

Opsional (workflow Sync Mini PC SSH saja):

| Secret | Nilai |
|--------|--------|
| `DEPLOY_HOST` | `100.100.17.92` |
| `DEPLOY_USER` | `seosementara` |
| `DEPLOY_SSH_PASSWORD` | Password Windows |
| `DEPLOY_PATH` | `C:/Seosementara` |

**Catatan:** GitHub **menolak** password/API key di dalam file kode (push protection). Secret hanya lewat menu **Secrets** di atas.
