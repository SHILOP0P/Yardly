"use client";

import { useState } from "react";
import type { ItemImage } from "@/shared/api/types";
import { resolveMediaUrl } from "@/shared/api/media";

export function ItemGallery({ images }: { images: ItemImage[] }) {
  const [idx, setIdx] = useState(0);
  const main = images[idx];

  if (!images.length) return <div className="border rounded-xl p-6 opacity-70">Нет фото</div>;

  return (
    <div className="space-y-3">
      <div className="border rounded-xl overflow-hidden">
        {/* пока просто <img> по url */}
        <img src={resolveMediaUrl(main.url)} alt="" className="w-full h-80 object-cover" />
      </div>

      <div className="flex gap-2 overflow-auto">
        {images.map((im, i) => (
          <button
            key={im.id}
            onClick={() => setIdx(i)}
            className={`border rounded-lg overflow-hidden w-20 h-20 flex-shrink-0 ${i === idx ? "opacity-100" : "opacity-60"}`}
            title={im.sort_order === 1 ? "Основная" : ""}
          >
            <img src={resolveMediaUrl(im.url)} alt="" className="w-full h-full object-cover" />
          </button>
        ))}
      </div>
    </div>
  );
}
