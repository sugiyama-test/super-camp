"use client";

import { useEffect, useState } from "react";
import { PageHeader } from "@/components/PageHeader";
import { ChecklistCard } from "@/components/ChecklistCard";
import { useChecklistStore } from "@/stores/useChecklistStore";
import { useRouter } from "next/navigation";

export default function ChecklistsPage() {
  const { checklists, loading, fetchChecklists, createChecklist, deleteChecklist } =
    useChecklistStore();
  const [title, setTitle] = useState("");
  const router = useRouter();

  useEffect(() => {
    fetchChecklists();
  }, [fetchChecklists]);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;
    const created = await createChecklist(title.trim());
    setTitle("");
    if (created) {
      router.push(`/checklists/${created.id}`);
    }
  };

  const handleDelete = (id: number) => {
    if (confirm("このチェックリストを削除しますか？")) {
      deleteChecklist(id);
    }
  };

  return (
    <div>
      <PageHeader
        title="持ち物チェックリスト"
        description="キャンプの持ち物を管理しましょう"
      />

      <form onSubmit={handleCreate} className="mt-6 flex gap-2">
        <input
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="新しいチェックリスト名..."
          className="flex-1 rounded-lg border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-green)] focus:border-transparent"
        />
        <button
          type="submit"
          disabled={!title.trim()}
          className="rounded-lg bg-[var(--camp-green)] px-4 py-2 text-sm font-medium text-white hover:opacity-90 disabled:opacity-50"
        >
          作成
        </button>
      </form>

      <div className="mt-6 space-y-3">
        {loading && checklists.length === 0 && (
          <p className="text-center text-gray-400 text-sm">読み込み中...</p>
        )}
        {!loading && checklists.length === 0 && (
          <p className="text-center text-gray-400 text-sm mt-8">
            チェックリストがありません。上のフォームから作成しましょう！
          </p>
        )}
        {checklists.map((checklist) => (
          <ChecklistCard
            key={checklist.id}
            checklist={checklist}
            onDelete={handleDelete}
          />
        ))}
      </div>
    </div>
  );
}
