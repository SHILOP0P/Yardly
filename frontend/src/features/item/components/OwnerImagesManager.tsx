"use client";

import { useState } from "react";
import { useAddItemImage, useDeleteItemImage, useItemImages } from "@/features/item/hooks";
import { resolveMediaUrl } from "@/shared/api/media";

export function OwnerImagesManager({ itemId }: { itemId: number }) {
  const imgsQ = useItemImages(itemId);
  const addM = useAddItemImage(itemId);
  const delM = useDeleteItemImage(itemId);

  const [file, setFile] = useState<File | null>(null);
  const images = imgsQ.data ?? [];

  return (
    <div className="border rounded-xl p-4 space-y-3">
      <div className="font-medium">Product photos (first photo is main)</div>

      <div className="flex gap-2">
        <input
          className="flex-1 border rounded-lg p-2"
          type="file"
          accept=".png,.jpg,.jpeg,image/png,image/jpeg"
          onChange={(e) => setFile(e.target.files?.[0] ?? null)}
        />
        <button
          className="border rounded-lg px-3"
          disabled={addM.isPending || !file}
          onClick={async () => {
            if (!file) return;
            await addM.mutateAsync(file);
            setFile(null);
          }}
        >
          Add
        </button>
      </div>

      {imgsQ.isLoading && <div className="opacity-70">Loading...</div>}
      {imgsQ.error && <div className="text-red-500 text-sm">Failed to load images</div>}

      <div className="space-y-2">
        {images.map((im) => (
          <div key={im.id} className="flex items-center gap-3 border rounded-lg p-2">
            <img src={resolveMediaUrl(im.url)} alt="" className="w-16 h-16 object-cover rounded-md" />
            <div className="flex-1">
              <div className="text-sm">
                #{im.sort_order} {im.sort_order === 1 ? "• main" : ""}
              </div>
              <div className="text-xs opacity-70 break-all">{im.url}</div>
            </div>
            <button className="border rounded-lg px-3 py-1" disabled={delM.isPending} onClick={() => delM.mutate(im.id)}>
              Delete
            </button>
          </div>
        ))}
        {!images.length && !imgsQ.isLoading && <div className="opacity-70">No photos yet</div>}
      </div>
    </div>
  );
}
