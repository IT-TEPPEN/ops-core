# OpsCore Frontend

OpsCoreフロントエンドは、React + TypeScript + Viteで実装された、運用手順書管理システムのWebインターフェースです。

## 技術スタック

- **フレームワーク**: React 18
- **言語**: TypeScript
- **ビルドツール**: Vite
- **ルーティング**: React Router
- **テスト**: Vitest + React Testing Library
- **リンター**: ESLint

## アプリケーション構成

### ページ構成

本アプリケーションは以下のページで構成されています：

#### 1. リポジトリ管理
- **リポジトリ一覧ページ** (`/repositories`)
  - 登録済みGitリポジトリの一覧表示
  - リポジトリの追加・編集・削除

#### 2. ドキュメント管理
- **ドキュメント一覧ページ** (`/documents`)
  - 公開ドキュメントの一覧表示
  - ドキュメントの公開・非公開設定
  - バージョン管理とロールバック
- **ドキュメント閲覧ページ** (`/documents/:id`)
  - Markdownコンテンツの表示
  - 変数入力フォーム（左ペイン）
  - 作業証跡記録パネル（右ペイン）

#### 3. 作業証跡管理
- **作業証跡一覧ページ** (`/execution-records`)
  - 過去の作業証跡の検索・フィルタリング
  - ステータス別の表示（進行中/完了/失敗）
- **作業証跡詳細ページ** (`/execution-records/:id`)
  - 作業証跡の詳細表示
  - 各ステップのメモと画面キャプチャ
  - 使用した変数値の確認

#### 4. ユーザー・グループ管理（管理者のみ）
- **ユーザー管理ページ** (`/users`)
  - ユーザーの一覧・追加・編集・削除
- **グループ管理ページ** (`/groups`)
  - グループの一覧・追加・編集・削除
  - メンバーの追加・削除

### 主要コンポーネント

#### 変数入力フォーム (`VariableForm`)
手順書に定義された変数を入力するためのフォームコンポーネント。

**機能:**
- 型別の入力フィールド（string/number/boolean/date）
- 必須項目のバリデーション
- デフォルト値の自動入力
- リアルタイムプレビュー

#### 作業証跡パネル (`ExecutionRecordPanel`)
作業証跡を記録するための右サイドパネルコンポーネント。

**機能:**
- 作業ステップの追加
- ステップごとのメモ入力
- 画面キャプチャのアップロード
- 作業ステータスの変更（進行中/完了/失敗）
- 作業証跡の共有設定

#### Markdownビューアー (`MarkdownViewer`)
Markdownコンテンツを表示するコンポーネント。

**機能:**
- Markdownのレンダリング
- シンタックスハイライト
- 変数置換後のコンテンツ表示

## 開発

### セットアップ

```bash
# 依存関係のインストール
npm install

# 開発サーバーの起動
npm run dev

# ブラウザで http://localhost:5173 を開く
```

### ビルド

```bash
# プロダクションビルド
npm run build

# ビルド結果のプレビュー
npm run preview
```

## Testing

This project uses Vitest and React Testing Library for testing.

### Running Tests

```bash
# Run all tests once
npm test

# Run tests in watch mode
npm run test:watch

# Run tests with coverage
npm run test:coverage
```

### Test Structure

- Unit tests for utilities are located in `src/utils/*.test.ts`
- Hook tests are located in `src/hooks/*.test.ts`
- Component tests are colocated with components (e.g., `src/App.test.tsx`, `src/pages/*.test.tsx`)

### Writing Tests

Tests use the following libraries:
- **Vitest**: Test runner and assertion library
- **React Testing Library**: For testing React components
- **@testing-library/jest-dom**: For additional DOM matchers

Example test:
```typescript
import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import MyComponent from "./MyComponent";

describe("MyComponent", () => {
  it("renders correctly", () => {
    render(
      <MemoryRouter>
        <MyComponent />
      </MemoryRouter>
    );
    expect(screen.getByText("Hello")).toBeInTheDocument();
  });
});
```

## Expanding the ESLint configuration

If you are developing a production application, we recommend updating the configuration to enable type-aware lint rules:

```js
export default tseslint.config({
  extends: [
    // Remove ...tseslint.configs.recommended and replace with this
    ...tseslint.configs.recommendedTypeChecked,
    // Alternatively, use this for stricter rules
    ...tseslint.configs.strictTypeChecked,
    // Optionally, add this for stylistic rules
    ...tseslint.configs.stylisticTypeChecked,
  ],
  languageOptions: {
    // other options...
    parserOptions: {
      project: ['./tsconfig.node.json', './tsconfig.app.json'],
      tsconfigRootDir: import.meta.dirname,
    },
  },
})
```

You can also install [eslint-plugin-react-x](https://github.com/Rel1cx/eslint-react/tree/main/packages/plugins/eslint-plugin-react-x) and [eslint-plugin-react-dom](https://github.com/Rel1cx/eslint-react/tree/main/packages/plugins/eslint-plugin-react-dom) for React-specific lint rules:

```js
// eslint.config.js
import reactX from 'eslint-plugin-react-x'
import reactDom from 'eslint-plugin-react-dom'

export default tseslint.config({
  plugins: {
    // Add the react-x and react-dom plugins
    'react-x': reactX,
    'react-dom': reactDom,
  },
  rules: {
    // other rules...
    // Enable its recommended typescript rules
    ...reactX.configs['recommended-typescript'].rules,
    ...reactDom.configs.recommended.rules,
  },
})
```
