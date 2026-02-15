"use client";

import { useEffect, useState } from "react";
import { PageHeader } from "@/components/PageHeader";
import { useLayoutStore } from "@/stores/useLayoutStore";
import { useRouter } from "next/navigation";
import Link from "next/link";

export default function LayoutsPage() {
  const { layouts, loading, fetchLayouts, createLayout, deleteLayout } =
    useLayoutStore();
  const [title, setTitle] = useState("");
  const router = useRouter();

  useEffect(() => {
    fetchLayouts();
  }, [fetchLayouts]);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;
    const created = await createLayout(title.trim());
    setTitle("");
    if (created) {
      router.push(`/layouts/${created.id}`);
    }
  };

  const handleDelete = (id: number) => {
    if (confirm("このレイアウトを削除しますか？")) {
      deleteLayout(id);
    }
  };

  return (
    <div>
      <PageHeader
        title="設営レイアウト"
        description="テント・タープの配置を計画しましょう"
      />

      <form onSubmit={handleCreate} className="mt-6 flex gap-2">
        <input
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="新しいレイアウト名..."
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
        {loading && layouts.length === 0 && (
          <p className="text-center text-gray-400 text-sm">読み込み中...</p>
        )}
        {!loading && layouts.length === 0 && (
          <p className="text-center text-gray-400 text-sm mt-8">
            レイアウトがありません。上のフォームから作成しましょう！
          </p>
        )}
        {layouts.map((layout) => (
          <div
            key={layout.id}
            className="flex items-center gap-3 rounded-xl bg-white p-4 shadow-sm"
          >
            <Link href={`/layouts/${layout.id}`} className="flex-1 min-w-0">
              <h3 className="font-medium text-gray-800 truncate">
                {layout.title}
              </h3>
              <p className="text-xs text-gray-500 mt-1">
                アイテム数: {layout.item_count}
              </p>
            </Link>
            <button
              onClick={() => handleDelete(layout.id)}
              className="text-gray-400 hover:text-red-500 p-1"
              aria-label="削除"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="18"
                height="18"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <path d="M3 6h18" />
                <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" />
                <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
              </svg>
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}
