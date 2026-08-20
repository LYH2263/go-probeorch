async function j(url, opts) {
  const r = await fetch(url, opts);
  const t = await r.text();
  try { return JSON.parse(t); } catch { return { raw: t, status: r.status }; }
}

async function refresh() {
  const stats = await j('/api/stats');
  const snap = await j('/api/snapshot');
  document.getElementById('stats-line').textContent =
    `registered=${stats.Registered||0} enabled=${stats.Enabled||0} runs=${stats.TotalRuns||0}`;
  const ul = document.getElementById('target-list');
  ul.innerHTML = '';
  (snap.Targets || []).forEach(t => {
    const li = document.createElement('li');
    li.textContent = `${t.ID} · ${t.Name||'-'} · ${t.Address} · rate=${(t.SuccessRate||0).toFixed(2)} · ${t.Enabled?'on':'off'}`;
    ul.appendChild(li);
  });
}

document.getElementById('btn-refresh').onclick = refresh;

document.getElementById('reg-form').onsubmit = async (e) => {
  e.preventDefault();
  const body = {
    name: document.getElementById('reg-name').value,
    address: document.getElementById('reg-addr').value,
    kind: document.getElementById('reg-kind').value,
    interval: document.getElementById('reg-interval').value || '30s',
  };
  const out = await j('/api/targets/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
  document.getElementById('reg-out').textContent = JSON.stringify(out, null, 2);
  refresh();
};

document.getElementById('run-form').onsubmit = async (e) => {
  e.preventDefault();
  const id = document.getElementById('run-id').value;
  const out = await j('/api/targets/run?id=' + encodeURIComponent(id));
  document.getElementById('run-out').textContent = JSON.stringify(out, null, 2);
  refresh();
};

refresh();
