/**
 * URL Platform API — diisi otomatis setelah deploy worker (platform-api-url.js)
 * atau manual: ?api=https://sse-platform.<subdomain>.workers.dev
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

  window.SSEO.links = {
    githubTokens: 'https://github.com/settings/tokens',
    githubSecrets: 'https://github.com/WBetEngine/Seosementara/settings/secrets/actions',
    githubRunners: 'https://github.com/WBetEngine/Seosementara/settings/actions/runners',
    githubActions: 'https://github.com/WBetEngine/Seosementara/actions',
    cfDashboard: 'https://dash.cloudflare.com/',
    cfApiTokens: 'https://dash.cloudflare.com/profile/api-tokens',
    cfZeroTrust: 'https://one.dash.cloudflare.com/',
    cfTunnels: 'https://one.dash.cloudflare.com/',
    cfPages: 'https://dash.cloudflare.com/?to=/:account/pages',
    cfDns: 'https://dash.cloudflare.com/',
    cfDocsTunnel: 'https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/',
    ghcr: 'https://github.com/WBetEngine/Seosementara/pkgs/container'
  };

  window.SSEO.asset = function (rel) {
    var p = rel.charAt(0) === '/' ? rel : '/' + rel;
    return (window.SSEO.basePath || '') + p;
  };
})();
