package tiffmedian

import (
	"context"
	_ "embed"

	"github.com/klippa-app/go-libtiff/internal/registry"
)

//go:embed tiffmedian.wasm
var wasmBinary []byte

func init() {
	registry.Register("tiffmedian", wasmBinary)
}

func Run(ctx context.Context, args []string) error {
	return registry.Run(ctx, "tiffmedian", args...)
}
