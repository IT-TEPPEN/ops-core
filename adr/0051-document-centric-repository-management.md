# ADR 0051: Document-Centric Repository Management and Publish Authorization

## Status

Proposed

## Context

OpsCoreは、外部Gitリポジトリ（GitHub/GitLab）からMarkdownファイルを取得し、独自機能（変数埋め込み、実行記録、証跡管理）を提供するシステムである。当初の設計では「リポジトリを事前登録してからファイルを公開する」フローを想定していたが、以下の課題が明らかになった：

### 課題1: リポジトリ事前登録の必要性が不明確

- ファイル公開時にURLからリポジトリ情報を自動抽出可能
- ユーザーにとって「リポジトリ登録」という概念が不自然
- コアバリュー（手順書公開・実行・証跡管理）には必須ではない

### 課題2: 権限モデルの危険性

現在の設計では、以下のようなセキュリティリスクが存在する：

```
シナリオ:
1. 新人エンジニアAさん: プライベートリポジトリ `company/secrets` に読み取り権限あり
2. AさんがOpsCoreで `secrets/sensitive-procedure.md` を「全員公開」
3. リポジトリOwnerの意図に反して機密情報が漏洩
```

**問題**: 「Git側で参照権限がある」≠「OpsCoreで公開する権限がある」

### 課題3: チームコラボレーションの考慮不足

- 運用業務は個人ではなくチームで実施される
- 個人のOAuth接続のみでは属人化が発生
- リポジトリの「チーム資産」としての管理が必要

## Decision

### 決定1: ドキュメント中心設計への変更

リポジトリを事前登録する概念を削除し、ファイル公開時に自動的にリポジトリ情報を記録する。

#### 変更前（リポジトリ事前登録）
```
1. リポジトリを登録（repositories テーブルに INSERT）
2. ファイルを選択
3. 公開（repository_id を指定）
```

#### 変更後（自動記録）
```
1. OAuth接続
2. ファイルをブラウズ・選択
3. 公開ボタン
   → ファイルURLからリポジトリ情報を自動抽出
   → repositories テーブルに自動UPSERT
   → documents に保存
```

#### データモデル

**repositoriesテーブル**: 「参照されたリポジトリの正規化テーブル」（統計・集約用）

```sql
repositories:
  - id UUID PRIMARY KEY
  - url TEXT UNIQUE NOT NULL          -- https://github.com/org/repo
  - name TEXT NOT NULL                -- org/repo
  - provider TEXT NOT NULL            -- github | gitlab
  - first_referenced_at TIMESTAMP
  - first_referenced_by UUID          -- users.id
```

**documentsテーブル**: 非正規化カラムを追加（高速アクセス用）

```sql
documents:
  - id UUID PRIMARY KEY
  - title TEXT
  - content TEXT
  - file_url TEXT NOT NULL           -- 元のGitファイルURL
  - repository_id UUID               -- repositories.id (JOIN用、nullable)
  - repository_url TEXT              -- 非正規化（高速アクセス用）
  - repository_name TEXT             -- 非正規化（表示用）
  - commit_hash TEXT
  - access_scope TEXT                -- public | private (OpsCore内の権限)
  - published_by UUID
  - published_at TIMESTAMP
```

### 決定2: 権限モデルの明確化

Git側の権限とOpsCore側の権限を明確に分離する。

#### Git側の権限（変更しない）

- OpsCoreはGit側の権限をバイパス・上書きしない
- ユーザーのOAuth接続で「アクセス可能なファイル」のみ取得可能

#### OpsCore側の権限（独立）

- 公開された「ドキュメント」の閲覧権のみ管理
- `access_scope`（public/private）で制御
- Git側の権限とは独立した、OpsCore内の閲覧制御

#### 具体例

```
ケース: プライベートリポジトリのファイルを公開

1. Aさん: GitHub `company/secrets` に読み取り権限あり
2. AさんがOpsCoreで `secrets/procedure.md` を取得
3. Aさんが「このファイルを公開」→ OpsCoreに保存
4. access_scope = 'public' で公開
5. Bさん（Git側では権限なし）がOpsCoreでドキュメントを閲覧可能

結果:
  - Git側: Bさんはリポジトリにアクセス不可（変わらず）
  - OpsCore側: Bさんはドキュメントを閲覧可能（Aさんが公開したため）
```

