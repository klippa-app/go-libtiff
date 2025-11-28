package tiffdither

import (
	"context"
	_ "embed"

	"github.com/klippa-app/go-libtiff/internal/registry"
)

//go:embed tiffdither.wasm
var wasmBinary []byte

func init() {
	registry.Register("tiffdither", wasmBinary)
}

func Run(ctx context.Context, args []string) error {
	return registry.Run(ctx, "tiffdither", args...)
}
