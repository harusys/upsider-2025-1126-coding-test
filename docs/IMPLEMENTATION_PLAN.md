# スーパー支払い君.com REST API 実装計画

## 概要
Golang + gin + PostgreSQLで請求書管理REST APIを実装する。
JWT認証を含む本格的なAPI実装を目指す。

## 技術スタック
- **言語**: Go 1.25.4
- **Webフレームワーク**: gin
- **データベース**: PostgreSQL 18
- **クエリビルダー**: sqlc (タイプセーフなコード生成)
- **マイグレーション**: sqldef (psqldef)
- **認証**: JWT (golang-jwt/jwt)
- **環境変数管理**: caarlos0/env v11 (構造体タグで環境変数を直接読み込み)
- **ロガー**: slog (Go標準ライブラリ)
- **API ドキュメント**: swag (コメントからOpenAPI自動生成)
- **テスト**: testify, httptest, gomock (モック生成)
- **Docker**: Docker Compose v2

## プロジェクト構造（Clean Architecture）
```
web-api-language-agnostic/
├── cmd/
│   └── api/
│       └── main.go                 # エントリーポイント
├── internal/
│   ├── domain/                     # ドメイン層（ビジネスロジック）
│   │   ├── entity/                 # エンティティ
│   │   ├── repository/             # リポジトリインターフェース
│   │   └── service/                # ドメインサービス
│   ├── usecase/                    # ユースケース層
│   │   ├── invoice/
│   │   └── auth/
│   ├── infrastructure/             # インフラ層
│   │   ├── database/               # DB接続・sqlc生成コード
│   │   │   ├── connection.go       # DB接続管理
│   │   │   └── sqlc/               # sqlc生成コード (gitignore)
│   │   │       ├── db.go
│   │   │       ├── models.go
│   │   │       ├── companies.sql.go
│   │   │       ├── users.sql.go
│   │   │       └── ...
│   │   ├── persistence/            # リポジトリ実装（sqlcラッパー）
│   │   └── security/               # JWT実装
│   └── controller/                 # コントローラー層
│       ├── auth/                   # 認証ハンドラ
│       │   ├── handler.go          # ハンドラ実装
│       │   ├── request.go          # リクエストDTO
│       │   └── response.go         # レスポンスDTO
│       ├── invoice/                # 請求書ハンドラ
│       │   ├── handler.go          # ハンドラ実装
│       │   ├── request.go          # リクエストDTO
│       │   └── response.go         # レスポンスDTO
│       └── middleware/             # ミドルウェア
│           ├── auth.go             # 認証ミドルウェア
│           └── error.go            # エラーハンドリング
├── db/
│   ├── schema.sql                  # sqldef用スキーマ定義
│   ├── queries/                    # sqlc用クエリ定義
│   │   ├── companies.sql
│   │   ├── users.sql
│   │   ├── vendors.sql
│   │   ├── vendor_bank_accounts.sql
│   │   └── invoices.sql
│   └── sqlc.yaml                   # sqlc設定ファイル
├── docs/
│   ├── swagger/                    # swag生成ファイル (gitignore)
│   │   ├── swagger.yaml
│   │   ├── swagger.json
│   │   └── docs.go
│   ├── architecture.md             # アーキテクチャ説明
│   └── coding-rules.md             # コーディングルール
├── tests/                          # 統合テスト
├── .gitignore                      # mock/, sqlc/, docs/swagger/ を含む
├── docker-compose.yml              # 環境変数も定義
├── Dockerfile
├── Makefile
└── README.md                       # 環境変数の説明を含む
```

## 実装フェーズ

### Phase 1: プロジェクトセットアップ
1. ✅ Goモジュール初期化
2. ✅ ディレクトリ構造作成
3. ✅ .gitignore作成（mock/, sqlc/, docs/swagger/等）
4. ✅ Docker Compose設定（PostgreSQL + API、環境変数定義）
5. ✅ Makefile作成（ビルド、マイグレーション、テスト、go generate、swag initコマンド）
6. ✅ コーディングルール作成（docs/coding-rules.md）

### Phase 2: データベース設計・マイグレーション
7. ✅ ERD設計
8. ✅ schema.sql作成（sqldef用）
   - companies（企業）
   - users（ユーザー）
   - vendors（支払先）
   - vendor_bank_accounts（支払先銀行口座）
   - invoices（請求書）
9. ✅ インデックス設計（支払期日検索用）
10. ✅ psqldef実行でマイグレーション

### Phase 3: sqlc設定・コード生成
11. ✅ sqlc.yaml設定
    - out: internal/infrastructure/database/sqlc
    - package: sqlc