これはOpsCoreが「公開済みコンテンツの配信プラットフォーム」であり、元のGit権限とは独立しているため問題ない。

### 決定3: 公開承認フローの段階的実装

リポジトリMaintainer/Ownerの承認を得てから公開する仕組みを段階的に実装する。

#### Phase 1: ロールチェック（MVP - 1-2週間）

**実装内容**:
- 公開時にGit APIでユーザーのロール（admin/maintain/write/read）を確認
- Maintainer/Ownerのみ直接公開可能
- それ以外のユーザーには警告表示 + 「承認リクエスト」ボタン

```typescript
async function canPublishDocument(userToken: string, repoUrl: string): Promise<boolean> {
  const userRole = await gitClient.getUserRole(userToken, repoUrl);
  // GitHub API: GET /repos/{owner}/{repo}/collaborators/{username}/permission

  return ['admin', 'maintain'].includes(userRole);
}
```

**メリット**:
- 実装が簡単（Git API呼び出しのみ）
- 既存のGit権限モデルをそのまま利用
- 最低限のセキュリティを即座に確保

**デメリット**:
- Maintainer/Ownerしか公開できない（柔軟性が低い）

#### Phase 2: 承認フロー（中期 - 1-2ヶ月）

**実装内容**:
- 承認リクエストテーブルの追加
- メール承認フロー
- ワンタイムトークンによる承認/拒否

**データベーススキーマ**:

```sql
-- 公開承認リクエスト
document_publish_requests:
  - id UUID PRIMARY KEY
  - file_url TEXT NOT NULL
  - requested_by UUID NOT NULL        -- users.id
  - repository_url TEXT NOT NULL
  - status TEXT NOT NULL              -- pending | approved | rejected
  - approval_token TEXT UNIQUE        -- ワンタイムトークン（UUID）
  - requested_at TIMESTAMP
  - expires_at TIMESTAMP              -- 7日後など

-- 承認者（リポジトリMaintainer/Owner）
document_publish_approvers:
  - request_id UUID
  - approver_email TEXT NOT NULL
  - approver_github_login TEXT        -- GitHub username
  - approved_at TIMESTAMP
  - PRIMARY KEY (request_id, approver_email)
```

**フロー**:

1. **ユーザーが公開リクエスト**:
   - リポジトリのMaintainer/Ownerを取得（GitHub API）
   - リクエストを保存
   - Maintainer/Ownerにメール送信（承認URL付き）

2. **Maintainer/Ownerが承認**:
   - メール内の承認URLをクリック
   - リクエストステータスを `approved` に更新
   - リクエスト者に通知メール送信

3. **リクエスト者が公開実行**:
   - 承認済みリクエストがあるか確認
   - なければMaintainer/Owner本人か確認
   - 公開処理を実行

**メリット**:
- 柔軟性が高い（Maintainer以外も公開可能）
- Git権限モデルを尊重
- メール承認なので特別なツール不要

**デメリット**:
- メール配信基盤が必要
- 承認フローの管理が必要

#### Phase 3: メタデータベース（長期 - 3-6ヶ月、オプション）

**実装内容**:
- リポジトリに `.opscore.yml` を配置
- 公開ポリシーを宣言的に管理

```yaml
# .opscore.yml (リポジトリルートに配置)
version: 1
policies:
  # パターン1: 特定ファイルを明示的に許可
  - path: "docs/procedures/*.md"
    publish: allowed
    default_scope: public

  # パターン2: ディレクトリ単位で制御
  - path: "runbooks/**/*.md"
    publish: allowed
    default_scope: private
    allowed_publishers:  # オプション: 特定ユーザーのみ許可
      - github:alice
      - github:bob

  # パターン3: 明示的に拒否
  - path: "internal/**/*"
    publish: denied
```

**メリット**:
- 宣言的で明確
- GitOps的なアプローチ（設定もGit管理）
- 承認フロー不要（Maintainerが `.opscore.yml` を編集するだけ）

**デメリット**:
- リポジトリにファイル追加が必要
- ユーザーがYAML設定を理解する必要がある

