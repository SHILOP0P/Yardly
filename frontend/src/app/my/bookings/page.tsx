"use client";

import Link from "next/link";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { bookingApi } from "@/shared/api/endpoints/booking";
import { useHydrateSession } from "@/shared/auth/useHydrateSession";

function Badge({ children }: { children: React.ReactNode }) {
  return (
    <span className="px-2 py-1 text-xs rounded-md bg-gray-100 text-gray-700">
      {children}
    </span>
  );
}

function formatDT(v?: string) {
  if (!v) return null;
  const d = new Date(v);
  if (Number.isNaN(d.getTime())) return v;
  return d.toLocaleString("ru-RU");
}

function daysBetweenInclusive(startISO?: string, endISO?: string) {
  if (!startISO || !endISO) return null;
  const s = new Date(startISO);
  const e = new Date(endISO);
  if (Number.isNaN(s.getTime()) || Number.isNaN(e.getTime())) return null;
  const ms = e.getTime() - s.getTime();
  if (ms < 0) return null;
  return Math.ceil(ms / (24 * 60 * 60 * 1000));
}

export default function MyBookingsPage() {
  useHydrateSession();

  const q = useQuery({
    queryKey: ["my", "bookings"],
    queryFn: bookingApi.listMy,
  });

  const qc = useQueryClient();
  const invalidate = () => qc.invalidateQueries({ queryKey: ["my", "bookings"] });

  const handoverMut = useMutation({
    mutationFn: (id: number) => bookingApi.handover(id),
    onSuccess: invalidate,
  });

  const returnMut = useMutation({
    mutationFn: (id: number) => bookingApi.return(id),
    onSuccess: invalidate,
  });

  const cancelMut = useMutation({
    mutationFn: (id: number) => bookingApi.cancel(id),
    onSuccess: invalidate,
  });

  if (q.isLoading) return <div className="p-6">Загрузка…</div>;
  if (q.error) return <div className="p-6">Ошибка / 401</div>;

  return (
    <div className="p-6 space-y-6">
      <h1 className="text-2xl font-semibold">Мои заявки</h1>

      {q.data?.length === 0 && (
        <div className="text-gray-500">Заявок пока нет</div>
      )}

      <div className="space-y-4">
        {q.data?.map((b: any) => {
          const st = b.status as string;
          const typ = b.type as string;

          // requester-side
          const canHandover = st === "approved" || st === "handover_pending";
          const canReturn = st === "in_use" || st === "return_pending";
          const canCancel = st === "requested" || st === "approved" || st === "handover_pending";

          const start = b.start ?? b.start_at;
          const end = b.end ?? b.end_at;

          const rentDays = daysBetweenInclusive(start, end);
          const hasPeriod = Boolean(start && end);

          return (
            <div key={b.id} className="border rounded-2xl p-5 bg-white space-y-4">
              <div className="flex items-start justify-between gap-4">
                <div>
                  <div className="text-lg font-semibold">Заявка #{b.id}</div>
                  <div className="mt-1 flex flex-wrap gap-2">
                    <Badge>Вещь #{b.item_id}</Badge>
                    <Badge>{typ}</Badge>
                    <Badge>{st}</Badge>
                  </div>
                </div>

                <Link
                  href={`/bookings/${b.id}/events`}
                  className="text-sm text-blue-600 hover:underline"
                >
                  События →
                </Link>
              </div>

              {hasPeriod && (
                <div className="text-sm text-gray-600">
                  Период:{" "}
                  <span className="font-medium">
                    {formatDT(start)} → {formatDT(end)}
                  </span>
                  {typeof rentDays === "number" && (
                    <span className="ml-2 text-gray-500">
                      ({rentDays} {rentDays === 1 ? "день" : rentDays < 5 ? "дня" : "дней"})
                    </span>
                  )}
                </div>
              )}

              <div className="flex flex-wrap gap-2 pt-2 border-t">
                <Link
                  className="border rounded-lg px-4 py-2 cursor-pointer hover:bg-gray-50 transition"
                  href={`/items/${b.item_id}`}
                >
                  Открыть вещь
                </Link>

                {canHandover && (
                  <button
                    className="border rounded-lg px-4 py-2 cursor-pointer hover:bg-gray-50 transition disabled:opacity-50 disabled:cursor-not-allowed"
                    disabled={handoverMut.isPending}
                    onClick={() => handoverMut.mutate(b.id)}
                  >
                    Подтвердить передачу
                  </button>
                )}

                {canReturn && (
                  <button
                    className="border rounded-lg px-4 py-2 cursor-pointer hover:bg-gray-50 transition disabled:opacity-50 disabled:cursor-not-allowed"
                    disabled={returnMut.isPending}
                    onClick={() => returnMut.mutate(b.id)}
                  >
                    Подтвердить возврат
                  </button>
                )}

                {canCancel && (
                  <button
                    className="border rounded-lg px-4 py-2 cursor-pointer text-red-600 hover:bg-red-50 transition disabled:opacity-50 disabled:cursor-not-allowed"
                    disabled={cancelMut.isPending}
                    onClick={() => cancelMut.mutate(b.id)}
                  >
                    Отменить
                  </button>
                )}
              </div>

              {(handoverMut.error || returnMut.error || cancelMut.error) ? (
                <div className="text-sm text-red-600">
                  {String(
                    (handoverMut.error as any)?.message ||
                      (returnMut.error as any)?.message ||
                      (cancelMut.error as any)?.message
                  )}
                </div>
              ) : null}
            </div>
          );
        })}
      </div>
    </div>
  );
}
