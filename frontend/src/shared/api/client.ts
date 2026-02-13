import { useSession } from "@/shared/auth/store";

export type ApiError = Error & {
  status?: number;
  payload?: any;
};

const BASE = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

function makeError(status: number, payload: any): ApiError {
  const e = new Error(payload?.error ?? payload?.message ?? `HTTP ${status}`) as ApiError;
  e.status = status;
  e.payload = payload;
  return e;
}

async function parseBody(res: Response) {
  const text = await res.text();
  if (!text) return null;
  try {
    return JSON.parse(text);
  } catch {
    return text;
  }
}

let refreshInFlight: Promise<void> | null = null;

async function doRefresh(): Promise<void> {
  const { setAccessToken, clear } = useSession.getState();

  const res = await fetch(`${BASE}/api/auth/refresh`, {
    method: "POST",
    credentials: "include",
  });

  const payload = await parseBody(res);

  if (!res.ok) {
    clear();
    throw makeError(res.status, payload);
  }

  const access = payload?.access_token;
  if (!access) {
    clear();
    throw makeError(500, { error: "invalid refresh response" });
  }

  setAccessToken(access);
}

export async function apiFetch<T>(
  path: string,
  init: RequestInit = {},
  opts: { auth?: boolean; retry401?: boolean } = { auth: true, retry401: true }
): Promise<T> {
  const { accessToken } = useSession.getState();

  const headers = new Headers(init.headers);
  if (!headers.has("Content-Type") && init.body && !(init.body instanceof FormData)) headers.set("Content-Type", "application/json");
  if (opts.auth !== false && accessToken) headers.set("Authorization", `Bearer ${accessToken}`);

  const res = await fetch(`${BASE}${path}`, { ...init, headers, credentials: "include" });

  if (res.status === 401 && opts.retry401 !== false && !path.startsWith("/api/auth/refresh")) {
    if (!refreshInFlight) {
      refreshInFlight = doRefresh().finally(() => (refreshInFlight = null));
    }
    await refreshInFlight;

    const { accessToken: newAccess } = useSession.getState();
    const headers2 = new Headers(init.headers);
    if (!headers2.has("Content-Type") && init.body && !(init.body instanceof FormData)) headers2.set("Content-Type", "application/json");
    if (opts.auth !== false && newAccess) headers2.set("Authorization", `Bearer ${newAccess}`);

    const res2 = await fetch(`${BASE}${path}`, { ...init, headers: headers2, credentials: "include" });
    const payload2 = await parseBody(res2);
    if (!res2.ok) throw makeError(res2.status, payload2);
    return payload2 as T;
  }

  const payload = await parseBody(res);
  if (!res.ok) throw makeError(res.status, payload);
  return payload as T;
}
