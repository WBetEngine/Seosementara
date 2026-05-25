/**
 * GitHub Pages base path + URL admin setelah onboarding.
 * Ubah adminUrlAfterComplete jika domain admin berbeda.
 */
(function () {
  var path = location.pathname.replace(/\/$/, '');
  var base = '';
  var marker = '/Seosementara';
  var idx = path.indexOf(marker);
  if (idx >= 0) {
    base = marker;
  } else {
    var parts = path.split('/').filter(Boolean);
    if (parts.length > 1) {
      base = '/' + parts.slice(0, -1).join('/');
    }
  }

  window.SSEO = window.SSEO || {};
  window.SSEO.basePath = base;
  window.SSEO.platformApiBase = '';
  window.SSEO.adminUrlAfterComplete =
    'https://seosementara.org/admin/login.html?from=onboarding';
  window.SSEO.onboardingPagesUrl = 'https://wbetengine.github.io/Seosementara/';
  window.SSEO.wizardStorageKey = 'sseo_onboarding_wizard';

  window.SSEO.asset = function (rel) {
    var p = rel.charAt(0) === '/' ? rel : '/' + rel;
    return (window.SSEO.basePath || '') + p;
  };
})();
