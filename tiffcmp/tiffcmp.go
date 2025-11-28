package tiffcmp

import (
	"context"
	_ "embed"

	"github.com/klippa-app/go-libtiff/internal/registry"
)

//go:embed tiffcmp.wasm
var wasmBinary []byte

func init() {
	registry.Register("tiffcmp", wasmBinary)
}

func Run(ctx context.Context, args []string) error {
	return registry.Run(ctx, "tiffcmp", args...)
}
