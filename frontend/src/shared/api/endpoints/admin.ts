import { apiFetch } from "@/shared/api/client";
import type {
  AdminBooking,
  AdminBookingEvent,
  AdminEvent,
  AdminItem,
  AdminListResp,
  AdminUser,
  UserRole,
} from "@/shared/api/types";

export type AdminListUsersParams = {
  q?: string;
  limit?: number;
  offset?: number;
};

export type AdminPatchUserPayload = {
  role?: UserRole;
  ban?: boolean;
  ban_reason?: string;
  ban_expires_at?: string | null;
};

export type AdminListBookingsParams = {
  status?: string;
  type?: string;
  item_id?: number;
  user_id?: number;
  limit?: number;
  offset?: number;
};

export type AdminListItemsParams = {
  q?: string;
  status?: string;
  mode?: string;
  include_deleted?: boolean;
  include_archived?: boolean;
  include_transferred?: boolean;
  limit?: number;
  offset?: number;
};

export type AdminPatchItemPayload = {
  title?: string;
  mode?: string;
  status?: string;
};

export type ModerationPayload = {
  reason?: string;
};

export type AdminListEventsParams = {
  entity_type?: string;
  entity_id?: number;
  actor_user_id?: number;
  limit?: number;
  offset?: number;
};

function withQuery(path: string, params: Record<string, unknown>) {
  const q = new URLSearchParams();
  Object.entries(params).forEach(([k, v]) => {
    if (v === undefined || v === null || v === "") return;
    if (typeof v === "boolean") {
      if (v) q.set(k, "true");
      return;
    }
    q.set(k, String(v));
  });
  const qs = q.toString();
  return qs ? `${path}?${qs}` : path;
}

export const adminApi = {
  users: {
    list: (params: AdminListUsersParams = {}) =>
      apiFetch<AdminListResp<AdminUser, "users">>(
        withQuery("/api/admin/users", params),
        { method: "GET" }
      ),
    get: (id: number) => apiFetch<AdminUser>(`/api/admin/users/${id}`, { method: "GET" }),
    patch: (id: number, payload: AdminPatchUserPayload) =>
      apiFetch<AdminUser>(`/api/admin/users/${id}`, {
        method: "PATCH",
        body: JSON.stringify(payload),
      }),
  },

  bookings: {
    list: (params: AdminListBookingsParams = {}) =>
      apiFetch<AdminListResp<AdminBooking, "bookings">>(
        withQuery("/api/admin/bookings", params),
        { method: "GET" }
      ),
    get: (id: number) => apiFetch<AdminBooking>(`/api/admin/bookings/${id}`, { method: "GET" }),
    events: (id: number, params: { limit?: number; offset?: number } = {}) =>
      apiFetch<AdminListResp<AdminBookingEvent, "events">>(
        withQuery(`/api/admin/bookings/${id}/events`, params),
        { method: "GET" }
      ),
  },

  items: {
    list: (params: AdminListItemsParams = {}) =>
      apiFetch<AdminListResp<AdminItem, "items">>(withQuery("/api/admin/items", params), { method: "GET" }),
    patch: (id: number, payload: AdminPatchItemPayload) =>
      apiFetch<AdminItem>(`/api/admin/items/${id}`, {
        method: "PATCH",
        body: JSON.stringify(payload),
      }),
    block: (id: number, payload: ModerationPayload = {}) =>
      apiFetch<AdminItem>(`/api/admin/items/${id}/block`, {
        method: "POST",
        body: JSON.stringify(payload),
      }),
    unblock: (id: number, payload: ModerationPayload = {}) =>
      apiFetch<AdminItem>(`/api/admin/items/${id}/unblock`, {
        method: "POST",
        body: JSON.stringify(payload),
      }),
    delete: (id: number, payload: ModerationPayload = {}) =>
      apiFetch<AdminItem>(`/api/admin/items/${id}/delete`, {
        method: "POST",
        body: JSON.stringify(payload),
      }),
  },

  events: {
    list: (params: AdminListEventsParams = {}) =>
      apiFetch<AdminListResp<AdminEvent, "events">>(withQuery("/api/admin/events", params), { method: "GET" }),
  },
};

