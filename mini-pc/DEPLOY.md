# Deploy mini PC (manual)

Semua deploy dilakukan **langsung di mini PC**. Tidak ada GitHub runner atau sync otomatis.

## Prasyarat

- Docker Desktop
- cloudflared tunnel → `api.apidevel.org`
- Folder `C:\Seosementara`

## Setup sekali

### 1. File `.env`

```powershell
notepad C:\Seosementara\.env
```

Salin dari `mini-pc/env.example` dan isi nilainya:

```env
DB_PASSWORD=
MASTER_ENCRYPTION_KEY=
SUPER_ADMIN_TOKEN=

CLOUDFLARE_API_KEY=
CLOUDFLARE_ACCOUNT_ID=
CLOUDFLARE_ACCOUNT_EMAIL=
```

### 2. File runtime

Salin atau clone ke `C:\Seosementara`:

```
C:\Seosementara\
  .env
  docker-compose.prod.yml
  Backend\migrations\*.up.sql
```

Clone (opsional):

```powershell
cd C:\
git clone https://github.com/WBetEngine/Seosementara.git Seosementara
Copy-Item C:\Seosementara\mini-pc\env.example C:\Seosementara\.env
notepad C:\Seosementara\.env
```

## Deploy / update

Jalankan di `C:\Seosementara`:

```powershell
cd C:\Seosementara
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d --force-recreate
curl http://localhost:8080/health
```

Docker Compose membaca variabel dari `.env` otomatis.

## Update image dari GHCR

Setelah push ke GitHub, workflow **Deploy Backend API** build image baru ke GHCR.
Di mini PC, jalankan perintah deploy di atas (`pull` + `up`).

## Admin UI (Cloudflare Workers)

Deploy terpisah — lihat `Frontend-admin/DEPLOY-CLOUDFLARE-PAGES.md`.

Pastikan `SUPER_ADMIN_TOKEN` di `.env` mini PC sama dengan token di `admin-config.js`.

## Cloudflare bootstrap

Konfigurasi Cloudflare (credentials, domain) lewat **admin panel** — tidak ada script otomatis.
