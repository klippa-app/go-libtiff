package raw2tiff

import (
	"context"
	_ "embed"

	"github.com/klippa-app/go-libtiff/internal/registry"
)

//go:embed raw2tiff.wasm
var wasmBinary []byte

func init() {
	registry.Register("raw2tiff", wasmBinary)
}

func Run(ctx context.Context, args []string) error {
	return registry.Run(ctx, "raw2tiff", args...)
}
