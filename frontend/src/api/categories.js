import { apiFetch } from "./http";

export function getCategories() {
  return apiFetch("/api/categories");
}