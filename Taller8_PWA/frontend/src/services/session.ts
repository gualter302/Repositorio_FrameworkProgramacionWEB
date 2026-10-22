export const SESSION_KEY = 'multicatalogo_sesion';
export const SESSION_EVENT = 'multicatalogo-session';
export function readSession() {
  try {
    const value = JSON.parse(localStorage.getItem(SESSION_KEY) ?? 'null');
    if (!value || typeof value.email !== 'string' || typeof value.token !== 'string') return null;
    const claims = JSON.parse(atob(value.token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')));
    if (!claims.exp || claims.exp * 1000 <= Date.now()) return null;
    return value as { email: string; token: string; nombre: string; rol: string };
  } catch { return null; }
}
export function expireSession() {
  localStorage.removeItem(SESSION_KEY);
  window.dispatchEvent(new Event(SESSION_EVENT));
}
