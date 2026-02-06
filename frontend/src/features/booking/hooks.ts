import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { bookingApi } from "@/shared/api/endpoints/booking";

export function useMyBookings() {
  return useQuery({ queryKey: ["booking", "my"], queryFn: bookingApi.listMy });
}

export function useMyItemsBookings() {
  return useQuery({ queryKey: ["booking", "myItems"], queryFn: bookingApi.listMyItems });
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
