"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { PageHeader } from "@/components/PageHeader";
import { ChecklistItemRow } from "@/components/ChecklistItemRow";
import { useChecklistStore, type ChecklistItem } from "@/stores/useChecklistStore";

export default function ChecklistDetailPage() {
  const params = useParams();
  const router = useRouter();
  const id = Number(params.id);
  const { currentChecklist, loading, fetchChecklist, addItem, updateItem, deleteItem } =
    useChecklistStore();
  const [newItemName, setNewItemName] = useState("");

  useEffect(() => {
    if (id) fetchChecklist(id);
  }, [id, fetchChecklist]);

  const handleAddItem = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newItemName.trim()) return;
    await addItem(id, newItemName.trim());
    setNewItemName("");
  };

  const handleToggle = (item: ChecklistItem) => {
    updateItem(id, item.id, {
      name: item.name,
      is_checked: !item.is_checked,
      quantity: item.quantity,
    });
  };

  const handleDeleteItem = (itemId: number) => {
    deleteItem(id, itemId);
  };

  if (loading && !currentChecklist) {
    return (
      <div className="text-center text-gray-400 text-sm mt-8">読み込み中...</div>
    );
  }

  if (!currentChecklist) {
    return (
      <div className="text-center text-gray-400 text-sm mt-8">
        チェックリストが見つかりません
      </div>
    );
  }

  const checkedCount = currentChecklist.items.filter((i) => i.is_checked).length;
  const totalCount = currentChecklist.items.length;

  return (
    <div>
      <button
        onClick={() => router.push("/checklists")}
        className="text-sm text-gray-500 hover:text-gray-700 mb-2"
      >
        &larr; 一覧に戻る
      </button>

      <PageHeader title={currentChecklist.title} />

      {totalCount > 0 && (
        <div className="mt-3 flex items-center gap-2">
          <div className="h-2 flex-1 rounded-full bg-gray-100 overflow-hidden">
            <div
              className="h-full rounded-full bg-[var(--camp-green)] transition-all"
              style={{ width: `${Math.round((checkedCount / totalCount) * 100)}%` }}
            />
          </div>
          <span className="text-xs text-gray-500">
            {checkedCount}/{totalCount}
          </span>
        </div>
      )}

      <form onSubmit={handleAddItem} className="mt-6 flex gap-2">
        <input
          type="text"
          value={newItemName}
          onChange={(e) => setNewItemName(e.target.value)}
          placeholder="アイテムを追加..."
          className="flex-1 rounded-lg border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-green)] focus:border-transparent"
        />
        <button
          type="submit"
          disabled={!newItemName.trim()}
          className="rounded-lg bg-[var(--camp-green)] px-4 py-2 text-sm font-medium text-white hover:opacity-90 disabled:opacity-50"
        >
          追加
        </button>
      </form>

      <div className="mt-4 rounded-xl bg-white p-4 shadow-sm">
        {currentChecklist.items.length === 0 ? (
          <p className="text-center text-gray-400 text-sm py-4">
            アイテムがありません。上のフォームから追加しましょう！
          </p>
        ) : (
          currentChecklist.items.map((item) => (
            <ChecklistItemRow
              key={item.id}
              item={item}
              onToggle={handleToggle}
              onDelete={handleDeleteItem}
            />
          ))
        )}
      </div>
    </div>
  );
}
