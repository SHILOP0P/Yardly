"use client";

import { useMemo, useState } from "react";
import Link from "next/link";
import { AdminGate } from "@/features/admin/AdminGate";
import { AdminNav } from "@/features/admin/AdminNav";
import { useAdminBookings } from "@/features/admin/hooks";

export default function AdminBookingsPage() {
  const [status, setStatus] = useState("");
  const [type, setType] = useState("");
  const [itemID, setItemID] = useState("");
  const [userID, setUserID] = useState("");
  const [limit, setLimit] = useState(20);
  const [offset, setOffset] = useState(0);

  const params = useMemo(
    () => ({
      status: status || undefined,
      type: type || undefined,
      item_id: itemID ? Number(itemID) : undefined,
      user_id: userID ? Number(userID) : undefined,
      limit,
      offset,
    }),
    [status, type, itemID, userID, limit, offset]
  );

  const bookingsQ = useAdminBookings(params);
  const bookings = bookingsQ.data?.bookings ?? [];

  return (
    <AdminGate>
      <div className="p-6 space-y-4">
        <h1 className="text-2xl font-semibold">Админка: бронирования</h1>
        <AdminNav />

        <div className="border rounded-xl p-4 space-y-3">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-2">
            <input
              value={status}
              onChange={(e) => {
                setOffset(0);
                setStatus(e.target.value);
              }}
              className="border rounded-lg px-3 py-2"
              placeholder="status"
            />
            <input
              value={type}
              onChange={(e) => {
                setOffset(0);
                setType(e.target.value);
              }}
              className="border rounded-lg px-3 py-2"
              placeholder="type (rent/buy/give)"
            />
            <input
              value={itemID}
              onChange={(e) => {
                setOffset(0);
                setItemID(e.target.value);
              }}
              className="border rounded-lg px-3 py-2"
              placeholder="item_id"
            />
            <input
              value={userID}
              onChange={(e) => {
                setOffset(0);
                setUserID(e.target.value);
              }}
              className="border rounded-lg px-3 py-2"
              placeholder="user_id"
            />
          </div>

          <div className="flex gap-2">
            <select
              value={limit}
              className="border rounded-lg px-3 py-2"
              onChange={(e) => {
                setOffset(0);
                setLimit(Number(e.target.value));
              }}
            >
              {[10, 20, 50, 100].map((v) => (
                <option key={v} value={v}>
                  limit {v}
                </option>
              ))}
            </select>
            <button
              className="border rounded-lg px-3 py-2 disabled:opacity-50"
              disabled={offset <= 0 || bookingsQ.isFetching}
              onClick={() => setOffset((v) => Math.max(0, v - limit))}
            >
              Назад
            </button>
            <button
              className="border rounded-lg px-3 py-2 disabled:opacity-50"
              disabled={bookings.length < limit || bookingsQ.isFetching}
              onClick={() => setOffset((v) => v + limit)}
            >
              Вперед
            </button>
            <div className="text-sm opacity-70 flex items-center">offset: {offset}</div>
          </div>
        </div>

        {bookingsQ.isLoading && <div>Загрузка...</div>}
        {bookingsQ.error && <div className="text-red-600">Ошибка: {(bookingsQ.error as Error).message}</div>}

        <div className="space-y-3">
          {bookings.map((b) => (
            <div key={b.id} className="border rounded-xl p-4">
              <div className="flex flex-wrap justify-between gap-2">
                <div>
                  <div className="font-medium">#{b.id} • {b.status}</div>
                  <div className="text-sm opacity-70">
                    type: {b.type} • item: #{b.item_id} • requester: #{b.requester_id} • owner: #{b.owner_id}
                  </div>
                </div>
                <div className="flex gap-2">
                  <Link className="border rounded-lg px-3 py-2 text-sm" href={`/admin/bookings/${b.id}`}>
                    Детали
                  </Link>
                  <Link className="border rounded-lg px-3 py-2 text-sm" href={`/admin/bookings/${b.id}/events`}>
                    Events
                  </Link>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </AdminGate>
  );
}

