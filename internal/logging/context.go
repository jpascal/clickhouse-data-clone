package logging

import (
	"context"
	"log/slog"
)

func SetCtx(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, "logger", logger)
}

func GetCtx(ctx context.Context) *slog.Logger {
	return ctx.Value("logger").(*slog.Logger)
}
