
# ドメインモデル設計見直し用チェックシート（2025-11-05更新）

## 概要

このチェックシートは、Issue #29（ドメインモデル図の設計見直し）に基づくドメインモデル再設計の実装前確認用です。
再設計内容については以下を参照してください：

- [ドメインモデル再設計サマリー](/workspaces/domain-model-redesign-summary.md)
- [ドキュメント管理コンテキスト README](/workspaces/docs/document-management/README.md)
- [ADR 0013: Document Variable Definition](/workspaces/adr/0013-document-variable-definition.md)
- [ADR 0014: Execution Record and Evidence Management](/workspaces/adr/0014-execution-record-and-evidence-management.md)

## チェック項目一覧

| #    | カテゴリ         | チェック項目                              | 対象ファイル                                                                                  | 内容                                                 | 修正要否 | 関連Issue | 備考                                                     |
| :--- | :--------------- | :---------------------------------------- | :-------------------------------------------------------------------------------------------- | :--------------------------------------------------- | :------- | :-------- | :------------------------------------------------------- |
| 1    | ADR              | ADR設計方針との整合性                     | adr/0007-backend-architecture-onion.md                                                        | Onionアーキテクチャ・依存関係の確認                  | 要確認   | #14       | internal/配下への移動が必要                              |
| 2    | ADR              | DBスキーマとの整合性                      | adr/0005-database-schema.md<br>adr/0004-database-specification.md                             | 新規エンティティのスキーマ定義追加                   | **必要** | #13, #36  | Document、DocumentVersion、ExecutionRecord等の追加が必要 |
| 3    | ADR              | 外部仕様との整合性                        | adr/0001-external-repository-markdown-structure.md<br>adr/0002-repository-access-method.md    | variablesフィールドの追加済み確認                    | **完了** | -         | ADR 0001は更新済み                                       |
| 4    | ADR              | managed_filesテーブルのADR作成            | （新規ADR作成）                                                                               | managed_filesの設計判断を文書化                      | **必要** | #19, #37  | 既存実装の後追いドキュメント化                           |
| 5    | ADR              | API仕様の整合性                           | adr/0003-backend-api-markdown-fetch.md<br>adr/0010-api-definition-generation-specification.md | 新規API仕様の追加・更新                              | **必要** | #16, 新規 | 変数入力、作業証跡、検索等のAPI追加                      |
| 6    | ドメインモデル   | 既存Entity修正                            | backend/domain/model/repository.go<br>backend/domain/model/file_node.go                       | 既存エンティティの影響確認                           | 要確認   | -         | FileNodeはDocumentに統合の可能性                         |
| 7    | ドメインモデル   | Document Entity新規作成                   | backend/domain/model/document.go（新規）                                                      | Document集約ルートの実装                             | **必要** | #31       | variables, versions等を含む                              |
| 8    | ドメインモデル   | DocumentVersion Entity新規作成            | backend/domain/model/document_version.go（新規）                                              | DocumentVersionエンティティの実装                    | **必要** | #31       | バージョン管理用                                         |
| 9    | ドメインモデル   | ExecutionRecord集約新規作成               | backend/domain/model/execution_record.go（新規）                                              | ExecutionRecord集約ルートの実装                      | **必要** | #32       | 作業証跡機能の中核                                       |
| 10   | ドメインモデル   | ExecutionStep Entity新規作成              | backend/domain/model/execution_step.go（新規）                                                | ExecutionStepエンティティの実装                      | **必要** | #32       | 作業ステップ記録                                         |
| 11   | ドメインモデル   | Attachment Entity新規作成                 | backend/domain/model/attachment.go（新規）                                                    | Attachmentエンティティの実装                         | **必要** | #32       | 画面キャプチャ等の添付ファイル                           |
| 12   | ドメインモデル   | User集約新規作成                          | backend/domain/model/user.go（新規）                                                          | User集約ルートの実装                                 | **必要** | #24, #33  | 認証機能と連携                                           |
| 13   | ドメインモデル   | Group集約新規作成                         | backend/domain/model/group.go（新規）                                                         | Group集約ルートの実装                                | **必要** | #33       | グループ管理機能                                         |
| 14   | ドメインモデル   | ViewHistory Entity新規作成                | backend/domain/model/view_history.go（新規）                                                  | ViewHistoryエンティティの実装                        | **必要** | #30       | 閲覧履歴記録                                             |
| 15   | ドメインモデル   | ViewStatistics集約新規作成                | backend/domain/model/view_statistics.go（新規）                                               | ViewStatistics集約ルートの実装                       | **必要** | #30       | 閲覧統計管理                                             |
| 16   | 値オブジェクト   | 新規値オブジェクト作成                    | backend/domain/model/value_objects.go（新規または分割）                                       | VariableDefinition, VariableValue, StorageType等     | **必要** | #31, #32  | 複数ファイルに分割も検討                                 |
| 17   | Domain Service   | Domain Serviceの設計・実装                | backend/domain/service/（新規追加検討）                                                       | 複数エンティティにまたがるロジック                   | 要検討   | -         | 必要に応じて追加                                         |
| 18   | テスト           | 新規Entityのテスト作成                    | backend/domain/model/*_test.go（新規）                                                        | 全新規エンティティのユニットテスト                   | **必要** | #17       | テストカバレッジ80%以上目標                              |
| 19   | Repository       | DocumentRepositoryインターフェース        | backend/domain/repository/document_repository.go（新規）                                      | Document集約の永続化インターフェース                 | **必要** | #31       | CRUD + バージョン管理                                    |
| 20   | Repository       | ExecutionRecordRepositoryインターフェース | backend/domain/repository/execution_record_repository.go（新規）                              | ExecutionRecord集約の永続化インターフェース          | **必要** | #32       | CRUD + 検索機能                                          |
| 21   | Repository       | UserRepositoryインターフェース            | backend/domain/repository/user_repository.go（新規）                                          | User集約の永続化インターフェース                     | **必要** | #24, #33  | 認証機能と連携                                           |
| 22   | Repository       | GroupRepositoryインターフェース           | backend/domain/repository/group_repository.go（新規）                                         | Group集約の永続化インターフェース                    | **必要** | #33       | メンバー管理含む                                         |
| 23   | Repository       | ViewHistoryRepositoryインターフェース     | backend/domain/repository/view_history_repository.go（新規）                                  | ViewHistory永続化インターフェース                    | **必要** | #30       | 閲覧履歴記録                                             |
| 24   | Repository       | ViewStatisticsRepositoryインターフェース  | backend/domain/repository/view_statistics_repository.go（新規）                               | ViewStatistics永続化インターフェース                 | **必要** | #30       | 統計情報管理                                             |
| 25   | Repository       | AttachmentRepositoryインターフェース      | backend/domain/repository/attachment_repository.go（新規）                                    | Attachment永続化インターフェース（ストレージ抽象化） | **必要** | #32, #34  | local/S3/MinIO対応                                       |
| 26   | Infrastructure   | ストレージ抽象化層の実装                  | backend/infrastructure/storage/（新規）                                                       | ローカル/S3/MinIO対応のストレージ層                  | **必要** | #34       | AttachmentRepository実装で使用                           |
| 27   | Infrastructure   | Repository実装（永続化層）                | backend/infrastructure/persistence/*_repository_impl.go（新規）                               | 各Repositoryインターフェースの実装                   | **必要** | 新規      | PostgreSQL実装                                           |
| 28   | Infrastructure   | アクセストークン暗号化                    | backend/infrastructure/persistence/repository_repository_impl.go                              | アクセストークンの暗号化・復号化                     | **必要** | #12       | セキュリティ対応                                         |
| 29   | マイグレーション | 新規テーブルのマイグレーション            | backend/infrastructure/persistence/migrations/（新規）                                        | documents, document_versions, execution_records等    | **必要** | #35       | 10個以上の新規テーブル                                   |
| 30   | Application      | 変数入力Usecase                           | backend/application/usecase/（新規）                                                          | 変数定義取得、バリデーション                         | **必要** | #38       | ADR 0013に基づく                                         |
| 31   | Application      | 作業証跡Usecase                           | backend/application/usecase/execution_record_usecase.go（新規）                               | 作業証跡CRUD、検索、共有                             | **必要** | #39       | ADR 0014に基づく                                         |
| 32   | Application      | ドキュメント管理Usecase                   | backend/application/usecase/document_usecase.go（新規）                                       | 公開、非公開、バージョン管理、ロールバック           | **必要** | #40       | 既存リポジトリUsecaseとの統合検討                        |
| 33   | Application      | ユーザー・グループ管理Usecase             | backend/application/usecase/user_usecase.go等（新規）                                         | ユーザー・グループCRUD                               | **必要** | #24, #44  | 認証機能と連携                                           |
| 34   | Application      | 閲覧履歴・統計Usecase                     | backend/application/usecase/view_usecase.go（新規）                                           | 閲覧記録、統計更新、履歴取得                         | **必要** | #42       | -                                                        |
| 32   | Application      | ドキュメント管理Usecase                   | backend/application/usecase/document_usecase.go（新規）                                       | 公開、非公開、バージョン管理、ロールバック           | **必要** | 新規      | 既存リポジトリUsecaseとの統合検討                        |
| 33   | Application      | ユーザー・グループ管理Usecase             | backend/application/usecase/user_usecase.go等（新規）                                         | ユーザー・グループCRUD                               | **必要** | #24       | 認証機能と連携                                           |
| 34   | Application      | 閲覧履歴・統計Usecase                     | backend/application/usecase/view_usecase.go（新規）                                           | 閲覧記録、統計更新、履歴取得                         | **必要** | 新規      | -                                                        |
| 35   | Application      | DTOの追加・整理                           | backend/application/dto/（新規または追加）                                                    | 新規API用のDTO定義                                   | **必要** | #15       | application/dto/レイヤー作成                             |
| 36   | Interfaces       | 変数入力API Handler                       | backend/interfaces/api/handlers/variable_handler.go（新規）                                   | 変数定義取得API                                      | **必要** | #38       | -                                                        |
| 37   | Interfaces       | 作業証跡API Handler                       | backend/interfaces/api/handlers/execution_handler.go（新規）                                  | 作業証跡CRUD、検索、共有API                          | **必要** | #39       | -                                                        |
| 38   | Interfaces       | ドキュメント管理API Handler               | backend/interfaces/api/handlers/document_handler.go（新規または更新）                         | 公開、非公開、バージョン管理API                      | **必要** | #40       | 既存repository_handlerとの統合検討                       |
| 39   | Interfaces       | 添付ファイルアップロードAPI Handler       | backend/interfaces/api/handlers/attachment_handler.go（新規）                                 | 画像アップロード、取得、削除API                      | **必要** | #41       | multipart/form-data対応                                  |
| 40   | Interfaces       | ユーザー・グループAPI Handler             | backend/interfaces/api/handlers/user_handler.go等（新規）                                     | ユーザー・グループCRUD API                           | **必要** | #24, #44  | 認証ミドルウェアと連携                                   |
| 41   | Frontend         | 変数入力フォームコンポーネント            | frontend/src/components/VariableForm.tsx（新規）                                              | 変数入力UI                                           | **必要** | #38       | 型別の入力フィールド                                     |
| 42   | Frontend         | 作業証跡記録コンポーネント                | frontend/src/components/ExecutionRecordPanel.tsx（新規）                                      | 右ペインの作業証跡UI                                 | **必要** | #39       | ステップ追加、画像アップロード                           |
| 43   | Frontend         | 作業証跡一覧・検索ページ                  | frontend/src/pages/ExecutionRecordsPage.tsx（新規）                                           | 作業証跡の検索・フィルタリング                       | **必要** | #39       | -                                                        |
| 44   | Frontend         | 作業証跡詳細ページ                        | frontend/src/pages/ExecutionRecordDetailPage.tsx（新規）                                      | 過去の作業証跡閲覧                                   | **必要** | #39       | -                                                        |
| 45   | Frontend         | ドキュメント公開管理ページ                | frontend/src/pages/DocumentManagementPage.tsx（新規）                                         | 公開設定、バージョン管理UI                           | **必要** | #40       | -                                                        |
| 46   | Frontend         | フロントエンドテスト                      | frontend/src/**/*.test.tsx（新規）                                                            | 新規コンポーネントのテスト                           | **必要** | #18, #46  | Vitest + React Testing Library                           |
| 47   | ファイル構造     | internal/配下への移動                     | backend/内の全ファイル                                                                        | ADR 0007に準拠したフォルダ構造                       | **必要** | #14       | domain → internal/domain等                               |
| 48   | セキュリティ     | 認証・認可機能の実装                      | backend/interfaces/api/middleware/auth.go等（新規）                                           | JWT/Session認証、RBAC                                | **必要** | #24       | Phase 2の先行実装                                        |
| 49   | ドキュメント     | API仕様書の更新                           | docs/api/（新規または更新）                                                                   | Swagger/OpenAPI仕様の更新                            | **必要** | #45       | 新規APIの仕様書化                                        |
| 50   | ドキュメント     | README更新                                | README.md                                                                                     | Phase 2機能の追加、実装状況の更新                    | **必要** | #47       | -                                                        |

