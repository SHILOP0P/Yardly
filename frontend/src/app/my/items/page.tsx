"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { itemsApi } from "@/shared/api/endpoints/items";
import { useHydrateSession } from "@/shared/auth/useHydrateSession";

export default function MyItemsPage() {
  useHydrateSession();

  const q = useQuery({
    queryKey: ["my", "items"],
    queryFn: () => itemsApi.myItems(),
  });

  if (q.isLoading) return <div className="p-6">Загрузка...</div>;
  if (q.error) return <div className="p-6">Ошибка / 401</div>;

  return (
    <div className="p-6 space-y-4">
      <h1 className="text-2xl font-semibold">Мои вещи</h1>
      <Link className="border rounded-lg px-3 py-2" href="/my/items/new">
        + Создать товар
      </Link>
      <div className="space-y-3">
        {q.data?.map((it: any) => (
          <div key={it.id} className="border rounded-xl p-4 space-y-2">
            <div className="font-medium">{it.title}</div>
            <div className="text-sm opacity-70">{it.mode} • {it.status}</div>

            <div className="flex flex-wrap gap-2">
              <Link className="border rounded-lg px-3 py-2" href={`/items/${it.id}`}>Открыть</Link>
              <Link className="border rounded-lg px-3 py-2" href={`/items/${it.id}/upcoming`}>Upcoming</Link>
              <Link className="border rounded-lg px-3 py-2" href={`/my/items/booking-requests`}>Заявки</Link>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
