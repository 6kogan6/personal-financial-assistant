import { apiFetch } from "./http";

export function getSummary(month) {
  return apiFetch(`/api/reports/summary?month=${encodeURIComponent(month)}`);
}