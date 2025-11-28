package tiffinfo

import (
	"context"
	_ "embed"

	"github.com/klippa-app/go-libtiff/internal/registry"
)

//go:embed tiffinfo.wasm
var wasmBinary []byte

func init() {
	registry.Register("tiffinfo", wasmBinary)
}

func Run(ctx context.Context, args []string) error {
	return registry.Run(ctx, "tiffinfo", args...)
}
