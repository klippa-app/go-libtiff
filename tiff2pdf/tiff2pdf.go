package tiff2pdf

import (
	"context"
	_ "embed"

	"github.com/klippa-app/go-libtiff/internal/registry"
)

//go:embed tiff2pdf.wasm
var wasmBinary []byte

func init() {
	registry.Register("tiff2pdf", wasmBinary)
}

func Run(ctx context.Context, args []string) error {
	return registry.Run(ctx, "tiff2pdf", args...)
}
