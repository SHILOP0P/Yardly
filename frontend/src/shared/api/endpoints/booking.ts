import { apiFetch } from "@/shared/api/client";

export type BookingType = "rent" | "buy" | "give";

export const bookingApi = {
  create: (itemId: number, payload: { type: BookingType; start_at?: string; end_at?: string }) =>
    apiFetch<any>(`/api/items/${itemId}/bookings`, { method: "POST", body: JSON.stringify(payload) }),

  listMy: () => apiFetch<any[]>(`/api/my/bookings`, { method: "GET" }),
  listMyItems: () => apiFetch<any[]>(`/api/my/items/bookings`, { method: "GET" }),

  approve: (id: number) => apiFetch<void>(`/api/bookings/${id}/approve`, { method: "POST" }),
  handover: (id: number) => apiFetch<void>(`/api/bookings/${id}/handover`, { method: "POST" }),
  return: (id: number) => apiFetch<void>(`/api/bookings/${id}/return`, { method: "POST" }),
  cancel: (id: number) => apiFetch<void>(`/api/bookings/${id}/cancel`, { method: "POST" }),

  events: (id: number) => apiFetch<any[]>(`/api/bookings/${id}/events`, { method: "GET" }),

  upcomingByItem: (itemId: number) =>
    apiFetch<any>(`/api/items/${itemId}/bookings/upcoming`, { method: "GET" }, { auth: false }),
};
