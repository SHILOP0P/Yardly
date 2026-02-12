import { apiFetch } from "@/shared/api/client";

export type BookingType = "rent" | "buy" | "give";

// (опционально) если у тебя есть union по статусам — подставь его сюда
export type BookingStatus = string;

export type BookingListParams = {
  statuses?: BookingStatus[]; // мультистатусность
  limit?: number;
  offset?: number;
};

export type BookingListResp<T = any> = {
  items: T[];
  limit: number;
  offset: number;
};

function buildQuery(params?: BookingListParams) {
  const qs = new URLSearchParams();
  const statuses = params?.statuses ?? [];
  for (const s of statuses) qs.append("status", String(s)); // status=...&status=...

  if (typeof params?.limit === "number") qs.set("limit", String(params.limit));
  if (typeof params?.offset === "number") qs.set("offset", String(params.offset));

  const str = qs.toString();
  return str ? `?${str}` : "";
}

export const bookingApi = {
  create: (itemId: number, payload: { type: BookingType; start_at?: string; end_at?: string }) =>
    apiFetch<any>(`/api/items/${itemId}/bookings`, { method: "POST", body: JSON.stringify(payload) }),

  listBusyForItem: (itemId: number) =>
    apiFetch<any[]>(`/api/items/${itemId}/bookings`, { method: "GET" }, { auth: false }),

  listMy: (params?: BookingListParams) =>
    apiFetch<BookingListResp>(`/api/my/bookings${buildQuery(params)}`, { method: "GET" }),

  listMyItems: (params?: BookingListParams) =>
    apiFetch<BookingListResp>(`/api/my/items/bookings${buildQuery(params)}`, { method: "GET" }),

  listMyItemsBookingRequests: (params?: { types?: BookingType[]; limit?: number; offset?: number }) => {
    const qs = new URLSearchParams();
    for (const t of params?.types ?? []) qs.append("type", t);
    if (typeof params?.limit === "number") qs.set("limit", String(params.limit));
    if (typeof params?.offset === "number") qs.set("offset", String(params.offset));
    const s = qs.toString();
    return apiFetch<BookingListResp>(`/api/my/items/booking-requests${s ? `?${s}` : ""}`, { method: "GET" });
  },

  approve: (id: number) => apiFetch<void>(`/api/bookings/${id}/approve`, { method: "POST" }),
  handover: (id: number) => apiFetch<void>(`/api/bookings/${id}/handover`, { method: "POST" }),
  return: (id: number) => apiFetch<void>(`/api/bookings/${id}/return`, { method: "POST" }),
  cancel: (id: number) => apiFetch<void>(`/api/bookings/${id}/cancel`, { method: "POST" }),

  events: (id: number) => apiFetch<any[]>(`/api/bookings/${id}/events`, { method: "GET" }),

  upcomingByItem: (itemId: number) =>
    apiFetch<any>(`/api/items/${itemId}/bookings/upcoming`, { method: "GET" }, { auth: false }),
};
