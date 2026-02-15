"use client";

export type ItemTemplate = {
  type: string;
  label: string;
  width: number;
  height: number;
  color: string;
};

export const ITEM_TEMPLATES: ItemTemplate[] = [
  { type: "tent", label: "テント", width: 80, height: 80, color: "#2D5016" },
  { type: "tarp", label: "タープ", width: 120, height: 60, color: "#6B8E23" },
  { type: "car", label: "車", width: 60, height: 100, color: "#4A5568" },
  { type: "firepit", label: "焚き火台", width: 40, height: 40, color: "#D2691E" },
  { type: "table", label: "テーブル", width: 70, height: 40, color: "#8B4513" },
  { type: "chair", label: "チェア", width: 30, height: 30, color: "#A0522D" },
];

type Props = {
  onAdd: (template: ItemTemplate) => void;
};

export function LayoutItemPalette({ onAdd }: Props) {
  return (
    <div className="flex flex-wrap gap-2">
      {ITEM_TEMPLATES.map((t) => (
        <button
          key={t.type}
          onClick={() => onAdd(t)}
          className="flex items-center gap-1.5 rounded-lg border border-gray-200 bg-white px-3 py-1.5 text-sm hover:bg-gray-50 active:bg-gray-100"
        >
          <span
            className="inline-block h-3 w-3 rounded-sm"
            style={{ backgroundColor: t.color }}
          />
          {t.label}
        </button>
      ))}
    </div>
  );
}
