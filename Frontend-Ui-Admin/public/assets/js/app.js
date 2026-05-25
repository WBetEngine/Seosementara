(function () {
  'use strict';

  function qs(sel, root) {
    return (root || document).querySelector(sel);
  }

  function qsa(sel, root) {
    return Array.from((root || document).querySelectorAll(sel));
  }

  window.openDrawer = function () {
    var drawer = qs('#app-drawer');
    var backdrop = qs('#drawer-backdrop');
    if (drawer) {
      drawer.classList.add('is-open');
      drawer.setAttribute('aria-hidden', 'false');
    }
    if (backdrop) backdrop.hidden = false;
    document.body.classList.add('drawer-open');
  };

  window.closeDrawer = function () {
    var drawer = qs('#app-drawer');
    var backdrop = qs('#drawer-backdrop');
    if (drawer) {
      drawer.classList.remove('is-open', 'drawer--wide');
      drawer.setAttribute('aria-hidden', 'true');
      drawer.innerHTML = '';
    }
    if (backdrop) backdrop.hidden = true;
    document.body.classList.remove('drawer-open');
  };

  window.showToast = function (message, type) {
    type = type || 'info';
    var toast = qs('#toast');
    if (!toast) return;
    toast.innerHTML =
      '<div class="toast toast--' +
      type +
      '" role="status">' +
      message +
      '</div>';
    toast.classList.add('is-visible');
    setTimeout(function () {
      toast.classList.remove('is-visible');
    }, 4000);
  };

  window.apiUrl = function (path) {
    var base = window.SSEO && window.SSEO.apiBase ? window.SSEO.apiBase : '/mock-api';
    if (path.indexOf('/api/') === 0) return base + path;
    return path;
  };

  function initSidebar() {
    var toggle = qs('[data-sidebar-toggle]');
    var sidebar = qs('#sidebar');
    if (toggle && sidebar) {
      toggle.addEventListener('click', function () {
        sidebar.classList.toggle('is-open');
        document.body.classList.toggle('sidebar-open');
      });
    }

    qsa('.nav-group__toggle').forEach(function (btn) {
      btn.addEventListener('click', function () {
        var group = btn.closest('.nav-group');
        if (group) group.classList.toggle('is-collapsed');
      });
    });

    var path = location.pathname.replace(/\.html$/, '');
    qsa('.nav-link').forEach(function (link) {
      var href = link.getAttribute('href') || '';
      if (href && path.endsWith(href.replace(/\.html$/, '').replace(/^\//, ''))) {
        link.classList.add('is-active');
        var group = link.closest('.nav-group');
        if (group) group.classList.remove('is-collapsed');
      }
    });
  }

  function initBackdrop() {
    var backdrop = qs('#drawer-backdrop');
    if (backdrop) {
      backdrop.addEventListener('click', closeDrawer);
    }
    document.addEventListener('keydown', function (e) {
      if (e.key === 'Escape') closeDrawer();
    });
  }

  function initDemoForms() {
    document.body.addEventListener('submit', function (e) {
      var form = e.target;
      if (!form.matches('[data-demo-form]')) return;
      e.preventDefault();
      var data = {};
      new FormData(form).forEach(function (v, k) {
        data[k] = v;
      });
      var key = form.getAttribute('data-storage-key');
      if (key) {
        try {
          localStorage.setItem(key, JSON.stringify(data));
        } catch (err) {}
      }
      showToast('Tersimpan (mode demo — backend belum aktif)', 'success');
      if (form.hasAttribute('data-close-drawer')) closeDrawer();
    });
  }

  function initHtmxHooks() {
    document.body.addEventListener('htmx:afterSwap', function (e) {
      if (e.detail.target && e.detail.target.id === 'app-drawer') {
        openDrawer();
      }
    });
    document.body.addEventListener('htmx:configRequest', function (e) {
      var path = e.detail.path;
      if (path && path.indexOf('/api/') === 0 && window.SSEO.apiMode === 'mock') {
        e.detail.path = window.SSEO.apiBase + path;
      }
      var domainId = localStorage.getItem('active_managed_domain_id');
      if (domainId) e.detail.headers['X-Managed-Domain-ID'] = domainId;
    });
  }

  function captureOnboardingReturn() {
    if (location.search.indexOf('from=onboarding') < 0) return;
    try {
      sessionStorage.setItem(
        (window.SSEO && window.SSEO.setupCompleteKey) || 'sseo_setup_complete',
        '1'
      );
    } catch (e) {}
    try {
      var u = new URL(location.href);
      u.searchParams.delete('from');
      history.replaceState(null, '', u.pathname + u.search);
    } catch (err) {}
  }

  function initSetupBanner() {
    if (location.pathname.indexOf('/admin/') < 0) return;
    if (location.pathname.indexOf('login') >= 0) return;
    var key = (window.SSEO && window.SSEO.setupCompleteKey) || 'sseo_setup_complete';
    try {
      if (sessionStorage.getItem(key) === '1') return;
    } catch (e) {}
    var url =
      (window.SSEO && window.SSEO.onboardingUrl) ||
      'https://wbetengine.github.io/Seosementara/';
    var el = document.createElement('div');
    el.className = 'alert alert--warning setup-infra-banner';
    el.setAttribute('role', 'alert');
    el.innerHTML =
      '<strong>Setup infrastruktur belum selesai.</strong> Lanjutkan di ' +
      '<a href="' +
      url +
      '" target="_blank" rel="noopener">GitHub Pages onboarding</a>.';
    var main = document.querySelector('.admin-main');
    if (main) {
      main.insertBefore(el, main.firstChild);
    } else {
      document.body.insertBefore(el, document.body.firstChild);
    }
  }

  document.addEventListener('DOMContentLoaded', function () {
    captureOnboardingReturn();
    initSidebar();
    initBackdrop();
    initDemoForms();
    initHtmxHooks();
    initSetupBanner();
  });
})();
