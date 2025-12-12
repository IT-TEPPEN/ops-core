# Deployment Guide

このドキュメントでは、OpsCoreのデプロイ手順を説明します。

## 目次

- [デプロイ方式](#デプロイ方式)
- [環境構成](#環境構成)
- [前提条件](#前提条件)
- [Docker Composeデプロイ](#docker-composeデプロイ)
- [Kubernetesデプロイ](#kubernetesデプロイ)
- [環境変数の設定](#環境変数の設定)
- [データベースのセットアップ](#データベースのセットアップ)
- [ストレージの設定](#ストレージの設定)
- [デプロイ後の確認](#デプロイ後の確認)

## デプロイ方式

OpsCoreは、以下のデプロイ方式をサポートしています：

1. **Docker Compose**: 開発環境・小規模環境向け
2. **Kubernetes**: 本番環境・大規模環境向け（将来対応予定）

## 環境構成

### 推奨環境

#### 開発環境
- CPU: 2コア以上
- メモリ: 4GB以上
- ストレージ: 20GB以上

#### 本番環境
- CPU: 4コア以上
- メモリ: 8GB以上
- ストレージ: 50GB以上（添付ファイル保存先による）

### コンポーネント構成

```
┌─────────────────┐
│  Load Balancer  │ (オプション)
└────────┬────────┘
         │
┌────────▼────────┐
│   Frontend      │ (React)
│   (Nginx/Static)│
└────────┬────────┘
         │
┌────────▼────────┐
│   Backend       │ (Go)
└────┬─────┬──────┘
     │     │
     │     └──────┐
     │            │
┌────▼────┐  ┌───▼─────┐
│PostgreSQL│  │ Storage │ (S3/MinIO/Local)
└─────────┘  └─────────┘
```

## 前提条件

### 必要なソフトウェア

- Docker: 20.10以上
- Docker Compose: 2.0以上
- Git

### ポートの確認

以下のポートが使用可能であることを確認してください：

- `8080`: バックエンドAPI
- `5173`: フロントエンド（開発時）
- `5432`: PostgreSQL
- `9000`: MinIO（使用する場合）

## Docker Composeデプロイ

### 1. リポジトリのクローン

```bash
git clone https://github.com/IT-TEPPEN/ops-core.git
cd ops-core
```

### 2. 環境変数ファイルの作成

`.env`ファイルを作成します：

```bash
cp .env.example .env
```

`.env`ファイルを編集：

```bash
# データベース設定
DB_HOST=postgres
DB_PORT=5432
DB_USER=opscore
DB_PASSWORD=your_secure_password_here
DB_NAME=opscore
DB_SSLMODE=disable

# 暗号化キー（32バイト）
ENCRYPTION_KEY=your_32_byte_encryption_key_here

# ストレージ設定
STORAGE_TYPE=local  # local, s3, minio
LOCAL_STORAGE_PATH=/app/storage

# バックエンドAPI URL
API_BASE_URL=http://localhost:8080

# フロントエンド URL
FRONTEND_URL=http://localhost:5173
```

### 3. Docker Composeの起動

```bash
# ビルドと起動
docker compose up -d

# ログの確認
docker compose logs -f

# 起動確認
docker compose ps
```

### 4. データベースマイグレーション

```bash
# マイグレーションの実行
docker compose exec backend migrate -path /app/migrations -database "postgres://opscore:password@postgres:5432/opscore?sslmode=disable" up
```

### 5. 動作確認

- フロントエンド: http://localhost:5173
- バックエンドAPI: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html

### 6. 停止とクリーンアップ

```bash
# 停止
docker compose down

# データも含めて削除
docker compose down -v
```

## Kubernetesデプロイ

（将来対応予定）

Kubernetesへのデプロイは、以下のマニフェストを使用します：

- `k8s/deployment.yaml`: アプリケーションのDeployment
- `k8s/service.yaml`: Service定義
- `k8s/ingress.yaml`: Ingress設定
- `k8s/configmap.yaml`: 環境変数の設定
- `k8s/secret.yaml`: シークレット情報

## 環境変数の設定

### 必須環境変数

#### バックエンド

```bash
# データベース接続
DB_HOST=localhost
DB_PORT=5432
DB_USER=opscore
DB_PASSWORD=secure_password
DB_NAME=opscore
DB_SSLMODE=require  # 本番環境では require

# 暗号化キー（32バイト）
ENCRYPTION_KEY=<32バイトのランダム文字列>

# ストレージ設定
STORAGE_TYPE=s3  # local, s3, minio
```

#### フロントエンド

```bash
# APIエンドポイント
VITE_API_BASE_URL=https://api.example.com
```

### オプション環境変数

#### S3ストレージ使用時

```bash
STORAGE_TYPE=s3
S3_REGION=ap-northeast-1
S3_BUCKET=opscore-attachments
AWS_ACCESS_KEY_ID=<アクセスキー>
AWS_SECRET_ACCESS_KEY=<シークレットキー>
```

#### MinIOストレージ使用時

```bash
STORAGE_TYPE=minio
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=opscore-attachments
MINIO_USE_SSL=false
```

### 暗号化キーの生成

```bash
# Linuxの場合
openssl rand -base64 32

# または
head -c 32 /dev/urandom | base64
```

## データベースのセットアップ

### PostgreSQLのインストール

#### Docker Composeを使用する場合

`docker-compose.yml`に含まれているため、個別のインストールは不要です。

#### 独立したPostgreSQLを使用する場合

```bash
# PostgreSQL 14のインストール（Ubuntu）
sudo apt update
sudo apt install postgresql-14

# データベースとユーザーの作成
sudo -u postgres psql
CREATE DATABASE opscore;
CREATE USER opscore WITH ENCRYPTED PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE opscore TO opscore;
```

### マイグレーションの実行

```bash
# migrate CLIのインストール
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# マイグレーションの実行
migrate -path backend/migrations -database "postgres://opscore:password@localhost:5432/opscore?sslmode=disable" up

# マイグレーションバージョンの確認
migrate -path backend/migrations -database "postgres://opscore:password@localhost:5432/opscore?sslmode=disable" version
```

### バックアップとリストア

詳細は [../operations/BACKUP.md](../operations/BACKUP.md) を参照してください。

## ストレージの設定

### ローカルファイルシステム

開発環境向けです。本番環境では推奨されません。

```bash
STORAGE_TYPE=local
LOCAL_STORAGE_PATH=/app/storage
```

ディレクトリの作成：

```bash
mkdir -p /app/storage
chmod 755 /app/storage
```

### AWS S3

本番環境向けの推奨設定です。

```bash
STORAGE_TYPE=s3
S3_REGION=ap-northeast-1
S3_BUCKET=opscore-attachments
AWS_ACCESS_KEY_ID=<アクセスキー>
AWS_SECRET_ACCESS_KEY=<シークレットキー>
```

S3バケットの作成：

```bash
aws s3 mb s3://opscore-attachments --region ap-northeast-1

# バケットポリシーの設定（プライベートアクセス）
aws s3api put-bucket-policy --bucket opscore-attachments --policy file://bucket-policy.json
```

### MinIO

オンプレミス環境向けのS3互換ストレージです。

```bash
STORAGE_TYPE=minio
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=opscore-attachments
MINIO_USE_SSL=false
```

MinIOの起動：

```bash
docker run -d \
  --name minio \
  -p 9000:9000 \
  -p 9001:9001 \
  -e MINIO_ROOT_USER=minioadmin \
  -e MINIO_ROOT_PASSWORD=minioadmin \
  -v /data/minio:/data \
  minio/minio server /data --console-address ":9001"
```

## デプロイ後の確認

### ヘルスチェック

```bash
# バックエンドのヘルスチェック
curl http://localhost:8080/health

# 期待される応答
{
  "status": "ok",
  "database": "connected",
  "storage": "available"
}
```

### ログの確認

```bash
# Docker Composeの場合
docker compose logs -f backend
docker compose logs -f frontend

# 特定のコンテナのログ
docker logs -f <container_id>
```

### データベース接続の確認

```bash
# PostgreSQLに接続
docker compose exec postgres psql -U opscore -d opscore

# テーブル一覧の確認
\dt

# 接続の確認
SELECT 1;
```

### ストレージの確認

```bash
# ローカルストレージの場合
ls -la /app/storage

# S3の場合
aws s3 ls s3://opscore-attachments

# MinIOの場合
# http://localhost:9001 にアクセスしてコンソールで確認
```

## トラブルシューティング

### データベース接続エラー

```bash
# PostgreSQLの起動確認
docker compose ps postgres

# PostgreSQLのログ確認
docker compose logs postgres

# 接続テスト
psql -h localhost -U opscore -d opscore
```

### ストレージアクセスエラー

```bash
# S3バケットの確認
aws s3 ls s3://opscore-attachments

# IAMポリシーの確認
aws iam get-user-policy --user-name opscore-user --policy-name S3Access
```

### メモリ不足

```bash
# コンテナのリソース使用状況確認
docker stats

# メモリ制限の調整（docker-compose.yml）
services:
  backend:
    mem_limit: 2g
```

## セキュリティ考慮事項

### 1. 暗号化キーの管理

- 環境変数で設定し、コードにハードコードしない
- 定期的にローテーション
- 32バイト以上のランダム文字列を使用

### 2. データベース認証

- 強固なパスワードを使用
- 本番環境ではSSL接続を有効化（`sslmode=require`）
- 最小権限の原則に従う

### 3. APIアクセス

- HTTPS通信を使用（本番環境）
- CORS設定を適切に構成
- レート制限の実装（将来）

### 4. ファイアウォール設定

```bash
# 必要なポートのみ開放
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw allow 22/tcp    # SSH（管理用）
```

## スケーリング

### 垂直スケーリング（スケールアップ）

リソースの増強：

```yaml
# docker-compose.yml
services:
  backend:
    cpus: '2.0'
    mem_limit: 4g
```

### 水平スケーリング（スケールアウト）

複数インスタンスの起動：

```bash
docker compose up -d --scale backend=3
```

ロードバランサーの設定が必要です（Nginx、HAProxyなど）。

## 更新とロールバック

### アプリケーションの更新

```bash
# 最新版の取得
git pull origin main

# イメージの再ビルド
docker compose build

# 再起動
docker compose up -d
```

### ロールバック

```bash
# 特定バージョンへのロールバック
git checkout <previous_version_tag>
docker compose build
docker compose up -d
```

## 監視とアラート

詳細は [../operations/MONITORING.md](../operations/MONITORING.md) を参照してください。

## 関連ドキュメント

- [運用ガイド](../operations/MONITORING.md)
- [バックアップ手順](../operations/BACKUP.md)
- [セキュリティガイド](../../adr/0015-backend-custom-error-design.md)
