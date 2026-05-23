/**
 * Platform setup: GitHub Environment production + Workers Secrets.
 */
(function () {
  "use strict";

  function api(path, options) {
    return fetch(path, {
      headers: { "Content-Type": "application/json", Accept: "application/json" },
      ...options,
    }).then(function (res) {
      return res.json().then(function (data) {
        if (!res.ok) throw new Error(data.error || res.statusText);
        return data;
      });
    });
  }

  function showMsg(el, text, ok) {
    if (!el) return;
    el.textContent = text;
    el.hidden = false;
    el.className = ok ? "alert alert-success" : "alert alert-error";
  }

  window.SeosementaraPlatform = {
    saveBootstrap: function (form, msgEl) {
      var fd = new FormData(form);
      return api("/admin/api/platform/bootstrap", {
        method: "POST",
        body: JSON.stringify({
          github_pat: fd.get("github_pat"),
          global_api_key: fd.get("global_api_key"),
          account_email: fd.get("account_email"),
          account_id: fd.get("account_id"),
          super_admin_token: fd.get("super_admin_token") || undefined,
        }),
      }).then(function (r) {
        showMsg(
          msgEl,
          (r.message || "Bootstrap OK") +
            " — Environment: " +
            (r.secrets_updated || []).join(", "),
          true
        );
        return r;
      });
    },
    saveInfra: function (form, msgEl) {
      var fd = new FormData(form);
      return api("/admin/api/platform/infra", {
        method: "POST",
        body: JSON.stringify({
          db_password: fd.get("db_password"),
          master_encryption_key: fd.get("master_encryption_key"),
          super_admin_token: fd.get("super_admin_token") || undefined,
        }),
      }).then(function (r) {
        showMsg(
          msgEl,
          "Environment production: " + r.secrets_updated.join(", ") + " — deploy dipicu",
          true
        );
        return r;
      });
    },
    saveCloudflare: function (form, msgEl) {
      var fd = new FormData(form);
      var authType = fd.get("auth_type") || "global_api_key";
      var body = { auth_type: authType, account_id: fd.get("account_id") || undefined };
      if (authType === "global_api_key") {
        body.global_api_key = fd.get("global_api_key");
        body.account_email = fd.get("account_email");
      } else {
        body.api_token = fd.get("api_token");
        body.account_email = fd.get("account_email") || undefined;
      }
      return api("/admin/api/platform/cloudflare/credentials", {
        method: "POST",
        body: JSON.stringify(body),
      }).then(function (r) {
        var extra = r.github_environment_synced && r.github_environment_synced.length
          ? " + GitHub: " + r.github_environment_synced.join(", ")
          : "";
        showMsg(msgEl, (r.message || "OK") + extra, true);
        return r;
      });
    },
    loadStatus: function (el) {
      return api("/admin/api/platform/setup/status").then(function (s) {
        if (!el) return s;
        el.innerHTML =
          "<strong>Status</strong><br/>" +
          "GitHub PAT di Worker: " +
          (s.github_token_stored ? "tersimpan" : "belum — isi Bootstrap") +
          "<br/>Environment: <code>" +
          (s.github_environment || "production") +
          "</code><br/>Cloudflare Worker: " +
          (s.cloudflare_configured ? "OK" : "belum");
        return s;
      });
    },
  };
})();