## 修正要否の凡例

- **必要**: 実装必須
- **完了**: 既に対応済み
- 要確認: 詳細調査が必要
- 要検討: 実装の是非を判断する必要がある

## 作成済みIssue一覧

以下のIssueを作成しました：

1. **#30** - 【親Issue】ドメインモデル再設計の実装
2. **#31** - Document集約の実装（Domain Model + Repository）
3. **#32** - ExecutionRecord集約の実装（Domain Model + Repository）
4. **#33** - User/Group集約の実装（Domain Model + Repository）
5. **#34** - ストレージ抽象化層の実装（local/S3/MinIO対応）
6. **#35** - 新規エンティティのDBマイグレーション実装
7. **#36** - DBスキーマADRの更新（ADR 0005）
8. **#37** - managed_filesテーブルのADR作成（ADR 0015）
9. **#38** - 変数入力機能の実装（Usecase + Handler + Frontend）
10. **#39** - 作業証跡機能の実装（ExecutionRecord管理）
11. **#40** - ドキュメント管理機能の実装（Document + Version管理）
12. **#41** - 添付ファイル管理機能の実装（画面キャプチャ・証跡添付）
13. **#42** - 閲覧履歴・統計機能の実装（ViewHistory + ViewStatistics）
14. **#43** - フロントエンド共通コンポーネントの実装
15. **#44** - グループ管理機能の実装（Ops-Core独自グループ）
16. **#45** - API仕様書の更新（Swagger）
17. **#46** - テストコードの実装（単体テスト・統合テスト）
18. **#47** - ドキュメント・README・ADRの更新

