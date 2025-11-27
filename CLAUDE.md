# CLAUDE.md

このファイルはClaude Codeがプロジェクトを理解するためのエントリーポイントです。

## プロジェクト概要

**スーパー支払い君.com API** - B2B向け請求書管理・支払い処理REST APIサービス

## ドキュメント

作業前に以下を確認してください：

1. **[docs/ASSIGNMENT.md](docs/ASSIGNMENT.md)** - 課題要件（最優先で確認）
2. **[docs/IMPLEMENTATION_PLAN.md](docs/IMPLEMENTATION_PLAN.md)** - 実装計画・進捗管理
3. **[docs/CODING_RULES.md](docs/CODING_RULES.md)** - コーディング規約・API設計ルール
4. **[README.md](README.md)** - セットアップ手順・API仕様

## クイックリファレンス

```bash
make build      # ビルド
make run        # 起動
make test       # テスト
make generate   # コード生成 (mock + sqlc)
make fmt        # フォーマット
make lint       # リント
make migrate    # マイグレーション
make docker-up  # Docker起動
```

## 技術スタック要点

- Go 1.25.4 / PostgreSQL 18 / gin
- sqlc（クエリビルダー）、sqldef（マイグレーション）
- caarlos0/env（環境変数直接読み込み、.env不使用）
- golangci-lint v2（フォーマット・リント統一）
