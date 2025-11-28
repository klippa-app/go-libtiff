package tiff2rgba

import (
	"context"
	_ "embed"

	"github.com/klippa-app/go-libtiff/internal/registry"
)

//go:embed tiff2rgba.wasm
var wasmBinary []byte

func init() {
	registry.Register("tiff2rgba", wasmBinary)
}

func Run(ctx context.Context, args []string) error {
	return registry.Run(ctx, "tiff2rgba", args...)
}
