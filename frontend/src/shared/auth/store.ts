import { create } from "zustand";

type Session = {
  accessToken: string | null;
  setAccessToken: (access: string) => void;
  clear: () => void;
  hydrateFromStorage: () => void;
};

const LS_KEY = "yardly.access.v1";

export const useSession = create<Session>((set) => ({
  accessToken: null,

  setAccessToken: (access) => {
    set({ accessToken: access });
    try {
      localStorage.setItem(LS_KEY, access);
    } catch {}
  },

  clear: () => {
    set({ accessToken: null });
    try {
      localStorage.removeItem(LS_KEY);
    } catch {}
  },

  hydrateFromStorage: () => {
    try {
      const access = localStorage.getItem(LS_KEY);
      if (access) set({ accessToken: access });
    } catch {}
  },
}));
