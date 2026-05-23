/**
 * Platform setup: GitHub Secrets → Docker + Cloudflare Workers Secrets.
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
          "GitHub Secrets diperbarui. Deploy Mini PC dipicu: " + r.secrets_updated.join(", "),
          true
        );
        return r;
      });
    },
    saveCloudflare: function (form, msgEl, testOnly) {
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
        showMsg(msgEl, testOnly ? "Koneksi OK" : r.message || "Tersimpan di Workers Secrets", true);
        return r;
      });
    },
    loadStatus: function (el) {
      return api("/admin/api/platform/setup/status").then(function (s) {
        if (!el) return s;
        el.innerHTML =
          "<strong>Status platform</strong><br/>" +
          "GitHub bootstrap: " +
          (s.github_bootstrap ? "OK" : "GITHUB_SETUP_TOKEN belum di Worker") +
          "<br/>Cloudflare: " +
          (s.cloudflare_configured ? "terkonfigurasi" : "belum");
        return s;
      });
    },
  };
})();
