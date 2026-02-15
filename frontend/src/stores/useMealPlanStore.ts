import { create } from "zustand";
import { apiFetch } from "@/lib/api";

export type MealPlan = {
  id: number;
  user_id: number;
  title: string;
  meal_type: string;
  servings: number;
  notes: string;
  created_at: string;
  updated_at: string;
};

export type CreateMealPlanData = {
  title: string;
  meal_type: string;
  servings: number;
  notes: string;
};

type MealPlanState = {
  mealPlans: MealPlan[];
  loading: boolean;

  fetchMealPlans: () => Promise<void>;
  createMealPlan: (data: CreateMealPlanData) => Promise<MealPlan | null>;
  updateMealPlan: (id: number, data: CreateMealPlanData) => Promise<void>;
  deleteMealPlan: (id: number) => Promise<void>;
};

export const useMealPlanStore = create<MealPlanState>((set, get) => ({
  mealPlans: [],
  loading: false,

  fetchMealPlans: async () => {
    set({ loading: true });
    try {
      const data = await apiFetch<MealPlan[]>("/api/meal-plans");
      set({ mealPlans: data, loading: false });
    } catch {
      set({ loading: false });
    }
  },

  createMealPlan: async (data: CreateMealPlanData) => {
    try {
      const result = await apiFetch<MealPlan>("/api/meal-plans", {
        method: "POST",
        body: JSON.stringify(data),
      });
      await get().fetchMealPlans();
      return result;
    } catch {
      return null;
    }
  },

  updateMealPlan: async (id: number, data: CreateMealPlanData) => {
    try {
      await apiFetch<MealPlan>(`/api/meal-plans/${id}`, {
        method: "PUT",
        body: JSON.stringify(data),
      });
      await get().fetchMealPlans();
    } catch {
      // ignore
    }
  },

  deleteMealPlan: async (id: number) => {
    try {
      await apiFetch(`/api/meal-plans/${id}`, { method: "DELETE" });
      set({ mealPlans: get().mealPlans.filter((m) => m.id !== id) });
    } catch {
      // ignore
    }
  },
}));
