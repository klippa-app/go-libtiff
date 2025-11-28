package tiff2ps

import (
	"context"
	_ "embed"

	"github.com/klippa-app/go-libtiff/internal/registry"
)

//go:embed tiff2ps.wasm
var wasmBinary []byte

func init() {
	registry.Register("tiff2ps", wasmBinary)
}

func Run(ctx context.Context, args []string) error {
	return registry.Run(ctx, "tiff2ps", args...)
}
