"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { authApi } from "@/shared/api/endpoints/auth";
import { useSession } from "@/shared/auth/store";

export default function RegisterPage() {
  const router = useRouter();
  const setAccessToken = useSession((s) => s.setAccessToken);

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [err, setErr] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setErr(null);
    setLoading(true);
    try {
      await authApi.register({
        email,
        password,
        first_name: firstName,
        last_name: lastName.trim() || undefined,
      });

      const t = await authApi.login(email, password);
      setAccessToken(t.access_token);
      router.push("/items");
    } catch (e: any) {
      setErr(e?.message ?? "register failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-6">
      <form onSubmit={onSubmit} className="w-full max-w-sm space-y-3 border rounded-xl p-4">
        <h1 className="text-xl font-semibold">Регистрация</h1>

        <input
          className="w-full border rounded-lg p-2"
          placeholder="first_name"
          value={firstName}
          onChange={(e) => setFirstName(e.target.value)}
        />
        <input
          className="w-full border rounded-lg p-2"
          placeholder="last_name (optional)"
          value={lastName}
          onChange={(e) => setLastName(e.target.value)}
        />
        <input className="w-full border rounded-lg p-2" placeholder="email" value={email} onChange={(e) => setEmail(e.target.value)} />
        <input
          className="w-full border rounded-lg p-2"
          placeholder="password (min 8)"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />

        {err && <div className="text-sm text-red-500">{err}</div>}

        <button className="w-full border rounded-lg p-2" disabled={loading || !firstName.trim()}>
          {loading ? "..." : "Создать аккаунт"}
        </button>

        <button type="button" className="w-full border rounded-lg p-2 opacity-80" onClick={() => router.push("/auth/login")}>
          Назад ко входу
        </button>
      </form>
    </div>
  );
}

