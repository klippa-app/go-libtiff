package tiffdump

import (
	"context"
	_ "embed"

	"github.com/klippa-app/go-libtiff/internal/registry"
)

//go:embed tiffdump.wasm
var wasmBinary []byte

func init() {
	registry.Register("tiffdump", wasmBinary)
}

func Run(ctx context.Context, args []string) error {
	return registry.Run(ctx, "tiffdump", args...)
}
