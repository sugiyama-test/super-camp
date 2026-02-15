"use client";

import { useEffect, useState, useCallback, use } from "react";
import { useLayoutStore } from "@/stores/useLayoutStore";
import type { LayoutItem, LayoutData } from "@/stores/useLayoutStore";
import dynamic from "next/dynamic";

const LayoutCanvas = dynamic(() => import("@/components/LayoutCanvas").then((m) => m.LayoutCanvas), {
  ssr: false,
  loading: () => <div className="h-[400px] rounded-lg border border-gray-200 bg-white flex items-center justify-center text-gray-400 text-sm">キャンバス読み込み中...</div>,
});
import {
  LayoutItemPalette,
  type ItemTemplate,
} from "@/components/LayoutItemPalette";
import Link from "next/link";

const CANVAS_WIDTH = 600;
const CANVAS_HEIGHT = 400;

function parseLayoutData(dataStr: string): LayoutData {
  try {
    const parsed = JSON.parse(dataStr);
    return {
      width: parsed.width ?? CANVAS_WIDTH,
      height: parsed.height ?? CANVAS_HEIGHT,
      items: Array.isArray(parsed.items) ? parsed.items : [],
    };
  } catch {
    return { width: CANVAS_WIDTH, height: CANVAS_HEIGHT, items: [] };
  }
}

export default function LayoutEditorPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  const layoutId = Number(id);
  const { currentLayout, loading, fetchLayout, updateLayout } =
    useLayoutStore();
  const [items, setItems] = useState<LayoutItem[]>([]);
  const [title, setTitle] = useState("");
  const [saving, setSaving] = useState(false);
  const [dirty, setDirty] = useState(false);

  useEffect(() => {
    fetchLayout(layoutId);
  }, [fetchLayout, layoutId]);

  useEffect(() => {
    if (currentLayout) {
      const data = parseLayoutData(currentLayout.data);
      setItems(data.items);
      setTitle(currentLayout.title);
    }
  }, [currentLayout]);

  const handleItemsChange = useCallback((newItems: LayoutItem[]) => {
    setItems(newItems);
    setDirty(true);
  }, []);

  const handleAddItem = useCallback(
    (template: ItemTemplate) => {
      const newItem: LayoutItem = {
        id: `item-${Date.now()}`,
        type: template.type,
        label: template.label,
        x: Math.round(CANVAS_WIDTH / 2 - template.width / 2),
        y: Math.round(CANVAS_HEIGHT / 2 - template.height / 2),
        width: template.width,
        height: template.height,
        rotation: 0,
      };
      setItems((prev) => [...prev, newItem]);
      setDirty(true);
    },
    []
  );

  const handleSave = async () => {
    setSaving(true);
    const data: LayoutData = {
      width: CANVAS_WIDTH,
      height: CANVAS_HEIGHT,
      items,
    };
    await updateLayout(layoutId, title, data);
    setDirty(false);
    setSaving(false);
  };

  if (loading && !currentLayout) {
    return (
      <div className="text-center text-gray-400 text-sm mt-8">
        読み込み中...
      </div>
    );
  }

  if (!currentLayout) {
    return (
      <div className="text-center mt-8">
        <p className="text-gray-500">レイアウトが見つかりません</p>
        <Link href="/layouts" className="text-sm text-[var(--camp-green)] mt-2 inline-block">
          一覧に戻る
        </Link>
      </div>
    );
  }

  return (
    <div>
      <div className="flex items-center gap-2 mb-4">
        <Link
          href="/layouts"
          className="text-gray-500 hover:text-gray-700 p-1"
          aria-label="戻る"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="20"
            height="20"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <path d="M19 12H5" />
            <path d="M12 19l-7-7 7-7" />
          </svg>
        </Link>
        <input
          type="text"
          value={title}
          onChange={(e) => {
            setTitle(e.target.value);
            setDirty(true);
          }}
          className="flex-1 text-lg font-bold text-gray-800 bg-transparent border-b border-transparent focus:border-gray-300 focus:outline-none px-1 py-0.5"
        />
        <button
          onClick={handleSave}
          disabled={saving || !dirty}
          className="rounded-lg bg-[var(--camp-green)] px-4 py-1.5 text-sm font-medium text-white hover:opacity-90 disabled:opacity-50"
        >
          {saving ? "保存中..." : "保存"}
        </button>
      </div>

      <div className="mb-3">
        <p className="text-xs text-gray-500 mb-2">アイテムを追加:</p>
        <LayoutItemPalette onAdd={handleAddItem} />
      </div>

      <div className="overflow-x-auto">
        <LayoutCanvas
          items={items}
          width={CANVAS_WIDTH}
          height={CANVAS_HEIGHT}
          onItemsChange={handleItemsChange}
        />
      </div>

      <p className="text-xs text-gray-400 mt-2 text-center">
        アイテムをドラッグして移動 / 選択してリサイズ・回転
      </p>
    </div>
  );
}
