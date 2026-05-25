const SESSION_TTL_SEC = 3600;

export async function createSession(kv, pat, login) {
  const id = crypto.randomUUID();
  await kv.put(
    `session:${id}`,
    JSON.stringify({
      login,
      pat_hint: pat.slice(0, 7) + '…',
      created: Date.now()
    }),
    { expirationTtl: SESSION_TTL_SEC }
  );
  await kv.put(`session_pat:${id}`, pat, { expirationTtl: SESSION_TTL_SEC });
  return id;
}

export async function getPat(kv, sessionId) {
  if (!sessionId) return null;
  return kv.get(`session_pat:${sessionId}`);
}

export async function getStatus(kv) {
  const raw = await kv.get('setup:status');
  if (!raw) {
    return defaultStatus();
  }
  try {
    return { ...defaultStatus(), ...JSON.parse(raw) };
  } catch {
    return defaultStatus();
  }
}

export async function patchStatus(kv, patch) {
  const current = await getStatus(kv);
  const next = { ...current, ...patch, updated_at: new Date().toISOString() };
  await kv.put('setup:status', JSON.stringify(next));
  return next;
}

function defaultStatus() {
  return {
    bootstrap_complete: false,
    cf_worker_ok: false,
    github_pat_ok: false,
    cloudflare_ok: false,
    ssh_ok: false,
    runner_ok: false,
    tunnel_ok: false,
    database_ok: false,
    deploy_backend_ok: false,
    deploy_pages_ok: false,
    worker_url: null
  };
}

export async function saveCfMeta(kv, meta) {
  await kv.put('setup:cf_meta', JSON.stringify(meta));
}

export async function getCfMeta(kv) {
  const raw = await kv.get('setup:cf_meta');
  return raw ? JSON.parse(raw) : null;
}
