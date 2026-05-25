(function () {
  'use strict';

  var MAX_STEP = 8;

  var STEP_FIELDS = {
    1: ['github_pat'],
    2: ['cf_token', 'cf_account', 'cf_zone', 'primary_domain'],
    3: ['ssh_host', 'ssh_port', 'ssh_user', 'ssh_secret'],
    4: ['runner_label'],
    5: ['tunnel_name'],
    6: ['db_password', 'master_key'],
    7: [],
    8: []
  };

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
      if (!v || !v.trim()) return 'API Token / Global API Key wajib diisi.';
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
        return 'Format domain tidak valid (contoh: seosementara.org).';
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
    }, 4000);
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

  function bindRealtimeValidation() {
    Object.keys(VALIDATORS).forEach(function (name) {
      var input = getInput(name);
      if (!input) return;
      input.addEventListener('input', function () {
        validateField(name, true);
        updateStepStatus();
      });
      input.addEventListener('blur', function () {
        validateField(name, true);
        updateStepStatus();
      });
    });
  }

  function initInfoIcons() {
    qsa('.info-icon').forEach(function (btn) {
      var tipText = btn.getAttribute('data-tip') || '';
      var group = btn.closest('.form-group');
      var tipEl = group ? qs('.info-tip', group) : null;

      function toggleTip(open) {
        if (!tipEl) return;
        var show = open !== undefined ? open : !tipEl.classList.contains('is-open');
        tipEl.classList.toggle('is-open', show);
        tipEl.textContent = show ? tipText : '';
        btn.setAttribute('aria-expanded', show ? 'true' : 'false');
      }

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
        toggleTip();
      });
    });

    document.addEventListener('click', function (e) {
      if (!e.target.closest('.info-icon') && !e.target.closest('.info-tip')) {
        qsa('.info-tip.is-open').forEach(function (t) {
          t.classList.remove('is-open');
          t.textContent = '';
        });
        qsa('.info-icon').forEach(function (b) {
          b.setAttribute('aria-expanded', 'false');
        });
      }
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

  function initWizard() {
    var form = qs('#wizard-form');
    if (!form) return;

    var step = 1;
    var prev = qs('#btn-prev');
    var next = qs('#btn-next');
    var statusEl = qs('#step-status');

    function updateStepStatus() {
      if (!statusEl) return;
      if (step >= MAX_STEP) {
        statusEl.textContent = '';
        statusEl.className = 'step-status';
        return;
      }
      var ok = validateStep(step, false);
      statusEl.textContent = ok ? 'Langkah siap dilanjutkan' : 'Lengkapi field yang ditandai merah';
      statusEl.className = 'step-status' + (ok ? ' is-ok' : '');
      if (next) next.disabled = !ok;
    }

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
      if (s < MAX_STEP) {
        validateStep(s, false);
        updateStepStatus();
      }
    }

    if (prev) {
      prev.addEventListener('click', function () {
        if (step > 1) show(step - 1);
      });
    }

    if (next) {
      next.addEventListener('click', function () {
        if (!validateStep(step, true)) {
          showToast('Perbaiki field yang belum valid sebelum lanjut.', 'error');
          var first = qs('.form-control.is-invalid', form);
          if (first) first.focus();
          return;
        }
        if (step < MAX_STEP) show(step + 1);
      });
    }

    qsa('[data-test-step]').forEach(function (btn) {
      btn.addEventListener('click', function () {
        if (!validateStep(step, true)) {
          showToast('Lengkapi form yang valid dulu.', 'error');
          return;
        }
        var label = btn.getAttribute('data-test-step') || 'Koneksi';
        showToast(label + ': OK (demo — Workers API belum aktif)', 'success');
      });
    });

    var finish = qs('#btn-finish');
    if (finish) {
      finish.addEventListener('click', function (e) {
        e.preventDefault();
        try {
          var data = {};
          new FormData(form).forEach(function (v, k) {
            if (v) data[k] = v;
          });
          localStorage.setItem(window.SSEO.wizardStorageKey, JSON.stringify(data));
        } catch (err) {}

        showToast('Setup selesai — membuka admin…', 'success');
        var url =
          (window.SSEO && window.SSEO.adminUrlAfterComplete) ||
          'https://seosementara.org/admin/login.html?from=onboarding';
        setTimeout(function () {
          window.location.href = url;
        }, 1200);
      });
    }

    show(1);
  }

  document.addEventListener('DOMContentLoaded', function () {
    var css = qs('#onboarding-css');
    if (css && window.SSEO && window.SSEO.asset) {
      css.href = window.SSEO.asset('/assets/css/onboarding.css');
    }
    initExternalLinks();
    initInfoIcons();
    bindRealtimeValidation();
    initWizard();
  });

  window.showToast = showToast;
})();
