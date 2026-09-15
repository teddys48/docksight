const API_BASE = '/api';

export async function fetchSystemStats() {
  const res = await fetch(`${API_BASE}/system/stats`);
  if (!res.ok) throw new Error('Failed to fetch system stats');
  return res.json();
}

export async function fetchSystemHistory(limit = 60) {
  const res = await fetch(`${API_BASE}/system/history?limit=${limit}`);
  if (!res.ok) throw new Error('Failed to fetch history');
  return res.json();
}

export async function fetchContainers() {
  const res = await fetch(`${API_BASE}/containers`);
  if (!res.ok) throw new Error('Failed to fetch containers');
  return res.json();
}

export async function performContainerAction(id, action, force = false) {
  const url = `${API_BASE}/containers/${id}/${action}${force ? '?force=true' : ''}`;
  const res = await fetch(url, { method: 'POST' });
  const data = await res.json();
  if (!res.ok || !data.success) {
    throw new Error(data.message || `Action ${action} failed`);
  }
  return data;
}

export async function inspectContainer(id) {
  const res = await fetch(`${API_BASE}/containers/${id}/inspect`);
  if (!res.ok) throw new Error('Failed to inspect container');
  return res.json();
}

export async function fetchImages() {
  const res = await fetch(`${API_BASE}/images`);
  if (!res.ok) throw new Error('Failed to fetch images');
  return res.json();
}

export async function fetchVolumes() {
  const res = await fetch(`${API_BASE}/volumes`);
  if (!res.ok) throw new Error('Failed to fetch volumes');
  return res.json();
}

// Formatting helpers
export function formatBytes(bytes, decimals = 2) {
  if (!bytes || bytes === 0) return '0 B';
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

export function formatUptime(createdUnix) {
  if (!createdUnix) return 'N/A';
  const now = Math.floor(Date.now() / 1000);
  const diff = now - createdUnix;
  if (diff < 60) return `${diff}s`;
  if (diff < 3600) return `${Math.floor(diff / 60)}m`;
  if (diff < 86400) return `${Math.floor(diff / 3600)}h`;
  return `${Math.floor(diff / 86400)}d`;
}

export function formatDate(createdUnix) {
  if (!createdUnix) return 'N/A';
  if (typeof createdUnix === 'string') {
    const d = new Date(createdUnix);
    return isNaN(d.getTime()) ? createdUnix : d.toLocaleString();
  }
  return new Date(createdUnix * 1000).toLocaleString();
}
