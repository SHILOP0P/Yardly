"use client";

import { useParams } from "next/navigation";
import { useItem, useItemImages } from "@/features/item/hooks";
import { ItemGallery } from "@/features/item/components/ItemGallery";
import Link from "next/link";


// ✅ добавь импорты favorites
import { useIsFavorite, useToggleFavorite } from "@/features/favorite/hooks";

import { AvailabilityBookingCard } from "@/features/booking/AvailabilityBookingCard";

export default function ItemPage() {
  const params = useParams<{ id: string }>();
  const id = Number(params.id);

  const itemQ = useItem(id);
  const imgsQ = useItemImages(id);

  // ✅ favorites hooks (после id)
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
        <div className="flex flex-wrap gap-2 mt-3">
            <Link className="border rounded-lg px-3 py-2" href={`/users/${it.owner_id}/items`}>
                Вещи владельца
            </Link>

            <Link href={`/items/${it.id}`}>Открыть</Link>
            <Link href={`/items/${it.id}/upcoming`}>Upcoming</Link>
            <Link href={`/my/items/bookings`}>Заявки</Link>


            <Link className="border rounded-lg px-3 py-2" href={`/items/${it.id}/upcoming`}>
                Upcoming по вещи
            </Link>
        </div>

        <div className="text-sm opacity-70">{it.mode} • {it.status}</div>

        {/* ✅ ВОТ СЮДА: кнопка избранного (после заголовка) */}
        <div className="mt-3">
          <button
            className="border rounded-lg px-3 py-2"
            disabled={add.isPending || remove.isPending}
            onClick={() => (isFav ? remove.mutate() : add.mutate())}
            title="Избранное"
          >
            {isFav ? "♥ В избранном" : "♡ В избранное"}
          </button>
        </div>
      </div>

      {imgsQ.data ? <ItemGallery images={imgsQ.data} /> : <div className="opacity-70">Фото...</div>}

      <AvailabilityBookingCard itemId={id} mode={it.mode} />

      {/* дальше сюда подключим availability + booking CTA */}
    </div>
  );
}
