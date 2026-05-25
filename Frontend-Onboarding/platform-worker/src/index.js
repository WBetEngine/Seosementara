import { corsHeaders, jsonResponse, readJson } from './cors.js';
import * as gh from './github.js';
import * as cf from './cloudflare.js';
import {
  createSession,
  getPat,
  getStatus,
  patchStatus,
  saveCfMeta,
  getCfMeta
} from './session.js';

const BASE = '/admin/api/platform';

export default {
  async fetch(request, env, ctx) {
    const url = new URL(request.url);

    if (request.method === 'OPTIONS') {
      return new Response(null, { status: 204, headers: corsHeaders(request, env) });
    }

    if (!url.pathname.startsWith(BASE)) {
      return jsonResponse({ error: 'Not found' }, 404, request, env);
    }

    const path = url.pathname.slice(BASE.length) || '/';

    try {
      if (path === '/setup/status' && request.method === 'GET') {
        return handleStatus(request, env);
      }
      if (path === '/bootstrap/cloudflare/verify' && request.method === 'POST') {
        return handleBootstrapCfVerify(request, env);
      }
      if (path === '/setup/initial' && request.method === 'POST') {
        return handleInitialSetup(request, env);
      }
      if (path === '/github/pat' && request.method === 'POST') {
        return handleGithubPat(request, env);
      }
      if (path === '/cloudflare/credentials/test' && request.method === 'POST') {
        return handleCfTest(request, env);
      }
      if (path === '/cloudflare/credentials' && request.method === 'POST') {
        return handleCfSave(request, env);
      }
      if (path === '/infra/ssh/test' && request.method === 'POST') {
        return handleSshTest(request, env);
      }
      if (path === '/github/runner/register' && request.method === 'POST') {
        return handleRunnerRegister(request, env);
      }
      if (path === '/cloudflare/tunnel/create' && request.method === 'POST') {
        return handleTunnelCreate(request, env);
      }
      if (path === '/infra/database' && request.method === 'POST') {
        return handleDatabase(request, env);
      }
      if (path === '/deploy/backend' && request.method === 'POST') {
        return handleDeploy(request, env, 'bootstrap-deploy-backend');
      }
      if (path === '/deploy/admin-pages' && request.method === 'POST') {
        return handleDeploy(request, env, 'bootstrap-deploy-admin-pages');
      }
      if (path === '/deploy/public-pages' && request.method === 'POST') {
        return handleDeploy(request, env, 'bootstrap-deploy-public-pages');
      }

      return jsonResponse({ error: 'Not found', path }, 404, request, env);
    } catch (e) {
      return jsonResponse(
        { ok: false, error: e.message || 'Internal error' },
        500,
        request,
        env
      );
    }
  }
};

async function handleStatus(request, env) {
  const status = await getStatus(env.SETUP_KV);
  return jsonResponse({ ok: true, status }, 200, request, env);
}

async function handleBootstrapCfVerify(request, env) {
  const body = await readJson(request);
  const token = body?.cf_token?.trim();
  const accountId = body?.cf_account?.trim();

  if (!token) {
    return jsonResponse({ ok: false, error: 'cf_token wajib' }, 400, request, env);
  }
  if (!accountId || !/^[a-f0-9]{32}$/i.test(accountId)) {
    return jsonResponse({ ok: false, error: 'cf_account wajib (32 hex)' }, 400, request, env);
  }

  const verify = await cf.verifyToken(token);
  if (!verify.ok) {
    return jsonResponse({ ok: false, error: verify.error }, 400, request, env);
  }

  const acc = await cf.getAccount(token, accountId);
  if (!acc.ok) {
    return jsonResponse({ ok: false, error: acc.error }, 400, request, env);
  }

  await env.SETUP_KV.put('setup:cf_token', token);
  await saveCfMeta(env.SETUP_KV, { account_id: accountId });
  await patchStatus(env.SETUP_KV, { cf_worker_ok: true, cf_account_name: acc.name });

  return jsonResponse(
    {
      ok: true,
      message: 'Token Cloudflare valid untuk deploy Platform Worker',
      account_name: acc.name
    },
    200,
    request,
    env
  );
}

