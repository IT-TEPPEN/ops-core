# Testing Guide

このドキュメントでは、OpsCoreプロジェクトにおけるテスト戦略とテストの書き方を説明します。

## 目次

- [テスト戦略](#テスト戦略)
- [バックエンドテスト](#バックエンドテスト)
- [フロントエンドテスト](#フロントエンドテスト)
- [統合テスト](#統合テスト)
- [テストカバレッジ](#テストカバレッジ)

## テスト戦略

OpsCoreでは、以下のテスト戦略を採用しています：

### テストピラミッド

```
        /\
       /  \  E2E Tests（最小）
      /----\
     /      \ Integration Tests（中）
    /--------\
   /          \ Unit Tests（最大）
  /------------\
```

1. **ユニットテスト（Unit Tests）**: 最も多く、各関数・メソッド・コンポーネントの単体テスト
2. **統合テスト（Integration Tests）**: 中程度、複数モジュールの連携テスト
3. **E2Eテスト（End-to-End Tests）**: 最小限、主要なユーザーフローのテスト（将来実装予定）

### テストの品質基準

- **カバレッジ**: 80%以上を目標
- **テストの独立性**: 各テストは独立して実行可能
- **テストの速度**: ユニットテストは高速（1秒未満/テスト）
- **テストの可読性**: テスト名で何をテストしているか明確に

## バックエンドテスト

バックエンドテストの詳細は [ADR 0009: Backend Testing Strategy](../../adr/0009-backend-testing-strategy.md) を参照してください。

### テスト環境のセットアップ

```bash
cd backend

# 依存関係のインストール
go mod download

# テストの実行
go test ./...
```

### ユニットテスト

#### ドメイン層のテスト

ドメイン層のテストは、ビジネスロジックの正確性を検証します。

**例: エンティティのテスト**

```go
// internal/document/domain/entity/document_test.go
package entity

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestDocument_Publish(t *testing.T) {
    tests := []struct {
        name        string
        document    *Document
        commitHash  string
        content     string
        wantErr     bool
        errContains string
    }{
        {
            name: "正常系: ドキュメントの公開に成功",
            document: &Document{
                ID:          NewDocumentID(),
                IsPublished: false,
                Versions:    []DocumentVersion{},
            },
            commitHash: "abc123",
            content:    "# Test Document",
            wantErr:    false,
        },
        {
            name: "異常系: 既に公開済み",
            document: &Document{
                ID:          NewDocumentID(),
                IsPublished: true,
            },
            commitHash:  "abc123",
            content:     "# Test Document",
            wantErr:     true,
            errContains: "already published",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.document.Publish(tt.commitHash, tt.content)
            if tt.wantErr {
                assert.Error(t, err)
                if tt.errContains != "" {
                    assert.Contains(t, err.Error(), tt.errContains)
                }
            } else {
                assert.NoError(t, err)
                assert.True(t, tt.document.IsPublished)
            }
        })
    }
}
```

**例: 値オブジェクトのテスト**

```go
// internal/document/domain/value_object/document_id_test.go
package value_object

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestNewDocumentID(t *testing.T) {
    id1 := NewDocumentID()
    id2 := NewDocumentID()
    
    assert.NotEmpty(t, id1.String())
    assert.NotEqual(t, id1, id2)
}

func TestDocumentID_Equals(t *testing.T) {
    id1 := NewDocumentID()
    id2 := id1
    id3 := NewDocumentID()
    
    assert.True(t, id1.Equals(id2))
    assert.False(t, id1.Equals(id3))
}
```

#### アプリケーション層のテスト

アプリケーション層のテストは、ユースケースの正確性を検証します。モックリポジトリを使用します。

**例: ユースケースのテスト**

```go
// internal/document/application/usecase/document_usecase_test.go
package usecase

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestDocumentUsecase_CreateDocument(t *testing.T) {
    // モックリポジトリの準備
    mockRepo := new(MockDocumentRepository)
    usecase := NewDocumentUsecase(mockRepo)

    tests := []struct {
        name    string
        input   CreateDocumentInput
        setup   func()
        wantErr bool
    }{
        {
            name: "正常系: ドキュメント作成成功",
            input: CreateDocumentInput{
                RepositoryID: "repo-123",
                FilePath:     "docs/example.md",
                Title:        "Example",
            },
            setup: func() {
                mockRepo.On("Save", mock.Anything, mock.Anything).
                    Return(nil).Once()
            },
            wantErr: false,
        },
        {
            name: "異常系: リポジトリIDが空",
            input: CreateDocumentInput{
                FilePath: "docs/example.md",
                Title:    "Example",
            },
            setup:   func() {},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.setup()
            
            _, err := usecase.CreateDocument(context.Background(), tt.input)
            
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
            
            mockRepo.AssertExpectations(t)
        })
    }
}
```

#### インフラストラクチャ層のテスト

インフラストラクチャ層のテストは、実際のデータベースまたはテスト用DBを使用します。

**例: リポジトリ実装のテスト**

```go
// internal/document/infrastructure/persistence/document_repository_impl_test.go
package persistence

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestDocumentRepositoryImpl_Save(t *testing.T) {
    // テスト用DBのセットアップ
    db := setupTestDB(t)
    defer teardownTestDB(t, db)
    
    repo := NewDocumentRepositoryImpl(db)
    
    doc := &entity.Document{
        ID:           value_object.NewDocumentID(),
        RepositoryID: "repo-123",
        FilePath:     "docs/example.md",
        Title:        "Example",
    }
    
    // 保存
    err := repo.Save(context.Background(), doc)
    require.NoError(t, err)
    
    // 取得して検証
    saved, err := repo.FindByID(context.Background(), doc.ID)
    require.NoError(t, err)
    assert.Equal(t, doc.ID, saved.ID)
    assert.Equal(t, doc.Title, saved.Title)
}
```

### テストの実行

```bash
# 全テストの実行
go test ./...

# カバレッジ付き
go test -cover ./...

# カバレッジレポート生成
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# 特定パッケージのテスト
go test ./internal/document/...

# ベンチマークテスト
go test -bench=. ./...

# レースコンディションの検出
go test -race ./...
```

### モックの作成

モックは、`gomock`または手動で作成します。

```go
// 手動モックの例
type MockDocumentRepository struct {
    mock.Mock
}

func (m *MockDocumentRepository) Save(ctx context.Context, doc *entity.Document) error {
    args := m.Called(ctx, doc)
    return args.Error(0)
}

func (m *MockDocumentRepository) FindByID(ctx context.Context, id value_object.DocumentID) (*entity.Document, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.Document), args.Error(1)
}
```

## フロントエンドテスト

### テスト環境のセットアップ

```bash
cd frontend

# 依存関係のインストール
npm install

# テストの実行
npm test
```

### ユニットテスト

#### ユーティリティ関数のテスト

```typescript
// src/utils/formatDate.test.ts
import { describe, it, expect } from "vitest";
import { formatDate } from "./formatDate";

describe("formatDate", () => {
  it("日付を正しくフォーマットする", () => {
    const date = new Date("2024-01-15T10:30:00Z");
    const formatted = formatDate(date);
    expect(formatted).toBe("2024-01-15 10:30");
  });

  it("nullの場合は空文字を返す", () => {
    const formatted = formatDate(null);
    expect(formatted).toBe("");
  });
});
```

#### カスタムフックのテスト

```typescript
// src/hooks/useDocuments.test.ts
import { describe, it, expect, beforeEach } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { useDocuments } from "./useDocuments";

describe("useDocuments", () => {
  it("ドキュメント一覧を取得する", async () => {
    const { result } = renderHook(() => useDocuments());

    await waitFor(() => {
      expect(result.current.documents).toHaveLength(2);
    });

    expect(result.current.loading).toBe(false);
    expect(result.current.error).toBeNull();
  });
});
```

#### コンポーネントのテスト

```typescript
// src/components/VariableForm.test.tsx
import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { VariableForm } from "./VariableForm";

describe("VariableForm", () => {
  const mockVariables = [
    {
      name: "server_name",
      label: "Server Name",
      type: "string",
      required: true,
      defaultValue: "server01",
    },
  ];

  it("変数フォームを正しく表示する", () => {
    render(
      <VariableForm 
        variables={mockVariables} 
        onSubmit={() => {}} 
      />
    );

    expect(screen.getByLabelText("Server Name")).toBeInTheDocument();
  });

  it("デフォルト値が入力されている", () => {
    render(
      <VariableForm 
        variables={mockVariables} 
        onSubmit={() => {}} 
      />
    );

    const input = screen.getByLabelText("Server Name") as HTMLInputElement;
    expect(input.value).toBe("server01");
  });

  it("必須項目が未入力の場合はエラーを表示", async () => {
    const mockOnSubmit = vi.fn();
    
    render(
      <VariableForm 
        variables={mockVariables} 
        onSubmit={mockOnSubmit} 
      />
    );

    const input = screen.getByLabelText("Server Name") as HTMLInputElement;
    fireEvent.change(input, { target: { value: "" } });

    const submitButton = screen.getByRole("button", { name: "Submit" });
    fireEvent.click(submitButton);

    expect(await screen.findByText(/required/i)).toBeInTheDocument();
    expect(mockOnSubmit).not.toHaveBeenCalled();
  });
});
```

### テストの実行

```bash
# 全テストの実行
npm test

# ウォッチモード（開発時）
npm run test:watch

# カバレッジ付き
npm run test:coverage

# UI上でカバレッジを確認
npm run test:coverage
# coverage/index.htmlをブラウザで開く
```

## 統合テスト

統合テストは、複数のモジュールが正しく連携することを検証します。

### バックエンド統合テスト

```go
// backend/cmd/api_tester/endpoint_tests.go
func TestFullDocumentWorkflow(t *testing.T) {
    // テスト用サーバーの起動
    server := setupTestServer(t)
    defer server.Close()
    
    // 1. ドキュメントの作成
    doc := createDocument(t, server.URL)
    assert.NotEmpty(t, doc.ID)
    
    // 2. ドキュメントの公開
    err := publishDocument(t, server.URL, doc.ID)
    assert.NoError(t, err)
    
    // 3. ドキュメントの取得
    retrieved := getDocument(t, server.URL, doc.ID)
    assert.True(t, retrieved.IsPublished)
    
    // 4. ドキュメントの削除
    err = deleteDocument(t, server.URL, doc.ID)
    assert.NoError(t, err)
}
```

### E2Eテスト（将来実装予定）

Playwrightなどを使用したブラウザ自動テストを予定しています。

```typescript
// e2e/document-workflow.spec.ts (将来)
import { test, expect } from '@playwright/test';

test('ドキュメント管理のフルワークフロー', async ({ page }) => {
  // ログイン
  await page.goto('http://localhost:5173/login');
  await page.fill('input[name="username"]', 'admin');
  await page.fill('input[name="password"]', 'password');
  await page.click('button[type="submit"]');

  // ドキュメント作成
  await page.goto('http://localhost:5173/documents/new');
  await page.fill('input[name="title"]', 'Test Document');
  await page.click('button[type="submit"]');

  // 作成確認
  await expect(page.locator('text=Test Document')).toBeVisible();
});
```

## テストカバレッジ

### カバレッジ目標

- **全体**: 80%以上
- **ドメイン層**: 90%以上（ビジネスロジックの中核）
- **アプリケーション層**: 85%以上
- **インフラストラクチャ層**: 70%以上
- **インターフェース層**: 75%以上

### カバレッジの確認

#### バックエンド

```bash
cd backend
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

#### フロントエンド

```bash
cd frontend
npm run test:coverage
# coverage/index.htmlをブラウザで開く
```

### カバレッジレポートの見方

- **緑色**: テストでカバーされている
- **赤色**: テストでカバーされていない
- **黄色**: 部分的にカバーされている

未カバーの重要な箇所には優先的にテストを追加してください。

## ベストプラクティス

### テストの命名

- **明確**: テスト名で何をテストしているか明確に
- **構造化**: "正常系"/"異常系"などで分類

```go
func TestDocument_Publish_Success(t *testing.T) { }
func TestDocument_Publish_AlreadyPublished_Error(t *testing.T) { }
```

### テストの独立性

- 各テストは独立して実行可能
- テストの実行順序に依存しない
- グローバル状態を変更しない

### テストのメンテナンス

- テストコードも本番コードと同様に品質を保つ
- 重複コードはヘルパー関数にまとめる
- テストが失敗したら、まずテスト自体が正しいか確認

## 参考資料

- [ADR 0009: Backend Testing Strategy](../../adr/0009-backend-testing-strategy.md)
- [Effective Go](https://go.dev/doc/effective_go)
- [React Testing Library](https://testing-library.com/docs/react-testing-library/intro/)
- [Vitest](https://vitest.dev/)
