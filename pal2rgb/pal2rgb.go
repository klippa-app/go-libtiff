package pal2rgb

import (
	"context"
	_ "embed"

	"github.com/klippa-app/go-libtiff/internal/registry"
)

//go:embed pal2rgb.wasm
var wasmBinary []byte

func init() {
	registry.Register("pal2rgb", wasmBinary)
}

func Run(ctx context.Context, args []string) error {
	return registry.Run(ctx, "pal2rgb", args...)
}