async function handleInitialSetup(request, env) {
  const body = await readJson(request);
  const pat = body?.github_pat?.trim();
  const token = body?.cf_token?.trim();
  const accountId = body?.cf_account?.trim();

  if (!pat) {
    return jsonResponse({ ok: false, error: 'github_pat wajib' }, 400, request, env);
  }
  if (!token || !accountId) {
    return jsonResponse(
      { ok: false, error: 'cf_token dan cf_account dari langkah 1 wajib' },
      400,
      request,
      env
    );
  }

  const check = await gh.validatePat(pat);
  if (!check.ok) {
    return jsonResponse({ ok: false, error: check.error }, 401, request, env);
  }

  const verify = await cf.verifyToken(token);
  if (!verify.ok) {
    return jsonResponse({ ok: false, error: verify.error }, 400, request, env);
  }

  await gh.setRepoSecret(env.GITHUB_REPO, pat, 'CLOUDFLARE_API_TOKEN', token);
  await gh.setRepoSecret(env.GITHUB_REPO, pat, 'CLOUDFLARE_ACCOUNT_ID', accountId);

  await env.SETUP_KV.put('setup:cf_token', token);
  await saveCfMeta(env.SETUP_KV, { account_id: accountId });

  await gh.dispatchWorkflow(env.GITHUB_REPO, pat, 'bootstrap-deploy-platform-worker', {
    account_id: accountId
  });

  const sessionId = await createSession(env.SETUP_KV, pat, check.login);
  await patchStatus(env.SETUP_KV, {
    github_pat_ok: true,
    github_login: check.login,
    cf_worker_ok: true,
    worker_deploy_dispatched: true
  });

  const workerUrl =
    env.WORKER_PUBLIC_URL ||
    (env.WORKER_SUBDOMAIN
      ? `https://sse-platform.${env.WORKER_SUBDOMAIN}.workers.dev`
      : null);

  return jsonResponse(
    {
      ok: true,
      session_id: sessionId,
      login: check.login,
      message:
        'GitHub PAT valid. CLOUDFLARE_* disimpan ke Secrets & workflow Deploy Platform Worker dipicu. Refresh halaman setelah Actions hijau.',
      actions_url: `https://github.com/${env.GITHUB_REPO}/actions`,
      worker_url: workerUrl
    },
    200,
    request,
    env
  );
}

async function handleGithubPat(request, env) {
  const body = await readJson(request);
  const pat = body?.github_pat?.trim();
  if (!pat) {
    return jsonResponse({ ok: false, error: 'github_pat wajib' }, 400, request, env);
  }

  const check = await gh.validatePat(pat);
  if (!check.ok) {
    return jsonResponse({ ok: false, error: check.error }, 401, request, env);
  }

  const sessionId = await createSession(env.SETUP_KV, pat, check.login);
  await patchStatus(env.SETUP_KV, { github_pat_ok: true, github_login: check.login });

  return jsonResponse(
    {
      ok: true,
      session_id: sessionId,
      login: check.login,
      message: 'GitHub PAT valid'
    },
    200,
    request,
    env
  );
}

async function requirePat(request, env, body) {
  const sessionId = request.headers.get('X-Setup-Session') || body?.session_id;
  let pat = body?.github_pat?.trim();
  if (!pat && sessionId) {
    pat = await getPat(env.SETUP_KV, sessionId);
  }
  if (!pat) {
    return { error: jsonResponse({ ok: false, error: 'Butuh GitHub PAT atau X-Setup-Session' }, 401, request, env) };
  }
  return { pat, sessionId };
}

