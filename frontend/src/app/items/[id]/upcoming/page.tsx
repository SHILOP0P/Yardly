"use client";

import { useParams } from "next/navigation";
import { useUpcomingByItem } from "@/features/booking/hooks";

export default function ItemUpcomingPage() {
  const params = useParams<{ id: string }>();
  const itemId = Number(params.id);

  const q = useUpcomingByItem(itemId);

  if (q.isLoading) return <div className="p-6">Загрузка...</div>;
  if (q.error) return <div className="p-6">Ошибка</div>;

  return (
    <div className="p-6 space-y-4">
      <h1 className="text-2xl font-semibold">Upcoming по item #{itemId}</h1>

      <pre className="border rounded-xl p-4 text-sm overflow-auto">
        {JSON.stringify(q.data, null, 2)}
      </pre>
    </div>
  );
}
