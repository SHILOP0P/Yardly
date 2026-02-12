import Link from "next/link";

export default function HomePage() {
  return (
    <div className="p-6 space-y-4">
      <h1 className="text-2xl font-semibold">Yardly</h1>
      <p className="opacity-70">Обмен вещами внутри ЖК/кампуса: аренда, покупка, отдать бесплатно.</p>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        <Link className="border rounded-xl p-4" href="/items">
          <div className="font-medium">Смотреть ленту</div>
          <div className="text-sm opacity-70">GET /api/items</div>
        </Link>

        <Link className="border rounded-xl p-4" href="/my/items">
          <div className="font-medium">Мои вещи</div>
          <div className="text-sm opacity-70">GET /api/my/items + создание + фото</div>
        </Link>

        <Link className="border rounded-xl p-4" href="/my/bookings">
          <div className="font-medium">Мои заявки</div>
          <div className="text-sm opacity-70">GET /api/my/bookings</div>
        </Link>

        <Link className="border rounded-xl p-4" href="/my/items/bookings">
          <div className="font-medium">Заявки на мои вещи</div>
          <div className="text-sm opacity-70">approve/handover/return</div>
        </Link>

        <Link className="border rounded-xl p-4" href="/admin">
          <div className="font-medium">Админка</div>
          <div className="text-sm opacity-70">Пользователи, брони, items, audit</div>
        </Link>
      </div>
    </div>
  );
}

