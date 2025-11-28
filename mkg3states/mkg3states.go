package mkg3states

import (
	"context"
	_ "embed"

	"github.com/klippa-app/go-libtiff/internal/registry"
)

//go:embed mkg3states.wasm
var wasmBinary []byte

func init() {
	registry.Register("mkg3states", wasmBinary)
}

func Run(ctx context.Context, args []string) error {
	return registry.Run(ctx, "mkg3states", args...)
}
