# データベースER図

現在のマイグレーションファイル (000001〜000008) から生成したER図です。

## ER図全体（マイグレーション順）

```mermaid
flowchart TB
    subgraph m001["000001: 認証基盤"]
        users[users]
        refresh_tokens[refresh_tokens]
        user_identities[user_identities]
        user_profiles[user_profiles]
        primary_user_identity[primary_user_identity]
        
        users --> refresh_tokens
        users --> user_identities
        users --> user_profiles
        users --> primary_user_identity
        user_identities --> primary_user_identity
    end
    
    subgraph m002["000002: Gitプロバイダー管理"]
        git_providers[git_providers]
    end
    
    subgraph m003["000003: OAuth接続"]
        oauth_connections[oauth_connections]
        
        users -.-> oauth_connections
        git_providers --> oauth_connections
    end
    
    subgraph m004["000004: リポジトリ・ドキュメント"]
        repositories[repositories]
        managed_files[managed_files]
        documents[documents]
        
        users -.-> repositories
        repositories --> managed_files
        managed_files --> documents
        oauth_connections -.-> documents
    end
    
    subgraph m005["000005: バージョン管理・タグ"]
        document_versions[document_versions]
        document_current_version[document_current_version]
        tags[tags]
        document_tags[document_tags]
        
        documents --> document_versions
        documents --> document_current_version
        document_versions --> document_current_version
        oauth_connections -.-> document_versions
        document_versions --> document_tags
        tags --> document_tags
    end
    
    subgraph m006["000006: 実行記録"]
        execution_records[execution_records]
        execution_steps[execution_steps]
        attachments[attachments]
        
        documents -.-> execution_records
        document_versions -.-> execution_records
        users -.-> execution_records
        execution_records --> execution_steps
        execution_records --> attachments
        execution_steps --> attachments
        users -.-> attachments
    end
    
    subgraph m007["000007: 閲覧履歴・統計"]
        view_history[view_history]
        view_statistics[view_statistics]
        
        documents -.-> view_history
        documents --> view_statistics
        users -.-> view_history
    end
    
    subgraph m008["000008: グループ（未実装）"]
        groups[groups]
        user_groups[user_groups]
        
        users -.-> user_groups
        groups --> user_groups
    end
    
    style m001 fill:#e1f5ff
    style m002 fill:#e8f5e9
    style m003 fill:#fff3e0
    style m004 fill:#f3e5f5
    style m005 fill:#fce4ec
    style m006 fill:#fff9c4
    style m007 fill:#e0f2f1
    style m008 fill:#f5f5f5
```

**凡例**:
- 実線矢印 (→): 同じマイグレーション内の外部キー関係
- 点線矢印 (-.->): 前のマイグレーションで定義されたテーブルへの外部キー関係

## リレーション詳細図

