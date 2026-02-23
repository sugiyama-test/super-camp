import { create } from "zustand";
import { apiFetch } from "@/lib/api";

export type Gear = {
  id: number;
  user_id: number;
  name: string;
  category: string;
  brand: string;
  weight_grams: number | null;
  notes: string;
  created_at: string;
  updated_at: string;
};

export type CreateGearData = {
  name: string;
  category: string;
  brand: string;
  weight_grams: number | null;
  notes: string;
};

type GearState = {
  gears: Gear[];
  loading: boolean;

  fetchGears: () => Promise<void>;
  createGear: (data: CreateGearData) => Promise<Gear | null>;
  updateGear: (id: number, data: CreateGearData) => Promise<void>;
  deleteGear: (id: number) => Promise<void>;
};

export const useGearStore = create<GearState>((set, get) => ({
  gears: [],
  loading: false,

  fetchGears: async () => {
    set({ loading: true });
    try {
      const data = await apiFetch<Gear[]>("/api/gears");
      set({ gears: data, loading: false });
    } catch {
      set({ loading: false });
    }
  },

  createGear: async (data: CreateGearData) => {
    try {
      const result = await apiFetch<Gear>("/api/gears", {
        method: "POST",
        body: JSON.stringify(data),
      });
      await get().fetchGears();
      return result;
    } catch {
      return null;
    }
  },

  updateGear: async (id: number, data: CreateGearData) => {
    try {
      await apiFetch<Gear>(`/api/gears/${id}`, {
        method: "PUT",
        body: JSON.stringify(data),
      });
      await get().fetchGears();
    } catch {
      // ignore
    }
  },

  deleteGear: async (id: number) => {
    try {
      await apiFetch(`/api/gears/${id}`, { method: "DELETE" });
      set({ gears: get().gears.filter((g) => g.id !== id) });
    } catch {
      // ignore
    }
  },
}));
