"use client";

import Link from "next/link";
import { useHydrateSession } from "@/shared/auth/useHydrateSession";
import { useMyBookings } from "@/features/booking/hooks";
import { useBookingActions } from "@/features/booking/hooks";

export default function MyBookingsPage() {
  useHydrateSession();
  const q = useMyBookings();
  const actions = useBookingActions();

  if (q.isLoading) return <div className="p-6">Загрузка...</div>;
  if (q.error) return <div className="p-6">Ошибка / 401 если не залогинен</div>;

  return (
    <div className="p-6 space-y-4">
      <h1 className="text-2xl font-semibold">Мои заявки</h1>

      <div className="space-y-3">
        {q.data?.map((b: any) => (
          <div key={b.id} className="border rounded-xl p-4 space-y-2">
            <div className="font-medium">
              Booking #{b.id} • item #{b.item_id} • {b.type} • {b.status}
            </div>

            {b.start_at && b.end_at && (
              <div className="text-sm opacity-70">
                {b.start_at} .. {b.end_at}
              </div>
            )}

            <div className="flex gap-2">
              <Link className="border rounded-lg px-3 py-2" href={`/bookings/${b.id}/events`}>
                События
              </Link>

              <button
                className="border rounded-lg px-3 py-2"
                disabled={actions.cancel.isPending}
                onClick={() => actions.cancel.mutate(b.id)}
              >
                Отменить
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
