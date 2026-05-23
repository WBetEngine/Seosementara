# Deploy mini PC (zero-touch)

Anda **tidak perlu** buka mini PC setelah runner terpasang.

## Arsitektur

```text
push main → Build image (GHCR) → Deploy Mini PC (self-hosted runner)
                                      ↓
                              docker compose (env dari GitHub Secrets)
                              tanpa file .env di disk
```

| Komponen | Lokasi |
|----------|--------|
| Kode | GitHub `main` |
| Image API | GHCR |
| Secret | GitHub → Settings → Secrets |
| Runtime | Mini PC (Docker saja) |
| Admin UI | Cloudflare Workers |

## Setup sekali

### 1. GitHub Secrets

| Secret | Wajib |
|--------|-------|
| `DB_PASSWORD` | Ya |
| `MASTER_ENCRYPTION_KEY` | Ya |
| `SUPER_ADMIN_TOKEN` | Ya |
| `CLOUDFLARE_API_KEY` | Bootstrap CF |
| `CLOUDFLARE_ACCOUNT_ID` | Bootstrap CF |
| `CLOUDFLARE_ACCOUNT_EMAIL` | Bootstrap CF + deploy Workers (Global API Key) |
| `CLOUDFLARE_TUNNEL_ID` | Opsional |

### 2. Self-hosted runner (mini PC, Administrator)

```powershell
cd C:\Seosementara
.\scripts\install-github-runner.ps1
```

Token dari: **Settings → Actions → Runners → New self-hosted runner**

## Setelah setup

Edit kode di GitHub → push `main` → deploy otomatis.

**Tidak ada** file `.env` di mini PC — Docker Compose membaca environment dari runner.

## File di mini PC

```
C:\Seosementara\
  docker-compose.prod.yml
  Backend\migrations\
  scripts\mini-pc-remote-deploy.ps1
  scripts\bootstrap-cloudflare.ps1
```

## Infra sekali (sudah)

- Docker Desktop
- cloudflared tunnel → api.apidevel.org
- OpenSSH (opsional, tidak dipakai workflow utama)
