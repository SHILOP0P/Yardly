"use client";

import { useState } from "react";
import { useAddItemImage, useDeleteItemImage, useItemImages } from "@/features/item/hooks";

export function OwnerImagesManager({ itemId }: { itemId: number }) {
  const imgsQ = useItemImages(itemId);
  const addM = useAddItemImage(itemId);
  const delM = useDeleteItemImage(itemId);

  const [url, setUrl] = useState("");

  const images = imgsQ.data ?? [];

  return (
    <div className="border rounded-xl p-4 space-y-3">
      <div className="font-medium">Фото товара (порядок как добавлял, 1-е — основное)</div>

      <div className="flex gap-2">
        <input
          className="flex-1 border rounded-lg p-2"
          placeholder="https://... (пока добавляем URL)"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
        />
        <button
          className="border rounded-lg px-3"
          disabled={addM.isPending || !url.trim()}
          onClick={async () => {
            await addM.mutateAsync(url.trim());
            setUrl("");
          }}
        >
          Добавить
        </button>
      </div>

      {imgsQ.isLoading && <div className="opacity-70">Загрузка...</div>}
      {imgsQ.error && <div className="text-red-500 text-sm">Ошибка загрузки фото</div>}

      <div className="space-y-2">
        {images.map((im) => (
          <div key={im.id} className="flex items-center gap-3 border rounded-lg p-2">
            <img src={im.url} alt="" className="w-16 h-16 object-cover rounded-md" />
            <div className="flex-1">
              <div className="text-sm">
                #{im.sort_order} {im.sort_order === 1 ? "• основная" : ""}
              </div>
              <div className="text-xs opacity-70 break-all">{im.url}</div>
            </div>
            <button
              className="border rounded-lg px-3 py-1"
              disabled={delM.isPending}
              onClick={() => delM.mutate(im.id)}
            >
              Удалить
            </button>
          </div>
        ))}
        {!images.length && !imgsQ.isLoading && <div className="opacity-70">Фото пока нет</div>}
      </div>

      <div className="text-xs opacity-70">
        Если удалить основную (№1), backend сдвинет порядок и новая №1 станет основной автоматически.
      </div>
    </div>
  );
}
