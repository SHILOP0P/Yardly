"use client";

import React from "react";
import Link from "next/link";
import { useHydrateSession } from "@/shared/auth/useHydrateSession";

// Если ты уже сделал хук — используй его:
import { useMyItemsBookings } from "@/features/booking/hooks";

// Если у тебя пока нет такого хука, можно временно заменить на bookingApi.listMyItems
// import { useQuery } from "@tanstack/react-query";
// import { bookingApi } from "@/shared/api/endpoints/booking";

const ALL_STATUSES = [
  "requested",
  "approved",
  "handover_pending",
  "in_use",
  "return_pending",
  "completed",
  "declined",
  "cancelled",
  "expired",
] as const;

const DEFAULT_ACTIVE_STATUSES = [
  "requested",
  "approved",
  "handover_pending",
  "in_use",
  "return_pending",
] as const;

type Status = (typeof ALL_STATUSES)[number] | string;

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

function ruDays(n: number) {
  const mod10 = n % 10;
  const mod100 = n % 100;
  if (mod10 === 1 && mod100 !== 11) return "день";
  if (mod10 >= 2 && mod10 <= 4 && !(mod100 >= 12 && mod100 <= 14)) return "дня";
  return "дней";
}

export default function MyItemsBookingsPage() {
  useHydrateSession();

  const [statuses, setStatuses] = React.useState<Status[]>([...DEFAULT_ACTIVE_STATUSES]);
  const [limit, setLimit] = React.useState<number>(20);
  const [offset, setOffset] = React.useState<number>(0);

  // Вариант 1: через твой hook (рекомендую)
  const q = useMyItemsBookings({ statuses, limit, offset });

  // Вариант 2: если нет hook — раскомментируй и используй напрямую:
  // const q = useQuery({
  //   queryKey: ["booking", "myItems", { statuses, limit, offset }],
  //   queryFn: () => bookingApi.listMyItems({ statuses, limit, offset } as any),
  // });

  const items = q.data?.items ?? [];
  const hasPrev = offset > 0;
  const hasNext = items.length === limit;

  const toggleStatus = (s: Status) => {
    setOffset(0);
    setStatuses((prev) => (prev.includes(s) ? prev.filter((x) => x !== s) : [...prev, s]));
  };

  const setActivePreset = () => {
    setOffset(0);
    setStatuses([...DEFAULT_ACTIVE_STATUSES]);
  };

  const setAllPreset = () => {
    setOffset(0);
    setStatuses([...ALL_STATUSES]);
  };

  const clearPreset = () => {
    setOffset(0);
    setStatuses([]);
  };

  if (q.isLoading) return <div className="p-6">Загрузка…</div>;

  if (q.error) {
    const msg = (q.error as any)?.message;
    return (
      <div className="p-6 text-red-600">
        Ошибка: {msg ? String(msg) : "ошибка загрузки"}
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      <h1 className="text-2xl font-semibold">Заявки на мои вещи</h1>

      {/* Фильтры + пагинация */}
      <div className="border rounded-2xl p-4 bg-white space-y-3">
        <div className="flex flex-wrap items-center gap-2">
          <div className="text-sm font-medium text-gray-700 mr-2">Статусы:</div>

          <button className="border rounded-lg px-3 py-1 text-sm hover:bg-gray-50" onClick={setActivePreset}>
            Актуальные
          </button>
          <button className="border rounded-lg px-3 py-1 text-sm hover:bg-gray-50" onClick={setAllPreset}>
            Все
          </button>
          <button className="border rounded-lg px-3 py-1 text-sm hover:bg-gray-50" onClick={clearPreset}>
            Сбросить
          </button>
        </div>

        <div className="flex flex-wrap gap-2">
          {ALL_STATUSES.map((s) => {
            const on = statuses.includes(s);
            return (
              <button
                key={s}
                className={`border rounded-lg px-3 py-1 text-sm transition ${
                  on ? "bg-gray-100" : "hover:bg-gray-50"
                }`}
                onClick={() => toggleStatus(s)}
              >
                {s}
              </button>
            );
          })}
        </div>

        <div className="flex flex-wrap items-center gap-3">
          <div className="text-sm text-gray-600">
            limit:
            <select
              className="ml-2 border rounded-lg px-2 py-1"
              value={limit}
              onChange={(e) => {
                setOffset(0);
                setLimit(Number(e.target.value));
              }}
            >
              {[10, 20, 50, 100].map((v) => (
                <option key={v} value={v}>
                  {v}
                </option>
              ))}
            </select>
          </div>

          <div className="text-sm text-gray-600">
            offset: <span className="font-medium">{offset}</span>
          </div>

          <div className="ml-auto flex gap-2">
            <button
              className="border rounded-lg px-4 py-2 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              disabled={!hasPrev || q.isFetching}
              onClick={() => setOffset((o) => Math.max(0, o - limit))}
            >
              ← Назад
            </button>
            <button
              className="border rounded-lg px-4 py-2 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              disabled={!hasNext || q.isFetching}
              onClick={() => setOffset((o) => o + limit)}
            >
              Вперёд →
            </button>
          </div>
        </div>
      </div>

      {items.length === 0 && <div className="text-gray-500">Заявок по выбранным статусам нет</div>}

      <div className="space-y-4">
        {items.map((b: any) => {
          const st = String(b.status);
          const typ = String(b.type);

          // у тебя поля start/end
          const start = b.start;
          const end = b.end;
          const d = daysBetweenInclusive(start, end);

          return (
            <div key={b.id} className="border rounded-2xl p-5 bg-white space-y-4">
              <div className="flex items-start justify-between gap-4">
                <div>
                  <div className="text-lg font-semibold">Заявка #{b.id}</div>

                  <div className="mt-1 flex flex-wrap gap-2">
                    <Badge>Вещь #{b.item_id}</Badge>
                    <Badge>{typ}</Badge>
                    <Badge>{st}</Badge>
                    <Badge>Requester #{b.requester_id}</Badge>
                  </div>
                </div>

                <Link href={`/bookings/${b.id}/events`} className="text-sm text-blue-600 hover:underline">
                  События →
                </Link>
              </div>

              {start && end && (
                <div className="text-sm text-gray-600">
                  Период:{" "}
                  <span className="font-medium">
                    {formatDT(start)} → {formatDT(end)}
                  </span>
                  {typeof d === "number" && (
                    <span className="ml-2 text-gray-500">
                      ({d} {ruDays(d)})
                    </span>
                  )}
                </div>
              )}

              {/* Owner-side действия можно добавить, когда скажешь какие у тебя эндпоинты:
                  approve/decline/handover-confirm/return-confirm и т.д.
              */}
              <div className="flex flex-wrap gap-2 pt-2 border-t">
                <Link className="border rounded-lg px-4 py-2 hover:bg-gray-50" href={`/items/${b.item_id}`}>
                  Открыть вещь
                </Link>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
