"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { authApi } from "@/shared/api/endpoints/auth";
import { useSession } from "@/shared/auth/store";
import { useHydrateSession } from "@/shared/auth/useHydrateSession";

export default function LoginPage() {
  useHydrateSession();

  const router = useRouter();
  const setTokens = useSession((s) => s.setAccessToken);

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [err, setErr] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setErr(null);
    setLoading(true);
    try {
      const t = await authApi.login(email, password);
      setTokens(t.access_token);
      router.push("/items");
    } catch (e: any) {
      setErr(e?.message ?? "login failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-6">
      <form onSubmit={onSubmit} className="w-full max-w-sm space-y-3 border rounded-xl p-4">
        <h1 className="text-xl font-semibold">Вход</h1>

        <input className="w-full border rounded-lg p-2" placeholder="email" value={email} onChange={(e) => setEmail(e.target.value)} />
        <input className="w-full border rounded-lg p-2" placeholder="password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} />

        {err && <div className="text-sm text-red-500">{err}</div>}

        <button className="w-full border rounded-lg p-2" disabled={loading}>
          {loading ? "..." : "Войти"}
        </button>

        <button
          type="button"
          className="w-full border rounded-lg p-2 opacity-80"
          onClick={() => router.push("/auth/register")}
        >
          Регистрация
        </button>
      </form>
    </div>
  );
}
