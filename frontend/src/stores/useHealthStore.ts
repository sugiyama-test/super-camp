import { create } from "zustand";
import { apiFetch } from "@/lib/api";

type HealthState = {
  status: string | null;
  loading: boolean;
  check: () => Promise<void>;
};

export const useHealthStore = create<HealthState>((set) => ({
  status: null,
  loading: false,
  check: async () => {
    set({ loading: true });
    try {
      const data = await apiFetch<{ status: string }>("/api/health");
      set({ status: data.status, loading: false });
    } catch {
      set({ status: "error", loading: false });
    }
  },
}));
