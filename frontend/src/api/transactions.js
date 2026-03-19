import { apiFetch } from "./http";

function buildQuery(params = {}) {
  const search = new URLSearchParams();

  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== null && value !== "") {
      search.set(key, value);
    }
  }

  const qs = search.toString();
  return qs ? `?${qs}` : "";
}

export async function getTransactions(params = {}) {
  return apiFetch(`/api/transactions${buildQuery(params)}`);
}

export async function createTransaction(payload) {
  return apiFetch("/api/transactions", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export async function deleteTransactionRequest(id) {
  const accessToken =
    localStorage.getItem("access_token") ||
    localStorage.getItem("token") ||
    sessionStorage.getItem("access_token") ||
    sessionStorage.getItem("token") ||
    "";

  const headers = {};
  if (accessToken) {
    headers.Authorization = `Bearer ${accessToken}`;
  }

  const response = await fetch(`/api/transactions/${id}`, {
    method: "DELETE",
    headers,
    credentials: "include",
  });

  if (!response.ok) {
    let message = "Не удалось удалить транзакцию";

    try {
      const text = await response.text();
      if (text) {
        message = text;
      }
    } catch {
      // ignore
    }

    throw new Error(message);
  }

  return true;
}
