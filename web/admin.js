const loginForm = document.getElementById('loginForm');
const loginMessage = document.getElementById('loginMessage');
const dashboard = document.getElementById('dashboard');
const loginCard = document.getElementById('loginCard');
const feedbackBody = document.getElementById('feedbackBody');
const stats = document.getElementById('stats');
const summaryBox = document.getElementById('summaryBox');
const pageInfo = document.getElementById('pageInfo');
const themeToggle = document.getElementById('themeToggle');

let token = localStorage.getItem('adminToken') || '';
let page = 1;
let total = 0;

function applyTheme(theme) {
  document.documentElement.setAttribute('data-theme', theme);
  localStorage.setItem('theme', theme);
}

themeToggle.addEventListener('click', () => {
  const current = document.documentElement.getAttribute('data-theme') || 'light';
  applyTheme(current === 'light' ? 'dark' : 'light');
});
applyTheme(localStorage.getItem('theme') || 'light');

function authHeaders() {
  return { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` };
}

function renderStats(data) {
  stats.innerHTML = '';
  const items = [
    ['Total Feedback', data.totalFeedback || 0],
    ['Open Items', data.openItems || 0],
    ['Avg Priority', Number(data.averagePriority || 0).toFixed(1)],
    ['Most Common Tag', data.mostCommonTag || '-'],
  ];
  for (const [label, value] of items) {
    const card = document.createElement('div');
    card.className = 'stat';
    card.innerHTML = `<span>${label}</span><b>${value}</b>`;
    stats.appendChild(card);
  }
}

function sentimentBadge(value) {
  const sentiment = value || 'Neutral';
  return `<span class="badge ${sentiment}">${sentiment}</span>`;
}

async function loadFeedback() {
  const params = new URLSearchParams({
    page: String(page),
    limit: '10',
    category: document.getElementById('filterCategory').value,
    status: document.getElementById('filterStatus').value,
    sortBy: document.getElementById('sortBy').value,
    search: document.getElementById('search').value,
  });

  const response = await fetch(`/api/feedback?${params.toString()}`, { headers: authHeaders() });
  const result = await response.json();
  if (!response.ok || !result.success) {
    feedbackBody.innerHTML = `<tr><td colspan="7">${result.error || 'Failed to load data'}</td></tr>`;
    return;
  }

  const items = result.data.items || [];
  total = result.data.total || 0;
  pageInfo.textContent = `Page ${page} of ${Math.max(1, Math.ceil(total / 10))}`;
  renderStats(result.data.stats || {});

  feedbackBody.innerHTML = '';
  for (const item of items) {
    const row = document.createElement('tr');
    row.innerHTML = `
      <td>${item.title}</td>
      <td>${item.category}</td>
      <td>${sentimentBadge(item.ai_sentiment)}</td>
      <td>${item.ai_priority || '-'}</td>
      <td>
        <select data-action="status" data-id="${item.id}">
          <option ${item.status === 'New' ? 'selected' : ''}>New</option>
          <option ${item.status === 'In Review' ? 'selected' : ''}>In Review</option>
          <option ${item.status === 'Resolved' ? 'selected' : ''}>Resolved</option>
        </select>
      </td>
      <td>${new Date(item.createdAt).toLocaleDateString()}</td>
      <td>
        <button class="button secondary" data-action="reanalyze" data-id="${item.id}">Re-run AI</button>
        <button class="button secondary" data-action="delete" data-id="${item.id}">Delete</button>
      </td>`;
    feedbackBody.appendChild(row);
  }
}

async function updateStatus(id, status) {
  await fetch(`/api/feedback/${id}`, { method: 'PATCH', headers: authHeaders(), body: JSON.stringify({ status }) });
  loadFeedback();
}

async function deleteFeedback(id) {
  await fetch(`/api/feedback/${id}`, { method: 'DELETE', headers: authHeaders() });
  loadFeedback();
}

async function reanalyze(id) {
  await fetch(`/api/feedback/${id}/reanalyze`, { method: 'POST', headers: authHeaders() });
  loadFeedback();
}

feedbackBody.addEventListener('change', (event) => {
  const target = event.target;
  if (target.dataset.action === 'status') {
    updateStatus(target.dataset.id, target.value);
  }
});

feedbackBody.addEventListener('click', (event) => {
  const target = event.target;
  if (target.dataset.action === 'delete') {
    deleteFeedback(target.dataset.id);
  }
  if (target.dataset.action === 'reanalyze') {
    reanalyze(target.dataset.id);
  }
});

loginForm.addEventListener('submit', async (event) => {
  event.preventDefault();
  loginMessage.textContent = 'Signing in...';

  const response = await fetch('/api/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      email: document.getElementById('adminEmail').value,
      password: document.getElementById('adminPassword').value,
    }),
  });
  const result = await response.json();
  if (!response.ok || !result.success) {
    loginMessage.textContent = result.error || 'Login failed';
    return;
  }

  token = result.data.token;
  localStorage.setItem('adminToken', token);
  loginCard.hidden = true;
  dashboard.hidden = false;
  loadFeedback();
});

document.getElementById('refresh').addEventListener('click', loadFeedback);
document.getElementById('filterCategory').addEventListener('change', () => { page = 1; loadFeedback(); });
document.getElementById('filterStatus').addEventListener('change', () => { page = 1; loadFeedback(); });
document.getElementById('sortBy').addEventListener('change', loadFeedback);
document.getElementById('search').addEventListener('input', () => { page = 1; loadFeedback(); });

document.getElementById('prevPage').addEventListener('click', () => {
  if (page > 1) {
    page -= 1;
    loadFeedback();
  }
});

document.getElementById('nextPage').addEventListener('click', () => {
  if (page < Math.ceil(total / 10)) {
    page += 1;
    loadFeedback();
  }
});

document.getElementById('loadSummary').addEventListener('click', async () => {
  const response = await fetch('/api/feedback/summary', { headers: authHeaders() });
  const result = await response.json();
  summaryBox.textContent = result?.data?.summary || result.error || 'No summary';
});

if (token) {
  loginCard.hidden = true;
  dashboard.hidden = false;
  loadFeedback();
}