```mermaid
erDiagram
    %% 000001: 認証基盤
    users ||--o{ refresh_tokens : "has"
    users ||--o{ user_identities : "has"
    users ||--|| user_profiles : "has profile"
    users ||--|| primary_user_identity : "has primary"
    user_identities ||--o{ primary_user_identity : "selected as"
    
    %% 000002-003: Gitプロバイダー・OAuth
    git_providers ||--o{ oauth_connections : "configured for"
    users ||--o{ oauth_connections : "has"
    
    %% 000004: リポジトリ・ドキュメント
    users ||--o{ repositories : "created"
    repositories ||--o{ managed_files : "contains"
    managed_files ||--o{ documents : "bound to"
    oauth_connections ||--o{ documents : "manages"
    users ||--o{ documents : "created"
    
    %% 000005: バージョン管理・タグ
    documents ||--o{ document_versions : "has versions"
    documents ||--|| document_current_version : "current version"
    document_versions ||--|| document_current_version : "stored in"
    oauth_connections ||--o{ document_versions : "fetched by"
    tags ||--o{ document_tags : "applied to"
    document_versions ||--o{ document_tags : "tagged with"
    
    %% 000006: 実行記録
    documents ||--o{ execution_records : "executed"
    document_versions ||--o{ execution_records : "executed version"
    users ||--o{ execution_records : "executor"
    execution_records ||--o{ execution_steps : "contains"
    execution_records ||--o{ attachments : "has"
    execution_steps ||--o{ attachments : "attached to"
    users ||--o{ attachments : "uploaded"
    
    %% 000007: 閲覧履歴・統計
    documents ||--o{ view_history : "viewed"
    documents ||--|| view_statistics : "stats"
    users ||--o{ view_history : "viewer"
    
    %% 000008: グループ（未実装）
    groups ||--o{ user_groups : "contains"
    users ||--o{ user_groups : "member of"
    users {
        UUID id PK
        VARCHAR primary_email UK
        VARCHAR display_name
        TEXT picture_url
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
        TIMESTAMPTZ last_login_at
    }
    
    refresh_tokens {
        UUID id PK
        UUID user_id FK
        VARCHAR token_hash UK
        VARCHAR jti UK
        TIMESTAMPTZ expires_at
        TIMESTAMPTZ created_at
        TIMESTAMPTZ last_used_at
        TIMESTAMPTZ revoked_at
        BOOLEAN is_revoked
        BOOLEAN remember_me
    }
    
    user_identities {
        UUID id PK
        UUID user_id FK
        VARCHAR provider
        VARCHAR provider_user_id UK
        VARCHAR email
        VARCHAR name
        TEXT picture_url
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
        TIMESTAMPTZ last_used_at
    }
    
    groups {
        UUID id PK
        VARCHAR name UK
        TEXT description
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    
    user_groups {
        UUID user_id PK "FK"
        UUID group_id PK "FK"
        TIMESTAMPTZ joined_at
    }
    
    oauth_connections {
        UUID id PK
        UUID user_id FK
        VARCHAR provider
        VARCHAR provider_host
        JSONB provider_metadata
        TEXT access_token_encrypted
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    
    repositories {
        UUID id PK
        VARCHAR provider_repository_id UK
        VARCHAR provider_repository_owner
        VARCHAR provider_repository_name
        VARCHAR url UK
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    
    managed_files {
        UUID id PK
        UUID repository_id FK
        TEXT file_path
        VARCHAR branch
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    
    documents {
        UUID id PK
        UUID managed_file_id FK
        UUID connection_id FK
        BOOLEAN is_auto_update
        BOOLEAN enabled
        UUID created_by FK
        TIMESTAMPTZ created_at
    }
    
    document_versions {
        UUID id PK
        UUID document_id FK
        INTEGER version_no
        VARCHAR commit_hash
        UUID connection_by FK
        TIMESTAMPTZ created_at
    }
    
    document_current_version {
        UUID document_id PK "FK"
        UUID document_version_id FK
        TEXT title
        TEXT description
        VARCHAR document_type
        TEXT content
    }
    
    tags {
        UUID id PK
        VARCHAR name UK
    }
    
    document_tags {
        UUID document_version_id PK "FK"
        UUID tag_id PK "FK"
    }
    
    execution_records {
        UUID id PK
        UUID document_id FK
        UUID document_version_id FK
        UUID executor_id FK
        VARCHAR title
        JSONB variable_values
        TEXT notes
        VARCHAR status
        VARCHAR access_scope
        TIMESTAMPTZ started_at
        TIMESTAMPTZ completed_at
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    
    execution_steps {
        UUID id PK
        UUID execution_record_id FK
        INTEGER step_number
        TEXT description
        TEXT notes
        TIMESTAMPTZ executed_at
    }
    
    attachments {
        UUID id PK
        UUID execution_record_id FK
        UUID execution_step_id FK
        VARCHAR file_name
        BIGINT file_size
        VARCHAR mime_type
        VARCHAR storage_type
        TEXT storage_path
        UUID uploaded_by FK
        TIMESTAMPTZ uploaded_at
    }
    
    view_history {
        UUID id PK
        UUID document_id FK
        UUID user_id FK
        INET ip_address
        TEXT user_agent
        TIMESTAMPTZ viewed_at
    }
    
    view_statistics {
        UUID document_id PK "FK"
        BIGINT total_views
        INTEGER unique_users
        TIMESTAMPTZ last_viewed_at
        TIMESTAMPTZ updated_at
    }
```

## 領域別ER図（詳細）

### 1. 認証・認可システム（000001）

