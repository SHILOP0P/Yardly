"use client";

import { useParams } from "next/navigation";
import { useBookingEvents } from "@/features/booking/hooks";
import { useHydrateSession } from "@/shared/auth/useHydrateSession";

export default function BookingEventsPage() {
  useHydrateSession();
  const params = useParams<{ id: string }>();
  const id = Number(params.id);

  const q = useBookingEvents(id);

  if (q.isLoading) return <div className="p-6">Загрузка...</div>;
  if (q.error) return <div className="p-6">Ошибка / 401</div>;

  return (
    <div className="p-6 space-y-4">
      <h1 className="text-2xl font-semibold">События брони #{id}</h1>

      <div className="space-y-2">
        {q.data?.map((e: any, idx: number) => (
          <div key={e.id ?? idx} className="border rounded-xl p-3">
            <div className="font-medium">{e.type ?? "event"}</div>
            <div className="text-sm opacity-70">{e.created_at ?? ""}</div>
            {e.meta && <pre className="text-xs opacity-80 mt-2">{JSON.stringify(e.meta, null, 2)}</pre>}
          </div>
        ))}
      </div>
    </div>
  );
}
