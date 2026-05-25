const CF_API = 'https://api.cloudflare.com/client/v4';

function cfHeaders(token) {
  return {
    Authorization: `Bearer ${token}`,
    'Content-Type': 'application/json'
  };
}

export async function verifyToken(token) {
  const res = await fetch(`${CF_API}/user/tokens/verify`, {
    method: 'GET',
    headers: cfHeaders(token)
  });
  const data = await res.json();
  if (!res.ok || !data.success) {
    return {
      ok: false,
      error: (data.errors && data.errors[0]?.message) || `Token tidak valid (${res.status})`
    };
  }
  return { ok: true, status: data.result?.status };
}

export async function getAccount(token, accountId) {
  const res = await fetch(`${CF_API}/accounts/${accountId}`, {
    headers: cfHeaders(token)
  });
  const data = await res.json();
  if (!res.ok || !data.success) {
    return { ok: false, error: 'Account ID tidak ditemukan atau token tidak punya akses' };
  }
  return { ok: true, name: data.result?.name };
}

export async function getZone(token, zoneId) {
  const res = await fetch(`${CF_API}/zones/${zoneId}`, {
    headers: cfHeaders(token)
  });
  const data = await res.json();
  if (!res.ok || !data.success) {
    return { ok: false, error: 'Zone ID tidak valid' };
  }
  return {
    ok: true,
    name: data.result?.name,
    status: data.result?.status
  };
}

export async function createTunnel(token, accountId, name) {
  const res = await fetch(`${CF_API}/accounts/${accountId}/cfd_tunnel`, {
    method: 'POST',
    headers: cfHeaders(token),
    body: JSON.stringify({ name, config_src: 'cloudflare' })
  });
  const data = await res.json();
  if (!res.ok || !data.success) {
    const msg = (data.errors && data.errors[0]?.message) || res.statusText;
    return { ok: false, error: msg };
  }
  return {
    ok: true,
    tunnel_id: data.result?.id,
    tunnel_name: data.result?.name
  };
}
