# OpsCoreへのコントリビューション

OpsCoreプロジェクトへのコントリビューションに興味をお持ちいただき、ありがとうございます！

## 📚 詳細なガイドライン

包括的な開発ガイドラインは、以下のドキュメントを参照してください：

**[開発ガイドライン（完全版）](docs/development/CONTRIBUTING.md)**

上記のドキュメントには、以下の詳細情報が記載されています：
- 開発環境のセットアップ
- 開発フロー
- コーディング規約（Go、TypeScript）
- テストの書き方と実行方法
- コミットメッセージ規約
- プルリクエストの作成方法
- ADRの書き方

## 🚀 クイックスタート

### 必要な環境

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- Docker & Docker Compose

### 開発環境の起動

```bash
# リポジトリのクローン
git clone https://github.com/IT-TEPPEN/ops-core.git
cd ops-core

# データベースの起動（Docker Compose）
docker compose up -d

# 環境変数の設定
export ENCRYPTION_KEY="dev-key-12345678901234567890123"  # 32 bytes

# バックエンドの起動
cd backend
go mod download
go run cmd/server/main.go

# フロントエンドの起動（別ターミナル）
cd frontend
npm install
npm run dev
```

アプリケーションへのアクセス：
- フロントエンド: http://localhost:5173
- バックエンドAPI: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html

## 📝 開発フロー（概要）

1. **Issueの確認**: 作業前に関連Issueを確認・作成
2. **ブランチ作成**: `git checkout -b feature/issue-number-description`
3. **開発**: コードを実装し、テストを追加
4. **テスト**: `go test ./...` または `npm test` で検証
5. **コミット**: 下記のコミットメッセージ規約に従ってコミット
6. **プッシュ**: `git push origin feature/issue-number-description`
7. **PR作成**: developブランチへのプルリクエストを作成
8. **レビュー**: レビュアーからのフィードバックに対応

## 🏗️ アーキテクチャ

本プロジェクトは**Onion Architecture**を採用しています。

- **バックエンド**: Go + Echo + PostgreSQL
- **フロントエンド**: React + TypeScript + Vite
- **詳細**: [ADRディレクトリ](adr/)および[バックエンドREADME](backend/README.md)を参照

## ✅ コミットメッセージ規約

```
<type>(<scope>): <subject>
```

**Type:**
- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメント
- `refactor`: リファクタリング
- `test`: テスト追加・修正
- `chore`: ビルド・ツール変更

**例:**
```
feat(document): 変数入力機能の追加

ドキュメントに定義された変数を入力できるフォームを実装。

Closes #38
```

## 🧪 テスト

```bash
# バックエンド
cd backend
go test ./...
go test -cover ./...

# フロントエンド
cd frontend
npm test
npm run test:coverage
```

詳細は[テストガイド](docs/development/TESTING.md)を参照してください。

## 🔒 セキュリティ

セキュリティ上の問題を発見した場合は、公開Issueではなく、メンテナーに直接連絡してください。

## 📖 その他のドキュメント

開発に役立つその他のドキュメント：

- **開発ガイド**
  - [テストガイド](docs/development/TESTING.md)
  - [API開発ガイド](docs/development/API.md)
- **アーキテクチャ**
  - [システム概要](docs/architecture/system-overview.md)
  - [バックエンドアーキテクチャ](backend/README.md)
  - [ADR一覧](adr/)
- **ユーザーガイド**
  - [ユーザーマニュアル](docs/user-guide/README.md)

## 💬 質問・サポート

- まず[ドキュメント](docs/)を確認
- 既存の[Issues](https://github.com/IT-TEPPEN/ops-core/issues)を検索
- 解決しない場合は新規Issueを作成

## 📄 ライセンス

このプロジェクトに貢献することで、あなたの貢献がプロジェクトのライセンスに従うことに同意したものとみなされます。

---

**詳細情報**: [開発ガイドライン（完全版）](docs/development/CONTRIBUTING.md)を必ずご確認ください。
