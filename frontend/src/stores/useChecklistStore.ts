import { create } from "zustand";
import { apiFetch } from "@/lib/api";

export type ChecklistItem = {
  id: number;
  checklist_id: number;
  name: string;
  is_checked: boolean;
  quantity: number;
  sort_order: number;
  created_at: string;
  updated_at: string;
};

export type ChecklistSummary = {
  id: number;
  user_id: number;
  title: string;
  item_count: number;
  checked_count: number;
  created_at: string;
  updated_at: string;
};

export type ChecklistDetail = {
  id: number;
  user_id: number;
  title: string;
  items: ChecklistItem[];
  created_at: string;
  updated_at: string;
};

type ChecklistState = {
  checklists: ChecklistSummary[];
  currentChecklist: ChecklistDetail | null;
  loading: boolean;

  fetchChecklists: () => Promise<void>;
  fetchChecklist: (id: number) => Promise<void>;
  createChecklist: (title: string) => Promise<ChecklistSummary | null>;
  deleteChecklist: (id: number) => Promise<void>;
  addItem: (checklistId: number, name: string, quantity?: number) => Promise<void>;
  updateItem: (checklistId: number, itemId: number, data: { name: string; is_checked: boolean; quantity: number }) => Promise<void>;
  deleteItem: (checklistId: number, itemId: number) => Promise<void>;
};

export const useChecklistStore = create<ChecklistState>((set, get) => ({
  checklists: [],
  currentChecklist: null,
  loading: false,

  fetchChecklists: async () => {
    set({ loading: true });
    try {
      const data = await apiFetch<ChecklistSummary[]>("/api/checklists");
      set({ checklists: data, loading: false });
    } catch {
      set({ loading: false });
    }
  },

  fetchChecklist: async (id: number) => {
    set({ loading: true });
    try {
      const data = await apiFetch<ChecklistDetail>(`/api/checklists/${id}`);
      set({ currentChecklist: data, loading: false });
    } catch {
      set({ loading: false });
    }
  },

  createChecklist: async (title: string) => {
    try {
      const data = await apiFetch<ChecklistSummary>("/api/checklists", {
        method: "POST",
        body: JSON.stringify({ title }),
      });
      await get().fetchChecklists();
      return data;
    } catch {
      return null;
    }
  },

  deleteChecklist: async (id: number) => {
    try {
      await apiFetch(`/api/checklists/${id}`, { method: "DELETE" });
      set({ checklists: get().checklists.filter((c) => c.id !== id) });
    } catch {
      // ignore
    }
  },

  addItem: async (checklistId: number, name: string, quantity = 1) => {
    try {
      await apiFetch(`/api/checklists/${checklistId}/items`, {
        method: "POST",
        body: JSON.stringify({ name, quantity }),
      });
      await get().fetchChecklist(checklistId);
    } catch {
      // ignore
    }
  },

  updateItem: async (checklistId: number, itemId: number, data) => {
    try {
      await apiFetch(`/api/checklists/${checklistId}/items/${itemId}`, {
        method: "PUT",
        body: JSON.stringify(data),
      });
      await get().fetchChecklist(checklistId);
    } catch {
      // ignore
    }
  },

  deleteItem: async (checklistId: number, itemId: number) => {
    try {
      await apiFetch(`/api/checklists/${checklistId}/items/${itemId}`, {
        method: "DELETE",
      });
      await get().fetchChecklist(checklistId);
    } catch {
      // ignore
    }
  },
}));
