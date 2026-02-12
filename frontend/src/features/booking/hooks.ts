import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { bookingApi, type BookingListParams } from "@/shared/api/endpoints/booking";

export function useMyBookings(params: BookingListParams) {
  return useQuery({
    queryKey: ["booking", "my", params],
    queryFn: () => bookingApi.listMy(params),
  });
}

export function useMyItemsBookings(params: BookingListParams) {
  return useQuery({
    queryKey: ["booking", "myItems", params],
    queryFn: () => bookingApi.listMyItems(params),
  });
}

export function useBookingEvents(id: number) {
  return useQuery({
    queryKey: ["booking", id, "events"],
    queryFn: () => bookingApi.events(id),
    enabled: Number.isFinite(id) && id > 0,
  });
}

export function useUpcomingByItem(itemId: number) {
  return useQuery({
    queryKey: ["item", itemId, "upcoming"],
    queryFn: () => bookingApi.upcomingByItem(itemId),
    enabled: Number.isFinite(itemId) && itemId > 0,
  });
}

export function useBusyByItem(itemId: number) {
  return useQuery({
    queryKey: ["item", itemId, "busyBookings"],
    queryFn: () => bookingApi.listBusyForItem(itemId),
    enabled: Number.isFinite(itemId) && itemId > 0,
  });
}

export function useMyItemsBookingRequests(params: { types?: Array<"rent" | "buy" | "give">; limit?: number; offset?: number }) {
  return useQuery({
    queryKey: ["booking", "myItemsRequests", params],
    queryFn: () => bookingApi.listMyItemsBookingRequests(params),
  });
}

export function useBookingActions() {
  const qc = useQueryClient();

  const wrap = (fn: (id: number) => Promise<any>) =>
    useMutation({
      mutationFn: fn,
      onSuccess: () => {
        qc.invalidateQueries({ queryKey: ["booking", "my"] });
        qc.invalidateQueries({ queryKey: ["booking", "myItems"] });
      },
    });

  return {
    approve: wrap(bookingApi.approve),
    handover: wrap(bookingApi.handover),
    returnB: wrap(bookingApi.return),
    cancel: wrap(bookingApi.cancel),
  };
}
