#!/usr/bin/env bash
# Dipanggil dari GitHub Actions sebelum wrangler deploy
set -euo pipefail

cd "$(dirname "$0")/.."

if [ -z "${CLOUDFLARE_API_TOKEN:-}" ]; then
  echo "::error::CLOUDFLARE_API_TOKEN kosong. Isi GitHub Secret CLOUDFLARE_API_TOKEN atau workflow input cloudflare_api_token."
  exit 1
fi

if [ -z "${CLOUDFLARE_ACCOUNT_ID:-}" ]; then
  echo "::error::CLOUDFLARE_ACCOUNT_ID kosong. Isi secret atau workflow input cloudflare_account_id."
  exit 1
fi

export CLOUDFLARE_API_TOKEN
export CLOUDFLARE_ACCOUNT_ID

parse_kv_id() {
  node -e "
    const raw = process.argv[1] || '';
    try {
      const arr = JSON.parse(raw);
      if (!Array.isArray(arr)) process.exit(0);
      const n = arr.find((x) => x && x.title === 'SETUP_KV');
      if (n && n.id) process.stdout.write(n.id);
    } catch (e) {}
  " "$1"
}

if [ -n "${PLATFORM_KV_ID:-}" ]; then
  sed -i "s/PLACEHOLDER_KV_ID/${PLATFORM_KV_ID}/" wrangler.toml
  echo "KV: memakai secret PLATFORM_KV_ID"
elif ! grep -q 'PLACEHOLDER_KV_ID' wrangler.toml; then
  echo "KV: id sudah ada di wrangler.toml"
else
  echo "KV: mencari namespace SETUP_KV..."
  LIST_RAW=$(npx wrangler kv namespace list 2>/dev/null || true)
  KV_ID=$(parse_kv_id "$LIST_RAW")

  if [ -z "$KV_ID" ]; then
    echo "KV: membuat namespace SETUP_KV (--update-config)..."
    npx wrangler kv namespace create SETUP_KV --update-config 2>&1 || true
  fi

  if grep -q 'PLACEHOLDER_KV_ID' wrangler.toml; then
    LIST_RAW=$(npx wrangler kv namespace list 2>/dev/null || true)
    KV_ID=$(parse_kv_id "$LIST_RAW")
    if [ -n "$KV_ID" ]; then
      sed -i "s/PLACEHOLDER_KV_ID/${KV_ID}/" wrangler.toml
      echo "KV: memakai id $KV_ID"
    fi
  fi
fi

if grep -q 'PLACEHOLDER_KV_ID' wrangler.toml; then
  echo "::error::KV SETUP_KV belum siap."
  echo "Token Cloudflare perlu izin: Account Settings Read, Workers Scripts Edit, Workers KV Storage Edit."
  echo "Manual: cd Frontend-Onboarding/platform-worker && npx wrangler kv namespace create SETUP_KV --update-config"
  echo "Simpan id di GitHub → Settings → Secrets → PLATFORM_KV_ID"
  exit 1
fi

FINAL_ID=$(grep -E '^id = ' wrangler.toml | grep -oE '[a-f0-9]{32}' | head -1)
echo "::notice title=SETUP_KV::KV id=${FINAL_ID} — opsional: simpan sebagai secret PLATFORM_KV_ID"
