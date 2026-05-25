import sodium from 'libsodium-wrappers';

const GH_API = 'https://api.github.com';

export async function validatePat(pat) {
  const res = await fetch(`${GH_API}/user`, {
    headers: githubHeaders(pat)
  });
  if (!res.ok) {
    const err = await res.text();
    return { ok: false, error: `GitHub PAT tidak valid (${res.status}): ${err.slice(0, 200)}` };
  }
  const user = await res.json();
  return { ok: true, login: user.login };
}

function githubHeaders(pat) {
  return {
    Authorization: `Bearer ${pat}`,
    Accept: 'application/vnd.github+json',
    'X-GitHub-Api-Version': '2022-11-28',
    'User-Agent': 'sse-platform-worker'
  };
}

export async function setRepoSecret(repo, pat, secretName, secretValue) {
  await sodium.ready;
  const [owner, name] = repo.split('/');
  const keyRes = await fetch(
    `${GH_API}/repos/${owner}/${name}/actions/secrets/public-key`,
    { headers: githubHeaders(pat) }
  );
  if (!keyRes.ok) {
    throw new Error(`Gagal ambil public key GitHub (${keyRes.status})`);
  }
  const keyData = await keyRes.json();
  const messageBytes = sodium.from_string(secretValue);
  const keyBytes = sodium.from_base64(keyData.key, sodium.base64_variants.ORIGINAL);
  const encryptedBytes = sodium.crypto_box_seal(messageBytes, keyBytes);
  const encrypted = sodium.to_base64(encryptedBytes, sodium.base64_variants.ORIGINAL);

  const putRes = await fetch(
    `${GH_API}/repos/${owner}/${name}/actions/secrets/${secretName}`,
    {
      method: 'PUT',
      headers: { ...githubHeaders(pat), 'Content-Type': 'application/json' },
      body: JSON.stringify({
        encrypted_value: encrypted,
        key_id: keyData.key_id
      })
    }
  );
  if (!putRes.ok) {
    const err = await putRes.text();
    throw new Error(`Gagal simpan secret ${secretName} (${putRes.status}): ${err.slice(0, 200)}`);
  }
  return { ok: true };
}

export async function dispatchWorkflow(repo, pat, eventType, clientPayload) {
  const [owner, name] = repo.split('/');
  const res = await fetch(`${GH_API}/repos/${owner}/${name}/dispatches`, {
    method: 'POST',
    headers: { ...githubHeaders(pat), 'Content-Type': 'application/json' },
    body: JSON.stringify({
      event_type: eventType,
      client_payload: clientPayload || {}
    })
  });
  if (!res.ok) {
    const err = await res.text();
    throw new Error(`Gagal trigger workflow (${res.status}): ${err.slice(0, 300)}`);
  }
  return { ok: true };
}
