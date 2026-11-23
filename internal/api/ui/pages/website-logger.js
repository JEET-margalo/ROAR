(() => {
  // Browser diagnostics are intentionally compact and redacted server-side.
  // They are useful for exporting a reproducible UI failure to an AI reviewer.
  let lastSent = 0;
  const send = (entry) => {
    const now = Date.now();
    if (now - lastSent < 750) return;
    lastSent = now;
    const payload = { ...entry, url: location.href, user_agent: navigator.userAgent };
    const body = JSON.stringify(payload);
    if (navigator.sendBeacon) {
      try { navigator.sendBeacon('/api/ui/logs', new Blob([body], { type: 'application/json' })); return; } catch (_) { /* fetch fallback */ }
    }
    fetch('/api/ui/logs', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body, keepalive: true }).catch(() => {});
  };
  window.addEventListener('error', (event) => send({ type: 'error', message: event.message || 'window error', source: event.filename, line: event.lineno, column: event.colno }));
  window.addEventListener('unhandledrejection', (event) => send({ type: 'unhandledrejection', message: String(event.reason?.stack || event.reason || 'unhandled rejection') }));
  window.ROARWebsiteLogger = { log: (message, context = {}) => send({ type: 'manual', message, context }) };
  document.addEventListener('DOMContentLoaded', () => {
    const button = document.getElementById('website-log-export');
    if (!button) return;
    button.addEventListener('click', async () => {
      try {
        const token = window.state?._authAccessToken || window.localTokenGet?.();
        const headers = token ? { Authorization: `Bearer ${token}` } : {};
        const response = await fetch('/api/ui/logs/export', { headers });
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        const blob = await response.blob();
        const url = URL.createObjectURL(blob);
        const link = document.createElement('a'); link.href = url; link.download = 'roar-website.log'; link.click();
        URL.revokeObjectURL(url);
      } catch (error) {
        window.showToast?.('error', 'Log export failed', error.message || String(error));
      }
    });
  });
})();
