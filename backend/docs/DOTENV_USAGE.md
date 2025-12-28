# Golang で .env ファイルを使用する方法

## 1. godotenv ライブラリのインストール

```bash
go get github.com/joho/godotenv
```

## 2. .env ファイルの作成

プロジェクトルートに `.env` ファイルを作成：

```bash
# /workspaces/backend/.env
DATABASE_URL=postgres://user:password@localhost:5432/dbname
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
GITLAB_CLIENT_ID=your_gitlab_client_id
GITLAB_CLIENT_SECRET=your_gitlab_client_secret
OAUTH_REDIRECT_URI=http://localhost:5173/oauth/callback
ENCRYPTION_KEY=your-32-character-encryption-key
LOG_LEVEL=INFO
```

## 3. main.go で .env をロード

```go
package main

import (
    "log"
    "os"
    
    "github.com/joho/godotenv"
)

func main() {
    // .envファイルを読み込む
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found")
    }
    
    // 環境変数を取得
    dbURL := os.Getenv("DATABASE_URL")
    clientID := os.Getenv("GITHUB_CLIENT_ID")
    
    log.Printf("Database URL: %s", dbURL)
}
```

## 4. カスタムパスから読み込む

```go
// 特定のファイルから読み込む
err := godotenv.Load(".env.production")

// 複数のファイルから読み込む（後のファイルが上書き）
err := godotenv.Load(".env", ".env.local")
```

## 5. 既存の環境変数を上書きしない

```go
// 既存の環境変数がある場合は上書きしない
err := godotenv.Overload(".env")  // 常に上書き
err := godotenv.Load(".env")      // 既存の環境変数を保持
```

## 6. マップとして読み込む

```go
// .envファイルをマップとして読み込む
envMap, err := godotenv.Read(".env")
if err != nil {
    log.Fatal("Error reading .env file")
}

// マップから値を取得
dbURL := envMap["DATABASE_URL"]
```

## 7. 文字列から読み込む

```go
content := `
DATABASE_URL=postgres://localhost:5432/db
API_KEY=secret123
`

envMap, err := godotenv.Unmarshal(content)
```

## 8. .gitignore に追加

```bash
# .gitignore
.env
.env.local
.env.*.local
```

機密情報は `.env` ファイルに保存し、Gitにはコミットしないようにしてください。

## 9. 環境別の設定

```go
env := os.Getenv("GO_ENV")
if env == "" {
    env = "development"
}

envFile := fmt.Sprintf(".env.%s", env)
godotenv.Load(envFile, ".env")  // フォールバック
```

## 10. 本番環境での使用

本番環境では、システムの環境変数を直接使用することを推奨します。
.envファイルは開発環境でのみ使用し、本番環境では：

- Kubernetes Secrets
- AWS Systems Manager Parameter Store
- HashiCorp Vault
- Docker secrets

などを使用してください。

## 現在のプロジェクトでの実装

`/workspaces/backend/cmd/server/main.go` で既に実装済みです：

```go
// Load .env file if it exists
if err := godotenv.Load(); err != nil {
    log.Println("No .env file found, using system environment variables")
} else {
    log.Println("Loaded .env file successfully")
}

// その後、通常通り os.Getenv() で取得
databaseUrl := os.Getenv("DATABASE_URL")
githubClientID := os.Getenv("GITHUB_CLIENT_ID")
```

これで、`.env`ファイルがあれば自動的に読み込まれ、なければシステムの環境変数を使用します。
