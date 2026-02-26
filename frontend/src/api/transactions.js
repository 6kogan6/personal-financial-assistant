import { apiFetch } from "./http";

export function getTransactions(params = {}) {
  const qs = new URLSearchParams(params).toString();
  const url = qs ? `/api/transactions?${qs}` : "/api/transactions";
  return apiFetch(url);
}

export function createTransaction(tx) {
  return apiFetch("/api/transactions", {
    method: "POST",
    body: JSON.stringify(tx),
  });
}