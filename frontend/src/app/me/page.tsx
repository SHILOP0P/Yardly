"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { authApi } from "@/shared/api/endpoints/auth";
import { useHydrateSession } from "@/shared/auth/useHydrateSession";

export default function MePage() {
  useHydrateSession();

  const meQ = useQuery({
    queryKey: ["me"],
    queryFn: () => authApi.me(),
  });

  if (meQ.isLoading) return <div className="p-6">Загрузка...</div>;
  if (meQ.error) return <div className="p-6">Ошибка / 401 если не залогинен</div>;

  const me = meQ.data;

  return (
    <div className="p-6 space-y-4">
      <h1 className="text-2xl font-semibold">Профиль</h1>

      <div className="border rounded-xl p-4">
        <div className="text-sm opacity-70">Вы вошли как:</div>
        <div className="font-medium">
          id: {me?.id ?? "—"} • {me?.email ?? "—"}
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        <Link className="border rounded-xl p-4 hover:bg-white/5" href="/my/items">
          <div className="font-medium">Мои товары</div>
          <div className="text-sm opacity-70">GET /api/my/items + создание + фото</div>
        </Link>

        <Link className="border rounded-xl p-4 hover:bg-white/5" href="/my/bookings">
          <div className="font-medium">Мои заявки</div>
          <div className="text-sm opacity-70">GET /api/my/bookings + cancel</div>
        </Link>

        <Link className="border rounded-xl p-4 hover:bg-white/5" href="/my/items/bookings">
          <div className="font-medium">Заявки на мои вещи</div>
          <div className="text-sm opacity-70">approve / handover / return</div>
        </Link>

        <Link className="border rounded-xl p-4 hover:bg-white/5" href="/my/favorites">
          <div className="font-medium">Избранное</div>
          <div className="text-sm opacity-70">GET /api/my/favorites</div>
        </Link>
      </div>

      <Link className="inline-flex border rounded-lg px-3 py-2" href="/settings/security">
        Безопасность (logout / logout_all)
      </Link>
    </div>
  );
}
