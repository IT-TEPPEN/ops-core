package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	repository "opscore/backend/internal/git_repository/application/usecase"
	"opscore/backend/internal/git_repository/infrastructure/encryption"
	"opscore/backend/internal/git_repository/infrastructure/git"
	"opscore/backend/internal/git_repository/infrastructure/persistence"
	"opscore/backend/internal/git_repository/interfaces/api/handlers"
)

// Base path for cloning repositories
const baseClonePath = "./cloned_repos" // TODO: Make this configurable

// SlogLoggerAdapter は slog.Logger を handlers.Logger インターフェースに適応させる
type SlogLoggerAdapter struct {
	logger *slog.Logger
}

func (a *SlogLoggerAdapter) Info(msg string, args ...any) {
	a.logger.Info(msg, args...)
}

func (a *SlogLoggerAdapter) Error(msg string, args ...any) {
	a.logger.Error(msg, args...)
}

func (a *SlogLoggerAdapter) Debug(msg string, args ...any) {
	a.logger.Debug(msg, args...)
}

func (a *SlogLoggerAdapter) Warn(msg string, args ...any) {
	a.logger.Warn(msg, args...)
}

// provideGitManager creates a GitManager instance.
func provideGitManager() (git.GitManager, error) {
	return git.NewGithubApiManager(baseClonePath)
}

// provideAppLogger creates a structured logger based on ADR-0008.
func provideAppLogger() *slog.Logger {
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
		AddSource: true,
		Level:     logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().UTC().Format(time.RFC3339))
			}
			return a
		},
	})

	return slog.New(jsonHandler)
}

// provideHandlerLogger adapts slog.Logger to the handlers.Logger interface.
func provideHandlerLogger() handlers.Logger {
	return &SlogLoggerAdapter{logger: provideAppLogger()}
}

// provideEncryptor creates an Encryptor from the environment variable.
func provideEncryptor() (*encryption.Encryptor, error) {
	keyStr := os.Getenv("ENCRYPTION_KEY")
	if keyStr == "" {
		slog.Warn("ENCRYPTION_KEY not set, using development default key. DO NOT USE IN PRODUCTION!")
		keyStr = "dev-key-123456789012345678901234"
	}

	key := []byte(keyStr)
	if len(key) != 32 {
		return nil, encryption.ErrInvalidKey
	}

	return encryption.NewEncryptor(key)
}

// InitializeAPI initializes all dependencies for the API handlers, using Postgres.
func InitializeAPI(db *pgxpool.Pool) (*handlers.RepositoryHandler, error) {
	// Create encryptor
	encryptor, err := provideEncryptor()
	if err != nil {
		return nil, err
	}

	// Create repository (persistence layer)
	repositoryRepository := persistence.NewPostgresRepository(db, encryptor)

	// Create git manager
	gitManager, err := provideGitManager()
	if err != nil {
		return nil, err
	}

	// Create use case
	repositoryUseCase := repository.NewRepositoryUseCase(repositoryRepository, gitManager)

	// Create logger
	logger := provideHandlerLogger()

	// Create and return handler
	repositoryHandler := handlers.NewRepositoryHandler(repositoryUseCase, logger)
	return repositoryHandler, nil
}
