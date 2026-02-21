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
│   ├── layout.tsx              # ルートレイアウト（BottomTabBar含む、Geistフォント）
│   ├── page.tsx                # ホーム画面（4機能へのダッシュボード）
│   ├── globals.css             # グローバルCSS・テーマカラー
│   ├── checklists/
│   │   ├── page.tsx            # チェックリスト一覧（作成・削除・進捗表示）
│   │   └── [id]/page.tsx       # チェックリスト詳細（アイテム追加・チェック・数量管理・進捗バー）
│   ├── layouts/
│   │   ├── page.tsx            # レイアウト一覧（作成・削除・アイテム数表示）
│   │   └── [id]/page.tsx       # レイアウトエディタ（ドラッグ&ドロップ・リサイズ・回転・グリッド表示）
│   ├── fire-logs/
│   │   └── page.tsx            # 焚き火ログ（一覧+作成/編集フォーム）
│   └── meal-plans/
│       └── page.tsx            # キャンプ飯プランナー（一覧+作成/編集フォーム、食事タイプ別グループ表示）
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

## 画面一覧

| 画面 | パス | 説明 |
|------|------|------|
| ホーム | `/` | 4機能へのダッシュボード（カードUI） |
| チェックリスト一覧 | `/checklists` | チェックリスト作成・一覧・進捗表示・削除 |
| チェックリスト詳細 | `/checklists/[id]` | アイテムの追加・チェック・数量管理・削除、進捗バー |
| レイアウト一覧 | `/layouts` | レイアウト作成・一覧・削除 |
| レイアウトエディタ | `/layouts/[id]` | キャンバス上でアイテムをドラッグ&ドロップ・リサイズ・回転 |
| 焚き火ログ | `/fire-logs` | 焚き火記録の一覧表示と作成/編集フォーム |
| キャンプ飯プランナー | `/meal-plans` | 食事メニューの一覧表示と作成/編集フォーム |

## UI特徴

- モバイルファーストデザイン（max-width: 480px）
- 固定下部ナビゲーションバー（5タブ: ホーム・チェックリスト・レイアウト・焚き火ログ・キャンプ飯）
- キャンプテーマのカラースキーム（グリーン・ブラウン・オレンジ・クリーム）
- Geistフォント使用
