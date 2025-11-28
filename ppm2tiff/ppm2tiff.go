package ppm2tiff

import (
	"context"
	_ "embed"

	"github.com/klippa-app/go-libtiff/internal/registry"
)

//go:embed ppm2tiff.wasm
var wasmBinary []byte

func init() {
	registry.Register("ppm2tiff", wasmBinary)
}

func Run(ctx context.Context, args []string) error {
	return registry.Run(ctx, "ppm2tiff", args...)
}
