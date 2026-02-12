"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const links = [
  { href: "/admin", label: "Обзор" },
  { href: "/admin/users", label: "Пользователи" },
  { href: "/admin/bookings", label: "Бронирования" },
  { href: "/admin/items", label: "Items" },
  { href: "/admin/events", label: "Аудит" },
];

export function AdminNav() {
  const pathname = usePathname();

  return (
    <nav className="flex flex-wrap gap-2">
      {links.map((l) => {
        const active = pathname === l.href || pathname.startsWith(l.href + "/");
        return (
          <Link
            key={l.href}
            href={l.href}
            className={`border rounded-lg px-3 py-2 text-sm ${active ? "bg-muted" : "hover:bg-muted/60"}`}
          >
            {l.label}
          </Link>
        );
      })}
    </nav>
  );
}

