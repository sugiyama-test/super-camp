"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const tabs = [
  { href: "/", label: "ホーム", icon: "🏕" },
  { href: "/checklists", label: "持ち物", icon: "✅" },
  { href: "/layouts", label: "レイアウト", icon: "🗺" },
  { href: "/fire-logs", label: "焚き火", icon: "🔥" },
  { href: "/meal-plans", label: "キャンプ飯", icon: "🍳" },
];

export function BottomTabBar() {
  const pathname = usePathname();

  return (
    <nav className="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 z-50">
      <div className="mx-auto max-w-lg flex justify-around">
        {tabs.map((tab) => {
          const isActive = pathname === tab.href;
          return (
            <Link
              key={tab.href}
              href={tab.href}
              className={`flex flex-col items-center py-2 px-3 text-xs
                ${isActive ? "text-[var(--camp-green)] font-bold" : "text-gray-500"}`}
            >
              <span className="text-xl">{tab.icon}</span>
              <span>{tab.label}</span>
            </Link>
          );
        })}
      </div>
    </nav>
  );
}