async function handleCfTest(request, env) {
  const body = await readJson(request);
  const auth = await requirePat(request, env, body);
  if (auth.error) return auth.error;

  const token = body?.cf_token?.trim();
  const accountId = body?.cf_account?.trim();
  const zoneId = body?.cf_zone?.trim();
  const domain = body?.primary_domain?.trim();

  if (!token) {
    return jsonResponse({ ok: false, error: 'cf_token wajib' }, 400, request, env);
  }

  const verify = await cf.verifyToken(token);
  if (!verify.ok) {
    return jsonResponse({ ok: false, error: verify.error }, 400, request, env);
  }

  let accountName = null;
  let zoneName = null;
  let zoneStatus = null;

  if (accountId) {
    const acc = await cf.getAccount(token, accountId);
    if (!acc.ok) {
      return jsonResponse({ ok: false, error: acc.error }, 400, request, env);
    }
    accountName = acc.name;
  }

  if (zoneId) {
    const zone = await cf.getZone(token, zoneId);
    if (!zone.ok) {
      return jsonResponse({ ok: false, error: zone.error }, 400, request, env);
    }
    zoneName = zone.name;
    zoneStatus = zone.status;
  }

  return jsonResponse(
    {
      ok: true,
      message: 'Token Cloudflare valid',
      account_name: accountName,
      zone_name: zoneName,
      zone_status: zoneStatus,
      primary_domain: domain
    },
    200,
    request,
    env
  );
}

async function handleCfSave(request, env) {
  const body = await readJson(request);
  const auth = await requirePat(request, env, body);
  if (auth.error) return auth.error;

  const token = body?.cf_token?.trim();
  if (!token) {
    return jsonResponse({ ok: false, error: 'cf_token wajib' }, 400, request, env);
  }

  const verify = await cf.verifyToken(token);
  if (!verify.ok) {
    return jsonResponse({ ok: false, error: verify.error }, 400, request, env);
  }

  const accountId = body?.cf_account?.trim();
  const zoneId = body?.cf_zone?.trim();

  await env.SETUP_KV.put('setup:cf_token', token);
  await saveCfMeta(env.SETUP_KV, {
    account_id: accountId,
    zone_id: zoneId,
    primary_domain: body?.primary_domain?.trim()
  });

  await gh.setRepoSecret(env.GITHUB_REPO, auth.pat, 'CLOUDFLARE_API_TOKEN', token);
  if (accountId) {
    await gh.setRepoSecret(env.GITHUB_REPO, auth.pat, 'CLOUDFLARE_ACCOUNT_ID', accountId);
  }

  await patchStatus(env.SETUP_KV, { cloudflare_ok: true });

  return jsonResponse(
    {
      ok: true,
      message: 'Kredensial Cloudflare tersimpan (KV + GitHub Secrets untuk Actions)'
    },
    200,
    request,
    env
  );
}

async function handleSshTest(request, env) {
  const body = await readJson(request);
  const auth = await requirePat(request, env, body);
  if (auth.error) return auth.error;

  const host = body?.ssh_host?.trim();
  const port = body?.ssh_port?.trim() || '22';
  const user = body?.ssh_user?.trim();
  const secret = body?.ssh_secret?.trim();

  if (!host || !user || !secret) {
    return jsonResponse({ ok: false, error: 'ssh_host, ssh_user, ssh_secret wajib' }, 400, request, env);
  }

  await gh.setRepoSecret(env.GITHUB_REPO, auth.pat, 'BOOTSTRAP_SSH_HOST', host);
  await gh.setRepoSecret(env.GITHUB_REPO, auth.pat, 'BOOTSTRAP_SSH_PORT', port);
  await gh.setRepoSecret(env.GITHUB_REPO, auth.pat, 'BOOTSTRAP_SSH_USER', user);
  await gh.setRepoSecret(env.GITHUB_REPO, auth.pat, 'BOOTSTRAP_SSH_PASSWORD', secret);

  await gh.dispatchWorkflow(env.GITHUB_REPO, auth.pat, 'bootstrap-ssh-test', {
    host,
    port,
    user
  });

  await patchStatus(env.SETUP_KV, { ssh_ok: true, ssh_test_dispatched: true });

  return jsonResponse(
    {
      ok: true,
      message: 'Test SSH dijalankan via GitHub Actions (bootstrap-ssh-test). Cek tab Actions.',
      actions_url: `https://github.com/${env.GITHUB_REPO}/actions`
    },
    200,
    request,
    env
  );
}

