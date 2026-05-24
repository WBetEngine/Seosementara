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
    var title = "Runner belum siap";
    var hint =
      "Jalankan sekali di mini PC (PowerShell): cd C:\\actions-runner; .\\run.cmd — biarkan terbuka, lalu Simpan Infra lagi (otomatis pasang service + deploy).";
    if (r.reason === "runners_offline") {
      title = "Runner offline";
      hint =
        "Buka C:\\actions-runner\\run.cmd atau pasang Windows Service lewat tombol di bawah / Simpan Infra.";
    }
    if (r.reason === "github_pat_missing") {
      title = "Isi Bootstrap Platform dulu";
      hint = "PAT harus punya Administration write untuk registration token otomatis.";
    }

    container.hidden = false;
    container.className = "alert alert-error";
    container.innerHTML =
      "<strong>" +
      esc(title) +
      "</strong>" +
      '<p style="margin:0.5rem 0 0.75rem;font-size:0.9rem">' +
      esc(hint) +
      "</p>" +
      '<button type="button" class="btn btn-secondary" id="btn-install-runner-service" style="margin-top:0.25rem">' +
      "Pasang runner service otomatis (GitHub Actions)" +
      "</button>" +
      '<p style="margin:0.5rem 0 0;font-size:0.8rem;color:var(--text-muted)">' +
      "Tanpa SSH — workflow jalan di mini PC, token registrasi dari PAT Bootstrap." +
      "</p>";

    var btn = document.getElementById("btn-install-runner-service");
    if (btn && !btn._bound) {
      btn._bound = true;
      btn.addEventListener("click", function () {
        btn.disabled = true;
        api("/admin/api/platform/runner/install-service", { method: "POST" })
          .then(function (res) {
            alert(res.message || "Workflow dipicu");
            return SeosementaraPlatform.refreshSetupUI();
          })
          .catch(function (e) {
            alert(e.message);
          })
          .finally(function () {
            btn.disabled = false;
          });
      });
    }
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
      "<br/>Simpan Infra: otomatis runner service + deploy Docker (via GitHub Actions)" +
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
        var extra = r.runner_service_triggered ? " + Install Runner Service" : "";
        showMsg(
          msgEl,
          (r.message || "Bootstrap OK") + extra + " — " + (r.secrets_updated || []).join(", "),
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
        var wf = (r.workflows_triggered || []).join(", ") || "deploy-mini-pc.yml";
        showMsg(msgEl, (r.message || "OK") + " Workflows: " + wf, true);
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
