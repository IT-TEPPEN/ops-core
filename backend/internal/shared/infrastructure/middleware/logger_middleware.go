package middleware

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ContextKey type for logger and request ID
type ContextKey string

const (
	LoggerKey    ContextKey = "logger"
	RequestIDKey ContextKey = "request_id"
)

// Logger provides a structured logger middleware for Gin
func Logger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// リクエスト開始時刻
		start := time.Now()

		// 一意のリクエストIDを生成
		requestID := uuid.NewString()

		// リクエストIDをヘッダーとコンテキストに設定
		c.Header("X-Request-ID", requestID)
		c.Set(string(RequestIDKey), requestID)

		// リクエスト専用のロガーをコンテキストに設定
		requestLogger := logger.With(
			slog.String("request_id", requestID),
			slog.String("ip", c.ClientIP()),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("user_agent", c.Request.UserAgent()),
		)
		c.Set(string(LoggerKey), requestLogger)

		// リクエスト開始ログ
		requestLogger.Info("Request started")

		// 次のハンドラーに処理を移譲
		c.Next()

		// レスポンス情報
		status := c.Writer.Status()
		latency := time.Since(start)

		// エラーログ取得（もしあれば）
		var logFunc func(msg string, args ...any)
		if status >= 500 {
			logFunc = requestLogger.Error
		} else if status >= 400 {
			logFunc = requestLogger.Warn
		} else {
			logFunc = requestLogger.Info
		}

		// センシティブデータをマスクする（URLクエリパラメータからのパスワードなど）
		// セキュリティ強化: クエリパラメータからセンシティブデータを削除したURLを使用
		sanitizedPath := sanitizeURL(c.Request.URL.String())

		// リクエスト終了ログ
		logFunc("Request completed",
			slog.Int("status", status),
			slog.String("latency", latency.String()),
			slog.String("sanitized_path", sanitizedPath),
			slog.Int("response_size", c.Writer.Size()),
		)
	}
}

// GetLogger retrieves the logger from gin context
func GetLogger(c *gin.Context) *slog.Logger {
	if logger, exists := c.Get(string(LoggerKey)); exists {
		if l, ok := logger.(*slog.Logger); ok {
			return l
		}
	}
	// Fallback to default logger
	return slog.Default()
}

// GetRequestID retrieves the request ID from gin context
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(string(RequestIDKey)); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// ContextWithLogger adds logger to context
func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

// LoggerFromContext retrieves logger from context
func LoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(LoggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// センシティブデータをマスクする関数
func sanitizeURL(urlString string) string {
	u, err := url.Parse(urlString)
	if err != nil {
		// URLのパース失敗時は元の文字列を返す
		return urlString
	}

	// センシティブと考えられるクエリパラメータのリスト
	sensitiveParams := []string{
		"password", "token", "key", "secret", "auth", "credential", "api_key",
		"apikey", "access_token", "refresh_token", "private", "pwd", "passwd",
	}

	if u.RawQuery != "" {
		q := u.Query()
		for _, param := range sensitiveParams {
			if q.Has(param) {
				q.Set(param, "[REDACTED]")
			}
		}
		u.RawQuery = q.Encode()
	}

	return u.String()
}