```mermaid
erDiagram
    users ||--o{ refresh_tokens : "1:N"
    users ||--o{ user_identities : "1:N"
    users ||--o| user_profiles : "1:0..1"
    users ||--|| primary_user_identity : "1:1"
    user_identities ||--o{ primary_user_identity : "1:N"
    
    users {
        UUID id PK
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
        TIMESTAMPTZ last_login_at
    }
    
    refresh_tokens {
        UUID id PK
        UUID user_id FK
        VARCHAR token_hash UK
        VARCHAR jti UK
        TIMESTAMPTZ expires_at
        BOOLEAN is_revoked
        BOOLEAN remember_me
    }
    
    user_identities {
        UUID id PK
        UUID user_id FK
        VARCHAR provider
        VARCHAR provider_user_id
        VARCHAR email
        VARCHAR name
        TEXT picture_url
    }
    
    user_profiles {
        UUID user_id PK "FK"
        VARCHAR display_name
        TEXT picture_url
        TEXT bio
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    
    primary_user_identity {
        UUID user_id PK "FK"
        UUID identity_id FK
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
```

### 2. Gitプロバイダー管理（000002）

```mermaid
erDiagram
    git_providers {
        UUID id PK
        VARCHAR provider_type
        VARCHAR provider_name UK
        VARCHAR base_url
        VARCHAR api_url
        BOOLEAN is_saas
        VARCHAR client_id
        TEXT client_secret_encrypted
        BOOLEAN enabled
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
```

### 3. OAuth接続管理（000003）

```mermaid
erDiagram
    users ||--o{ oauth_connections : "1:N"
    git_providers ||--o{ oauth_connections : "1:N"
    
    oauth_connections {
        UUID id PK
        UUID user_id FK
        UUID git_provider_id FK
        VARCHAR provider
        VARCHAR provider_host
        JSONB provider_metadata
        TEXT access_token_encrypted
        TEXT refresh_token_encrypted
        TIMESTAMPTZ token_expires_at
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
        TIMESTAMPTZ last_used_at
    }
```

### 4. リポジトリ・ファイル管理（000004）

```mermaid
erDiagram
    users ||--o{ repositories : "1:N"
    repositories ||--o{ managed_files : "1:N"
    managed_files ||--o{ documents : "1:N"
    oauth_connections ||--o{ documents : "1:N"
    users ||--o{ documents : "1:N"
    
    repositories {
        UUID id PK
        VARCHAR provider_repository_id UK
        VARCHAR provider_repository_owner
        VARCHAR provider_repository_name
        VARCHAR url UK
        UUID created_by FK
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    
    managed_files {
        UUID id PK
        UUID repository_id FK
        TEXT file_path
        VARCHAR branch
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    
    documents {
        UUID id PK
        UUID managed_file_id FK
        UUID connection_id FK
        BOOLEAN is_auto_update
        BOOLEAN enabled
        UUID created_by FK
        TIMESTAMPTZ created_at
    }
```

### 5. ドキュメントバージョン管理（000005）

```mermaid
erDiagram
    documents ||--o{ document_versions : "1:N"
    documents ||--|| document_current_version : "1:1"
    document_versions ||--|| document_current_version : "1:1"
    document_versions ||--o{ document_tags : "N:M"
    tags ||--o{ document_tags : "N:M"
    oauth_connections ||--o{ document_versions : "1:N"
    
    document_versions {
        UUID id PK
        UUID document_id FK
        INTEGER version_no
        VARCHAR commit_hash
        UUID connection_by FK
        TIMESTAMPTZ created_at
    }
    
    document_current_version {
        UUID document_id PK "FK"
        UUID document_version_id FK
        TEXT title
        TEXT description
        VARCHAR document_type
        TEXT content
    }
    
    tags {
        UUID id PK
        VARCHAR name UK
    }
    
    document_tags {
        UUID document_version_id PK "FK"
        UUID tag_id PK "FK"
    }
```

### 6. 実行記録システム（000006）
    documents {
        UUID id PK
        UUID managed_file_id FK
        UUID connection_id FK
        BOOLEAN is_auto_update
        BOOLEAN enabled
        UUID created_by FK
        TIMESTAMPTZ created_at
    }
    
    document_versions {
        UUID id PK
        UUID document_id FK
        INTEGER version_no
        VARCHAR commit_hash
        UUID connection_by FK
        TIMESTAMPTZ created_at
    }
    
    document_current_version {
        UUID document_id PK "FK"
        UUID document_version_id FK
        TEXT title
        TEXT description
        VARCHAR document_type
        TEXT content
    }
    
    tags {
        UUID id PK
        VARCHAR name UK
    }
    
    document_tags {
        UUID document_version_id PK "FK"
        UUID tag_id PK "FK"
    }
