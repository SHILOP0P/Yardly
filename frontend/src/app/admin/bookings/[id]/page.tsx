"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { AdminGate } from "@/features/admin/AdminGate";
import { AdminNav } from "@/features/admin/AdminNav";
import { useAdminBooking } from "@/features/admin/hooks";

export default function AdminBookingDetailsPage() {
  const params = useParams<{ id: string }>();
  const id = Number(params.id);
  const bookingQ = useAdminBooking(id);

  return (
    <AdminGate>
      <div className="p-6 space-y-4">
        <h1 className="text-2xl font-semibold">Админка: бронь #{id}</h1>
        <AdminNav />

        <div className="flex gap-2">
          <Link className="border rounded-lg px-3 py-2" href="/admin/bookings">
            ← К списку
          </Link>
          <Link className="border rounded-lg px-3 py-2" href={`/admin/bookings/${id}/events`}>
            Events
          </Link>
        </div>

        {bookingQ.isLoading && <div>Загрузка...</div>}
        {bookingQ.error && <div className="text-red-600">Ошибка: {(bookingQ.error as Error).message}</div>}

        {bookingQ.data && (
          <div className="border rounded-xl p-4">
            <pre className="text-sm overflow-auto">{JSON.stringify(bookingQ.data, null, 2)}</pre>
          </div>
        )}
      </div>
    </AdminGate>
  );
}

