"use client";

import { useMemo, useState } from "react";
import Link from "next/link";
import { AdminGate } from "@/features/admin/AdminGate";
import { AdminNav } from "@/features/admin/AdminNav";
import { useAdminActions, useAdminUsers } from "@/features/admin/hooks";
import type { UserRole } from "@/shared/api/types";

const ROLES: UserRole[] = ["user", "admin", "superadmin"];

export default function AdminUsersPage() {
  const [q, setQ] = useState("");
  const [search, setSearch] = useState("");
  const [limit, setLimit] = useState(20);
  const [offset, setOffset] = useState(0);

  const params = useMemo(() => ({ q: search, limit, offset }), [search, limit, offset]);
  const usersQ = useAdminUsers(params);
  const actions = useAdminActions();

  const users = usersQ.data?.users ?? [];
  const canNext = users.length === limit;
  const canPrev = offset > 0;

  return (
    <AdminGate>
      <div className="p-6 space-y-4">
        <h1 className="text-2xl font-semibold">Админка: пользователи</h1>
        <AdminNav />

        <form
          className="border rounded-xl p-4 flex flex-wrap items-center gap-3"
          onSubmit={(e) => {
            e.preventDefault();
            setOffset(0);
            setSearch(q);
          }}
        >
          <input
            value={q}
            onChange={(e) => setQ(e.target.value)}
            placeholder="Поиск по email"
            className="border rounded-lg px-3 py-2 min-w-64"
          />
          <select
            value={limit}
            onChange={(e) => {
              setOffset(0);
              setLimit(Number(e.target.value));
            }}
            className="border rounded-lg px-3 py-2"
          >
            {[10, 20, 50, 100].map((v) => (
              <option key={v} value={v}>
                limit {v}
              </option>
            ))}
          </select>
          <button className="border rounded-lg px-3 py-2">Искать</button>
          <button
            type="button"
            className="border rounded-lg px-3 py-2"
            onClick={() => {
              setQ("");
              setSearch("");
              setOffset(0);
            }}
          >
            Сброс
          </button>

          <div className="ml-auto flex gap-2">
            <button
              type="button"
              className="border rounded-lg px-3 py-2 disabled:opacity-50"
              disabled={!canPrev || usersQ.isFetching}
              onClick={() => setOffset((s) => Math.max(0, s - limit))}
            >
              Назад
            </button>
            <button
              type="button"
              className="border rounded-lg px-3 py-2 disabled:opacity-50"
              disabled={!canNext || usersQ.isFetching}
              onClick={() => setOffset((s) => s + limit)}
            >
              Вперед
            </button>
          </div>
        </form>

        {usersQ.isLoading && <div>Загрузка...</div>}
        {usersQ.error && <div className="text-red-600">Ошибка: {(usersQ.error as Error).message}</div>}

        <div className="space-y-3">
          {users.map((u) => (
            <div key={u.id} className="border rounded-xl p-4 space-y-3">
              <div className="flex flex-wrap items-center justify-between gap-2">
                <div>
                  <div className="font-medium">#{u.id} {u.email}</div>
                  <div className="text-sm opacity-70">
                    role: {u.role}
                    {u.banned_at ? ` • banned_at: ${new Date(u.banned_at).toLocaleString("ru-RU")}` : ""}
                  </div>
                </div>
                <Link className="border rounded-lg px-3 py-2 text-sm" href={`/admin/users/${u.id}`}>
                  Детали
                </Link>
              </div>

              <div className="flex flex-wrap gap-2">
                {ROLES.map((role) => (
                  <button
                    key={role}
                    className="border rounded-lg px-3 py-2 text-sm disabled:opacity-50"
                    disabled={actions.patchUser.isPending || role === u.role}
                    onClick={() => actions.patchUser.mutate({ id: u.id, payload: { role } })}
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
                      id: u.id,
                      payload: {
                        ban: true,
                        ban_reason: reason || undefined,
                      },
                    });
                  }}
                >
                  Ban
                </button>
                <button
                  className="border rounded-lg px-3 py-2 text-sm disabled:opacity-50"
                  disabled={actions.patchUser.isPending}
                  onClick={() => actions.patchUser.mutate({ id: u.id, payload: { ban: false } })}
                >
                  Unban
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </AdminGate>
  );
}

