# OpsCore

OpsCoreは、運用の中核を担うシステムである。

## 概要

OpsCoreは、GitHubやGitLabなどの外部リポジトリに保存された運用手順書（Markdownファイル）を集約し、Web上で閲覧できるようにするシステムです。運用者はこのシステムを通じて各種手順書やナレッジベースにアクセスし、日々の運用業務を実施します。

## 開発フェーズと機能

本プロジェクトは段階的に機能を実装していきます。

### Phase 1: 基本的なMarkdown閲覧機能（✅ 実装完了）

Phase 1では、外部リポジトリからMarkdownファイルを取得・表示する基本機能を提供します。

#### 実装済み機能

1. **リポジトリ管理**
   - GitHubリポジトリの登録・一覧表示・削除
   - アクセストークンの設定・更新
   - **アクセストークンの暗号化（AES-256-GCM）**
   - リポジトリ情報の永続化（PostgreSQL）

2. **ドキュメント管理**
   - リポジトリからのMarkdownファイル一覧取得
   - ドキュメントの公開・非公開設定
   - ドキュメントのバージョン管理（コミットハッシュと連番管理）
   - 公開範囲の設定（public/private）
   - 自動更新設定（リポジトリ更新時の自動反映）

3. **変数入力機能**
   - 手順書への変数定義（Frontmatter形式）
   - 変数の型指定（string/number/boolean/date）
   - 必須/任意の設定とデフォルト値
   - 変数値の入力UIと置換表示

4. **Markdown表示**
   - 選択された複数のMarkdownファイルを統合して表示
   - ブログ形式での閲覧UI
   - 変数置換後のコンテンツ表示

5. **バックエンドアーキテクチャ**
   - Onion Architecture（ドメイン駆動設計）
   - コンテキスト分割（Git Repository, Document, Execution Record, User, View統計）
   - Repository パターンによるデータアクセス
   - 依存性注入（Wire）
   - カスタムエラー設計とエラーハンドリング
   - ロギング（zap）
   - データベースマイグレーション（golang-migrate）

6. **フロントエンド**
   - React + TypeScript
   - リポジトリ管理画面
   - ドキュメント管理画面
   - 変数入力フォーム
   - Markdown閲覧画面

### Phase 2: 作業証跡と検索機能（🚧 実装中）

Phase 2では、運用業務の記録と効率化を支援する機能を追加します。

#### 実装済み機能

1. **作業証跡管理**
   - 手順書実行セッションの記録
   - 各ステップへのメモと画面キャプチャ添付
   - 作業ステータス管理（in_progress/completed/failed）
   - 作業証跡の共有機能
   - ローカルファイルシステム/S3/MinIOへの添付ファイル保存

2. **ユーザー・グループ管理**
   - ユーザー管理（admin/userロール）
   - グループ管理とメンバーシップ
   - アクセス制御とアクセス範囲の設定

3. **閲覧履歴・統計**
   - ドキュメント閲覧履歴の記録
   - 閲覧統計情報（総閲覧数、ユニークユーザー数、最終閲覧日時）

#### 計画中の機能

1. **ユーザー認証・認可**
   - ユーザー登録・ログイン機能
   - ロールベースアクセス制御（RBAC）
   - リポジトリごとのアクセス権限管理

2. **検索・フィルタリング**
   - 手順書のタイトル・タグによる検索
   - カテゴリ別フィルタリング
   - 最近閲覧した手順書の表示
   - 作業証跡の検索・フィルタリング

### Phase 3: AI連携と自動化（将来構想）

Phase 3では、生成AIを活用した高度な運用支援機能を提供します。

#### 構想中の機能

1. **生成AIとの連携**
   - 手順書の内容に基づいた質疑応答
   - 作業手順の要約生成
   - トラブルシューティング支援

2. **作業の自動化**
   - 定型作業の自動実行
   - チェックリスト形式での作業進捗管理
   - 作業完了の自動検証

3. **改善提案**
   - 作業ログの分析による改善点の提示
   - 頻繁に参照される手順書の最適化提案
   - 作業時間の傾向分析

## アーキテクチャ

詳細な設計判断については、`adr/`ディレクトリ内のArchitecture Decision Records（ADR）を参照してください。

### バックエンド

- 言語: Go
- アーキテクチャ: Onion Architecture
- データベース: PostgreSQL
- 主要ライブラリ: Echo（Webフレームワーク）、Wire（DI）、zap（ロギング）

### フロントエンド

- 言語: TypeScript
- フレームワーク: React
- ビルドツール: Vite

## ドキュメント

詳細なドキュメントは以下を参照してください：

- **開発ガイド**
  - [開発ガイドライン（CONTRIBUTING.md）](docs/development/CONTRIBUTING.md)
  - [テストガイド](docs/development/TESTING.md)
  - [API開発ガイド](docs/development/API.md)
- **ユーザーマニュアル**
  - [ユーザーガイド](docs/user-guide/README.md)
  - [ドキュメント管理](docs/user-guide/document-management.md)
  - [作業証跡記録](docs/user-guide/execution-record.md)
  - [変数入力機能](docs/user-guide/variable-input.md)
  - [グループ管理](docs/user-guide/group-management.md)
- **アーキテクチャ**
  - [システム概要](docs/architecture/system-overview.md)
  - [DBスキーマ図](docs/architecture/database-schema.md)
  - [API処理フロー](docs/architecture/api-flow.md)
  - [バックエンドアーキテクチャ](backend/README.md)
- **運用ガイド**
  - [デプロイ手順](docs/deployment/README.md)
  - [監視・運用](docs/operations/MONITORING.md)
  - [バックアップ・リストア](docs/operations/BACKUP.md)

## セットアップ

開発環境のセットアップ方法や貢献ガイドラインについては、[CONTRIBUTING.md](docs/development/CONTRIBUTING.md)を参照してください。

### 必要な環境

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- Docker & Docker Compose（開発環境）

### クイックスタート

```bash
# 開発環境の起動（Docker Compose）
docker compose up -d

# 暗号化キーの設定（開発環境用）
export ENCRYPTION_KEY="dev-key-12345678901234567890123"  # 32 bytes

# バックエンドの起動
cd backend
go run cmd/server/main.go

# フロントエンドの起動
cd frontend
npm install
npm run dev
```

**重要**: 本番環境では、セキュアな32バイトの暗号化キーを生成し、環境変数 `ENCRYPTION_KEY` に設定してください。詳細は [backend/internal/git_repository/infrastructure/encryption/README.md](backend/internal/git_repository/infrastructure/encryption/README.md) を参照してください。

## ライセンス

TBD

