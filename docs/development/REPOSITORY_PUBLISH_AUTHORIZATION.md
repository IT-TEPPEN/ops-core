# Repository Publish Authorization Implementation Guide

このドキュメントは、ADR 0051「Document-Centric Repository Management and Publish Authorization」の実装ガイドです。

## 概要

OpsCoreでは、外部Gitリポジトリからファイルを公開する際に、リポジトリMaintainer/Ownerの承認を得る仕組みを段階的に実装します。

## 設計原則

### 1. Git権限とOpsCore権限の分離

```
Git側の権限:
  - ファイルの読み取り・編集権限
  - OpsCoreはこれをバイパスしない

OpsCore側の権限:
  - 公開されたドキュメントの閲覧権限
  - access_scope (public/private) で制御
```

### 2. ドキュメント中心設計

- リポジトリの事前登録は不要
- ファイル公開時にリポジトリ情報を自動記録
- `repositories` テーブルは正規化・統計用

### 3. 段階的実装

| Phase | 期間 | 内容 |
|-------|------|------|
| Phase 1 | 1-2週間 | ロールチェック（MVP） |
| Phase 2 | 1-2ヶ月 | 承認フロー |
| Phase 3 | 3-6ヶ月 | メタデータベース（オプション） |

## Phase 1: ロールチェック（MVP）

### 目標

- Maintainer/Ownerのみ直接公開可能
- それ以外のユーザーには警告表示

### データフロー

```
1. ユーザーが「公開」ボタンをクリック
2. バックエンドがGit APIでロールを確認
3. admin/maintain → 公開許可
4. write/read → エラー「Maintainer/Ownerの承認が必要です」
```

### バックエンド実装

#### 1. Git Client インターフェース追加

```go
// internal/git_repository/domain/repository/git_client.go
type GitClient interface {
    FetchFile(ctx context.Context, token, fileURL string) (content, commitHash string, err error)

    // 追加: ユーザーのリポジトリロールを取得
    GetUserRole(ctx context.Context, token, repoURL string) (string, error)
    // 戻り値: "admin" | "maintain" | "write" | "read" | "none"
}
```

#### 2. GitHub/GitLab実装

```go
// internal/git_repository/infrastructure/external/github_client.go
func (c *GitHubClient) GetUserRole(ctx context.Context, token, repoURL string) (string, error) {
    owner, repo := extractOwnerRepo(repoURL)
    username := c.getAuthenticatedUsername(ctx, token)

    // GitHub API: GET /repos/{owner}/{repo}/collaborators/{username}/permission
    endpoint := fmt.Sprintf("https://api.github.com/repos/%s/%s/collaborators/%s/permission",
        owner, repo, username)

    req, _ := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Accept", "application/vnd.github+json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return "", fmt.Errorf("failed to get user role: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode == 404 {
        return "none", nil  // アクセス権なし
    }

    var result struct {
        Permission string `json:"permission"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", err
    }

    return result.Permission, nil
}
```

```go
// internal/git_repository/infrastructure/external/gitlab_client.go
func (c *GitLabClient) GetUserRole(ctx context.Context, token, repoURL string) (string, error) {
    projectID := extractProjectID(repoURL)

    // GitLab API: GET /projects/{id}/members/all?user_ids[]={user_id}
    endpoint := fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/members/all", projectID)

    // ... 実装

    // GitLab access_level → 統一ロールにマッピング
    // 50 (Owner) → "admin"
    // 40 (Maintainer) → "maintain"
    // 30 (Developer) → "write"
    // 20 (Reporter) → "read"
    // 10 (Guest) → "read"
}
```

#### 3. Usecase更新

```go
// internal/document/application/usecase/publish_document.go
type PublishDocumentUsecase struct {
    docRepo     domain.DocumentRepository
    repoRepo    gitrepo.RepositoryRepository
    gitClient   gitrepo.GitClient
    oauthRepo   oauth.ConnectionRepository
    // ...
}

func (uc *PublishDocumentUsecase) Execute(ctx context.Context, req PublishDocumentRequest) error {
    // 1. OAuth接続取得
    token, err := uc.oauthRepo.GetUserToken(ctx, req.UserID, req.Provider)
    if err != nil {
        return fmt.Errorf("OAuth token not found: %w", err)
    }

    // 2. ユーザーロール確認（Phase 1の追加部分）
    repoURL := extractRepositoryURL(req.FileURL)
    userRole, err := uc.gitClient.GetUserRole(ctx, token, repoURL)
    if err != nil {
        return fmt.Errorf("failed to get user role: %w", err)
    }

    // 3. 権限チェック
    if !canPublish(userRole) {
        return ErrPublishNotAuthorized{
            UserRole: userRole,
            Message: "リポジトリのMaintainer/Ownerのみがドキュメントを公開できます",
        }
    }

    // 4. ファイル取得・公開処理（既存ロジック）
    // ...
}

