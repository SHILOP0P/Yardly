"use client";

import { useMemo, useState } from "react";
import Link from "next/link";
import { AdminGate } from "@/features/admin/AdminGate";
import { AdminNav } from "@/features/admin/AdminNav";
import { useAdminEvents } from "@/features/admin/hooks";

export default function AdminEventsPage() {
  const [entityType, setEntityType] = useState("");
  const [entityID, setEntityID] = useState("");
  const [actorUserID, setActorUserID] = useState("");
  const [limit, setLimit] = useState(20);
  const [offset, setOffset] = useState(0);

  const params = useMemo(
    () => ({
      entity_type: entityType || undefined,
      entity_id: entityID ? Number(entityID) : undefined,
      actor_user_id: actorUserID ? Number(actorUserID) : undefined,
      limit,
      offset,
    }),
    [entityType, entityID, actorUserID, limit, offset]
  );

  const eventsQ = useAdminEvents(params);
  const events = eventsQ.data?.events ?? [];

  return (
    <AdminGate>
      <div className="p-6 space-y-4">
        <h1 className="text-2xl font-semibold">Админка: audit events</h1>
        <AdminNav />

        <div className="border rounded-xl p-4 space-y-3">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-2">
            <input
              value={entityType}
              onChange={(e) => {
                setOffset(0);
                setEntityType(e.target.value);
              }}
              className="border rounded-lg px-3 py-2"
              placeholder="entity_type (user/item/...)"
            />
            <input
              value={entityID}
              onChange={(e) => {
                setOffset(0);
                setEntityID(e.target.value);
              }}
              className="border rounded-lg px-3 py-2"
              placeholder="entity_id"
            />
            <input
              value={actorUserID}
              onChange={(e) => {
                setOffset(0);
                setActorUserID(e.target.value);
              }}
              className="border rounded-lg px-3 py-2"
              placeholder="actor_user_id"
            />
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
              disabled={offset <= 0 || eventsQ.isFetching}
              onClick={() => setOffset((v) => Math.max(0, v - limit))}
            >
              Назад
            </button>
            <button
              className="border rounded-lg px-3 py-2 disabled:opacity-50"
              disabled={events.length < limit || eventsQ.isFetching}
              onClick={() => setOffset((v) => v + limit)}
            >
              Вперед
            </button>
            <div className="text-sm opacity-70 flex items-center">offset: {offset}</div>
          </div>
        </div>

        {eventsQ.isLoading && <div>Загрузка...</div>}
        {eventsQ.error && <div className="text-red-600">Ошибка: {(eventsQ.error as Error).message}</div>}

        <div className="space-y-3">
          {events.map((e) => (
            <div key={e.id} className="border rounded-xl p-4 space-y-2">
              <div className="font-medium">{e.action}</div>
              <div className="text-sm opacity-70">
                #{e.id} • actor #{e.actor_user_id} • {e.entity_type} #{e.entity_id}
              </div>
              <div className="text-sm opacity-70">{new Date(e.created_at).toLocaleString("ru-RU")}</div>
              {e.reason && <div className="text-sm">reason: {e.reason}</div>}
              {e.meta != null && <pre className="text-xs overflow-auto">{JSON.stringify(e.meta, null, 2)}</pre>}

              {e.entity_type === "user" && (
                <Link className="text-sm underline" href={`/admin/users/${e.entity_id}`}>
                  Открыть пользователя
                </Link>
              )}
            </div>
          ))}
        </div>
      </div>
    </AdminGate>
  );
}

