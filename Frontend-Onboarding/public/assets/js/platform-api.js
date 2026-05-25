(function () {
  'use strict';

  function apiBase() {
    if (typeof window.__SSE_PLATFORM_API__ === 'string' && window.__SSE_PLATFORM_API__) {
      return window.__SSE_PLATFORM_API__.replace(/\/$/, '');
    }
    if (window.SSEO && window.SSEO.platformApiBase) {
      return window.SSEO.platformApiBase.replace(/\/$/, '');
    }
    var params = new URLSearchParams(location.search);
    var q = params.get('api');
    if (q) return q.replace(/\/$/, '');
    return '';
  }

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
        'Platform API belum dikonfigurasi. Deploy platform-worker dulu, atau buka dengan ?api=https://sse-platform.<account>.workers.dev'
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
