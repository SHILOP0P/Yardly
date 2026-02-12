"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { AdminGate } from "@/features/admin/AdminGate";
import { AdminNav } from "@/features/admin/AdminNav";
import { useAdminActions, useAdminUser } from "@/features/admin/hooks";
import type { UserRole } from "@/shared/api/types";

const ROLES: UserRole[] = ["user", "admin", "superadmin"];

export default function AdminUserDetailsPage() {
  const params = useParams<{ id: string }>();
  const id = Number(params.id);
  const userQ = useAdminUser(id);
  const actions = useAdminActions();

  return (
    <AdminGate>
      <div className="p-6 space-y-4">
        <h1 className="text-2xl font-semibold">Админка: пользователь #{id}</h1>
        <AdminNav />
        <Link href="/admin/users" className="inline-flex border rounded-lg px-3 py-2">
          ← К списку
        </Link>

        {userQ.isLoading && <div>Загрузка...</div>}
        {userQ.error && <div className="text-red-600">Ошибка: {(userQ.error as Error).message}</div>}

        {userQ.data && (
          <div className="border rounded-xl p-4 space-y-4">
            <div className="space-y-1">
              <div className="font-medium">{userQ.data.email}</div>
              <div className="text-sm opacity-70">id: {userQ.data.id}</div>
              <div className="text-sm opacity-70">role: {userQ.data.role}</div>
              <div className="text-sm opacity-70">banned_at: {userQ.data.banned_at ?? "—"}</div>
              <div className="text-sm opacity-70">ban_reason: {userQ.data.ban_reason ?? "—"}</div>
              <div className="text-sm opacity-70">created_at: {new Date(userQ.data.created_at).toLocaleString("ru-RU")}</div>
              <div className="text-sm opacity-70">updated_at: {new Date(userQ.data.updated_at).toLocaleString("ru-RU")}</div>
            </div>

            <div className="flex flex-wrap gap-2">
              {ROLES.map((role) => (
                <button
                  key={role}
                  className="border rounded-lg px-3 py-2 text-sm disabled:opacity-50"
                  disabled={actions.patchUser.isPending}
                  onClick={() => actions.patchUser.mutate({ id, payload: { role } })}
                >
                  role: {role}
                </button>
              ))}
            </div>

            <div className="flex flex-wrap gap-2">
              <button
                className="border rounded-lg px-3 py-2 text-sm disabled:opacity-50"
                disabled={actions.patchUser.isPending}
                onClick={() => {
                  const reason = prompt("Причина бана (опционально):") ?? "";
                  actions.patchUser.mutate({
                    id,
                    payload: { ban: true, ban_reason: reason || undefined },
                  });
                }}
              >
                Ban
              </button>
              <button
                className="border rounded-lg px-3 py-2 text-sm disabled:opacity-50"
                disabled={actions.patchUser.isPending}
                onClick={() => actions.patchUser.mutate({ id, payload: { ban: false } })}
              >
                Unban
              </button>
            </div>
          </div>
        )}
      </div>
    </AdminGate>
  );
}

