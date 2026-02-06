import { apiFetch } from "@/shared/api/client";
import type { Tokens } from "@/shared/api/types";

export const authApi = {
  register: (email: string, password: string) =>
    apiFetch<Tokens>(
      "/api/auth/register",
      {
        method: "POST",
        body: JSON.stringify({ email, password }),
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

  me: () => apiFetch<any>("/api/users/me"),
};
