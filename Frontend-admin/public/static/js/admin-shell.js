/**
 * Seosementara Admin shell — sidebar + universal drawer (#app-drawer)
 * UI prototype; backend HTMX endpoints besok.
 */
(function () {
  function qs(sel, root) {
    return (root || document).querySelector(sel);
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

  window.openDrawer = openDrawer;
  window.closeDrawer = closeDrawer;

  document.addEventListener("DOMContentLoaded", function () {
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
  });

  document.body.addEventListener("htmx:afterSwap", function (ev) {
    if (ev.detail.target && ev.detail.target.id === "app-drawer") {
      openDrawer();
      initDrawerTabs(ev.detail.target);
    }
  });

  document.body.addEventListener("click", function (e) {
    if (e.target.closest("[data-close-drawer]")) {
      e.preventDefault();
      closeDrawer();
    }
  });

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

  document.body.addEventListener("htmx:afterSwap", function (ev) {
    if (ev.detail.target && ev.detail.target.id === "main") {
      var active = document.querySelector(".nav-group__items a.is-active, .nav-group__items button.is-active");
      if (window.innerWidth < 1024) closeSidebar();
    }
  });
})();
