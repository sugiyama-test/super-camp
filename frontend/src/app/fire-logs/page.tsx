"use client";

import { useEffect, useState } from "react";
import { PageHeader } from "@/components/PageHeader";
import { useFireLogStore } from "@/stores/useFireLogStore";
import type { CreateFireLogData, FireLog } from "@/stores/useFireLogStore";

const WOOD_TYPES = ["広葉樹", "針葉樹", "薪（種類不明）", "炭", "その他"];

function formatDate(dateStr: string) {
  const d = new Date(dateStr);
  return `${d.getFullYear()}/${String(d.getMonth() + 1).padStart(2, "0")}/${String(d.getDate()).padStart(2, "0")}`;
}

function toInputDate(dateStr: string) {
  const d = new Date(dateStr);
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
}

const emptyForm: CreateFireLogData = {
  date: new Date().toISOString().slice(0, 10),
  location: "",
  wood_type: "",
  duration_minutes: 60,
  notes: "",
  temperature: null,
};

export default function FireLogsPage() {
  const { fireLogs, loading, fetchFireLogs, createFireLog, updateFireLog, deleteFireLog } =
    useFireLogStore();
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [form, setForm] = useState<CreateFireLogData>(emptyForm);

  useEffect(() => {
    fetchFireLogs();
  }, [fetchFireLogs]);

  const resetForm = () => {
    setForm(emptyForm);
    setEditingId(null);
    setShowForm(false);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.date) return;
    if (editingId) {
      await updateFireLog(editingId, form);
    } else {
      await createFireLog(form);
    }
    resetForm();
  };

  const handleEdit = (log: FireLog) => {
    setForm({
      date: toInputDate(log.date),
      location: log.location,
      wood_type: log.wood_type,
      duration_minutes: log.duration_minutes,
      notes: log.notes,
      temperature: log.temperature,
    });
    setEditingId(log.id);
    setShowForm(true);
  };

  const handleDelete = (id: number) => {
    if (confirm("この焚き火ログを削除しますか？")) {
      deleteFireLog(id);
    }
  };

  return (
    <div>
      <PageHeader
        title="焚き火ログ"
        description="焚き火の記録を残しましょう"
      />

      <div className="mt-6">
        {!showForm ? (
          <button
            onClick={() => { setForm(emptyForm); setEditingId(null); setShowForm(true); }}
            className="w-full rounded-lg bg-[var(--camp-orange)] px-4 py-2 text-sm font-medium text-white hover:opacity-90"
          >
            + 新しい焚き火ログ
          </button>
        ) : (
          <form onSubmit={handleSubmit} className="rounded-xl bg-white p-4 shadow-sm space-y-3">
            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="block text-xs text-gray-500 mb-1">日付 *</label>
                <input
                  type="date"
                  value={form.date}
                  onChange={(e) => setForm({ ...form, date: e.target.value })}
                  className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-orange)]"
                  required
                />
              </div>
              <div>
                <label className="block text-xs text-gray-500 mb-1">場所</label>
                <input
                  type="text"
                  value={form.location}
                  onChange={(e) => setForm({ ...form, location: e.target.value })}
                  placeholder="キャンプ場名..."
                  className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-orange)]"
                />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="block text-xs text-gray-500 mb-1">薪の種類</label>
                <select
                  value={form.wood_type}
                  onChange={(e) => setForm({ ...form, wood_type: e.target.value })}
                  className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-orange)]"
                >
                  <option value="">選択...</option>
                  {WOOD_TYPES.map((w) => (
                    <option key={w} value={w}>{w}</option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-xs text-gray-500 mb-1">時間（分）</label>
                <input
                  type="number"
                  value={form.duration_minutes}
                  onChange={(e) => setForm({ ...form, duration_minutes: Number(e.target.value) })}
                  min={0}
                  className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-orange)]"
                />
              </div>
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-1">気温（℃）</label>
              <input
                type="number"
                value={form.temperature ?? ""}
                onChange={(e) => setForm({ ...form, temperature: e.target.value ? Number(e.target.value) : null })}
                placeholder="任意"
                step="0.1"
                className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-orange)]"
              />
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-1">メモ</label>
              <textarea
                value={form.notes}
                onChange={(e) => setForm({ ...form, notes: e.target.value })}
                placeholder="焚き火の様子、薪の感想など..."
                rows={2}
                className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-orange)]"
              />
            </div>
            <div className="flex gap-2">
              <button
                type="submit"
                className="flex-1 rounded-lg bg-[var(--camp-orange)] px-4 py-2 text-sm font-medium text-white hover:opacity-90"
              >
                {editingId ? "更新" : "記録する"}
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
        {loading && fireLogs.length === 0 && (
          <p className="text-center text-gray-400 text-sm">読み込み中...</p>
        )}
        {!loading && fireLogs.length === 0 && (
          <p className="text-center text-gray-400 text-sm mt-8">
            焚き火ログがありません。上のボタンから記録しましょう！
          </p>
        )}
        {fireLogs.map((log) => (
          <div key={log.id} className="rounded-xl bg-white p-4 shadow-sm">
            <div className="flex items-start justify-between">
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2">
                  <span className="text-sm font-medium text-gray-800">
                    {formatDate(log.date)}
                  </span>
                  {log.location && (
                    <span className="text-xs text-gray-500">@ {log.location}</span>
                  )}
                </div>
                <div className="mt-1 flex flex-wrap gap-2">
                  {log.wood_type && (
                    <span className="inline-block rounded-full bg-orange-100 text-orange-700 px-2 py-0.5 text-xs">
                      {log.wood_type}
                    </span>
                  )}
                  <span className="inline-block rounded-full bg-gray-100 text-gray-600 px-2 py-0.5 text-xs">
                    {log.duration_minutes}分
                  </span>
                  {log.temperature != null && (
                    <span className="inline-block rounded-full bg-blue-100 text-blue-600 px-2 py-0.5 text-xs">
                      {log.temperature}℃
                    </span>
                  )}
                </div>
                {log.notes && (
                  <p className="mt-2 text-xs text-gray-600 whitespace-pre-wrap">{log.notes}</p>
                )}
              </div>
              <div className="flex gap-1 ml-2">
                <button
                  onClick={() => handleEdit(log)}
                  className="text-gray-400 hover:text-gray-600 p-1"
                  aria-label="編集"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                    <path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z" />
                  </svg>
                </button>
                <button
                  onClick={() => handleDelete(log.id)}
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
