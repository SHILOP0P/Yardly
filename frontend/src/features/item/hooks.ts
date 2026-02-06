import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { itemsApi } from "@/shared/api/endpoints/items";

export function useItemsList(mode?: any) {
  return useQuery({
    queryKey: ["items", "list", mode ?? null],
    queryFn: () => itemsApi.list(mode ? { mode } : undefined),
  });
}

export function useItem(id: number) {
  return useQuery({
    queryKey: ["items", "byId", id],
    queryFn: () => itemsApi.getById(id),
    enabled: Number.isFinite(id) && id > 0,
  });
}

export function useItemImages(itemId: number) {
  return useQuery({
    queryKey: ["items", itemId, "images"],
    queryFn: () => itemsApi.images.list(itemId),
    enabled: Number.isFinite(itemId) && itemId > 0,
  });
}

export function useAddItemImage(itemId: number) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (url: string) => itemsApi.images.add(itemId, url),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["items", itemId, "images"] }),
  });
}

export function useDeleteItemImage(itemId: number) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (imageId: number) => itemsApi.images.delete(itemId, imageId),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["items", itemId, "images"] }),
  });
}
