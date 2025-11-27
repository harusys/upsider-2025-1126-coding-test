# スーパー支払い君.com API

請求書管理・支払い処理を行うREST APIサービス。

## 技術スタック

- **言語**: Go 1.25.4
- **フレームワーク**: gin
- **データベース**: PostgreSQL 18
- **クエリビルダー**: sqlc
- **マイグレーション**: sqldef (psqldef)
- **認証**: JWT
- **ロガー**: slog
- **APIドキュメント**: swag (OpenAPI自動生成)
- **テスト**: testify, gomock
- **リンター**: golangci-lint v2

## 環境変数

| 変数名 | 説明 | デフォルト | 必須 |
|--------|------|------------|------|
| `DB_HOST` | データベースホスト | `localhost` | |
| `DB_PORT` | データベースポート | `5432` | |
| `DB_USER` | データベースユーザー | - | ✓ |
| `DB_PASSWORD` | データベースパスワード | - | ✓ |
| `DB_NAME` | データベース名 | - | ✓ |
| `DB_SSLMODE` | SSL モード | `disable` | |
| `JWT_SECRET` | JWT署名用シークレット | - | ✓ |
| `PORT` | APIサーバーポート | `8080` | |

## セットアップ

### 開発ツールのインストール

```bash
make tools
```

### Docker で起動

```bash
# コンテナ起動
make docker-up

# マイグレーション実行
make migrate

# ログ確認
make docker-logs
```

### ローカルで起動

```bash
# PostgreSQL が起動している状態で
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=super_shiharai
export JWT_SECRET=your-secret-key

# マイグレーション
make migrate

# 起動
make run
```

## 開発

### コード生成

```bash
# モック + sqlc 生成
make generate

# モックのみ
make generate-mocks

# sqlcのみ
make generate-sqlc

# Swagger ドキュメント生成
make swagger
```

### フォーマット・リント

```bash
# フォーマット (golangci-lint v2)
make fmt

# リント
make lint
```

設定ファイル: [.golangci-lint.yml](.golangci-lint.yml)

### テスト

```bash
# テスト実行
make test

# カバレッジレポート付き
make test-coverage
```

### ビルド

```bash
make build
```

## API エンドポイント

### 認証

| メソッド | エンドポイント | 説明 |
|----------|----------------|------|
| POST | `/api/auth/register` | ユーザー登録 |
| POST | `/api/auth/login` | ログイン |
| POST | `/api/auth/refresh` | トークン更新 |

### 請求書

| メソッド | エンドポイント | 説明 | 認証 |
|----------|----------------|------|------|
| POST | `/api/invoices` | 請求書作成 | 必須 |
| GET | `/api/invoices` | 請求書一覧取得 | 必須 |

#### GET /api/invoices クエリパラメータ

| パラメータ | 説明 | 例 |
|------------|------|-----|
| `start_date` | 支払期日の開始日 | `2024-01-01` |
| `end_date` | 支払期日の終了日 | `2024-12-31` |

## API 使用例

### ユーザー登録

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "山田太郎",
    "company_name": "株式会社ABC"
  }'
```

### ログイン

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### 請求書作成

```bash
curl -X POST http://localhost:8080/api/invoices \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "vendor_id": 1,
    "vendor_bank_account_id": 1,
    "issue_date": "2024-01-15",
    "payment_amount": 10000,
    "due_date": "2024-02-15"
  }'
```

### 請求書一覧取得

```bash
curl "http://localhost:8080/api/invoices?start_date=2024-01-01&end_date=2024-12-31" \
  -H "Authorization: Bearer <token>"
```

## Swagger UI

サーバー起動後、以下のURLでAPI仕様を確認できます：

```
http://localhost:8080/swagger/index.html
```

## ディレクトリ構造

```
.
├── cmd/api/              # エントリーポイント
├── internal/
│   ├── domain/           # ドメイン層
│   │   ├── entity/       # エンティティ
│   │   ├── repository/   # リポジトリインターフェース
│   │   └── service/      # ドメインサービス
│   ├── usecase/          # ユースケース層
│   ├── infrastructure/   # インフラ層
│   │   ├── database/     # DB接続・sqlc
│   │   ├── persistence/  # リポジトリ実装
│   │   └── security/     # JWT
│   └── controller/       # コントローラー層
│       ├── auth/         # 認証ハンドラ
│       ├── invoice/      # 請求書ハンドラ
│       └── middleware/   # ミドルウェア
├── db/
│   ├── schema.sql        # スキーマ定義
│   ├── queries/          # sqlcクエリ
│   └── sqlc.yaml         # sqlc設定
├── docs/
│   ├── ASSIGNMENT.md         # 課題説明
│   ├── IMPLEMENTATION_PLAN.md # 実装計画
│   ├── CODING_RULES.md       # コーディングルール
│   └── swagger/              # 自動生成
└── tests/                    # 統合テスト
```

## ドキュメント

- [課題説明](docs/ASSIGNMENT.md)
- [実装計画](docs/IMPLEMENTATION_PLAN.md)
- [コーディングルール](docs/CODING_RULES.md)