```

### 6. 実行記録システム（000006）

```mermaid
erDiagram
    documents ||--o{ execution_records : "1:N"
    document_versions ||--o{ execution_records : "1:N"
    users ||--o{ execution_records : "1:N"
    execution_records ||--o{ execution_steps : "1:N"
    execution_records ||--o{ attachments : "1:N"
    execution_steps ||--o{ attachments : "1:N"
    users ||--o{ attachments : "1:N"
    
    execution_records {
        UUID id PK
        UUID document_id FK
        UUID document_version_id FK
        UUID executor_id FK
        VARCHAR title
        JSONB variable_values
        TEXT notes
        VARCHAR status
        VARCHAR access_scope
        TIMESTAMPTZ started_at
        TIMESTAMPTZ completed_at
    }
    
    execution_steps {
        UUID id PK
        UUID execution_record_id FK
        INTEGER step_number
        TEXT description
        TEXT notes
        TIMESTAMPTZ executed_at
    }
    
    attachments {
        UUID id PK
        UUID execution_record_id FK
        UUID execution_step_id FK
        VARCHAR file_name
        BIGINT file_size
        VARCHAR mime_type
        VARCHAR storage_type
        TEXT storage_path
        UUID uploaded_by FK
        TIMESTAMPTZ uploaded_at
    }
```

### 7. 閲覧履歴・統計システム（000007）
    }
    
    attachments {
        UUID id PK
        UUID execution_record_id FK
        UUID execution_step_id FK
        VARCHAR file_name
        BIGINT file_size
        VARCHAR mime_type
        VARCHAR storage_type
        TEXT storage_path
        UUID uploaded_by FK
        TIMESTAMPTZ uploaded_at
    }
```

### 7. 閲覧履歴・統計システム（000007）

```mermaid
erDiagram
    documents ||--o{ view_history : "1:N"
    documents ||--|| view_statistics : "1:1"
    users ||--o{ view_history : "1:N"
    
    view_history {
        UUID id PK
        UUID document_id FK
        UUID user_id FK
        INET ip_address
        TEXT user_agent
        TIMESTAMPTZ viewed_at
    }
    
    view_statistics {
        UUID document_id PK "FK"
        BIGINT total_views
        INTEGER unique_users
        TIMESTAMPTZ last_viewed_at
        TIMESTAMPTZ updated_at
    }
```

### 8. グループ管理（000008）- 未実装

```mermaid
erDiagram
    users ||--o{ user_groups : "N:M"
    groups ||--o{ user_groups : "N:M"
    
    groups {
        UUID id PK
        VARCHAR name UK
        TEXT description
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    
    user_groups {
        UUID user_id PK "FK"
        UUID group_id PK "FK"
        TIMESTAMPTZ joined_at
    }
```

## テーブル一覧

| #   | テーブル名               | 用途                                     | マイグレーション |
| --- | ------------------------ | ---------------------------------------- | ---------------- |
| 1   | users                    | ユーザー管理（認証のみ）                 | 000001           |
| 2   | refresh_tokens           | リフレッシュトークン管理                 | 000001           |
| 3   | user_identities          | 外部プロバイダー認証                     | 000001           |
| 4   | user_profiles            | ユーザープロフィール（カスタマイズ可能） | 000001           |
| 5   | primary_user_identity    | デフォルトidentity選択                   | 000001           |
| 6   | git_providers            | Gitプロバイダー設定                      | 000002           |
| 7   | oauth_connections        | OAuth接続管理                            | 000003           |
| 8   | repositories             | リポジトリ管理                           | 000004           |
| 9   | managed_files            | 管理対象ファイル                         | 000004           |
| 10  | documents                | ドキュメント論理管理                     | 000004           |
| 11  | document_versions        | ドキュメントバージョン管理               | 000005           |
| 12  | document_current_version | ドキュメント最新バージョン               | 000005           |
| 13  | tags                     | タグマスター                             | 000005           |
| 14  | document_tags            | ドキュメント・タグ中間テーブル           | 000005           |
| 15  | execution_records        | 実行記録                                 | 000006           |
| 16  | execution_steps          | 実行ステップ                             | 000006           |
| 17  | attachments              | 添付ファイル                             | 000006           |
| 18  | view_history             | 閲覧履歴                                 | 000007           |
| 19  | view_statistics          | 閲覧統計                                 | 000007           |
| 20  | groups                   | グループ管理                             | 000008（未実装） |
| 21  | user_groups              | ユーザー・グループ中間テーブル           | 000008（未実装） |

**合計: 21テーブル**
