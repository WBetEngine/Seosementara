/** Mode API: 'mock' (fase 0) | 'live' (backend Go) */
window.SSEO = window.SSEO || {};
window.SSEO.apiMode = 'mock';
window.SSEO.apiBase = window.SSEO.apiMode === 'mock' ? '/mock-api' : '';
window.SSEO.appName = 'Seosementara';
window.SSEO.bootstrapKey = 'sseo_bull_bootstrap_v1';
