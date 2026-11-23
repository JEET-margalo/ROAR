(() => {
  const esc = (v) => window.esc ? window.esc(v) : String(v ?? '');
  async function load() {
    const body = document.getElementById('website-logs-body');
    if (!body) return;
    try {
      const data = await window.apiFetch('/api/ui/logs');
      const entries = data.entries || [];
      if (!entries.length) {
        body.innerHTML = '<div class="empty-state"><div class="empty-title">No browser diagnostics</div><div class="empty-detail">Errors and rejected requests will appear here while you use the dashboard.</div></div>';
        return;
      }
      body.innerHTML = `<div class="website-log-list">${entries.slice().reverse().map((e) => {
        const level = String(e.type || 'event');
        const time = e.server_time || e.time || '';
        const message = e.message || 'No message';
        return `<article class="website-log-entry"><div class="website-log-meta"><span class="badge ${level === 'error' ? 'badge-dead' : 'badge-starting'}">${esc(level)}</span><time>${esc(time)}</time></div><pre>${esc(message)}</pre>${e.source ? `<div class="website-log-source">${esc(e.source)}:${esc(e.line || '')}</div>` : ''}</article>`;
      }).join('')}</div>`;
    } catch (error) {
      body.innerHTML = `<div class="empty-state"><div class="empty-title">Could not load logs</div><div class="empty-detail">${esc(error.message || error)}</div></div>`;
    }
  }
  function exportLogs() {
    const token = window.state?._authAccessToken || window.localTokenGet?.();
    fetch('/api/ui/logs/export', { headers: token ? { Authorization: `Bearer ${token}` } : {} }).then(async (r) => {
      if (!r.ok) throw new Error(`HTTP ${r.status}`);
      const blob = await r.blob(); const url = URL.createObjectURL(blob); const a = document.createElement('a');
      a.href = url; a.download = 'roar-website.log'; a.click(); URL.revokeObjectURL(url);
    }).catch((e) => window.showToast?.('error', 'Log export failed', e.message));
  }
  window.WebsiteLogsPage = { load, exportLogs };
  document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('website-logs-refresh')?.addEventListener('click', load);
    document.getElementById('website-logs-export')?.addEventListener('click', exportLogs);
  });
})();
