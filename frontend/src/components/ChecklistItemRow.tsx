"use client";

import type { ChecklistItem } from "@/stores/useChecklistStore";

type Props = {
  item: ChecklistItem;
  onToggle: (item: ChecklistItem) => void;
  onDelete: (itemId: number) => void;
};

export function ChecklistItemRow({ item, onToggle, onDelete }: Props) {
  return (
    <div className="flex items-center gap-3 py-3 border-b border-gray-100 last:border-0">
      <button
        onClick={() => onToggle(item)}
        className={`w-6 h-6 rounded-md border-2 flex items-center justify-center transition-colors
          ${item.is_checked
            ? "bg-[var(--camp-green)] border-[var(--camp-green)] text-white"
            : "border-gray-300 hover:border-[var(--camp-green)]"
          }`}
        aria-label={item.is_checked ? "チェック解除" : "チェック"}
      >
        {item.is_checked && (
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round">
            <polyline points="20 6 9 17 4 12" />
          </svg>
        )}
      </button>
      <span
        className={`flex-1 text-sm ${item.is_checked ? "line-through text-gray-400" : "text-gray-700"}`}
      >
        {item.name}
        {item.quantity > 1 && (
          <span className="ml-1 text-xs text-gray-400">x{item.quantity}</span>
        )}
      </span>
      <button
        onClick={() => onDelete(item.id)}
        className="text-gray-300 hover:text-red-500 p-1"
        aria-label="削除"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <line x1="18" y1="6" x2="6" y2="18" /><line x1="6" y1="6" x2="18" y2="18" />
        </svg>
      </button>
    </div>
  );
}
