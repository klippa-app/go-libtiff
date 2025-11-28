package fax2tiff

import (
	"context"
	_ "embed"

	"github.com/klippa-app/go-libtiff/internal/registry"
)

//go:embed fax2tiff.wasm
var wasmBinary []byte

func init() {
	registry.Register("fax2tiff", wasmBinary)
}

func Run(ctx context.Context, args []string) error {
	return registry.Run(ctx, "fax2tiff", args...)
}
