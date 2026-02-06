import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { favoriteApi } from "@/shared/api/endpoints/favorite";

export function useIsFavorite(itemId: number) {
  return useQuery({
    queryKey: ["favorite", "is", itemId],
    queryFn: () => favoriteApi.isFavorite(itemId),
    enabled: Number.isFinite(itemId) && itemId > 0,
  });
}

export function useToggleFavorite(itemId: number) {
  const qc = useQueryClient();

  const add = useMutation({
    mutationFn: () => favoriteApi.add(itemId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["favorite", "is", itemId] });
      qc.invalidateQueries({ queryKey: ["favorite", "my"] });
    },
  });

  const remove = useMutation({
    mutationFn: () => favoriteApi.remove(itemId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["favorite", "is", itemId] });
      qc.invalidateQueries({ queryKey: ["favorite", "my"] });
    },
  });

  return { add, remove };
}

export function useMyFavorites() {
  return useQuery({
    queryKey: ["favorite", "my"],
    queryFn: () => favoriteApi.my(),
  });
}
