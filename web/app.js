const form = document.getElementById('feedbackForm');
const message = document.getElementById('formMessage');
const description = document.getElementById('description');
const descCount = document.getElementById('descCount');
const themeToggle = document.getElementById('themeToggle');

function applyTheme(theme) {
  document.documentElement.setAttribute('data-theme', theme);
  localStorage.setItem('theme', theme);
}

themeToggle.addEventListener('click', () => {
  const current = document.documentElement.getAttribute('data-theme') || 'light';
  applyTheme(current === 'light' ? 'dark' : 'light');
});

applyTheme(localStorage.getItem('theme') || 'light');

description.addEventListener('input', () => {
  descCount.textContent = String(description.value.length);
});

form.addEventListener('submit', async (event) => {
  event.preventDefault();
  message.textContent = 'Submitting...';

  const payload = {
    title: document.getElementById('title').value.trim(),
    description: description.value.trim(),
    category: document.getElementById('category').value,
    name: document.getElementById('name').value.trim(),
    email: document.getElementById('email').value.trim(),
  };

  if (!payload.title || payload.description.length < 20) {
    message.textContent = 'Title is required and description must be at least 20 characters.';
    return;
  }

  try {
    const response = await fetch('/api/feedback', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });
    const result = await response.json();
    if (!response.ok || !result.success) {
      message.textContent = result.error || 'Submission failed.';
      return;
    }
    message.textContent = 'Feedback submitted successfully.';
    form.reset();
    descCount.textContent = '0';
  } catch (error) {
    message.textContent = 'Network error. Please try again.';
  }
});

