"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { AdminGate } from "@/features/admin/AdminGate";
import { AdminNav } from "@/features/admin/AdminNav";
import { useAdminBookingEvents } from "@/features/admin/hooks";

export default function AdminBookingEventsPage() {
  const params = useParams<{ id: string }>();
  const id = Number(params.id);
  const eventsQ = useAdminBookingEvents(id, { limit: 100, offset: 0 });

  return (
    <AdminGate>
      <div className="p-6 space-y-4">
        <h1 className="text-2xl font-semibold">Админка: события брони #{id}</h1>
        <AdminNav />

        <div className="flex gap-2">
          <Link className="border rounded-lg px-3 py-2" href="/admin/bookings">
            ← К списку
          </Link>
          <Link className="border rounded-lg px-3 py-2" href={`/admin/bookings/${id}`}>
            Карточка брони
          </Link>
        </div>

        {eventsQ.isLoading && <div>Загрузка...</div>}
        {eventsQ.error && <div className="text-red-600">Ошибка: {(eventsQ.error as Error).message}</div>}

        <div className="space-y-3">
          {(eventsQ.data?.events ?? []).map((e) => (
            <div key={e.id} className="border rounded-xl p-4 space-y-2">
              <div className="font-medium">{e.action}</div>
              <div className="text-sm opacity-70">
                id: {e.id} • actor: {e.actor_user_id ?? "system"} • {new Date(e.created_at).toLocaleString("ru-RU")}
              </div>
              <div className="text-sm opacity-70">
                {e.from_status ?? "—"} → {e.to_status ?? "—"}
              </div>
              {e.meta != null && <pre className="text-xs overflow-auto">{JSON.stringify(e.meta, null, 2)}</pre>}
            </div>
          ))}
        </div>
      </div>
    </AdminGate>
  );
}

