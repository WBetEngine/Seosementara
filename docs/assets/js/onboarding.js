(function () {
  'use strict';

  var MAX_STEP = 8;

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

  function applyAssetUrls() {
    qsa('[data-asset]').forEach(function (el) {
      var rel = el.getAttribute('data-asset');
      if (rel && window.SSEO && window.SSEO.asset) {
        el.setAttribute('href', window.SSEO.asset(rel));
      }
    });
  }

  function initWizard() {
    var form = qs('#wizard-form');
    if (!form) return;

    var step = 1;
    var prev = qs('#btn-prev');
    var next = qs('#btn-next');

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
        next.textContent = s >= MAX_STEP ? 'Selesai' : 'Lanjut';
        next.style.display = s >= MAX_STEP ? 'none' : '';
      }
    }

    if (prev) {
      prev.addEventListener('click', function () {
        if (step > 1) show(step - 1);
      });
    }

    if (next) {
      next.addEventListener('click', function () {
        if (step < MAX_STEP) show(step + 1);
      });
    }

    qsa('[data-test-step]').forEach(function (btn) {
      btn.addEventListener('click', function () {
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
        } catch (e) {}

        show(MAX_STEP);
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
    applyAssetUrls();
    var css = qs('#onboarding-css');
    if (css && window.SSEO && window.SSEO.asset) {
      css.href = window.SSEO.asset('/assets/css/onboarding.css');
    }
    initWizard();
  });

  window.showToast = showToast;
})();
