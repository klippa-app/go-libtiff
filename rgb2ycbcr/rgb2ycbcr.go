package rgb2ycbcr

import (
	"context"
	_ "embed"

	"github.com/klippa-app/go-libtiff/internal/registry"
)

//go:embed rgb2ycbcr.wasm
var wasmBinary []byte

func init() {
	registry.Register("rgb2ycbcr", wasmBinary)
}

func Run(ctx context.Context, args []string) error {
	return registry.Run(ctx, "rgb2ycbcr", args...)
}
