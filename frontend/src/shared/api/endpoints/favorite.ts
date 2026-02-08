import { apiFetch } from "@/shared/api/client";
import type { FavoriteItem, Item } from "@/shared/api/types";

export const favoriteApi = {
  add: (itemId: number) => apiFetch<void>(`/api/items/${itemId}/favorite`, { method: "POST" }),
  remove: (itemId: number) => apiFetch<void>(`/api/items/${itemId}/favorite`, { method: "DELETE" }),
  isFavorite: (itemId: number) => apiFetch<{ is_favorite: boolean }>(`/api/items/${itemId}/favorite`, { method: "GET" }),
  my: (params?: { limit?: number; offset?: number }) => {
    const q = new URLSearchParams();
    if (params?.limit != null) q.set("limit", String(params.limit));
    if (params?.offset != null) q.set("offset", String(params.offset));
    const qs = q.toString();
    return apiFetch<FavoriteItem[]>(`/api/my/favorites${qs ? `?${qs}` : ""}`, { method: "GET" });
  },
};
