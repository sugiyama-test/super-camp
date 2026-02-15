"use client";

import { useEffect, useState } from "react";
import { PageHeader } from "@/components/PageHeader";
import { useMealPlanStore } from "@/stores/useMealPlanStore";
import type { CreateMealPlanData, MealPlan } from "@/stores/useMealPlanStore";

const MEAL_TYPES: { value: string; label: string }[] = [
  { value: "breakfast", label: "朝食" },
  { value: "lunch", label: "昼食" },
  { value: "dinner", label: "夕食" },
  { value: "snack", label: "おやつ" },
];

function getMealTypeLabel(type: string) {
  return MEAL_TYPES.find((t) => t.value === type)?.label ?? type;
}

const emptyForm: CreateMealPlanData = {
  title: "",
  meal_type: "dinner",
  servings: 2,
  notes: "",
};

export default function MealPlansPage() {
  const { mealPlans, loading, fetchMealPlans, createMealPlan, updateMealPlan, deleteMealPlan } =
    useMealPlanStore();
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [form, setForm] = useState<CreateMealPlanData>(emptyForm);

  useEffect(() => {
    fetchMealPlans();
  }, [fetchMealPlans]);

  const resetForm = () => {
    setForm(emptyForm);
    setEditingId(null);
    setShowForm(false);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.title.trim()) return;
    if (editingId) {
      await updateMealPlan(editingId, form);
    } else {
      await createMealPlan(form);
    }
    resetForm();
  };

  const handleEdit = (plan: MealPlan) => {
    setForm({
      title: plan.title,
      meal_type: plan.meal_type,
      servings: plan.servings,
      notes: plan.notes,
    });
    setEditingId(plan.id);
    setShowForm(true);
  };

  const handleDelete = (id: number) => {
    if (confirm("このメニューを削除しますか？")) {
      deleteMealPlan(id);
    }
  };

  // Group by meal type
  const grouped = MEAL_TYPES.map((type) => ({
    ...type,
    plans: mealPlans.filter((p) => p.meal_type === type.value),
  })).filter((g) => g.plans.length > 0);

  return (
    <div>
      <PageHeader
        title="キャンプ飯プランナー"
        description="キャンプの食事を計画しましょう"
      />

      <div className="mt-6">
        {!showForm ? (
          <button
            onClick={() => { setForm(emptyForm); setEditingId(null); setShowForm(true); }}
            className="w-full rounded-lg bg-[var(--camp-brown)] px-4 py-2 text-sm font-medium text-white hover:opacity-90"
          >
            + 新しいメニュー
          </button>
        ) : (
          <form onSubmit={handleSubmit} className="rounded-xl bg-white p-4 shadow-sm space-y-3">
            <div>
              <label className="block text-xs text-gray-500 mb-1">メニュー名 *</label>
              <input
                type="text"
                value={form.title}
                onChange={(e) => setForm({ ...form, title: e.target.value })}
                placeholder="カレーライス、BBQ..."
                className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-brown)]"
                required
              />
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="block text-xs text-gray-500 mb-1">食事タイプ</label>
                <select
                  value={form.meal_type}
                  onChange={(e) => setForm({ ...form, meal_type: e.target.value })}
                  className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-brown)]"
                >
                  {MEAL_TYPES.map((t) => (
                    <option key={t.value} value={t.value}>{t.label}</option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-xs text-gray-500 mb-1">人数</label>
                <input
                  type="number"
                  value={form.servings}
                  onChange={(e) => setForm({ ...form, servings: Number(e.target.value) })}
                  min={1}
                  className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-brown)]"
                />
              </div>
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-1">メモ・レシピ</label>
              <textarea
                value={form.notes}
                onChange={(e) => setForm({ ...form, notes: e.target.value })}
                placeholder="材料、作り方のメモ..."
                rows={3}
                className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-brown)]"
              />
            </div>
            <div className="flex gap-2">
              <button
                type="submit"
                className="flex-1 rounded-lg bg-[var(--camp-brown)] px-4 py-2 text-sm font-medium text-white hover:opacity-90"
              >
                {editingId ? "更新" : "追加する"}
              </button>
              <button
                type="button"
                onClick={resetForm}
                className="rounded-lg border border-gray-300 px-4 py-2 text-sm text-gray-600 hover:bg-gray-50"
              >
                キャンセル
              </button>
            </div>
          </form>
        )}
      </div>

      <div className="mt-6 space-y-4">
        {loading && mealPlans.length === 0 && (
          <p className="text-center text-gray-400 text-sm">読み込み中...</p>
        )}
        {!loading && mealPlans.length === 0 && (
          <p className="text-center text-gray-400 text-sm mt-8">
            メニューがありません。上のボタンから追加しましょう！
          </p>
        )}
        {grouped.map((group) => (
          <div key={group.value}>
            <h3 className="text-sm font-medium text-gray-600 mb-2">{group.label}</h3>
            <div className="space-y-2">
              {group.plans.map((plan) => (
                <div key={plan.id} className="rounded-xl bg-white p-4 shadow-sm">
                  <div className="flex items-start justify-between">
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2">
                        <span className="font-medium text-gray-800 text-sm">{plan.title}</span>
                        <span className="text-xs text-gray-500">{plan.servings}人前</span>
                      </div>
                      {plan.notes && (
                        <p className="mt-1 text-xs text-gray-600 whitespace-pre-wrap">{plan.notes}</p>
                      )}
                    </div>
                    <div className="flex gap-1 ml-2">
                      <button
                        onClick={() => handleEdit(plan)}
                        className="text-gray-400 hover:text-gray-600 p-1"
                        aria-label="編集"
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                          <path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z" />
                        </svg>
                      </button>
                      <button
                        onClick={() => handleDelete(plan.id)}
                        className="text-gray-400 hover:text-red-500 p-1"
                        aria-label="削除"
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                          <path d="M3 6h18" /><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" /><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
                        </svg>
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
