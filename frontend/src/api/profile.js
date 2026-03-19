import { apiFetch } from "./http";

export async function getProfile() {
  return apiFetch("/api/me");
}
