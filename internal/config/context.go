package config

import (
	"context"
)

func SetCtx(ctx context.Context, cfg *Config) context.Context {
	return context.WithValue(ctx, "config", cfg)
}

func GetCtx(ctx context.Context) *Config {
	return ctx.Value("config").(*Config)
}
