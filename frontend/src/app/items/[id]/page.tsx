"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { useItem, useItemImages } from "@/features/item/hooks";
import { ItemGallery } from "@/features/item/components/ItemGallery";
import { useIsFavorite, useToggleFavorite } from "@/features/favorite/hooks";
import { AvailabilityBookingCard } from "@/features/booking/AvailabilityBookingCard";
import { useBusyByItem } from "@/features/booking/hooks";

export default function ItemPage() {
  const params = useParams<{ id: string }>();
  const id = Number(params.id);

  const itemQ = useItem(id);
  const imgsQ = useItemImages(id);
  const busyQ = useBusyByItem(id);
  const favQ = useIsFavorite(id);
  const { add, remove } = useToggleFavorite(id);

  if (itemQ.isLoading) return <div className="p-6">Загрузка...</div>;
  if (itemQ.error) return <div className="p-6">Ошибка загрузки товара</div>;
  const it = itemQ.data!;

  const isFav = favQ.data?.is_favorite === true;

  return (
    <div className="p-6 space-y-4">
      <div>
        <h1 className="text-2xl font-semibold">{it.title}</h1>
        <div className="text-sm opacity-70">
          {it.mode} • {it.status}
        </div>

        <div className="flex flex-wrap gap-2 mt-3">
          <Link className="border rounded-lg px-3 py-2" href={`/users/${it.owner_id}/items`}>
            Вещи владельца
          </Link>
          <Link className="border rounded-lg px-3 py-2" href={`/items/${it.id}/upcoming`}>
            Upcoming
          </Link>
        </div>

        <div className="mt-3">
          <button
            className="border rounded-lg px-3 py-2"
            disabled={add.isPending || remove.isPending}
            onClick={() => (isFav ? remove.mutate() : add.mutate())}
          >
            {isFav ? "В избранном" : "В избранное"}
          </button>
        </div>
      </div>

      {imgsQ.data ? <ItemGallery images={imgsQ.data} /> : <div className="opacity-70">Фото...</div>}

      <div className="border rounded-xl p-4 space-y-2">
        <div className="font-medium">Занятые интервалы (GET /api/items/{id}/bookings)</div>
        {busyQ.isLoading && <div className="text-sm opacity-70">Загрузка...</div>}
        {busyQ.error && <div className="text-sm text-red-600">Ошибка: {(busyQ.error as Error).message}</div>}
        {(busyQ.data ?? []).length === 0 && !busyQ.isLoading && (
          <div className="text-sm opacity-70">Нет активных занятых интервалов.</div>
        )}
        <div className="space-y-2">
          {(busyQ.data ?? []).map((b: any) => (
            <div key={b.id} className="text-sm border rounded-lg p-2">
              booking #{b.id} • {b.status} • {b.start ?? "—"} → {b.end ?? "—"}
            </div>
          ))}
        </div>
      </div>

      <AvailabilityBookingCard itemId={id} mode={it.mode} />
    </div>
  );
}
