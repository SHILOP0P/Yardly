"use client";

import Link from "next/link";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useHydrateSession } from "@/shared/auth/useHydrateSession";
import { useMyItemsBookings } from "@/features/booking/hooks"; // если у тебя так называется
import { bookingApi } from "@/shared/api/endpoints/booking";

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
  // если бронь на сутки обычно end = start + N дней (как у тебя), то считаем дни как разницу по датам
  const days = Math.ceil(ms / (24 * 60 * 60 * 1000));
  return days;
}

export default function MyItemsBookingsPage() {
  useHydrateSession();

  const q = useMyItemsBookings();
  const qc = useQueryClient();

  const invalidate = () => {
    qc.invalidateQueries({ queryKey: ["my", "items", "bookings"] }); // если у тебя другой ключ — поправь
    // иногда проще так:
    qc.invalidateQueries();
  };

  const approveMut = useMutation({
    mutationFn: (id: number) => bookingApi.approve(id),
    onSuccess: invalidate,
  });

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
  if (q.error) return <div className="p-6">Ошибка / 401 если не залогинен</div>;

  return (
    <div className="p-6 space-y-6">
      <h1 className="text-2xl font-semibold">Заявки на мои вещи</h1>

      {q.data?.length === 0 && (
        <div className="text-gray-500">Заявок пока нет</div>
      )}

      <div className="space-y-4">
        {q.data?.map((b: any) => {
          const st = b.status as string;
          const typ = b.type as string;

          // Owner-side логика видимости кнопок
          const canApprove = st === "requested";

          // Подтверждение передачи (handover) для owner:
          // rent: approved или handover_pending
          // buy/give: approved или handover_pending (у тебя общий эндпоинт handover, судя по UI)
          const canHandover =
            st === "approved" || st === "handover_pending";

          // Подтверждение возврата (только rent)
          const canReturn =
            typ === "rent" && (st === "in_use" || st === "return_pending");

          // Cancel на owner-side обычно НЕ должен быть (это спорно),
          // но раз у тебя есть кнопка, показываем только когда ещё не завершено.
          const canCancel =
            st === "requested" || st === "approved" || st === "handover_pending";

          const rentDays = daysBetweenInclusive(b.start, b.end);
          const hasPeriod = Boolean(b.start && b.end);

          return (
            <div key={b.id} className="border rounded-2xl p-5 bg-white space-y-4">
              {/* Header */}
              <div className="flex items-start justify-between gap-4">
                <div>
                  <div className="text-lg font-semibold">Заявка #{b.id}</div>

                  <div className="mt-1 flex flex-wrap gap-2">
                    <Badge>Вещь #{b.item_id}</Badge>
                    <Badge>Запросил: #{b.requester_id}</Badge>
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

              {/* Period */}
              {hasPeriod && (
                <div className="text-sm text-gray-600">
                  Период аренды:{" "}
                  <span className="font-medium">
                    {formatDT(b.start)} → {formatDT(b.end)}
                  </span>
                  {typeof rentDays === "number" && (
                    <span className="ml-2 text-gray-500">
                      ({rentDays} {rentDays === 1 ? "день" : rentDays < 5 ? "дня" : "дней"})
                    </span>
                  )}
                </div>
              )}

              {/* Actions */}
              <div className="flex flex-wrap gap-2 pt-2 border-t">
                <Link
                  className="border rounded-lg px-4 py-2 cursor-pointer hover:bg-gray-50 transition"
                  href={`/items/${b.item_id}`}
                >
                  Открыть вещь
                </Link>

                {canApprove && (
                  <button
                    className="border rounded-lg px-4 py-2 cursor-pointer hover:bg-gray-50 transition disabled:opacity-50 disabled:cursor-not-allowed"
                    disabled={approveMut.isPending}
                    onClick={() => approveMut.mutate(b.id)}
                  >
                    Одобрить
                  </button>
                )}

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

              {/* Errors */}
              {(approveMut.error || handoverMut.error || returnMut.error || cancelMut.error) ? (
                <div className="text-sm text-red-600">
                  {String(
                    (approveMut.error as any)?.message ||
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
