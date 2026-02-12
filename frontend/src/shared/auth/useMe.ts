"use client";

import { useQuery } from "@tanstack/react-query";
import { authApi } from "@/shared/api/endpoints/auth";
import { useSession } from "@/shared/auth/store";

export function useMe() {
  const accessToken = useSession((s) => s.accessToken);

  return useQuery({
    queryKey: ["me"],
    queryFn: () => authApi.me(),
    enabled: Boolean(accessToken),
  });
}

export function useIsAdmin() {
  const meQ = useMe();
  const role = meQ.data?.role;
  const isAdmin = role === "admin" || role === "superadmin";
  return { ...meQ, isAdmin };
}

