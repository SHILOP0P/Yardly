"use client";

import { useParams } from "next/navigation";
import { useBookingEvents } from "@/features/booking/hooks";
import { useHydrateSession } from "@/shared/auth/useHydrateSession";

type Status = string;

type Event = {
  id: number;
  booking_id: number;
  actor_user_id?: number | null;
  action: string;
  from_status?: Status | null;
  to_status?: Status | null;
  meta?: any;
  created_at: string; // JSON time
};

function formatWhen(iso: string) {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return iso;
  return d.toLocaleString("ru-RU");
}

function actorLabel(e: Event) {
  if (e.actor_user_id == null) return "система";
  return `user #${e.actor_user_id}`;
}

/**
 * Здесь ты можешь привести "action" к нормальным названиям.
 * Я не знаю твой полный список action, поэтому:
 * - даю хорошие дефолты
 * - добавляю несколько самых вероятных
 */
function actionTitle(action: string) {
  switch (action) {
    case "created":
    case "booking_created":
      return "Создана заявка";
    case "approved":
    case "booking_approved":
      return "Владелец одобрил";
    case "declined":
    case "booking_declined":
      return "Отклонено";
    case "canceled":
    case "booking_canceled":
      return "Отменено";
    case "handover_confirmed":
    case "handover":
      return "Подтверждение передачи";
    case "return_confirmed":
    case "return":
      return "Подтверждение возврата";
    case "expired":
    case "booking_expired":
      return "Истекло по дедлайну";
    case "completed":
    case "booking_completed":
      return "Завершено";
    default:
      // fallback: красиво показать неизвестный action
      return `Событие: ${action}`;
  }
}

function statusChip(s?: Status | null) {
  if (!s) return null;
  return (
    <span className="inline-flex items-center rounded-full border px-2 py-0.5 text-xs">
      {s}
    </span>
  );
}

function prettyMeta(meta: any) {
  if (meta == null) return null;

  // Самый полезный кейс из твоего скрина: { "by": "owner" | "requester" }
  if (typeof meta === "object" && !Array.isArray(meta) && meta.by) {
    const by =
      meta.by === "owner" ? "владелец" :
      meta.by === "requester" ? "заявитель" :
      String(meta.by);

    return (
      <div className="text-sm">
        <span className="opacity-70">Кто подтвердил:</span> {by}
      </div>
    );
  }

  // fallback: показываем JSON
  return (
    <pre className="mt-2 text-xs opacity-80 overflow-auto border rounded-lg p-2">
      {JSON.stringify(meta, null, 2)}
    </pre>
  );
}

export default function BookingEventsPage() {
  useHydrateSession();

  const params = useParams<{ id: string }>();
  const id = Number(params.id);

  const q = useBookingEvents(id);

  if (q.isLoading) return <div className="p-6">Загрузка...</div>;
  if (q.error) return <div className="p-6">Ошибка / 401</div>;

  const events = (q.data ?? []) as Event[];

  // Сортировка на всякий случай (если бек уже отсортировал — не повредит)
  events.sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime());

  return (
    <div className="p-6 space-y-4">
      <div>
        <h1 className="text-2xl font-semibold">События брони #{id}</h1>
        <div className="text-sm opacity-70">
          Всего событий: {events.length}
        </div>
      </div>

      <div className="space-y-3">
        {events.map((e) => (
          <div key={e.id} className="border rounded-xl p-4">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <div className="font-medium">{actionTitle(e.action)}</div>
              <div className="text-xs opacity-70">{formatWhen(e.created_at)}</div>
            </div>

            <div className="mt-2 flex flex-wrap items-center gap-2 text-sm">
              <span className="opacity-70">Актор:</span> {actorLabel(e)}
              <span className="opacity-40">•</span>
              <span className="opacity-70">Action:</span>{" "}
              <span className="font-mono text-xs">{e.action}</span>
            </div>

            {(e.from_status || e.to_status) && (
              <div className="mt-2 flex flex-wrap items-center gap-2 text-sm">
                <span className="opacity-70">Статус:</span>
                {statusChip(e.from_status)}
                <span className="opacity-50">→</span>
                {statusChip(e.to_status)}
              </div>
            )}

            <div className="mt-3">{prettyMeta(e.meta)}</div>
          </div>
        ))}
      </div>

      {events.length === 0 && (
        <div className="opacity-70">Нет событий</div>
      )}
    </div>
  );
}
