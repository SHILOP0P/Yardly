import { apiFetch } from "@/shared/api/client";

export type Availability = {
  item_id: number;
  from: string; // YYYY-MM-DD
  to: string;   // YYYY-MM-DD
  timezone: string;
  is_in_use_now: boolean;
  busy: { start: string; end: string }[];
};

export const availabilityApi = {
  get: (itemId: number, from: string, to: string) =>
    apiFetch<Availability>(`/api/items/${itemId}/availability?from=${encodeURIComponent(from)}&to=${encodeURIComponent(to)}`, {
      method: "GET",
    }, { auth: false }),
};
