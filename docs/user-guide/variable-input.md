# 変数入力機能ガイド

このガイドでは、OpsCoreの変数入力機能について説明します。

## 目次

- [変数入力機能とは](#変数入力機能とは)
- [変数の型](#変数の型)
- [変数の入力方法](#変数の入力方法)
- [バリデーション](#バリデーション)
- [変数の保存と再利用](#変数の保存と再利用)

## 変数入力機能とは

変数入力機能は、手順書を実行する際に、環境やケースに応じたパラメータを動的に入力できる機能です。

### メリット

- **再利用性**: 1つの手順書で複数の環境に対応
- **正確性**: 入力値の型チェックとバリデーション
- **作業証跡**: 使用した変数値を記録

### 使用例

```markdown
# サーバーデプロイ手順

## 前提条件
デプロイ先: {{server_name}}
デプロイ日時: {{deploy_date}}

## 手順
1. {{server_name}}にSSH接続
   ```bash
   ssh admin@{{server_name}}.example.com
   ```

2. アプリケーションのデプロイ
   バージョン: {{app_version}}
```

上記の手順書では、`server_name`、`deploy_date`、`app_version`が変数として定義されています。

## 変数の型

OpsCoreは、以下の4つの変数型をサポートしています：

### 1. String（文字列）

任意の文字列を入力できます。

```yaml
variables:
  - name: server_name
    label: サーバー名
    type: string
    required: true
    defaultValue: server01
```

**入力例**: `production-web-01`

### 2. Number（数値）

数値のみを入力できます。

```yaml
variables:
  - name: port_number
    label: ポート番号
    type: number
    required: true
    defaultValue: 8080
```

**入力例**: `8080`

### 3. Boolean（真偽値）

true/falseの2値を選択できます。

```yaml
variables:
  - name: enable_ssl
    label: SSL有効化
    type: boolean
    required: true
    defaultValue: true
```

**入力例**: チェックボックスでON/OFF

### 4. Date（日付）

日付を入力できます。

```yaml
variables:
  - name: deploy_date
    label: デプロイ日時
    type: date
    required: true
    defaultValue: "2024-01-15"
```

**入力例**: カレンダーから選択または手入力 `2024-01-15`

## 変数の入力方法

### 基本的な入力手順

```
1. ドキュメントを開く
2. 左ペインの変数入力フォームを確認
3. 各変数に値を入力
4. 「適用」ボタンをクリック
5. 中央ペインで変数が置換された手順書を確認
```

### 入力フォームの見方

```
┌─────────────────────────────┐
│ 変数入力                     │
├─────────────────────────────┤
│ サーバー名 *                 │
│ [server01            ]      │
│ デプロイ先のサーバー名       │
├─────────────────────────────┤
│ ポート番号 *                 │
│ [8080                ]      │
│ アプリケーションポート       │
├─────────────────────────────┤
│ SSL有効化                    │
│ [✓] 有効にする               │
├─────────────────────────────┤
│ [適用] [リセット]            │
└─────────────────────────────┘

* = 必須項目
```

### 型別の入力方法

#### String（文字列）

```
- テキストボックスに直接入力
- 特殊文字も入力可能
- 最大長: 通常1000文字まで
```

#### Number（数値）

```
- 数値のみ入力可能
- 小数点も使用可能
- 範囲指定がある場合、その範囲内で入力
```

#### Boolean（真偽値）

```
- チェックボックスでON/OFF
- または、ドロップダウンでtrue/false選択
```

#### Date（日付）

```
- カレンダーウィジェットから選択
- または、手入力（YYYY-MM-DD形式）
- 時刻が必要な場合は時刻選択も表示
```

## バリデーション

### 必須項目のチェック

必須項目（`required: true`）が未入力の場合、エラーメッセージが表示されます。

```
エラー: サーバー名は必須です
```

「適用」ボタンは、全ての必須項目が入力されるまで無効化されます。

### 型のチェック

入力値が指定された型と一致しない場合、エラーメッセージが表示されます。

```
エラー: ポート番号は数値で入力してください
```

### カスタムバリデーション

変数定義に`validation`が指定されている場合、そのルールに従ってチェックされます。

```yaml
variables:
  - name: email
    label: メールアドレス
    type: string
    required: true
    validation:
      pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
      message: "有効なメールアドレスを入力してください"
```

## 変数の保存と再利用

### 入力値の保存

変数入力値は、作業証跡と共に自動的に保存されます。

```
1. 変数を入力
2. 「適用」ボタンをクリック
3. 作業を開始（右ペインで作業証跡を作成）
4. 作業証跡に入力値が記録される
```

### 過去の入力値の再利用

```
1. ドキュメントを開く
2. 左ペインで「過去の入力値」をクリック
3. 作業証跡の一覧から選択
4. 「この値を使用」ボタンをクリック
5. 変数フォームに自動入力される
```

### テンプレートの作成（将来実装予定）

よく使用する変数の組み合わせをテンプレートとして保存できます。

```
1. 変数を入力
2. 「テンプレートとして保存」ボタンをクリック
3. テンプレート名を入力
4. 「保存」ボタンをクリック
```

## 実際の使用例

### 例1: サーバーデプロイ

**変数定義**:
```yaml
variables:
  - name: environment
    label: 環境
    type: string
    required: true
    defaultValue: staging
  - name: server_name
    label: サーバー名
    type: string
    required: true
  - name: app_version
    label: アプリバージョン
    type: string
    required: true
  - name: enable_backup
    label: バックアップを取る
    type: boolean
    defaultValue: true
```

**入力例**:
- 環境: `production`
- サーバー名: `web-prod-01`
- アプリバージョン: `v2.5.0`
- バックアップを取る: `ON`

**置換後の手順書**:
```markdown
# デプロイ手順（production環境）

## 対象サーバー
- サーバー名: web-prod-01
- バージョン: v2.5.0

## 手順
1. バックアップの取得（有効）
2. web-prod-01へのデプロイ...
```

### 例2: データベースメンテナンス

**変数定義**:
```yaml
variables:
  - name: db_name
    label: データベース名
    type: string
    required: true
  - name: maintenance_date
    label: メンテナンス日時
    type: date
    required: true
  - name: vacuum_full
    label: VACUUM FULLを実行
    type: boolean
    defaultValue: false
```

**入力例**:
- データベース名: `opscore_production`
- メンテナンス日時: `2024-02-01`
- VACUUM FULLを実行: `OFF`

## ベストプラクティス

### 1. デフォルト値の設定

よく使用する値はデフォルト値として設定すると便利です。

```yaml
variables:
  - name: port
    defaultValue: 8080  # 推奨値を設定
```

### 2. 分かりやすいラベルと説明

```yaml
variables:
  - name: replica_count
    label: レプリカ数
    description: 起動するコンテナの数。負荷に応じて調整してください。
```

### 3. 適切な型の選択

- 数値は`number`型を使用（文字列混入を防止）
- ON/OFFの選択は`boolean`型を使用
- 日付は`date`型を使用（形式の統一）

## トラブルシューティング

### Q: 変数が置換されない

**A**: 以下を確認してください：
- 変数名が正しいか（大文字小文字も一致させる）
- `{{variable_name}}`の形式で記載されているか
- 「適用」ボタンをクリックしたか

### Q: バリデーションエラーが消えない

**A**: 以下を確認してください：
- 入力値が型と一致しているか
- 必須項目が全て入力されているか
- カスタムバリデーションの条件を満たしているか

### Q: デフォルト値が表示されない

**A**: 以下を確認してください：
- Frontmatterに`defaultValue`が設定されているか
- ブラウザのキャッシュをクリアしてみる

## 関連ドキュメント

- [ユーザーガイド](./README.md)
- [ドキュメント管理](./document-management.md)
- [作業証跡記録](./execution-record.md)
- [ADR 0013: Document Variable Definition](../../adr/0013-document-variable-definition.md)
