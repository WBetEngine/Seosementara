/**
 * Cloudflare API — simpan Global API Key ke Workers Secrets via CF API.
 */

function cfHeaders(email, apiKey) {
  return {
    "Content-Type": "application/json",
    "X-Auth-Email": email,
    "X-Auth-Key": apiKey,
  };
}

async function cfFetch(accountId, email, apiKey, path, init = {}) {
  const res = await fetch(`https://api.cloudflare.com/client/v4${path}`, {
    ...init,
    headers: { ...cfHeaders(email, apiKey), ...(init.headers || {}) },
  });
  const body = await res.json();
  if (!body.success) {
    const err = body.errors?.[0]?.message || res.statusText;
    throw new Error(`Cloudflare API: ${err}`);
  }
  return body.result;
}

export async function testCloudflareCredentials(authType, payload) {
  if (authType === "global_api_key") {
    const { global_api_key, account_email } = payload;
    if (!global_api_key || !account_email) {
      throw new Error("global_api_key dan account_email wajib");
    }
    const res = await fetch("https://api.cloudflare.com/client/v4/user", {
      headers: cfHeaders(account_email, global_api_key),
    });
    const body = await res.json();
    if (!body.success) {
      throw new Error(body.errors?.[0]?.message || "Global API Key tidak valid");
    }
    return { ok: true, auth_type: authType };
  }
  const { api_token } = payload;
  if (!api_token) throw new Error("api_token wajib");
  const res = await fetch("https://api.cloudflare.com/client/v4/user/tokens/verify", {
    headers: { Authorization: `Bearer ${api_token}` },
  });
  const body = await res.json();
  if (!body.success) {
    throw new Error(body.errors?.[0]?.message || "Token tidak valid");
  }
  return { ok: true, auth_type: authType };
}

export async function putWorkerSecret(accountId, email, apiKey, scriptName, secretName, secretValue) {
  const res = await fetch(
    `https://api.cloudflare.com/client/v4/accounts/${accountId}/workers/scripts/${scriptName}/secrets`,
    {
      method: "PUT",
      headers: cfHeaders(email, apiKey),
      body: JSON.stringify({
        name: secretName,
        text: secretValue,
        type: "secret_text",
      }),
    }
  );
  const body = await res.json();
  if (!body.success) {
    const err = body.errors?.[0]?.message || res.statusText;
    throw new Error(`Simpan Workers Secret: ${err}`);
  }
  return body.result;
}

export async function saveCloudflareCredentials(env, payload, syncGitHubFn) {
  const authType = payload.auth_type || "global_api_key";
  await testCloudflareCredentials(authType, payload);

  const accountId = payload.account_id || env.CF_ACCOUNT_ID;
  const scriptName = env.WORKER_SCRIPT_NAME || "seosementara";
  if (!accountId) {
    throw new Error("account_id wajib");
  }

  let githubSynced = [];

  if (authType === "global_api_key") {
    const { global_api_key, account_email } = payload;
    await putWorkerSecret(accountId, account_email, global_api_key, scriptName, "CF_GLOBAL_API_KEY", global_api_key);
    await putWorkerSecret(accountId, account_email, global_api_key, scriptName, "CF_ACCOUNT_EMAIL", account_email);
    await putWorkerSecret(accountId, account_email, global_api_key, scriptName, "CF_ACCOUNT_ID", accountId);
    if (syncGitHubFn) {
      githubSynced = await syncGitHubFn(env, {
        global_api_key,
        account_email,
        account_id: accountId,
      });
    }
  } else {
    const { api_token } = payload;
    const email = payload.account_email || env.CF_BOOTSTRAP_EMAIL || "";
    const bootstrapKey = env.CF_BOOTSTRAP_GLOBAL_KEY || "";
    if (!email || !bootstrapKey) {
      throw new Error("Simpan API Token butuh CF_BOOTSTRAP_GLOBAL_KEY + email di Worker (sekali, via CI)");
    }
    await putWorkerSecret(accountId, email, bootstrapKey, scriptName, "CF_API_TOKEN", api_token);
  }

  return {
    ok: true,
    configured: true,
    auth_type: authType,
    account_id: accountId,
    github_environment_synced: githubSynced,
    message: "Workers Secrets" + (githubSynced.length ? " + GitHub Environment production" : ""),
  };
}

export async function cloudflareConfigured(env) {
  return Boolean(env.CF_GLOBAL_API_KEY || env.CF_API_TOKEN);
}
