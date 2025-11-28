package libtiff

import (
	"context"

	"github.com/tetratelabs/wazero"
)

type Config struct {
	CompilationCache wazero.CompilationCache
	FSConfig         wazero.FSConfig
	Debug            bool
}

type configCtxKey struct{}

func ConfigInContext(ctx context.Context, config *Config) context.Context {
	return context.WithValue(ctx, configCtxKey{}, config)
}

func ConfigFromContext(ctx context.Context) *Config {
	raw := ctx.Value(configCtxKey{})
	if raw == nil {
		return nil
	}

	value, ok := raw.(*Config)
	if !ok {
		return nil
	}

	return value
}
