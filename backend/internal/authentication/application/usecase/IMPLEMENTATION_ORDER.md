# Authentication Usecase Implementation Order

このドキュメントでは、`internal/authentication/application/usecase`の実装順序を定義します。

## 実装方針

- 全てのusecaseインターフェースのExecuteメソッドは `Execute(ctx context.Context, dto XXXRequest) (result, error)` の形式に統一
- 依存関係を考慮し、基盤となる機能から順に実装
- 各usecaseには対応するユニットテストを作成

---

## Phase 1: コア認証機能 (最優先) 🔴

システムの基盤となる機能で、他の全てのusecaseが依存します。

### 1. AuthenticateUser (authenticate_user.go)

**優先度**: 最高  
**依存**: UserRepository のみ

```go
Execute(ctx context.Context, dto dto.ProviderUserInfo) (*dto.AuthenticationResult, error)
```

**実装ポイント**:
- Provider情報からUser/Identityを作成または更新
- 新規ユーザー登録と既存ユーザー認証の両方を処理
- 初回ログイン時のユーザー作成ロジック
- LastLoginAtの更新

**テストシナリオ**:
- 新規ユーザーの作成
- 既存ユーザーの認証
- 無効なProvider情報の処理
- Identity情報の更新

---

### 2. CreateSession (create_session.go)

**優先度**: 最高  
**依存**: SessionRepository のみ

```go
Execute(ctx context.Context, dto dto.CreateSessionRequest) (*dto.SessionInfo, error)
```

**リクエストDTO**:
```go
type CreateSessionRequest struct {
    UserID      string      `json:"user_id"`
    SessionType SessionType `json:"session_type"`
}
```

**実装ポイント**:
- SessionTypeに応じた有効期限設定
  - SessionTypeShort: 7日間
  - SessionTypeLong: 30日間
- セッショントークンの生成とハッシュ化
- JTI (JWT ID) の生成

**テストシナリオ**:
- Short sessionの作成
- Long sessionの作成
- 無効なuserIDの処理
- トークン生成エラーの処理

---

### 3. GetAuthenticatedUser (get_authenticated_user.go)

**優先度**: 最高  
**依存**: SessionRepository, UserRepository

```go
Execute(ctx context.Context, dto dto.SessionIDRequest) (*dto.UserInfo, error)
```

**リクエストDTO**:
```go
type SessionIDRequest struct {
    SessionID string `json:"session_id"`
}
```

**実装ポイント**:
- セッションの有効性検証
  - 期限チェック (ExpiresAt)
  - Revoke状態チェック
  - IsValid() メソッドの利用
- ユーザー情報の取得と返却

**テストシナリオ**:
- 有効なセッションでのユーザー取得
- 期限切れセッションの処理
- Revokeされたセッションの処理
- 存在しないセッションの処理

---

## Phase 2: セッション管理 (高優先) 🟠

実運用で必須の機能です。

### 4. RefreshSession (refresh_session.go)

**優先度**: 高  
**依存**: SessionRepository のみ

```go
Execute(ctx context.Context, dto dto.SessionIDRequest) (*dto.SessionInfo, error)
```

**リクエストDTO**:
```go
type SessionIDRequest struct {
    SessionID string `json:"session_id"`
}
```

**実装ポイント**:
- LastUsedAtの更新
- セッションの有効性チェック
- 期限延長は行わない（ExpiresAtは変更しない）

**テストシナリオ**:
- 有効なセッションのリフレッシュ
- 期限切れセッションの処理
- Revokeされたセッションの処理

---

### 5. RevokeSession (revoke_session.go)

**優先度**: 高  
**依存**: SessionRepository のみ

```go
Execute(ctx context.Context, dto dto.SessionIDRequest) error
```

**リクエストDTO**:
```go
type SessionIDRequest struct {
    SessionID string `json:"session_id"`
}
```

