# ADR 0008: Backend Logging Strategy

## Status

Accepted

## Context

In application operations, logs are essential for debugging, monitoring, auditing, and performance analysis. Currently, there are no clearly defined logging rules or formats, which can lead to the following issues:

* Inconsistent log formats make parsing and searching difficult.
* Necessary information might be missing from logs, delaying root cause analysis during incidents.
* Unnecessary or sensitive information might be logged, increasing log volume or creating security risks.
* Ambiguous use of log levels makes it hard to determine the severity of log entries.

To address these challenges and improve the application's maintainability, reliability, and observability, a consistent logging strategy needs to be defined.

## Decision

We will adopt the following logging strategy for the backend application:

1. **Logging Library:** Use Go's standard `log/slog` package. This reduces external dependencies and facilitates structured logging.

2. **Log Format:** Adopt **JSON format**. This simplifies integration with log aggregation and analysis tools (e.g., Datadog, Splunk, ELK Stack) and is suitable for machine processing.

3. **Log Levels:** Define and appropriately use the following four levels. The log level should be configurable externally (e.g., via environment variables), with `INFO` and above being the default for production environments.
   * `DEBUG`: Detailed debugging information for development. Typically disabled in production.
   * `INFO`: Information indicating the normal operation of the application (e.g., request start/end, major processing steps).
   * `WARN`: Situations that require attention but are not immediate errors (e.g., potential problems, deprecated operations).
   * `ERROR`: Failures or unexpected errors during processing. Should include the error object and stack trace where possible.

4. **Structured Logging:** All logs must be structured logs with key-value attributes. This facilitates searching, filtering, and aggregation.

5. **Mandatory Log Fields:** Every log entry must include at least the following fields:
   * `time`: Timestamp of the log event (RFC3339 format, UTC).
   * `level`: Log level (e.g., "INFO", "ERROR").
   * `msg`: The main log message.
   * `source`: Source file name and line number (configured via `slog.HandlerOptions`) - primarily useful during development/debugging.

6. **Recommended Log Fields:** Include the following information when relevant to the context:
   * `request_id`: A unique identifier for the request. Generate at the start of request processing and propagate via `context.Context`. This allows easy tracing of logs related to a specific request.
   * `user_id`: The ID of the user performing the action (if authenticated).
   * `error`: For `ERROR` level logs, include details from the error object (message, type). Include stack traces if helpful.
   * Other context-specific key-value pairs (e.g., `order_id`, `resource_type`).

7. **Request-Scoped Loggers:** For operations within a specific context (like handling an HTTP request), create a logger instance with common attributes (e.g., `request_id`) and pass it via `context.Context`.

8. **Sensitive Information Masking:** Do not log sensitive information such as passwords, API keys, or personal data. If necessary, mask the data or design the logging to exclude it.

9. **Logging Locations:** Log at the following points:
   * Request start and end (including processing duration).
   * Execution points of significant business logic.
   * Interactions with external services (summary of request/response).
   * Error handling points.
   * Important initialization steps (e.g., configuration loading).

**Implementation Example (using `log/slog`):**

* **`internal/interfaces/http/middleware/request_logger.go`** (New File)

    ```go
    package middleware

    import (
        "context"
        "log/slog"
        "net/http"
        "time"

        "github.com/google/uuid"
    )

    // ContextKey type for logger and request ID
    type contextKey string
    const LoggerKey contextKey = "logger"
    const RequestIDKey contextKey = "request_id"

    // RequestLoggerMiddleware adds a request ID and a request-scoped logger to the context.
    func RequestLoggerMiddleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            requestID := uuid.NewString()
            // Create a logger instance with the request_id attribute
            logger := slog.Default().With("request_id", requestID)
            // Add the logger and request_id to the context
            ctx := context.WithValue(r.Context(), LoggerKey, logger)
            ctx = context.WithValue(ctx, RequestIDKey, requestID)

            logger.InfoContext(ctx, "Request started", slog.String("method", r.Method), slog.String("path", r.URL.Path))
            startTime := time.Now()

            // Call the next handler with the updated context
            next.ServeHTTP(w, r.WithContext(ctx))

            duration := time.Since(startTime)
            // Use slog attributes for structured data
            logger.InfoContext(ctx, "Request finished", slog.Duration("duration", duration))
        })
    }
    ```

