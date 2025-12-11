# System Overview

このドキュメントでは、OpsCoreのシステム全体のアーキテクチャと構成を説明します。

## 目次

- [システムアーキテクチャ](#システムアーキテクチャ)
- [コンポーネント構成](#コンポーネント構成)
- [データフロー](#データフロー)
- [技術スタック](#技術スタック)
- [デプロイ構成](#デプロイ構成)

## システムアーキテクチャ

### 全体構成図

```mermaid
graph TB
    subgraph "External Services"
        GitHub[GitHub/GitLab<br/>Repository]
        S3[AWS S3 /  MinIO<br/>Object Storage]
    end

    subgraph "Frontend Layer"
        Browser[Web Browser]
        React[React Application<br/>Vite + TypeScript]
    end

    subgraph "Backend Layer"
        API[API Server<br/>Go + Echo]
        
        subgraph "Domain Contexts"
            GitRepo[Git Repository<br/>Context]
            Document[Document<br/>Context]
            Execution[Execution Record<br/>Context]
            User[User<br/>Context]
            ViewStats[View Statistics<br/>Context]
        end
    end

    subgraph "Data Layer"
        PostgreSQL[(PostgreSQL<br/>Database)]
        Storage[File Storage<br/>Local/S3/MinIO]
    end

    Browser -->|HTTP/HTTPS| React
    React -->|REST API| API
    API -->|Query/Command| GitRepo
    API -->|Query/Command| Document
    API -->|Query/Command| Execution
    API -->|Query/Command| User
    API -->|Query/Command| ViewStats
    
    GitRepo -->|Fetch Files| GitHub
    Document -->|Read/Write| PostgreSQL
    Execution -->|Read/Write| PostgreSQL
    Execution -->|Upload/Download| Storage
    User -->|Read/Write| PostgreSQL
    ViewStats -->|Read/Write| PostgreSQL
    
    Storage -->|Store Objects| S3
    
    style React fill:#61dafb
    style API fill:#00add8
    style PostgreSQL fill:#336791
    style GitHub fill:#181717
    style S3 fill:#ff9900
```

## コンポーネント構成

### フロントエンド

#### React Application

- **技術**: React 18 + TypeScript + Vite
- **状態管理**: React Hooks（useState, useEffect）
- **ルーティング**: React Router
- **UIライブラリ**: カスタムコンポーネント

**主要ページ**:
```
/                           ホーム
/repositories               リポジトリ一覧
/documents                  ドキュメント一覧
/documents/:id              ドキュメント詳細
/execution-records          作業証跡一覧
/execution-records/:id      作業証跡詳細
/users                      ユーザー管理（管理者のみ）
/groups                     グループ管理（管理者のみ）
```

### バックエンド

#### API Server（Go + Echo）

- **アーキテクチャ**: Onion Architecture
- **パターン**: DDD（Domain-Driven Design）
- **コンテキスト分割**: 6つの境界づけられたコンテキスト

#### コンテキスト一覧

1. **Git Repository Context**
   - 外部Gitリポジトリとの連携
   - ファイル一覧の取得
   - アクセストークンの暗号化管理

2. **Document Context**
   - ドキュメントのライフサイクル管理
   - バージョン管理
   - 変数定義管理
   - アクセス制御

3. **Execution Record Context**
   - 作業証跡の記録
   - ステップ管理
   - 添付ファイル管理

4. **User Context**
   - ユーザー管理
   - グループ管理
   - 認証・認可（将来実装）

5. **View History Context**
   - 閲覧履歴の記録

6. **View Statistics Context**
   - 閲覧統計の集計

#### レイヤー構成

```
Context/
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

### データベース

#### PostgreSQL

- **バージョン**: 14+
- **用途**: 構造化データの永続化
- **主要テーブル**: 17テーブル

**テーブルグループ**:
```
- リポジトリ管理: repositories
- ドキュメント管理: documents, document_versions
- 作業証跡: execution_records, execution_steps, attachments
- ユーザー管理: users, groups, user_groups
- 閲覧管理: view_histories, view_statistics
```

### ストレージ

#### ファイルストレージ

- **ローカル**: 開発環境向け
- **AWS S3**: 本番環境向け
- **MinIO**: オンプレミス環境向け

**保存データ**:
- 作業証跡の画面キャプチャ
- 添付ファイル

## データフロー

### ドキュメント閲覧フロー

```mermaid
sequenceDiagram
    participant Browser
    participant React
    participant API
    participant Document
    participant PostgreSQL
    participant GitHub

    Browser->>React: ドキュメント一覧表示
    React->>API: GET /api/v1/documents
    API->>Document: ListDocuments()
    Document->>PostgreSQL: SELECT documents
    PostgreSQL-->>Document: documents data
    Document-->>API: DocumentList
    API-->>React: JSON response
    React-->>Browser: ドキュメント一覧

    Browser->>React: ドキュメント選択
    React->>API: GET /api/v1/documents/:id
    API->>Document: GetDocument(id)
    Document->>PostgreSQL: SELECT document
    PostgreSQL-->>Document: document data
    Document->>GitHub: Fetch markdown content
    GitHub-->>Document: markdown content
    Document-->>API: Document with content
    API-->>React: JSON response
    React-->>Browser: ドキュメント表示
```

### 作業証跡記録フロー

```mermaid
sequenceDiagram
    participant Browser
    participant React
    participant API
    participant Execution
    participant PostgreSQL
    participant Storage

    Browser->>React: 作業証跡開始
    React->>API: POST /api/v1/execution-records
    API->>Execution: CreateExecutionRecord()
    Execution->>PostgreSQL: INSERT execution_record
    PostgreSQL-->>Execution: record_id
    Execution-->>API: ExecutionRecord
    API-->>React: JSON response
    React-->>Browser: 作業証跡ID

    Browser->>React: ステップ追加
    React->>API: POST /api/v1/execution-records/:id/steps
    API->>Execution: AddStep()
    Execution->>PostgreSQL: INSERT execution_step
    PostgreSQL-->>Execution: step_id
    Execution-->>API: ExecutionStep
    API-->>React: JSON response
    React-->>Browser: ステップ記録完了

    Browser->>React: 画像アップロード
    React->>API: POST /api/v1/execution-records/:id/attachments<br/>(multipart/form-data)
    API->>Execution: UploadAttachment()
    Execution->>Storage: Store file
    Storage-->>Execution: file_path
    Execution->>PostgreSQL: INSERT attachment
    PostgreSQL-->>Execution: attachment_id
    Execution-->>API: Attachment
    API-->>React: JSON response
    React-->>Browser: アップロード完了
```

### 認証フロー（将来実装）

```mermaid
sequenceDiagram
    participant Browser
    participant React
    participant API
    participant User
    participant PostgreSQL

    Browser->>React: ログイン画面
    React->>API: POST /api/v1/auth/login<br/>{username, password}
    API->>User: Authenticate()
    User->>PostgreSQL: SELECT user
    PostgreSQL-->>User: user data
    User->>User: Verify password
    User-->>API: JWT token
    API-->>React: {token, user}
    React->>React: Store token
    React-->>Browser: リダイレクト

    Browser->>React: 保護されたリソース
    React->>API: GET /api/v1/documents<br/>Authorization: Bearer <token>
    API->>API: Verify JWT
    API->>User: GetUserByID()
    User->>PostgreSQL: SELECT user
    PostgreSQL-->>User: user data
    User-->>API: User
    API->>Document: ListDocuments(user)
    Document->>PostgreSQL: SELECT documents
    PostgreSQL-->>Document: documents
    Document-->>API: DocumentList
    API-->>React: JSON response
    React-->>Browser: ドキュメント一覧
```

## 技術スタック

### フロントエンド

| カテゴリ | 技術 | 用途 |
|---------|------|------|
| フレームワーク | React 18 | UI構築 |
| 言語 | TypeScript | 型安全な開発 |
| ビルドツール | Vite | 高速な開発環境 |
| ルーティング | React Router | SPA routing |
| テスト | Vitest + RTL | 単体・統合テスト |
| リンター | ESLint | コード品質 |

### バックエンド

| カテゴリ | 技術 | 用途 |
|---------|------|------|
| 言語 | Go 1.21+ | API実装 |
| フレームワーク | Echo v4 | Webフレームワーク |
| DI | Wire | 依存性注入 |
| ORM | pgx | PostgreSQLクライアント |
| ロギング | zap | 構造化ロギング |
| マイグレーション | golang-migrate | DBマイグレーション |
| API仕様 | swaggo | OpenAPI/Swagger |
| テスト | testify | アサーション |

### インフラストラクチャ

| カテゴリ | 技術 | 用途 |
|---------|------|------|
| データベース | PostgreSQL 14+ | 構造化データ |
| ストレージ | S3/MinIO | オブジェクトストレージ |
| コンテナ | Docker | アプリケーション実行環境 |
| オーケストレーション | Docker Compose | 開発環境 |
| 監視 | Prometheus/Grafana | メトリクス監視 |
| ログ | Loki | ログ集約 |

## デプロイ構成

### 開発環境

```mermaid
graph TB
    subgraph "Developer Machine"
        Frontend[Frontend<br/>npm run dev<br/>:5173]
        Backend[Backend<br/>go run main.go<br/>:8080]
    end

    subgraph "Docker Compose"
        PostgreSQL[(PostgreSQL<br/>:5432)]
        MinIO[MinIO<br/>:9000]
    end

    Frontend -->|API Call| Backend
    Backend -->|Query| PostgreSQL
    Backend -->|Store| MinIO
```

### 本番環境（想定）

```mermaid
graph TB
    subgraph "Load Balancer"
        LB[Load Balancer<br/>ALB/NLB]
    end

    subgraph "Application Servers"
        Frontend1[Frontend<br/>Nginx + Static]
        Frontend2[Frontend<br/>Nginx + Static]
        Backend1[Backend<br/>Container]
        Backend2[Backend<br/>Container]
    end

    subgraph "Data Layer"
        RDS[(RDS PostgreSQL<br/>Multi-AZ)]
        S3[(S3 Bucket)]
    end

    subgraph "Monitoring"
        Prometheus[Prometheus]
        Grafana[Grafana]
        Loki[Loki]
    end

    Internet --> LB
    LB --> Frontend1
    LB --> Frontend2
    Frontend1 --> Backend1
    Frontend2 --> Backend2
    Backend1 --> RDS
    Backend2 --> RDS
    Backend1 --> S3
    Backend2 --> S3
    
    Backend1 -.->|Metrics| Prometheus
    Backend2 -.->|Metrics| Prometheus
    Backend1 -.->|Logs| Loki
    Backend2 -.->|Logs| Loki
    Prometheus -.->|Data| Grafana
    Loki -.->|Data| Grafana
```

### Kubernetesデプロイ（将来）

```mermaid
graph TB
    subgraph "Kubernetes Cluster"
        subgraph "Ingress"
            Ingress[Ingress Controller]
        end

        subgraph "Frontend"
            FrontendDep[Frontend Deployment<br/>3 replicas]
            FrontendSvc[Frontend Service]
        end

        subgraph "Backend"
            BackendDep[Backend Deployment<br/>3 replicas]
            BackendSvc[Backend Service]
        end

        subgraph "Data"
            PostgresStateful[PostgreSQL StatefulSet]
            PostgresService[PostgreSQL Service]
            PVC[Persistent Volume Claims]
        end
    end

    subgraph "External Services"
        S3External[AWS S3]
    end

    Internet --> Ingress
    Ingress --> FrontendSvc
    FrontendSvc --> FrontendDep
    FrontendDep --> BackendSvc
    BackendSvc --> BackendDep
    BackendDep --> PostgresService
    PostgresService --> PostgresStateful
    PostgresStateful --> PVC
    BackendDep --> S3External
```

## スケーラビリティ

### 水平スケーリング

- **フロントエンド**: 複数インスタンスの起動（ステートレス）
- **バックエンド**: 複数インスタンスの起動（ステートレス）
- **データベース**: Read Replicaの追加

### 垂直スケーリング

- **バックエンド**: CPUとメモリの増強
- **データベース**: インスタンスサイズの拡大

### キャッシング戦略（将来実装）

- **CDN**: フロントエンドの静的ファイル
- **Redis**: セッション情報、頻繁にアクセスされるデータ
- **アプリケーションレベル**: インメモリキャッシュ

## セキュリティ

### ネットワークセキュリティ

- **HTTPS**: TLS 1.2以上
- **ファイアウォール**: 必要なポートのみ開放
- **VPC**: プライベートネットワーク

### アプリケーションセキュリティ

- **認証**: JWT（将来実装）
- **認可**: RBAC（Role-Based Access Control）
- **入力検証**: バリデーション層での実施
- **暗号化**: アクセストークンの暗号化（AES-256-GCM）

### データセキュリティ

- **データベース**: SSL/TLS接続
- **ストレージ**: サーバーサイド暗号化
- **バックアップ**: 暗号化されたバックアップ

## パフォーマンス

### 目標値

| メトリクス | 目標値 |
|-----------|--------|
| API応答時間（95パーセンタイル） | < 2秒 |
| ページロード時間 | < 3秒 |
| スループット | > 100 req/s |
| エラー率 | < 1% |

### 最適化戦略

1. **データベースクエリ**: インデックス最適化
2. **N+1問題**: Eager Loading
3. **ファイルサイズ**: 圧縮と最適化
4. **並行処理**: Goroutineの活用

## 関連ドキュメント

- [データベーススキーマ](./database-schema.md)
- [API処理フロー](./api-flow.md)
- [バックエンドアーキテクチャ](../../backend/README.md)
- [デプロイガイド](../deployment/README.md)