12. ✅ SQLクエリ定義（db/queries/*.sql）
    - CRUD操作
    - 期間指定検索
    - JOIN操作
13. ✅ sqlc generate実行
    - internal/infrastructure/database/sqlc/ 配下に生成

### Phase 4: ドメイン層実装
14. ✅ エンティティ定義
    - Company, User, Vendor, VendorBankAccount, Invoice
15. ✅ リポジトリインターフェース定義
    - repository/interface.go に定義
    - //go:generate mockgen コメント追加
16. ✅ ドメインサービス実装
    - 請求金額計算ロジック（手数料4% + 消費税10%）
    - バリデーションロジック

### Phase 5: インフラ層実装
17. ✅ データベース接続（database/sql + pgx）
18. ✅ リポジトリ実装（sqlc生成コードをラップ）
19. ✅ JWT実装
    - トークン生成
    - トークン検証
    - リフレッシュトークン対応

### Phase 6: ユースケース層実装
20. ✅ ユースケースインターフェース定義（//go:generate付き）
21. ✅ 認証ユースケース実装
    - ユーザー登録
    - ログイン
    - トークン更新
22. ✅ 請求書ユースケース実装
    - 請求書作成
    - 請求書一覧取得（期間指定）

### Phase 7: コントローラー層実装
23. ✅ ミドルウェア実装（controller/middleware）
    - 認証ミドルウェア
    - エラーハンドリングミドルウェア
24. ✅ 認証ハンドラ実装（controller/auth）
    - handler.go: POST /api/auth/register, login, refresh + swagコメント
    - request.go: リクエストDTO + バリデーション
    - response.go: レスポンスDTO
25. ✅ 請求書ハンドラ実装（controller/invoice）
    - handler.go: POST /api/invoices, GET /api/invoices + swagコメント
    - request.go: リクエストDTO + バリデーション
    - response.go: レスポンスDTO
26. ✅ main.goにswag一般情報コメント追加

### Phase 8: テスト実装
27. ✅ go generate実行（gomockでモック生成）
28. ✅ 単体テスト
    - ドメインサービステスト（請求金額計算）
    - ユースケーステスト（gomockで生成したモックリポジトリ使用）
29. ✅ 統合テスト
    - APIエンドポイントテスト
    - テストデータベース使用
    - 認証フローテスト

### Phase 9: ドキュメント作成
30. ✅ swag init実行（docs/swagger/に自動生成）
31. ✅ README.md
    - プロジェクト概要
    - 環境変数の説明
    - セットアップ手順
    - go generate / swag init の使い方
    - API使用例（curl）
    - Swagger UI アクセス方法
    - テスト実行方法
32. ✅ アーキテクチャドキュメント（docs/architecture.md）

### Phase 10: 品質向上・最適化
33. ✅ パフォーマンス最適化
    - データベースインデックス確認
    - N+1問題対策（sqlcでJOIN使用）
    - コネクションプーリング設定
34. ✅ セキュリティ対策
    - パスワードハッシュ化（bcrypt）
    - SQL injection対策（sqlcのプレースホルダー）
    - CORS設定
    - Rate limiting（オプション）
35. ✅ エラーハンドリング改善
    - 適切なHTTPステータスコード
    - エラーレスポンス統一
    - slogによる構造化ロギング実装

### Phase 11: 最終確認
36. ✅ コードレビュー観点チェック
    - SOLID原則遵守
    - 責務分離
    - 可読性
    - コミット粒度
37. ✅ 全テスト実行
38. ✅ Docker環境での動作確認
39. ✅ ドキュメント最終確認（Swagger UI確認含む）

## 評価基準対応チェックリスト

### ✅ クラス・メソッド・構造体の責務
- Clean Architectureによる明確な層分離
- 各層のインターフェース定義
- 単一責任原則の遵守

### ✅ SOLID原則・アーキテクチャ
- **S**ingle Responsibility: 各構造体は単一の責務
- **O**pen/Closed: インターフェースによる拡張性
- **L**iskov Substitution: インターフェースの適切な実装
- **I**nterface Segregation: 小さく分割されたインターフェース
- **D**ependency Inversion: 依存性注入の活用

### ✅ コーディングスタイル・可読性
- gofmtによるフォーマット
- golangci-lintによるリント
- 適切なコメント・命名規則

### ✅ コミット管理
- 機能単位での小さなコミット
- わかりやすいコミットメッセージ
- レビューしやすい差分

### ✅ 認証・認可
- JWT認証の実装
- 企業単位でのデータ分離
- トークンの適切な管理

### ✅ 秘匿情報の取り扱い
- .envによる環境変数管理
- .gitignoreで秘匿情報を除外
- パスワードのハッシュ化

### ✅ テストコード
- 単体テスト（カバレッジ80%以上目標）
- 統合テスト
- テストデータの適切な管理

### ✅ パフォーマンス
- データベースインデックス
- コネクションプーリング
- 適切なクエリ設計

### ✅ エラー処理
- カスタムエラー型
- 適切なHTTPステータスコード
- エラーログの記録

### ✅ ドキュメント
- README.md（セットアップ・実行方法）
- API仕様書（OpenAPI）
- アーキテクチャドキュメント

## データモデル詳細

### 請求金額計算式
```
支払金額: payment_amount
手数料率: 4%
消費税率: 10%

手数料 = payment_amount * 0.04
手数料消費税 = 手数料 * 0.10
請求金額 = payment_amount + 手数料 + 手数料消費税
        = payment_amount * (1 + 0.04 * 1.10)
        = payment_amount * 1.044

例: 10,000円 → 10,440円
```

### ステータス遷移
```
未処理 (pending) → 処理中 (processing) → 支払い済み (paid)
                                      ↘ エラー (error)
```

## セキュリティ考慮事項

1. **認証**: JWT Bearer Token
2. **パスワード**: bcryptハッシュ化（cost: 12）
3. **環境変数**: 機密情報は環境変数で管理（caarlos0/env使用）
4. **SQL Injection**: sqlcのプレースホルダー使用
5. **CORS**: 適切なオリジン設定
6. **ログ**: パスワード等の機密情報を出力しない

## Git戦略

- **ブランチ**: main（本番相当）
- **コミット粒度**:
  - セットアップ系: まとめてOK
  - 機能実装: Phase単位またはそれより細かく
  - リファクタリング: 別コミット
  - テスト: 実装と同時 or 別コミット

## データベース管理（sqlc + sqldef）

### sqlcの利点
- **タイプセーフ**: SQLからGoコードを生成、コンパイル時に型チェック
- **パフォーマンス**: ORMのオーバーヘッドなし、生成されたコードは最適化済み
- **明示的**: SQLを直接書くため、実行されるクエリが明確
- **メンテナンス性**: スキーマ変更時もSQLとコードが同期

### sqlc 設定例 (db/sqlc.yaml)
```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "queries/"
    schema: "schema.sql"
    gen:
      go:
        package: "sqlc"
        out: "../internal/infrastructure/database/sqlc"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_interface: true
        emit_exact_table_names: false
```

### sqldefの利点
- **宣言的**: 理想のスキーマを定義するだけで差分を自動適用
- **安全**: DRY RUNで確認可能
- **シンプル**: マイグレーションファイルの履歴管理不要

## 環境変数管理（caarlos0/env）

- **caarlos0/env v11** を使用して環境変数を構造体にマッピング
- .envファイルは使用せず、環境変数から直接読み込み
- Docker Compose や Makefile で環境変数を設定

**設定例**:
```go
type Config struct {
    DBHost     string `env:"DB_HOST" envDefault:"localhost"`
    DBPort     int    `env:"DB_PORT" envDefault:"5432"`
    DBUser     string `env:"DB_USER,required"`
    DBPassword string `env:"DB_PASSWORD,required"`
    DBName     string `env:"DB_NAME,required"`
    JWTSecret  string `env:"JWT_SECRET,required"`
    Port       int    `env:"PORT" envDefault:"8080"`
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    if err := env.Parse(cfg); err != nil {
        return nil, err
    }
    return cfg, nil
}
```

## swag によるAPI仕様書生成

- **swaggo/swag** を使用してコメントからOpenAPI仕様書を自動生成
- ハンドラ関数に `@Summary`, `@Description`, `@Accept`, `@Produce`, `@Param`, `@Success`, `@Failure` 等のコメントを記述
- `swag init` でdocs/swagger/配下に生成
- Swagger UI でブラウザから確認可能

**例**:
```go
// @Summary ユーザー登録
// @Description 新規ユーザーを登録する
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "登録情報"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
    // ...
}
```

**main.goの一般情報**:
```go
// @title スーパー支払い君 API
// @version 1.0
// @description 請求書管理・支払い処理API
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
```

## モック生成

- **gomock** を使用してインターフェースのモックを生成
- `go generate` で自動生成
- 生成されたモックは gitignore

## 想定タイムライン（2-3時間）

- Phase 1-2: 30分（セットアップ、DB設計、sqldef）
- Phase 3: 15分（sqlc設定、コード生成）
- Phase 4-5: 30分（ドメイン、インフラ実装）
- Phase 6-7: 45分（ユースケース、API実装）
- Phase 8: 30分（テスト実装、go generate）
- Phase 9-11: 15分（ドキュメント、最終確認）

## 次のステップ

1. Phase 1から順次実装開始
2. 各Phase完了時にコミット
3. 疑問点があればその都度確認
4. テストを実行して動作確認

---
