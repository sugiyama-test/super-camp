# Super Camp Backend

Go (Golang) で実装されたREST APIサーバー。

## 技術

- Go 1.25
- [chi](https://github.com/go-chi/chi) - HTTPルーター
- [pgx/v5](https://github.com/jackc/pgx) - PostgreSQLドライバー
- [Air](https://github.com/air-verse/air) - ホットリロード
- [golang-migrate](https://github.com/golang-migrate/migrate) - DBマイグレーション

## ディレクトリ構成

```
backend/
├── cmd/api/main.go          # エントリーポイント
├── internal/
│   ├── config/config.go     # 環境変数読み込み
│   ├── handler/             # HTTPハンドラー
│   │   ├── health.go        # ヘルスチェック
│   │   ├── checklist.go     # チェックリストAPI
│   │   ├── layout.go        # レイアウトAPI
│   │   ├── firelog.go       # 焚き火ログAPI
│   │   └── mealplan.go      # キャンプ飯API
│   ├── model/models.go      # データモデル
│   ├── repository/          # DB操作
│   │   ├── checklist.go
│   │   ├── layout.go
│   │   ├── firelog.go
│   │   └── mealplan.go
│   └── router/router.go     # ルーティング定義
├── migrations/              # DBマイグレーションファイル
│   ├── 000001_init_schema.up.sql
│   ├── 000001_init_schema.down.sql
│   ├── 000002_seed_user.up.sql
│   └── 000002_seed_user.down.sql
├── Dockerfile
├── .air.toml                # Airホットリロード設定
├── go.mod
└── go.sum
```

## アーキテクチャ

```
Router → Handler → Repository → PostgreSQL
```

- **Handler**: HTTPリクエスト/レスポンスの処理、バリデーション
- **Repository**: SQL実行、データ永続化
- **Model**: データ構造体の定義

## ローカル開発（Docker Compose外）

```bash
# 環境変数
export DATABASE_URL="postgres://supercamp:supercamp@localhost:5434/supercamp_dev?sslmode=disable"
export PORT=8080

# ビルド＆実行
go run ./cmd/api

# テスト
go test ./...

# ビルド確認
go build ./...
```

## DB スキーマ

マイグレーションで管理。主要テーブル:

| テーブル | 説明 |
|---------|------|
| users | ユーザー |
| checklists | チェックリスト |
| checklist_items | チェックリストのアイテム |
| layouts | 設営レイアウト（data列にJSONB） |
| fire_logs | 焚き火ログ |
| meal_plans | キャンプ飯プラン |
| gears | ギア（将来用） |
| campsites | キャンプ場（将来用） |

## ミドルウェア

- `middleware.Logger` - リクエストログ
- `middleware.Recoverer` - パニックリカバリー
- `middleware.RequestID` - リクエストID付与
- `cors.Handler` - CORS（localhost:3000を許可）
