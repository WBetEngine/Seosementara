/** Mode API: 'mock' (fase 0) | 'live' (backend Go) */
window.SSEO = window.SSEO || {};
window.SSEO.apiMode = 'mock';
window.SSEO.apiBase = window.SSEO.apiMode === 'mock' ? '/mock-api' : '';
window.SSEO.appName = 'Seosementara';
/** URL GitHub Pages onboarding (first boot) */
window.SSEO.onboardingUrl = 'https://wbetengine.github.io/Seosementara/';
/** sessionStorage: diset saat redirect ?from=onboarding dari wizard */
window.SSEO.setupCompleteKey = 'sseo_setup_complete';
