import { PageHeader } from "@/components/PageHeader";

export default function Home() {
  return (
    <div>
      <PageHeader
        title="スーパーキャンプアプリ"
        description="キャンプをもっと楽しく便利に"
      />
      <div className="mt-8 grid grid-cols-2 gap-4">
        <a
          href="/checklists"
          className="flex flex-col items-center gap-2 rounded-xl bg-white p-6 shadow-sm hover:shadow-md transition-shadow"
        >
          <span className="text-3xl">✅</span>
          <span className="text-sm font-medium text-gray-700">持ち物チェックリスト</span>
        </a>
        <a
          href="/layouts"
          className="flex flex-col items-center gap-2 rounded-xl bg-white p-6 shadow-sm hover:shadow-md transition-shadow"
        >
          <span className="text-3xl">🗺</span>
          <span className="text-sm font-medium text-gray-700">設営レイアウト</span>
        </a>
        <a
          href="/fire-logs"
          className="flex flex-col items-center gap-2 rounded-xl bg-white p-6 shadow-sm hover:shadow-md transition-shadow"
        >
          <span className="text-3xl">🔥</span>
          <span className="text-sm font-medium text-gray-700">焚き火ログ</span>
        </a>
        <a
          href="/meal-plans"
          className="flex flex-col items-center gap-2 rounded-xl bg-white p-6 shadow-sm hover:shadow-md transition-shadow"
        >
          <span className="text-3xl">🍳</span>
          <span className="text-sm font-medium text-gray-700">キャンプ飯プランナー</span>
        </a>
        <a
          href="/gears"
          className="flex flex-col items-center gap-2 rounded-xl bg-white p-6 shadow-sm hover:shadow-md transition-shadow"
        >
          <span className="text-3xl">🎒</span>
          <span className="text-sm font-medium text-gray-700">ギア管理</span>
        </a>
        <a
          href="/campsites"
          className="flex flex-col items-center gap-2 rounded-xl bg-white p-6 shadow-sm hover:shadow-md transition-shadow"
        >
          <span className="text-3xl">⛺</span>
          <span className="text-sm font-medium text-gray-700">キャンプ場</span>
        </a>
      </div>
    </div>
  );
}