func canPublish(role string) bool {
    return role == "admin" || role == "maintain"
}
```

#### 4. エラー定義

```go
// internal/document/application/error/errors.go
type ErrPublishNotAuthorized struct {
    UserRole string
    Message  string
}

func (e ErrPublishNotAuthorized) Error() string {
    return fmt.Sprintf("PUBLISH_NOT_AUTHORIZED: %s (role: %s)", e.Message, e.UserRole)
}

func (e ErrPublishNotAuthorized) Code() string {
    return "APP_PUBLISH_NOT_AUTHORIZED"
}
```

### フロントエンド実装

#### 1. API型定義

```typescript
// features/document/infrastructure/services/types.ts
interface PublishDocumentErrorResponse {
  error: {
    code: 'APP_PUBLISH_NOT_AUTHORIZED';
    message: string;
    details?: {
      user_role: 'write' | 'read' | 'none';
    };
  };
}
```

#### 2. エラーハンドリング

```typescript
// features/document/presentation/hooks/usePublishDocument.ts
export function usePublishDocument() {
  const mutation = useMutation({
    mutationFn: async (data: PublishDocumentInput) => {
      return await documentCommandService.publish(data);
    },
    onError: (error: ApiError) => {
      if (error.code === 'APP_PUBLISH_NOT_AUTHORIZED') {
        toast.error(
          'ドキュメントを公開するには、リポジトリのMaintainer/Owner権限が必要です。' +
          '\n\n権限がない場合は、Maintainerに承認をリクエストできます。',
          {
            duration: 8000,
            action: {
              label: '承認をリクエスト',
              onClick: () => {
                // Phase 2で実装
                console.log('承認リクエスト機能は開発中です');
              },
            },
          }
        );
      }
    },
  });

  return mutation;
}
```

#### 3. UI表示

```typescript
// features/document/presentation/components/PublishButton.tsx
export function PublishButton({ fileUrl }: Props) {
  const { data: userRole, isLoading } = useGitUserRole(fileUrl);
  const { mutate: publish } = usePublishDocument();

  if (isLoading) {
    return <Button disabled>確認中...</Button>;
  }

  // Maintainer/Ownerなら直接公開可能
  if (userRole === 'admin' || userRole === 'maintain') {
    return (
      <Button onClick={() => publish({ fileUrl })}>
        公開する
      </Button>
    );
  }

  // それ以外は警告 + 今後の承認リクエスト機能を示唆
  return (
    <div>
      <Alert variant="warning">
        <AlertTitle>権限不足</AlertTitle>
        <AlertDescription>
          このリポジトリのMaintainer/Owner権限が必要です。
          <br />
          現在のロール: <Badge>{userRole}</Badge>
        </AlertDescription>
      </Alert>

      <Button disabled className="mt-2">
        承認をリクエスト（近日実装予定）
      </Button>
    </div>
  );
}
```

### テスト

#### バックエンドテスト

```go
// internal/document/application/usecase/publish_document_test.go
func TestPublishDocument_RoleCheck(t *testing.T) {
    tests := []struct {
        name        string
        userRole    string
        expectError bool
    }{
        {"admin can publish", "admin", false},
        {"maintainer can publish", "maintain", false},
        {"writer cannot publish", "write", true},
        {"reader cannot publish", "read", true},
        {"no access cannot publish", "none", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Mock setup
            mockGitClient := &MockGitClient{
                GetUserRoleFunc: func(ctx context.Context, token, repoURL string) (string, error) {
                    return tt.userRole, nil
                },
            }

            uc := NewPublishDocumentUsecase(mockDocRepo, mockRepoRepo, mockGitClient, mockOAuthRepo)

            err := uc.Execute(context.Background(), PublishDocumentRequest{
                UserID:   "user-123",
                FileURL:  "https://github.com/org/repo/blob/main/doc.md",
                Provider: "github",
            })

            if tt.expectError {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), "PUBLISH_NOT_AUTHORIZED")
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

#### フロントエンドテスト

```typescript
// features/document/presentation/components/PublishButton.test.tsx
describe('PublishButton', () => {
  it('shows publish button for maintainer', () => {
    mockUseGitUserRole.mockReturnValue({ data: 'maintain', isLoading: false });

    render(<PublishButton fileUrl="https://github.com/org/repo/blob/main/doc.md" />);

    expect(screen.getByRole('button', { name: '公開する' })).toBeEnabled();
  });

  it('shows warning for writer', () => {
    mockUseGitUserRole.mockReturnValue({ data: 'write', isLoading: false });

    render(<PublishButton fileUrl="https://github.com/org/repo/blob/main/doc.md" />);

    expect(screen.getByText(/権限不足/i)).toBeInTheDocument();
    expect(screen.getByText(/現在のロール: write/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /承認をリクエスト/i })).toBeDisabled();
  });
});
```

---

## Phase 2: 承認フロー（詳細は後日追加）

### データベーススキーマ

```sql
-- migration: 000XXX_create_publish_requests.up.sql

CREATE TABLE document_publish_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_url TEXT NOT NULL,
    requested_by UUID NOT NULL REFERENCES users(id),
    repository_url TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('pending', 'approved', 'rejected')),
    approval_token TEXT UNIQUE NOT NULL,
    requested_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,

    INDEX idx_publish_requests_status (status),
    INDEX idx_publish_requests_token (approval_token),
    INDEX idx_publish_requests_requested_by (requested_by)
);

CREATE TABLE document_publish_approvers (
    request_id UUID NOT NULL REFERENCES document_publish_requests(id) ON DELETE CASCADE,
    approver_email TEXT NOT NULL,
    approver_github_login TEXT,
    approved_at TIMESTAMP,

    PRIMARY KEY (request_id, approver_email)
);
```

### API エンドポイント

```
POST   /api/documents/publish-requests     # 承認リクエスト作成
GET    /api/documents/publish-requests     # 自分のリクエスト一覧
GET    /api/documents/publish-requests/:id # リクエスト詳細
POST   /api/approve/:token                 # 承認（メールリンク用）
POST   /api/reject/:token                  # 拒否（メールリンク用）
```

### 実装詳細

（Phase 2実装時に詳細を追加）

---

## Phase 3: メタデータベース（詳細は後日追加）

### .opscore.yml 仕様

```yaml
version: 1
policies:
  - path: "docs/procedures/*.md"
    publish: allowed
    default_scope: public
    allowed_publishers:  # オプション
      - github:alice
      - github:bob
```

### 実装詳細

（Phase 3実装時に詳細を追加）

---

## セキュリティ考慮事項

### 1. トークン管理

- OAuth接続のアクセストークンは暗号化保存（AES-256-GCM）
- 承認トークンは UUID v4 でワンタイム使用
- 承認トークンは7日で自動失効

### 2. レートリミット

- Git API呼び出しにレートリミット適用
- ロールチェックは結果をキャッシュ（5分）

### 3. 監査ログ

- 全ての公開リクエスト・承認・拒否を記録
- `audit_logs` テーブルに保存

```sql
INSERT INTO audit_logs (action, user_id, resource_type, resource_id, metadata)
VALUES ('DOCUMENT_PUBLISHED', $1, 'document', $2, jsonb_build_object(
    'file_url', $3,
    'repository_url', $4,
    'user_role', $5
));
```

---

## トラブルシューティング

### 問題: ロールチェックが失敗する

**症状**: `failed to get user role: 403 Forbidden`

**原因**: OAuth接続のスコープ不足

**解決策**:
1. GitHub: `repo` スコープが必要
2. GitLab: `read_api` スコープが必要
3. OAuth接続を再作成して正しいスコープを付与

### 問題: Maintainerなのに公開できない

**症状**: `PUBLISH_NOT_AUTHORIZED` エラー

**デバッグ**:
```bash
# バックエンドログで確認
grep "GetUserRole" /var/log/opscore/app.log

# 期待される出力
{"level":"debug","msg":"GetUserRole","repo":"org/repo","role":"maintain"}
```

**原因**:
- Git APIのレスポンスがキャッシュされている
- OAuth接続が古い

**解決策**:
- OAuth接続を再作成
- キャッシュをクリア

---

## 参考リンク

- [ADR 0051: Document-Centric Repository Management](/adr/0051-document-centric-repository-management.md)
- [GitHub API: Collaborators](https://docs.github.com/en/rest/collaborators/collaborators)
- [GitLab API: Project members](https://docs.gitlab.com/ee/api/members.html)
- [ADR 0002: Repository Access Method (OAuth 2.0)](/adr/0002-repository-access-method.md)
- [ADR 0050: User Authentication (Multi-Provider OIDC)](/adr/0050-user-authentication-google-oidc.md)
