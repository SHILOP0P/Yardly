"use client";

import Link from "next/link";
import { useHydrateSession } from "@/shared/auth/useHydrateSession";
import { useMyFavorites } from "@/features/favorite/hooks";

export default function MyFavoritesPage() {
  useHydrateSession();
  const q = useMyFavorites();

  if (q.isLoading) return <div className="p-6">Загрузка...</div>;
  if (q.error) return <div className="p-6">Ошибка. Если не залогинен — будет 401.</div>;

  return (
    <div className="p-6 space-y-4">
      <h1 className="text-2xl font-semibold">Моё избранное</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        {q.data?.map((it) => (
        <Link
            key={it.id}                         // ✅ key обязательно и уникальный
            href={`/items/${it.id}`}
            className="border rounded-xl p-4 hover:bg-white/5"
        >
            <div className="font-medium">{it.title}</div>
            <div className="text-sm opacity-70">
            {it.mode} • {it.status}
            </div>
        </Link>
        ))}

      </div>
    </div>
  );
}
