package tiffset

import (
	"context"
	_ "embed"

	"github.com/klippa-app/go-libtiff/internal/registry"
)

//go:embed tiffset.wasm
var wasmBinary []byte

func init() {
	registry.Register("tiffset", wasmBinary)
}

func Run(ctx context.Context, args []string) error {
	return registry.Run(ctx, "tiffset", args...)
}
