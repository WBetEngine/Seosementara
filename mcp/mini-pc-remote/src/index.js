#!/usr/bin/env node
/**
 * MCP server: remote ops mini PC Seosementara via SSH (Tailscale).
 *
 * Env wajib:
 *   MINI_PC_SSH_HOST  — hostname Tailscale, mis. mini-pc.tail1234.ts.net
 *   MINI_PC_SSH_USER  — user Windows, mis. Administrator
 *
 * Env opsional:
 *   MINI_PC_SSH_PORT       — default 22
 *   MINI_PC_SSH_KEY_PATH   — path private key (OpenSSH)
 *   MINI_PC_SSH_PASSWORD   — hindari; pakai key + Tailscale SSH ACL
 *   MINI_PC_SSH_WRAPPER    — mis. "tailscale ssh" → tailscale ssh user@host --
 */
import { Client } from "ssh2";
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { z } from "zod";
import fs from "node:fs";

const PS = "powershell.exe -NoProfile -NonInteractive -Command";

function sshConfig() {
  const host = process.env.MINI_PC_SSH_HOST;
  const username = process.env.MINI_PC_SSH_USER;
  if (!host || !username) {
    throw new Error(
      "Set MINI_PC_SSH_HOST dan MINI_PC_SSH_USER (Tailscale hostname + user Windows)",
    );
  }
  const cfg = {
    host,
    port: Number(process.env.MINI_PC_SSH_PORT || 22),
    username,
    readyTimeout: 30000,
  };
  const keyPath = process.env.MINI_PC_SSH_KEY_PATH;
  if (keyPath) {
    cfg.privateKey = fs.readFileSync(keyPath);
  } else if (process.env.MINI_PC_SSH_PASSWORD) {
    cfg.password = process.env.MINI_PC_SSH_PASSWORD;
  }
  return cfg;
}

function execSsh(command, timeoutMs = 120000) {
  return new Promise((resolve, reject) => {
    const conn = new Client();
    let stdout = "";
    let stderr = "";
    const timer = setTimeout(() => {
      conn.end();
      reject(new Error(`SSH timeout setelah ${timeoutMs}ms`));
    }, timeoutMs);

    conn
      .on("ready", () => {
        conn.exec(command, (err, stream) => {
          if (err) {
            clearTimeout(timer);
            conn.end();
            reject(err);
            return;
          }
          stream
            .on("close", (code) => {
              clearTimeout(timer);
              conn.end();
              resolve({ stdout, stderr, code: code ?? 0 });
            })
            .on("data", (d) => {
              stdout += d.toString();
            });
          stream.stderr.on("data", (d) => {
            stderr += d.toString();
          });
        });
      })
      .on("error", (e) => {
        clearTimeout(timer);
        reject(e);
      })
      .connect(sshConfig());
  });
}

