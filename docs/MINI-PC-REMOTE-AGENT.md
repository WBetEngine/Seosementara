# Akses Mini PC dari Cursor AI Agent (Tailscale SSH + MCP)

Cloud Agent Cursor **tidak otomatis** punya akses ke PC rumah Anda.  
Agar agent bisa cek runner, Docker, dan health **tanpa RDP manual**, hubungkan lewat **Tailscale SSH** + **MCP server** repo ini.

## Prasyarat di mini PC (Windows)

1. **Tailscale** terpasang + login (SSH enabled di admin console Tailscale)
2. **OpenSSH Server** Windows aktif **atau** Tailscale SSH (`tailscale set --ssh`)
3. Hostname Tailscale dicatat, mis. `mini-pc.tailxxxxx.ts.net`

Verifikasi dari laptop Anda:

```bash
ssh Administrator@mini-pc.tailxxxxx.ts.net "hostname"
# atau
tailscale ssh Administrator@mini-pc "hostname"
```

## Pasang MCP di Cursor (Desktop)

1. Clone repo / buka folder `mcp/mini-pc-remote`
2. `npm install`
3. Cursor → **Settings → MCP → Add server**

Contoh `.cursor/mcp.json` (di mesin **Anda**, bukan di cloud):

```json
{
  "mcpServers": {
    "seosementara-mini-pc": {
      "command": "node",
      "args": ["/path/to/Seosementara/mcp/mini-pc-remote/src/index.js"],
      "env": {
        "MINI_PC_SSH_HOST": "mini-pc.tailxxxxx.ts.net",
        "MINI_PC_SSH_USER": "Administrator",
        "MINI_PC_SSH_KEY_PATH": "/home/you/.ssh/id_ed25519",
        "PLATFORM_GITHUB_PAT": "ghp_xxx_optional_for_runner_install"
      }
    }
  }
}
```

> **Jangan commit** PAT atau private key ke GitHub. Simpan hanya di MCP env lokal Cursor.

## Tools MCP yang tersedia

| Tool | Fungsi |
|------|--------|
| `mini_pc_ping` | Tes SSH |
| `mini_pc_health` | localhost:8080 + api.apidevel.org |
| `mini_pc_docker_status` | `docker compose ps` |
| `mini_pc_runner_status` | Service GitHub runner |
| `mini_pc_cloudflared_status` | Service cloudflared |
| `mini_pc_install_runner_service` | Pasang runner service (butuh PAT) |
| `mini_pc_full_diagnostics` | Script lengkap |
| `mini_pc_run_powershell` | Perintah diagnostic terbatas |

Setelah MCP aktif, minta agent: *"Jalankan mini_pc_full_diagnostics"*.

## Alternatif tanpa MCP: GitHub Actions

Workflow **Mini PC Diagnostics** jalan di self-hosted runner (butuh runner **online**):

```text
GitHub → Actions → Mini PC Diagnostics → Run workflow
```

Hasil muncul di job summary (JSON).

## Cloud Agent vs Cursor Desktop

| | Cloud Agent (background) | Cursor Desktop + MCP |
|--|---------------------------|----------------------|
| Akses Tailscale SSH | ❌ default tidak | ✅ jika MCP dikonfigurasi |
| Trigger workflow | ❌ butuh token Actions write | ✅ dari browser GitHub |
| Cek API publik | ✅ api.apidevel.org/health | ✅ |

Untuk otomatisasi penuh ala teman Anda: **runner Windows Service** + **MCP atau workflow diagnostics**.

## Checklist koneksi

```
[ ] tailscale ssh / ssh ke mini PC dari laptop OK
[ ] MCP seosementara-mini-pc hijau di Cursor
[ ] mini_pc_ping → hostname mini PC
[ ] mini_pc_health → local=ok remote=ok
[ ] runner_service → Running
[ ] cloudflared → Running
```

Lihat juga: [mini-pc/DEPLOY.md](../mini-pc/DEPLOY.md)
