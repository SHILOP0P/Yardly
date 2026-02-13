"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { useItem, useItemImages } from "@/features/item/hooks";
import { ItemGallery } from "@/features/item/components/ItemGallery";
import { formatDealMode, formatItemStatus } from "@/features/item/presentation";
import { useIsFavorite, useToggleFavorite } from "@/features/favorite/hooks";
import { AvailabilityBookingCard } from "@/features/booking/AvailabilityBookingCard";
import { useBusyByItem } from "@/features/booking/hooks";

type BusyBooking = {
  id: number;
  status?: string;
  start?: string | null;
  end?: string | null;
};

export default function ItemPage() {
  const params = useParams<{ id: string }>();
  const id = Number(params.id);

  const itemQ = useItem(id);
  const imgsQ = useItemImages(id);
  const busyQ = useBusyByItem(id);
  const favQ = useIsFavorite(id);
  const { add, remove } = useToggleFavorite(id);

  if (itemQ.isLoading) return <div className="p-6">Loading...</div>;
  if (itemQ.error) return <div className="p-6">Failed to load product</div>;
  if (!itemQ.data) return <div className="p-6">Product not found</div>;

  const it = itemQ.data;
  const isFav = favQ.data?.is_favorite === true;
  const imageCount = imgsQ.data?.length ?? it.images?.length ?? 0;

  return (
    <div className="p-6 space-y-4">
      <div>
        <h1 className="text-2xl font-semibold">{it.title}</h1>
        <div className="text-sm opacity-70">
          {formatDealMode(it.mode)} ({it.mode}) | {formatItemStatus(it.status)} ({it.status})
        </div>

        <div className="flex flex-wrap gap-2 mt-3">
          <Link className="border rounded-lg px-3 py-2" href={`/users/${it.owner_id}/items`}>
            Owner products
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
            {isFav ? "In favorites" : "Add to favorites"}
          </button>
        </div>
      </div>

      <section className="border rounded-xl p-4 space-y-3">
        <div className="font-medium">Product model fields</div>
        <div className="grid grid-cols-[140px_1fr] gap-y-2 text-sm">
          <div className="opacity-70">id</div>
          <div>{it.id}</div>

          <div className="opacity-70">owner_id</div>
          <div>{it.owner_id}</div>

          <div className="opacity-70">title</div>
          <div>{it.title}</div>

          <div className="opacity-70">mode</div>
          <div>{it.mode}</div>

          <div className="opacity-70">status</div>
          <div>{it.status}</div>

          <div className="opacity-70">description</div>
          <div>{it.description || "-"}</div>

          <div className="opacity-70">price</div>
          <div>{it.price}</div>

          <div className="opacity-70">deposit</div>
          <div>{it.deposit}</div>

          <div className="opacity-70">location</div>
          <div>{it.location || "-"}</div>

          <div className="opacity-70">category</div>
          <div>{it.category || "-"}</div>

          <div className="opacity-70">images_count</div>
          <div>{imageCount}</div>
        </div>
      </section>

      {imgsQ.data ? <ItemGallery images={imgsQ.data} /> : <div className="opacity-70">Loading images...</div>}

      {(imgsQ.data?.length ?? 0) > 0 && (
        <section className="border rounded-xl p-4 space-y-2">
          <div className="font-medium">Image model fields</div>
          <div className="space-y-2">
            {imgsQ.data?.map((im) => (
              <div key={im.id} className="border rounded-lg p-2 text-sm">
                <div>id: {im.id}</div>
                <div>item_id: {im.item_id}</div>
                <div>url: <span className="break-all">{im.url}</span></div>
                <div>sort_order: {im.sort_order}</div>
                <div>created_at: {im.created_at}</div>
              </div>
            ))}
          </div>
        </section>
      )}

      <div className="border rounded-xl p-4 space-y-2">
        <div className="font-medium">Busy intervals (GET /api/items/{id}/bookings)</div>
        {busyQ.isLoading && <div className="text-sm opacity-70">Loading...</div>}
        {busyQ.error && <div className="text-sm text-red-600">Error: {(busyQ.error as Error).message}</div>}
        {(busyQ.data ?? []).length === 0 && !busyQ.isLoading && (
          <div className="text-sm opacity-70">No active busy intervals.</div>
        )}
        <div className="space-y-2">
          {(busyQ.data ?? []).map((b: BusyBooking) => (
            <div key={b.id} className="text-sm border rounded-lg p-2">
              booking #{b.id} | {b.status ?? "-"} | {b.start ?? "-"} {"->"} {b.end ?? "-"}
            </div>
          ))}
        </div>
      </div>

      <AvailabilityBookingCard itemId={id} mode={it.mode} />
    </div>
  );
}
