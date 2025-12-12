# Contributing to OpsCore

OpsCoreへのコントリビューションをご検討いただき、ありがとうございます。このドキュメントでは、開発に参加する際のガイドラインを説明します。

## 目次

- [開発環境のセットアップ](#開発環境のセットアップ)
- [開発フロー](#開発フロー)
- [コーディング規約](#コーディング規約)
- [テスト](#テスト)
- [コミットメッセージ](#コミットメッセージ)
- [プルリクエスト](#プルリクエスト)
- [Issue作成](#issue作成)

## 開発環境のセットアップ

### 必要な環境

- **Go**: 1.21以上
- **Node.js**: 18以上
- **PostgreSQL**: 14以上
- **Docker & Docker Compose**: 最新版（推奨）

### クイックスタート

1. **リポジトリのクローン**

```bash
git clone https://github.com/IT-TEPPEN/ops-core.git
cd ops-core
```

2. **Docker Composeでの起動（推奨）**

```bash
# データベースの起動
docker compose up -d

# 環境変数の設定
export ENCRYPTION_KEY="dev-key-12345678901234567890123"  # 32 bytes
export STORAGE_TYPE="local"  # or "s3" or "minio"
```

3. **バックエンドの起動**

```bash
cd backend

# 依存関係のインストール
go mod download

# データベースマイグレーション（初回のみ）
# migrate -path db/migrations -database "postgres://user:pass@localhost:5432/opscore?sslmode=disable" up

# サーバーの起動
go run cmd/server/main.go
```

4. **フロントエンドの起動**

```bash
cd frontend

# 依存関係のインストール
npm install

# 開発サーバーの起動
npm run dev
```

5. **アプリケーションへのアクセス**

- フロントエンド: http://localhost:5173
- バックエンドAPI: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html

### 環境変数

以下の環境変数を設定する必要があります：

#### バックエンド

```bash
# 必須
ENCRYPTION_KEY="your-32-byte-encryption-key-here"  # アクセストークン暗号化キー

# データベース（デフォルト値）
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="opscore"
DB_PASSWORD="opscore"
DB_NAME="opscore"
DB_SSLMODE="disable"

# ストレージ設定
STORAGE_TYPE="local"  # local, s3, minio
LOCAL_STORAGE_PATH="./storage"

# S3設定（STORAGE_TYPE=s3の場合）
S3_REGION="ap-northeast-1"
S3_BUCKET="opscore-attachments"
AWS_ACCESS_KEY_ID="your-access-key"
AWS_SECRET_ACCESS_KEY="your-secret-key"

# MinIO設定（STORAGE_TYPE=minioの場合）
MINIO_ENDPOINT="localhost:9000"
MINIO_ACCESS_KEY="minioadmin"
MINIO_SECRET_KEY="minioadmin"
MINIO_BUCKET="opscore-attachments"
MINIO_USE_SSL="false"
```

#### フロントエンド

```bash
# APIエンドポイント
VITE_API_BASE_URL="http://localhost:8080"
```

## 開発フロー

### ブランチ戦略

- `main`: 本番環境用の安定版ブランチ
- `develop`: 開発用のブランチ
- `feature/*`: 新機能開発用ブランチ
- `bugfix/*`: バグ修正用ブランチ
- `hotfix/*`: 緊急修正用ブランチ

### 開発の流れ

1. **Issueの作成または確認**
   - 作業を始める前に、関連するIssueが存在することを確認
   - Issueが無い場合は、新規作成して内容を明確化

2. **ブランチの作成**

```bash
# developブランチから作業ブランチを作成
git checkout develop
git pull origin develop
git checkout -b feature/issue-number-short-description
```

3. **開発とテスト**
   - コードを実装
   - テストを追加・実行
   - リンターでコードスタイルをチェック

4. **コミット**
   - 論理的な単位でコミット
   - [コミットメッセージ規約](#コミットメッセージ)に従う

5. **プッシュとプルリクエスト**

```bash
git push origin feature/issue-number-short-description
```

6. **コードレビュー**
   - レビュアーからのフィードバックに対応
   - 必要に応じて修正・追加コミット

7. **マージ**
   - レビュー承認後、developブランチにマージ

## コーディング規約

### Go（バックエンド）

#### 基本方針

- **Effective Go**および**Go Code Review Comments**に従う
- `gofmt`で自動フォーマット
- `golangci-lint`でリンターチェック

#### Onion Architectureの遵守

詳細は [ADR 0007: Backend Architecture - Onion Architecture](../../adr/0007-backend-architecture-onion.md) を参照してください。

**依存関係のルール:**
- ドメイン層は他の層に依存しない
- アプリケーション層はドメイン層にのみ依存
- インフラストラクチャ層はドメイン層とアプリケーション層に依存
- インターフェース層は全ての層に依存可能

**ディレクトリ構成:**

```
internal/
└── context_name/
    ├── domain/              # ドメイン層
    │   ├── entity/         # エンティティ
    │   ├── value_object/   # 値オブジェクト
    │   ├── repository/     # リポジトリインターフェース
    │   └── error/          # ドメインエラー
    ├── application/        # アプリケーション層
    │   ├── usecase/        # ユースケース
    │   ├── dto/            # DTO
    │   └── error/          # アプリケーションエラー
    ├── infrastructure/     # インフラストラクチャ層
    │   └── persistence/    # データベース実装
    └── interfaces/         # インターフェース層
        └── api/
            ├── handlers/   # HTTPハンドラー
            └── schema/     # APIスキーマ
```

#### 命名規則

- **パッケージ名**: 小文字の単数形（例: `document`, `user`）
- **ファイル名**: スネークケース（例: `document_usecase.go`）
- **型名**: PascalCase（例: `DocumentID`, `ExecutionRecord`）
- **関数/メソッド名**: PascalCase（例: `CreateDocument`, `GetByID`）
- **変数名**: camelCase（例: `documentID`, `executionRecord`）
- **定数**: PascalCaseまたはUPPER_SNAKE_CASE（例: `MaxRetries`, `DEFAULT_TIMEOUT`）

#### エラーハンドリング

詳細は [ADR 0015: Backend Custom Error Design](../../adr/0015-backend-custom-error-design.md) を参照してください。

- 各層で独自のエラー型を定義
- `errors.Is`、`errors.As`でエラーチェック
- エラーをラップする際は`fmt.Errorf("context: %w", err)`を使用
- エラーコードを付与（`DOMAIN_XXX`, `APP_XXX`, `INFRA_XXX`, `API_XXX`）

#### テストの書き方

- テストファイルは`*_test.go`
- テーブルドリブンテストを推奨
- モックは`gomock`または手動で作成
- テストカバレッジは80%以上を目標

```go
func TestDocumentUsecase_CreateDocument(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateDocumentInput
        want    *Document
        wantErr bool
    }{
        // テストケース
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // テスト実装
        })
    }
}
```

### TypeScript（フロントエンド）

#### 基本方針

- **TypeScript**の型安全性を最大限活用
- **ESLint**でリンターチェック
- **Prettier**で自動フォーマット（設定ファイルに従う）

#### ディレクトリ構成

```
src/
├── components/      # 再利用可能なコンポーネント
├── pages/          # ページコンポーネント
├── hooks/          # カスタムフック
├── utils/          # ユーティリティ関数
├── types/          # 型定義
├── services/       # API通信
└── App.tsx         # アプリケーションルート
```

#### 命名規則

- **コンポーネント**: PascalCase（例: `VariableForm.tsx`, `ExecutionRecordPanel.tsx`）
- **フック**: camelCaseで`use`プレフィックス（例: `useDocuments.ts`）
- **ユーティリティ**: camelCase（例: `formatDate.ts`）
- **型定義**: PascalCase（例: `Document`, `ExecutionRecord`）

#### コンポーネントの書き方

- 関数コンポーネントを使用
- `React.FC`は使用しない（不要）
- Propsの型は別途定義

```typescript
type VariableFormProps = {
  variables: VariableDefinition[];
  onSubmit: (values: VariableValue[]) => void;
};

export function VariableForm({ variables, onSubmit }: VariableFormProps) {
  // 実装
}
```

#### テストの書き方

- テストファイルは`*.test.tsx`または`*.test.ts`
- Vitestを使用
- React Testing Libraryでコンポーネントテスト

```typescript
import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { MyComponent } from "./MyComponent";

describe("MyComponent", () => {
  it("renders correctly", () => {
    render(
      <MemoryRouter>
        <MyComponent />
      </MemoryRouter>
    );
    expect(screen.getByText("Hello")).toBeInTheDocument();
  });
});
```

## テスト

### バックエンドテスト

```bash
cd backend

# 全テストの実行
go test ./...

# カバレッジ付き
go test -cover ./...

# 特定パッケージのテスト
go test ./internal/document/...
```

### フロントエンドテスト

```bash
cd frontend

# 全テストの実行
npm test

# ウォッチモード
npm run test:watch

# カバレッジ付き
npm run test:coverage
```

詳細は [TESTING.md](./TESTING.md) を参照してください。

## コミットメッセージ

### フォーマット

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type

- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメントのみの変更
- `style`: コードの意味に影響しない変更（空白、フォーマット等）
- `refactor`: リファクタリング
- `test`: テストの追加・修正
- `chore`: ビルドプロセスやツールの変更

### 例

```
feat(document): 変数入力機能の追加

ドキュメントに定義された変数を入力できるフォームを実装。
型別のバリデーションとデフォルト値の設定に対応。

Closes #38
```

## プルリクエスト

### 作成前のチェックリスト

- [ ] 関連するIssueが存在する
- [ ] テストが追加され、全て通過している
- [ ] リンターエラーがない
- [ ] ドキュメントが更新されている（必要な場合）
- [ ] コミットメッセージが規約に従っている

### PRテンプレート

```markdown
## 概要
このPRの目的を簡潔に説明

## 変更内容
- 変更点1
- 変更点2

## 関連Issue
Closes #issue-number

## テスト
- [ ] 単体テスト追加
- [ ] 統合テスト追加
- [ ] 手動テスト実施

## スクリーンショット（UI変更の場合）
（スクリーンショットを添付）
```

## Issue作成

### Issueテンプレート

#### バグ報告

```markdown
## バグの概要
バグの内容を簡潔に説明

## 再現手順
1. ステップ1
2. ステップ2
3. ...

## 期待される動作
本来どうあるべきか

## 実際の動作
実際にどうなったか

## 環境
- OS: 
- ブラウザ: 
- バージョン: 
```

#### 機能リクエスト

```markdown
## 機能の概要
実装したい機能の概要

## 背景・目的
なぜこの機能が必要か

## 提案する実装
具体的な実装案

## 代替案
他に考えられる方法
```

## ADR（Architecture Decision Records）

重要な設計判断は、ADRとして文書化してください。

ADRの書き方は [ADR 0000: ADR Writing Guidelines](../../adr/0000-adr-writing-guidelines.md) を参照してください。

## セキュリティ

セキュリティ上の問題を発見した場合は、公開Issueではなく、メンテナーに直接連絡してください。

## ライセンス

このプロジェクトに貢献することで、あなたの貢献がプロジェクトのライセンスに従うことに同意したものとみなされます。

## 質問・サポート

- **ドキュメント**: まず関連ドキュメントを確認
- **Issue**: 既存のIssueを検索
- **新規Issue**: 解決しない場合は新規Issueを作成

ご協力ありがとうございます！
