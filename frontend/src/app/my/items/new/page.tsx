"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useMutation } from "@tanstack/react-query";
import { itemsApi } from "@/shared/api/endpoints/items";
import { useHydrateSession } from "@/shared/auth/useHydrateSession";

export default function NewItemPage() {
  useHydrateSession();
  const router = useRouter();

  const [title, setTitle] = useState("");
  const [mode, setMode] = useState<"sale" | "rent" | "free" | "sale_rent">("rent");

  const mut = useMutation({
    mutationFn: () => itemsApi.create({ title, mode }),
    onSuccess: (it: any) => router.push(`/items/${it.id}`),
  });

  return (
    <div className="p-6 space-y-4">
      <h1 className="text-2xl font-semibold">Создать товар</h1>

      <div className="space-y-2">
        <div className="text-sm opacity-70">Название</div>
        <input
          className="border rounded-lg px-3 py-2 w-full"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="Например: Дрель Bosch"
        />
      </div>

      <div className="space-y-2">
        <div className="text-sm opacity-70">Режим</div>
        <select
          className="border rounded-lg px-3 py-2"
          value={mode}
          onChange={(e) => setMode(e.target.value as any)}
        >
          <option value="rent">Сдать в аренду</option>
          <option value="sale">Продать</option>
          <option value="free">Отдать бесплатно</option>
          <option value="sale_rent">Продать или сдать</option>
        </select>
      </div>

      <button
        className="border rounded-lg px-3 py-2"
        onClick={() => mut.mutate()}
        disabled={mut.isPending || !title.trim()}
      >
        Создать
      </button>

      {mut.error ? (
        <div className="text-sm text-red-600">{String((mut.error as any).message)}</div>
      ) : null}
    </div>
  );
}
