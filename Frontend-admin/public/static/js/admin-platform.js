/**
 * Platform setup: GitHub Environment production + Workers Secrets + runner check.
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

  function esc(s) {
    return String(s)
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;");
  }

  function renderRunnerBanner(container, status) {
    if (!container) return;
    var r = status.runner || {};
    if (!r.needs_runner_setup) {
      container.hidden = true;
      container.innerHTML = "";
      return;
    }
    var url = (status.runner_install && status.runner_install.script_url) || "";
    var runnersNew =
      (status.runner_install && status.runner_install.runners_new_url) ||
      "https://github.com/WBetEngine/Seosementara/settings/actions/runners/new";
    var ps =
      "mkdir C:\\Seosementara\\scripts -Force\n" +
      'Invoke-WebRequest -Uri "' +
      url +
      '" -OutFile C:\\Seosementara\\scripts\\install-github-runner.ps1\n' +
      "C:\\Seosementara\\scripts\\install-github-runner.ps1";
    var title = "Self-hosted runner belum siap";
    if (r.reason === "runners_offline") title = "Runner terdaftar tapi offline";
    if (r.reason === "github_pat_missing")
      title = "Isi Bootstrap Platform dulu — lalu cek runner otomatis";

    container.hidden = false;
    container.className = "alert alert-error";
    container.innerHTML =
      "<strong>" +
      esc(title) +
      "</strong>" +
      '<p style="margin:0.5rem 0 0.75rem;font-size:0.9rem">' +
      "Mini PC hanya Docker. Pasang runner sekali di mini PC — PowerShell <strong>Administrator</strong>, tanpa clone repo." +
      "</p>" +
      '<p style="margin:0 0 0.35rem;font-size:0.85rem">Token registrasi: ' +
      '<a href="' +
      esc(runnersNew) +
      '" target="_blank" rel="noopener">GitHub → New self-hosted runner</a></p>' +
      '<pre style="margin:0;padding:0.75rem;background:#1e1e1e;color:#e5e5e5;border-radius:6px;font-size:0.8rem;overflow:auto;white-space:pre-wrap;user-select:all">' +
      esc(ps) +
      "</pre>" +
      '<p style="margin:0.5rem 0 0;font-size:0.8rem;color:var(--text-muted)">' +
      "Setelah selesai, runner harus <strong>Idle</strong> di GitHub. Muat ulang halaman ini." +
      "</p>";
  }

  function renderStatusSummary(el, s) {
    if (!el) return;
    var r = s.runner || {};
    var runnerLine = r.checked
      ? "Runner online: " + r.online_count + " / " + r.total_count
      : "Runner: belum dicek (isi Bootstrap dulu)";
    if (r.checked && r.runners && r.runners.length) {
      runnerLine +=
        " (" +
        r.runners
          .map(function (x) {
            return x.name + "=" + x.status;
          })
          .join(", ") +
        ")";
    }
    el.innerHTML =
      "<strong>Status platform</strong><br/>" +
      "GitHub PAT: " +
      (s.github_token_stored ? "tersimpan" : "belum") +
      "<br/>Bootstrap: " +
      (s.bootstrap_complete ? "OK" : "belum") +
      "<br/>Infra secrets: " +
      (s.infra_complete ? "OK" : "belum") +
      "<br/>" +
      runnerLine +
      "<br/>Environment: <code>" +
      esc(s.github_environment || "production") +
      "</code>";
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
          (r.message || "Bootstrap OK") + " — " + (r.secrets_updated || []).join(", "),
          true
        );
        return SeosementaraPlatform.refreshSetupUI();
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
        return SeosementaraPlatform.refreshSetupUI();
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
        var extra =
          r.github_environment_synced && r.github_environment_synced.length
            ? " + GitHub: " + r.github_environment_synced.join(", ")
            : "";
        showMsg(msgEl, (r.message || "OK") + extra, true);
        return r;
      });
    },
    refreshSetupUI: function () {
      return api("/admin/api/platform/setup/status").then(function (s) {
        renderStatusSummary(document.getElementById("platform-status"), s);
        renderRunnerBanner(document.getElementById("runner-setup-banner"), s);
        return s;
      });
    },
    loadStatus: function () {
      return SeosementaraPlatform.refreshSetupUI();
    },
  };
})();
