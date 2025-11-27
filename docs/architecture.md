# アーキテクチャ設計

## 概要

本プロジェクトは **Clean Architecture** を採用し、関心の分離と依存性の方向を明確にしています。

## レイヤー構成

```
┌──────────────────────────────────────────────────────────┐
│                    Controller Layer                       │
│              (HTTP ハンドラ、ミドルウェア)                  │
└──────────────────────────────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────┐
│                     Usecase Layer                         │
│                 (ビジネスロジック)                          │
└──────────────────────────────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────┐
│                     Domain Layer                          │
│           (エンティティ、リポジトリIF、サービス)             │
└──────────────────────────────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────┐
│                  Infrastructure Layer                     │
│            (DB接続、リポジトリ実装、JWT)                    │
└──────────────────────────────────────────────────────────┘
```

## 依存性の方向

```
Controller → Usecase → Domain ← Infrastructure
```

- **Domain層** は他の層に依存しない（最も内側）
- **Infrastructure層** は Domain層のインターフェースを実装
- **Usecase層** は Domain層のインターフェースに依存
- **Controller層** は Usecase層のインターフェースに依存

## 各層の責務

### Domain Layer (`internal/domain/`)

最も内側のレイヤー。ビジネスルールとエンティティを定義。

| パッケージ | 責務 |
|-----------|------|
| `entity/` | ビジネスエンティティ（User, Invoice, Vendor等） |
| `repository/` | リポジトリインターフェース定義 |
| `service/` | ドメインサービス（手数料計算等） |

```go
// entity/invoice.go - ビジネスルールをエンティティに持たせる
type Invoice struct {
    ID            int64
    PaymentAmount int64
    Fee           int64
    Tax           int64
    TotalAmount   int64
    Status        InvoiceStatus
    // ...
}
```

### Usecase Layer (`internal/usecase/`)

アプリケーション固有のビジネスロジック。

| パッケージ | 責務 |
|-----------|------|
| `auth/` | 認証ユースケース（登録、ログイン、トークン更新） |
| `invoice/` | 請求書ユースケース（作成、一覧、詳細） |

```go
// usecase/invoice/usecase.go
type Usecase interface {
    Create(ctx context.Context, input *CreateInput) (*entity.Invoice, error)
    List(ctx context.Context, input *ListInput) ([]*entity.Invoice, error)
    GetByID(ctx context.Context, companyID, invoiceID int64) (*entity.Invoice, error)
}
```

### Infrastructure Layer (`internal/infrastructure/`)

外部システムとの接続を担当。

| パッケージ | 責務 |
|-----------|------|
| `database/` | PostgreSQL接続、sqlc生成コード |
| `persistence/` | リポジトリ実装 |
| `security/` | JWT サービス |

```go
// persistence/invoice_repository.go - Domain層のIFを実装
type InvoiceRepository struct {
    pool    *pgxpool.Pool
    queries *database.Queries
}

func (r *InvoiceRepository) Create(ctx context.Context, inv *entity.Invoice) error {
    // sqlcを使用したDB操作
}
```

### Controller Layer (`internal/controller/`)

HTTP リクエスト/レスポンスの処理。

| パッケージ | 責務 |
|-----------|------|
| `auth/` | 認証エンドポイント |
| `invoice/` | 請求書エンドポイント |
| `middleware/` | 認証ミドルウェア、エラーハンドリング |

```go
// controller/invoice/handler.go
func (h *Handler) Create(c *gin.Context) {
    // 1. リクエストバリデーション
    // 2. ユースケース呼び出し
    // 3. レスポンス変換
}
```

## データフロー

### 請求書作成の例

```
1. HTTP Request
       ↓
2. Controller: リクエストバリデーション、DTOからInput変換
       ↓
3. Usecase: ビジネスロジック実行
   - 取引先存在確認
   - 手数料計算（Domain Service）
   - 請求書エンティティ作成
       ↓
4. Repository: DB永続化
       ↓
5. Controller: エンティティからResponse変換
       ↓
6. HTTP Response
```

## 認証フロー

```
┌─────────────────────────────────────────────────────────────┐
│                      認証フロー                              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Register/Login → JWT発行 → Access Token + Refresh Token   │
│                                                             │
│  Protected API → Auth Middleware → Token検証 → Handler      │
│                                                             │
│  Token期限切れ → Refresh Token → 新Token発行                 │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### トークン仕様

| トークン | 有効期限 | 用途 |
|---------|---------|------|
| Access Token | 15分 | API認証 |
| Refresh Token | 7日 | Access Token更新 |

## エラーハンドリング

### レイヤー別エラー処理

| レイヤー | 責務 |
|---------|------|
| Domain | ビジネスエラー定義（ErrNotFound等） |
| Usecase | ビジネスエラー発生・変換 |
| Controller | HTTPステータスコードへの変換 |
| Middleware | 予期せぬエラーのキャッチ |

```go
// domain/errors.go
var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
)

// controller - エラーをHTTPステータスに変換
if errors.Is(err, domain.ErrNotFound) {
    c.JSON(http.StatusNotFound, ...)
}
```

## テスト戦略

### テストピラミッド

```
        /\
       /  \
      / E2E \        tests/api_test.go
     /--------\
    /          \
   / Integration \   usecase/*_test.go
  /--------------\
 /                \
/      Unit        \ domain/service/*_test.go
                     controller/*_test.go
```

| テスト種別 | 対象 | ツール |
|-----------|------|--------|
| Unit | Domain Service, Handler | testify, gomock |
| Integration | Usecase | testify, gomock |
| E2E | API全体 | testify/suite |

## セキュリティ対策

| 対策 | 実装 |
|------|------|
| パスワードハッシュ化 | bcrypt |
| SQL Injection防止 | sqlc（プレースホルダー） |
| JWT署名 | HS256 |
| 認可 | company_id によるマルチテナント分離 |

## パフォーマンス考慮

| 項目 | 対策 |
|------|------|
| コネクションプーリング | pgxpool使用 |
| N+1問題 | sqlcでJOIN使用 |
| インデックス | 検索条件カラムにINDEX |
