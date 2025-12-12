# Backup and Restore Guide

このドキュメントでは、OpsCoreのバックアップとリストア手順を説明します。

## 目次

- [バックアップ戦略](#バックアップ戦略)
- [データベースのバックアップ](#データベースのバックアップ)
- [添付ファイルのバックアップ](#添付ファイルのバックアップ)
- [設定ファイルのバックアップ](#設定ファイルのバックアップ)
- [リストア手順](#リストア手順)
- [災害復旧計画](#災害復旧計画)
- [自動バックアップの設定](#自動バックアップの設定)

## バックアップ戦略

### バックアップ対象

OpsCoreでは、以下のデータをバックアップする必要があります：

1. **PostgreSQLデータベース**: 全ての構造化データ
2. **添付ファイル**: 作業証跡の画面キャプチャ等
3. **設定ファイル**: 環境変数、設定ファイル
4. **暗号化キー**: アクセストークン暗号化に使用

### バックアップポリシー

#### 本番環境

| データ種別 | バックアップ頻度 | 保持期間 | 保存先 |
|-----------|---------------|----------|--------|
| データベース（フル） | 日次 | 30日 | S3/別サーバー |
| データベース（増分） | 1時間毎 | 7日 | S3/別サーバー |
| 添付ファイル | 日次 | 90日 | S3/別サーバー |
| 設定ファイル | 変更時 | 無期限 | Git/S3 |

#### 開発環境

| データ種別 | バックアップ頻度 | 保持期間 |
|-----------|---------------|----------|
| データベース | 週次 | 14日 |
| 添付ファイル | 不要 | - |

### バックアップの3-2-1ルール

- **3**: データのコピーを3つ保持
- **2**: 2つの異なるメディアに保存
- **1**: 1つはオフサイト（別拠点）に保存

## データベースのバックアップ

### pg_dumpによるバックアップ

#### フルバックアップ

```bash
#!/bin/bash
# backup-database.sh

# 設定
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="opscore"
DB_NAME="opscore"
BACKUP_DIR="/backup/postgres"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/opscore_${DATE}.sql.gz"

# バックアップディレクトリの作成
mkdir -p ${BACKUP_DIR}

# バックアップの実行
PGPASSWORD="${DB_PASSWORD}" pg_dump \
  -h ${DB_HOST} \
  -p ${DB_PORT} \
  -U ${DB_USER} \
  -d ${DB_NAME} \
  --format=custom \
  --compress=9 \
  -f ${BACKUP_FILE}

# 成功確認
if [ $? -eq 0 ]; then
    echo "Backup completed successfully: ${BACKUP_FILE}"
    
    # S3へのアップロード（オプション）
    aws s3 cp ${BACKUP_FILE} s3://opscore-backups/postgres/
    
    # 古いバックアップの削除（30日以上前）
    find ${BACKUP_DIR} -name "opscore_*.sql.gz" -mtime +30 -delete
else
    echo "Backup failed"
    exit 1
fi
```

#### Docker環境でのバックアップ

```bash
# Docker Composeを使用している場合
docker compose exec -T postgres pg_dump \
  -U opscore \
  -d opscore \
  --format=custom \
  --compress=9 \
  > backup_$(date +%Y%m%d).dump
```

### 継続的アーカイブとポイントインタイムリカバリ（PITR）

本番環境では、WAL（Write-Ahead Log）を使用した継続的アーカイブを推奨します。

#### postgresql.confの設定

```conf
# WALアーカイブの有効化
wal_level = replica
archive_mode = on
archive_command = 'aws s3 cp %p s3://opscore-backups/wal/%f'
archive_timeout = 300  # 5分

# WALの保持
max_wal_size = 4GB
min_wal_size = 1GB
```

#### ベースバックアップの作成

```bash
#!/bin/bash
# create-base-backup.sh

BACKUP_DIR="/backup/postgres/base"
DATE=$(date +%Y%m%d_%H%M%S)

# ベースバックアップの作成
pg_basebackup \
  -h localhost \
  -U opscore \
  -D ${BACKUP_DIR}/${DATE} \
  -Ft \
  -z \
  -P \
  -X fetch

# S3へのアップロード
tar -czf ${BACKUP_DIR}/${DATE}.tar.gz -C ${BACKUP_DIR} ${DATE}
aws s3 cp ${BACKUP_DIR}/${DATE}.tar.gz s3://opscore-backups/postgres/base/

# ローカルファイルの削除
rm -rf ${BACKUP_DIR}/${DATE}
```

## 添付ファイルのバックアップ

### ローカルストレージの場合

```bash
#!/bin/bash
# backup-attachments.sh

SOURCE_DIR="/app/storage"
BACKUP_DIR="/backup/attachments"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/attachments_${DATE}.tar.gz"

# バックアップの作成
tar -czf ${BACKUP_FILE} -C ${SOURCE_DIR} .

# S3へのアップロード
aws s3 cp ${BACKUP_FILE} s3://opscore-backups/attachments/

# 古いバックアップの削除（90日以上前）
find ${BACKUP_DIR} -name "attachments_*.tar.gz" -mtime +90 -delete
```

### S3ストレージの場合

S3を使用している場合、バックアップは以下の方法で実現できます：

#### 1. S3バージョニングの有効化

```bash
aws s3api put-bucket-versioning \
  --bucket opscore-attachments \
  --versioning-configuration Status=Enabled
```

#### 2. S3レプリケーションの設定

```json
{
  "Role": "arn:aws:iam::123456789012:role/s3-replication-role",
  "Rules": [
    {
      "Status": "Enabled",
      "Priority": 1,
      "Filter": {},
      "Destination": {
        "Bucket": "arn:aws:s3:::opscore-attachments-backup",
        "ReplicationTime": {
          "Status": "Enabled",
          "Time": {
            "Minutes": 15
          }
        }
      }
    }
  ]
}
```

#### 3. S3ライフサイクルポリシー

```json
{
  "Rules": [
    {
      "Id": "archive-old-files",
      "Status": "Enabled",
      "Transitions": [
        {
          "Days": 90,
          "StorageClass": "GLACIER"
        }
      ],
      "NoncurrentVersionExpiration": {
        "NoncurrentDays": 30
      }
    }
  ]
}
```

## 設定ファイルのバックアップ

### 環境変数のバックアップ

```bash
#!/bin/bash
# backup-config.sh

BACKUP_DIR="/backup/config"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p ${BACKUP_DIR}

# .envファイルのバックアップ（暗号化）
openssl enc -aes-256-cbc -salt \
  -in .env \
  -out ${BACKUP_DIR}/env_${DATE}.enc \
  -pass pass:${ENCRYPTION_PASSWORD}

# S3へのアップロード
aws s3 cp ${BACKUP_DIR}/env_${DATE}.enc s3://opscore-backups/config/
```

### Gitによる設定管理

設定ファイルはGitで管理し、シークレット情報は別途管理します：

```bash
# .gitignore
.env
*.key
secrets/
```

シークレット情報は、AWS Secrets Manager、HashiCorp Vault等で管理することを推奨します。

## リストア手順

### データベースのリストア

#### フルバックアップからのリストア

```bash
#!/bin/bash
# restore-database.sh

BACKUP_FILE="$1"  # 引数でバックアップファイルを指定

if [ -z "$BACKUP_FILE" ]; then
    echo "Usage: $0 <backup_file>"
    exit 1
fi

# S3からダウンロード（必要な場合）
if [[ $BACKUP_FILE == s3://* ]]; then
    LOCAL_FILE="/tmp/$(basename $BACKUP_FILE)"
    aws s3 cp $BACKUP_FILE $LOCAL_FILE
    BACKUP_FILE=$LOCAL_FILE
fi

# 既存データベースの削除（注意！）
PGPASSWORD="${DB_PASSWORD}" dropdb \
  -h ${DB_HOST} \
  -p ${DB_PORT} \
  -U ${DB_USER} \
  ${DB_NAME}

# データベースの再作成
PGPASSWORD="${DB_PASSWORD}" createdb \
  -h ${DB_HOST} \
  -p ${DB_PORT} \
  -U ${DB_USER} \
  ${DB_NAME}

# リストアの実行
PGPASSWORD="${DB_PASSWORD}" pg_restore \
  -h ${DB_HOST} \
  -p ${DB_PORT} \
  -U ${DB_USER} \
  -d ${DB_NAME} \
  -v \
  ${BACKUP_FILE}

if [ $? -eq 0 ]; then
    echo "Restore completed successfully"
else
    echo "Restore failed"
    exit 1
fi
```

#### ポイントインタイムリカバリ

```bash
#!/bin/bash
# pitr-restore.sh

TARGET_TIME="2024-01-15 10:30:00"  # リストア対象時刻
BASE_BACKUP="/backup/postgres/base/20240115_000000"
WAL_ARCHIVE="s3://opscore-backups/wal"

# PostgreSQLの停止
systemctl stop postgresql

# データディレクトリのクリア
rm -rf /var/lib/postgresql/14/main/*

# ベースバックアップのリストア
tar -xzf ${BASE_BACKUP}.tar.gz -C /var/lib/postgresql/14/main/

# recovery.confの作成
cat > /var/lib/postgresql/14/main/recovery.conf <<EOF
restore_command = 'aws s3 cp ${WAL_ARCHIVE}/%f %p'
recovery_target_time = '${TARGET_TIME}'
recovery_target_action = 'promote'
EOF

# PostgreSQLの起動
systemctl start postgresql

# リカバリの確認
tail -f /var/log/postgresql/postgresql-14-main.log
```

### 添付ファイルのリストア

```bash
#!/bin/bash
# restore-attachments.sh

BACKUP_FILE="$1"
RESTORE_DIR="/app/storage"

if [ -z "$BACKUP_FILE" ]; then
    echo "Usage: $0 <backup_file>"
    exit 1
fi

# S3からダウンロード（必要な場合）
if [[ $BACKUP_FILE == s3://* ]]; then
    LOCAL_FILE="/tmp/$(basename $BACKUP_FILE)"
    aws s3 cp $BACKUP_FILE $LOCAL_FILE
    BACKUP_FILE=$LOCAL_FILE
fi

# バックアップディレクトリの作成
mkdir -p ${RESTORE_DIR}

# リストアの実行
tar -xzf ${BACKUP_FILE} -C ${RESTORE_DIR}

if [ $? -eq 0 ]; then
    echo "Attachments restored successfully"
else
    echo "Restore failed"
    exit 1
fi
```

## 災害復旧計画（DRP）

### 目標復旧時間（RTO）と目標復旧時点（RPO）

| 環境 | RTO（復旧時間） | RPO（データ損失） |
|------|---------------|----------------|
| 本番環境 | 4時間以内 | 1時間以内 |
| ステージング | 24時間以内 | 24時間以内 |
| 開発環境 | 制限なし | 制限なし |

### 災害復旧手順

#### 1. 影響評価

- システム停止の原因特定
- データ損失の有無確認
- 影響範囲の確認

#### 2. 復旧作業

```bash
#!/bin/bash
# disaster-recovery.sh

echo "=== Disaster Recovery Process ==="

# 1. 最新のバックアップ確認
echo "Step 1: Checking latest backups..."
LATEST_DB_BACKUP=$(aws s3 ls s3://opscore-backups/postgres/ | sort | tail -n 1 | awk '{print $4}')
LATEST_ATTACHMENTS=$(aws s3 ls s3://opscore-backups/attachments/ | sort | tail -n 1 | awk '{print $4}')

echo "Latest DB backup: ${LATEST_DB_BACKUP}"
echo "Latest attachments backup: ${LATEST_ATTACHMENTS}"

# 2. インフラストラクチャの復旧
echo "Step 2: Restoring infrastructure..."
docker compose up -d postgres

# 3. データベースのリストア
echo "Step 3: Restoring database..."
./restore-database.sh s3://opscore-backups/postgres/${LATEST_DB_BACKUP}

# 4. 添付ファイルのリストア
echo "Step 4: Restoring attachments..."
./restore-attachments.sh s3://opscore-backups/attachments/${LATEST_ATTACHMENTS}

# 5. アプリケーションの起動
echo "Step 5: Starting application..."
docker compose up -d backend frontend

# 6. ヘルスチェック
echo "Step 6: Health check..."
sleep 10
curl -f http://localhost:8080/health || echo "Health check failed"

echo "=== Recovery process completed ==="
```

#### 3. 検証

- ヘルスチェックの確認
- 主要機能のテスト
- データ整合性の確認

#### 4. 通知

- ステークホルダーへの復旧完了通知
- インシデントレポートの作成

## 自動バックアップの設定

### Cronによる定期実行

```bash
# crontab -e

# データベースの日次フルバックアップ（毎日2:00）
0 2 * * * /opt/opscore/scripts/backup-database.sh >> /var/log/opscore/backup.log 2>&1

# 添付ファイルの日次バックアップ（毎日3:00）
0 3 * * * /opt/opscore/scripts/backup-attachments.sh >> /var/log/opscore/backup.log 2>&1

# データベースの増分バックアップ（毎時）
0 * * * * /opt/opscore/scripts/backup-database-incremental.sh >> /var/log/opscore/backup.log 2>&1

# バックアップの検証（毎日4:00）
0 4 * * * /opt/opscore/scripts/verify-backup.sh >> /var/log/opscore/backup-verify.log 2>&1
```

### バックアップの検証スクリプト

```bash
#!/bin/bash
# verify-backup.sh

echo "=== Backup Verification ==="

# 最新のバックアップファイルを取得
LATEST_BACKUP=$(aws s3 ls s3://opscore-backups/postgres/ | sort | tail -n 1 | awk '{print $4}')

# テスト用データベースへのリストア
TEST_DB="opscore_test_restore"

# データベースの作成
PGPASSWORD="${DB_PASSWORD}" createdb \
  -h ${DB_HOST} \
  -U ${DB_USER} \
  ${TEST_DB}

# リストアの実行
aws s3 cp s3://opscore-backups/postgres/${LATEST_BACKUP} /tmp/test_backup.dump
PGPASSWORD="${DB_PASSWORD}" pg_restore \
  -h ${DB_HOST} \
  -U ${DB_USER} \
  -d ${TEST_DB} \
  /tmp/test_backup.dump

# 検証
if [ $? -eq 0 ]; then
    echo "Backup verification successful: ${LATEST_BACKUP}"
    
    # 基本的なクエリテスト
    PGPASSWORD="${DB_PASSWORD}" psql \
      -h ${DB_HOST} \
      -U ${DB_USER} \
      -d ${TEST_DB} \
      -c "SELECT COUNT(*) FROM documents;" > /dev/null
    
    if [ $? -eq 0 ]; then
        echo "Database integrity check passed"
    else
        echo "Database integrity check failed"
        exit 1
    fi
else
    echo "Backup verification failed: ${LATEST_BACKUP}"
    exit 1
fi

# テストデータベースの削除
PGPASSWORD="${DB_PASSWORD}" dropdb \
  -h ${DB_HOST} \
  -U ${DB_USER} \
  ${TEST_DB}

rm /tmp/test_backup.dump
```

## バックアップの監視

### バックアップ失敗時のアラート

```bash
#!/bin/bash
# backup-with-alert.sh

# バックアップの実行
./backup-database.sh

if [ $? -ne 0 ]; then
    # Slackへの通知
    curl -X POST https://hooks.slack.com/services/XXX/YYY/ZZZ \
      -H 'Content-Type: application/json' \
      -d '{
        "text": "❌ Database backup failed",
        "channel": "#alerts"
      }'
    
    # メール通知
    echo "Database backup failed at $(date)" | \
      mail -s "OpsCore Backup Failure" admin@example.com
fi
```

### バックアップサイズの監視

```bash
#!/bin/bash
# monitor-backup-size.sh

BACKUP_DIR="/backup/postgres"
MAX_SIZE_GB=100  # GB

# バックアップディレクトリのサイズ取得
CURRENT_SIZE_GB=$(du -sb ${BACKUP_DIR} | awk '{print int($1/1024/1024/1024)}')

if [ ${CURRENT_SIZE_GB} -gt ${MAX_SIZE_GB} ]; then
    echo "⚠️ Backup directory size exceeds threshold: ${CURRENT_SIZE_GB}GB / ${MAX_SIZE_GB}GB"
    # アラート送信
fi
```

## ベストプラクティス

1. **定期的なリストアテスト**: 月次でリストア手順を実際に実行
2. **バックアップの暗号化**: 機密データの保護
3. **オフサイトバックアップ**: 別リージョン・別拠点への保存
4. **バックアップの検証**: 自動検証スクリプトの実行
5. **ドキュメント管理**: 復旧手順書の最新化

## 関連ドキュメント

- [デプロイガイド](../deployment/README.md)
- [監視・運用ガイド](./MONITORING.md)
- [データベース仕様（ADR 0005）](../../adr/0005-database-schema.md)
