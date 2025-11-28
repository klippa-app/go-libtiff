package tiff2bw

import (
	"context"
	_ "embed"

	"github.com/klippa-app/go-libtiff/internal/registry"
)

//go:embed tiff2bw.wasm
var wasmBinary []byte

func init() {
	registry.Register("tiff2bw", wasmBinary)
}

func Run(ctx context.Context, args []string) error {
	return registry.Run(ctx, "tiff2bw", args...)
}