async function handleRunnerRegister(request, env) {
  const body = await readJson(request);
  const auth = await requirePat(request, env, body);
  if (auth.error) return auth.error;

  const label = body?.runner_label?.trim() || 'mini-pc';

  await gh.dispatchWorkflow(env.GITHUB_REPO, auth.pat, 'bootstrap-register-runner', {
    runner_label: label
  });

  await patchStatus(env.SETUP_KV, { runner_ok: true, runner_dispatched: true });

  return jsonResponse(
    {
      ok: true,
      message: 'Workflow register runner dipicu. Ikuti log di GitHub Actions.',
      actions_url: `https://github.com/${env.GITHUB_REPO}/actions`
    },
    200,
    request,
    env
  );
}

async function handleTunnelCreate(request, env) {
  const body = await readJson(request);
  const auth = await requirePat(request, env, body);
  if (auth.error) return auth.error;

  const token = (await env.SETUP_KV.get('setup:cf_token')) || body?.cf_token?.trim();
  const meta = await getCfMeta(env.SETUP_KV);
  const accountId = body?.cf_account?.trim() || meta?.account_id;
  const tunnelName = body?.tunnel_name?.trim() || 'sse-production';

  if (!token || !accountId) {
    return jsonResponse(
      { ok: false, error: 'Simpan Cloudflare credentials dulu (langkah 3 — Zone & domain)' },
      400,
      request,
      env
    );
  }

  const created = await cf.createTunnel(token, accountId, tunnelName);
  if (!created.ok) {
    return jsonResponse({ ok: false, error: created.error }, 400, request, env);
  }

  await patchStatus(env.SETUP_KV, {
    tunnel_ok: true,
    tunnel_id: created.tunnel_id,
    tunnel_name: created.tunnel_name
  });

  await gh.dispatchWorkflow(env.GITHUB_REPO, auth.pat, 'bootstrap-install-tunnel', {
    tunnel_id: created.tunnel_id,
    tunnel_name: created.tunnel_name
  });

  return jsonResponse(
    {
      ok: true,
      message: 'Tunnel dibuat via Cloudflare API. Install connector via workflow bootstrap-install-tunnel.',
      tunnel_id: created.tunnel_id,
      tunnel_name: created.tunnel_name
    },
    200,
    request,
    env
  );
}

async function handleDatabase(request, env) {
  const body = await readJson(request);
  const auth = await requirePat(request, env, body);
  if (auth.error) return auth.error;

  const dbPassword = body?.db_password;
  const masterKey = body?.master_key;

  if (!dbPassword || dbPassword.length < 12) {
    return jsonResponse({ ok: false, error: 'db_password minimal 12 karakter' }, 400, request, env);
  }
  if (!masterKey || masterKey.length < 16) {
    return jsonResponse({ ok: false, error: 'master_key minimal 16 karakter' }, 400, request, env);
  }

  await gh.setRepoSecret(env.GITHUB_REPO, auth.pat, 'DB_PASSWORD', dbPassword);
  await gh.setRepoSecret(env.GITHUB_REPO, auth.pat, 'MASTER_ENCRYPTION_KEY', masterKey);

  await patchStatus(env.SETUP_KV, { database_ok: true });

  return jsonResponse(
    { ok: true, message: 'DB_PASSWORD dan MASTER_ENCRYPTION_KEY disimpan ke GitHub Secrets' },
    200,
    request,
    env
  );
}

async function handleDeploy(request, env, eventType) {
  const body = await readJson(request);
  const auth = await requirePat(request, env, body);
  if (auth.error) return auth.error;

  await gh.dispatchWorkflow(env.GITHUB_REPO, auth.pat, eventType, body || {});

  const patch =
    eventType === 'bootstrap-deploy-backend'
      ? { deploy_backend_ok: true }
      : eventType === 'bootstrap-deploy-admin-pages'
        ? { deploy_pages_ok: true }
        : { deploy_public_ok: true };

  await patchStatus(env.SETUP_KV, patch);

  const status = await getStatus(env.SETUP_KV);
  if (status.github_pat_ok && status.cloudflare_ok && status.database_ok) {
    await patchStatus(env.SETUP_KV, { bootstrap_complete: true });
  }

  return jsonResponse(
    {
      ok: true,
      message: `Workflow ${eventType} dipicu`,
      actions_url: `https://github.com/${env.GITHUB_REPO}/actions`
    },
    200,
    request,
    env
  );
}
