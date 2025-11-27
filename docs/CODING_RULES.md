# コーディングルール

## 1. プロジェクト構造

### Clean Architecture
```
internal/
├── domain/          # ドメイン層（ビジネスロジック）
├── usecase/         # ユースケース層（アプリケーションロジック）
├── infrastructure/  # インフラ層（DB、外部サービス）
└── controller/      # コントローラー層（HTTP）
```

### 依存関係の方向
- controller → usecase → domain
- infrastructure → domain
- domain は他の層に依存しない

## 2. 命名規則

### パッケージ名
- 小文字、単数形
- 短く明確に（`auth`, `invoice`, `user`）

### ファイル名
- snake_case（`invoice_repository.go`）
- テストファイルは `_test.go` サフィックス

### 構造体・インターフェース
- PascalCase
- インターフェースは動詞 + er（`InvoiceRepository`, `Authenticator`）
- 実装は具体的な名前（`PostgresInvoiceRepository`）

### 関数・メソッド
- PascalCase（公開）/ camelCase（非公開）
- 動詞から始める（`CreateInvoice`, `findByID`）

### 変数
- camelCase
- 短くても意味が明確に（`inv` より `invoice`）
- ループカウンタは `i, j, k` 可

### 定数
- PascalCase（公開）
- グループ化する場合は `const ( ... )`

## 3. エラーハンドリング

### カスタムエラー
```go
type AppError struct {
    Code    string
    Message string
    Err     error
}
```

### エラーラップ
```go
if err != nil {
    return fmt.Errorf("failed to create invoice: %w", err)
}
```

### パニックしない
- パニックは初期化時のみ
- それ以外はエラーを返す

## 4. ロギング

### slog使用
```go
slog.Info("invoice created",
    slog.Int64("id", invoice.ID),
    slog.Int64("company_id", invoice.CompanyID),
)
```

### ログレベル
- `Debug`: 開発時の詳細情報
- `Info`: 正常な処理の記録
- `Warn`: 注意が必要な状況
- `Error`: エラー発生時

### 機密情報
- パスワード、トークンはログに出力しない

## 5. テスト

### テストファイル配置
- 同じパッケージ内に `_test.go`

### テスト関数名
```go
func TestCreateInvoice_Success(t *testing.T)
func TestCreateInvoice_InvalidAmount(t *testing.T)
```

### テーブルドリブンテスト
```go
tests := []struct {
    name    string
    input   CreateInvoiceInput
    want    *Invoice
    wantErr bool
}{
    // ...
}
```

### モック
- gomock使用
- `//go:generate mockgen` でモック生成

## 6. API設計

### エンドポイント
- RESTful
- 小文字、ハイフン区切り（`/api/invoices`）

### HTTPメソッド
- GET: 取得
- POST: 作成
- PUT: 全体更新
- PATCH: 部分更新
- DELETE: 削除

### 正常レスポンス
単一取得時：
```json
{
  "id": 1,
  "company_name": "株式会社ABC",
  "amount": 10440,
  "created_at": "2024-01-01T00:00:00Z"
}
```

複数取得時：
```json
[
  { "id": 1, "amount": 10440, ... },
  { "id": 2, "amount": 20880, ... }
]
```

### エラーレスポンス
```json
{
  "error": {
    "code": "INVALID_INPUT",
    "message": "支払金額は1以上である必要があります"
  }
}
```

## 7. データベース

### テーブル名
- 小文字、snake_case、複数形（`invoices`, `vendor_bank_accounts`）

### カラム名
- 小文字、snake_case
- 主キーは `id`
- 外部キーは `{table}_id`
- タイムスタンプは `created_at`, `updated_at`

### インデックス
- 検索条件に使うカラムにはインデックス
- 命名: `idx_{table}_{column}`

## 8. セキュリティ

### パスワード
- bcryptでハッシュ化（cost: 12）

### JWT
- 適切な有効期限設定
- シークレットは環境変数から取得

### SQL Injection
- sqlcのプレースホルダー使用（自動的に防御）

### バリデーション
- 入力は必ずバリデーション
- `go-playground/validator` 使用

## 9. コード品質

### golangci-lint v2
フォーマットとリントは `golangci-lint v2` で統一管理。

設定ファイル: [.golangci-lint.yml](../.golangci-lint.yml)

```bash
# フォーマット
golangci-lint fmt ./...

# リント
golangci-lint run ./...
```

### 有効なフォーマッター
- `gofmt`: 標準フォーマット
- `goimports`: import整理
- `gci`: import順序
- `golines`: 行長制限
- `swaggo`: Swagger コメント整形

### 主要なリンター
- `gosec`: セキュリティチェック
- `gocritic`: コード品質
- `gocyclo` / `cyclop`: 複雑度
- `err113`: エラーハンドリング
- `testifylint`: テストコード品質

### コメント
- 公開APIにはGoDoc形式のコメント
- 複雑なロジックには説明コメント

## 10. Git

### コミットメッセージ
```
<type>: <description>

<body>
```

### タイプ
- `feat`: 新機能
- `fix`: バグ修正
- `refactor`: リファクタリング
- `docs`: ドキュメント
- `test`: テスト
- `chore`: その他

### コミット粒度
- 1つの論理的な変更 = 1コミット
- レビューしやすいサイズ
