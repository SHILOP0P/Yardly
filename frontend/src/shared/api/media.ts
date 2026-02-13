const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export function resolveMediaUrl(url: string): string {
  if (!url) return "";
  if (/^https?:\/\//i.test(url) || url.startsWith("data:")) return url;
  if (url.startsWith("/")) return `${API_BASE}${url}`;
  return `${API_BASE}/${url}`;
}
