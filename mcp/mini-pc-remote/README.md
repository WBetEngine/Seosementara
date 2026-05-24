# MCP — remote mini PC via SSH (Tailscale)

Jalankan perintah diagnostic/deploy di mini PC Windows dari Cursor AI.

## Setup

```bash
npm install
```

Lihat panduan lengkap: [docs/MINI-PC-REMOTE-AGENT.md](../docs/MINI-PC-REMOTE-AGENT.md)

## Env

| Variable | Wajib | Keterangan |
|----------|-------|------------|
| `MINI_PC_SSH_HOST` | Ya | Hostname Tailscale |
| `MINI_PC_SSH_USER` | Ya | User Windows |
| `MINI_PC_SSH_KEY_PATH` | Disarankan | Private key OpenSSH |
| `PLATFORM_GITHUB_PAT` | Opsional | Untuk `mini_pc_install_runner_service` |

## Cursor

Konfigurasi produksi: `.cursor/mcp.json` (gitignored — isi host/user/password mini PC).
