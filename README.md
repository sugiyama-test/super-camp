# Super Camp App

キャンプの準備・計画・記録・振り返りをサポートするWebアプリケーション。
ソロキャンパーからファミリーキャンプまで対応。

## 機能一覧

| 機能 | 説明 |
|------|------|
| 持ち物チェックリスト | キャンプの持ち物を管理。アイテムの追加・チェック・削除、進捗バー表示 |
| 設営レイアウト | テント・タープ・車・焚き火台等をキャンバス上にドラッグ&ドロップで配置。JSON形式で保存・復元 |
| 焚き火ログ | 焚き火の記録（日付・場所・薪の種類・時間・気温・メモ） |
| キャンプ飯プランナー | 食事メニューの計画（朝食/昼食/夕食/おやつ、人数、レシピメモ） |

## 技術スタック

| レイヤー | 技術 |
|---------|------|
| フロントエンド | Next.js 16 (App Router), TypeScript, Tailwind CSS v4, Zustand, react-konva |
| バックエンド | Go (Golang), chi router, pgx/v5 |
| データベース | PostgreSQL 16 |
| インフラ | Docker Compose |

## プロジェクト構成

```
super-camp/
├── frontend/          # Next.js フロントエンド
├── backend/           # Go バックエンド
├── docker-compose.yml # 開発環境定義（DB, API, Frontend, テストDB）
├── Makefile           # 開発用コマンド
└── .env.example       # 環境変数テンプレート
```

## セットアップ

### 前提条件
- Docker & Docker Compose

### 起動

```bash
# 全サービス起動（DB + API + Frontend）
make up-build

# マイグレーション実行
make migrate-up

# ヘルスチェック
make health
```

### アクセス

| サービス | URL |
|---------|-----|
| フロントエンド | http://localhost:3000 |
| API | http://localhost:8081/api |
| PostgreSQL | localhost:5434 (user: supercamp / pass: supercamp / db: supercamp_dev) |

### よく使うコマンド

```bash
# Docker
make up              # サービス起動
make up-build        # ビルドして起動
make down            # サービス停止
make logs            # 全ログ表示
make logs-api        # APIログのみ
make logs-frontend   # フロントエンドログのみ

# マイグレーション
make migrate-up      # マイグレーション適用
make migrate-down    # マイグレーション1つ戻す

# テスト
make test-db-up      # テスト用DB起動
make test-db-migrate # テスト用DBにマイグレーション適用
make api-test        # バックエンド全テスト実行
make api-test-unit   # ハンドラーのユニットテスト
make api-test-integration  # リポジトリの結合テスト
make frontend-lint   # フロントエンドLint
make frontend-build  # フロントエンドビルド

# その他
make health          # ヘルスチェック
```

## API エンドポイント

### チェックリスト
| Method | Path | 説明 |
|--------|------|------|
| GET | `/api/checklists` | 一覧 |
| POST | `/api/checklists` | 作成 |
| GET | `/api/checklists/{id}` | 詳細（アイテム付き） |
| DELETE | `/api/checklists/{id}` | 削除 |
| POST | `/api/checklists/{id}/items` | アイテム追加 |
| PUT | `/api/checklists/{id}/items/{itemID}` | アイテム更新 |
| DELETE | `/api/checklists/{id}/items/{itemID}` | アイテム削除 |

### レイアウト
| Method | Path | 説明 |
|--------|------|------|
| GET | `/api/layouts` | 一覧 |
| POST | `/api/layouts` | 作成 |
| GET | `/api/layouts/{id}` | 詳細 |
| PUT | `/api/layouts/{id}` | 更新 |
| DELETE | `/api/layouts/{id}` | 削除 |

### 焚き火ログ
| Method | Path | 説明 |
|--------|------|------|
| GET | `/api/fire-logs` | 一覧 |
| POST | `/api/fire-logs` | 作成 |
| GET | `/api/fire-logs/{id}` | 詳細 |
| PUT | `/api/fire-logs/{id}` | 更新 |
| DELETE | `/api/fire-logs/{id}` | 削除 |

### キャンプ飯プランナー
| Method | Path | 説明 |
|--------|------|------|
| GET | `/api/meal-plans` | 一覧 |
| POST | `/api/meal-plans` | 作成 |
| GET | `/api/meal-plans/{id}` | 詳細 |
| PUT | `/api/meal-plans/{id}` | 更新 |
| DELETE | `/api/meal-plans/{id}` | 削除 |

## テスト

バックエンドにはハンドラーとリポジトリの両方にテストを実装済み。

- **ユニットテスト**: ハンドラー層のテスト（モックリポジトリ使用）
- **結合テスト**: リポジトリ層のテスト（テスト用PostgreSQLコンテナ使用、ポート5435）

```bash
# テスト用DB起動 → マイグレーション → テスト実行
make test-db-up
make test-db-migrate
make api-test
```

## 備考

- 認証機能は未実装。現在は固定ユーザー（user_id=1）で動作
- DBスキーマには `users`, `gears`, `campsites` テーブルも定義済み（将来の機能拡張用）
