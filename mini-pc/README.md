# Mini PC — runtime saja (Opsi A)

Folder ini menjelaskan file **minimal** di komputer rumah. **Kode Go tetap di GitHub**; mini PC hanya menarik **image Docker** dari GHCR.

## File yang perlu ada di `C:\Seosementara`

| File / folder | Wajib? | Catatan |
|---------------|--------|---------|
| `docker-compose.prod.yml` | Ya | Dari repo (root) |
| `.env` | Ya | Secret lokal, jangan commit |
| `Backend/migrations/` | Ya (DB baru) | Hanya SQL init Postgres |
| `Backend/` source, `Frontend-admin/` | **Tidak** | Sudah dibundle di image |

## Setup sekali

### 1. Login GHCR (repo private)

Buat Personal Access Token GitHub: scope `read:packages`.

```powershell
docker login ghcr.io -u NAMA_GITHUB_ANDA
# Password = token (bukan password akun)
```

### 2. `.env` di root

Salin dari `Backend/env.example`, isi `DB_PASSWORD`, `MASTER_ENCRYPTION_KEY`, `SUPER_ADMIN_TOKEN`.

### 3. Jalankan

```powershell
cd C:\Seosementara
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d
curl http://localhost:8080/health
```

## Update setelah push ke GitHub

**Otomatis (SSH):** lihat `Frontend-admin/DEPLOY-OPSI-A.md` — set variable `MINI_PC_DEPLOY=true` + secrets.

**Manual:**

```powershell
cd C:\Seosementara
docker compose -f docker-compose.prod.yml pull api
docker compose -f docker-compose.prod.yml up -d api
```

Atau jalankan `scripts/mini-pc-pull.ps1`.