## 関連Issue

- #29: ドメインモデル図の設計見直し（本チェックシートの起点）
- #24: ユーザー認証・認可機能の設計と実装
- #19: managed_filesテーブルに関するADR作成
- #18: フロントエンドのテスト環境構築
- #17: GitManager実装のテスト追加
- #16: ADR 0003のAPI仕様を現在の実装に合わせて更新
- #15: application/dto/レイヤーの追加
- #14: backend配下のフォルダ構造をADR 0007に合わせて再構成
- #13: データベーススキーマのADR更新
- #12: アクセストークンの暗号化実装

## 実装の優先順位（推奨）

### Phase 1: 基盤整備

1. フォルダ構造の再構成（#14）
2. DTOレイヤーの追加（#15）
3. アクセストークン暗号化（#12）
4. ADR更新（#13, #19, API仕様等）

### Phase 2: ドメインモデル実装

1. 値オブジェクトの実装
2. Document集約の実装
3. DocumentVersion実装
4. 新規Repositoryインターフェース定義
5. DBマイグレーション作成
6. Repository実装（永続化層）

### Phase 3: コア機能実装

1. ドキュメント管理Usecase + Handler
2. 変数入力機能（Usecase + Handler + Frontend）
3. バージョン管理機能
4. 閲覧履歴・統計機能

### Phase 4: 作業証跡機能

1. ExecutionRecord集約の実装
2. ストレージ抽象化層の実装
3. 添付ファイル機能
4. 作業証跡Usecase + Handler
5. 作業証跡Frontend実装

### Phase 5: ユーザー管理

1. User/Group集約の実装
2. 認証・認可機能（#24）
3. ユーザー管理Usecase + Handler + Frontend

### Phase 6: テスト・ドキュメント

1. 全Entityのテスト追加（#17, #18）
2. API仕様書更新
3. README更新

---

**最終更新日**: 2025-11-05  
**更新者**: GitHub Copilot  
**関連Issue**: #29

