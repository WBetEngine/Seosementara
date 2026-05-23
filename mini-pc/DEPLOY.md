# Deploy mini PC (manual)

Deploy **langsung di mini PC** dengan file `.env`. Tidak ada runner atau sync otomatis.

## Struktur folder

```
C:\Seosementara\
  .env                          ← salin dari mini-pc/env.production
  docker-compose.prod.yml
  Backend\migrations\
  mini-pc\
    env.example                 ← template kosong
    env.production              ← credential lengkap (gitignore)
    DEPLOY.md
```

## Setup sekali

```powershell
cd C:\Seosementara
Copy-Item mini-pc\env.production .env
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d --force-recreate
curl http://localhost:8080/health
```

File `mini-pc/env.production` berisi credential production lengkap (file lokal di repo, di-gitignore — tidak masuk GitHub). Alternatif: salin `mini-pc/env.example` ke `.env` dan isi manual.

## Update image

Setelah push ke GitHub, workflow **Deploy Backend API** build image baru ke GHCR. Di mini PC:

```powershell
cd C:\Seosementara
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d --force-recreate
```

## Dev lokal (laptop)

Salin `mini-pc/env.example` ke `.env` di **root repo**, uncomment bagian DEV, lalu:

```powershell
docker compose up -d --build
```

## Admin UI

Deploy Workers: lihat `Frontend-admin/DEPLOY-CLOUDFLARE-PAGES.md`.

`SUPER_ADMIN_TOKEN` di `.env` mini PC harus sama dengan token di `admin-config.js`.

## Prasyarat

- Docker Desktop
- cloudflared tunnel → `api.apidevel.org`
