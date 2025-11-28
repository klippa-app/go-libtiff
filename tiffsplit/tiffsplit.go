package tiffsplit

import (
	"context"
	_ "embed"

	"github.com/klippa-app/go-libtiff/internal/registry"
)

//go:embed tiffsplit.wasm
var wasmBinary []byte

func init() {
	registry.Register("tiffsplit", wasmBinary)
}

func Run(ctx context.Context, args []string) error {
	return registry.Run(ctx, "tiffsplit", args...)
}
