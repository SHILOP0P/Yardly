"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { itemsApi } from "@/shared/api/endpoints/items";
import { useHydrateSession } from "@/shared/auth/useHydrateSession";
import { formatDealMode, formatItemStatus } from "@/features/item/presentation";

export default function MyItemsPage() {
  useHydrateSession();

  const q = useQuery({
    queryKey: ["my", "items"],
    queryFn: () => itemsApi.myItems(),
  });

  if (q.isLoading) return <div className="p-6">Loading...</div>;
  if (q.error) return <div className="p-6">Error / 401</div>;

  return (
    <div className="p-6 space-y-4">
      <h1 className="text-2xl font-semibold">My products</h1>
      <Link className="border rounded-lg px-3 py-2 inline-flex" href="/my/items/new">
        + Create product
      </Link>

      <div className="space-y-3">
        {q.data?.map((it) => (
          <div key={it.id} className="border rounded-xl p-4 space-y-2">
            <div className="font-medium">{it.title}</div>
            <div className="text-sm opacity-70">{formatDealMode(it.mode)} | {formatItemStatus(it.status)}</div>
            {it.description ? <div className="text-sm">{it.description}</div> : null}
            <div className="text-sm opacity-80">price: {it.price} | deposit: {it.deposit}</div>
            <div className="text-sm opacity-70">{it.category || "-"} | {it.location || "-"}</div>
            <div className="text-xs opacity-60">id: {it.id} | owner_id: {it.owner_id}</div>

            <div className="flex flex-wrap gap-2">
              <Link className="border rounded-lg px-3 py-2" href={`/items/${it.id}`}>
                Open
              </Link>
              <Link className="border rounded-lg px-3 py-2" href={`/items/${it.id}/upcoming`}>
                Upcoming
              </Link>
              <Link className="border rounded-lg px-3 py-2" href="/my/items/booking-requests">
                Booking requests
              </Link>
              <Link className="border rounded-lg px-3 py-2" href="/my/items/bookings">
                All bookings
              </Link>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
