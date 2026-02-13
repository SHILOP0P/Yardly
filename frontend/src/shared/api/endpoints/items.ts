import { apiFetch } from "@/shared/api/client";
import type { DealMode, Item, ItemImage } from "@/shared/api/types";

export type ItemListParams = {
  mode?: DealMode;
  category?: string;
  location?: string;
  min_price?: number;
  max_price?: number;
  limit?: number;
  offset?: number;
};

export type CreateItemDto = {
  title: string;
  mode: DealMode;
  description?: string;
  price?: number;
  deposit?: number;
  location?: string;
  category?: string;
};

export const itemsApi = {
  list: (params?: ItemListParams) => {
    const q = new URLSearchParams();
    if (params?.mode) q.set("mode", params.mode);
    if (params?.category) q.set("category", params.category);
    if (params?.location) q.set("location", params.location);
    if (params?.min_price != null) q.set("min_price", String(params.min_price));
    if (params?.max_price != null) q.set("max_price", String(params.max_price));
    if (params?.limit != null) q.set("limit", String(params.limit));
    if (params?.offset != null) q.set("offset", String(params.offset));
    const qs = q.toString();
    return apiFetch<Item[]>(`/api/items${qs ? `?${qs}` : ""}`, { method: "GET" }, { auth: false });
  },

  getById: (id: number) => apiFetch<Item>(`/api/items/${id}`, { method: "GET" }, { auth: false }),

  images: {
    list: (itemId: number) =>
      apiFetch<ItemImage[]>(`/api/items/${itemId}/images`, { method: "GET" }, { auth: false }),

    add: (itemId: number, file: File) => {
      const body = new FormData();
      body.append("file", file);
      return apiFetch<ItemImage>(`/api/items/${itemId}/images`, { method: "POST", body });
    },

    delete: (itemId: number, imageId: number) =>
      apiFetch<void>(`/api/items/${itemId}/images/${imageId}`, { method: "DELETE" }),
  },

  byOwner: (ownerId: number, params?: ItemListParams) => {
    const q = new URLSearchParams();
    if (params?.mode) q.set("mode", params.mode);
    if (params?.category) q.set("category", params.category);
    if (params?.location) q.set("location", params.location);
    if (params?.min_price != null) q.set("min_price", String(params.min_price));
    if (params?.max_price != null) q.set("max_price", String(params.max_price));
    if (params?.limit != null) q.set("limit", String(params.limit));
    if (params?.offset != null) q.set("offset", String(params.offset));
    const qs = q.toString();
    return apiFetch<Item[]>(`/api/users/${ownerId}/items${qs ? `?${qs}` : ""}`, { method: "GET" }, { auth: false });
    },

  myItems: (params?: ItemListParams) => {
    const q = new URLSearchParams();
    if (params?.mode) q.set("mode", params.mode);
    if (params?.category) q.set("category", params.category);
    if (params?.location) q.set("location", params.location);
    if (params?.min_price != null) q.set("min_price", String(params.min_price));
    if (params?.max_price != null) q.set("max_price", String(params.max_price));
    if (params?.limit != null) q.set("limit", String(params.limit));
    if (params?.offset != null) q.set("offset", String(params.offset));
    const qs = q.toString();
    return apiFetch<Item[]>(`/api/my/items${qs ? `?${qs}` : ""}`, { method: "GET" });
  },

  create: (dto: CreateItemDto) =>
    apiFetch<Item>(`/api/items`, { method: "POST", body: JSON.stringify(dto) }),


};

