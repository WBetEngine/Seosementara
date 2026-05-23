/**
 * GitHub Actions secrets + workflow dispatch (admin → GitHub → Docker inject).
 */

import sodium from "tweetsodium";

const DEFAULT_REPO = "WBetEngine/Seosementara";
const DEPLOY_WORKFLOW = "deploy-mini-pc.yml";

function repoPath(env) {
  return env.GITHUB_REPO || DEFAULT_REPO;
}

async function ghFetch(env, path, init = {}) {
  const token = env.GITHUB_SETUP_TOKEN;
  if (!token) {
    throw new Error("GITHUB_SETUP_TOKEN belum dikonfigurasi di Workers Secrets (bootstrap CI)");
  }
  const res = await fetch(`https://api.github.com${path}`, {
    ...init,
    headers: {
      Accept: "application/vnd.github+json",
      Authorization: `Bearer ${token}`,
      "X-GitHub-Api-Version": "2022-11-28",
      ...(init.headers || {}),
    },
  });
  const text = await res.text();
  let body = null;
  if (text) {
    try {
      body = JSON.parse(text);
    } catch {
      body = { raw: text };
    }
  }
  if (!res.ok) {
    const msg = body?.message || body?.raw || res.statusText;
    throw new Error(`GitHub API ${res.status}: ${msg}`);
  }
  return body;
}

function encryptSecret(publicKeyB64, secretValue) {
  const messageBytes = new TextEncoder().encode(secretValue);
  const keyBytes = Uint8Array.from(atob(publicKeyB64), (c) => c.charCodeAt(0));
  const encryptedBytes = sodium.seal(messageBytes, keyBytes);
  let binary = "";
  for (let i = 0; i < encryptedBytes.length; i++) {
    binary += String.fromCharCode(encryptedBytes[i]);
  }
  return btoa(binary);
}

export async function putRepoSecret(env, name, value) {
  const [owner, repo] = repoPath(env).split("/");
  const pub = await ghFetch(env, `/repos/${owner}/${repo}/actions/secrets/public-key`);
  const encrypted = encryptSecret(pub.key, value);
  await ghFetch(env, `/repos/${owner}/${repo}/actions/secrets/${name}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      encrypted_value: encrypted,
      key_id: pub.key_id,
    }),
  });
}

export async function triggerMiniPcDeploy(env) {
  const [owner, repo] = repoPath(env).split("/");
  await ghFetch(env, `/repos/${owner}/${repo}/actions/workflows/${DEPLOY_WORKFLOW}/dispatches`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ ref: "main" }),
  });
}

export async function saveInfraSecrets(env, payload) {
  const { db_password, master_encryption_key, super_admin_token } = payload;
  if (!db_password || !master_encryption_key) {
    throw new Error("db_password dan master_encryption_key wajib");
  }
  await putRepoSecret(env, "DB_PASSWORD", db_password);
  await putRepoSecret(env, "MASTER_ENCRYPTION_KEY", master_encryption_key);
  if (super_admin_token) {
    await putRepoSecret(env, "SUPER_ADMIN_TOKEN", super_admin_token);
  }
  await triggerMiniPcDeploy(env);
  return {
    ok: true,
    secrets_updated: ["DB_PASSWORD", "MASTER_ENCRYPTION_KEY"].concat(
      super_admin_token ? ["SUPER_ADMIN_TOKEN"] : []
    ),
    deploy_triggered: true,
  };
}

export async function githubConfigured(env) {
  return Boolean(env.GITHUB_SETUP_TOKEN);
}
