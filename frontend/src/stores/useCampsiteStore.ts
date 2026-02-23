import { create } from "zustand";
import { apiFetch } from "@/lib/api";

export type Campsite = {
  id: number;
  name: string;
  address: string;
  latitude: number | null;
  longitude: number | null;
  notes: string;
  created_at: string;
  updated_at: string;
};

export type CreateCampsiteData = {
  name: string;
  address: string;
  latitude: number | null;
  longitude: number | null;
  notes: string;
};

type CampsiteState = {
  campsites: Campsite[];
  loading: boolean;

  fetchCampsites: () => Promise<void>;
  createCampsite: (data: CreateCampsiteData) => Promise<Campsite | null>;
  updateCampsite: (id: number, data: CreateCampsiteData) => Promise<void>;
  deleteCampsite: (id: number) => Promise<void>;
};

export const useCampsiteStore = create<CampsiteState>((set, get) => ({
  campsites: [],
  loading: false,

  fetchCampsites: async () => {
    set({ loading: true });
    try {
      const data = await apiFetch<Campsite[]>("/api/campsites");
      set({ campsites: data, loading: false });
    } catch {
      set({ loading: false });
    }
  },

  createCampsite: async (data: CreateCampsiteData) => {
    try {
      const result = await apiFetch<Campsite>("/api/campsites", {
        method: "POST",
        body: JSON.stringify(data),
      });
      await get().fetchCampsites();
      return result;
    } catch {
      return null;
    }
  },

  updateCampsite: async (id: number, data: CreateCampsiteData) => {
    try {
      await apiFetch<Campsite>(`/api/campsites/${id}`, {
        method: "PUT",
        body: JSON.stringify(data),
      });
      await get().fetchCampsites();
    } catch {
      // ignore
    }
  },

  deleteCampsite: async (id: number) => {
    try {
      await apiFetch(`/api/campsites/${id}`, { method: "DELETE" });
      set({ campsites: get().campsites.filter((c) => c.id !== id) });
    } catch {
      // ignore
    }
  },
}));