**実装ポイント**:
- セッションのRevoke処理
- RevokedAtタイムスタンプの設定
- 既にRevokeされているセッションの冪等性

**テストシナリオ**:
- セッションのRevoke
- 既にRevokeされたセッションの処理
- 存在しないセッションの処理

---

## Phase 3: マルチID管理 (中優先) 🟡

複数プロバイダー連携機能です。

### 6. LinkNewIdentity (link_new_identity.go)

**優先度**: 中  
**依存**: UserRepository のみ

```go
Execute(ctx context.Context, dto dto.LinkIdentityRequest) (*dto.IdentityInfo, error)
```

**リクエストDTO**:
```go
type LinkIdentityRequest struct {
    UserID       string           `json:"user_id"`
    ProviderInfo ProviderUserInfo `json:"provider_info"`
}
```

**実装ポイント**:
- 既存Identity重複チェック
  - 同じユーザーへの重複チェック (domain_err.ErrDuplicateIdentity)
  - 他のユーザーへの重複チェック (application_err.ErrIdentityAlreadyLinked)
- 新しいIdentityの作成と追加
- Primaryでない場合は自動的に非Primary設定

**テストシナリオ**:
- 新しいIdentityの追加
- 同じProviderのIdentity重複エラー
- 他のユーザーに紐づくIdentityの処理
- 存在しないユーザーへのリンク試行

---

### 7. ListUserIdentities (list_user_identities.go)

**優先度**: 中  
**依存**: UserRepository のみ

```go
Execute(ctx context.Context, dto dto.UserIDRequest) ([]dto.IdentityInfo, error)
```

**リクエストDTO**:
```go
type UserIDRequest struct {
    UserID string `json:"user_id"`
}
```

**実装ポイント**:
- ユーザーの全Identityを取得
- Primary Identityの識別
- 空のリストも正常なレスポンス

**テストシナリオ**:
- 複数Identityを持つユーザーの取得
- 1つのIdentityのみのユーザー
- 存在しないユーザーの処理

---

### 8. SetPrimaryIdentity (set_primary_identity.go)

**優先度**: 中  
**依存**: UserRepository のみ

```go
Execute(ctx context.Context, dto dto.SetPrimaryIdentityRequest) error
```

**リクエストDTO**:
```go
type SetPrimaryIdentityRequest struct {
    UserID     string `json:"user_id"`
    IdentityID string `json:"identity_id"`
}
```

**実装ポイント**:
- 前のPrimaryを自動的にdemote
- 新しいPrimaryの設定
- Identity所有権の検証

**テストシナリオ**:
- Primary Identityの変更
- 他のユーザーのIdentity指定エラー
- 存在しないIdentityの指定
- 既にPrimaryのIdentityの指定（冪等性）

---

### 9. UnlinkIdentity (unlink_identity.go)

**優先度**: 中  
**依存**: UserRepository のみ

```go
Execute(ctx context.Context, dto dto.UnlinkIdentityRequest) error
```

**リクエストDTO**:
```go
type UnlinkIdentityRequest struct {
    UserID     string `json:"user_id"`
    IdentityID string `json:"identity_id"`
}
```

**実装ポイント**:
- 最低1つのIdentity保持チェック
- Primary Identity削除不可チェック
- Identity所有権の検証

**テストシナリオ**:
- 非PrimaryなIdentityの削除
- 最後のIdentity削除試行エラー
- Primary Identity削除試行エラー
- 他のユーザーのIdentity削除試行エラー

---

## Phase 4: プロフィール・セキュリティ (低優先) 🟢

ユーザー体験向上のための機能です。

### 10. UpdateUserProfile (update_user_profile.go)

**優先度**: 低  
**依存**: UserRepository のみ

```go
Execute(ctx context.Context, dto dto.UserProfile) (*dto.UserProfile, error)
```

**実装ポイント**:
- DisplayName、PictureURLの更新
- 空文字列の扱い（Primary Identityの値を使用）
- 更新後のプロフィールを返却

