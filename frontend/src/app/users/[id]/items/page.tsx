"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { itemsApi } from "@/shared/api/endpoints/items";
import { useQuery } from "@tanstack/react-query";

export default function UserItemsPage() {
  const params = useParams<{ id: string }>();
  const userId = Number(params.id);

  const q = useQuery({
    queryKey: ["user", userId, "items"],
    queryFn: () => itemsApi.byOwner(userId),
    enabled: userId > 0,
  });

  if (q.isLoading) return <div className="p-6">Загрузка...</div>;
  if (q.error) return <div className="p-6">Ошибка</div>;

  return (
    <div className="p-6 space-y-4">
      <h1 className="text-2xl font-semibold">Вещи пользователя #{userId}</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        {q.data?.map((it: any) => (
          <Link key={it.id} href={`/items/${it.id}`} className="border rounded-xl p-4 hover:bg-white/5">
            <div className="font-medium">{it.title}</div>
            <div className="text-sm opacity-70">{it.mode} • {it.status}</div>
          </Link>
        ))}
      </div>
    </div>
  );
}
