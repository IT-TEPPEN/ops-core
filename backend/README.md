# OpsCore Backend

OpsCoreバックエンドは、Go言語で実装された、運用手順書管理システムのサーバーサイドアプリケーションです。

## アーキテクチャ概要

### Onion Architecture

本システムは、Onion Architecture（オニオンアーキテクチャ）に基づいて設計されています。詳細は [ADR 0007](../adr/0007-backend-architecture-onion.md) を参照してください。

```
backend/
├── cmd/                    # エントリーポイント
│   ├── server/            # APIサーバー
│   └── tools/             # ツール類
├── internal/              # プライベートコード
│   ├── document/          # ドキュメント管理コンテキスト
│   ├── execution_record/  # 作業証跡コンテキスト
│   ├── git_repository/    # Gitリポジトリ管理コンテキスト
│   ├── user/             # ユーザー管理コンテキスト
│   ├── view_history/     # 閲覧履歴コンテキスト
│   ├── view_statistics/  # 閲覧統計コンテキスト
│   └── shared/           # 共通モジュール
├── infrastructure/        # インフラストラクチャ層（ストレージ等）
└── docs/                 # API仕様書（Swagger）
```

### コンテキスト分割

ドメイン駆動設計（DDD）の考え方に基づき、以下のコンテキストに分割されています：

#### 1. Git Repository Context
外部Gitリポジトリとの連携を担当。

**主要エンティティ:**
- Repository（集約ルート）: Gitリポジトリ情報の管理

**主要機能:**
- リポジトリの登録・更新・削除
- アクセストークンの暗号化管理
- Gitファイル一覧の取得

#### 2. Document Context
ドキュメントのライフサイクル管理を担当。

**主要エンティティ:**
- Document（集約ルート）: ドキュメント情報とバージョン管理
- DocumentVersion: ドキュメントの特定バージョン

**主要値オブジェクト:**
- DocumentID, VersionID, RepositoryID
- DocumentSource（FilePath + CommitHash）
- DocumentType（procedure/knowledge）
- Tag, VariableDefinition
- AccessScope（public/private + 共有設定）

**主要機能:**
- ドキュメントの公開・非公開
- バージョン管理とロールバック
- 変数定義の管理
- アクセス範囲の設定
- 自動更新設定

**参考ADR:**
- [ADR 0013: Document Variable Definition](../adr/0013-document-variable-definition.md)
- [ADR 0016: Document Domain Model Design](../adr/0016-document-domain-model-design.md)
- [ADR 0017: Application Data Validation](../adr/0017-application-data-validation.md)

#### 3. Execution Record Context
手順書実行の作業証跡を管理。

**主要エンティティ:**
- ExecutionRecord（集約ルート）: 作業証跡の管理
- ExecutionStep: 作業ステップの記録
- Attachment: 画面キャプチャなどの添付ファイル

**主要値オブジェクト:**
- ExecutionRecordID, AttachmentID
- VariableValue（変数の入力値）
- StorageType（local/s3/minio）

**主要機能:**
- 作業証跡の作成・更新
- 作業ステップの記録
- 画面キャプチャの添付
- 作業ステータス管理（in_progress/completed/failed）
- 作業証跡の共有
- 作業証跡の検索・フィルタリング

**参考ADR:**
- [ADR 0014: Execution Record and Evidence Management](../adr/0014-execution-record-and-evidence-management.md)

#### 4. User Context
ユーザーとグループの管理を担当。

**主要エンティティ:**
- User（集約ルート）: ユーザー情報
- Group（集約ルート）: グループ情報

**主要値オブジェクト:**
- UserID, GroupID

**主要機能:**
- ユーザーのCRUD操作
- グループのCRUD操作
- グループメンバーシップ管理
- ロールベース認可（admin/user）

#### 5. View History Context
ドキュメントの閲覧履歴を記録。

**主要エンティティ:**
- ViewHistory: 閲覧履歴レコード

**主要機能:**
- 閲覧履歴の記録
- 閲覧履歴の取得

#### 6. View Statistics Context
ドキュメントの閲覧統計を管理。

**主要エンティティ:**
- ViewStatistics（集約ルート）: 閲覧統計情報

**主要機能:**
- 閲覧数の集計
- ユニークユーザー数の管理
- 最終閲覧日時の記録

### レイヤー構成

各コンテキストは以下のレイヤーで構成されています：

