import { create } from "zustand";
import { apiFetch } from "@/lib/api";

export type LayoutItem = {
  id: string;
  type: string;
  label: string;
  x: number;
  y: number;
  width: number;
  height: number;
  rotation: number;
};

export type LayoutData = {
  width: number;
  height: number;
  items: LayoutItem[];
};

export type LayoutSummary = {
  id: number;
  user_id: number;
  title: string;
  item_count: number;
  created_at: string;
  updated_at: string;
};

export type LayoutDetail = {
  id: number;
  user_id: number;
  title: string;
  data: string;
  created_at: string;
  updated_at: string;
};

type LayoutState = {
  layouts: LayoutSummary[];
  currentLayout: LayoutDetail | null;
  loading: boolean;

  fetchLayouts: () => Promise<void>;
  fetchLayout: (id: number) => Promise<void>;
  createLayout: (title: string) => Promise<LayoutSummary | null>;
  updateLayout: (id: number, title: string, data: LayoutData) => Promise<void>;
  deleteLayout: (id: number) => Promise<void>;
};

export const useLayoutStore = create<LayoutState>((set, get) => ({
  layouts: [],
  currentLayout: null,
  loading: false,

  fetchLayouts: async () => {
    set({ loading: true });
    try {
      const data = await apiFetch<LayoutSummary[]>("/api/layouts");
      set({ layouts: data, loading: false });
    } catch {
      set({ loading: false });
    }
  },

  fetchLayout: async (id: number) => {
    set({ loading: true });
    try {
      const data = await apiFetch<LayoutDetail>(`/api/layouts/${id}`);
      set({ currentLayout: data, loading: false });
    } catch {
      set({ loading: false });
    }
  },

  createLayout: async (title: string) => {
    try {
      const data = await apiFetch<LayoutSummary>("/api/layouts", {
        method: "POST",
        body: JSON.stringify({ title }),
      });
      await get().fetchLayouts();
      return data;
    } catch {
      return null;
    }
  },

  updateLayout: async (id: number, title: string, data: LayoutData) => {
    try {
      const result = await apiFetch<LayoutDetail>(`/api/layouts/${id}`, {
        method: "PUT",
        body: JSON.stringify({ title, data }),
      });
      set({ currentLayout: result });
    } catch {
      // ignore
    }
  },

  deleteLayout: async (id: number) => {
    try {
      await apiFetch(`/api/layouts/${id}`, { method: "DELETE" });
      set({ layouts: get().layouts.filter((l) => l.id !== id) });
    } catch {
      // ignore
    }
  },
}));
