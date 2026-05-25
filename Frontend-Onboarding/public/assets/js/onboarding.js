(function () {
  'use strict';

  var MAX_STEP = 9;
  var CF_BOOTSTRAP_KEY = 'sse_cf_bootstrap';
  var API = function () {
    return window.SSEOPlatform;
  };

  var STEP_FIELDS = {
    1: ['cf_token', 'cf_account'],
    2: ['github_pat'],
    3: ['cf_zone', 'primary_domain'],
    4: ['ssh_host', 'ssh_port', 'ssh_user', 'ssh_secret'],
    5: ['runner_label'],
    6: ['tunnel_name'],
    7: ['db_password', 'master_key'],
    8: [],
    9: []
  };

  function stepNeedsApi(stepNum) {
    return stepNum >= 3;
  }

  function saveCfBootstrapLocal(form) {
    var payload = API().formPayload(form);
    try {
      localStorage.setItem(
        CF_BOOTSTRAP_KEY,
        JSON.stringify({
          cf_token: payload.cf_token || '',
          cf_account: payload.cf_account || ''
        })
      );
    } catch (e) {}
  }

  function loadCfBootstrapIntoForm(form) {
    try {
      var raw = localStorage.getItem(CF_BOOTSTRAP_KEY);
      if (!raw) return;
      var data = JSON.parse(raw);
      if (data.cf_token && getInput('cf_token')) getInput('cf_token').value = data.cf_token;
      if (data.cf_account && getInput('cf_account')) getInput('cf_account').value = data.cf_account;
    } catch (e) {}
  }

  var VALIDATORS = {
    github_pat: function (v) {
      if (!v || !v.trim()) return 'GitHub PAT wajib diisi.';
      if (!/^(ghp_|github_pat_|gho_)/i.test(v.trim())) {
        return 'Format token harus diawali ghp_, github_pat_, atau gho_.';
      }
      if (v.trim().length < 20) return 'Token terlalu pendek.';
      return '';
    },
    cf_token: function (v) {
      if (!v || !v.trim()) return 'API Token wajib diisi.';
      if (v.trim().length < 20) return 'Token minimal 20 karakter.';
      return '';
    },
    cf_account: function (v) {
      if (!v || !v.trim()) return 'Account ID wajib diisi.';
      if (!/^[a-f0-9]{32}$/i.test(v.trim())) return 'Account ID harus 32 karakter hex.';
      return '';
    },
    cf_zone: function (v) {
      if (!v || !v.trim()) return 'Zone ID wajib diisi.';
      if (!/^[a-f0-9]{32}$/i.test(v.trim())) return 'Zone ID harus 32 karakter hex.';
      return '';
    },
    primary_domain: function (v) {
      if (!v || !v.trim()) return 'Domain utama wajib diisi.';
      if (!/^[a-z0-9]([a-z0-9-]*[a-z0-9])?(\.[a-z0-9]([a-z0-9-]*[a-z0-9])?)+$/i.test(v.trim())) {
        return 'Format domain tidak valid.';
      }
      return '';
    },
    ssh_host: function (v) {
      if (!v || !v.trim()) return 'Host / IP wajib diisi.';
      var t = v.trim();
      if (
        !/^(\d{1,3}\.){3}\d{1,3}$/.test(t) &&
        !/^[a-z0-9]([a-z0-9.-]*[a-z0-9])?$/i.test(t)
      ) {
        return 'Masukkan IP atau hostname yang valid.';
      }
      return '';
    },
    ssh_port: function (v) {
      var n = parseInt(v, 10);
      if (!v || isNaN(n)) return 'Port wajib diisi.';
      if (n < 1 || n > 65535) return 'Port harus 1–65535.';
      return '';
    },
    ssh_user: function (v) {
      if (!v || !v.trim()) return 'User SSH wajib diisi.';
      return '';
    },
    ssh_secret: function (v) {
      if (!v || !v.trim()) return 'Password atau private key wajib diisi.';
      if (v.length < 4) return 'Kredensial terlalu pendek.';
      return '';
    },
    runner_label: function (v) {
      if (!v || !v.trim()) return 'Label runner wajib diisi.';
      if (!/^[a-z0-9][a-z0-9-_]*$/i.test(v.trim())) {
        return 'Hanya huruf, angka, strip, dan underscore.';
      }
      return '';
    },
    tunnel_name: function (v) {
      if (!v || !v.trim()) return 'Nama tunnel wajib diisi.';
      if (!/^[a-z0-9][a-z0-9-_]*$/i.test(v.trim())) {
        return 'Nama tunnel: huruf kecil, angka, strip, underscore.';
      }
      return '';
    },
    db_password: function (v) {
      if (!v || v.length < 12) return 'Password DB minimal 12 karakter.';
      return '';
    },
    master_key: function (v) {
      if (!v || v.length < 16) return 'Master encryption key minimal 16 karakter.';
      return '';
    }
  };

  function qs(sel, root) {
    return (root || document).querySelector(sel);
  }

  function qsa(sel, root) {
    return Array.from((root || document).querySelectorAll(sel));
  }

  function showToast(message, type) {
    type = type || 'info';
    var toast = qs('#toast');
    if (!toast) return;
    toast.innerHTML =
      '<div class="toast toast--' + type + '" role="status">' + message + '</div>';
    toast.classList.add('is-visible');
    setTimeout(function () {
      toast.classList.remove('is-visible');
    }, 5000);
  }

  function getInput(name) {
    return qs('#wizard-form [name="' + name + '"]');
  }

  function getErrorEl(name) {
    return qs('[data-error-for="' + name + '"]');
  }

  function validateField(name, showMsg) {
    var input = getInput(name);
    if (!input || !VALIDATORS[name]) return true;
    var msg = VALIDATORS[name](input.value);
    var err = getErrorEl(name);
    if (showMsg !== false && err) err.textContent = msg;
    input.classList.toggle('is-invalid', !!msg);
    input.classList.toggle('is-valid', !msg && input.value.trim() !== '');
    input.setAttribute('aria-invalid', msg ? 'true' : 'false');
    return !msg;
  }

  function validateStep(stepNum, showMsg) {
    var fields = STEP_FIELDS[stepNum] || [];
    var ok = true;
    fields.forEach(function (name) {
      if (!validateField(name, showMsg)) ok = false;
    });
    return ok;
  }

  function setButtonLoading(btn, loading, label) {
    if (!btn) return;
    btn.disabled = loading;
    if (loading) {
      btn.dataset.prevText = btn.textContent;
      btn.textContent = label || 'Memproses…';
    } else if (btn.dataset.prevText) {
      btn.textContent = btn.dataset.prevText;
      delete btn.dataset.prevText;
    }
  }

  function bindRealtimeValidation() {
    Object.keys(VALIDATORS).forEach(function (name) {
      var input = getInput(name);
      if (!input) return;
      input.addEventListener('input', function () {
        validateField(name, true);
        if (typeof updateStepStatus === 'function') updateStepStatus();
      });
      input.addEventListener('blur', function () {
        validateField(name, true);
        if (typeof updateStepStatus === 'function') updateStepStatus();
      });
    });
  }

  function initInfoIcons() {
    qsa('.info-icon').forEach(function (btn) {
      var tipText = btn.getAttribute('data-tip') || '';
      var group = btn.closest('.form-group') || btn.closest('.api-config');
      var tipEl = group ? qs('.info-tip', group) : null;

      btn.addEventListener('click', function (e) {
        e.preventDefault();
        qsa('.info-tip.is-open').forEach(function (t) {
          if (t !== tipEl) {
            t.classList.remove('is-open');
            t.textContent = '';
          }
        });
        qsa('.info-icon[aria-expanded="true"]').forEach(function (b) {
          if (b !== btn) b.setAttribute('aria-expanded', 'false');
        });
        if (tipEl) {
          var open = !tipEl.classList.contains('is-open');
          tipEl.classList.toggle('is-open', open);
          tipEl.textContent = open ? tipText : '';
          btn.setAttribute('aria-expanded', open ? 'true' : 'false');
        }
      });
    });
  }

  function initPasswordToggles() {
    qsa('.pw-toggle').forEach(function (btn) {
      var wrap = btn.closest('.form-control-wrap');
      var input = wrap ? qs('.form-control', wrap) : null;
      if (!input) {
        var id = btn.getAttribute('data-pw-toggle');
        input = id ? qs('#' + id) : null;
      }
      if (!input) return;

      btn.addEventListener('click', function (e) {
        e.preventDefault();
        var show = input.type === 'password';
        input.type = show ? 'text' : 'password';
        btn.classList.toggle('is-visible', show);
        btn.setAttribute('aria-pressed', show ? 'true' : 'false');
        btn.setAttribute('aria-label', show ? 'Sembunyikan isi' : 'Tampilkan isi');
      });
    });
  }

  function initExternalLinks() {
    var L = window.SSEO && window.SSEO.links;
    if (!L) return;
    qsa('[data-ext]').forEach(function (a) {
      var key = a.getAttribute('data-ext');
      if (L[key]) a.href = L[key];
    });
  }

  function refreshApiBanner() {
    var banner = qs('#api-banner');
    var connected = qs('#api-banner-connected');
    var setup = qs('#api-banner-setup');
    var urlEl = qs('#api-banner-url');
    var input = qs('#platform-api-input');
    if (!banner) return;

    var base = API() && API().apiBase();
    if (base) {
      banner.className = 'alert alert--info';
      if (connected) connected.hidden = false;
      if (setup) setup.hidden = true;
      if (urlEl) urlEl.textContent = base;
      if (input && !input.value) input.value = base;
    } else {
      banner.className = 'alert alert--warning';
      if (connected) connected.hidden = true;
      if (setup) setup.hidden = false;
      try {
        var saved = localStorage.getItem('sse_platform_api_url');
        if (input && saved && !input.value) input.value = saved;
      } catch (e) {}
    }
    if (typeof updateStepStatus === 'function') updateStepStatus();
  }

  function initApiBanner() {
    var deployLink = qs('#link-deploy-worker');
    var L = window.SSEO && window.SSEO.links;
    if (deployLink && L && L.platformWorkerDeploy) {
      deployLink.href = L.platformWorkerDeploy;
    }

    var saveBtn = qs('#btn-api-save-test');
    var changeBtn = qs('#btn-api-change');
    var input = qs('#platform-api-input');

    if (saveBtn && input) {
      saveBtn.addEventListener('click', async function () {
        var url = input.value.trim();
        if (!url) {
          showToast('Masukkan URL Platform Worker.', 'error');
          input.focus();
          return;
        }
        if (!/^https?:\/\//i.test(url)) {
          showToast('URL harus diawali http:// atau https://', 'error');
          input.focus();
          return;
        }
        setButtonLoading(saveBtn, true, 'Menguji…');
        try {
          API().setApiBase(url);
          await API().testConnection();
          showToast('Platform API terhubung.', 'success');
          refreshApiBanner();
        } catch (err) {
          API().setApiBase('');
          showToast(err.message || String(err), 'error');
        } finally {
          setButtonLoading(saveBtn, false);
        }
      });
    }

    if (changeBtn) {
      changeBtn.addEventListener('click', function () {
        var connected = qs('#api-banner-connected');
        var setup = qs('#api-banner-setup');
        var banner = qs('#api-banner');
        if (connected) connected.hidden = true;
        if (setup) setup.hidden = false;
        if (banner) banner.className = 'alert alert--warning';
        if (input) {
          input.value = (API() && API().apiBase()) || '';
          input.focus();
        }
        API().setApiBase('');
      });
    }

    window.refreshPlatformApiUi = refreshApiBanner;
    refreshApiBanner();
  }

  async function checkPlatformStatus() {
    try {
      var res = await API().getStatus();
      if (res.status && res.status.bootstrap_complete) {
        showToast('Bootstrap sudah selesai menurut server.', 'success');
      }
    } catch (e) {
      /* API belum deploy */
    }
  }

  var updateStepStatus;

  function initWizard() {
    var form = qs('#wizard-form');
    if (!form) return;

    var step = 1;
    var prev = qs('#btn-prev');
    var next = qs('#btn-next');
    var statusEl = qs('#step-status');

    updateStepStatus = function () {
      if (!statusEl) return;
      if (step >= MAX_STEP) {
        statusEl.textContent = '';
        statusEl.className = 'step-status';
        return;
      }
      var ok = validateStep(step, false);
      var needsApi = stepNeedsApi(step);
      var apiOk = !needsApi || (API() && API().apiBase());
      if (step === 1) {
        statusEl.textContent = ok
          ? 'Langkah 1 siap — Lanjut ke GitHub PAT'
          : 'Isi CLOUDFLARE_API_TOKEN dan CLOUDFLARE_ACCOUNT_ID';
      } else if (step === 2) {
        statusEl.textContent = ok
          ? 'Deploy worker via tombol utama, lalu hubungkan URL API di banner'
          : 'Isi GitHub PAT';
      } else {
        statusEl.textContent = !apiOk
          ? 'Hubungkan Platform API (banner) setelah deploy worker'
          : ok
            ? 'Langkah siap — gunakan Test atau Lanjut'
            : 'Lengkapi field yang ditandai merah';
      }
      statusEl.className = 'step-status' + (ok && apiOk ? ' is-ok' : '');
      if (next) next.disabled = !ok || !apiOk;
    };

    function show(s) {
      step = s;
      qsa('.wizard-panel', form).forEach(function (p) {
        p.hidden = p.getAttribute('data-panel') !== String(s);
      });
      qsa('.wizard-step').forEach(function (el) {
        var n = +el.getAttribute('data-step');
        el.classList.toggle('is-active', n === s);
        el.classList.toggle('is-done', n < s);
      });
      if (prev) prev.disabled = s <= 1;
      if (next) {
        next.textContent = 'Lanjut';
        next.style.display = s >= MAX_STEP ? 'none' : '';
      }
      var nav = qs('#wizard-nav');
      if (nav) nav.style.display = s >= MAX_STEP ? 'none' : '';
      if (s < MAX_STEP) updateStepStatus();
    }

    if (prev) {
      prev.addEventListener('click', function () {
        if (step > 1) show(step - 1);
      });
    }

    if (next) {
      next.addEventListener('click', function () {
        if (!validateStep(step, true)) {
          showToast('Perbaiki field yang belum valid.', 'error');
          var first = qs('.form-control.is-invalid', form);
          if (first) first.focus();
          return;
        }
        if (step < MAX_STEP) show(step + 1);
      });
    }

    qsa('[data-action]').forEach(function (btn) {
      btn.addEventListener('click', async function () {
        if (!validateStep(step, true)) {
          showToast('Lengkapi form yang valid dulu.', 'error');
          return;
        }
        var action = btn.getAttribute('data-action');
        var payload = API().formPayload(form);
        setButtonLoading(btn, true);
        try {
          var result;
          switch (action) {
            case 'cf-bootstrap-verify':
              saveCfBootstrapLocal(form);
              if (!API().apiBase()) {
                showToast(
                  'Worker belum online — simpan kredensial, deploy via langkah 2 atau workflow manual, lalu verifikasi lagi.',
                  'warning'
                );
                return;
              }
              result = await API().bootstrapVerifyCf({
                cf_token: payload.cf_token,
                cf_account: payload.cf_account
              });
              showToast(
                (result.message || 'Cloudflare OK') +
                  (result.account_name ? ' — ' + result.account_name : ''),
                'success'
              );
              break;
            case 'cf-bootstrap-local':
              saveCfBootstrapLocal(form);
              showToast('Kredensial Cloudflare disimpan sementara di browser.', 'success');
              break;
            case 'github-initial-setup':
              saveCfBootstrapLocal(form);
              if (!payload.cf_token || !payload.cf_account) {
                showToast('Lengkapi langkah 1 (Cloudflare token + account) dulu.', 'error');
                return;
              }
              if (!API().apiBase()) {
                var deployUrl =
                  window.SSEO && window.SSEO.links && window.SSEO.links.platformWorkerDeploy;
                showToast(
                  'Deploy pertama: jalankan workflow Deploy Platform Worker (isi token + account dari langkah 1). Setelah hijau, tempel URL worker di banner lalu ulangi tombol ini.',
                  'warning'
                );
                if (deployUrl) window.open(deployUrl, '_blank', 'noopener,noreferrer');
                return;
              }
              result = await API().initialSetup({
                github_pat: payload.github_pat,
                cf_token: payload.cf_token,
                cf_account: payload.cf_account
              });
              if (result.worker_url) API().setApiBase(result.worker_url);
              showToast(
                (result.message || 'Setup OK') + (result.login ? ' — ' + result.login : ''),
                'success'
              );
              break;
            case 'github-pat':
              if (!API().apiBase()) {
                showToast('Hubungkan Platform API dulu (banner atas).', 'error');
                return;
              }
              result = await API().saveGithubPat(payload.github_pat);
              showToast('GitHub PAT valid — login: ' + (result.login || ''), 'success');
              break;
            case 'cf-test':
              result = await API().testCloudflare(payload);
              showToast(
                'Cloudflare OK' +
                  (result.zone_name ? ' — zone: ' + result.zone_name : ''),
                'success'
              );
              break;
            case 'cf-save':
              result = await API().saveCloudflare(payload);
              showToast(result.message || 'Cloudflare tersimpan', 'success');
              break;
            case 'ssh-test':
              result = await API().testSsh(payload);
              showToast(result.message || 'Test SSH dipicu', 'success');
              break;
            case 'runner':
              result = await API().registerRunner(payload);
              showToast(result.message || 'Register runner dipicu', 'success');
              break;
            case 'tunnel':
              result = await API().createTunnel(payload);
              showToast(
                (result.message || 'Tunnel dibuat') +
                  (result.tunnel_id ? ' ID: ' + result.tunnel_id : ''),
                'success'
              );
              break;
            case 'database':
              result = await API().saveDatabase(payload);
              showToast(result.message || 'Secrets database tersimpan', 'success');
              break;
            case 'deploy-all':
              await API().deployBackend();
              showToast('Deploy backend dipicu', 'success');
              await API().deployAdminPages();
              showToast('Deploy admin Pages dipicu', 'success');
              await API().deployPublicPages();
              showToast('Deploy publik Pages dipicu — cek GitHub Actions', 'success');
              break;
            default:
              throw new Error('Aksi tidak dikenal');
          }
        } catch (err) {
          showToast(err.message || String(err), 'error');
        } finally {
          setButtonLoading(btn, false);
        }
      });
    });

    var finish = qs('#btn-finish');
    if (finish) {
      finish.addEventListener('click', async function (e) {
        e.preventDefault();
        setButtonLoading(finish, true, 'Membuka admin…');
        try {
          var st = await API().getStatus();
          if (!st.status || !st.status.bootstrap_complete) {
            showToast(
              'Deploy belum lengkap — jalankan langkah 8 (Deploy) dan cek Actions hijau.',
              'warning'
            );
          }
          try {
            localStorage.setItem(
              window.SSEO.wizardStorageKey,
              JSON.stringify(API().formPayload(form))
            );
          } catch (err) {}
          var url =
            (window.SSEO && window.SSEO.adminUrlAfterComplete) ||
            'https://seosementara.org/admin/login.html?from=onboarding';
          window.location.href = url;
        } catch (err) {
          showToast(err.message, 'error');
          setButtonLoading(finish, false);
        }
      });
    }

    show(1);
  }

  document.addEventListener('DOMContentLoaded', function () {
    var css = qs('#onboarding-css');
    if (css && window.SSEO && window.SSEO.asset) {
      css.href = window.SSEO.asset('/assets/css/onboarding.css') + '?v=20250524';
    }
    initExternalLinks();
    initInfoIcons();
    initPasswordToggles();
    bindRealtimeValidation();
    var form = qs('#wizard-form');
    if (form) loadCfBootstrapIntoForm(form);
    initApiBanner();
    var linkStep1 = qs('#link-deploy-worker-step1');
    var L = window.SSEO && window.SSEO.links;
    if (linkStep1 && L && L.platformWorkerDeploy) linkStep1.href = L.platformWorkerDeploy;
    initWizard();
    checkPlatformStatus();
  });

  window.showToast = showToast;
})();
