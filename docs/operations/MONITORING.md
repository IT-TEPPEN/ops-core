# Monitoring and Operations Guide

このドキュメントでは、OpsCoreの監視と運用に関するガイドラインを説明します。

## 目次

- [監視戦略](#監視戦略)
- [メトリクス監視](#メトリクス監視)
- [ログ管理](#ログ管理)
- [ヘルスチェック](#ヘルスチェック)
- [アラート設定](#アラート設定)
- [パフォーマンスモニタリング](#パフォーマンスモニタリング)
- [トラブルシューティング](#トラブルシューティング)

## 監視戦略

OpsCoreの監視は、以下の階層で実施します：

### 監視レベル

1. **インフラストラクチャ監視**: サーバー、ネットワーク、ストレージ
2. **アプリケーション監視**: API応答時間、エラー率
3. **ビジネスメトリクス監視**: ユーザー数、ドキュメント数、作業証跡数

### 監視ツール（推奨）

- **Prometheus**: メトリクス収集
- **Grafana**: メトリクス可視化
- **Loki**: ログ集約
- **Alertmanager**: アラート管理

## メトリクス監視

### 収集するメトリクス

#### システムメトリクス

```
# CPUメトリクス
node_cpu_seconds_total
process_cpu_seconds_total

# メモリメトリクス
node_memory_MemTotal_bytes
node_memory_MemAvailable_bytes
process_resident_memory_bytes

# ディスクメトリクス
node_disk_read_bytes_total
node_disk_written_bytes_total
node_filesystem_avail_bytes
```

#### アプリケーションメトリクス

```
# HTTPリクエストメトリクス
http_requests_total{method="GET", endpoint="/api/v1/documents", status="200"}
http_request_duration_seconds{method="GET", endpoint="/api/v1/documents"}

# データベースメトリクス
db_connections_total
db_connections_active
db_query_duration_seconds

# ビジネスメトリクス
documents_total
execution_records_total
users_active_total
storage_usage_bytes
```

### Prometheusの設定

#### prometheus.yml

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'opscore-backend'
    static_configs:
      - targets: ['backend:8080']
    metrics_path: '/metrics'

  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres-exporter:9187']

  - job_name: 'node'
    static_configs:
      - targets: ['node-exporter:9100']
```

#### バックエンドへのメトリクスエンドポイント追加

```go
// cmd/server/main.go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func setupMetrics(e *echo.Echo) {
    // Prometheusメトリクスエンドポイント
    e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
    
    // カスタムメトリクスの登録
    registerCustomMetrics()
}

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
        },
        []string{"method", "endpoint"},
    )
)
```

## ログ管理

### ログレベル

OpsCoreでは、以下のログレベルを使用します：

- **DEBUG**: デバッグ情報（開発環境のみ）
- **INFO**: 一般的な情報
- **WARN**: 警告（エラーではないが注意が必要）
- **ERROR**: エラー（処理は継続）
- **FATAL**: 致命的エラー（アプリケーション停止）

### ログフォーマット

```json
{
  "level": "info",
  "timestamp": "2024-01-15T10:30:00.000Z",
  "logger": "document.usecase",
  "message": "Document created successfully",
  "document_id": "doc-123",
  "user_id": "user-456",
  "request_id": "req-xyz789"
}
```

### ログの出力先

#### 開発環境
- 標準出力（stdout）

#### 本番環境
- ファイル: `/var/log/opscore/app.log`
- ログローテーション: 日次、7日間保持
- 集約: Loki、Elasticsearch等

### ログ設定

```go
// internal/shared/infrastructure/logger/logger.go
import "go.uber.org/zap"

func NewLogger(env string) (*zap.Logger, error) {
    var config zap.Config
    
    if env == "production" {
        config = zap.NewProductionConfig()
        config.OutputPaths = []string{
            "stdout",
            "/var/log/opscore/app.log",
        }
    } else {
        config = zap.NewDevelopmentConfig()
    }
    
    return config.Build()
}
```

### Lokiの設定

```yaml
# docker-compose.yml
services:
  loki:
    image: grafana/loki:latest
    ports:
      - "3100:3100"
    volumes:
      - ./loki-config.yaml:/etc/loki/loki-config.yaml
    command: -config.file=/etc/loki/loki-config.yaml

  promtail:
    image: grafana/promtail:latest
    volumes:
      - /var/log:/var/log
      - ./promtail-config.yaml:/etc/promtail/promtail-config.yaml
    command: -config.file=/etc/promtail/promtail-config.yaml
```

## ヘルスチェック

### ヘルスチェックエンドポイント

```
GET /health
```

レスポンス例：

```json
{
  "status": "ok",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0",
  "checks": {
    "database": {
      "status": "ok",
      "response_time_ms": 5
    },
    "storage": {
      "status": "ok",
      "type": "s3"
    }
  }
}
```

### ヘルスチェックの実装

```go
// internal/shared/interfaces/api/handlers/health_handler.go
type HealthHandler struct {
    db      *sql.DB
    storage storage.Storage
}

func (h *HealthHandler) GetHealth(c echo.Context) error {
    result := HealthResponse{
        Status:    "ok",
        Timestamp: time.Now(),
        Version:   version.Version,
        Checks:    make(map[string]CheckResult),
    }
    
    // データベースチェック
    if err := h.checkDatabase(); err != nil {
        result.Status = "degraded"
        result.Checks["database"] = CheckResult{
            Status: "error",
            Error:  err.Error(),
        }
    } else {
        result.Checks["database"] = CheckResult{Status: "ok"}
    }
    
    // ストレージチェック
    if err := h.checkStorage(); err != nil {
        result.Status = "degraded"
        result.Checks["storage"] = CheckResult{
            Status: "error",
            Error:  err.Error(),
        }
    } else {
        result.Checks["storage"] = CheckResult{Status: "ok"}
    }
    
    statusCode := http.StatusOK
    if result.Status != "ok" {
        statusCode = http.StatusServiceUnavailable
    }
    
    return c.JSON(statusCode, result)
}
```

### Readiness / Liveness Probe（Kubernetes）

```yaml
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: backend
    image: opscore-backend:latest
    livenessProbe:
      httpGet:
        path: /health
        port: 8080
      initialDelaySeconds: 30
      periodSeconds: 10
    readinessProbe:
      httpGet:
        path: /ready
        port: 8080
      initialDelaySeconds: 5
      periodSeconds: 5
```

## アラート設定

### アラートルール

#### prometheus-alerts.yml

```yaml
groups:
  - name: opscore_alerts
    interval: 30s
    rules:
      # APIエラー率が高い
      - alert: HighErrorRate
        expr: |
          rate(http_requests_total{status=~"5.."}[5m]) / 
          rate(http_requests_total[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value }}%"

      # データベース接続エラー
      - alert: DatabaseDown
        expr: up{job="postgres"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Database is down"
          description: "PostgreSQL is not responding"

      # メモリ使用率が高い
      - alert: HighMemoryUsage
        expr: |
          (node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / 
          node_memory_MemTotal_bytes > 0.9
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High memory usage"
          description: "Memory usage is {{ $value }}%"

      # ディスク使用率が高い
      - alert: HighDiskUsage
        expr: |
          (node_filesystem_size_bytes - node_filesystem_avail_bytes) / 
          node_filesystem_size_bytes > 0.85
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High disk usage"
          description: "Disk usage is {{ $value }}%"

      # API応答時間が遅い
      - alert: SlowAPIResponse
        expr: |
          histogram_quantile(0.95, 
            rate(http_request_duration_seconds_bucket[5m])
          ) > 2
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "Slow API response time"
          description: "95th percentile response time is {{ $value }}s"
```

### Alertmanagerの設定

```yaml
# alertmanager.yml
global:
  resolve_timeout: 5m

route:
  group_by: ['alertname', 'severity']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 12h
  receiver: 'team-email'

receivers:
  - name: 'team-email'
    email_configs:
      - to: 'team@example.com'
        from: 'alertmanager@example.com'
        smarthost: 'smtp.example.com:587'
        auth_username: 'alertmanager@example.com'
        auth_password: 'password'

  - name: 'slack'
    slack_configs:
      - api_url: 'https://hooks.slack.com/services/XXX/YYY/ZZZ'
        channel: '#alerts'
        title: '{{ .GroupLabels.alertname }}'
        text: '{{ range .Alerts }}{{ .Annotations.description }}{{ end }}'
```

## パフォーマンスモニタリング

### キーパフォーマンス指標（KPI）

#### 1. レスポンスタイム

- 目標: 95%のリクエストが2秒以内
- 測定: `http_request_duration_seconds`

#### 2. スループット

- 目標: 100 req/s以上
- 測定: `rate(http_requests_total[1m])`

#### 3. エラー率

- 目標: 1%未満
- 測定: `rate(http_requests_total{status=~"5.."}[5m])`

#### 4. データベースクエリ時間

- 目標: 95%のクエリが100ms以内
- 測定: `db_query_duration_seconds`

### Grafanaダッシュボード

#### ダッシュボードの作成

1. **システム概要**
   - CPU使用率
   - メモリ使用率
   - ディスク使用率
   - ネットワークトラフィック

2. **アプリケーション概要**
   - リクエスト数（時系列）
   - エラー率（時系列）
   - レスポンスタイム（パーセンタイル）
   - アクティブコネクション数

3. **ビジネスメトリクス**
   - ドキュメント数
   - 作業証跡数
   - アクティブユーザー数
   - ストレージ使用量

#### Grafanaダッシュボード設定例

```json
{
  "dashboard": {
    "title": "OpsCore Overview",
    "panels": [
      {
        "title": "Request Rate",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])"
          }
        ]
      },
      {
        "title": "Error Rate",
        "targets": [
          {
            "expr": "rate(http_requests_total{status=~\"5..\"}[5m])"
          }
        ]
      }
    ]
  }
}
```

## トラブルシューティング

### よくある問題と対処法

#### 1. APIレスポンスが遅い

**診断:**
```bash
# レスポンスタイムの確認
curl -w "\nTime: %{time_total}s\n" http://localhost:8080/api/v1/documents

# データベースクエリの確認
docker compose exec postgres psql -U opscore -d opscore -c "SELECT * FROM pg_stat_activity WHERE state = 'active';"
```

**対処:**
- インデックスの最適化
- クエリの最適化
- キャッシュの導入

#### 2. メモリリーク

**診断:**
```bash
# メモリ使用量の確認
docker stats

# Goのプロファイリング
curl http://localhost:8080/debug/pprof/heap > heap.pprof
go tool pprof heap.pprof
```

**対処:**
- プロファイリング結果の分析
- リソースリークの修正
- メモリ制限の調整

#### 3. データベース接続エラー

**診断:**
```bash
# PostgreSQLの状態確認
docker compose exec postgres pg_isready

# 接続数の確認
docker compose exec postgres psql -U opscore -d opscore -c "SELECT count(*) FROM pg_stat_activity;"
```

**対処:**
- 接続プールの設定見直し
- `max_connections`の増加
- 長時間実行クエリのキル

#### 4. ストレージ容量不足

**診断:**
```bash
# ディスク使用量の確認
df -h

# S3使用量の確認
aws s3 ls s3://opscore-attachments --recursive --summarize
```

**対処:**
- 古い添付ファイルのアーカイブ
- ストレージ容量の拡張
- ライフサイクルポリシーの設定

## 運用ベストプラクティス

### 1. 定期メンテナンス

- **日次**: ログのローテーション、バックアップ確認
- **週次**: データベースのVACUUM、統計情報更新
- **月次**: セキュリティアップデート、パフォーマンスレビュー

### 2. キャパシティプランニング

- ストレージ使用量の予測
- ユーザー増加に伴うスケーリング計画
- データベースパフォーマンスの監視

### 3. インシデント対応

1. **検知**: アラートによる異常検知
2. **初動**: 影響範囲の確認、ステークホルダーへの通知
3. **対応**: 原因調査と修復
4. **事後**: ポストモーテムの実施、再発防止策の策定

## 関連ドキュメント

- [デプロイガイド](../deployment/README.md)
- [バックアップ手順](./BACKUP.md)
- [ロギング戦略（ADR 0008）](../../adr/0008-backend-logging-strategy.md)
