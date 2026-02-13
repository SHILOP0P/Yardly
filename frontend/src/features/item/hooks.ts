import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { itemsApi, type ItemListParams } from "@/shared/api/endpoints/items";

export function useItemsList(params?: ItemListParams) {
  return useQuery({
    queryKey: ["items", "list", params ?? null],
    queryFn: () => itemsApi.list(params),
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
    mutationFn: (file: File) => itemsApi.images.add(itemId, file),
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