### 決定4: お気に入り機能（オプション）

頻繁に使うリポジトリをブックマークする軽量な機能を提供。

```sql
-- ユーザーの個人的なブックマーク
user_repository_bookmarks:
  - id UUID PRIMARY KEY
  - user_id UUID NOT NULL
  - repository_url TEXT NOT NULL
  - custom_name TEXT                 -- ユーザーが付けた名前（オプション）
  - created_at TIMESTAMP

  UNIQUE(user_id, repository_url)
```

これは単なるUI上のショートカットであり、権限管理や強制登録とは無関係。

## Consequences

### メリット

1. **UXの向上**:
   - リポジトリ登録という概念が不要
   - OAuth接続 → ファイル選択 → 公開のシンプルなフロー
   - お気に入り機能で頻繁に使うリポジトリへの高速アクセス

2. **セキュリティの向上**:
   - Git権限モデルを尊重（バイパスしない）
   - Maintainer/Ownerの承認を経て公開
   - 監査ログで追跡可能

3. **チームコラボレーション**:
   - リポジトリが「チーム資産」として認識される
   - 承認フローでチーム内の合意形成

4. **拡張性**:
   - 将来のWebhook機能追加時にrepositories テーブルを活用可能
   - メタデータベース（`.opscore.yml`）への移行パスが明確

### デメリット・リスク

1. **開発コスト**:
   - Phase 2の承認フロー実装に1-2ヶ月必要
   - メール配信基盤の構築が必要

2. **ユーザー教育**:
   - 「Git権限」と「OpsCore権限」の違いを理解してもらう必要
   - Phase 3のメタデータファイルはYAML知識が必要

3. **初期制限**:
   - Phase 1ではMaintainer/Ownerのみ公開可能
   - Phase 2実装までは柔軟性が低い

### マイグレーション

既存の `repositories` テーブルに登録済みデータがある場合:

```sql
-- Step 1: カラム追加（非正規化）
ALTER TABLE documents ADD COLUMN repository_url TEXT;
ALTER TABLE documents ADD COLUMN repository_name TEXT;

-- Step 2: 既存データを埋める
UPDATE documents d
SET
  repository_url = r.url,
  repository_name = r.name
FROM repositories r
WHERE d.repository_id = r.id;

-- Step 3: repositories テーブルのカラム変更
ALTER TABLE repositories RENAME COLUMN registered_by TO first_referenced_by;
ALTER TABLE repositories RENAME COLUMN registered_at TO first_referenced_at;

-- Step 4: URL に UNIQUE 制約追加
ALTER TABLE repositories ADD CONSTRAINT repositories_url_unique UNIQUE (url);
```

## Implementation Plan

### Phase 1: MVP（1-2週間）

- [ ] Git APIロールチェック実装
- [ ] Maintainer/Owner判定ロジック
- [ ] UI上の警告表示
- [ ] 「承認リクエスト機能は今後追加予定」メッセージ

### Phase 2: 承認フロー（1-2ヶ月）

- [ ] `document_publish_requests` テーブル追加
- [ ] `document_publish_approvers` テーブル追加
- [ ] マイグレーションファイル作成
- [ ] メール配信基盤構築
- [ ] 承認リクエストAPI実装
- [ ] 承認/拒否エンドポイント実装
- [ ] メールテンプレート作成
- [ ] UI実装（リクエストボタン、ステータス表示）

### Phase 3: メタデータベース（3-6ヶ月、オプション）

- [ ] `.opscore.yml` スキーマ定義
- [ ] ポリシーパーサー実装
- [ ] ポリシー検証ロジック
- [ ] UI実装（ポリシー設定サポート）
- [ ] ドキュメント作成

## References

- ADR 0002: Repository Access Method (OAuth 2.0)
- ADR 0050: User Authentication (Multi-Provider OIDC)
- GitHub API: [Check permissions](https://docs.github.com/en/rest/collaborators/collaborators#get-repository-permissions-for-a-user)
- GitLab API: [Project members](https://docs.gitlab.com/ee/api/members.html)

## Related Issues

- #90: 動作しない箇所のバグ修正（リポジトリ関連の設計見直し）
