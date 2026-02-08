import { apiFetch } from "@/shared/api/client";
import type { DealMode, Item, ItemImage } from "@/shared/api/types";

export const itemsApi = {
  list: (params?: { mode?: DealMode; limit?: number; offset?: number }) => {
    const q = new URLSearchParams();
    if (params?.mode) q.set("mode", params.mode);
    if (params?.limit != null) q.set("limit", String(params.limit));
    if (params?.offset != null) q.set("offset", String(params.offset));
    const qs = q.toString();
    return apiFetch<Item[]>(`/api/items${qs ? `?${qs}` : ""}`, { method: "GET" }, { auth: false });
  },

  getById: (id: number) => apiFetch<Item>(`/api/items/${id}`, { method: "GET" }, { auth: false }),

  images: {
    list: (itemId: number) =>
      apiFetch<ItemImage[]>(`/api/items/${itemId}/images`, { method: "GET" }, { auth: false }),

    add: (itemId: number, url: string) =>
      apiFetch<ItemImage>(`/api/items/${itemId}/images`, {
        method: "POST",
        body: JSON.stringify({ url }),
      }),

    delete: (itemId: number, imageId: number) =>
      apiFetch<void>(`/api/items/${itemId}/images/${imageId}`, { method: "DELETE" }),
  },

  byOwner: (ownerId: number, params?: { limit?: number; offset?: number }) => {
    const q = new URLSearchParams();
    if (params?.limit != null) q.set("limit", String(params.limit));
    if (params?.offset != null) q.set("offset", String(params.offset));
    const qs = q.toString();
    return apiFetch<any[]>(`/api/users/${ownerId}/items${qs ? `?${qs}` : ""}`, { method: "GET" }, { auth: false });
    },

  myItems: (params?: { limit?: number; offset?: number }) => {
    const q = new URLSearchParams();
    if (params?.limit != null) q.set("limit", String(params.limit));
    if (params?.offset != null) q.set("offset", String(params.offset));
    const qs = q.toString();
    return apiFetch<any[]>(`/api/my/items${qs ? `?${qs}` : ""}`, { method: "GET" });
  },

  create: (dto: { title: string; mode: string }) =>
  apiFetch(`/api/items`, { method: "POST", body: JSON.stringify(dto) }),


};

