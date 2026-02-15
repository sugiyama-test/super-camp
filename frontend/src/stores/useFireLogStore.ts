import { create } from "zustand";
import { apiFetch } from "@/lib/api";

export type FireLog = {
  id: number;
  user_id: number;
  date: string;
  location: string;
  wood_type: string;
  duration_minutes: number;
  notes: string;
  temperature: number | null;
  campsite_id: number | null;
  created_at: string;
  updated_at: string;
};

export type CreateFireLogData = {
  date: string;
  location: string;
  wood_type: string;
  duration_minutes: number;
  notes: string;
  temperature: number | null;
};

type FireLogState = {
  fireLogs: FireLog[];
  loading: boolean;

  fetchFireLogs: () => Promise<void>;
  createFireLog: (data: CreateFireLogData) => Promise<FireLog | null>;
  updateFireLog: (id: number, data: CreateFireLogData) => Promise<void>;
  deleteFireLog: (id: number) => Promise<void>;
};

export const useFireLogStore = create<FireLogState>((set, get) => ({
  fireLogs: [],
  loading: false,

  fetchFireLogs: async () => {
    set({ loading: true });
    try {
      const data = await apiFetch<FireLog[]>("/api/fire-logs");
      set({ fireLogs: data, loading: false });
    } catch {
      set({ loading: false });
    }
  },

  createFireLog: async (data: CreateFireLogData) => {
    try {
      const result = await apiFetch<FireLog>("/api/fire-logs", {
        method: "POST",
        body: JSON.stringify(data),
      });
      await get().fetchFireLogs();
      return result;
    } catch {
      return null;
    }
  },

  updateFireLog: async (id: number, data: CreateFireLogData) => {
    try {
      await apiFetch<FireLog>(`/api/fire-logs/${id}`, {
        method: "PUT",
        body: JSON.stringify(data),
      });
      await get().fetchFireLogs();
    } catch {
      // ignore
    }
  },

  deleteFireLog: async (id: number) => {
    try {
      await apiFetch(`/api/fire-logs/${id}`, { method: "DELETE" });
      set({ fireLogs: get().fireLogs.filter((l) => l.id !== id) });
    } catch {
      // ignore
    }
  },
}));
