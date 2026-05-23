/**
 * GitHub Environment secrets + workflow dispatch.
 * Token: Workers Secret GITHUB_SETUP_TOKEN, atau PAT sekali dari form bootstrap.
 */

import sodium from "tweetsodium";

const DEFAULT_REPO = "WBetEngine/Seosementara";
const DEFAULT_ENV = "production";
const DEPLOY_WORKFLOW = "deploy-mini-pc.yml";
const DEPLOY_ADMIN_WORKFLOW = "deploy-admin.yml";

function repoPath(env) {
  return env.GITHUB_REPO || DEFAULT_REPO;
}

function ghEnvName(env) {
  return env.GITHUB_ENVIRONMENT || DEFAULT_ENV;
}

function splitRepo(path) {
  const [owner, repo] = path.split("/");
  return { owner, repo };
}

async function ghFetch(token, path, init = {}) {
  if (!token) {
    throw new Error("GitHub token belum ada — isi di Bootstrap Platform (PAT)");
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

function resolveToken(env, override) {
  return override || env.GITHUB_SETUP_TOKEN || "";
}

export async function ensureEnvironment(token, env, envName) {
  const { owner, repo } = splitRepo(repoPath(env));
  await ghFetch(token, `/repos/${owner}/${repo}/environments/${envName}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({}),
  });
}

export async function putEnvironmentSecret(env, name, value, tokenOverride) {
  const token = resolveToken(env, tokenOverride);
  const envName = ghEnvName(env);
  const { owner, repo } = splitRepo(repoPath(env));
  await ensureEnvironment(token, env, envName);
  const pub = await ghFetch(
    token,
    `/repos/${owner}/${repo}/environments/${envName}/secrets/public-key`
  );
  const encrypted = encryptSecret(pub.key, value);
  await ghFetch(token, `/repos/${owner}/${repo}/environments/${envName}/secrets/${name}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      encrypted_value: encrypted,
      key_id: pub.key_id,
    }),
  });
}

export async function triggerWorkflow(env, workflowFile, tokenOverride) {
  const token = resolveToken(env, tokenOverride);
  const { owner, repo } = splitRepo(repoPath(env));
  await ghFetch(token, `/repos/${owner}/${repo}/actions/workflows/${workflowFile}/dispatches`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ ref: "main" }),
  });
}

export async function triggerMiniPcDeploy(env, tokenOverride) {
  await triggerWorkflow(env, DEPLOY_WORKFLOW, tokenOverride);
}

export async function triggerAdminDeploy(env, tokenOverride) {
  await triggerWorkflow(env, DEPLOY_ADMIN_WORKFLOW, tokenOverride);
}

/** Bootstrap: PAT + Cloudflare → GitHub Environment production + Workers Secrets */
export async function saveBootstrap(env, payload, putWorkerSecretFn) {
  const {
    github_pat,
    global_api_key,
    account_email,
    account_id,
    super_admin_token,
  } = payload;

  if (!github_pat) throw new Error("github_pat wajib (PAT: repo secrets + actions write)");
  if (!global_api_key || !account_email || !account_id) {
    throw new Error("global_api_key, account_email, account_id wajib");
  }

  const scriptName = env.WORKER_SCRIPT_NAME || "seosementara";
  const updated = [];

  await putEnvironmentSecret(env, "GITHUB_SETUP_TOKEN", github_pat, github_pat);
  updated.push("GITHUB_SETUP_TOKEN");
  await putEnvironmentSecret(env, "CLOUDFLARE_API_KEY", global_api_key, github_pat);
  updated.push("CLOUDFLARE_API_KEY");
  await putEnvironmentSecret(env, "CLOUDFLARE_ACCOUNT_EMAIL", account_email, github_pat);
  updated.push("CLOUDFLARE_ACCOUNT_EMAIL");
  await putEnvironmentSecret(env, "CLOUDFLARE_ACCOUNT_ID", account_id, github_pat);
  updated.push("CLOUDFLARE_ACCOUNT_ID");
  if (super_admin_token) {
    await putEnvironmentSecret(env, "SUPER_ADMIN_TOKEN", super_admin_token, github_pat);
    updated.push("SUPER_ADMIN_TOKEN");
  }

  if (putWorkerSecretFn) {
    await putWorkerSecretFn(account_id, account_email, global_api_key, scriptName, "GITHUB_SETUP_TOKEN", github_pat);
    await putWorkerSecretFn(account_id, account_email, global_api_key, scriptName, "CF_GLOBAL_API_KEY", global_api_key);
    await putWorkerSecretFn(account_id, account_email, global_api_key, scriptName, "CF_ACCOUNT_EMAIL", account_email);
    await putWorkerSecretFn(account_id, account_email, global_api_key, scriptName, "CF_ACCOUNT_ID", account_id);
  }

  await triggerAdminDeploy(env, github_pat);

  return {
    ok: true,
    environment: ghEnvName(env),
    secrets_updated: updated,
    worker_secrets_updated: putWorkerSecretFn
      ? ["GITHUB_SETUP_TOKEN", "CF_GLOBAL_API_KEY", "CF_ACCOUNT_EMAIL", "CF_ACCOUNT_ID"]
      : [],
    deploy_admin_triggered: true,
    message: "Bootstrap tersimpan di GitHub Environment production + Workers Secrets",
  };
}

