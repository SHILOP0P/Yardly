"use client";

import { useMemo, useState } from "react";
import { AdminGate } from "@/features/admin/AdminGate";
import { AdminNav } from "@/features/admin/AdminNav";
import { useAdminActions, useAdminItems } from "@/features/admin/hooks";

const statusOptions = ["", "active", "in_use", "archived", "deleted", "transferred"];
const modeOptions = ["", "sale", "rent", "free", "sale_rent"];

export default function AdminItemsPage() {
  const [q, setQ] = useState("");
  const [status, setStatus] = useState("");
  const [mode, setMode] = useState("");
  const [includeDeleted, setIncludeDeleted] = useState(false);
  const [includeArchived, setIncludeArchived] = useState(false);
  const [includeTransferred, setIncludeTransferred] = useState(false);
  const [limit, setLimit] = useState(20);
  const [offset, setOffset] = useState(0);

  const params = useMemo(
    () => ({
      q: q || undefined,
      status: status || undefined,
      mode: mode || undefined,
      include_deleted: includeDeleted || undefined,
      include_archived: includeArchived || undefined,
      include_transferred: includeTransferred || undefined,
      limit,
      offset,
    }),
    [q, status, mode, includeDeleted, includeArchived, includeTransferred, limit, offset]
  );

  const itemsQ = useAdminItems(params);
  const actions = useAdminActions();
  const items = itemsQ.data?.items ?? [];

  return (
    <AdminGate>
      <div className="p-6 space-y-4">
        <h1 className="text-2xl font-semibold">Админка: items</h1>
        <AdminNav />

        <div className="border rounded-xl p-4 space-y-3">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-2">
            <input
              value={q}
              onChange={(e) => {
                setOffset(0);
                setQ(e.target.value);
              }}
              className="border rounded-lg px-3 py-2"
              placeholder="поиск по title"
            />
            <select
              value={status}
              onChange={(e) => {
                setOffset(0);
                setStatus(e.target.value);
              }}
              className="border rounded-lg px-3 py-2"
            >
              {statusOptions.map((s) => (
                <option key={s} value={s}>
                  {s || "status: any"}
                </option>
              ))}
            </select>
            <select
              value={mode}
              onChange={(e) => {
                setOffset(0);
                setMode(e.target.value);
              }}
              className="border rounded-lg px-3 py-2"
            >
              {modeOptions.map((m) => (
                <option key={m} value={m}>
                  {m || "mode: any"}
                </option>
              ))}
            </select>
          </div>

          <div className="flex flex-wrap gap-3 text-sm">
            <label className="flex items-center gap-2">
              <input type="checkbox" checked={includeDeleted} onChange={(e) => setIncludeDeleted(e.target.checked)} />
              include_deleted
            </label>
            <label className="flex items-center gap-2">
              <input type="checkbox" checked={includeArchived} onChange={(e) => setIncludeArchived(e.target.checked)} />
              include_archived
            </label>
            <label className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={includeTransferred}
                onChange={(e) => setIncludeTransferred(e.target.checked)}
              />
              include_transferred
            </label>
          </div>

          <div className="flex gap-2">
            <select
              value={limit}
              className="border rounded-lg px-3 py-2"
              onChange={(e) => {
                setOffset(0);
                setLimit(Number(e.target.value));
              }}
            >
              {[10, 20, 50, 100].map((v) => (
                <option key={v} value={v}>
                  limit {v}
                </option>
              ))}
            </select>
            <button
              className="border rounded-lg px-3 py-2 disabled:opacity-50"
              disabled={offset <= 0 || itemsQ.isFetching}
              onClick={() => setOffset((v) => Math.max(0, v - limit))}
            >
              Назад
            </button>
            <button
              className="border rounded-lg px-3 py-2 disabled:opacity-50"
              disabled={items.length < limit || itemsQ.isFetching}
              onClick={() => setOffset((v) => v + limit)}
            >
              Вперед
            </button>
            <div className="text-sm opacity-70 flex items-center">offset: {offset}</div>
          </div>
        </div>

        {itemsQ.isLoading && <div>Загрузка...</div>}
        {itemsQ.error && <div className="text-red-600">Ошибка: {(itemsQ.error as Error).message}</div>}

        <div className="space-y-3">
          {items.map((it) => (
            <div key={it.id} className="border rounded-xl p-4 space-y-3">
              <div className="flex flex-wrap items-center justify-between gap-2">
                <div>
                  <div className="font-medium">
                    #{it.id} • {it.title}
                  </div>
                  <div className="text-sm opacity-70">
                    owner #{it.owner_id} • mode: {it.mode} • status: {it.status}
                  </div>
                  <div className="text-sm opacity-70">
                    blocked_at: {it.blocked_at ?? "—"} • block_reason: {it.block_reason ?? "—"}
                  </div>
                </div>
              </div>

              <div className="flex flex-wrap gap-2">
                <button
                  className="border rounded-lg px-3 py-2 text-sm disabled:opacity-50"
                  disabled={actions.patchItem.isPending}
                  onClick={() => {
                    const title = prompt("Новый title:", it.title);
                    if (!title || title === it.title) return;
                    actions.patchItem.mutate({ id: it.id, payload: { title } });
                  }}
                >
                  Изменить title
                </button>
                <button
                  className="border rounded-lg px-3 py-2 text-sm disabled:opacity-50"
                  disabled={actions.patchItem.isPending}
                  onClick={() => {
                    const next = prompt("Новый mode: sale/rent/free/sale_rent", it.mode);
                    if (!next || next === it.mode) return;
                    actions.patchItem.mutate({ id: it.id, payload: { mode: next } });
                  }}
                >
                  Изменить mode
                </button>
                <button
                  className="border rounded-lg px-3 py-2 text-sm disabled:opacity-50"
                  disabled={actions.patchItem.isPending}
                  onClick={() => {
                    const next = prompt("Новый status: active/in_use/archived/transferred", it.status);
                    if (!next || next === it.status) return;
                    actions.patchItem.mutate({ id: it.id, payload: { status: next } });
                  }}
                >
                  Изменить status
                </button>
              </div>

              <div className="flex flex-wrap gap-2">
                <button
                  className="border rounded-lg px-3 py-2 text-sm disabled:opacity-50"
                  disabled={actions.blockItem.isPending}
                  onClick={() => {
                    const reason = prompt("Причина блокировки (опционально):") ?? "";
                    actions.blockItem.mutate({ id: it.id, payload: { reason: reason || undefined } });
                  }}
                >
                  Block
                </button>
                <button
                  className="border rounded-lg px-3 py-2 text-sm disabled:opacity-50"
                  disabled={actions.unblockItem.isPending}
                  onClick={() => {
                    const reason = prompt("Причина разблокировки (опционально):") ?? "";
                    actions.unblockItem.mutate({ id: it.id, payload: { reason: reason || undefined } });
                  }}
                >
                  Unblock
                </button>
                <button
                  className="border rounded-lg px-3 py-2 text-sm text-red-700 disabled:opacity-50"
                  disabled={actions.deleteItem.isPending}
                  onClick={() => {
                    const reason = prompt("Причина удаления (опционально):") ?? "";
                    actions.deleteItem.mutate({ id: it.id, payload: { reason: reason || undefined } });
                  }}
                >
                  Delete (soft)
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </AdminGate>
  );
}

