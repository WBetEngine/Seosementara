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

  function checkBootstrapRedirect() {
    if (location.pathname.indexOf('bootstrap') >= 0 || location.pathname.indexOf('login') >= 0) return;
    try {
      var done = localStorage.getItem(window.SSEO.bootstrapKey);
      if (!done && location.pathname.indexOf('/admin/') >= 0) {
        /* opsional: redirect ke wizard pertama kali */
      }
    } catch (e) {}
  }

  document.addEventListener('DOMContentLoaded', function () {
    initSidebar();
    initBackdrop();
    initDemoForms();
    initHtmxHooks();
    checkBootstrapRedirect();
  });
})();