* **Usage in Handler/Usecase (Conceptual Example)**

    ```go
    package handler // or usecase

    import (
        "context"
        "log/slog"
        "YOUR_PROJECT/internal/interfaces/http/middleware" // Import middleware package for keys
        "errors"
    )

    func handleSomething(ctx context.Context /*, ... other params */) error {
        // Retrieve the logger from context, fallback to default if not found
        logger, ok := ctx.Value(middleware.LoggerKey).(*slog.Logger)
        if !ok {
            // This should ideally not happen if middleware is applied correctly
            logger = slog.Default()
            logger.Warn("Logger not found in context, using default")
        }

        // Example: Retrieve user ID from context (assuming it's set elsewhere, e.g., by auth middleware)
        userID, _ := ctx.Value("user_id").(string) // Use appropriate key
        if userID != "" {
            logger = logger.With(slog.String("user_id", userID))
        }

        logger.InfoContext(ctx, "Processing request in handler/usecase")

        // ... business logic ...
        err := callSomeService(ctx)
        if err != nil {
            logger.ErrorContext(ctx, "Service call failed", slog.Any("error", err))
            return err // Propagate error
        }

        logger.DebugContext(ctx, "Handler/usecase processing successful")
        return nil
    }

    // Dummy function representing a call to another layer/service
    func callSomeService(ctx context.Context) error {
        logger, ok := ctx.Value(middleware.LoggerKey).(*slog.Logger)
        if !ok {
            logger = slog.Default()
        }
        logger.DebugContext(ctx, "Calling some service...")
        // Simulate an error
        if time.Now().Unix()%2 == 0 {
             return errors.New("service unavailable")
        }
        return nil
    }
    ```

* **`cmd/server/main.go` (Logger Setup and Middleware Application)**

    ```go
    package main

    import (
        "log/slog"
        "net/http"
        "os"
        "time"

        "YOUR_PROJECT/internal/interfaces/http/handler" // Assuming handlers are here
        "YOUR_PROJECT/internal/interfaces/http/middleware"
        // ... other necessary imports (DI, config, etc.)
    )

    func main() {
        // Configure slog handler (JSON format, add source, set level)
        logLevel := new(slog.LevelVar)
        envLevel := os.Getenv("LOG_LEVEL")
        switch envLevel {
        case "DEBUG":
            logLevel.Set(slog.LevelDebug)
        case "WARN":
            logLevel.Set(slog.LevelWarn)
        case "ERROR":
            logLevel.Set(slog.LevelError)
        default:
            logLevel.Set(slog.LevelInfo)
        }

        jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
            AddSource: true, // Include source file and line number
            Level:     logLevel,
            ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
                // Customize time format to RFC3339 UTC
                if a.Key == slog.TimeKey {
                    a.Value = slog.StringValue(a.Value.Time().UTC().Format(time.RFC3339))
                }
                return a
            },
        })
        slog.SetDefault(slog.New(jsonHandler))

        slog.Info("Logger initialized", slog.String("log_level", logLevel.Level().String()))

        // --- Dependency Injection (Example) ---
        // ... Initialize repositories, usecases, handlers ...
        // exampleHandler := handler.NewExampleHandler(...) // Your actual handler

        // --- HTTP Server Setup ---
        mux := http.NewServeMux()

        // Example handler registration with middleware
        // Assume exampleHandler is your actual http.Handler implementation
        // mux.Handle("/example", middleware.RequestLoggerMiddleware(exampleHandler))

        // Apply middleware globally if using a framework like Gin or Chi
        // Example with a basic mux:
        // finalHandler := middleware.RequestLoggerMiddleware(mux) // Apply middleware to the main mux

        slog.Info("Starting server on :8080")
        // if err := http.ListenAndServe(":8080", finalHandler); err != nil { // Use the handler with middleware
        //     slog.Error("Server failed to start", slog.Any("error", err))
        //     os.Exit(1)
        // }
    }
    ```

## Consequences

* **Pros:**
  * **Improved Debugging:** Structured logs and `request_id` facilitate faster troubleshooting and tracing.
  * **Enhanced Monitoring:** JSON format simplifies log aggregation, analysis, and alerting with monitoring tools.
  * **Increased Maintainability:** Consistent logging rules improve code readability and make it easier for developers to add and understand logs.
  * **Standardization:** Using the Go standard library reduces learning curves and dependencies.
* **Cons:**
  * **Initial Setup:** Requires some initial effort for `slog` configuration and middleware implementation.
  * **Log Volume:** Detailed logging might increase log storage costs (manageable via log level adjustments).
  * **Performance Overhead:** High-frequency logging might introduce minor performance overhead (usually negligible).
