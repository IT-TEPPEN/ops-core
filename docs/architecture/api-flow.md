# API Processing Flow

このドキュメントでは、OpsCoreの主要なAPI処理フローについて説明します。

## 目次

- [リクエスト処理の流れ](#リクエスト処理の流れ)
- [主要APIフロー](#主要apiフロー)
- [エラーハンドリングフロー](#エラーハンドリングフロー)
- [認証・認可フロー](#認証認可フロー)

## リクエスト処理の流れ

### 基本的な処理フロー

```mermaid
sequenceDiagram
    participant Client
    participant Middleware
    participant Handler
    participant Usecase
    participant Domain
    participant Repository
    participant DB

    Client->>Middleware: HTTP Request
    Middleware->>Middleware: ロギング
    Middleware->>Middleware: 認証チェック（将来）
    Middleware->>Handler: 処理を委譲
    Handler->>Handler: リクエスト検証
    Handler->>Handler: DTOへの変換
    Handler->>Usecase: ユースケース実行
    Usecase->>Domain: ドメインロジック実行
    Domain->>Repository: データ操作
    Repository->>DB: クエリ実行
    DB-->>Repository: 結果
    Repository-->>Domain: エンティティ
    Domain-->>Usecase: ビジネスロジック結果
    Usecase-->>Handler: DTO
    Handler->>Handler: レスポンス変換
    Handler-->>Middleware: HTTP Response
    Middleware-->>Client: JSON Response
```

### レイヤーの責務

| レイヤー | 責務 |
|---------|------|
| Middleware | ロギング、認証、CORS、レート制限 |
| Handler | リクエスト検証、DTO変換、エラーマッピング |
| Usecase | ユースケース実装、トランザクション管理 |
| Domain | ビジネスロジック、ドメインルール |
| Repository | データアクセス、永続化 |

## 主要APIフロー

### 1. ドキュメント公開フロー

```mermaid
sequenceDiagram
    participant Client
    participant Handler as DocumentHandler
    participant Usecase as DocumentUsecase
    participant Domain as Document Entity
    participant GitRepo as GitRepository
    participant DocRepo as DocumentRepository
    participant DB

    Client->>Handler: POST /api/v1/documents/:id/publish
    Handler->>Handler: リクエスト検証
    Handler->>Usecase: PublishDocument(id)
    
    Usecase->>DocRepo: FindByID(id)
    DocRepo->>DB: SELECT document
    DB-->>DocRepo: document row
    DocRepo-->>Usecase: Document Entity
    
    Usecase->>Domain: Publish(commitHash, content)
    Domain->>Domain: ビジネスルール検証
    Domain->>Domain: 新バージョン作成
    Domain-->>Usecase: Updated Document
    
    Usecase->>GitRepo: FetchContent(filePath, commitHash)
    GitRepo->>GitRepo: GitHub API呼び出し
    GitRepo-->>Usecase: Markdown Content
    
    Usecase->>DocRepo: Save(document)
    DocRepo->>DB: BEGIN TRANSACTION
    DocRepo->>DB: UPDATE documents
    DocRepo->>DB: INSERT document_versions
    DocRepo->>DB: COMMIT
    DB-->>DocRepo: Success
    DocRepo-->>Usecase: Saved Document
    
    Usecase-->>Handler: DocumentDTO
    Handler->>Handler: レスポンス変換
    Handler-->>Client: 200 OK + JSON
```

### 2. 作業証跡記録フロー

```mermaid
sequenceDiagram
    participant Client
    participant Handler as ExecutionHandler
    participant Usecase as ExecutionUsecase
    participant Domain as ExecutionRecord
    participant ExecRepo as ExecutionRepository
    participant Storage
    participant DB

    Note over Client,DB: 1. 作業証跡の作成
    Client->>Handler: POST /api/v1/execution-records
    Handler->>Usecase: CreateExecutionRecord(input)
    Usecase->>Domain: New ExecutionRecord()
    Domain->>Domain: 初期化・検証
    Domain-->>Usecase: ExecutionRecord Entity
    Usecase->>ExecRepo: Save(record)
    ExecRepo->>DB: INSERT execution_record
    DB-->>ExecRepo: record_id
    ExecRepo-->>Usecase: Saved Record
    Usecase-->>Handler: ExecutionRecordDTO
    Handler-->>Client: 201 Created

    Note over Client,DB: 2. ステップの追加
    Client->>Handler: POST /api/v1/execution-records/:id/steps
    Handler->>Usecase: AddStep(recordID, stepInput)
    Usecase->>ExecRepo: FindByID(recordID)
    ExecRepo->>DB: SELECT execution_record
    DB-->>ExecRepo: record row
    ExecRepo-->>Usecase: ExecutionRecord
    Usecase->>Domain: AddStep(stepNumber, description)
    Domain->>Domain: ステップ追加
    Domain-->>Usecase: Updated Record
    Usecase->>ExecRepo: Save(record)
    ExecRepo->>DB: INSERT execution_step
    DB-->>ExecRepo: Success
    ExecRepo-->>Usecase: Saved Record
    Usecase-->>Handler: ExecutionStepDTO
    Handler-->>Client: 201 Created

    Note over Client,DB: 3. 添付ファイルのアップロード
    Client->>Handler: POST /api/v1/execution-records/:id/attachments
    Handler->>Handler: Multipart解析
    Handler->>Usecase: UploadAttachment(recordID, file)
    Usecase->>Storage: Store(file)
    Storage->>Storage: S3/MinIO/Local保存
    Storage-->>Usecase: storagePath
    Usecase->>Domain: AttachFile(attachment)
    Domain->>Domain: 添付ファイル追加
    Domain-->>Usecase: Updated Record
    Usecase->>ExecRepo: Save(record)
    ExecRepo->>DB: INSERT attachment
    DB-->>ExecRepo: Success
    ExecRepo-->>Usecase: Saved Record
    Usecase-->>Handler: AttachmentDTO
    Handler-->>Client: 201 Created
```

### 3. 変数入力・レンダリングフロー

```mermaid
sequenceDiagram
    participant Client
    participant Handler as DocumentHandler
    participant Usecase as VariableUsecase
    participant Domain as Document
    participant DocRepo as DocumentRepository
    participant DB

    Note over Client,DB: 1. 変数定義の取得
    Client->>Handler: GET /api/v1/documents/:id/variables
    Handler->>Usecase: GetVariableDefinitions(docID)
    Usecase->>DocRepo: FindByID(docID)
    DocRepo->>DB: SELECT document, document_version
    DB-->>DocRepo: document + version data
    DocRepo-->>Usecase: Document Entity
    Usecase->>Domain: GetVariables()
    Domain-->>Usecase: VariableDefinitions[]
    Usecase-->>Handler: VariableDefinitionDTOs
    Handler-->>Client: 200 OK + Variables

    Note over Client,DB: 2. 変数値のバリデーション
    Client->>Handler: POST /api/v1/documents/:id/variables/validate
    Handler->>Usecase: ValidateVariableValues(docID, values)
    Usecase->>DocRepo: FindByID(docID)
    DocRepo->>DB: SELECT document
    DB-->>DocRepo: document data
    DocRepo-->>Usecase: Document Entity
    Usecase->>Domain: ValidateVariables(values)
    Domain->>Domain: 型チェック
    Domain->>Domain: 必須チェック
    Domain->>Domain: カスタムバリデーション
    Domain-->>Usecase: ValidationResult
    Usecase-->>Handler: ValidationDTO
    Handler-->>Client: 200 OK + Result

    Note over Client,DB: 3. 変数置換後のコンテンツ取得
    Client->>Handler: POST /api/v1/documents/:id/render
    Handler->>Usecase: RenderDocument(docID, values)
    Usecase->>DocRepo: FindByID(docID)
    DocRepo->>DB: SELECT document_version
    DB-->>DocRepo: version data
    DocRepo-->>Usecase: DocumentVersion
    Usecase->>Domain: RenderContent(values)
    Domain->>Domain: 変数置換処理
    Domain->>Domain: {{variable}} → value
    Domain-->>Usecase: RenderedContent
    Usecase-->>Handler: RenderedContentDTO
    Handler-->>Client: 200 OK + Content
```

### 4. 閲覧統計フロー

```mermaid
sequenceDiagram
    participant Client
    participant Handler as DocumentHandler
    participant ViewUsecase
    participant StatsUsecase
    participant ViewRepo as ViewHistoryRepository
    participant StatsRepo as ViewStatisticsRepository
    participant DB

    Client->>Handler: POST /api/v1/documents/:id/view
    Handler->>ViewUsecase: RecordView(docID, userID, ipAddr)
    
    par 閲覧履歴の記録
        ViewUsecase->>ViewRepo: Save(viewHistory)
        ViewRepo->>DB: INSERT view_history
        DB-->>ViewRepo: Success
        ViewRepo-->>ViewUsecase: Saved
    and 閲覧統計の更新
        ViewUsecase->>StatsUsecase: IncrementView(docID, userID)
        StatsUsecase->>StatsRepo: FindByDocumentID(docID)
        StatsRepo->>DB: SELECT view_statistics
        DB-->>StatsRepo: stats data
        StatsRepo-->>StatsUsecase: ViewStatistics
        StatsUsecase->>StatsRepo: Update(stats)
        StatsRepo->>DB: UPDATE view_statistics
        DB-->>StatsRepo: Success
        StatsRepo-->>StatsUsecase: Updated
        StatsUsecase-->>ViewUsecase: Success
    end
    
    ViewUsecase-->>Handler: Success
    Handler-->>Client: 204 No Content
```

## エラーハンドリングフロー

### エラー伝播の流れ

```mermaid
sequenceDiagram
    participant Client
    participant Handler
    participant Usecase
    participant Domain
    participant Repository
    participant DB

    Client->>Handler: HTTP Request
    Handler->>Usecase: Process()
    Usecase->>Domain: BusinessLogic()
    Domain->>Domain: Validation Error!
    Domain-->>Usecase: DomainError
    Usecase->>Usecase: Wrap Error
    Usecase-->>Handler: ApplicationError
    Handler->>Handler: Map to HTTP Error
    Handler->>Handler: Create ErrorResponse
    Handler-->>Client: 400 Bad Request<br/>+ Error JSON

    Note over Client,DB: 別のケース: DB エラー
    Client->>Handler: HTTP Request
    Handler->>Usecase: Process()
    Usecase->>Repository: Query()
    Repository->>DB: SQL
    DB-->>Repository: Connection Error!
    Repository-->>Usecase: InfraError
    Usecase->>Usecase: Wrap Error
    Usecase-->>Handler: ApplicationError
    Handler->>Handler: Map to HTTP Error
    Handler->>Handler: Create ErrorResponse
    Handler-->>Client: 500 Internal Error<br/>+ Error JSON
```

### エラーレスポンス例

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": [
      {
        "field": "title",
        "message": "title is required"
      }
    ]
  }
}
```

## 認証・認可フロー

### JWT認証フロー（将来実装）

```mermaid
sequenceDiagram
    participant Client
    participant AuthMiddleware
    participant Handler
    participant Usecase
    participant UserRepo
    participant DB

    Note over Client,DB: 1. ログイン
    Client->>Handler: POST /api/v1/auth/login
    Handler->>Usecase: Authenticate(username, password)
    Usecase->>UserRepo: FindByEmail(email)
    UserRepo->>DB: SELECT user
    DB-->>UserRepo: user data
    UserRepo-->>Usecase: User Entity
    Usecase->>Usecase: Verify Password
    Usecase->>Usecase: Generate JWT
    Usecase-->>Handler: JWT Token
    Handler-->>Client: 200 OK + {token, user}

    Note over Client,DB: 2. 保護されたリソースへのアクセス
    Client->>AuthMiddleware: GET /api/v1/documents<br/>Authorization: Bearer <token>
    AuthMiddleware->>AuthMiddleware: Extract Token
    AuthMiddleware->>AuthMiddleware: Verify JWT Signature
    AuthMiddleware->>AuthMiddleware: Check Expiration
    AuthMiddleware->>UserRepo: FindByID(userID)
    UserRepo->>DB: SELECT user
    DB-->>UserRepo: user data
    UserRepo-->>AuthMiddleware: User Entity
    AuthMiddleware->>AuthMiddleware: Set User in Context
    AuthMiddleware->>Handler: Request + User Context
    Handler->>Handler: Check Authorization
    Handler->>Usecase: ListDocuments(user)
    Usecase-->>Handler: Documents
    Handler-->>Client: 200 OK + Documents
```

## トランザクション管理

### トランザクション境界

```go
// Usecase層でトランザクション管理
func (u *DocumentUsecase) PublishDocument(ctx context.Context, id string) error {
    // トランザクション開始
    tx, err := u.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // ビジネスロジック実行
    doc, err := u.docRepo.FindByID(ctx, tx, id)
    if err != nil {
        return err
    }

    if err := doc.Publish(); err != nil {
        return err
    }

    if err := u.docRepo.Save(ctx, tx, doc); err != nil {
        return err
    }

    // コミット
    return tx.Commit()
}
```

### トランザクション分離レベル

- **デフォルト**: READ COMMITTED
- **特殊ケース**: SERIALIZABLE（統計更新など）

## パフォーマンス最適化

### 1. N+1問題の回避

```go
// 悪い例: N+1クエリ
documents := repo.FindAll()
for _, doc := range documents {
    versions := repo.FindVersions(doc.ID) // N回のクエリ
}

// 良い例: Eager Loading
documents := repo.FindAllWithVersions() // 1回のJOINクエリ
```

### 2. ページネーション

```go
// リクエスト
GET /api/v1/documents?page=2&per_page=20

// レスポンス
{
  "documents": [...],
  "pagination": {
    "page": 2,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

### 3. 並行処理

```go
// Goroutineを使用した並行処理
func (u *Usecase) Process() error {
    var wg sync.WaitGroup
    errCh := make(chan error, 2)

    wg.Add(2)
    go func() {
        defer wg.Done()
        if err := u.task1(); err != nil {
            errCh <- err
        }
    }()

    go func() {
        defer wg.Done()
        if err := u.task2(); err != nil {
            errCh <- err
        }
    }()

    wg.Wait()
    close(errCh)

    for err := range errCh {
        if err != nil {
            return err
        }
    }
    return nil
}
```

## 関連ドキュメント

- [システム概要](./system-overview.md)
- [データベーススキーマ](./database-schema.md)
- [API開発ガイド](../development/API.md)
- [ADR 0015: Backend Custom Error Design](../../adr/0015-backend-custom-error-design.md)
