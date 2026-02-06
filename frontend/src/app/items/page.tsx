"use client";

import Link from "next/link";
import { useItemsList } from "@/features/item/hooks";

export default function ItemsPage() {
  const { data, isLoading, error } = useItemsList();

  if (isLoading) return <div className="p-6">Загрузка...</div>;
  if (error) return <div className="p-6">Ошибка</div>;

  return (
    <div className="p-6 space-y-4">
      <h1 className="text-2xl font-semibold">Лента вещей</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        {data?.map((it) => (
          <Link key={it.id} href={`/items/${it.id}`} className="border rounded-xl p-4 hover:bg-white/5">
            <div className="font-medium">{it.title}</div>
            <div className="text-sm opacity-70">{it.mode} • {it.status}</div>
          </Link>
        ))}
      </div>
    </div>
  );
}
