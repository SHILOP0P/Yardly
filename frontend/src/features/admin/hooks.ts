"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  adminApi,
  type AdminListBookingsParams,
  type AdminListEventsParams,
  type AdminListItemsParams,
  type AdminListUsersParams,
  type AdminPatchItemPayload,
  type AdminPatchUserPayload,
  type ModerationPayload,
} from "@/shared/api/endpoints/admin";

export function useAdminUsers(params: AdminListUsersParams) {
  return useQuery({
    queryKey: ["admin", "users", params],
    queryFn: () => adminApi.users.list(params),
  });
}

export function useAdminUser(id: number) {
  return useQuery({
    queryKey: ["admin", "users", id],
    queryFn: () => adminApi.users.get(id),
    enabled: Number.isFinite(id) && id > 0,
  });
}

export function useAdminBookings(params: AdminListBookingsParams) {
  return useQuery({
    queryKey: ["admin", "bookings", params],
    queryFn: () => adminApi.bookings.list(params),
  });
}

export function useAdminBooking(id: number) {
  return useQuery({
    queryKey: ["admin", "bookings", id],
    queryFn: () => adminApi.bookings.get(id),
    enabled: Number.isFinite(id) && id > 0,
  });
}

export function useAdminBookingEvents(id: number, params: { limit?: number; offset?: number }) {
  return useQuery({
    queryKey: ["admin", "bookings", id, "events", params],
    queryFn: () => adminApi.bookings.events(id, params),
    enabled: Number.isFinite(id) && id > 0,
  });
}

export function useAdminItems(params: AdminListItemsParams) {
  return useQuery({
    queryKey: ["admin", "items", params],
    queryFn: () => adminApi.items.list(params),
  });
}

export function useAdminEvents(params: AdminListEventsParams) {
  return useQuery({
    queryKey: ["admin", "events", params],
    queryFn: () => adminApi.events.list(params),
  });
}

export function useAdminActions() {
  const qc = useQueryClient();

  const invalidateUsers = () => {
    qc.invalidateQueries({ queryKey: ["admin", "users"] });
    qc.invalidateQueries({ queryKey: ["admin", "events"] });
  };

  const invalidateItems = () => {
    qc.invalidateQueries({ queryKey: ["admin", "items"] });
    qc.invalidateQueries({ queryKey: ["admin", "events"] });
  };

  const patchUser = useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: AdminPatchUserPayload }) => adminApi.users.patch(id, payload),
    onSuccess: invalidateUsers,
  });

  const patchItem = useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: AdminPatchItemPayload }) => adminApi.items.patch(id, payload),
    onSuccess: invalidateItems,
  });

  const blockItem = useMutation({
    mutationFn: ({ id, payload }: { id: number; payload?: ModerationPayload }) =>
      adminApi.items.block(id, payload ?? {}),
    onSuccess: invalidateItems,
  });

  const unblockItem = useMutation({
    mutationFn: ({ id, payload }: { id: number; payload?: ModerationPayload }) =>
      adminApi.items.unblock(id, payload ?? {}),
    onSuccess: invalidateItems,
  });

  const deleteItem = useMutation({
    mutationFn: ({ id, payload }: { id: number; payload?: ModerationPayload }) =>
      adminApi.items.delete(id, payload ?? {}),
    onSuccess: invalidateItems,
  });

  return { patchUser, patchItem, blockItem, unblockItem, deleteItem };
}