**テストシナリオ**:
- DisplayNameの更新
- PictureURLの更新
- 両方の更新
- 存在しないユーザーの処理

---

### 11. GetUserSessions (get_user_sessions.go)

**優先度**: 低  
**依存**: SessionRepository のみ

```go
Execute(ctx context.Context, dto dto.UserIDRequest) ([]dto.SessionInfo, error)
```

**リクエストDTO**:
```go
type UserIDRequest struct {
    UserID string `json:"user_id"`
}
```

**実装ポイント**:
- アクティブセッションのみフィルタ
  - IsValid() == true
  - IsRevoked() == false
- セッション情報のDTO変換

**テストシナリオ**:
- 複数のアクティブセッション取得
- 期限切れ・Revokeされたセッションの除外
- セッションがない場合

---

### 12. RevokeAllUserSessions (revoke_all_user_sessions.go)

**優先度**: 低  
**依存**: SessionRepository のみ

```go
Execute(ctx context.Context, dto dto.UserIDRequest) (int, error)
```

**リクエストDTO**:
```go
type UserIDRequest struct {
    UserID string `json:"user_id"`
}
```

**実装ポイント**:
- バルクRevoke操作
- Repository層の RevokeAllByUserID() メソッドを使用
- Revokeされたセッション数を返却

**テストシナリオ**:
- 複数セッションの一括Revoke
- セッションがない場合
- 既にRevokeされているセッションの扱い

---

## 実装時の注意事項

### 1. エラーハンドリング

- Domain層のエラーはそのまま伝播
- Application層固有のエラーは application_err を使用
- エラーメッセージには適切なコンテキスト情報を含める

### 2. DTO変換

- Entity → DTO 変換は各usecase内でヘルパー関数として実装
- 共通化が必要になった場合は `application/mapper/` パッケージに移行検討

### 3. トランザクション管理

- Usecase層ではトランザクションを意識しない
- Repository実装が内部でトランザクションを管理する前提

### 4. テスト戦略

- 各usecaseに対応する `*_test.go` を作成
- Repository はモックを使用
- 正常系、異常系、境界値のテストを含める

---

## 実装例パターン

### 基本構造

```go
package usecase

import (
    "context"
    "opscore/backend/internal/authentication/application/dto"
    "opscore/backend/internal/authentication/application/err"
    "opscore/backend/internal/authentication/domain/repository"
)

type authenticateUserImpl struct {
    userRepo repository.UserRepository
}

func NewAuthenticateUser(userRepo repository.UserRepository) AuthenticateUser {
    return &authenticateUserImpl{
        userRepo: userRepo,
    }
}

func (uc *authenticateUserImpl) Execute(ctx context.Context, dto dto.ProviderUserInfo) (*dto.AuthenticationResult, error) {
    // 1. Input validation
    if err := uc.validateInput(dto); err != nil {
        return nil, err
    }

    // 2. Business logic (use Domain entities)
    // ...

    // 3. Return DTO
    return &dto.AuthenticationResult{
        // ...
    }, nil
}

func (uc *authenticateUserImpl) validateInput(info dto.ProviderUserInfo) error {
    if info.Provider == "" {
        return application_err.NewMissingRequiredField("provider")
    }
    // ...
    return nil
}
```

---

## まとめ

実装順序:
1. **Phase 1** (AuthenticateUser → CreateSession → GetAuthenticatedUser): システムの基盤
2. **Phase 2** (RefreshSession → RevokeSession): セッション管理
3. **Phase 3** (LinkNewIdentity → ListUserIdentities → SetPrimaryIdentity → UnlinkIdentity): マルチID管理
4. **Phase 4** (UpdateUserProfile → GetUserSessions → RevokeAllUserSessions): プロフィール・セキュリティ

この順序で実装することで、依存関係の問題なく段階的に機能を構築できます。
