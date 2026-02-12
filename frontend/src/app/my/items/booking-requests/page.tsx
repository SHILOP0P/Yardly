"use client";

import React from "react";
import Link from "next/link";
import { useHydrateSession } from "@/shared/auth/useHydrateSession";
import { useMyItemsBookingRequests } from "@/features/booking/hooks";

const ALL_TYPES = ["rent", "buy", "give"] as const;
type BookingType = (typeof ALL_TYPES)[number];

export default function MyItemsBookingRequestsPage() {
  useHydrateSession();

  const [types, setTypes] = React.useState<BookingType[]>(["rent", "buy", "give"]);
  const [limit, setLimit] = React.useState(20);
  const [offset, setOffset] = React.useState(0);

  const q = useMyItemsBookingRequests({ types, limit, offset });
  const items = q.data?.items ?? [];

  const toggleType = (t: BookingType) => {
    setOffset(0);
    setTypes((prev) => (prev.includes(t) ? prev.filter((x) => x !== t) : [...prev, t]));
  };

  return (
    <div className="p-6 space-y-5">
      <h1 className="text-2xl font-semibold">Запросы на мои вещи (by type)</h1>

      <div className="border rounded-xl p-4 space-y-3">
        <div className="text-sm opacity-70">Фильтр `type` для `GET /api/my/items/booking-requests`</div>
        <div className="flex flex-wrap gap-2">
          {ALL_TYPES.map((t) => {
            const on = types.includes(t);
            return (
              <button
                key={t}
                className={`border rounded-lg px-3 py-2 text-sm ${on ? "bg-muted" : ""}`}
                onClick={() => toggleType(t)}
              >
                {t}
              </button>
            );
          })}
        </div>

        <div className="flex gap-2">
          <select
            value={limit}
            onChange={(e) => {
              setOffset(0);
              setLimit(Number(e.target.value));
            }}
            className="border rounded-lg px-3 py-2"
          >
            {[10, 20, 50, 100].map((v) => (
              <option key={v} value={v}>
                limit {v}
              </option>
            ))}
          </select>
          <button
            className="border rounded-lg px-3 py-2 disabled:opacity-50"
            disabled={offset <= 0 || q.isFetching}
            onClick={() => setOffset((v) => Math.max(0, v - limit))}
          >
            Назад
          </button>
          <button
            className="border rounded-lg px-3 py-2 disabled:opacity-50"
            disabled={items.length < limit || q.isFetching}
            onClick={() => setOffset((v) => v + limit)}
          >
            Вперед
          </button>
        </div>
      </div>

      {q.isLoading && <div>Загрузка...</div>}
      {q.error && <div className="text-red-600">Ошибка: {(q.error as Error).message}</div>}

      <div className="space-y-3">
        {items.map((b: any) => (
          <div key={b.id} className="border rounded-xl p-4">
            <div className="font-medium">
              #{b.id} • {b.type} • {b.status}
            </div>
            <div className="text-sm opacity-70">
              item #{b.item_id} • requester #{b.requester_id}
            </div>
            <div className="mt-2 flex gap-2">
              <Link href={`/items/${b.item_id}`} className="border rounded-lg px-3 py-2 text-sm">
                Вещь
              </Link>
              <Link href={`/bookings/${b.id}/events`} className="border rounded-lg px-3 py-2 text-sm">
                Events
              </Link>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

