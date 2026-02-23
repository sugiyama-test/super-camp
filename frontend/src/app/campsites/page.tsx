"use client";

import { useEffect, useState } from "react";
import { PageHeader } from "@/components/PageHeader";
import { useCampsiteStore } from "@/stores/useCampsiteStore";
import type { CreateCampsiteData, Campsite } from "@/stores/useCampsiteStore";

const emptyForm: CreateCampsiteData = {
  name: "",
  address: "",
  latitude: null,
  longitude: null,
  notes: "",
};

export default function CampsitesPage() {
  const { campsites, loading, fetchCampsites, createCampsite, updateCampsite, deleteCampsite } =
    useCampsiteStore();
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [form, setForm] = useState<CreateCampsiteData>(emptyForm);

  useEffect(() => {
    fetchCampsites();
  }, [fetchCampsites]);

  const resetForm = () => {
    setForm(emptyForm);
    setEditingId(null);
    setShowForm(false);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.name.trim()) return;
    if (editingId) {
      await updateCampsite(editingId, form);
    } else {
      await createCampsite(form);
    }
    resetForm();
  };

  const handleEdit = (campsite: Campsite) => {
    setForm({
      name: campsite.name,
      address: campsite.address,
      latitude: campsite.latitude,
      longitude: campsite.longitude,
      notes: campsite.notes,
    });
    setEditingId(campsite.id);
    setShowForm(true);
  };

  const handleDelete = (id: number) => {
    if (confirm("このキャンプ場を削除しますか？")) {
      deleteCampsite(id);
    }
  };

  return (
    <div>
      <PageHeader
        title="キャンプ場"
        description="お気に入りのキャンプ場を管理しましょう"
      />

      <div className="mt-6">
        {!showForm ? (
          <button
            onClick={() => { setForm(emptyForm); setEditingId(null); setShowForm(true); }}
            className="w-full rounded-lg bg-[var(--camp-orange)] px-4 py-2 text-sm font-medium text-white hover:opacity-90"
          >
            + 新しいキャンプ場
          </button>
        ) : (
          <form onSubmit={handleSubmit} className="rounded-xl bg-white p-4 shadow-sm space-y-3">
            <div>
              <label className="block text-xs text-gray-500 mb-1">名前 *</label>
              <input
                type="text"
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                placeholder="キャンプ場名..."
                className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-orange)]"
                required
              />
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-1">住所</label>
              <input
                type="text"
                value={form.address}
                onChange={(e) => setForm({ ...form, address: e.target.value })}
                placeholder="住所を入力..."
                className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-orange)]"
              />
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="block text-xs text-gray-500 mb-1">緯度</label>
                <input
                  type="number"
                  value={form.latitude ?? ""}
                  onChange={(e) => setForm({ ...form, latitude: e.target.value ? Number(e.target.value) : null })}
                  placeholder="35.6762"
                  step="any"
                  className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-orange)]"
                />
              </div>
              <div>
                <label className="block text-xs text-gray-500 mb-1">経度</label>
                <input
                  type="number"
                  value={form.longitude ?? ""}
                  onChange={(e) => setForm({ ...form, longitude: e.target.value ? Number(e.target.value) : null })}
                  placeholder="139.6503"
                  step="any"
                  className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-orange)]"
                />
              </div>
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-1">メモ</label>
              <textarea
                value={form.notes}
                onChange={(e) => setForm({ ...form, notes: e.target.value })}
                placeholder="キャンプ場の特徴、設備など..."
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
        {loading && campsites.length === 0 && (
          <p className="text-center text-gray-400 text-sm">読み込み中...</p>
        )}
        {!loading && campsites.length === 0 && (
          <p className="text-center text-gray-400 text-sm mt-8">
            キャンプ場が登録されていません。上のボタンから追加しましょう！
          </p>
        )}
        {campsites.map((campsite) => (
          <div key={campsite.id} className="rounded-xl bg-white p-4 shadow-sm">
            <div className="flex items-start justify-between">
              <div className="flex-1 min-w-0">
                <span className="text-sm font-medium text-gray-800">
                  {campsite.name}
                </span>
                {campsite.address && (
                  <p className="mt-0.5 text-xs text-gray-500">{campsite.address}</p>
                )}
                <div className="mt-1 flex flex-wrap gap-2">
                  {campsite.latitude != null && campsite.longitude != null && (
                    <span className="inline-block rounded-full bg-blue-100 text-blue-600 px-2 py-0.5 text-xs">
                      {campsite.latitude.toFixed(4)}, {campsite.longitude.toFixed(4)}
                    </span>
                  )}
                </div>
                {campsite.notes && (
                  <p className="mt-2 text-xs text-gray-600 whitespace-pre-wrap">{campsite.notes}</p>
                )}
              </div>
              <div className="flex gap-1 ml-2">
                <button
                  onClick={() => handleEdit(campsite)}
                  className="text-gray-400 hover:text-gray-600 p-1"
                  aria-label="編集"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                    <path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z" />
                  </svg>
                </button>
                <button
                  onClick={() => handleDelete(campsite.id)}
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