export async function saveInfraSecrets(env, payload) {
  const { db_password, master_encryption_key, super_admin_token } = payload;
  if (!db_password || !master_encryption_key) {
    throw new Error("db_password dan master_encryption_key wajib");
  }
  const updated = ["DB_PASSWORD", "MASTER_ENCRYPTION_KEY"];
  await putEnvironmentSecret(env, "DB_PASSWORD", db_password);
  await putEnvironmentSecret(env, "MASTER_ENCRYPTION_KEY", master_encryption_key);
  if (super_admin_token) {
    await putEnvironmentSecret(env, "SUPER_ADMIN_TOKEN", super_admin_token);
    updated.push("SUPER_ADMIN_TOKEN");
  }
  await triggerMiniPcDeploy(env);
  return {
    ok: true,
    environment: ghEnvName(env),
    secrets_updated: updated,
    deploy_triggered: true,
  };
}

/** Sync Cloudflare deploy secrets ke GitHub Environment (untuk wrangler CI) */
export async function syncCloudflareToGitHub(env, { global_api_key, account_email, account_id }) {
  const updated = [];
  if (global_api_key) {
    await putEnvironmentSecret(env, "CLOUDFLARE_API_KEY", global_api_key);
    updated.push("CLOUDFLARE_API_KEY");
  }
  if (account_email) {
    await putEnvironmentSecret(env, "CLOUDFLARE_ACCOUNT_EMAIL", account_email);
    updated.push("CLOUDFLARE_ACCOUNT_EMAIL");
  }
  if (account_id) {
    await putEnvironmentSecret(env, "CLOUDFLARE_ACCOUNT_ID", account_id);
    updated.push("CLOUDFLARE_ACCOUNT_ID");
  }
  return updated;
}

export async function githubConfigured(env) {
  return Boolean(env.GITHUB_SETUP_TOKEN);
}

/** Cek self-hosted runner via GitHub API (butuh PAT di Worker). */
export async function getRunnerStatus(env) {
  const token = resolveToken(env);
  if (!token) {
    return {
      checked: false,
      needs_runner_setup: true,
      reason: "github_pat_missing",
      online_count: 0,
      total_count: 0,
      runners: [],
    };
  }
  try {
    const { owner, repo } = splitRepo(repoPath(env));
    const data = await ghFetch(token, `/repos/${owner}/${repo}/actions/runners?per_page=30`);
    const runners = (data.runners || []).map((r) => ({
      name: r.name,
      status: r.status,
      os: r.os,
      busy: r.busy,
    }));
    const online = runners.filter((r) => r.status === "online");
    return {
      checked: true,
      needs_runner_setup: runners.length === 0 || online.length === 0,
      reason: runners.length === 0 ? "no_runners" : online.length === 0 ? "runners_offline" : null,
      online_count: online.length,
      total_count: runners.length,
      runners,
    };
  } catch (e) {
    return {
      checked: false,
      needs_runner_setup: true,
      reason: "api_error",
      error: e.message,
      online_count: 0,
      total_count: 0,
      runners: [],
    };
  }
}

/** Nama secret di Environment production (tanpa nilai). */
export async function listEnvironmentSecretNames(env) {
  const token = resolveToken(env);
  if (!token) return [];
  try {
    const envName = ghEnvName(env);
    const { owner, repo } = splitRepo(repoPath(env));
    const data = await ghFetch(
      token,
      `/repos/${owner}/${repo}/environments/${envName}/secrets?per_page=50`
    );
    return (data.secrets || []).map((s) => s.name);
  } catch {
    return [];
  }
}

export async function getPlatformSetupStatus(env) {
  const github_token_stored = await githubConfigured(env);
  const secretNames = github_token_stored ? await listEnvironmentSecretNames(env) : [];
  const runner = await getRunnerStatus(env);
  const envName = ghEnvName(env);
  const repo = repoPath(env);

  return {
    github_token_stored,
    github_repo: repo,
    github_environment: envName,
    environment_secrets: secretNames,
    bootstrap_complete: secretNames.includes("GITHUB_SETUP_TOKEN") && secretNames.includes("CLOUDFLARE_API_KEY"),
    infra_complete: secretNames.includes("DB_PASSWORD") && secretNames.includes("MASTER_ENCRYPTION_KEY"),
    runner,
    runner_install: {
      script_url: `https://raw.githubusercontent.com/${repo}/main/scripts/install-github-runner.ps1`,
      runners_new_url: `https://github.com/${repo}/settings/actions/runners/new`,
      expected_runner_name: "mini-pc-seosementara",
    },
  };
}
