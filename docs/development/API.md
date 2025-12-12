# API Development Guide

このドキュメントでは、OpsCoreのAPI開発に関するガイドラインを説明します。

## 目次

- [API設計原則](#api設計原則)
- [RESTful API設計](#restful-api設計)
- [エンドポイント一覧](#エンドポイント一覧)
- [リクエスト・レスポンス形式](#リクエストレスポンス形式)
- [エラーハンドリング](#エラーハンドリング)
- [認証・認可](#認証認可)
- [API仕様書](#api仕様書)

## API設計原則

OpsCoreのAPIは、以下の原則に基づいて設計されています：

1. **RESTful設計**: リソース指向の設計
2. **一貫性**: 命名規則とレスポンス形式の統一
3. **セキュリティ**: 認証・認可とデータ保護
4. **バージョニング**: API互換性の維持
5. **ドキュメント**: OpenAPI/Swagger仕様の提供

## RESTful API設計

### HTTPメソッドの使い分け

| HTTPメソッド | 用途 | べき等性 | 安全性 |
|-------------|------|---------|--------|
| GET | リソースの取得 | ✅ | ✅ |
| POST | リソースの作成 | ❌ | ❌ |
| PUT | リソースの完全置換 | ✅ | ❌ |
| PATCH | リソースの部分更新 | ❌ | ❌ |
| DELETE | リソースの削除 | ✅ | ❌ |

### URLパターン

```
# コレクション取得
GET /api/v1/documents

# 単一リソース取得
GET /api/v1/documents/{id}

# リソース作成
POST /api/v1/documents

# リソース更新
PUT /api/v1/documents/{id}
PATCH /api/v1/documents/{id}

# リソース削除
DELETE /api/v1/documents/{id}

# サブリソース
GET /api/v1/documents/{id}/versions
POST /api/v1/documents/{id}/publish

# アクション
POST /api/v1/execution-records/{id}/complete
```

### 命名規則

- **URL**: ケバブケース（例: `/execution-records`）
- **JSONキー**: スネークケース（例: `document_id`, `created_at`）
- **複数形**: コレクションは複数形（例: `/documents`, `/users`）

## エンドポイント一覧

### Git Repository API

#### リポジトリ管理

```
GET    /api/v1/repositories           # リポジトリ一覧取得
POST   /api/v1/repositories           # リポジトリ登録
GET    /api/v1/repositories/{id}      # リポジトリ取得
PUT    /api/v1/repositories/{id}      # リポジトリ更新
DELETE /api/v1/repositories/{id}      # リポジトリ削除

GET    /api/v1/repositories/{id}/files # リポジトリ内のファイル一覧取得
```

### Document API

#### ドキュメント管理

```
GET    /api/v1/documents              # ドキュメント一覧取得
POST   /api/v1/documents              # ドキュメント作成
GET    /api/v1/documents/{id}         # ドキュメント取得
PUT    /api/v1/documents/{id}         # ドキュメント更新
DELETE /api/v1/documents/{id}         # ドキュメント削除

POST   /api/v1/documents/{id}/publish # ドキュメント公開
POST   /api/v1/documents/{id}/unpublish # ドキュメント非公開

GET    /api/v1/documents/{id}/versions # バージョン一覧取得
POST   /api/v1/documents/{id}/rollback # バージョンロールバック
```

#### 変数管理

```
GET    /api/v1/documents/{id}/variables # 変数定義取得
POST   /api/v1/documents/{id}/variables/validate # 変数値のバリデーション
POST   /api/v1/documents/{id}/render    # 変数置換後のコンテンツ取得
```

### Execution Record API

#### 作業証跡管理

```
GET    /api/v1/execution-records       # 作業証跡一覧取得
POST   /api/v1/execution-records       # 作業証跡作成
GET    /api/v1/execution-records/{id}  # 作業証跡取得
PUT    /api/v1/execution-records/{id}  # 作業証跡更新
DELETE /api/v1/execution-records/{id}  # 作業証跡削除

POST   /api/v1/execution-records/{id}/steps # ステップ追加
POST   /api/v1/execution-records/{id}/complete # 作業完了
POST   /api/v1/execution-records/{id}/fail # 作業失敗
```

#### 添付ファイル管理

```
POST   /api/v1/execution-records/{id}/attachments # 添付ファイルアップロード
GET    /api/v1/execution-records/{id}/attachments/{attachment_id} # 添付ファイル取得
DELETE /api/v1/execution-records/{id}/attachments/{attachment_id} # 添付ファイル削除
```

### User & Group API

#### ユーザー管理

```
GET    /api/v1/users                   # ユーザー一覧取得
POST   /api/v1/users                   # ユーザー作成
GET    /api/v1/users/{id}              # ユーザー取得
PUT    /api/v1/users/{id}              # ユーザー更新
DELETE /api/v1/users/{id}              # ユーザー削除
```

#### グループ管理

```
GET    /api/v1/groups                  # グループ一覧取得
POST   /api/v1/groups                  # グループ作成
GET    /api/v1/groups/{id}             # グループ取得
PUT    /api/v1/groups/{id}             # グループ更新
DELETE /api/v1/groups/{id}             # グループ削除

POST   /api/v1/groups/{id}/members     # メンバー追加
DELETE /api/v1/groups/{id}/members/{user_id} # メンバー削除
```

### View History & Statistics API

#### 閲覧履歴

```
GET    /api/v1/documents/{id}/view-history # 閲覧履歴取得
POST   /api/v1/documents/{id}/view         # 閲覧記録
```

#### 閲覧統計

```
GET    /api/v1/documents/{id}/statistics   # 閲覧統計取得
```

## リクエスト・レスポンス形式

### リクエスト

#### Content-Type

```
Content-Type: application/json
```

#### リクエストボディ例

```json
// POST /api/v1/documents
{
  "repository_id": "repo-123",
  "file_path": "docs/example.md",
  "title": "Example Document",
  "type": "procedure",
  "tags": ["deploy", "production"],
  "access_scope": {
    "type": "private",
    "shared_with": ["user-456", "group-789"]
  }
}
```

### レスポンス

#### 成功レスポンス

```json
// 200 OK - 単一リソース
{
  "id": "doc-123",
  "repository_id": "repo-123",
  "title": "Example Document",
  "type": "procedure",
  "is_published": false,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}

// 200 OK - コレクション
{
  "documents": [
    {
      "id": "doc-123",
      "title": "Document 1"
    },
    {
      "id": "doc-456",
      "title": "Document 2"
    }
  ],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  }
}

// 201 Created - リソース作成
{
  "id": "doc-123",
  "repository_id": "repo-123",
  "title": "Example Document",
  "created_at": "2024-01-15T10:30:00Z"
}

// 204 No Content - 削除成功
// (レスポンスボディなし)
```

#### エラーレスポンス

```json
// 400 Bad Request - バリデーションエラー
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "title",
        "message": "title is required"
      },
      {
        "field": "type",
        "message": "type must be 'procedure' or 'knowledge'"
      }
    ]
  }
}

// 401 Unauthorized - 認証エラー
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Authentication required"
  }
}

// 403 Forbidden - 認可エラー
{
  "error": {
    "code": "FORBIDDEN",
    "message": "You don't have permission to access this resource"
  }
}

// 404 Not Found - リソース未発見
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Document not found",
    "resource_type": "document",
    "resource_id": "doc-123"
  }
}

// 409 Conflict - 競合エラー
{
  "error": {
    "code": "CONFLICT",
    "message": "Document is already published"
  }
}

// 500 Internal Server Error - サーバーエラー
{
  "error": {
    "code": "INTERNAL_SERVER_ERROR",
    "message": "An unexpected error occurred",
    "request_id": "req-xyz789"
  }
}
```

### ステータスコード

| コード | 意味 | 用途 |
|-------|------|------|
| 200 | OK | リソースの取得・更新成功 |
| 201 | Created | リソースの作成成功 |
| 204 | No Content | リソースの削除成功 |
| 400 | Bad Request | リクエストが不正 |
| 401 | Unauthorized | 認証エラー |
| 403 | Forbidden | 認可エラー |
| 404 | Not Found | リソースが見つからない |
| 409 | Conflict | リソースの競合 |
| 422 | Unprocessable Entity | バリデーションエラー |
| 500 | Internal Server Error | サーバーエラー |

## エラーハンドリング

詳細は [ADR 0015: Backend Custom Error Design](../../adr/0015-backend-custom-error-design.md) を参照してください。

### エラーコード体系

```
# ドメイン層エラー
DOMAIN_001: Invalid entity state
DOMAIN_002: Business rule violation

# アプリケーション層エラー
APP_001: Validation error
APP_002: Not found
APP_003: Already exists

# インフラストラクチャ層エラー
INFRA_001: Database error
INFRA_002: External service error

# API層エラー
API_001: Invalid request
API_002: Authentication failed
API_003: Authorization failed
```

### エラーレスポンスの実装

```go
// internal/document/interfaces/error/mapper.go
func MapError(err error) echo.HTTPError {
    var domainErr *domain.DocumentError
    var appErr *application.ApplicationError
    
    switch {
    case errors.As(err, &domainErr):
        return mapDomainError(domainErr)
    case errors.As(err, &appErr):
        return mapApplicationError(appErr)
    default:
        return echo.NewHTTPError(
            http.StatusInternalServerError,
            ErrorResponse{
                Code:    "INTERNAL_SERVER_ERROR",
                Message: "An unexpected error occurred",
            },
        )
    }
}
```

## 認証・認可

### 認証（Authentication）

#### JWT認証（計画中）

```
Authorization: Bearer <JWT_TOKEN>
```

#### セッション認証（計画中）

```
Cookie: session_id=<SESSION_ID>
```

### 認可（Authorization）

#### ロールベースアクセス制御（RBAC）

- **admin**: 全てのリソースへのアクセス権限
- **user**: 自分のリソースと共有されたリソースへのアクセス権限

#### リソースベースアクセス制御

```go
// ドキュメントのアクセス制御例
func (h *DocumentHandler) GetDocument(c echo.Context) error {
    documentID := c.Param("id")
    userID := getUserIDFromContext(c)
    
    doc, err := h.usecase.GetDocument(c.Request().Context(), documentID)
    if err != nil {
        return err
    }
    
    // アクセス権限のチェック
    if !doc.IsAccessibleBy(userID) {
        return echo.NewHTTPError(http.StatusForbidden, "Access denied")
    }
    
    return c.JSON(http.StatusOK, doc)
}
```

## API仕様書

### Swagger/OpenAPIの使用

OpsCoreでは、Swagger/OpenAPI 3.0を使用してAPI仕様を定義しています。

#### Swagger UIへのアクセス

開発サーバーを起動後、以下のURLでSwagger UIにアクセスできます：

```
http://localhost:8080/swagger/index.html
```

#### Swagger仕様の生成

```bash
cd backend

# swaggoのインストール
go install github.com/swaggo/swag/cmd/swag@latest

# Swagger仕様の生成
swag init -g cmd/server/main.go -o docs

# 生成されるファイル
# - docs/docs.go
# - docs/swagger.json
# - docs/swagger.yaml
```

#### Swaggerアノテーションの書き方

```go
// @Summary ドキュメント一覧取得
// @Description 公開されているドキュメントの一覧を取得します
// @Tags documents
// @Accept json
// @Produce json
// @Param page query int false "ページ番号" default(1)
// @Param per_page query int false "1ページあたりの件数" default(20)
// @Success 200 {object} ListDocumentsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/documents [get]
func (h *DocumentHandler) ListDocuments(c echo.Context) error {
    // 実装
}
```

### API仕様書のメンテナンス

1. **コードとの同期**: Swaggerアノテーションを最新に保つ
2. **例の提供**: リクエスト・レスポンスの例を記載
3. **エラーケース**: 全てのエラーケースを文書化
4. **バージョン管理**: API変更時はバージョン番号を更新

## クエリパラメータ

### ページネーション

```
GET /api/v1/documents?page=2&per_page=20
```

- `page`: ページ番号（1から開始、デフォルト: 1）
- `per_page`: 1ページあたりの件数（デフォルト: 20、最大: 100）

### フィルタリング

```
GET /api/v1/documents?type=procedure&tag=deploy
```

### ソート

```
GET /api/v1/documents?sort_by=created_at&order=desc
```

- `sort_by`: ソート対象のフィールド
- `order`: `asc`（昇順）または`desc`（降順）

### 検索

```
GET /api/v1/documents?q=keyword
```

- `q`: 検索キーワード

## ベストプラクティス

### 1. べき等性の保証

PUT、DELETE、GETは必ずべき等に実装してください。

### 2. レート制限

将来的に、API使用量の制限を実装予定：

```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640000000
```

### 3. キャッシュ制御

```
Cache-Control: max-age=3600, must-revalidate
ETag: "33a64df551425fcc55e4d42a148795d9f25f89d4"
```

### 4. CORS対応

```go
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{"http://localhost:5173"},
    AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
}))
```

### 5. ロギング

全てのAPIリクエストをログに記録：

```
[INFO] 2024-01-15 10:30:00 GET /api/v1/documents 200 45ms user=user-123
```

## テスト

APIエンドポイントのテストは、統合テストとして実装します。

```go
func TestDocumentAPI_CreateDocument(t *testing.T) {
    server := setupTestServer(t)
    defer server.Close()
    
    payload := map[string]interface{}{
        "repository_id": "repo-123",
        "title": "Test Document",
    }
    
    resp, err := makeRequest(
        "POST",
        server.URL+"/api/v1/documents",
        payload,
    )
    
    assert.NoError(t, err)
    assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
```

## 関連ドキュメント

- [ADR 0003: Backend API Markdown Fetch](../../adr/0003-backend-api-markdown-fetch.md)
- [ADR 0010: API Definition Generation Specification](../../adr/0010-api-definition-generation-specification.md)
- [ADR 0015: Backend Custom Error Design](../../adr/0015-backend-custom-error-design.md)
- [バックエンドアーキテクチャ](../../backend/README.md)
- [テストガイド](./TESTING.md)
