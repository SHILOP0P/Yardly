"use client";

import type { ReactNode } from "react";
import Link from "next/link";
import { useHydrateSession } from "@/shared/auth/useHydrateSession";
import { useSession } from "@/shared/auth/store";
import { useIsAdmin } from "@/shared/auth/useMe";

type Props = {
  children: ReactNode;
};

export function AdminGate({ children }: Props) {
  useHydrateSession();
  const accessToken = useSession((s) => s.accessToken);
  const meQ = useIsAdmin();

  if (!accessToken) {
    return (
      <div className="p-6 space-y-3">
        <h1 className="text-2xl font-semibold">Админка</h1>
        <p className="opacity-70">Нужна авторизация.</p>
        <Link href="/auth/login" className="inline-flex border rounded-lg px-3 py-2">
          Войти
        </Link>
      </div>
    );
  }

  if (meQ.isLoading) {
    return <div className="p-6">Загрузка профиля...</div>;
  }

  if (meQ.error) {
    return <div className="p-6 text-red-600">Не удалось проверить роль: {(meQ.error as Error).message}</div>;
  }

  if (!meQ.isAdmin) {
    return (
      <div className="p-6 space-y-3">
        <h1 className="text-2xl font-semibold">Админка</h1>
        <p className="text-red-600">Доступ запрещен. Нужна роль admin/superadmin.</p>
      </div>
    );
  }

  return <>{children}</>;
}

