# Zero-touch mini PC

**Visi:** Setelah setup sekali, Anda hanya edit GitHub — mini PC update sendiri.

## Kenapa `.env` tidak bisa di file GitHub?

GitHub **memblokir push** jika ada API key/password di kode (push protection).  
Solusi: secret disimpan di **GitHub Secrets** (menu web), workflow yang rakit file `.env` di mini PC.

## Alur setelah setup

```text
Edit kode / template di GitHub → push main
  → Build image (GHCR)
  → Runner di mini PC: copy file + buat .env dari Secrets + docker compose up
```

Anda **tidak buka** mini PC lagi.

## Setup sekali (2 langkah)

### A. Isi GitHub Secrets

**Settings → Secrets and variables → Actions → New repository secret**

Daftar lengkap: `mini-pc/GITHUB-SECRETS-LIST.md`

### B. Pasang runner di mini PC (PowerShell Admin, sekali)

```powershell
cd C:\Seosementara
# unduh install-github-runner.ps1 dari main, lalu:
.\scripts\install-github-runner.ps1
```

Token registration dari: **Settings → Actions → Runners → New self-hosted runner**

## File di repo (bukan secret)

| File | Fungsi |
|------|--------|
| `mini-pc/production.env.template` | Bentuk `.env` — nilai dari Secrets |
| `docker-compose.prod.yml` | Docker stack |
| `Backend/migrations/*.sql` | DB schema |

## Workflow utama

**Deploy Mini PC (local)** — `runs-on: self-hosted` — jalan otomatis tiap push `main`.

## Akses agent ke mini PC

Cloud agent **tidak bisa** SSH ke Tailscale IP Anda secara andal.  
Runner **self-hosted** di mini PC = GitHub "jalan di dalam" mini PC tanpa SSH dari internet.

## Sebelum rilis ke pasar

Rotasi semua Secrets; pertimbangkan repo private + hapus history secret lama.
