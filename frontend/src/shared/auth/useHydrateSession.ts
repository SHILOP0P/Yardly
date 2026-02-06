"use client";

import { useEffect } from "react";
import { useSession } from "./store";

export function useHydrateSession() {
  const hydrate = useSession((s) => s.hydrateFromStorage);
  useEffect(() => hydrate(), [hydrate]);
}
