"use client";

import Link from "next/link";
import { AdminGate } from "@/features/admin/AdminGate";
import { AdminNav } from "@/features/admin/AdminNav";
import { useMe } from "@/shared/auth/useMe";

export default function AdminPage() {
  const meQ = useMe();

  return (
    <AdminGate>
      <div className="p-6 space-y-4">
        <h1 className="text-2xl font-semibold">Админ-панель</h1>
        <AdminNav />

        <div className="border rounded-xl p-4">
          <div className="text-sm opacity-70">Текущий оператор</div>
          <div className="font-medium">
            {meQ.data?.email} ({meQ.data?.role})
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
          <Link href="/admin/users" className="border rounded-xl p-4 hover:bg-muted/70">
            <div className="font-medium">Пользователи</div>
            <div className="text-sm opacity-70">GET/PATCH /api/admin/users*</div>
          </Link>
          <Link href="/admin/bookings" className="border rounded-xl p-4 hover:bg-muted/70">
            <div className="font-medium">Бронирования</div>
            <div className="text-sm opacity-70">GET /api/admin/bookings*</div>
          </Link>
          <Link href="/admin/items" className="border rounded-xl p-4 hover:bg-muted/70">
            <div className="font-medium">Items</div>
            <div className="text-sm opacity-70">GET/PATCH/moderation /api/admin/items*</div>
          </Link>
          <Link href="/admin/events" className="border rounded-xl p-4 hover:bg-muted/70">
            <div className="font-medium">Admin Events</div>
            <div className="text-sm opacity-70">GET /api/admin/events</div>
          </Link>
        </div>
      </div>
    </AdminGate>
  );
}