function ps(script) {
  const escaped = script.replace(/"/g, '\\"');
  return `${PS} "${escaped}"`;
}

function textResult(obj) {
  return {
    content: [{ type: "text", text: typeof obj === "string" ? obj : JSON.stringify(obj, null, 2) }],
  };
}

const server = new McpServer({
  name: "seosementara-mini-pc",
  version: "1.0.0",
});

server.tool(
  "mini_pc_ping",
  "Tes koneksi SSH ke mini PC (hostname, uptime)",
  {},
  async () => {
    const r = await execSsh(ps("$env:COMPUTERNAME; (Get-Date).ToString('o')"), 15000);
    return textResult({ ok: r.code === 0, code: r.code, stdout: r.stdout.trim(), stderr: r.stderr.trim() });
  },
);

server.tool(
  "mini_pc_health",
  "Cek health API Go di localhost:8080 dan api.apidevel.org",
  {},
  async () => {
    const script = [
      "$local = try { (curl.exe -sf http://localhost:8080/health) } catch { 'FAIL' }",
      "$remote = try { (curl.exe -sf https://api.apidevel.org/health) } catch { 'FAIL' }",
      "Write-Output \"local=$local remote=$remote\"",
    ].join("; ");
    const r = await execSsh(ps(script), 30000);
    return textResult({ code: r.code, stdout: r.stdout.trim(), stderr: r.stderr.trim() });
  },
);

server.tool(
  "mini_pc_docker_status",
  "docker compose ps di C:\\Seosementara",
  {},
  async () => {
    const r = await execSsh(
      ps("Set-Location C:\\Seosementara; docker compose -f docker-compose.prod.yml ps"),
      60000,
    );
    return textResult({ code: r.code, stdout: r.stdout, stderr: r.stderr });
  },
);

server.tool(
  "mini_pc_runner_status",
  "Status GitHub Actions self-hosted runner (service + proses)",
  {},
  async () => {
    const script = [
      "Get-Service | Where-Object { $_.Name -like 'actions.runner.*' } | Format-List Name,Status,StartType",
      "Get-Process -Name Runner.Listener -ErrorAction SilentlyContinue | Select-Object Id,ProcessName",
    ].join("; ");
    const r = await execSsh(ps(script), 30000);
    return textResult({ code: r.code, stdout: r.stdout, stderr: r.stderr });
  },
);

server.tool(
  "mini_pc_cloudflared_status",
  "Status service cloudflared",
  {},
  async () => {
    const r = await execSsh(
      ps("Get-Service cloudflared -ErrorAction SilentlyContinue | Format-List Name,Status,StartType"),
      30000,
    );
    return textResult({ code: r.code, stdout: r.stdout, stderr: r.stderr });
  },
);

server.tool(
  "mini_pc_install_runner_service",
  "Pasang runner Windows Service (butuh PLATFORM_GITHUB_PAT di env MCP server)",
  {},
  async () => {
    const pat = process.env.PLATFORM_GITHUB_PAT;
    if (!pat) {
      return textResult({ error: "Set PLATFORM_GITHUB_PAT di env MCP server (Cursor mcp.json)" });
    }
    const script = [
      "$env:PLATFORM_GITHUB_PAT = '" + pat.replace(/'/g, "''") + "'",
      "Set-Location C:\\Seosementara\\scripts",
      "& .\\install-github-runner-service-auto.ps1 -GitHubPat $env:PLATFORM_GITHUB_PAT",
    ].join("; ");
    const r = await execSsh(ps(script), 180000);
    return textResult({ code: r.code, stdout: r.stdout, stderr: r.stderr });
  },
);

server.tool(
  "mini_pc_deploy",
  "Jalankan mini-pc-deploy.ps1 (secret DB harus sudah di GitHub Environment)",
  {},
  async () => {
    return textResult({
      error: "Deploy penuh butuh DB_PASSWORD dari GitHub Secrets.",
      hint: "Trigger GitHub Actions Deploy Mini PC saat runner online, atau set MINI_PC_DEPLOY=1 + secrets di workflow.",
    });
  },
);

server.tool(
  "mini_pc_run_powershell",
  "Jalankan perintah PowerShell terbatas di mini PC (read-only diagnostics)",
  {
    command: z
      .string()
      .describe("Perintah PS — hanya allowlist: Get-Service, docker, curl.exe, Get-Process"),
  },
  async ({ command }) => {
    const blocked = /(Remove-|Set-Content|Invoke-WebRequest|Start-Process|format\s+c:)/i;
    if (blocked.test(command)) {
      return textResult({ error: "Perintah diblokir demi keamanan. Pakai tool khusus." });
    }
    const r = await execSsh(ps(command), 120000);
    return textResult({ code: r.code, stdout: r.stdout, stderr: r.stderr });
  },
);

server.tool(
  "mini_pc_full_diagnostics",
  "Jalankan semua cek sekaligus (health, docker, runner, cloudflared)",
  {},
  async () => {
    const scriptPath = "C:\\Seosementara\\scripts\\mini-pc-diagnostics.ps1";
    const r = await execSsh(ps(`if (Test-Path '${scriptPath}') { & '${scriptPath}' } else { Write-Output 'Script belum di-sync — jalankan Deploy Mini PC workflow' }`), 120000);
    return textResult({ code: r.code, stdout: r.stdout, stderr: r.stderr });
  },
);

async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
}

main().catch((e) => {
  console.error(e);
  process.exit(1);
});
