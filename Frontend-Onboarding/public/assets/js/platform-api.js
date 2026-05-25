(function () {
  'use strict';

  var LS_KEY = 'sse_platform_api_url';

  function normalizeUrl(url) {
    return (url || '').trim().replace(/\/$/, '');
  }

  function readStored() {
    try {
      return normalizeUrl(localStorage.getItem(LS_KEY) || '');
    } catch (e) {
      return '';
    }
  }

  function persistBase(url) {
    var u = normalizeUrl(url);
    window.__SSE_PLATFORM_API__ = u;
    if (window.SSEO) window.SSEO.platformApiBase = u;
    try {
      if (u) localStorage.setItem(LS_KEY, u);
      else localStorage.removeItem(LS_KEY);
    } catch (e) {}
    if (typeof window.refreshPlatformApiUi === 'function') {
      window.refreshPlatformApiUi();
    }
    return u;
  }

  function apiBase() {
    if (typeof window.__SSE_PLATFORM_API__ === 'string' && window.__SSE_PLATFORM_API__) {
      return normalizeUrl(window.__SSE_PLATFORM_API__);
    }
    if (window.SSEO && window.SSEO.platformApiBase) {
      return normalizeUrl(window.SSEO.platformApiBase);
    }
    var params = new URLSearchParams(location.search);
    var q = params.get('api');
    if (q) return normalizeUrl(q);
    return readStored();
  }

  (function initApiUrl() {
    var params = new URLSearchParams(location.search);
    var q = params.get('api');
    if (q) {
      persistBase(q);
      return;
    }
    if (!window.__SSE_PLATFORM_API__) {
      var stored = readStored();
      if (stored) window.__SSE_PLATFORM_API__ = stored;
    }
  })();

  function sessionId() {
    try {
      return sessionStorage.getItem('sse_setup_session') || '';
    } catch (e) {
      return '';
    }
  }

  function setSession(id) {
    try {
      sessionStorage.setItem('sse_setup_session', id);
    } catch (e) {}
  }

  async function request(method, path, body) {
    var base = apiBase();
    if (!base) {
      throw new Error(
        'Platform API belum dikonfigurasi. Isi URL Worker di kotak kuning di atas (Simpan & tes koneksi), jalankan workflow Deploy Platform Worker di GitHub Actions, atau buka dengan ?api=https://sse-platform.<account>.workers.dev'
      );
    }

    var headers = { 'Content-Type': 'application/json' };
    var sid = sessionId();
    if (sid) headers['X-Setup-Session'] = sid;

    var res = await fetch(base + '/admin/api/platform' + path, {
      method: method,
      headers: headers,
      body: body ? JSON.stringify(body) : undefined
    });

    var data = null;
    try {
      data = await res.json();
    } catch (e) {
      data = { ok: false, error: await res.text() };
    }

    if (!res.ok || data.ok === false) {
      var msg = data.error || data.message || 'HTTP ' + res.status;
      throw new Error(msg);
    }
    return data;
  }

  window.SSEOPlatform = {
    apiBase: apiBase,
    setApiBase: persistBase,
    testConnection: function () {
      return request('GET', '/setup/status');
    },
    getStatus: function () {
      return request('GET', '/setup/status');
    },
    saveGithubPat: function (github_pat) {
      return request('POST', '/github/pat', { github_pat: github_pat }).then(function (d) {
        if (d.session_id) setSession(d.session_id);
        return d;
      });
    },
    testCloudflare: function (payload) {
      return request('POST', '/cloudflare/credentials/test', payload);
    },
    saveCloudflare: function (payload) {
      return request('POST', '/cloudflare/credentials', payload);
    },
    testSsh: function (payload) {
      return request('POST', '/infra/ssh/test', payload);
    },
    registerRunner: function (payload) {
      return request('POST', '/github/runner/register', payload);
    },
    createTunnel: function (payload) {
      return request('POST', '/cloudflare/tunnel/create', payload);
    },
    saveDatabase: function (payload) {
      return request('POST', '/infra/database', payload);
    },
    deployBackend: function () {
      return request('POST', '/deploy/backend', {});
    },
    deployAdminPages: function () {
      return request('POST', '/deploy/admin-pages', {});
    },
    deployPublicPages: function () {
      return request('POST', '/deploy/public-pages', {});
    },
    formPayload: function (form) {
      var data = {};
      new FormData(form).forEach(function (v, k) {
        if (v !== '') data[k] = v;
      });
      return data;
    }
  };
})();
