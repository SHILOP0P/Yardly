"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useSession } from "@/shared/auth/store";

function NavLink({ href, label }: { href: string; label: string }) {
  const p = usePathname();
  const active = p === href || p.startsWith(href + "/");
  return (
    <Link
      href={href}
      className={`px-3 py-2 rounded-lg border text-sm ${
        active ? "opacity-100" : "opacity-70 hover:opacity-100"
      }`}
    >
      {label}
    </Link>
  );
}

export function AppHeader() {
  const access = useSession((s) => s.accessToken);

  return (
    <header className="sticky top-0 z-50 border-b bg-background/80 backdrop-blur">
      <div className="mx-auto max-w-5xl px-4 py-3 flex items-center justify-between gap-3">
        <Link href="/" className="font-semibold">
          Yardly
        </Link>

        <nav className="flex flex-wrap gap-2">
          <NavLink href="/items" label="Лента" />
          <NavLink href="/my/favorites" label="Избранное" />
          <NavLink href="/my/items" label="Мои вещи" />
          <NavLink href="/my/bookings" label="Мои заявки" />
          <NavLink href="/my/items/bookings" label="Заявки на мои вещи" />
          <NavLink href="/settings/security" label="Безопасность" />
          <NavLink href="/me" label="Профиль" />
        </nav>

        <div className="text-xs opacity-70">
          {access ? "auth: ON" : "auth: OFF"}
        </div>
      </div>
    </header>
  );
}
