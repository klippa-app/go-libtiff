package tiffcrop

import (
	"context"
	_ "embed"

	"github.com/klippa-app/go-libtiff/internal/registry"
)

//go:embed tiffcrop.wasm
var wasmBinary []byte

func init() {
	registry.Register("tiffcrop", wasmBinary)
}

func Run(ctx context.Context, args []string) error {
	return registry.Run(ctx, "tiffcrop", args...)
}
