"use client";

import { useSession } from "@/shared/auth/store";
import { apiFetch } from "@/shared/api/client";

export default function SecurityPage() {
  const clear = useSession((s) => s.clear);

  return (
    <div className="p-6 space-y-4">
      <h1 className="text-2xl font-semibold">Безопасность</h1>

      <button
        className="border rounded-lg px-3 py-2"
        onClick={async () => {
          await apiFetch("/api/auth/logout", { method: "POST" }, { auth: false });
          clear();
          location.href = "/auth/login";
        }}
      >
        Выйти
      </button>

      <button
        className="border rounded-lg px-3 py-2"
        onClick={async () => {
          await apiFetch("/api/auth/logout_all", { method: "POST" });
          clear();
          location.href = "/auth/login";
        }}
      >
        Выйти со всех устройств
      </button>
    </div>
  );
}
