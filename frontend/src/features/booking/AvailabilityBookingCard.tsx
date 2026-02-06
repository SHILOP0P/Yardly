"use client";

import { useMemo, useState } from "react";
import { useQuery, useMutation } from "@tanstack/react-query";
import { availabilityApi } from "@/shared/api/endpoints/availability";
import { bookingApi, BookingType } from "@/shared/api/endpoints/booking";

function overlaps(aStart: string, aEnd: string, bStart: string, bEnd: string) {
  // строки YYYY-MM-DD сравниваются лексикографически корректно
  return !(aEnd < bStart || bEnd < aStart);
}

export function AvailabilityBookingCard({ itemId, mode }: { itemId: number; mode: string }) {
  // диапазон для календаря (на месяц вперёд)
  const today = new Date();
  const yyyy = today.getFullYear();
  const mm = String(today.getMonth() + 1).padStart(2, "0");
  const dd = String(today.getDate()).padStart(2, "0");
  const from = `${yyyy}-${mm}-${dd}`;

  const toDate = new Date(today);
  toDate.setDate(toDate.getDate() + 31);
  const to = `${toDate.getFullYear()}-${String(toDate.getMonth() + 1).padStart(2, "0")}-${String(toDate.getDate()).padStart(2, "0")}`;

  const availQ = useQuery({
    queryKey: ["availability", itemId, from, to],
    queryFn: () => availabilityApi.get(itemId, from, to),
    enabled: itemId > 0,
  });

  const [start, setStart] = useState(from);
  const [end, setEnd] = useState(from);

  const busy = availQ.data?.busy ?? [];

  const isRangeValid = useMemo(() => {
    if (!start || !end) return false;
    if (end < start) return false;
    for (const b of busy) {
      if (overlaps(start, end, b.start, b.end)) return false;
    }
    return true;
  }, [start, end, busy]);

  const createM = useMutation({
    mutationFn: async () => {
      if (mode === "rent" || mode === "sale_rent") {
        return bookingApi.create(itemId, { type: "rent", start_at: start, end_at: end });
      }
      if (mode === "sale") return bookingApi.create(itemId, { type: "buy" });
      return bookingApi.create(itemId, { type: "give" });
    },
  });

  return (
    <div className="border rounded-xl p-4 space-y-3">
      <div className="font-medium">Бронирование</div>

      {availQ.isLoading && <div className="opacity-70">Загрузка доступности...</div>}
      {availQ.error && <div className="text-sm text-red-500">Не удалось загрузить доступность</div>}

      {(mode === "rent" || mode === "sale_rent") && (
        <>
          <div className="text-sm opacity-70">
            Выбирай даты (формат YYYY-MM-DD). Занятые дни запрещены к выбору.
          </div>

          <div className="flex gap-2">
            <div className="flex-1">
              <div className="text-xs opacity-70 mb-1">Start</div>
              <input className="w-full border rounded-lg p-2" type="date" value={start} onChange={(e) => setStart(e.target.value)} />
            </div>
            <div className="flex-1">
              <div className="text-xs opacity-70 mb-1">End</div>
              <input className="w-full border rounded-lg p-2" type="date" value={end} onChange={(e) => setEnd(e.target.value)} />
            </div>
          </div>

          {!isRangeValid && (
            <div className="text-sm text-red-500">
              Диапазон некорректен или пересекается с занятыми датами.
            </div>
          )}

          {busy.length > 0 && (
            <div className="text-xs opacity-70">
              Занято: {busy.map((b) => `${b.start}..${b.end}`).join(", ")}
            </div>
          )}
        </>
      )}

      {(mode === "sale" || mode === "free") && (
        <div className="text-sm opacity-70">
          Даты не нужны. Нажми кнопку — создастся заявка {mode === "sale" ? "на покупку" : "на получение"}.
        </div>
      )}

      <button
        className="border rounded-lg px-3 py-2"
        disabled={createM.isPending || ((mode === "rent" || mode === "sale_rent") && !isRangeValid)}
        onClick={() => createM.mutate()}
      >
        {createM.isPending ? "..." : "Создать заявку"}
      </button>

      {createM.isError && <div className="text-sm text-red-500">{(createM.error as any)?.message ?? "Ошибка"}</div>}
      {createM.isSuccess && <div className="text-sm text-green-600">Заявка создана</div>}
    </div>
  );
}
