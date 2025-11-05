# ドメインモデル再設計サマリー

## 概要

このドキュメントは、Ops-Coreのドキュメント管理コンテキストにおけるドメインモデルの再設計をまとめたものです。

## 再設計の目的

既存のドメインモデルをユースケースベースで見直し、以下の新機能を追加しました：

1. **手順書の変数定義機能**
   - 手順書実行時にパラメータを入力できる仕組み
   - 変数定義（name, label, description, type, required, defaultValue）
   - Frontmatterでの変数定義、本文での変数参照（`{{variable_name}}`形式）

2. **作業証跡管理機能**
   - 手順書実行セッション単位での証跡記録
   - 各ステップへのメモと画面キャプチャ添付
   - ローカルファイルシステムまたはS3/MinIOへの添付ファイル保存
   - 実行者・管理者・共有先ユーザー/グループによる閲覧制御
   - 作業証跡の検索・フィルタリング

3. **ユーザー・グループ管理機能**
   - ユーザー管理（admin/userロール）
   - グループ管理とメンバーシップ
   - 公開範囲とアクセス制御

## 主要な変更点

### 新規エンティティ

- **ExecutionRecord（集約ルート）**: 手順書実行の証跡を管理
- **ExecutionStep**: 作業ステップごとの記録
- **Attachment**: 画面キャプチャなどの添付ファイル
- **User（集約ルート）**: ユーザー情報とグループ所属
- **Group（集約ルート）**: グループ情報とメンバー管理
- **ViewHistory**: ドキュメント閲覧履歴
- **ViewStatistics（集約ルート）**: ドキュメント閲覧統計

### 既存エンティティの拡張

- **Document**
  - `variables: [VariableDefinition]` を追加
  - `isAutoUpdate: boolean` を追加（自動更新設定）
  - `currentVersion: DocumentVersion` を追加
  - `versions: [DocumentVersion]` を追加（バージョン履歴）

- **DocumentVersion**（新規）
  - ドキュメントの特定バージョンを表現
  - `versionNumber: VersionNumber`（ファイルごとに1から連番）

### 新規値オブジェクト

- **VariableDefinition**: 変数定義（name, label, description, type, required, defaultValue）
- **VariableValue**: 変数の入力値（name, value）
- **StorageType**: ストレージタイプ（local, s3, minio）
- **IPAddress**: IPアドレス
- **UserAgent**: ユーザーエージェント
- **VersionNumber**: バージョン番号（連番）

## 追加されたADR

1. **ADR 0013: Document Variable Definition and Substitution**
   - 手順書の変数定義と置換の仕様
   - Frontmatterでの変数定義方法
   - 変数参照の構文（`{{variable_name}}`）
   - 変数入力UIの仕様

2. **ADR 0014: Execution Record and Evidence Management**
   - 作業証跡管理の仕様
   - ExecutionRecord、ExecutionStep、Attachmentの構造
   - ストレージ戦略（ローカル/S3）
   - アクセス制御と共有機能
   - 検索・フィルタリング機能

3. **ADR 0001の更新**
   - Frontmatterに`variables`フィールドを追加

## ドメインモデル図

以下の3つのMermaid図を作成しました：

1. **エンティティ関連図**: エンティティ間の関連を表現
2. **値オブジェクト関連図**: 主要な値オブジェクトの構造
3. **集約境界図**: 集約ルートと集約の境界を表現

## ユースケース更新

以下のカテゴリに分類し、番号を振り直しました：

- **リポジトリ管理**（4つ）
- **ドキュメント管理**（6つ）
- **閲覧管理**（2つ）
- **手順書実行と作業証跡**（6つ）
- **ユーザー・グループ管理**（3つ）

合計21のユースケースを整理しました。

## 集約設計

以下の6つの集約を定義しました：

1. **Repository集約**: リポジトリ情報の一貫性を保証
2. **Document集約**: ドキュメントとそのバージョン履歴の整合性を管理
3. **ExecutionRecord集約**: 作業証跡と関連ステップ・添付ファイルの整合性を保証
4. **User集約**: ユーザー情報とグループ所属の管理
5. **Group集約**: グループ情報とメンバー一覧の管理
6. **ViewStatistics集約**: 閲覧統計の集計と更新

## リポジトリインターフェース

以下の8つのリポジトリインターフェースを定義しました：

1. RepositoryRepository
2. DocumentRepository
3. ExecutionRecordRepository
4. UserRepository
5. GroupRepository
6. ViewHistoryRepository
7. ViewStatisticsRepository
8. AttachmentRepository（ストレージ抽象化）

## 実装上の考慮事項

### ストレージ抽象化

- ローカルファイルシステムとS3/MinIOの両方をサポート
- 環境変数による設定切り替え
- AttachmentRepositoryでストレージ層を抽象化

### アクセス制御

- ドキュメントと作業証跡の両方でAccessScopeを使用
- private/shared の2つのアクセスタイプ
- 管理者は全てのリソースにアクセス可能

### バージョン管理

- ドキュメントごとに独立したバージョン番号（1から連番）
- コミットハッシュとバージョン番号を紐付け
- ロールバック時は公開状態を維持

### 自動更新

- ファイル単位で自動更新を設定可能
- 自動更新時はメタ情報の変更も反映

## 次のステップ

1. データベーススキーマの設計（ADR 0005の更新）
2. API仕様の設計（新規ADR作成）
3. 実装の優先順位付け
4. テスト戦略の策定

## 関連ドキュメント

- [ドキュメント管理コンテキスト README](docs/document-management/README.md)
- [ADR 0001: External Repository Markdown Structure](adr/0001-external-repository-markdown-structure.md)
- [ADR 0013: Document Variable Definition and Substitution](adr/0013-document-variable-definition.md)
- [ADR 0014: Execution Record and Evidence Management](adr/0014-execution-record-and-evidence-management.md)
