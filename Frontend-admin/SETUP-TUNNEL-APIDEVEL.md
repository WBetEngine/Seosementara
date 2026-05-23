# Tunnel Cloudflare — api.apidevel.org → mini PC

Domain **apidevel.org** sudah Active di Cloudflare. Panduan ini menghubungkan **API Go** di mini PC (`localhost:8080`) ke internet via **Named Tunnel**.

## Arsitektur

```text
Admin UI (Workers)  ──HTTPS──►  api.apidevel.org  ──Tunnel──►  cloudflared (Windows)  ──►  localhost:8080
```

| Hostname | Fungsi |
|----------|--------|
| `api.apidevel.org` | Backend Go (Tunnel) |
| `admin.apidevel.org` | (Opsional nanti) Custom domain Workers |
| `seosementara.*.workers.dev` | Admin UI sekarang |

---

## Bagian A — Cloudflare Zero Trust (dashboard)

### 1. Buka Zero Trust

1. [dash.cloudflare.com](https://dash.cloudflare.com) → pilih **apidevel.org**
2. Menu kiri: **Zero Trust** (atau [one.dash.cloudflare.com](https://one.dash.cloudflare.com))
3. **Networks** → **Tunnels** → **Create a tunnel**

### 2. Buat tunnel

1. Pilih connector: **Cloudflared**
2. Nama tunnel: `seosementara-api`
3. **Save tunnel**
4. Di halaman install, pilih **Windows**
5. **Salin perintah** `cloudflared.exe service install eyJh...` (token panjang) — dipakai di Bagian B

### 3. Public Hostname (route)

Masih di wizard tunnel (atau **Tunnels → seosementara-api → Public Hostname**):

| Field | Nilai |
|-------|--------|
| Subdomain | `api` |
| Domain | `apidevel.org` |
| Type | **HTTP** |
| URL | `localhost:8080` |

Simpan. Cloudflare otomatis buat DNS `api.apidevel.org` → tunnel.

### 4. Verifikasi di DNS

**Websites → apidevel.org → DNS**:

Harus ada record **CNAME** `api` → `{tunnel-id}.cfargotunnel.com` (proxied).

---

## Bagian B — Mini PC Windows

### Prasyarat

```powershell
cd C:\Seosementara
docker compose -f docker-compose.prod.yml ps
curl.exe http://localhost:8080/health
```

Harus **`ok`** sebelum tunnel.

### 1. Install cloudflared

PowerShell **Administrator**:

```powershell
cd C:\Seosementara
.\scripts\install-cloudflared.ps1
```

Atau manual: unduh dari [developers.cloudflare.com/cloudflare-one/connections/connect-networks/downloads/](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/downloads/) → `cloudflared-windows-amd64.exe` → letakkan di `C:\Program Files\cloudflared\cloudflared.exe`

### 2. Install sebagai service (token dari dashboard)

PowerShell **Administrator** — ganti `TOKEN_DARI_DASHBOARD` dengan token dari langkah A.2:

```powershell
& "C:\Program Files\cloudflared\cloudflared.exe" service install TOKEN_DARI_DASHBOARD
Start-Service cloudflared
Get-Service cloudflared
```

Status harus **Running**.

### 3. Tes dari internet

```powershell
curl.exe https://api.apidevel.org/health
```

Harus: **`ok`**

---

## Bagian C — Sambungkan Admin UI (Workers)

### 1. Konfigurasi API base

File `Frontend-admin/public/static/js/admin-config.js` (salin dari `admin-config.example.js`):

```javascript
window.SEOSEMENTARA_API_BASE = "https://api.apidevel.org";
window.SEOSEMENTARA_SUPER_ADMIN_TOKEN = "SAMA_DENGAN_SUPER_ADMIN_TOKEN_DI_ENV_MINI_PC";
```

**Jangan commit** `admin-config.js` jika berisi token (sudah di `.gitignore`).

### 2. Deploy Workers

```powershell
cd Frontend-admin
npm run deploy
```

### 3. Tes admin

Buka admin Workers → **Settings → Cloudflare → Koneksi** — request HTMX harus ke `https://api.apidevel.org/api/...`

---

## Troubleshooting

| Gejala | Penyebab | Solusi |
|--------|----------|--------|
| 502 / error tunnel | API Docker mati | `docker compose ... up -d` |
| 404 di `/health` | Route hostname salah | URL harus `localhost:8080` |
| SSL error | DNS belum propagate | Tunggu 5–15 menit |
| 401 di admin API | Token tidak cocok | Samakan `SUPER_ADMIN_TOKEN` (.env) dan `admin-config.js` |
| cloudflared stopped | Service belum start | `Start-Service cloudflared` |

---

## Checklist

- [ ] Tunnel **seosementara-api** Active (connector hijau di Zero Trust)
- [ ] `curl https://api.apidevel.org/health` → `ok`
- [ ] `admin-config.js` + deploy Workers
- [ ] Tab Cloudflare di admin memanggil API (bukan mock)
