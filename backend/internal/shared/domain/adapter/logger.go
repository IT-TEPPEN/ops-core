package adapter

import (
	"context"
	"opscore/backend/internal/shared/typ"
)

type Logger interface {
	Info(ctx context.Context, message string, args typ.Attr)
	Warn(ctx context.Context, message string, args typ.Attr)
	Error(ctx context.Context, message string, args typ.Attr)
	Debug(ctx context.Context, message string, args typ.Attr)
}
