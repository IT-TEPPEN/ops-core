//go:build wireinject
// +build wireinject

package main

import (
	"log/slog"
	"os"
	"time"

	"opscore/backend/infrastructure/git"
	"opscore/backend/infrastructure/persistence"
	"opscore/backend/interfaces/api/handlers"
	repoUseCase "opscore/backend/usecases/repository"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
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

// provideGitManager is a Wire provider function for GitManager.
func provideGitManager() (git.GitManager, error) {
	// CLIベースのGitManagerからGitHub APIベースの実装に変更
	return git.NewGithubApiManager(baseClonePath)
}

// provideAppLogger is a Wire provider function for *slog.Logger
func provideAppLogger() *slog.Logger {
	// ADR-0008に従ったロガーの設定
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

	return slog.New(jsonHandler)
}

// provideHandlerLogger adapts slog.Logger to the handlers.Logger interface
func provideHandlerLogger() handlers.Logger {
	return &SlogLoggerAdapter{logger: provideAppLogger()}
}

// InitializeAPI initializes all dependencies for the API handlers, using Postgres.
func InitializeAPI(db *pgxpool.Pool) (*handlers.RepositoryHandler, error) {
	wire.Build(
		provideGitManager,
		provideHandlerLogger,
		persistence.NewPostgresRepository,
		repoUseCase.NewRepositoryUseCase,
		handlers.NewRepositoryHandler,
	)
	return nil, nil
}