```
context/
├── domain/              # ドメイン層（ビジネスロジック）
│   ├── entity/         # エンティティ
│   ├── value_object/   # 値オブジェクト
│   ├── repository/     # リポジトリインターフェース
│   └── error/          # ドメインエラー
├── application/        # アプリケーション層（ユースケース）
│   ├── usecase/        # ユースケース実装
│   ├── dto/            # データ転送オブジェクト
│   └── error/          # アプリケーションエラー
├── infrastructure/     # インフラストラクチャ層（永続化）
│   ├── persistence/    # データベース実装
│   └── ...
└── interfaces/         # インターフェース層（API）
    ├── api/
    │   ├── handlers/   # HTTPハンドラー
    │   └── schema/     # APIスキーマ
    └── error/          # HTTPエラー
```

## 技術スタック

- **言語**: Go 1.21+
- **Webフレームワーク**: Echo v4
- **依存性注入**: Wire
- **ロギング**: zap
- **データベース**: PostgreSQL 14+
- **マイグレーション**: golang-migrate
- **テスト**: 標準テストパッケージ + testify
- **API仕様**: Swagger/OpenAPI（swaggo）

## 主要な設計判断

### カスタムエラー設計

各層で独自のエラー型を定義し、適切なエラーハンドリングを実現しています。
詳細は [ADR 0015: Backend Custom Error Design](../adr/0015-backend-custom-error-design.md) を参照してください。

**エラーコード体系:**
- `DOMAIN_XXX`: ドメイン層エラー
- `APP_XXX`: アプリケーション層エラー
- `INFRA_XXX`: インフラストラクチャ層エラー
- `API_XXX`: API層エラー

### データベーススキーマ

PostgreSQLを使用し、以下の主要テーブルで構成されています：

- `repositories`: Gitリポジトリ情報
- `documents`: ドキュメント情報
- `document_versions`: ドキュメントバージョン履歴
- `execution_records`: 作業証跡
- `execution_steps`: 作業ステップ
- `attachments`: 添付ファイル
- `users`: ユーザー情報
- `groups`: グループ情報
- `user_groups`: ユーザー・グループ関連
- `view_histories`: 閲覧履歴
- `view_statistics`: 閲覧統計

詳細は [ADR 0005: Database Schema](../adr/0005-database-schema.md) を参照してください。

### ストレージ抽象化

添付ファイルの保存先を抽象化し、複数のストレージバックエンドに対応しています：

- **ローカルファイルシステム**: 開発環境向け
- **AWS S3**: 本番環境向け
- **MinIO**: オンプレミス環境向け

詳細は [infrastructure/storage/README.md](infrastructure/storage/README.md) を参照してください。

### セキュリティ

#### アクセストークンの暗号化

Gitリポジトリのアクセストークンは、AES-256-GCM方式で暗号化してデータベースに保存されます。
詳細は [internal/git_repository/infrastructure/encryption/README.md](internal/git_repository/infrastructure/encryption/README.md) を参照してください。

#### データバリデーション

アプリケーション層でデータバリデーションを実施し、不正なデータの混入を防ぎます。
詳細は [ADR 0017: Application Data Validation](../adr/0017-application-data-validation.md) を参照してください。

## 開発

### ビルドと実行

```bash
# 依存関係のインストール
go mod download

# ビルド
go build -o bin/server cmd/server/main.go

# 実行
./bin/server
```

### テストの実行

```bash
# 全テストの実行
go test ./...

# カバレッジ付きテスト
go test -cover ./...

# 特定パッケージのテスト
go test ./internal/document/...
```

### マイグレーション

データベースマイグレーションの詳細は [ADR 0006: Database Migration](../adr/0006-database-migration.md) を参照してください。

```bash
# マイグレーションの適用
migrate -path db/migrations -database "postgres://user:pass@localhost:5432/opscore?sslmode=disable" up

# マイグレーションのロールバック
migrate -path db/migrations -database "postgres://user:pass@localhost:5432/opscore?sslmode=disable" down 1
```

### API仕様書の生成

```bash
# Swagger仕様の生成
swag init -g cmd/server/main.go -o docs

# API仕様書の確認
# http://localhost:8080/swagger/index.html
```

## 関連ドキュメント

- [ADR（Architecture Decision Records）](../adr/)
- [API開発ガイド](../docs/development/API.md)
- [テストガイド](../docs/development/TESTING.md)
- [システム概要](../docs/architecture/system-overview.md)
