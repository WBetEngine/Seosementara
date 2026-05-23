# Deploy mini PC — GitHub pusat, tanpa file `.env`

## Arsitektur

```text
Admin Workers → GitHub Environment production → runner → Docker inject
Admin Workers → Workers Secrets (CF + GitHub PAT)
```

**Tidak perlu** buka GitHub Settings → Environments manual.

---

## Bootstrap (admin panel)

Buka: **Settings → Infra & GitHub**

### Langkah 1 — Bootstrap Platform (sekali)

| Field | Disimpan ke |
|-------|-------------|
| **GitHub PAT** | Environment `production` + Workers `GITHUB_SETUP_TOKEN` |
| **Global API Key** | Environment + Workers `CF_GLOBAL_API_KEY` |
| **Email + Account ID** | Environment + Workers |
| **SUPER_ADMIN_TOKEN** (opsional) | Environment `production` |

Klik **Simpan bootstrap** → otomatis trigger **Deploy Admin UI**.

### Langkah 2 — Infra mini PC

Isi `DB_PASSWORD`, `MASTER_ENCRYPTION_KEY` → **Simpan & deploy mini PC**  
→ Environment `production` + workflow Deploy Mini PC.

### Langkah 3 — Cloudflare (opsional update)

**Settings → Cloudflare → Koneksi** — update CF key (sync Workers + GitHub Environment jika PAT sudah ada).

---

## Mini PC sekali

```powershell
cd C:\Seosementara
git pull
.\scripts\install-github-runner.ps1   # Administrator
```

Docker + cloudflared harus sudah jalan.

---

## Deploy Workers pertama kali (tanpa PAT di GitHub)

Jika Environment `production` masih kosong, deploy Worker pertama:

```powershell
cd Frontend-admin
# Set CLOUDFLARE_* di shell, lalu:
npx wrangler deploy
```

Lalu lanjut **Bootstrap** di admin (langkah 1).

---

## Update rutin

Push `main` → GHCR build → Deploy Mini PC otomatis (runner).

Panduan arsitektur: [Plan/28-platform-github-workers.md](../Plan/28-platform-github-workers.md)
