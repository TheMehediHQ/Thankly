const API = import.meta.env.PUBLIC_API_URL || 'http://localhost:8080';

async function request(path: string, options: RequestInit = {}) {
  const token = localStorage.getItem('token');

  const res = await fetch(`${API}${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options.headers,
    },
  });

  if (res.status === 401) {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/login';
    return null;
  }

  return res;
}

export async function register(name: string, email: string, password: string) {
  const res = await request('/api/auth/register', {
    method: 'POST',
    body: JSON.stringify({ name, email, password }),
  });
  if (!res) return null;
  const data = await res.json();
  if (res.ok) {
    localStorage.setItem('token', data.token);
    localStorage.setItem('refresh_token', data.refresh_token);
    localStorage.setItem('user', JSON.stringify(data.user));
  }
  return { ok: res.ok, data };
}

export async function login(email: string, password: string) {
  const res = await request('/api/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  });
  if (!res) return null;
  const data = await res.json();
  if (res.ok) {
    localStorage.setItem('token', data.token);
    localStorage.setItem('refresh_token', data.refresh_token);
    localStorage.setItem('user', JSON.stringify(data.user));
  }
  return { ok: res.ok, data };
}

export async function getToday() {
  const res = await request('/api/journal/today');
  if (!res) return null;
  return res.json();
}

export async function createGratitude(content: string) {
  const res = await request('/api/gratitudes', {
    method: 'POST',
    body: JSON.stringify({ content }),
  });
  if (!res) return null;
  const data = await res.json();
  return { ok: res.ok, data };
}

export async function getHistory() {
  const res = await request('/api/journal/history');
  if (!res) return null;
  return res.json();
}

export async function getProfile() {
  const res = await request('/api/me');
  if (!res) return null;
  return res.json();
}

export function getUser() {
  const user = localStorage.getItem('user');
  return user ? JSON.parse(user) : null;
}

export function isLoggedIn() {
  return !!localStorage.getItem('token');
}

export function logout() {
  localStorage.removeItem('token');
  localStorage.removeItem('refresh_token');
  localStorage.removeItem('user');
  window.location.href = '/login';
}
