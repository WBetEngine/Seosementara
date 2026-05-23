import {
  saveInfraSecrets,
  saveBootstrap,
  githubConfigured,
  syncCloudflareToGitHub,
} from "./github.js";
import {
  saveCloudflareCredentials,
  cloudflareConfigured,
  putWorkerSecret,
} from "./cloudflare.js";

const JSON_HEADERS = { "Content-Type": "application/json; charset=utf-8" };

function json(data, status = 200) {
  return new Response(JSON.stringify(data), { status, headers: JSON_HEADERS });
}

async function readJson(request) {
  try {
    return await request.json();
  } catch {
    return null;
  }
}

async function handleAdminAPI(request, env) {
  const url = new URL(request.url);
  const path = url.pathname;

  if (request.method === "GET" && path === "/admin/api/platform/setup/status") {
    return json({
      github_token_stored: await githubConfigured(env),
      cloudflare_configured: await cloudflareConfigured(env),
      github_repo: env.GITHUB_REPO || "WBetEngine/Seosementara",
      github_environment: env.GITHUB_ENVIRONMENT || "production",
      worker_script: env.WORKER_SCRIPT_NAME || "seosementara",
    });
  }

  if (request.method === "POST" && path === "/admin/api/platform/bootstrap") {
    const body = await readJson(request);
    if (!body) return json({ error: "JSON invalid" }, 400);
    try {
      const result = await saveBootstrap(env, body, putWorkerSecret);
      return json(result);
    } catch (e) {
      return json({ error: e.message }, 400);
    }
  }

  if (request.method === "POST" && path === "/admin/api/platform/infra") {
    const body = await readJson(request);
    if (!body) return json({ error: "JSON invalid" }, 400);
    try {
      const result = await saveInfraSecrets(env, body);
      return json(result);
    } catch (e) {
      return json({ error: e.message }, 400);
    }
  }

  if (request.method === "POST" && path === "/admin/api/platform/cloudflare/credentials") {
    const body = await readJson(request);
    if (!body) return json({ error: "JSON invalid" }, 400);
    try {
      const syncFn = (await githubConfigured(env)) ? syncCloudflareToGitHub : null;
      const result = await saveCloudflareCredentials(env, body, syncFn);
      return json(result);
    } catch (e) {
      return json({ error: e.message }, 400);
    }
  }

  return json({ error: "not found" }, 404);
}

export default {
  async fetch(request, env) {
    const url = new URL(request.url);
    if (url.pathname.startsWith("/admin/api/platform/")) {
      return handleAdminAPI(request, env);
    }
    return env.ASSETS.fetch(request);
  },
};
