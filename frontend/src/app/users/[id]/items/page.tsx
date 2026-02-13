"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { itemsApi } from "@/shared/api/endpoints/items";
import { useQuery } from "@tanstack/react-query";
import { formatDealMode, formatItemStatus } from "@/features/item/presentation";

export default function UserItemsPage() {
  const params = useParams<{ id: string }>();
  const userId = Number(params.id);

  const q = useQuery({
    queryKey: ["user", userId, "items"],
    queryFn: () => itemsApi.byOwner(userId),
    enabled: userId > 0,
  });

  if (q.isLoading) return <div className="p-6">Loading...</div>;
  if (q.error) return <div className="p-6">Error</div>;

  return (
    <div className="p-6 space-y-4">
      <h1 className="text-2xl font-semibold">User products #{userId}</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        {q.data?.map((it) => (
          <Link key={it.id} href={`/items/${it.id}`} className="border rounded-xl p-4 hover:bg-white/5">
            <div className="font-medium">{it.title}</div>
            <div className="text-sm opacity-70">{formatDealMode(it.mode)} | {formatItemStatus(it.status)}</div>
            {it.description ? <div className="text-sm mt-1 line-clamp-2">{it.description}</div> : null}
            <div className="text-sm opacity-80 mt-1">
              price: {it.price} | deposit: {it.deposit}
            </div>
            <div className="text-sm opacity-70">
              {it.category || "-"} | {it.location || "-"}
            </div>
            <div className="text-xs opacity-60 mt-1">id: {it.id} | owner_id: {it.owner_id}</div>
          </Link>
        ))}
      </div>
    </div>
  );
}
