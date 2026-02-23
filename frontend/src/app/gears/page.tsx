"use client";

import { useEffect, useState } from "react";
import { PageHeader } from "@/components/PageHeader";
import { useGearStore } from "@/stores/useGearStore";
import type { CreateGearData, Gear } from "@/stores/useGearStore";

const CATEGORIES = ["シェルター", "寝具", "調理器具", "焚き火", "照明", "ファニチャー", "ツール", "その他"];

const emptyForm: CreateGearData = {
  name: "",
  category: "",
  brand: "",
  weight_grams: null,
  notes: "",
};

export default function GearsPage() {
  const { gears, loading, fetchGears, createGear, updateGear, deleteGear } =
    useGearStore();
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [form, setForm] = useState<CreateGearData>(emptyForm);

  useEffect(() => {
    fetchGears();
  }, [fetchGears]);

  const resetForm = () => {
    setForm(emptyForm);
    setEditingId(null);
    setShowForm(false);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.name.trim()) return;
    if (editingId) {
      await updateGear(editingId, form);
    } else {
      await createGear(form);
    }
    resetForm();
  };

  const handleEdit = (gear: Gear) => {
    setForm({
      name: gear.name,
      category: gear.category,
      brand: gear.brand,
      weight_grams: gear.weight_grams,
      notes: gear.notes,
    });
    setEditingId(gear.id);
    setShowForm(true);
  };

  const handleDelete = (id: number) => {
    if (confirm("このギアを削除しますか？")) {
      deleteGear(id);
    }
  };

  return (
    <div>
      <PageHeader
        title="ギア管理"
        description="キャンプギアを管理しましょう"
      />

      <div className="mt-6">
        {!showForm ? (
          <button
            onClick={() => { setForm(emptyForm); setEditingId(null); setShowForm(true); }}
            className="w-full rounded-lg bg-[var(--camp-orange)] px-4 py-2 text-sm font-medium text-white hover:opacity-90"
          >
            + 新しいギア
          </button>
        ) : (
          <form onSubmit={handleSubmit} className="rounded-xl bg-white p-4 shadow-sm space-y-3">
            <div>
              <label className="block text-xs text-gray-500 mb-1">名前 *</label>
              <input
                type="text"
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                placeholder="テント、タープ、寝袋..."
                className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-orange)]"
                required
              />
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="block text-xs text-gray-500 mb-1">カテゴリ</label>
                <select
                  value={form.category}
                  onChange={(e) => setForm({ ...form, category: e.target.value })}
                  className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-orange)]"
                >
                  <option value="">選択...</option>
                  {CATEGORIES.map((c) => (
                    <option key={c} value={c}>{c}</option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-xs text-gray-500 mb-1">ブランド</label>
                <input
                  type="text"
                  value={form.brand}
                  onChange={(e) => setForm({ ...form, brand: e.target.value })}
                  placeholder="Snow Peak..."
                  className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-orange)]"
                />
              </div>
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-1">重量（g）</label>
              <input
                type="number"
                value={form.weight_grams ?? ""}
                onChange={(e) => setForm({ ...form, weight_grams: e.target.value ? Number(e.target.value) : null })}
                placeholder="任意"
                min={0}
                className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-orange)]"
              />
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-1">メモ</label>
              <textarea
                value={form.notes}
                onChange={(e) => setForm({ ...form, notes: e.target.value })}
                placeholder="特徴、使用感など..."
                rows={2}
                className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-orange)]"
              />
            </div>
            <div className="flex gap-2">
              <button
                type="submit"
                className="flex-1 rounded-lg bg-[var(--camp-orange)] px-4 py-2 text-sm font-medium text-white hover:opacity-90"
              >
                {editingId ? "更新" : "登録する"}
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

      <div className="mt-6 space-y-3">
        {loading && gears.length === 0 && (
          <p className="text-center text-gray-400 text-sm">読み込み中...</p>
        )}
        {!loading && gears.length === 0 && (
          <p className="text-center text-gray-400 text-sm mt-8">
            ギアが登録されていません。上のボタンから追加しましょう！
          </p>
        )}
        {gears.map((gear) => (
          <div key={gear.id} className="rounded-xl bg-white p-4 shadow-sm">
            <div className="flex items-start justify-between">
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2">
                  <span className="text-sm font-medium text-gray-800">
                    {gear.name}
                  </span>
                  {gear.brand && (
                    <span className="text-xs text-gray-500">{gear.brand}</span>
                  )}
                </div>
                <div className="mt-1 flex flex-wrap gap-2">
                  {gear.category && (
                    <span className="inline-block rounded-full bg-green-100 text-green-700 px-2 py-0.5 text-xs">
                      {gear.category}
                    </span>
                  )}
                  {gear.weight_grams != null && (
                    <span className="inline-block rounded-full bg-gray-100 text-gray-600 px-2 py-0.5 text-xs">
                      {gear.weight_grams >= 1000
                        ? `${(gear.weight_grams / 1000).toFixed(1)}kg`
                        : `${gear.weight_grams}g`}
                    </span>
                  )}
                </div>
                {gear.notes && (
                  <p className="mt-2 text-xs text-gray-600 whitespace-pre-wrap">{gear.notes}</p>
                )}
              </div>
              <div className="flex gap-1 ml-2">
                <button
                  onClick={() => handleEdit(gear)}
                  className="text-gray-400 hover:text-gray-600 p-1"
                  aria-label="編集"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                    <path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z" />
                  </svg>
                </button>
                <button
                  onClick={() => handleDelete(gear.id)}
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
  );
}
