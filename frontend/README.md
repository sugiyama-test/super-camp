# Super Camp Frontend

Next.js (App Router) で実装されたフロントエンド。

## 技術

- Next.js 16 (App Router, Turbopack)
- TypeScript
- Tailwind CSS v4
- [Zustand](https://github.com/pmndrs/zustand) - 状態管理
- [react-konva](https://konvajs.org/docs/react/) - キャンバス描画（レイアウトエディタ）

## ディレクトリ構成

```
frontend/src/
├── app/
│   ├── layout.tsx              # ルートレイアウト（BottomTabBar含む）
│   ├── page.tsx                # ホーム画面
│   ├── globals.css             # グローバルCSS・テーマカラー
│   ├── checklists/
│   │   ├── page.tsx            # チェックリスト一覧
│   │   └── [id]/page.tsx       # チェックリスト詳細・編集
│   ├── layouts/
│   │   ├── page.tsx            # レイアウト一覧
│   │   └── [id]/page.tsx       # レイアウトエディタ（react-konva）
│   ├── fire-logs/
│   │   └── page.tsx            # 焚き火ログ（一覧+フォーム）
│   └── meal-plans/
│       └── page.tsx            # キャンプ飯プランナー（一覧+フォーム）
├── components/
│   ├── BottomTabBar.tsx        # 固定下部ナビゲーション
│   ├── PageHeader.tsx          # 共通ヘッダー
│   ├── ChecklistCard.tsx       # チェックリストカード
│   ├── ChecklistItemRow.tsx    # チェックリストアイテム行
│   ├── LayoutCanvas.tsx        # Konvaキャンバス
│   └── LayoutItemPalette.tsx   # レイアウトアイテム選択パレット
├── stores/
│   ├── useChecklistStore.ts    # チェックリストストア
│   ├── useLayoutStore.ts       # レイアウトストア
│   ├── useFireLogStore.ts      # 焚き火ログストア
│   ├── useMealPlanStore.ts     # キャンプ飯ストア
│   └── useHealthStore.ts       # ヘルスチェックストア
└── lib/
    └── api.ts                  # API通信ユーティリティ
```

## テーマカラー

| 変数 | 色 | 用途 |
|------|-----|------|
| `--camp-green` | #2D5016 | メインカラー（チェックリスト等） |
| `--camp-brown` | #8B4513 | キャンプ飯 |
| `--camp-orange` | #D2691E | 焚き火ログ |
| `--camp-cream` | #FAEBD7 | 背景色 |

## ローカル開発（Docker Compose外）

```bash
npm install
npm run dev    # http://localhost:3000
npm run build  # 本番ビルド
npm run lint   # ESLint
```

環境変数: `NEXT_PUBLIC_API_URL=http://localhost:8081`

## 状態管理パターン

各機能はZustandストアで管理。パターン:

```typescript
const useXxxStore = create<XxxState>((set, get) => ({
  items: [],
  loading: false,
  fetchItems: async () => { ... },
  createItem: async (data) => { ... },
  updateItem: async (id, data) => { ... },
  deleteItem: async (id) => { ... },
}));
```

API通信は `lib/api.ts` の `apiFetch` を使用（204レスポンス対応済み）。
