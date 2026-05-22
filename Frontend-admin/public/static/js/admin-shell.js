/**
 * Seosementara Admin shell v2 — sidebar, drawer, page tabs, nav aktif
 */
(function () {
  function qs(sel, root) {
    return (root || document).querySelector(sel);
  }

  function qsa(sel, root) {
    return Array.prototype.slice.call((root || document).querySelectorAll(sel));
  }

  function openDrawer() {
    var drawer = qs("#app-drawer");
    var backdrop = qs("#drawer-backdrop");
    if (drawer) {
      drawer.classList.add("is-open");
      drawer.setAttribute("aria-hidden", "false");
    }
    if (backdrop) backdrop.classList.add("is-visible");
    document.body.style.overflow = "hidden";
  }

  function closeDrawer() {
    var drawer = qs("#app-drawer");
    var backdrop = qs("#drawer-backdrop");
    if (drawer) {
      drawer.classList.remove("is-open");
      drawer.setAttribute("aria-hidden", "true");
      drawer.innerHTML = "";
    }
    if (backdrop) backdrop.classList.remove("is-visible");
    document.body.style.overflow = "";
  }

  function openSidebar() {
    var sb = qs("#admin-sidebar");
    var ov = qs("#sidebar-overlay");
    if (sb) sb.classList.add("is-open");
    if (ov) ov.classList.add("is-visible");
  }

  function closeSidebar() {
    var sb = qs("#admin-sidebar");
    var ov = qs("#sidebar-overlay");
    if (sb) sb.classList.remove("is-open");
    if (ov) ov.classList.remove("is-visible");
  }

  function setActiveNav(route) {
    if (!route) return;
    qsa(".nav-link[data-nav]").forEach(function (el) {
      var r = el.getAttribute("data-nav");
      var match =
        route === r ||
        (r !== "/" && route.indexOf(r) === 0);
      if (r === "/admin/plugins/pixel" && route.indexOf("/admin/plugins/pixel") === 0) {
        match = true;
      }
      el.classList.toggle("is-active", match);
    });
  }

  function initPageTabs(root) {
    var scope = root || document;
    var tabs = qsa("[data-page-tab]", scope);
    if (!tabs.length) return;

    tabs.forEach(function (tab) {
      if (tab._pageTabBound) return;
      tab._pageTabBound = true;
      tab.addEventListener("click", function () {
        var container = tab.closest(".page-tabs") || scope;
        qsa("[data-page-tab]", container).forEach(function (t) {
          t.classList.toggle("is-active", t === tab);
        });
      });
    });

    syncPageTabsFromUrl(scope);
  }

  function syncPageTabsFromUrl(scope) {
    var path = window.location.pathname;
    var key = null;

    var domainMap = {
      "/admin/domain/shared": "shared",
      "/admin/domain/add": "add",
      "/admin/domain/all": "all",
      "/admin/domain": "mine",
    };
    if (domainMap[path]) {
      key = domainMap[path];
    } else if (path.indexOf("/admin/plugins/pixel/facebook") === 0) {
      var pixelMap = {
        "/admin/plugins/pixel/facebook/setup": "setup",
        "/admin/plugins/pixel/facebook/connection": "connection",
        "/admin/plugins/pixel/facebook/domains": "domains",
        "/admin/plugins/pixel/facebook/diagnostics": "diagnostics",
        "/admin/plugins/pixel/facebook/events": "events",
        "/admin/plugins/pixel/facebook/analytics": "analytics",
        "/admin/plugins/pixel/facebook": "overview",
      };
      key = pixelMap[path] || "overview";
    }

    if (!key) return;
    qsa("[data-page-tab]", scope).forEach(function (t) {
      t.classList.toggle("is-active", t.getAttribute("data-page-tab") === key);
    });
  }

  function initDrawerTabs(root) {
    var tabs = root.querySelectorAll("[data-drawer-tab]");
    var panels = root.querySelectorAll("[data-drawer-panel]");
    if (!tabs.length) return;
    tabs.forEach(function (tab) {
      tab.addEventListener("click", function () {
        var id = tab.getAttribute("data-drawer-tab");
        tabs.forEach(function (t) {
          t.classList.toggle("is-active", t === tab);
        });
        panels.forEach(function (p) {
          p.hidden = p.getAttribute("data-drawer-panel") !== id;
        });
      });
    });
  }

  window.openDrawer = openDrawer;
  window.closeDrawer = closeDrawer;
  window.setActiveNav = setActiveNav;

  document.addEventListener("DOMContentLoaded", function () {
    if (typeof htmx === "undefined") {
      console.error(
        "[Seosementara Admin] HTMX tidak termuat. Pastikan /static/js/htmx.min.js dapat diakses."
      );
      return;
    }

    document.body.addEventListener("htmx:responseError", function (ev) {
      console.error("[HTMX]", ev.detail.pathInfo.requestPath, ev.detail.xhr.status);
    });

    var menuBtn = qs("[data-toggle-sidebar]");
    var overlay = qs("#sidebar-overlay");
    if (menuBtn) menuBtn.addEventListener("click", openSidebar);
    if (overlay) overlay.addEventListener("click", closeSidebar);

    var backdrop = qs("#drawer-backdrop");
    if (backdrop) backdrop.addEventListener("click", closeDrawer);

    document.addEventListener("keydown", function (e) {
      if (e.key === "Escape") {
        closeDrawer();
        closeSidebar();
      }
    });

    qsa(".nav-group__title[data-collapse]").forEach(function (btn) {
      btn.addEventListener("click", function () {
        var expanded = btn.getAttribute("aria-expanded") !== "false";
        btn.setAttribute("aria-expanded", expanded ? "false" : "true");
        var list = btn.nextElementSibling;
        if (list) list.hidden = expanded;
      });
    });

    initPageTabs(document);
    var path = window.location.pathname;
    setActiveNav(path);
  });

  document.body.addEventListener("htmx:afterSwap", function (ev) {
    var target = ev.detail.target;
    if (!target) return;

    if (target.id === "app-drawer") {
      openDrawer();
      initDrawerTabs(target);
    }

    if (
      target.id === "main" ||
      target.id === "page-tab-panel" ||
      target.id === "pixel-tab-panel"
    ) {
      var root =
        target.id === "main"
          ? target
          : target.closest("#main") || document.getElementById("main") || document;
      initPageTabs(root);
      if (window.innerWidth < 1024) closeSidebar();
    }
  });

  document.body.addEventListener("htmx:pushedIntoHistory", function () {
    setActiveNav(window.location.pathname);
  });

  document.body.addEventListener("click", function (e) {
    if (e.target.closest("[data-close-drawer]")) {
      e.preventDefault();
      closeDrawer();
    }
    var nav = e.target.closest(".nav-link[data-nav]");
    if (nav) {
      setActiveNav(nav.getAttribute("data-nav"));
    }
  });
})();
