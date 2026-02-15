"use client";

import Link from "next/link";
import type { ChecklistSummary } from "@/stores/useChecklistStore";

type Props = {
  checklist: ChecklistSummary;
  onDelete: (id: number) => void;
};

export function ChecklistCard({ checklist, onDelete }: Props) {
  const progress =
    checklist.item_count > 0
      ? Math.round((checklist.checked_count / checklist.item_count) * 100)
      : 0;

  return (
    <div className="flex items-center gap-3 rounded-xl bg-white p-4 shadow-sm">
      <Link
        href={`/checklists/${checklist.id}`}
        className="flex-1 min-w-0"
      >
        <h3 className="font-medium text-gray-800 truncate">{checklist.title}</h3>
        <div className="mt-1 flex items-center gap-2">
          <div className="h-2 flex-1 rounded-full bg-gray-100 overflow-hidden">
            <div
              className="h-full rounded-full bg-[var(--camp-green)] transition-all"
              style={{ width: `${progress}%` }}
            />
          </div>
          <span className="text-xs text-gray-500 whitespace-nowrap">
            {checklist.checked_count}/{checklist.item_count}
          </span>
        </div>
      </Link>
      <button
        onClick={() => onDelete(checklist.id)}
        className="text-gray-400 hover:text-red-500 p-1"
        aria-label="削除"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <path d="M3 6h18" /><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" /><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
        </svg>
      </button>
    </div>
  );
}
