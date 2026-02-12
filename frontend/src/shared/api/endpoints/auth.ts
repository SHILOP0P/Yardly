import { apiFetch } from "@/shared/api/client";
import type { Me, RegisterResponse, Tokens } from "@/shared/api/types";

export const authApi = {
  register: (payload: { email: string; password: string; first_name: string; last_name?: string }) =>
    apiFetch<RegisterResponse>(
      "/api/auth/register",
      {
        method: "POST",
        body: JSON.stringify(payload),
      },
      { auth: false }
    ),

  login: (email: string, password: string) =>
    apiFetch<Tokens>(
      "/api/auth/login",
      {
        method: "POST",
        body: JSON.stringify({ email, password }),
      },
      { auth: false }
    ),

  me: () => apiFetch<Me>("/api/users/me"),
};
