import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useSession } from "@/shared/auth/store";
import { favoriteApi } from "@/shared/api/endpoints/favorite";

export function useIsFavorite(itemId: number) {
  const access = useSession((s) => s.accessToken);

  return useQuery({
    queryKey: ["favorite", itemId, "is"],
    queryFn: () => favoriteApi.isFavorite(itemId),
    enabled: !!access && itemId > 0, // ✅ ключевой фикс
    retry: false,
  });
}


export function useToggleFavorite(itemId: number) {
  const qc = useQueryClient();

  const add = useMutation({
    mutationFn: () => favoriteApi.add(itemId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["favorite", itemId, "is"] });
      qc.invalidateQueries({ queryKey: ["favorite", "my"] });
    },
  });

  const remove = useMutation({
    mutationFn: () => favoriteApi.remove(itemId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["favorite", itemId, "is"] });
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
