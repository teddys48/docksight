import { writable } from 'svelte/store';

// Theme state (dark/light)
const initialTheme = typeof window !== 'undefined' 
  ? (localStorage.getItem('theme') || 'dark') 
  : 'dark';

export const theme = writable(initialTheme);

if (typeof window !== 'undefined') {
  theme.subscribe(val => {
    localStorage.setItem('theme', val);
    if (val === 'dark') {
      document.documentElement.classList.add('dark');
      document.documentElement.classList.remove('light');
    } else {
      document.documentElement.classList.remove('dark');
      document.documentElement.classList.add('light');
    }
  });
}

// SSE Connection state
export const sseConnected = writable(false);

// Active Tab
export const activeTab = writable('dashboard');

// Toast notifications
export const toasts = writable([]);

export function addToast(message, type = 'info', duration = 4000) {
  const id = Date.now() + Math.random();
  toasts.update(all => [...all, { id, message, type }]);
  setTimeout(() => {
    toasts.update(all => all.filter(t => t.id !== id));
  }, duration);
}
