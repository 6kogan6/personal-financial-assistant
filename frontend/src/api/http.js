import { getToken } from "./token";

export async function apiFetch(path, options = {}) {
  const headers = new Headers(options.headers || {});
  headers.set("Content-Type", "application/json");

  const token = getToken();
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  const res = await fetch(path, { ...options, headers });

  // Попробуем прочитать JSON (и для ошибок тоже)
  let data = null;
  const text = await res.text();
  if (text) {
    try {
      data = JSON.parse(text);
    } catch {
      data = { raw: text };
    }
  }

  if (!res.ok) {
    const msg = data?.error || `HTTP ${res.status}`;
    throw new Error(msg);
  }

  return data;
}